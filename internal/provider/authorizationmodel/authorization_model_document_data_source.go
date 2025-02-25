package authorizationmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/encoding/protojson"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/language/pkg/go/transformer"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AuthorizationModelDocumentDataSource{}
var _ datasource.DataSourceWithConfigure = &AuthorizationModelDocumentDataSource{}

func NewAuthorizationModelDocumentDataSource() datasource.DataSource {
	return &AuthorizationModelDocumentDataSource{}
}

type AuthorizationModelDocumentDataSource struct{}

type AuthorizationModelDocumentDataSourceModel struct {
	Dsl         types.String              `tfsdk:"dsl"`
	ModFilePath types.String              `tfsdk:"mod_file_path"`
	Json        types.String              `tfsdk:"json"`
	Model       *CustomAuthorizationModel `tfsdk:"model"`

	Result types.String `tfsdk:"result"`
}

func (d *AuthorizationModelDocumentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_model_document"
}

func (d *AuthorizationModelDocumentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Generates an authorization model in JSON format for use with resources that expect authorization models such as ` + "`openfga_authorization_model`" + `.

Can be used to convert an authorization model from DSL format to JSON format. It is also possible to provide an authorization model in JSON format or as native Terraform object.

It will always generate a stable output that is not influenced by the format of the input data.

Using this data source to generate authorization models is optional. It is also valid to use literal JSON strings in your configuration.
`,

		Attributes: map[string]schema.Attribute{
			"dsl": schema.StringAttribute{
				MarkdownDescription: "An authorization model in DSL format. Conflicts with `json`, `model` and `mod_file_path` fields.",
				Optional:            true,
			},
			"mod_file_path": schema.StringAttribute{
				MarkdownDescription: "A file path to an `fga.mod` file. Conflicts with `json`, `model` and `dsl` fields.",
				Optional:            true,
			},
			"json": schema.StringAttribute{
				MarkdownDescription: "An authorization model in JSON format. Conflicts with `dsl`, `model` and `mod_file_path` fields.",
				Optional:            true,
			},
			"model": schema.SingleNestedAttribute{
				MarkdownDescription: "An authorization model as Terraform object. Conflicts with `dsl`, `json` and `mod_file_path` fields.",
				Optional:            true,
				Attributes:          CustomAuthorizationModelSchema(),
			},
			"result": schema.StringAttribute{
				MarkdownDescription: "The authorization model definition in a stable JSON format.",
				Computed:            true,
			},
		},
	}
}

func (p AuthorizationModelDocumentDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("dsl"),
			path.MatchRoot("json"),
			path.MatchRoot("model"),
			path.MatchRoot("mod_file_path"),
		),
	}
}

func (d *AuthorizationModelDocumentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *AuthorizationModelDocumentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AuthorizationModelDocumentDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	modelProto, err := extractAuthorizationModelProto(&state)
	if err != nil {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to extract model from state, got error: %s", err))
		return
	}

	json, err := marshalToSanitizedJson(modelProto)
	if err != nil {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to marshal model into sanitized JSON, got error: %s", err))
		return
	}

	state.Result = types.StringValue(json)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func extractAuthorizationModelProto(state *AuthorizationModelDocumentDataSourceModel) (*openfgav1.AuthorizationModel, error) {
	if state.Model != nil {
		return parseModelToAuthorizationModelProto(state.Model)
	}

	if state.ModFilePath.ValueString() != "" {
		return parseModFileToAuthorizationModelProto(state.ModFilePath.ValueString())
	}

	if state.Dsl.ValueString() != "" {
		return parseDslToAuthorizationModelProto(state.Dsl.ValueString())
	}

	if state.Json.ValueString() != "" {
		return parseJsonToAuthorizationModelProto(state.Json.ValueString())
	}

	return nil, fmt.Errorf("at least one of model, mod file path, DSL or JSON has to be provided")
}

func parseModelToAuthorizationModelProto(model *CustomAuthorizationModel) (*openfgav1.AuthorizationModel, error) {
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		return nil, fmt.Errorf("unable to transform custom model into JSON, got error: %s", err)
	}

	return parseJsonToAuthorizationModelProto(string(jsonBytes))
}

func parseDslToAuthorizationModelProto(dsl string) (*openfgav1.AuthorizationModel, error) {
	modelProto, err := transformer.TransformDSLToProto(dsl)
	if err != nil {
		return nil, fmt.Errorf("unable to transform DSL into model proto, got error: %s", err)
	}

	return modelProto, nil
}

func parseModFileToAuthorizationModelProto(modFilePath string) (*openfgav1.AuthorizationModel, error) {
	modFileBytes, err := os.ReadFile(modFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read mod file, got error: %s", err)
	}

	modFile, err := transformer.TransformModFile(string(modFileBytes))
	if err != nil {
		return nil, fmt.Errorf("unable to parse mod file, got error: %s", err)
	}

	modFileDirectory := filepath.Dir(modFilePath)

	moduleFiles := []transformer.ModuleFile{}
	for _, moduleFilePathProperty := range modFile.Contents.Value {
		moduleFilePath := filepath.Join(modFileDirectory, moduleFilePathProperty.Value)

		moduleFileBytes, err := os.ReadFile(moduleFilePath)
		if err != nil {
			return nil, fmt.Errorf("unable to read module file, got error: %s", err)
		}

		moduleFile := transformer.ModuleFile{
			Name:     moduleFilePath,
			Contents: string(moduleFileBytes),
		}

		moduleFiles = append(moduleFiles, moduleFile)
	}

	modelProto, err := transformer.TransformModuleFilesToModel(moduleFiles, modFile.Schema.Value)
	if err != nil {
		return nil, fmt.Errorf("unable to transform module files into model proto, got error: %s", err)
	}

	return modelProto, nil
}

func parseJsonToAuthorizationModelProto(json string) (*openfgav1.AuthorizationModel, error) {
	modelProto, err := transformer.LoadJSONStringToProto(json)
	if err != nil {
		return nil, fmt.Errorf("unable to transform JSON into model proto, got error: %s", err)
	}

	return modelProto, nil
}

func marshalToSanitizedJson(modelProto *openfgav1.AuthorizationModel) (string, error) {
	marshaller := protojson.MarshalOptions{EmitDefaultValues: true}

	unstableJsonBytes, err := marshaller.Marshal(modelProto)
	if err != nil {
		return "", fmt.Errorf("unable to marshal model proto into unstable JSON, got error: %s", err)
	}

	var sanitizedAuthorizationModel AuthorizationModelWithoutId
	err = json.Unmarshal(unstableJsonBytes, &sanitizedAuthorizationModel)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal unstable JSON into sanitized model, got error: %s", err)
	}

	stableJsonBytes, err := json.Marshal(sanitizedAuthorizationModel)
	if err != nil {
		return "", fmt.Errorf("unable to convert sanitized model into canonical JSON form, got error: %s", err)
	}

	return string(stableJsonBytes), nil
}
