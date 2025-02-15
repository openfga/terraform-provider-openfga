package relationshiptuple

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

type ContextModel struct {
	ContextJson jsontypes.Normalized `tfsdk:"context_json"`
}

func (model ContextModel) GetContextJson() string {
	return model.ContextJson.ValueString()
}

func (model ContextModel) GetContextMap() (*map[string]interface{}, error) {
	if model.ContextJson.ValueString() == "" {
		return nil, nil
	}

	var context map[string]interface{}
	err := json.Unmarshal([]byte(model.ContextJson.ValueString()), &context)
	if err != nil {
		return nil, fmt.Errorf("failed to parse context JSON, got error: %s", err)
	}

	return &context, nil
}

func NewContextModel(data *map[string]interface{}) *ContextModel {
	context := jsontypes.NewNormalizedNull()

	jsonBytes, err := json.Marshal(data)
	if err == nil {
		context = jsontypes.NewNormalizedValue(string(jsonBytes))
	}

	return &ContextModel{
		ContextJson: context,
	}
}
