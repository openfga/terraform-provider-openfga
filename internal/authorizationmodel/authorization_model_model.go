package authorizationmodel

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	openfga "github.com/openfga/go-sdk"
)

type AuthorizationModelModel struct {
	Id        types.String         `tfsdk:"id"`
	StoreId   types.String         `tfsdk:"store_id"`
	ModelJson jsontypes.Normalized `tfsdk:"model_json"`
}

type AuthorizationModelWithoutId struct {
	SchemaVersion   string                        `json:"schema_version" yaml:"schema_version"`
	TypeDefinitions []openfga.TypeDefinition      `json:"type_definitions" yaml:"type_definitions"`
	Conditions      *map[string]openfga.Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
}
