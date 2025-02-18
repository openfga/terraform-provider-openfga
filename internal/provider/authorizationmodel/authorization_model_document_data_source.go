package authorizationmodel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfga "github.com/openfga/go-sdk"
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
	Dsl   types.String              `tfsdk:"dsl"`
	Json  types.String              `tfsdk:"json"`
	Model *CustomAuthorizationModel `tfsdk:"model"`

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
				MarkdownDescription: "An authorization model in DSL format. Conflicts with `json` and `model` fields.",
				Optional:            true,
			},
			"json": schema.StringAttribute{
				MarkdownDescription: "An authorization model in JSON format. Conflicts with `dsl` and `model` fields.",
				Optional:            true,
			},
			"model": schema.SingleNestedAttribute{
				MarkdownDescription: "An authorization model as Terraform object. Conflicts with `dsl` and `json` fields.",
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

	model := state.Model
	jsonString := state.Json.ValueStringPointer()
	dslString := state.Dsl.ValueStringPointer()

	if model != nil {
		result, err := json.Marshal(state.Model)
		if err != nil {
			resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to transform model into JSON, got error: %s", err))
			return
		}

		jsonString = openfga.PtrString(string(result))
	}

	if jsonString != nil {
		result, err := transformer.TransformJSONStringToDSL(*jsonString)
		if err != nil {
			resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to transform JSON into DSL, got error: %s", err))
			return
		}

		dslString = result
	}

	if dslString == nil {
		resp.Diagnostics.AddError("Input Error", "DSL is undefined")
		return
	}

	// Transform DSL to canonical JSON form
	unstableResult, err := transformer.TransformDSLToJSON(*dslString)
	if err != nil {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to transform DSL into JSON, got error: %s", err))
		return
	}

	var tmp any
	err = json.Unmarshal([]byte(unstableResult), &tmp)
	if err != nil {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to bring JSON in canonical form, got error: %s", err))
		return
	}
	stableResultBytes, err := json.Marshal(tmp)
	if err != nil {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to bring JSON in canonical form, got error: %s", err))
		return
	}

	stableResult := string(stableResultBytes)

	state.Result = types.StringValue(stableResult)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
