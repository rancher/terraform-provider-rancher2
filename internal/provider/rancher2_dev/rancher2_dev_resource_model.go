package rancher2_dev

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type RancherDevResourceModel struct {
  Id               types.String  `tfsdk:"id"`
	NumberAttribute  types.Number  `tfsdk:"number_attribute"`
	StringAttribute  types.String  `tfsdk:"string_attribute"`
  Int32Attribute   types.Int32   `tfsdk:"int32_attribute"`
  BoolAttribute    types.Bool    `tfsdk:"bool_attribute"`
	Int64Attribute   types.Int64   `tfsdk:"int64_attribute"`
	Float64Attribute types.Float64 `tfsdk:"float64_attribute"`
	Float32Attribute types.Float32 `tfsdk:"float32_attribute"`
	ListAttribute    types.List    `tfsdk:"list_attribute"`
	SetAttribute     types.Set     `tfsdk:"set_attribute"`
	MapAttribute     types.Map     `tfsdk:"map_attribute"`
	NestedObject     types.Object  `tfsdk:"nested_object"`
	NestedObjectList types.List    `tfsdk:"nested_object_list"`
	NestedObjectMap  types.Map     `tfsdk:"nested_object_map"`
}
type NestedResourceModel struct {
	StringAttribute           types.String  `tfsdk:"string_attribute"`
	NestedNestedResourceModel types.Object `tfsdk:"nested_nested_object"`
}
type NestedNestedResourceModel struct {
	StringAttribute types.String `tfsdk:"string_attribute"`
	BoolAttribute   types.Bool  `tfsdk:"bool_attribute"`
}

// ToPlan returns a tfsdk.Plan with the data from the model.
func (m *RancherDevResourceModel) ToPlan(ctx context.Context, diags *diag.Diagnostics) tfsdk.Plan {
  tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel to tfsdk.Plan: %+v", pp.PrettyPrint(m)))
	if diags.HasError() {
		return tfsdk.Plan{}
	}
	r := NewRancherDevResource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	plan := tfsdk.Plan{
		Schema: s.Schema,
	}
	if diags.HasError() {
		return plan
	}
	dgs := plan.Set(ctx, m)
	if dgs.HasError() {
		diags.Append(dgs...)
	}
	return plan
}

// ToState returns a tfsdk.State with the data from the model.
func (m *RancherDevResourceModel) ToState(ctx context.Context, diags *diag.Diagnostics) tfsdk.State {
  tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel to tfsdk.State: %+v", pp.PrettyPrint(m)))
  if diags.HasError() {
		return tfsdk.State{}
	}
	r := NewRancherDevResource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	state := tfsdk.State{
		Schema: s.Schema,
	}
	dgs := state.Set(ctx, m)
	if dgs.HasError() {
		diags.Append(dgs...)
	}
	return state
}

// ToGoModel converts a RancherDevResourceModel to a RancherDevModel
//
// This is helpful when building a json representation of the model for API requests.
func (data *RancherDevResourceModel) ToGoModel(ctx context.Context) RancherDevModel {
  tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel to RancherDevModel: %+v", pp.PrettyPrint(data)))
	obj := RancherDevModel{
		// Adding required attributes by default because we know they can't be null or empty
		Id:              data.Id.ValueString(),
		StringAttribute: data.StringAttribute.ValueString(),
    NumberAttribute: data.NumberAttribute.ValueBigFloat(),
	}

  // Add read only attributes as well (ToGoModel can be used to convert state, not just plan and config)
  if !data.Int32Attribute.IsNull() && !data.Int32Attribute.IsUnknown() {
  	obj.Int32Attribute = data.Int32Attribute.ValueInt32()
  }
  if !data.BoolAttribute.IsNull() && !data.BoolAttribute.IsUnknown() {
		obj.BoolAttribute = data.BoolAttribute.ValueBool()
	}
	if !data.Int64Attribute.IsNull() && !data.Int64Attribute.IsUnknown() {
		obj.Int64Attribute = data.Int64Attribute.ValueInt64()
	}
	if !data.Float64Attribute.IsNull() && !data.Float64Attribute.IsUnknown() {
		obj.Float64Attribute = data.Float64Attribute.ValueFloat64()
	}
	if !data.Float32Attribute.IsNull() && !data.Float32Attribute.IsUnknown() {
		obj.Float32Attribute = data.Float32Attribute.ValueFloat32()
	}
  if !data.MapAttribute.IsNull() && !data.MapAttribute.IsUnknown() {
  	if len(data.MapAttribute.Elements()) > 0 {
  		obj.MapAttribute = map[string]string{}
  		data.MapAttribute.ElementsAs(ctx, &obj.MapAttribute, false)
  	}
  }
  if !data.ListAttribute.IsNull() && !data.ListAttribute.IsUnknown() {
    if len(data.ListAttribute.Elements()) > 0 {
  		obj.ListAttribute = []string{}
  		data.ListAttribute.ElementsAs(ctx, &obj.ListAttribute, false)
  	}
  }

  if !data.SetAttribute.IsNull() && !data.SetAttribute.IsUnknown() {
  	if len(data.SetAttribute.Elements()) > 0 {
      var setAttrs []string
  		sa := map[string]bool{}
  		data.SetAttribute.ElementsAs(ctx, &setAttrs, false)
  		for _, v := range setAttrs {
  			sa[v] = true
  		}
      obj.SetAttribute = sa
  	}
  }

	if !data.NestedObject.IsNull() && !data.NestedObject.IsUnknown() {
		var nestedObject NestedResourceModel
		data.NestedObject.As(ctx, &nestedObject, basetypes.ObjectAsOptions{})
		obj.NestedObject = NestedObject{}
		nestedObject.ToGoModel(ctx, &obj.NestedObject)
	}

	if !data.NestedObjectList.IsNull() && !data.NestedObjectList.IsUnknown() {
		var nestedObjectList []NestedResourceModel
		data.NestedObjectList.ElementsAs(ctx, &nestedObjectList, false)
    nested := make([]NestedObject, 0, len(nestedObjectList))
		for _, v := range nestedObjectList {
			nestedObject := NestedObject{}
			v.ToGoModel(ctx, &nestedObject)
			nested = append(nested, nestedObject)
		}
    obj.NestedObjectList = nested
	}

	if !data.NestedObjectMap.IsNull() && !data.NestedObjectMap.IsUnknown() {
		var nestedObjectMap map[string]NestedResourceModel
		data.NestedObjectMap.ElementsAs(ctx, &nestedObjectMap, false)
    nested := make(map[string]NestedObject, len(nestedObjectMap))
		for k, v := range nestedObjectMap {
			nestedObject := NestedObject{}
			v.ToGoModel(ctx, &nestedObject)
			nested[k] = nestedObject
		}
    obj.NestedObjectMap = nested
	}

	return obj
}


// Fills the target with the values set in the current NestedResourceModel.
func (r *NestedResourceModel) ToGoModel(ctx context.Context, target *NestedObject) diag.Diagnostics {
  tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel NestedResourceModel to RancherDevModel NestedObject: %+v", pp.PrettyPrint(r)))
	dgs := diag.Diagnostics{}

	if target == nil {
		dgs.AddError("target cannot be nil", "")
		return dgs
	}

	if r.StringAttribute.IsNull() || r.StringAttribute.IsUnknown() {
    dgs.AddError("string_attribute is required", "")
    return dgs
	}

  target.StringAttribute = r.StringAttribute.ValueString()

	if !r.NestedNestedResourceModel.IsNull() && !r.NestedNestedResourceModel.IsUnknown() {
		var rm NestedNestedResourceModel
		r.NestedNestedResourceModel.As(ctx, &rm, basetypes.ObjectAsOptions{})
		target.NestedNestedObject = NestedNestedObject{}
		diags := rm.ToGoModel(ctx, &target.NestedNestedObject)
		dgs.Append(diags...)
	}

	return dgs
}

// Fills the target with the values set in the current NestedNestedResourceModel.
func (r *NestedNestedResourceModel) ToGoModel(ctx context.Context, target *NestedNestedObject) diag.Diagnostics {
  tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel NestedNestedResourceModel to RancherDevModel NestedNestedObject: %+v", pp.PrettyPrint(r)))
	dgs := diag.Diagnostics{}

	if target == nil {
		dgs.AddError("target cannot be nil", "")
		return dgs
	}
	if r.StringAttribute.IsNull() || r.StringAttribute.IsUnknown() {
    dgs.AddError("string_attribute is required", "")
    return dgs
	}
  target.StringAttribute = r.StringAttribute.ValueString()
  if !r.BoolAttribute.IsNull() && !r.BoolAttribute.IsUnknown() {
		target.BoolAttribute = r.BoolAttribute.ValueBool()
	}

	return dgs
}
