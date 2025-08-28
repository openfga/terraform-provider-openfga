package authorizationmodel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"

	internalError "github.com/openfga/terraform-provider-openfga/internal/apierror"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AuthorizationModelResource{}
var _ resource.ResourceWithImportState = &AuthorizationModelResource{}

func NewAuthorizationModelResource() resource.Resource {
	return &AuthorizationModelResource{}
}

type AuthorizationModelResource struct {
	client *AuthorizationModelClient
}

type AuthorizationModelResourceModel struct {
	StoreId types.String `tfsdk:"store_id"`
	AuthorizationModelModel
}

func (r *AuthorizationModelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_model"
}

func (r *AuthorizationModelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Provides the ability to create and manage OpenFGA authorization models.

An authorization model defines one or more type definitions and optionally a set of conditions.

Together with relationship tuples, the authorization model determines whether a relationship exists between a user and an object.

~> We suggest using [` + "`openfga_authorization_model_document`" + `](../data-sources/authorization_model_document) when assigning a value to ` + "`model_json`" + `. This allows to use models in different formats (e.g. DSL, JSON, Terraform native) and prevents potential complications arising from formatting discrepancies, whitespace inconsistencies, and other nuances inherent to JSON.
`,

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the store this authorization model belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the authorization model.",
				Computed:            true,
			},
			"model_json": schema.StringAttribute{
				MarkdownDescription: "The authorization model definition in JSON format. Consider using [`openfga_authorization_model_document`](../data-sources/authorization_model_document) to set this field.",
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

	r.client = NewAuthorizationModelClient(client)
}

func (r *AuthorizationModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state AuthorizationModelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	authorizationModelModel, err := r.client.CreateAuthorizationModel(ctx, state.StoreId.ValueString(), state.AuthorizationModelModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create authorization model, got error: %s", err))
		return
	}

	state.AuthorizationModelModel = *authorizationModelModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AuthorizationModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AuthorizationModelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	authorizationModelModel, err := r.client.ReadAuthorizationModel(ctx, state.StoreId.ValueString(), state.AuthorizationModelModel)
	if err != nil {
		if internalError.IsStatusNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read authorization model, got error: %s", err))
		return
	}

	state.AuthorizationModelModel = *authorizationModelModel

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

	state := AuthorizationModelResourceModel{
		StoreId:                 types.StringValue(parts[0]),
		AuthorizationModelModel: *NewAuthorizationModelModel(parts[1]),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
