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
var _ datasource.DataSource = &AuthorizationModelsDataSource{}
var _ datasource.DataSourceWithConfigure = &AuthorizationModelsDataSource{}

func NewAuthorizationModelsDataSource() datasource.DataSource {
	return &AuthorizationModelsDataSource{}
}

type AuthorizationModelsDataSource struct {
	client *client.OpenFgaClient
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
		MarkdownDescription: "An authorization model combines one or more type definitions. This is used to define the permission model of a system.",

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store this authorization model belongs to",
				Required:            true,
			},
			"authorization_models": schema.ListNestedAttribute{
				MarkdownDescription: "List of OpenFGA authorization models",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the OpenFGA authorization model",
							Computed:            true,
						},
						"store_id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the OpenFGA store this authorization model belongs to",
							Computed:            true,
						},
						"model_json": schema.StringAttribute{
							MarkdownDescription: "The full authorization model definition in JSON format",
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

	d.client = client
}

func (d *AuthorizationModelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AuthorizationModelsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var authorizationModels []openfga.AuthorizationModel
	options := client.ClientReadAuthorizationModelsOptions{
		StoreId:           state.StoreId.ValueStringPointer(),
		ContinuationToken: nil,
	}

	for isLastPage := false; !isLastPage; isLastPage = options.ContinuationToken == nil {
		response, err := d.client.ReadAuthorizationModels(ctx).Options(options).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read authorization models, got error: %s", err))
			return
		}

		authorizationModels = append(authorizationModels, response.AuthorizationModels...)

		options.ContinuationToken = response.ContinuationToken
	}

	state.AuthorizationModels = []AuthorizationModelModel{}
	for _, authorizationModel := range authorizationModels {
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

		state.AuthorizationModels = append(state.AuthorizationModels, AuthorizationModelModel{
			Id:        types.StringValue(authorizationModel.Id),
			StoreId:   types.StringValue(*options.StoreId),
			ModelJson: jsontypes.NewNormalizedValue(string(jsonBytes)),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
