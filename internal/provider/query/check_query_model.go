package query

import (
	"github.com/mauriceackel/terraform-provider-openfga/internal/provider/relationshiptuple"
)

type CheckQueryModel struct {
	relationshiptuple.RelationshipTupleModel
	ContextualTuples *[]relationshiptuple.RelationshipTupleWithConditionModel `tfsdk:"contextual_tuples"`
	relationshiptuple.ContextModel
}

func (query CheckQueryModel) GetContextualTuples() []relationshiptuple.RelationshipTupleWithConditionModel {
	if query.ContextualTuples == nil {
		return []relationshiptuple.RelationshipTupleWithConditionModel{}
	}

	return *query.ContextualTuples
}

func NewCheckQueryModel(user string, relation string, object string, contextualTuples *[]relationshiptuple.RelationshipTupleWithConditionModel, context *map[string]interface{}) *CheckQueryModel {
	return &CheckQueryModel{
		RelationshipTupleModel: *relationshiptuple.NewRelationshipTupleModel(user, relation, object),
		ContextualTuples:       contextualTuples,
		ContextModel:           *relationshiptuple.NewContextModel(context),
	}
}
