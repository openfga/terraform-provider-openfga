package relationshiptuple

import (
	"context"
	"fmt"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

type RelationshipTupleClient struct {
	client *client.OpenFgaClient
}

func NewRelationshipTupleClient(client *client.OpenFgaClient) *RelationshipTupleClient {
	return &RelationshipTupleClient{client: client}
}

func (model RelationshipTupleWithConditionModel) ToCreateRequest() (*client.ClientWriteTuplesBody, error) {
	tuple, err := model.ToTupleWithCondition()
	if err != nil {
		return nil, err
	}

	return &client.ClientWriteTuplesBody{
		*tuple,
	}, nil
}

func (wrapper *RelationshipTupleClient) CreateRelationshipTuple(ctx context.Context, storeId string, model RelationshipTupleWithConditionModel) (*RelationshipTupleWithConditionModel, error) {
	options := client.ClientWriteOptions{
		StoreId: openfga.PtrString(storeId),
	}

	body, err := model.ToCreateRequest()
	if err != nil {
		return nil, err
	}

	response, err := wrapper.client.WriteTuples(ctx).Options(options).Body(*body).Execute()
	if err != nil {
		return nil, err
	}

	writeResult := response.Writes[0]
	if writeResult.Error != nil {
		return nil, writeResult.Error
	}

	tuple := writeResult.TupleKey

	return NewRelationshipTupleWithConditionModelFromTuple(&tuple), nil
}

func (model RelationshipTupleModel) ToReadRequest() *client.ClientReadRequest {
	tuple := model.ToTuple()

	return &client.ClientReadRequest{
		User:     openfga.PtrString(tuple.GetUser()),
		Relation: openfga.PtrString(tuple.GetRelation()),
		Object:   openfga.PtrString(tuple.GetObject()),
	}
}

func (wrapper *RelationshipTupleClient) ReadRelationshipTuple(ctx context.Context, storeId string, model RelationshipTupleModel) (*RelationshipTupleWithConditionModel, error) {
	options := client.ClientReadOptions{
		StoreId: openfga.PtrString(storeId),
	}

	body := model.ToReadRequest()

	response, err := wrapper.client.Read(ctx).Options(options).Body(*body).Execute()
	if err != nil {
		return nil, err
	}

	tuples := response.Tuples

	if len(tuples) != 1 {
		return nil, fmt.Errorf("expected one result but received: %d", len(tuples))
	}

	tuple := tuples[0].Key

	return NewRelationshipTupleWithConditionModelFromTuple(&tuple), nil
}

func (wrapper *RelationshipTupleClient) ListRelationshipTuples(ctx context.Context, storeId string, query *RelationshipTupleModel) (*[]RelationshipTupleWithConditionModel, error) {
	options := client.ClientReadOptions{
		StoreId:           openfga.PtrString(storeId),
		ContinuationToken: openfga.PtrString(""),
	}

	body := client.ClientReadRequest{}
	if query != nil {
		body = *query.ToReadRequest()
	}

	tuples := []openfga.TupleKey{}

	for isLastPage := false; !isLastPage; isLastPage = *options.ContinuationToken == "" {
		response, err := wrapper.client.Read(ctx).Options(options).Body(body).Execute()
		if err != nil {
			return nil, err
		}

		for _, element := range response.Tuples {
			tuples = append(tuples, element.Key)
		}

		options.ContinuationToken = openfga.PtrString(response.ContinuationToken)
	}

	relationshipTupleModels := []RelationshipTupleWithConditionModel{}
	for _, tuple := range tuples {
		relationshipTupleModels = append(
			relationshipTupleModels,
			*NewRelationshipTupleWithConditionModelFromTuple(&tuple),
		)
	}

	return &relationshipTupleModels, nil
}

func (model RelationshipTupleModel) ToDeleteRequest() *client.ClientDeleteTuplesBody {
	return &client.ClientDeleteTuplesBody{
		*model.ToTuple(),
	}
}

func (wrapper *RelationshipTupleClient) DeleteRelationshipTuple(ctx context.Context, storeId string, model RelationshipTupleWithConditionModel) error {
	options := client.ClientWriteOptions{
		StoreId: openfga.PtrString(storeId),
	}

	body := model.ToDeleteRequest()

	response, err := wrapper.client.DeleteTuples(ctx).Options(options).Body(*body).Execute()
	if err != nil {
		return err
	}

	deleteResult := response.Deletes[0]
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	return nil
}
