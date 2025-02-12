package authorizationmodel

import (
	"context"
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

	jsonString := state.Json.ValueStringPointer()
	dslString := state.Dsl.ValueStringPointer()

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
	result, err := transformer.TransformDSLToJSON(*dslString)
	if err != nil {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable transform DSL into JSON, got error: %s", err))
		return
	}

	state.Result = types.StringValue(result)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
