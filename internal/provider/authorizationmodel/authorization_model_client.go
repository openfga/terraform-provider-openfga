package authorizationmodel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

type AuthorizationModelClient struct {
	client *client.OpenFgaClient
}

func NewAuthorizationModelClient(client *client.OpenFgaClient) *AuthorizationModelClient {
	return &AuthorizationModelClient{client: client}
}

func (model AuthorizationModelModel) ToCreateRequest() (*client.ClientWriteAuthorizationModelRequest, error) {
	authorizationModel, err := model.ToAuthorizationModel()
	if err != nil {
		return nil, err
	}

	return &client.ClientWriteAuthorizationModelRequest{
		SchemaVersion:   authorizationModel.SchemaVersion,
		TypeDefinitions: authorizationModel.TypeDefinitions,
		Conditions:      authorizationModel.Conditions,
	}, nil
}

func (wrapper *AuthorizationModelClient) CreateAuthorizationModel(ctx context.Context, storeId string, model AuthorizationModelModel) (*AuthorizationModelModel, error) {
	options := client.ClientWriteAuthorizationModelOptions{
		StoreId: openfga.PtrString(storeId),
	}

	body, err := model.ToCreateRequest()
	if err != nil {
		return nil, err
	}

	response, err := wrapper.client.WriteAuthorizationModel(ctx).Options(options).Body(*body).Execute()
	if err != nil {
		return nil, err
	}

	authorizationModelModel := model

	authorizationModelModel.Id = types.StringValue(response.AuthorizationModelId)

	return &authorizationModelModel, nil
}

func (wrapper *AuthorizationModelClient) ReadAuthorizationModel(ctx context.Context, storeId string, model AuthorizationModelModel) (*AuthorizationModelModel, error) {
	options := client.ClientReadAuthorizationModelOptions{
		StoreId:              openfga.PtrString(storeId),
		AuthorizationModelId: openfga.PtrString(model.GetId()),
	}

	response, err := wrapper.client.ReadAuthorizationModel(ctx).Options(options).Execute()
	if err != nil {
		return nil, err
	}

	authorizationModel := *response.AuthorizationModel

	return NewAuthorizationModelModelFromAuthorizationModel(authorizationModel), nil
}

func (wrapper *AuthorizationModelClient) ReadLatestAuthorizationModel(ctx context.Context, storeId string) (*AuthorizationModelModel, error) {
	options := client.ClientReadLatestAuthorizationModelOptions{
		StoreId: openfga.PtrString(storeId),
	}

	response, err := wrapper.client.ReadLatestAuthorizationModel(ctx).Options(options).Execute()
	if err != nil {
		return nil, err
	}

	if response.AuthorizationModel == nil {
		return nil, fmt.Errorf("unable to find authorization model")
	}

	authorizationModel := *response.AuthorizationModel

	return NewAuthorizationModelModelFromAuthorizationModel(authorizationModel), nil
}

func (wrapper *AuthorizationModelClient) ListAuthorizationModels(ctx context.Context, storeId string) (*[]AuthorizationModelModel, error) {
	options := client.ClientReadAuthorizationModelsOptions{
		StoreId:           openfga.PtrString(storeId),
		ContinuationToken: openfga.PtrString(""),
	}

	authorizationModels := []openfga.AuthorizationModel{}

	for isLastPage := false; !isLastPage; isLastPage = *options.ContinuationToken == "" {
		response, err := wrapper.client.ReadAuthorizationModels(ctx).Options(options).Execute()
		if err != nil {
			return nil, err
		}

		authorizationModels = append(authorizationModels, response.AuthorizationModels...)

		if response.ContinuationToken != nil {
			options.ContinuationToken = response.ContinuationToken
		}
	}

	authorizationModelModels := []AuthorizationModelModel{}
	for _, authorizationModel := range authorizationModels {
		authorizationModelModels = append(
			authorizationModelModels,
			*NewAuthorizationModelModelFromAuthorizationModel(authorizationModel),
		)
	}

	return &authorizationModelModels, nil
}
