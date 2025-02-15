package relationshiptuple

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	openfga "github.com/openfga/go-sdk"
)

type RelationshipConditionModel struct {
	Name types.String `tfsdk:"name"`
	ContextModel
}

func (model RelationshipConditionModel) GetName() string {
	return model.Name.ValueString()
}

func (model *RelationshipConditionModel) ToCondition() (*openfga.RelationshipCondition, error) {
	if model == nil {
		return nil, nil
	}

	context, err := model.GetContextMap()
	if err != nil {
		return nil, err
	}

	return &openfga.RelationshipCondition{
		Name:    model.GetName(),
		Context: context,
	}, nil
}

func NewRelationshipConditionModel(name string, context *map[string]interface{}) *RelationshipConditionModel {
	contextModel := NewContextModel(context)

	return &RelationshipConditionModel{
		Name:         types.StringValue(name),
		ContextModel: *contextModel,
	}
}

func NewRelationshipConditionModelFromCondition(condition openfga.RelationshipCondition) *RelationshipConditionModel {
	contextModel := NewContextModel(condition.Context)

	return &RelationshipConditionModel{
		Name:         types.StringValue(condition.GetName()),
		ContextModel: *contextModel,
	}
}
