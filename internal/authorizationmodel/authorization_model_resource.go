package authorizationmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AuthorizationModelResource{}
var _ resource.ResourceWithImportState = &AuthorizationModelResource{}

func NewAuthorizationModelResource() resource.Resource {
	return &AuthorizationModelResource{}
}

type AuthorizationModelResource struct {
	client *client.OpenFgaClient
}

type AuthorizationModelResourceModel AuthorizationModelModel

func (r *AuthorizationModelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_model"
}

func (r *AuthorizationModelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "An authorization model combines one or more type definitions. This is used to define the permission model of a system.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA authorization model",
				Computed:            true,
			},
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store this authorization model belongs to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"model_json": schema.StringAttribute{
				MarkdownDescription: "The full authorization model definition in JSON format",
				Required:            true,
				CustomType:          jsontypes.NormalizedType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *AuthorizationModelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.OpenFgaClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.OpenFgaClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *AuthorizationModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state AuthorizationModelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var authorizationModel AuthorizationModelWithoutId
	err := json.Unmarshal([]byte(state.ModelJson.ValueString()), &authorizationModel)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Model JSON", fmt.Sprintf("Unable to parse model JSON, got error: %s", err))
		return
	}

	options := client.ClientWriteAuthorizationModelOptions{
		StoreId: state.StoreId.ValueStringPointer(),
	}

	body := client.ClientWriteAuthorizationModelRequest{
		SchemaVersion:   authorizationModel.SchemaVersion,
		TypeDefinitions: authorizationModel.TypeDefinitions,
		Conditions:      authorizationModel.Conditions,
	}

	response, err := r.client.WriteAuthorizationModel(ctx).Options(options).Body(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create authorization model, got error: %s", err))
		return
	}

	state.Id = types.StringValue(response.AuthorizationModelId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AuthorizationModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AuthorizationModelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	options := client.ClientReadAuthorizationModelOptions{
		StoreId:              state.StoreId.ValueStringPointer(),
		AuthorizationModelId: state.Id.ValueStringPointer(),
	}

	response, err := r.client.ReadAuthorizationModel(ctx).Options(options).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read authorization model, got error: %s", err))
		return
	}

	state.Id = types.StringValue(response.AuthorizationModel.Id)

	authorizationModelWithoutId := AuthorizationModelWithoutId{
		SchemaVersion:   response.AuthorizationModel.SchemaVersion,
		TypeDefinitions: response.AuthorizationModel.TypeDefinitions,
		Conditions:      response.AuthorizationModel.Conditions,
	}

	jsonBytes, err := json.Marshal(authorizationModelWithoutId)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Response Data", fmt.Sprintf("Unable to convert to model JSON, got error: %s", err))
		return
	}
	state.ModelJson = jsontypes.NewNormalizedValue(string(jsonBytes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AuthorizationModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported and should never be called
	resp.Diagnostics.AddError(
		"Client Error",
		"Unable to update authorization model, update operation is not supported",
	)
}

func (r *AuthorizationModelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Deletion is not possible, we treat it as a noop
}

func (r *AuthorizationModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")

	if len(parts) != 2 {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Input ID has to be in the format of <store_id>/<authorization_model_id>, but received: %s", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("store_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
