package authorizationmodel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AuthorizationModelDataSource{}
var _ datasource.DataSourceWithConfigure = &AuthorizationModelDataSource{}

func NewAuthorizationModelDataSource() datasource.DataSource {
	return &AuthorizationModelDataSource{}
}

type AuthorizationModelDataSource struct {
	client *client.OpenFgaClient
}

type AuthorizationModelDataSourceModel AuthorizationModelModel

func (d *AuthorizationModelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_model"
}

func (d *AuthorizationModelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "An authorization model combines one or more type definitions. This is used to define the permission model of a system.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA authorization model (leave blank to get latest model)",
				Optional:            true,
			},
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store this authorization model belongs to",
				Required:            true,
			},
			"model_json": schema.StringAttribute{
				MarkdownDescription: "The full authorization model definition in JSON format",
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

	d.client = client
}

func (d *AuthorizationModelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AuthorizationModelDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var authorizationModel *openfga.AuthorizationModel
	if state.Id.IsNull() {
		options := client.ClientReadLatestAuthorizationModelOptions{
			StoreId: state.StoreId.ValueStringPointer(),
		}

		response, err := d.client.ReadLatestAuthorizationModel(ctx).Options(options).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read latest authorization model, got error: %s", err))
			return
		}

		authorizationModel = response.AuthorizationModel
	} else {
		options := client.ClientReadAuthorizationModelOptions{
			StoreId:              state.StoreId.ValueStringPointer(),
			AuthorizationModelId: state.Id.ValueStringPointer(),
		}

		response, err := d.client.ReadAuthorizationModel(ctx).Options(options).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read authorization model, got error: %s", err))
			return
		}

		authorizationModel = response.AuthorizationModel
	}

	state.Id = types.StringValue(authorizationModel.Id)

	authorizationModelWithoutId := AuthorizationModelWithoutId{
		SchemaVersion:   authorizationModel.SchemaVersion,
		TypeDefinitions: authorizationModel.TypeDefinitions,
		Conditions:      authorizationModel.Conditions,
	}

	jsonBytes, err := json.Marshal(authorizationModelWithoutId)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Response Data", fmt.Sprintf("Unable to convert to model JSON, got error: %s", err))
		return
	}
	state.ModelJson = jsontypes.NewNormalizedValue(string(jsonBytes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
