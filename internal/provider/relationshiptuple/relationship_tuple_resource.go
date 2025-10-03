package relationshiptuple

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"
	internalError "github.com/openfga/terraform-provider-openfga/internal/apierror"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RelationshipTupleResource{}
var _ resource.ResourceWithImportState = &RelationshipTupleResource{}

func NewRelationshipTupleResource() resource.Resource {
	return &RelationshipTupleResource{}
}

type RelationshipTupleResource struct {
	client *RelationshipTupleClient
}

type RelationshipTupleResourceModel struct {
	StoreId              types.String `tfsdk:"store_id"`
	AuthorizationModelId types.String `tfsdk:"authorization_model_id"`
	RelationshipTupleWithConditionModel
}

func (r *RelationshipTupleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_tuple"
}

func (r *RelationshipTupleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Provides the ability to create and manage OpenFGA relationship tuples.

A relationship tuple consists of a user, a relation, and an object. Tuples may include an additional condition, which has to be fulfilled for the tuple to be considered.

Together with an authorization model, the relationship tuples determine whether a relationship exists between a user and an object.
`,

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the store this relationship tuple belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"authorization_model_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the authorization model this relationship tuple is related with. Can be left blank to refer to the latest authorization model.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The user of the relationship tuple.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"relation": schema.StringAttribute{
				MarkdownDescription: "The relation of the relationship tuple.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object of the relationship tuple.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"condition": schema.SingleNestedAttribute{
				MarkdownDescription: "A condition of the relationship tuple.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the condition.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"context_json": schema.StringAttribute{
						MarkdownDescription: "The (partial) context under which the condition is evaluated.",
						CustomType:          jsontypes.NormalizedType{},
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *RelationshipTupleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = NewRelationshipTupleClient(client)
}

func (r *RelationshipTupleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state RelationshipTupleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	relationshipTupleModel, err := r.client.CreateRelationshipTuple(ctx, state.StoreId.ValueString(), state.AuthorizationModelId.ValueStringPointer(), state.RelationshipTupleWithConditionModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create relationship tuple, got error: %s", err))
		return
	}

	state.RelationshipTupleWithConditionModel = *relationshipTupleModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RelationshipTupleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RelationshipTupleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	relationshipTupleModel, err := r.client.ReadRelationshipTuple(ctx, state.StoreId.ValueString(), state.RelationshipTupleModel)
	if err != nil {
		if internalError.IsExpectedOneResultError(err) {
			resp.Diagnostics.AddWarning(
				"relationship tuple not found",
				fmt.Sprintf("relationship tuple (user: %q, relation: %q, object: %q) no longer exists; removing from state.",
					state.User.ValueString(),
					state.Relation.ValueString(),
					state.Object.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read relationship tuple, got error: %s", err))
		return
	}

	state.RelationshipTupleWithConditionModel = *relationshipTupleModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RelationshipTupleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported and should never be called
	resp.Diagnostics.AddError(
		"Client Error",
		"Unable to update relationship tuple, update operation is not supported",
	)
}

func (r *RelationshipTupleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RelationshipTupleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRelationshipTuple(ctx, state.StoreId.ValueString(), state.AuthorizationModelId.ValueStringPointer(), state.RelationshipTupleWithConditionModel)
	if err != nil {
		if internalError.IsExpectedOneResultError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete relationship tuple, got error: %s", err))
		return
	}
}

func (r *RelationshipTupleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")

	var state RelationshipTupleResourceModel
	if len(parts) == 4 {
		state = RelationshipTupleResourceModel{
			StoreId:                             types.StringValue(parts[0]),
			RelationshipTupleWithConditionModel: *NewRelationshipTupleWithConditionModel(parts[1], parts[2], parts[3], nil),
		}
	} else if len(parts) == 5 {
		state = RelationshipTupleResourceModel{
			StoreId:                             types.StringValue(parts[0]),
			AuthorizationModelId:                types.StringValue(parts[1]),
			RelationshipTupleWithConditionModel: *NewRelationshipTupleWithConditionModel(parts[2], parts[3], parts[4], nil),
		}
	} else {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Input ID has to be in the format of <store_id>/<user>/<relation>/<object> or <store_id>/<authorization_model_id>/<user>/<relation>/<object>, but received: %s", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
