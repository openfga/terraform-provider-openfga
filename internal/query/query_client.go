package query

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

type QueryClient struct {
	client *client.OpenFgaClient
}

func NewQueryClient(client *client.OpenFgaClient) *QueryClient {
	return &QueryClient{client: client}
}

func (query CheckQueryModel) ToCheckRequest() (*client.ClientCheckRequest, error) {
	context, err := query.GetContextMap()
	if err != nil {
		return nil, err
	}

	contextualTuples := []client.ClientContextualTupleKey{}
	for _, contextualTupleModel := range query.GetContextualTuples() {
		contextualTuple, err := contextualTupleModel.ToTupleWithCondition()
		if err != nil {
			return nil, err
		}

		contextualTuples = append(contextualTuples, *contextualTuple)
	}

	return &client.ClientCheckRequest{
		User:             query.GetUser(),
		Relation:         query.GetRelation(),
		Object:           query.GetObject(),
		ContextualTuples: contextualTuples,
		Context:          context,
	}, nil
}

func (wrapper *QueryClient) Check(ctx context.Context, storeId string, authorizationModelId string, model CheckQueryModel) (types.Bool, error) {
	options := client.ClientCheckOptions{
		StoreId:              openfga.PtrString(storeId),
		AuthorizationModelId: openfga.PtrString(authorizationModelId),
	}

	body, err := model.ToCheckRequest()
	if err != nil {
		return types.BoolNull(), err
	}

	response, err := wrapper.client.Check(ctx).Options(options).Body(*body).Execute()
	if err != nil {
		return types.BoolNull(), err
	}

	return types.BoolValue(response.GetAllowed()), nil
}

func (query ListObjectsQueryModel) ToListObjectsRequest() (*client.ClientListObjectsRequest, error) {
	context, err := query.GetContextMap()
	if err != nil {
		return nil, err
	}

	contextualTuples := []client.ClientContextualTupleKey{}
	for _, contextualTupleModel := range query.GetContextualTuples() {
		contextualTuple, err := contextualTupleModel.ToTupleWithCondition()
		if err != nil {
			return nil, err
		}

		contextualTuples = append(contextualTuples, *contextualTuple)
	}

	return &client.ClientListObjectsRequest{
		User:             query.GetUser(),
		Relation:         query.GetRelation(),
		Type:             query.GetType(),
		ContextualTuples: contextualTuples,
		Context:          context,
	}, nil
}

func (wrapper *QueryClient) ListObjects(ctx context.Context, storeId string, authorizationModelId string, model ListObjectsQueryModel) (types.List, error) {
	options := client.ClientListObjectsOptions{
		StoreId:              openfga.PtrString(storeId),
		AuthorizationModelId: openfga.PtrString(authorizationModelId),
	}

	body, err := model.ToListObjectsRequest()
	if err != nil {
		return types.ListNull(types.StringType), err
	}

	response, err := wrapper.client.ListObjects(ctx).Options(options).Body(*body).Execute()
	if err != nil {
		return types.ListNull(types.StringType), err
	}

	elements := []attr.Value{}
	for _, object := range response.GetObjects() {
		elements = append(elements, types.StringValue(object))
	}

	return types.ListValueMust(types.StringType, elements), nil
}
