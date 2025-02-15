package query

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CheckQueryDataSource{}
var _ datasource.DataSourceWithConfigure = &CheckQueryDataSource{}

func NewCheckQueryDataSource() datasource.DataSource {
	return &CheckQueryDataSource{}
}

type CheckQueryDataSource struct {
	client *QueryClient
}

type CheckQueryDataSourceModel struct {
	StoreId              types.String `tfsdk:"store_id"`
	AuthorizationModelId types.String `tfsdk:"authorization_model_id"`

	CheckQueryModel

	Result types.Bool `tfsdk:"result"`
}

func (d *CheckQueryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_query"
}

func (d *CheckQueryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A 'check' query can be performed to establish whether a particular user has a specific relationship with a particular object.",

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store this query is run against",
				Required:            true,
			},
			"authorization_model_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA authorization model this query is run against",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The user of the query",
				Required:            true,
			},
			"relation": schema.StringAttribute{
				MarkdownDescription: "The relation to check for",
				Required:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object of the query",
				Required:            true,
			},
			"contextual_tuples": schema.ListNestedAttribute{
				MarkdownDescription: "The contextual tuples that should be considered for the query",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user": schema.StringAttribute{
							MarkdownDescription: "The user of the contextual relationship tuple",
							Required:            true,
						},
						"relation": schema.StringAttribute{
							MarkdownDescription: "The relation of the contextual relationship tuple",
							Required:            true,
						},
						"object": schema.StringAttribute{
							MarkdownDescription: "The object of the contextual relationship tuple",
							Required:            true,
						},
						"condition": schema.SingleNestedAttribute{
							MarkdownDescription: "A condition of the contextual relationship tuple",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "The name of the condition",
									Required:            true,
								},
								"context_json": schema.StringAttribute{
									MarkdownDescription: "The (partial) context under which the condition is evaluated",
									CustomType:          jsontypes.NormalizedType{},
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"context_json": schema.StringAttribute{
				MarkdownDescription: "The (partial) context under which the condition is evaluated",
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
			},
			"result": schema.BoolAttribute{
				MarkdownDescription: "Boolean value indicating whether the user has a relation to the object",
				Computed:            true,
			},
		},
	}
}

func (d *CheckQueryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = NewQueryClient(client)
}

func (d *CheckQueryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state CheckQueryDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Check(ctx, state.StoreId.ValueString(), state.AuthorizationModelId.ValueString(), state.CheckQueryModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to perform check query, got error: %s", err))
		return
	}

	state.Result = result

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
