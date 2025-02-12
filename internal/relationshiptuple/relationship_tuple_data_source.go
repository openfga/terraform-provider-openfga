package relationshiptuple

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RelationshipTupleDataSource{}
var _ datasource.DataSourceWithConfigure = &RelationshipTupleDataSource{}

func NewRelationshipTupleDataSource() datasource.DataSource {
	return &RelationshipTupleDataSource{}
}

type RelationshipTupleDataSource struct {
	client *client.OpenFgaClient
}

type RelationshipTupleDataSourceModel RelationshipTupleModel

func (d *RelationshipTupleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_tuple"
}

func (d *RelationshipTupleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A relationship tuple is a tuple consisting of a user, relation, and object. Tuples may add an optional condition.",

		Attributes: map[string]schema.Attribute{
			"store_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store this relationship tuple belongs to",
				Required:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The user of the OpenFGA relationship tuple",
				Required:            true,
			},
			"relation": schema.StringAttribute{
				MarkdownDescription: "The relation of the OpenFGA relationship tuple",
				Required:            true,
			},
			"object": schema.StringAttribute{
				MarkdownDescription: "The object of the OpenFGA relationship tuple",
				Required:            true,
			},
		},
	}
}

func (d *RelationshipTupleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RelationshipTupleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RelationshipTupleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

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

	response, err := d.client.Read(ctx).Options(options).Body(body).Execute()
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
