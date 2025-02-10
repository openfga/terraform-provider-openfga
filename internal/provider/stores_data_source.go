// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StoresDataSource{}
var _ datasource.DataSourceWithConfigure = &StoresDataSource{}

func NewStoresDataSource() datasource.DataSource {
	return &StoresDataSource{}
}

type StoresDataSource struct {
	client *client.OpenFgaClient
}

type StoresDataSourceModel struct {
	Stores []StoreModel `tfsdk:"stores"`
}

func (d *StoresDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stores"
}

func (d *StoresDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A store is a logical container for authorization data, including authorization models and tuples.",

		Attributes: map[string]schema.Attribute{
			"stores": schema.ListNestedAttribute{
				MarkdownDescription: "List of OpenFGA stores",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the OpenFGA store",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the OpenFGA store",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *StoresDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StoresDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state StoresDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var stores []openfga.Store
	options := client.ClientListStoresOptions{
		ContinuationToken: openfga.PtrString(""),
	}

	for isLastPage := false; !isLastPage; isLastPage = *options.ContinuationToken == "" {

		response, err := d.client.ListStores(ctx).Options(options).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read stores, got error: %s", err))
			return
		}

		stores = append(stores, response.Stores...)

		options.ContinuationToken = openfga.PtrString(response.ContinuationToken)
	}

	state.Stores = []StoreModel{}
	for _, store := range stores {
		state.Stores = append(state.Stores, StoreModel{
			Id:   types.StringValue(store.Id),
			Name: types.StringValue(store.Name),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
