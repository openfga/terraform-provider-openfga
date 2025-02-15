package relationshiptuple

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
var _ datasource.DataSource = &RelationshipTuplesDataSource{}
var _ datasource.DataSourceWithConfigure = &RelationshipTuplesDataSource{}

func NewRelationshipTuplesDataSource() datasource.DataSource {
	return &RelationshipTuplesDataSource{}
}

type RelationshipTuplesDataSource struct {
	client *RelationshipTupleClient
}

type RelationshipTuplesDataSourceModel struct {
	StoreId            types.String                          `tfsdk:"store_id"`
	Query              *RelationshipTupleModel               `tfsdk:"query"`
	RelationshipTuples []RelationshipTupleWithConditionModel `tfsdk:"relationship_tuples"`
}

func (d *RelationshipTuplesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_tuples"
}

func (d *RelationshipTuplesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A relationship tuple is a tuple consisting of a user, relation, and object. Tuples may add an optional condition.",

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store this relationship tuple belongs to",
				Required:            true,
			},
			"query": schema.SingleNestedAttribute{
				MarkdownDescription: "A query to filter the returned tuples (leave empty to read all tuples)",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						MarkdownDescription: "The user of the OpenFGA relationship tuple",
						Optional:            true,
					},
					"relation": schema.StringAttribute{
						MarkdownDescription: "The relation of the OpenFGA relationship tuple",
						Optional:            true,
					},
					"object": schema.StringAttribute{
						MarkdownDescription: "The object of the OpenFGA relationship tuple",
						Required:            true,
					},
				},
			},
			"relationship_tuples": schema.ListNestedAttribute{
				MarkdownDescription: "List of OpenFGA relationship tuples",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user": schema.StringAttribute{
							MarkdownDescription: "The user of the OpenFGA relationship tuple",
							Computed:            true,
						},
						"relation": schema.StringAttribute{
							MarkdownDescription: "The relation of the OpenFGA relationship tuple",
							Computed:            true,
						},
						"object": schema.StringAttribute{
							MarkdownDescription: "The object of the OpenFGA relationship tuple",
							Computed:            true,
						},
						"condition": schema.SingleNestedAttribute{
							MarkdownDescription: "A condition of the OpenFGA relationship tuple",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "The name of the condition",
									Computed:            true,
								},
								"context_json": schema.StringAttribute{
									MarkdownDescription: "The (partial) context under which the condition is evaluated",
									CustomType:          jsontypes.NormalizedType{},
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *RelationshipTuplesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = NewRelationshipTupleClient(client)
}

func (d *RelationshipTuplesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RelationshipTuplesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	relationshipTupleModels, err := d.client.ListRelationshipTuples(ctx, state.StoreId.ValueString(), state.Query)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read relationship tuples, got error: %s", err))
		return
	}

	state.RelationshipTuples = *relationshipTupleModels

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
