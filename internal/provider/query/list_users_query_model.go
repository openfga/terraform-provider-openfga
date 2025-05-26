package query

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/terraform-provider-openfga/internal/provider/relationshiptuple"
)

type ListUsersQueryModel struct {
	Type             types.String                                             `tfsdk:"type"`
	Relation         types.String                                             `tfsdk:"relation"`
	Object           types.String                                             `tfsdk:"object"`
	ContextualTuples *[]relationshiptuple.RelationshipTupleWithConditionModel `tfsdk:"contextual_tuples"`
	relationshiptuple.ContextModel
}

func (query ListUsersQueryModel) GetType() string {
	return query.Type.ValueString()
}

func (query ListUsersQueryModel) GetRelation() string {
	return query.Relation.ValueString()
}

func (query ListUsersQueryModel) GetObject() string {
	return query.Object.ValueString()
}

func (query ListUsersQueryModel) GetContextualTuples() []relationshiptuple.RelationshipTupleWithConditionModel {
	if query.ContextualTuples == nil {
		return []relationshiptuple.RelationshipTupleWithConditionModel{}
	}

	return *query.ContextualTuples
}

func NewListUsersQueryModel(type_ string, relation string, object string, contextualTuples *[]relationshiptuple.RelationshipTupleWithConditionModel, context *map[string]interface{}) *ListUsersQueryModel {
	return &ListUsersQueryModel{
		Type:             types.StringValue(type_),
		Relation:         types.StringValue(relation),
		Object:           types.StringValue(object),
		ContextualTuples: contextualTuples,
		ContextModel:     *relationshiptuple.NewContextModel(context),
	}
}
