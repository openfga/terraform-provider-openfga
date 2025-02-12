package relationshiptuple

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RelationshipTupleResource{}
var _ resource.ResourceWithImportState = &RelationshipTupleResource{}

func NewRelationshipTupleResource() resource.Resource {
	return &RelationshipTupleResource{}
}

type RelationshipTupleResource struct {
	client *client.OpenFgaClient
}

type RelationshipTupleResourceModel RelationshipTupleModel

func (r *RelationshipTupleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_tuple"
}

func (r *RelationshipTupleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A relationship tuple is a tuple consisting of a user, relation, and object. Tuples may add an optional condition.",

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store this relationship tuple belongs to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The user of the OpenFGA relationship tuple",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"relation": schema.StringAttribute{
				MarkdownDescription: "The relation of the OpenFGA relationship tuple",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object of the OpenFGA relationship tuple",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"condition": schema.SingleNestedAttribute{
				MarkdownDescription: "A condition of the OpenFGA relationship tuple",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the condition",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"context_json": schema.StringAttribute{
						MarkdownDescription: "The (partial) context under which the condition is evaluated",
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

	r.client = client
}

func (r *RelationshipTupleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state RelationshipTupleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	options := client.ClientWriteOptions{
		StoreId: openfga.PtrString(state.StoreId.ValueString()),
	}

	var (
		context   *map[string]interface{}
		condition *openfga.RelationshipCondition
	)
	if state.Condition != nil {
		if !state.Condition.Context.IsNull() {
			var result map[string]interface{}

			resp.Diagnostics.Append(state.Condition.Context.Unmarshal(&result)...)
			if resp.Diagnostics.HasError() {
				return
			}

			context = &result
		}

		condition = &openfga.RelationshipCondition{
			Name:    state.Condition.Name.ValueString(),
			Context: context,
		}
	}

	tuple := openfga.TupleKey{
		User:      state.User.ValueString(),
		Relation:  state.Relation.ValueString(),
		Object:    state.Object.ValueString(),
		Condition: condition,
	}

	body := client.ClientWriteRequest{
		Writes: []client.ClientTupleKey{
			tuple,
		},
	}

	response, err := r.client.Write(ctx).Options(options).Body(body).Execute()
	if err != nil || response.Writes[0].Error != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create relationship tuple, got error: %s", err))
		return
	}

	tupleResponse := response.Writes[0].TupleKey

	state.User = types.StringValue(tupleResponse.User)
	state.Relation = types.StringValue(tupleResponse.Relation)
	state.Object = types.StringValue(tupleResponse.Object)
	if tupleResponse.Condition != nil {
		context := jsontypes.NewNormalizedNull()
		if tupleResponse.Condition.Context != nil {
			jsonBytes, err := json.Marshal(tupleResponse.Condition.Context)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to convert relationship tuple condition context to JSON, got error: %s", err))
				return
			}

			context = jsontypes.NewNormalizedValue(string(jsonBytes))
		}

		state.Condition = &RelationshipTupleCondition{
			Name:    types.StringValue(tupleResponse.Condition.Name),
			Context: context,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RelationshipTupleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RelationshipTupleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	options := client.ClientReadOptions{
		StoreId: state.StoreId.ValueStringPointer(),
	}

	body := client.ClientReadRequest{
		User:     state.User.ValueStringPointer(),
		Relation: state.Relation.ValueStringPointer(),
		Object:   state.Object.ValueStringPointer(),
	}

	response, err := r.client.Read(ctx).Options(options).Body(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read relationship tuple, got error: %s", err))
		return
	}

	if len(response.Tuples) != 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read relationship tuple, expected one result but received: %d", len(response.Tuples)))
		return
	}

	tuple := response.Tuples[0]

	state.User = types.StringValue(tuple.Key.User)
	state.Relation = types.StringValue(tuple.Key.Relation)
	state.Object = types.StringValue(tuple.Key.Object)
	if tuple.Key.Condition != nil {
		context := jsontypes.NewNormalizedNull()
		if tuple.Key.Condition.Context != nil {
			jsonBytes, err := json.Marshal(tuple.Key.Condition.Context)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to convert relationship tuple condition context to JSON, got error: %s", err))
				return
			}

			context = jsontypes.NewNormalizedValue(string(jsonBytes))
		}

		state.Condition = &RelationshipTupleCondition{
			Name:    types.StringValue(tuple.Key.Condition.Name),
			Context: context,
		}
	}

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

	options := client.ClientWriteOptions{
		StoreId: openfga.PtrString(state.StoreId.ValueString()),
	}

	body := client.ClientWriteRequest{
		Deletes: []client.ClientTupleKeyWithoutCondition{
			{
				User:     state.User.ValueString(),
				Relation: state.Relation.ValueString(),
				Object:   state.Object.ValueString(),
			},
		},
	}

	response, err := r.client.Write(ctx).Options(options).Body(body).Execute()
	if err != nil || response.Deletes[0].Error != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete relationship tuple, got error: %s", err))
		return
	}
}

func (r *RelationshipTupleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")

	if len(parts) != 4 {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Input ID has to be in the format of <store_id>/<user>/<relation>/<object>, but received: %s", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("store_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("relation"), parts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object"), parts[3])...)
}
