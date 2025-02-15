package store

import (
	"context"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

type StoreClient struct {
	client *client.OpenFgaClient
}

func NewStoreClient(client *client.OpenFgaClient) *StoreClient {
	return &StoreClient{client: client}
}

func (model StoreModel) ToCreateRequest() *client.ClientCreateStoreRequest {
	store := model.ToStore()

	return &client.ClientCreateStoreRequest{
		Name: store.GetName(),
	}
}

func (wrapper *StoreClient) CreateStore(ctx context.Context, model StoreModel) (*StoreModel, error) {
	options := client.ClientCreateStoreOptions{}

	body := *model.ToCreateRequest()

	response, err := wrapper.client.CreateStore(ctx).Options(options).Body(body).Execute()
	if err != nil {
		return nil, err
	}

	return NewStoreModelFromStore(response), nil
}

func (wrapper *StoreClient) ReadStore(ctx context.Context, model StoreModel) (*StoreModel, error) {
	options := client.ClientGetStoreOptions{
		StoreId: openfga.PtrString(model.GetId()),
	}

	response, err := wrapper.client.GetStore(ctx).Options(options).Execute()
	if err != nil {
		return nil, err
	}

	return NewStoreModelFromStore(response), nil
}

func (wrapper *StoreClient) ListStores(ctx context.Context) (*[]StoreModel, error) {
	options := client.ClientListStoresOptions{
		ContinuationToken: openfga.PtrString(""),
	}

	stores := []openfga.Store{}

	for isLastPage := false; !isLastPage; isLastPage = *options.ContinuationToken == "" {
		response, err := wrapper.client.ListStores(ctx).Options(options).Execute()
		if err != nil {
			return nil, err
		}

		stores = append(stores, response.Stores...)

		options.ContinuationToken = openfga.PtrString(response.ContinuationToken)
	}

	storeModels := []StoreModel{}
	for _, store := range stores {
		storeModels = append(
			storeModels,
			*NewStoreModelFromStore(&store),
		)
	}

	return &storeModels, nil
}

func (wrapper *StoreClient) DeleteStore(ctx context.Context, model StoreModel) error {
	options := client.ClientDeleteStoreOptions{
		StoreId: openfga.PtrString(model.GetId()),
	}

	_, err := wrapper.client.DeleteStore(ctx).Options(options).Execute()
	if err != nil {
		return err
	}

	return nil
}
