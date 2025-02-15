package store

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	openfga "github.com/openfga/go-sdk"
)

type StoreModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (model StoreModel) GetId() string {
	return model.Id.ValueString()
}

func (model StoreModel) GetName() string {
	return model.Name.ValueString()
}

func (model StoreModel) ToStore() *openfga.Store {
	return &openfga.Store{
		Id:   model.GetId(),
		Name: model.GetName(),
	}
}

func NewStoreModel(id string, name string) *StoreModel {
	return &StoreModel{
		Id:   types.StringValue(id),
		Name: types.StringValue(name),
	}
}

type StoreInterface interface {
	GetId() string
	GetName() string
}

func NewStoreModelFromStore(store StoreInterface) *StoreModel {
	return &StoreModel{
		Id:   types.StringValue(store.GetId()),
		Name: types.StringValue(store.GetName()),
	}
}
