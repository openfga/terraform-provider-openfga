package store

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StoreDataSource{}
var _ datasource.DataSourceWithConfigure = &StoreDataSource{}

func NewStoreDataSource() datasource.DataSource {
	return &StoreDataSource{}
}

type StoreDataSource struct {
	client *client.OpenFgaClient
}

type StoreDataSourceModel StoreModel

func (d *StoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func (d *StoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A store is a logical container for authorization data, including authorization models and tuples.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the OpenFGA store",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the OpenFGA store",
				Computed:            true,
			},
		},
	}
}

func (d *StoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state StoreDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	options := client.ClientGetStoreOptions{
		StoreId: state.Id.ValueStringPointer(),
	}

	response, err := d.client.GetStore(ctx).Options(options).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read store, got error: %s", err))
		return
	}

	state.Id = types.StringValue(response.Id)
	state.Name = types.StringValue(response.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
