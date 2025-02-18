package query

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mauriceackel/terraform-provider-openfga/internal/provider/relationshiptuple"
)

type ListObjectsQueryModel struct {
	User             types.String                                             `tfsdk:"user"`
	Relation         types.String                                             `tfsdk:"relation"`
	Type             types.String                                             `tfsdk:"type"`
	ContextualTuples *[]relationshiptuple.RelationshipTupleWithConditionModel `tfsdk:"contextual_tuples"`
	relationshiptuple.ContextModel
}

func (query ListObjectsQueryModel) GetUser() string {
	return query.User.ValueString()
}

func (query ListObjectsQueryModel) GetRelation() string {
	return query.Relation.ValueString()
}

func (query ListObjectsQueryModel) GetType() string {
	return query.Type.ValueString()
}

func (query ListObjectsQueryModel) GetContextualTuples() []relationshiptuple.RelationshipTupleWithConditionModel {
	if query.ContextualTuples == nil {
		return []relationshiptuple.RelationshipTupleWithConditionModel{}
	}

	return *query.ContextualTuples
}

func NewListObjectsQueryModel(user string, relation string, type_ string, contextualTuples *[]relationshiptuple.RelationshipTupleWithConditionModel, context *map[string]interface{}) *ListObjectsQueryModel {
	return &ListObjectsQueryModel{
		User:             types.StringValue(user),
		Relation:         types.StringValue(relation),
		Type:             types.StringValue(type_),
		ContextualTuples: contextualTuples,
		ContextModel:     *relationshiptuple.NewContextModel(context),
	}
}
