package rancher2_dev

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

type RancherDevModel struct {
  Id               string                  `json:"id"`
	StringAttribute  string                  `json:"string_attribute,omitempty"`
  NumberAttribute  *big.Float              `json:"number_attribute,omitempty"`
	Int32Attribute   int32                   `json:"int32_attribute,omitempty"`
  BoolAttribute    bool                    `json:"bool_attribute,omitempty"`
	Int64Attribute   int64                   `json:"int64_attribute,omitempty"`
	Float64Attribute float64                 `json:"float64_attribute,omitempty"`
	Float32Attribute float32                 `json:"float32_attribute,omitempty"`
	ListAttribute    []string                `json:"list_attribute,omitempty"`
	SetAttribute     map[string]bool         `json:"set_attribute,omitempty"`
	MapAttribute     map[string]string       `json:"map_attribute,omitempty"`
	NestedObject     NestedObject            `json:"nested_object,omitempty"`
	NestedObjectList []NestedObject          `json:"nested_object_list,omitempty"`
	NestedObjectMap  map[string]NestedObject `json:"nested_object_map,omitempty"`
}

type NestedObject struct {
	StringAttribute    string             `json:"string_attribute,omitempty"`
	NestedNestedObject NestedNestedObject `json:"nested_nested_object,omitempty"`
}

type NestedNestedObject struct {
	StringAttribute string `json:"string_attribute,omitempty"`
	BoolAttribute   bool   `json:"bool_attribute,omitempty"`
}

// ToResource converts a RancherDevModel to a RancherDevResourceModel
//
// This is useful for processing json marshalled response bodies.
func (obj *RancherDevModel) ToResourceModel(ctx context.Context, diags *diag.Diagnostics) *RancherDevResourceModel {
  tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevModel to RancherDevResourceModel: %+v", pp.PrettyPrint(obj)))
	if diags.HasError() {
		return nil
	}

	// var err error
	var data RancherDevResourceModel

	// primitive types (string, bool, int, etc)
	data.Id = types.StringValue(obj.Id)
	data.StringAttribute = types.StringValue(obj.StringAttribute)
  data.NumberAttribute = types.NumberValue(obj.NumberAttribute)

	data.BoolAttribute = types.BoolValue(obj.BoolAttribute)
	data.Int32Attribute = types.Int32Value(obj.Int32Attribute)
	data.Int64Attribute = types.Int64Value(obj.Int64Attribute)
	data.Float64Attribute = types.Float64Value(obj.Float64Attribute)
	data.Float32Attribute = types.Float32Value(obj.Float32Attribute)
	if diags.HasError() {
		return &data
	}

  // map
  mapElems := make(map[string]attr.Value)
	for k, v := range obj.MapAttribute {
		mapElems[k] = basetypes.NewStringValue(v)
	}
  mapVal, d := basetypes.NewMapValue(types.StringType, mapElems)
  diags.Append(d...)
  if mapVal.IsNull() {
    diags.AddError("Map Creation Error", "basetypes.NewMapValue returned null")
  }
  if diags.HasError() {
    return &data
  }
  data.MapAttribute = mapVal

	// list
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
		return &data
	}
	data.ListAttribute = listVal

  // set
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
		return &data
	}
	data.SetAttribute = setVal

  // complex types (nested objects)
  var nestedNestedObjectAttrTypes = map[string]attr.Type{
  	"string_attribute": types.StringType,
  	"bool_attribute":   types.BoolType,
  }
  var nestedObjectAttrTypes = map[string]attr.Type{
  	"string_attribute":     types.StringType,
  	"nested_nested_object": types.ObjectType{
  		AttrTypes: nestedNestedObjectAttrTypes,
  	},
  }

	rm := NestedResourceModel{}
	d = obj.NestedObject.ToResourceModel(ctx, &rm)
	diags.Append(d...)

	nest, d := basetypes.NewObjectValueFrom(
		ctx,
		nestedObjectAttrTypes,
		rm,
	)
	diags.Append(d...)
	data.NestedObject = nest

	if diags.HasError() {
		return &data
	}

	// list of nested objects
	nestedObjectList := make([]NestedResourceModel, 0, len(obj.NestedObjectList))
	for _, v := range obj.NestedObjectList {
		r := NestedResourceModel{}
		d := v.ToResourceModel(ctx, &r)
		diags.Append(d...)
		if diags.HasError() {
			return &data
		}
		nestedObjectList = append(nestedObjectList, r)
	}
	nestl, d := basetypes.NewListValueFrom(
		ctx,
		types.ObjectType{
			AttrTypes: nestedObjectAttrTypes,
		},
		nestedObjectList,
	)
	diags.Append(d...)
	if diags.HasError() {
		return &data
	}
	data.NestedObjectList = nestl

  // map of nested objects
	nestedObjectMap := make(map[string]NestedResourceModel, len(obj.NestedObjectMap))
	for k, v := range obj.NestedObjectMap {
		r := NestedResourceModel{}
		d := v.ToResourceModel(ctx, &r)
		diags.Append(d...)
		if diags.HasError() {
			return &data
		}
		nestedObjectMap[k] = r
	}
	nestm, d := basetypes.NewMapValueFrom(
		ctx,
		types.ObjectType{AttrTypes: nestedObjectAttrTypes},
		nestedObjectMap,
	)
	diags.Append(d...)
	if diags.HasError() {
		return &data
	}
	data.NestedObjectMap = nestm

  err := validateData(&data)
	if err != nil {
		diags.AddError("Error validating data: ", err.Error())
	}

  tflog.Debug(ctx, fmt.Sprintf("Converted RancherDevModel to RancherDevResourceModel: %+v", pp.PrettyPrint(data)))
	return &data
}

// Fills the target with types appropriate for a resource model.
func (m *NestedObject) ToResourceModel(ctx context.Context, target *NestedResourceModel) diag.Diagnostics {
  tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevModel NestedObject to RancherDevResourceModel NestedResourceModel: %+v", pp.PrettyPrint(m)))
	dgs := diag.Diagnostics{}

  // string attribute required
  if target == nil {
		dgs.AddError("target cannot be nil", "")
		return dgs
	}
	target.StringAttribute = types.StringValue(m.StringAttribute)

  var nestedNestedObjectAttrTypes = map[string]attr.Type{
  	"string_attribute": types.StringType,
  	"bool_attribute":   types.BoolType,
  }

  nrm := NestedNestedResourceModel{}
	diags := m.NestedNestedObject.ToResourceModel(ctx, &nrm)
	dgs.Append(diags...)
	if dgs.HasError() {
		return dgs
	}

	objValue, diags := basetypes.NewObjectValueFrom(ctx, nestedNestedObjectAttrTypes, nrm)
	dgs.Append(diags...)
	if dgs.HasError() {
		return dgs
	}
	target.NestedNestedResourceModel = objValue

  tflog.Debug(ctx, fmt.Sprintf("Converted RancherDevModel NestedObject to RancherDevResourceModel NestedResourceModel: %+v", pp.PrettyPrint(target)))
	return dgs
}

// Fills the target with types appropriate for a resource model.
func (m *NestedNestedObject) ToResourceModel(ctx context.Context, target *NestedNestedResourceModel) diag.Diagnostics {
  tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevModel NestedNestedObject to RancherDevResourceModel NestedNestedResourceModel: %+v", pp.PrettyPrint(m)))
  dgs := diag.Diagnostics{}
  if target == nil {
		dgs.AddError("target cannot be nil", "")
		return dgs
	}
  target.StringAttribute = types.StringValue(m.StringAttribute)
	target.BoolAttribute = types.BoolValue(m.BoolAttribute)
  tflog.Debug(ctx, fmt.Sprintf("Converted RancherDevModel NestedNestedObject to RancherDevResourceModel NestedNestedResourceModel: %+v", pp.PrettyPrint(target)))
	return dgs
}
