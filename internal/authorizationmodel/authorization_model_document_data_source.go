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

	"github.com/openfga/go-sdk/client"
	"github.com/openfga/language/pkg/go/transformer"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AuthorizationModelDocumentDataSource{}
var _ datasource.DataSourceWithConfigure = &AuthorizationModelDocumentDataSource{}

func NewAuthorizationModelDocumentDataSource() datasource.DataSource {
	return &AuthorizationModelDocumentDataSource{}
}

type AuthorizationModelDocumentDataSource struct {
	client *client.OpenFgaClient
}

type AuthorizationModelDocumentDataSourceModel struct {
	Dsl  types.String `tfsdk:"dsl"`
	Json types.String `tfsdk:"json"`

	Result types.String `tfsdk:"result"`
}

func (d *AuthorizationModelDocumentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_model_document"
}

func (d *AuthorizationModelDocumentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "An authorization model document represents the configuration of an authorization model. It can be used to convert between DSL and JSON or to represent a JSON model in canonical form.",

		Attributes: map[string]schema.Attribute{
			"dsl": schema.StringAttribute{
				MarkdownDescription: "The authorization model in DSL format",
				Optional:            true,
			},
			"json": schema.StringAttribute{
				MarkdownDescription: "The authorization model in JSON format",
				Optional:            true,
			},
			"result": schema.StringAttribute{
				MarkdownDescription: "The resulting model in JSON format",
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
		),
	}
}

func (d *AuthorizationModelDocumentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.OpenFgaClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.OpenFgaClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *AuthorizationModelDocumentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AuthorizationModelDocumentDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var jsonData string
	if !state.Dsl.IsNull() {
		var err error
		jsonData, err = transformer.TransformDSLToJSON(state.Dsl.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable transform DSL into JSON, got error: %s", err))
			return
		}
	} else if !state.Json.IsNull() {
		jsonData = state.Json.ValueString()
	}

	// JSON validation
	authorizationModel, err := transformer.LoadJSONStringToProto(jsonData)
	if err != nil {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to parse JSON, got error: %s", err))
		return
	}

	// Convert into canonical JSON form
	jsonBytes, err := json.Marshal(authorizationModel)
	if err != nil {
		resp.Diagnostics.AddError("InternalError", fmt.Sprintf("Unable to convert model to JSON, got error: %s", err))
		return
	}

	state.Result = types.StringValue(string(jsonBytes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
