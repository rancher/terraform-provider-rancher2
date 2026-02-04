package rancher2_dev

import (
	"context"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type RancherDevModel struct {
	Id               string                  `json:"id"`
	BoolAttribute    bool                    `json:"bool_attribute,omitempty"`
	NumberAttribute  *big.Float              `json:"number_attribute,omitempty"`
	Int64Attribute   int64                   `json:"int64_attribute,omitempty"`
	Int32Attribute   int32                   `json:"int32_attribute,omitempty"`
	Float64Attribute float64                 `json:"float64_attribute,omitempty"`
	Float32Attribute float32                 `json:"float32_attribute,omitempty"`
	StringAttribute  string                  `json:"string_attribute,omitempty"`
	ListAttribute    []string                `json:"list_attribute,omitempty"`
	SetAttribute     map[string]bool         `json:"set_attribute,omitempty"`
	MapAttribute     map[string]string       `json:"map_attribute,omitempty"`
	NestedObject     NestedObject            `json:"nested_object"`
	NestedObjectList []NestedObject          `json:"nested_object_list,omitempty"`
	NestedObjectMap  map[string]NestedObject `json:"nested_object_map,omitempty"`
}

// ToResource converts a RancherDevModel to a RancherDevResourceModel
//
// This is useful for processing json marshalled response bodies.
func (obj *RancherDevModel) ToResourceModel(ctx context.Context, diags *diag.Diagnostics) *RancherDevResourceModel {
	if diags.HasError() {
		return nil
	}
	// var err error
	var data RancherDevResourceModel

	// primitive types (string, bool, int, etc)
	data.Id = types.StringValue(obj.Id)
	data.BoolAttribute = types.BoolValue(obj.BoolAttribute)
	data.Int32Attribute = types.Int32Value(obj.Int32Attribute)
	data.Int64Attribute = types.Int64Value(obj.Int64Attribute)
	data.Float64Attribute = types.Float64Value(obj.Float64Attribute)
	data.Float32Attribute = types.Float32Value(obj.Float32Attribute)
	data.StringAttribute = types.StringValue(obj.StringAttribute)
	data.NumberAttribute = types.NumberValue(obj.NumberAttribute)
	if diags.HasError() {
		return &data
	}

	// simple types (map, list, set)
	mapElems := make(map[string]attr.Value)
	if obj.MapAttribute != nil {
		for k, v := range obj.MapAttribute {
			mapElems[k] = basetypes.NewStringValue(v)
		}
	}
	mapVal, d := basetypes.NewMapValue(types.StringType, mapElems)
	diags.Append(d...)
	if mapVal.IsNull() {
		diags.AddError("Map Creation Error", "basetypes.NewMapValue returned null")
	}
	if diags.HasError() {
		diags.AddWarning("Progress:", "Error getting map from value.")
		return &data
	}
	data.MapAttribute = mapVal

	var listElems []attr.Value
	for _, v := range obj.ListAttribute {
		listElems = append(listElems, basetypes.NewStringValue(v))
	}
	listVal, d := basetypes.NewListValue(types.StringType, listElems)
	diags.Append(d...)
	if listVal.IsNull() {
		diags.AddError("List Creation Error", "basetypes.NewListValue returned null")
	}
	if diags.HasError() {
		diags.AddWarning("Progress:", "Error getting list from value.")
		return &data
	}
	data.ListAttribute = listVal

	var setAttributeElems []attr.Value
	for k := range obj.SetAttribute {
		setAttributeElems = append(setAttributeElems, basetypes.NewStringValue(k))
	}
	setVal, d := basetypes.NewSetValue(types.StringType, setAttributeElems)
	diags.Append(d...)
	if setVal.IsNull() {
		diags.AddError("Set Creation Error", "basetypes.NewSetValue returned null")
	}
	if diags.HasError() {
		diags.AddWarning("Progress:", "Error getting set from value.")
		return &data
	}
	data.SetAttribute = setVal

	// complex types (nested objects)
	data.NestedObject = NestedResourceModel{}
	d = obj.NestedObject.ToResourceModel(ctx, &data.NestedObject)
	diags.Append(d...)
	if diags.HasError() {
		diags.AddWarning("Progress:", "Error getting resource from nested object.")
		return &data
	}

	if len(obj.NestedObjectList) > 0 {
		for _, v := range obj.NestedObjectList {
			r := NestedResourceModel{}
			d = v.ToResourceModel(ctx, &r)
			diags.Append(d...)
			if diags.HasError() {
				diags.AddWarning("Progress:", "Error getting resource from list of nested objects.")
				return &data
			}
			data.NestedObjectList = append(data.NestedObjectList, r)
		}
	}

	if len(obj.NestedObjectMap) > 0 {
		data.NestedObjectMap = make(map[string]NestedResourceModel)

		for k, v := range obj.NestedObjectMap {
			r := NestedResourceModel{}
			d = v.ToResourceModel(ctx, &r)
			diags.Append(d...)
			if diags.HasError() {
				diags.AddWarning("Progress:", "Error getting resource from map of nested objects.")
				return &data
			}
			data.NestedObjectMap[k] = r
		}
	}

	err := validateData(&data)
	if err != nil {
		diags.AddError("Error validating data: ", err.Error())
	}

	return &data
}

type NestedObject struct {
	StringAttribute    string             `json:"string_attribute,omitempty"`
	NestedNestedObject NestedNestedObject `json:"nested_nested_object"`
}

// Fills the target with types appropriate for a resource model.
func (m *NestedObject) ToResourceModel(ctx context.Context, target *NestedResourceModel) diag.Diagnostics {
	dgs := diag.Diagnostics{}

	target.StringAttribute = types.StringValue(m.StringAttribute)
	diags := m.NestedNestedObject.ToResourceModel(ctx, &target.NestedNestedObject)
	dgs.Append(diags...)
	return dgs
}

type NestedNestedObject struct {
	StringAttribute string `json:"string_attribute,omitempty"`
	BoolAttribute   bool   `json:"bool_attribute,omitempty"`
}

// Fills the target with types appropriate for a resource model.
func (m *NestedNestedObject) ToResourceModel(ctx context.Context, target *NestedNestedResourceModel) diag.Diagnostics {
	dgs := diag.Diagnostics{}
	target.StringAttribute = types.StringValue(m.StringAttribute)
	target.BoolAttribute = types.BoolValue(m.BoolAttribute)
	return dgs
}
