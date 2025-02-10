package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StoreResource{}
var _ resource.ResourceWithImportState = &StoreResource{}

func NewStoreResource() resource.Resource {
	return &StoreResource{}
}

type StoreResource struct {
	client *client.OpenFgaClient
}

type StoreResourceModel StoreModel

func (r *StoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func (r *StoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A store is a logical container for authorization data, including authorization models and tuples.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the OpenFGA store",
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

	r.client = client
}

func (r *StoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state StoreResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := client.ClientCreateStoreRequest{
		Name: state.Name.ValueString(),
	}

	response, err := r.client.CreateStore(ctx).Body(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create store, got error: %s", err))
		return
	}

	state.Id = types.StringValue(response.Id)
	state.Name = types.StringValue(response.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StoreResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	options := client.ClientGetStoreOptions{
		StoreId: state.Id.ValueStringPointer(),
	}

	response, err := r.client.GetStore(ctx).Options(options).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read store, got error: %s", err))
		return
	}

	state.Id = types.StringValue(response.Id)
	state.Name = types.StringValue(response.Name)

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

	options := client.ClientDeleteStoreOptions{
		StoreId: state.Id.ValueStringPointer(),
	}

	_, err := r.client.DeleteStore(ctx).Options(options).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete store, got error: %s", err))
		return
	}
}

func (r *StoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
