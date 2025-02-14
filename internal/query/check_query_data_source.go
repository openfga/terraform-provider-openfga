package query

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CheckQueryDataSource{}
var _ datasource.DataSourceWithConfigure = &CheckQueryDataSource{}

func NewCheckQueryDataSource() datasource.DataSource {
	return &CheckQueryDataSource{}
}

type CheckQueryDataSource struct {
	client *client.OpenFgaClient
}

type CheckQueryDataSourceModel struct {
	StoreId              types.String         `tfsdk:"store_id"`
	AuthorizationModelId types.String         `tfsdk:"authorization_model_id"`
	User                 types.String         `tfsdk:"user"`
	Relation             types.String         `tfsdk:"relation"`
	Object               types.String         `tfsdk:"object"`
	ContextualTuples     *[]ContextualTuple   `tfsdk:"contextual_tuples"`
	Context              jsontypes.Normalized `tfsdk:"context_json"`

	Result types.Bool `tfsdk:"result"`
}

type ContextualTuple struct {
	User      types.String                `tfsdk:"user"`
	Relation  types.String                `tfsdk:"relation"`
	Object    types.String                `tfsdk:"object"`
	Condition *RelationshipTupleCondition `tfsdk:"condition"`
}

type RelationshipTupleCondition struct {
	Name    types.String         `tfsdk:"name"`
	Context jsontypes.Normalized `tfsdk:"context_json"`
}

func (d *CheckQueryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_query"
}

func (d *CheckQueryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A 'check' query can be called to establish whether a particular user has a specific relationship with a particular object.",

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store this query is run against",
				Required:            true,
			},
			"authorization_model_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA authorization model this query is run against",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The user of the query",
				Required:            true,
			},
			"relation": schema.StringAttribute{
				MarkdownDescription: "The relation to check for",
				Required:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object of the query",
				Required:            true,
			},
			"contextual_tuples": schema.ListNestedAttribute{
				MarkdownDescription: "The object of the query",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user": schema.StringAttribute{
							MarkdownDescription: "The user of the contextual relationship tuple",
							Required:            true,
						},
						"relation": schema.StringAttribute{
							MarkdownDescription: "The relation of the contextual relationship tuple",
							Required:            true,
						},
						"object": schema.StringAttribute{
							MarkdownDescription: "The object of the contextual relationship tuple",
							Required:            true,
						},
						"condition": schema.SingleNestedAttribute{
							MarkdownDescription: "A condition of the contextual relationship tuple",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "The name of the condition",
									Required:            true,
								},
								"context_json": schema.StringAttribute{
									MarkdownDescription: "The (partial) context under which the condition is evaluated",
									CustomType:          jsontypes.NormalizedType{},
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"context_json": schema.StringAttribute{
				MarkdownDescription: "The (partial) context under which the condition is evaluated",
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
			},
			"result": schema.BoolAttribute{
				MarkdownDescription: "The result of the check query",
				Computed:            true,
			},
		},
	}
}

func (d *CheckQueryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CheckQueryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state CheckQueryDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	options := client.ClientCheckOptions{
		StoreId:              state.StoreId.ValueStringPointer(),
		AuthorizationModelId: state.AuthorizationModelId.ValueStringPointer(),
	}

	var contextualTuples []client.ClientContextualTupleKey = []client.ClientContextualTupleKey{}
	if state.ContextualTuples != nil {
		for _, contextualTuple := range *state.ContextualTuples {
			var condition *openfga.RelationshipCondition
			if contextualTuple.Condition != nil {
				var context *map[string]interface{}
				if !contextualTuple.Condition.Context.IsNull() {
					var result map[string]interface{}

					resp.Diagnostics.Append(contextualTuple.Condition.Context.Unmarshal(&result)...)
					if resp.Diagnostics.HasError() {
						return
					}

					context = &result
				}

				condition = &openfga.RelationshipCondition{
					Name:    contextualTuple.Condition.Name.ValueString(),
					Context: context,
				}
			}

			contextualTuples = append(contextualTuples, client.ClientContextualTupleKey{
				User:      contextualTuple.User.ValueString(),
				Relation:  contextualTuple.Relation.ValueString(),
				Object:    contextualTuple.Object.ValueString(),
				Condition: condition,
			})
		}
	}

	var context *map[string]interface{}
	if !state.Context.IsNull() {
		var result map[string]interface{}

		resp.Diagnostics.Append(state.Context.Unmarshal(&result)...)
		if resp.Diagnostics.HasError() {
			return
		}

		context = &result
	}

	body := client.ClientCheckRequest{
		User:             state.User.ValueString(),
		Relation:         state.Relation.ValueString(),
		Object:           state.Object.ValueString(),
		ContextualTuples: contextualTuples,
		Context:          context,
	}

	response, err := d.client.Check(ctx).Options(options).Body(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to perform check query, got error: %s", err))
		return
	}

	state.Result = types.BoolValue(response.GetAllowed())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
