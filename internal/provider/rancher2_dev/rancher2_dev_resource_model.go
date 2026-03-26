package rancher2_dev

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

type RancherDevResourceModel struct {
	ID               types.String  `tfsdk:"id"`
	Identifier       types.String  `tfsdk:"identifier"`
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
	StringAttribute    types.String `tfsdk:"string_attribute"`
	NestedNestedObject types.Object `tfsdk:"nested_nested_object"`
}
type NestedNestedResourceModel struct {
	StringAttribute types.String `tfsdk:"string_attribute"`
	BoolAttribute   types.Bool   `tfsdk:"bool_attribute"`
}

// ToPlan returns a tfsdk.Plan with the data from the model.
func (m *RancherDevResourceModel) ToPlan(ctx context.Context, diags *diag.Diagnostics) tfsdk.Plan {
	tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel to tfsdk.Plan: \n%+v", pp.PrettyPrint(m)))
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
	tflog.Debug(ctx, fmt.Sprintf("Converted RancherDevResourceModel to tfsdk.Plan: \n%+v", pp.PrettyPrint(plan)))
	return plan
}

// ToState returns a tfsdk.State with the data from the model.
func (m *RancherDevResourceModel) ToState(ctx context.Context, diags *diag.Diagnostics) tfsdk.State {
	tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel to tfsdk.State: \n%+v", pp.PrettyPrint(m)))
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
	tflog.Debug(ctx, fmt.Sprintf("Converted RancherDevResourceModel to tfsdk.State: \n%+v", pp.PrettyPrint(state)))
	return state
}

// ToGoModel converts a RancherDevResourceModel to a RancherDevModel
//
// This is helpful when building a json representation of the model for API requests.
func (data *RancherDevResourceModel) ToGoModel(ctx context.Context) RancherDevModel {
	tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel to RancherDevModel: \n%+v", pp.PrettyPrint(data)))
	obj := RancherDevModel{}

	// NOTE: there is no point trying to leave out null/unknowns here, they will just be translated to Go zero values.
	// primitive
	obj.ID = data.ID.ValueString()                        // problem! How do we leave out values that shouldn't be going to the API?
	obj.Identifier = data.Identifier.ValueString()        // problem! How do we leave out values that shouldn't be going to the API?
	obj.Int32Attribute = data.Int32Attribute.ValueInt32() // problem! How do we leave out values that shouldn't be going to the API?
	obj.StringAttribute = data.StringAttribute.ValueString()
	obj.NumberAttribute = data.NumberAttribute.ValueBigFloat()
	obj.BoolAttribute = data.BoolAttribute.ValueBool()
	obj.Int64Attribute = data.Int64Attribute.ValueInt64()
	obj.Float64Attribute = data.Float64Attribute.ValueFloat64()
	obj.Float32Attribute = data.Float32Attribute.ValueFloat32()

	// simple
	// map
	obj.MapAttribute = map[string]string{}
	data.MapAttribute.ElementsAs(ctx, &obj.MapAttribute, false)
	// list
	obj.ListAttribute = []string{}
	data.ListAttribute.ElementsAs(ctx, &obj.ListAttribute, false)
	// set
	var setAttrs []string
	sa := map[string]bool{}
	data.SetAttribute.ElementsAs(ctx, &setAttrs, false)
	for _, v := range setAttrs {
		sa[v] = true
	}
	obj.SetAttribute = sa

	// complex types (nested objects)
	// nested object
	var nestedObject NestedResourceModel
	data.NestedObject.As(ctx, &nestedObject, basetypes.ObjectAsOptions{})
	obj.NestedObject = NestedObject{}
	nestedObject.ToGoModel(ctx, &obj.NestedObject)
	// list of nested objects
	var nestedObjectList []NestedResourceModel
	data.NestedObjectList.ElementsAs(ctx, &nestedObjectList, false)
	nestedList := make([]NestedObject, 0, len(nestedObjectList))
	for _, v := range nestedObjectList {
		nestedObject := NestedObject{}
		v.ToGoModel(ctx, &nestedObject)
		nestedList = append(nestedList, nestedObject)
	}
	obj.NestedObjectList = nestedList
	// map of nested objects
	var nestedObjectMap map[string]NestedResourceModel
	data.NestedObjectMap.ElementsAs(ctx, &nestedObjectMap, false)
	nestedMap := make(map[string]NestedObject, len(nestedObjectMap))
	for k, v := range nestedObjectMap {
		nestedObject := NestedObject{}
		v.ToGoModel(ctx, &nestedObject)
		nestedMap[k] = nestedObject
	}
	obj.NestedObjectMap = nestedMap

	tflog.Debug(ctx, fmt.Sprintf("Converted RancherDevResourceModel to RancherDevModel: \n%+v", pp.PrettyPrint(obj)))
	return obj
}

// Fills the target with the values set in the current NestedResourceModel.
func (r *NestedResourceModel) ToGoModel(ctx context.Context, target *NestedObject) diag.Diagnostics {
	tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel NestedResourceModel to RancherDevModel NestedObject: \n%+v", pp.PrettyPrint(r)))
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

	if !r.NestedNestedObject.IsNull() && !r.NestedNestedObject.IsUnknown() {
		var rm NestedNestedResourceModel
		r.NestedNestedObject.As(ctx, &rm, basetypes.ObjectAsOptions{})
		target.NestedNestedObject = NestedNestedObject{}
		diags := rm.ToGoModel(ctx, &target.NestedNestedObject)
		dgs.Append(diags...)
	}
	tflog.Debug(ctx, fmt.Sprintf("Converted RancherDevResourceModel NestedResourceModel to RancherDevModel NestedObject: \n%+v", pp.PrettyPrint(target)))
	return dgs
}

// Fills the target with the values set in the current NestedNestedResourceModel.
func (r *NestedNestedResourceModel) ToGoModel(ctx context.Context, target *NestedNestedObject) diag.Diagnostics {
	tflog.Debug(ctx, fmt.Sprintf("Converting RancherDevResourceModel NestedNestedResourceModel to RancherDevModel NestedNestedObject: \n%+v", pp.PrettyPrint(r)))
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
	tflog.Debug(ctx, fmt.Sprintf("Converted RancherDevResourceModel NestedNestedResourceModel to RancherDevModel NestedNestedObject: \n%+v", pp.PrettyPrint(target)))
	return dgs
}
