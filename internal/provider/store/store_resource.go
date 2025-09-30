package store

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/openfga/go-sdk/client"

	internalError "github.com/openfga/terraform-provider-openfga/internal/apierror"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StoreResource{}
var _ resource.ResourceWithImportState = &StoreResource{}

func NewStoreResource() resource.Resource {
	return &StoreResource{}
}

type StoreResource struct {
	client *StoreClient
}

type StoreResourceModel struct {
	StoreModel
}

func (r *StoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func (r *StoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Provides the ability to create and manage OpenFGA stores.

A store is a logical container for authorization data. Each store contains one or more versions of an authorization model and can contain various relationship tuples.

Separate stores can be created for separate authorization needs or isolated environments, e.g. development/ prod.
`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the store.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the store.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *StoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = NewStoreClient(client)
}

func (r *StoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state StoreResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storeModel, err := r.client.CreateStore(ctx, state.StoreModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create store, got error: %s", err))
		return
	}

	state.StoreModel = *storeModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StoreResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storeModel, err := r.client.ReadStore(ctx, state.StoreModel)
	if err != nil {
		if internalError.IsStatusNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Store not found",
				fmt.Sprintf("Store %q no longer exists; removing from state.", state.Id.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read store, got error: %s", err))
		return
	}

	state.StoreModel = *storeModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported and should never be called
	resp.Diagnostics.AddError(
		"Client Error",
		"Unable to update store, update operation is not supported",
	)
}

func (r *StoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StoreResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStore(ctx, state.StoreModel)
	if err != nil {
		if internalError.IsStatusNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete store, got error: %s", err))
		return
	}
}

func (r *StoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
