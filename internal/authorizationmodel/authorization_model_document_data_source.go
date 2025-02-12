package authorizationmodel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"
	"github.com/openfga/language/pkg/go/transformer"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AuthorizationModelDocumentDataSource{}
var _ datasource.DataSourceWithConfigure = &AuthorizationModelDocumentDataSource{}

func NewAuthorizationModelDocumentDataSource() datasource.DataSource {
	return &AuthorizationModelDocumentDataSource{}
}

type AuthorizationModelDocumentDataSource struct {
	client *client.OpenFgaClient
}

type AuthorizationModelDocumentDataSourceModel struct {
	Dsl  types.String `tfsdk:"dsl"`
	Json types.String `tfsdk:"json"`
	// Model *MyCustomModel       `tfsdk:"model"` // TODO: Type

	Result types.String `tfsdk:"result"`
}

func (d *AuthorizationModelDocumentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_model_document"
}

func (d *AuthorizationModelDocumentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "An authorization model document represents the configuration of an authorization model. It can be used to convert between DSL and JSON or to represent a JSON model in canonical form.",

		Attributes: map[string]schema.Attribute{
			"dsl": schema.StringAttribute{
				MarkdownDescription: "The authorization model in DSL format",
				Optional:            true,
			},
			"json": schema.StringAttribute{
				MarkdownDescription: "The authorization model in JSON format",
				Optional:            true,
			},
			// TODO: Fully represent the data type in TF
			// "model": schema.SingleNestedAttribute{
			// 	MarkdownDescription: "The authorization model as object",
			// 	Optional:            true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"schema_version": schema.StringAttribute{
			// 			MarkdownDescription: "The schema's version",
			// 			Required:            true,
			// 		},
			// 		"type_definitions": schema.ListNestedAttribute{
			// 			MarkdownDescription: "The schema's types",
			// 			Required:            true,
			// 			NestedObject: schema.NestedAttributeObject{
			// 				Attributes: map[string]schema.Attribute{
			// 					"type": schema.StringAttribute{
			// 						MarkdownDescription: "The type's name",
			// 						Required:            true,
			// 						Validators: []validator.String{
			// 							stringvalidator.LengthAtLeast(1),
			// 						},
			// 					},
			// 					"relations": schema.MapNestedAttribute{
			// 						MarkdownDescription: "The type's relations",
			// 						Optional:            true,
			// 						NestedObject: schema.NestedAttributeObject{
			// 							Attributes: map[string]schema.Attribute{
			// 								// TODO
			// 							},
			// 						},
			// 					},
			// 					"metadata": schema.MapNestedAttribute{
			// 						MarkdownDescription: "The type's metadata",
			// 						Optional:            true,
			// 						NestedObject: schema.NestedAttributeObject{
			// 							Attributes: map[string]schema.Attribute{
			// 								"relations": schema.MapNestedAttribute{
			// 									MarkdownDescription: "The relations metadata",
			// 									Optional:            true,
			// 									NestedObject: schema.NestedAttributeObject{
			// 										Attributes: map[string]schema.Attribute{
			// 											"directly_related_usetypes": schema.ListNestedAttribute{
			// 												MarkdownDescription: "List of related user types",
			// 												Optional:            true,
			// 												NestedObject: schema.NestedAttributeObject{
			// 													Attributes: map[string]schema.Attribute{
			// 														"type": schema.StringAttribute{
			// 															MarkdownDescription: "The type's name",
			// 															Required:            true,
			// 															Validators: []validator.String{
			// 																stringvalidator.LengthAtLeast(1),
			// 															},
			// 														},
			// 														"relation": schema.StringAttribute{
			// 															MarkdownDescription: "The relation's name",
			// 															Optional:            true,
			// 															Validators: []validator.String{
			// 																stringvalidator.LengthAtLeast(1),
			// 															},
			// 														},
			// 														"wildcard": schema.SingleNestedAttribute{
			// 															MarkdownDescription: "Set, if the relationship is based on a wildcard",
			// 															Optional:            true,
			// 															Attributes:          map[string]schema.Attribute{},
			// 														},
			// 														"condition": schema.StringAttribute{
			// 															MarkdownDescription: "The name of a condition that is enforced over the allowed relation",
			// 															Optional:            true,
			// 														},
			// 													},
			// 												},
			// 											},
			// 											"module": schema.StringAttribute{
			// 												MarkdownDescription: "The module name",
			// 												Optional:            true,
			// 												Validators: []validator.String{
			// 													stringvalidator.LengthAtLeast(1),
			// 												},
			// 											},
			// 											"source_info": schema.SingleNestedAttribute{
			// 												MarkdownDescription: "The source information",
			// 												Optional:            true,
			// 												Attributes: map[string]schema.Attribute{
			// 													"file": schema.StringAttribute{
			// 														MarkdownDescription: "The file source",
			// 														Optional:            true,
			// 														Validators: []validator.String{
			// 															stringvalidator.LengthAtLeast(1),
			// 														},
			// 													},
			// 												},
			// 											},
			// 										},
			// 									},
			// 								},
			// 								"module": schema.StringAttribute{
			// 									MarkdownDescription: "The module name",
			// 									Optional:            true,
			// 									Validators: []validator.String{
			// 										stringvalidator.LengthAtLeast(1),
			// 									},
			// 								},
			// 								"source_info": schema.SingleNestedAttribute{
			// 									MarkdownDescription: "The source information",
			// 									Optional:            true,
			// 									Attributes: map[string]schema.Attribute{
			// 										"file": schema.StringAttribute{
			// 											MarkdownDescription: "The file source",
			// 											Optional:            true,
			// 											Validators: []validator.String{
			// 												stringvalidator.LengthAtLeast(1),
			// 											},
			// 										},
			// 									},
			// 								},
			// 							},
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 		"conditions": schema.MapNestedAttribute{
			// 			MarkdownDescription: "The schema's conditions",
			// 			Optional:            true,
			// 			NestedObject: schema.NestedAttributeObject{
			// 				Attributes: map[string]schema.Attribute{
			// 					"name": schema.StringAttribute{
			// 						MarkdownDescription: "The condition name",
			// 						Required:            true,
			// 						Validators: []validator.String{
			// 							stringvalidator.LengthAtLeast(1),
			// 						},
			// 					},
			// 					"expression": schema.StringAttribute{
			// 						MarkdownDescription: "The condition expression (in Google CEL format)",
			// 						Required:            true,
			// 						Validators: []validator.String{
			// 							stringvalidator.LengthAtLeast(1),
			// 						},
			// 					},
			// 					"parameters": schema.MapNestedAttribute{
			// 						MarkdownDescription: "The condition parameters",
			// 						Required:            true,
			// 						NestedObject: schema.NestedAttributeObject{
			// 							Attributes: map[string]schema.Attribute{
			// 								"type_name": schema.StringAttribute{
			// 									MarkdownDescription: "The parameter's type",
			// 									Required:            true,
			// 									Validators: []validator.String{
			// 										stringvalidator.OneOf(
			// 											"type_name_int",
			// 											"type_name_uint",
			// 											"type_name_double",
			// 											"type_name_bool",
			// 											"type_name_bytes",
			// 											"type_name_string",
			// 											"type_name_duration",
			// 											"type_name_timestamp",
			// 											"type_name_any",
			// 											"type_name_list",
			// 											"type_name_map",
			// 											"type_name_ipaddress",
			// 										),
			// 									},
			// 								},
			// 								"generic_types": schema.ListAttribute{
			// 									MarkdownDescription: "The parameter's generic types",
			// 									ElementType:         types.StringType,
			// 									Optional:            true,
			// 									Validators: []validator.List{
			// 										listvalidator.ValueStringsAre(
			// 											stringvalidator.OneOf(
			// 												"type_name_int",
			// 												"type_name_uint",
			// 												"type_name_double",
			// 												"type_name_bool",
			// 												"type_name_bytes",
			// 												"type_name_string",
			// 												"type_name_duration",
			// 												"type_name_timestamp",
			// 												"type_name_any",
			// 												"type_name_list",
			// 												"type_name_map",
			// 												"type_name_ipaddress",
			// 											),
			// 										),
			// 									},
			// 								},
			// 							},
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			"result": schema.StringAttribute{
				MarkdownDescription: "The resulting model in JSON format",
				Computed:            true,
			},
		},
	}
}

func (p AuthorizationModelDocumentDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("dsl"),
			path.MatchRoot("json"),
			// path.MatchRoot("model"), // TODO
		),
	}
}

func (d *AuthorizationModelDocumentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.OpenFgaClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.OpenFgaClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *AuthorizationModelDocumentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AuthorizationModelDocumentDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// model := state.Model
	jsonString := state.Json.ValueStringPointer()
	dslString := state.Dsl.ValueStringPointer()

	// if model != nil {
	// 	result, err := json.Marshal(state.Model)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to transform model into JSON, got error: %s", err))
	// 		return
	// 	}

	// 	jsonString = openfga.PtrString(string(result))
	// }

	if jsonString != nil {
		result, err := transformer.TransformJSONStringToDSL(*jsonString)
		if err != nil {
			resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable to transform JSON into DSL, got error: %s", err))
			return
		}

		dslString = result
	}

	if dslString == nil {
		resp.Diagnostics.AddError("Input Error", "DSL is undefined")
		return
	}

	// Transform DSL to canonical JSON form
	result, err := transformer.TransformDSLToJSON(*dslString)
	if err != nil {
		resp.Diagnostics.AddError("Input Error", fmt.Sprintf("Unable transform DSL into JSON, got error: %s", err))
		return
	}

	state.Result = types.StringValue(result)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
