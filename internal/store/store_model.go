package store

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StoreModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
