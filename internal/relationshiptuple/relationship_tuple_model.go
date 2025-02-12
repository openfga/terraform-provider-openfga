package relationshiptuple

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RelationshipTupleModel struct {
	StoreId   types.String                `tfsdk:"store_id"`
	User      types.String                `tfsdk:"user"`
	Relation  types.String                `tfsdk:"relation"`
	Object    types.String                `tfsdk:"object"`
	Condition *RelationshipTupleCondition `tfsdk:"condition"`
}

type RelationshipTupleCondition struct {
	Name    types.String         `tfsdk:"name"`
	Context jsontypes.Normalized `tfsdk:"context_json"`
}
