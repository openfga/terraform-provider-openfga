package authorizationmodel

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	openfga "github.com/openfga/go-sdk"
)

type AuthorizationModelWithoutId struct {
	openfga.AuthorizationModel
	Id *interface{} `json:"-"`
}

func (o AuthorizationModelWithoutId) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["schema_version"] = o.SchemaVersion
	toSerialize["type_definitions"] = o.TypeDefinitions
	if o.Conditions != nil {
		toSerialize["conditions"] = o.Conditions
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(toSerialize)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type AuthorizationModelModel struct {
	Id        types.String         `tfsdk:"id"`
	ModelJson jsontypes.Normalized `tfsdk:"model_json"`
}

func (model AuthorizationModelModel) GetId() string {
	return model.Id.ValueString()
}

func (model AuthorizationModelModel) GetModelJson() string {
	return model.ModelJson.ValueString()
}

func (model AuthorizationModelModel) ToAuthorizationModel() (*openfga.AuthorizationModel, error) {
	var authorizationModel openfga.AuthorizationModel
	err := json.Unmarshal([]byte(model.GetModelJson()), &authorizationModel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model JSON, got error: %s", err)
	}

	authorizationModel.Id = model.GetId()

	return &authorizationModel, nil
}

func NewAuthorizationModelModel(id string) *AuthorizationModelModel {
	return &AuthorizationModelModel{
		Id:        types.StringValue(id),
		ModelJson: jsontypes.NewNormalizedNull(),
	}
}

func NewAuthorizationModelModelWithModelJson(id string, modelJson string) *AuthorizationModelModel {
	return &AuthorizationModelModel{
		Id:        types.StringValue(id),
		ModelJson: jsontypes.NewNormalizedValue(modelJson),
	}
}

func NewAuthorizationModelModelFromAuthorizationModel(authorizationModel openfga.AuthorizationModel) *AuthorizationModelModel {
	jsonBytes, _ := json.Marshal(AuthorizationModelWithoutId{AuthorizationModel: authorizationModel})

	return &AuthorizationModelModel{
		Id:        types.StringValue(authorizationModel.GetId()),
		ModelJson: jsontypes.NewNormalizedValue(string(jsonBytes)),
	}
}
