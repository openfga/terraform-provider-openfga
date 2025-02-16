package authorizationmodel

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

const maxRecursionDepth = 5

type CustomAuthorizationModel struct {
	SchemaVersion   string                      `tfsdk:"schema_version" json:"schema_version"`
	TypeDefinitions []CustomTypeDefinition      `tfsdk:"type_definitions" json:"type_definitions"`
	Conditions      *map[string]CustomCondition `tfsdk:"conditions" json:"conditions"`
}

type CustomTypeDefinition struct {
	Type      string                    `tfsdk:"type" json:"type"`
	Relations *map[string]CustomUserset `tfsdk:"relations" json:"relations"`
	Metadata  *CustomMetadata           `tfsdk:"metadata" json:"metadata"`
}

type CustomMetadata struct {
	Relations  *map[string]CustomRelationMetadata `tfsdk:"relations" json:"relations"`
	Module     *string                            `tfsdk:"module" json:"module"`
	SourceInfo *CustomSourceInfo                  `tfsdk:"source_info" json:"source_info"`
}

type CustomRelationMetadata struct {
	DirectlyRelatedUserTypes *[]CustomRelationReference `tfsdk:"directly_related_user_types" json:"directly_related_user_types"`
	Module                   *string                    `tfsdk:"module" json:"module"`
	SourceInfo               *CustomSourceInfo          `tfsdk:"source_info" json:"source_info"`
}

type CustomRelationReference struct {
	Type      string       `tfsdk:"type" json:"type"`
	Relation  *string      `tfsdk:"relation" json:"relation"`
	Wildcard  *EmptyObject `tfsdk:"wildcard" json:"wildcard"`
	Condition *string      `tfsdk:"condition" json:"condition"`
}

type CustomUserset struct {
	This            *EmptyObject          `tfsdk:"this" json:"this"`
	ComputedUserset *CustomObjectRelation `tfsdk:"computed_userset" json:"computed_userset"`
	TupleToUserset  *CustomTupleToUserset `tfsdk:"tuple_to_userset" json:"tuple_to_userset"`
	Union           *CustomUserset        `tfsdk:"union" json:"union"`
	Intersection    *CustomUserset        `tfsdk:"intersection" json:"intersection"`
	Difference      *CustomDifference     `tfsdk:"difference" json:"difference"`
}

type EmptyObject struct {
}

type CustomObjectRelation struct {
	Object   *string `tfsdk:"object" json:"object"`
	Relation *string `tfsdk:"relation" json:"relation"`
}

type CustomTupleToUserset struct {
	Tupleset        CustomObjectRelation `tfsdk:"tupleset" json:"tupleset"`
	ComputedUserset CustomObjectRelation `tfsdk:"computed_userset" json:"computed_userset"`
}

type CustomDifference struct {
	Base     CustomUserset `tfsdk:"base" json:"base"`
	Subtract CustomUserset `tfsdk:"subtract" json:"subtract"`
}

type CustomCondition struct {
	Name       string                                  `tfsdk:"name" json:"name"`
	Expression string                                  `tfsdk:"expression" json:"expression"`
	Parameters *map[string]CustomConditionParamTypeRef `tfsdk:"parameters" json:"parameters"`
	Metadata   *CustomConditionMetadata                `tfsdk:"metadata" json:"metadata"`
}

type CustomConditionParamTypeRef struct {
	TypeName     string                         `tfsdk:"type_name" json:"type_name"`
	GenericTypes *[]CustomConditionParamTypeRef `tfsdk:"generic_types" json:"generic_types"`
}

type CustomConditionMetadata struct {
	Module     *string           `tfsdk:"module" json:"module"`
	SourceInfo *CustomSourceInfo `tfsdk:"source_info" json:"source_info"`
}

type CustomSourceInfo struct {
	File *string `tfsdk:"file" json:"file"`
}

func CustomAuthorizationModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"schema_version": schema.StringAttribute{
			Required: true,
		},
		"type_definitions": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomTypeDefinitionSchema(),
			},
			Required: true,
		},
		"conditions": schema.MapNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomConditionSchema(),
			},
			Optional: true,
		},
	}
}

func CustomTypeDefinitionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Required: true,
		},
		"relations": schema.MapNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomUsersetSchema(0),
			},
			Optional: true,
		},
		"metadata": schema.SingleNestedAttribute{
			Attributes: CustomMetadataSchema(),
			Optional:   true,
		},
	}
}

func CustomMetadataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"relations": schema.MapNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomRelationMetadataSchema(),
			},
			Optional: true,
		},
		"module": schema.StringAttribute{
			Optional: true,
		},
		"source_info": schema.SingleNestedAttribute{
			Attributes: CustomSourceInfoSchema(),
			Optional:   true,
		},
	}
}

func CustomRelationMetadataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"directly_related_user_types": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomRelationReferenceSchema(),
			},
			Optional: true,
		},
		"module": schema.StringAttribute{
			Optional: true,
		},
		"source_info": schema.SingleNestedAttribute{
			Attributes: CustomSourceInfoSchema(),
			Optional:   true,
		},
	}
}

func CustomRelationReferenceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Required: true,
		},
		"relation": schema.StringAttribute{
			Optional: true,
		},
		"wildcard": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{},
			Optional:   true,
		},
		"condition": schema.StringAttribute{
			Optional: true,
		},
	}
}

func CustomUsersetSchema(depth int) map[string]schema.Attribute {
	if depth >= maxRecursionDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		"this": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{},
			Optional:   true,
		},
		"computed_userset": schema.SingleNestedAttribute{
			Attributes: CustomObjectRelationSchema(),
			Optional:   true,
		},
		"tuple_to_userset": schema.SingleNestedAttribute{
			Attributes: CustomTupleToUsersetSchema(),
			Optional:   true,
		},
		"union": schema.SingleNestedAttribute{
			Attributes: CustomUsersetSchema(depth + 1),
			Optional:   true,
		},
		"intersection": schema.SingleNestedAttribute{
			Attributes: CustomUsersetSchema(depth + 1),
			Optional:   true,
		},
		"difference": schema.SingleNestedAttribute{
			Attributes: CustomDifferenceSchema(depth),
			Optional:   true,
		},
	}
}

func CustomObjectRelationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"object": schema.StringAttribute{
			Optional: true,
		},
		"relation": schema.StringAttribute{
			Optional: true,
		},
	}
}

func CustomTupleToUsersetSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"tupleset": schema.SingleNestedAttribute{
			Attributes: CustomObjectRelationSchema(),
			Required:   true,
		},
		"computed_userset": schema.SingleNestedAttribute{
			Attributes: CustomObjectRelationSchema(),
			Required:   true,
		},
	}
}

func CustomDifferenceSchema(depth int) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"base": schema.SingleNestedAttribute{
			Attributes: CustomUsersetSchema(depth + 1),
			Required:   true,
		},
		"subtract": schema.SingleNestedAttribute{
			Attributes: CustomUsersetSchema(depth + 1),
			Required:   true,
		},
	}
}

func CustomConditionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required: true,
		},
		"expression": schema.StringAttribute{
			Required: true,
		},
		"parameters": schema.MapNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomConditionParamTypeRefSchema(0),
			},
			Optional: true,
		},
		"metadata": schema.SingleNestedAttribute{
			Attributes: CustomConditionMetadataSchema(),
			Optional:   true,
		},
	}
}

func CustomConditionParamTypeRefSchema(depth int) map[string]schema.Attribute {
	if depth >= maxRecursionDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		"type_name": schema.StringAttribute{
			Required: true,
		},
		"generic_types": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: CustomConditionParamTypeRefSchema(depth + 1),
			},
			Optional: true,
		},
	}
}

func CustomConditionMetadataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"module": schema.StringAttribute{
			Optional: true,
		},
		"source_info": schema.SingleNestedAttribute{
			Attributes: CustomSourceInfoSchema(),
			Optional:   true,
		},
	}
}

func CustomSourceInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"file": schema.StringAttribute{
			Optional: true,
		},
	}
}
