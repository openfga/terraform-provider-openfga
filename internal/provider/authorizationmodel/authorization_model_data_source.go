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
var _ datasource.DataSource = &AuthorizationModelDataSource{}
var _ datasource.DataSourceWithConfigure = &AuthorizationModelDataSource{}

func NewAuthorizationModelDataSource() datasource.DataSource {
	return &AuthorizationModelDataSource{}
}

type AuthorizationModelDataSource struct {
	client *AuthorizationModelClient
}

type AuthorizationModelDataSourceModel struct {
	StoreId types.String `tfsdk:"store_id"`
	AuthorizationModelModel
}

func (d *AuthorizationModelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_model"
}

func (d *AuthorizationModelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides the ability to retrieve details of an existing OpenFGA authorization model.",

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the store this authorization model belongs to.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the authorization model. Can be left blank to retrieve the latest authorization model.",
				Optional:            true,
			},
			"model_json": schema.StringAttribute{
				MarkdownDescription: "The authorization model definition in JSON format.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
		},
	}
}

func (d *AuthorizationModelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AuthorizationModelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AuthorizationModelDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		authorizationModelModel *AuthorizationModelModel
		err                     error
	)
	if state.Id.IsNull() {
		authorizationModelModel, err = d.client.ReadLatestAuthorizationModel(ctx, state.StoreId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read latest authorization model, got error: %s", err))
			return
		}
	} else {
		authorizationModelModel, err = d.client.ReadAuthorizationModel(ctx, state.StoreId.ValueString(), state.AuthorizationModelModel)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read authorization model, got error: %s", err))
			return
		}
	}

	state.AuthorizationModelModel = *authorizationModelModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
