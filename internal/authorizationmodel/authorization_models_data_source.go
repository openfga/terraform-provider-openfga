package authorizationmodel

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
var _ datasource.DataSource = &AuthorizationModelsDataSource{}
var _ datasource.DataSourceWithConfigure = &AuthorizationModelsDataSource{}

func NewAuthorizationModelsDataSource() datasource.DataSource {
	return &AuthorizationModelsDataSource{}
}

type AuthorizationModelsDataSource struct {
	client *AuthorizationModelClient
}

type AuthorizationModelsDataSourceModel struct {
	StoreId             types.String              `tfsdk:"store_id"`
	AuthorizationModels []AuthorizationModelModel `tfsdk:"authorization_models"`
}

func (d *AuthorizationModelsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_models"
}

func (d *AuthorizationModelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides the ability to list and retrieve details of existing authorization models in a specific store.",

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the store to list authorization models for.",
				Required:            true,
			},
			"authorization_models": schema.ListNestedAttribute{
				MarkdownDescription: "List of existing authorization models in the specific store.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the authorization model.",
							Computed:            true,
						},
						"model_json": schema.StringAttribute{
							MarkdownDescription: "The authorization model definition in JSON format.",
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
						},
					},
				},
			},
		},
	}
}

func (d *AuthorizationModelsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = NewAuthorizationModelClient(client)
}

func (d *AuthorizationModelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AuthorizationModelsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	authorizationModelModels, err := d.client.ListAuthorizationModels(ctx, state.StoreId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read authorization models, got error: %s", err))
		return
	}

	state.AuthorizationModels = *authorizationModelModels

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
