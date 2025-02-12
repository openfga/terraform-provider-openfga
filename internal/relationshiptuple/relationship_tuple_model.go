package relationshiptuple

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RelationshipTupleModel struct {
	StoreId  types.String `tfsdk:"store_id"`
	User     types.String `tfsdk:"user"`
	Relation types.String `tfsdk:"relation"`
	Object   types.String `tfsdk:"object"`
}
