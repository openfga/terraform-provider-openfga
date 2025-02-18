package store

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StoreDataSource{}
var _ datasource.DataSourceWithConfigure = &StoreDataSource{}

func NewStoreDataSource() datasource.DataSource {
	return &StoreDataSource{}
}

type StoreDataSource struct {
	client *StoreClient
}

type StoreDataSourceModel struct {
	StoreModel
}

func (d *StoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func (d *StoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides the ability to retrieve details of an existing OpenFGA store.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the store.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the store.",
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

	d.client = NewStoreClient(client)
}

func (d *StoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state StoreDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storeModel, err := d.client.ReadStore(ctx, state.StoreModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read store, got error: %s", err))
		return
	}

	state.StoreModel = *storeModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
