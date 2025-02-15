package relationshiptuple

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	openfga "github.com/openfga/go-sdk"
)

type RelationshipTupleModel struct {
	User     types.String `tfsdk:"user"`
	Relation types.String `tfsdk:"relation"`
	Object   types.String `tfsdk:"object"`
}

func (model RelationshipTupleModel) GetUser() string {
	return model.User.ValueString()
}

func (model RelationshipTupleModel) GetRelation() string {
	return model.Relation.ValueString()
}

func (model RelationshipTupleModel) GetObject() string {
	return model.Object.ValueString()
}

func (model RelationshipTupleModel) ToTuple() *openfga.TupleKeyWithoutCondition {
	return &openfga.TupleKeyWithoutCondition{
		User:     model.GetUser(),
		Relation: model.GetRelation(),
		Object:   model.GetObject(),
	}
}

func NewRelationshipTupleModel(user string, relation string, object string) *RelationshipTupleModel {
	return &RelationshipTupleModel{
		User:     types.StringValue(user),
		Relation: types.StringValue(relation),
		Object:   types.StringValue(object),
	}
}

type TupleInterface interface {
	GetUser() string
	GetRelation() string
	GetObject() string
}

func NewRelationshipTupleModelFromTuple(tuple TupleInterface) *RelationshipTupleModel {
	return &RelationshipTupleModel{
		User:     types.StringValue(tuple.GetUser()),
		Relation: types.StringValue(tuple.GetRelation()),
		Object:   types.StringValue(tuple.GetObject()),
	}
}

type RelationshipTupleWithConditionModel struct {
	RelationshipTupleModel
	Condition *RelationshipConditionModel `tfsdk:"condition"`
}

func (model RelationshipTupleWithConditionModel) GetCondition() *RelationshipConditionModel {
	return model.Condition
}

func (model RelationshipTupleWithConditionModel) ToTupleWithCondition() (*openfga.TupleKey, error) {
	condition, err := model.GetCondition().ToCondition()
	if err != nil {
		return nil, err
	}

	return &openfga.TupleKey{
		User:      model.GetUser(),
		Relation:  model.GetRelation(),
		Object:    model.GetObject(),
		Condition: condition,
	}, nil
}

func NewRelationshipTupleWithConditionModel(user string, relation string, object string, condition *RelationshipConditionModel) *RelationshipTupleWithConditionModel {
	return &RelationshipTupleWithConditionModel{
		RelationshipTupleModel: *NewRelationshipTupleModel(user, relation, object),
		Condition:              condition,
	}
}

type TupleWithConditionInterface interface {
	TupleInterface
	HasCondition() bool
	GetCondition() openfga.RelationshipCondition
}

func NewRelationshipTupleWithConditionModelFromTuple(tuple TupleWithConditionInterface) *RelationshipTupleWithConditionModel {
	var condition *RelationshipConditionModel
	if tuple.HasCondition() {
		condition = NewRelationshipConditionModelFromCondition(tuple.GetCondition())
	}

	return &RelationshipTupleWithConditionModel{
		RelationshipTupleModel: *NewRelationshipTupleModelFromTuple(tuple),
		Condition:              condition,
	}
}
