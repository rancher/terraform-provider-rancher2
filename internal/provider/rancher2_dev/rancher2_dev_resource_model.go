package rancher2_dev

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

type RancherDevResourceModel struct {
	Id               types.String                   `tfsdk:"id"`
	UserToken        types.String                   `tfsdk:"user_token"`
	BoolAttribute    types.Bool                     `tfsdk:"bool_attribute"`
	NumberAttribute  types.Number                   `tfsdk:"number_attribute"`
	Int64Attribute   types.Int64                    `tfsdk:"int64_attribute"`
	Int32Attribute   types.Int32                    `tfsdk:"int32_attribute"`
	Float64Attribute types.Float64                  `tfsdk:"float64_attribute"`
	Float32Attribute types.Float32                  `tfsdk:"float32_attribute"`
	StringAttribute  types.String                   `tfsdk:"string_attribute"`
	ListAttribute    types.List                     `tfsdk:"list_attribute"`
	SetAttribute     types.Set                      `tfsdk:"set_attribute"`
	MapAttribute     types.Map                      `tfsdk:"map_attribute"`
	NestedObject     NestedResourceModel            `tfsdk:"nested_object"`
	NestedObjectList []NestedResourceModel          `tfsdk:"nested_object_list"`
	NestedObjectMap  map[string]NestedResourceModel `tfsdk:"nested_object_map"`
}

// ToPlan returns a tfsdk.Plan with the data from the model.
func (m *RancherDevResourceModel) ToPlan(ctx context.Context, diags *diag.Diagnostics) tfsdk.Plan {
	if diags.HasError() {
		return tfsdk.Plan{}
	}
	// diags.AddWarning("Model given ToPlan:", fmt.Sprintf("%#v", m))
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
	if diags.HasError() {
		return tfsdk.State{}
	}
	diags.AddWarning("Model given ToState:", pp.PrettyPrint(m))
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
	obj := RancherDevModel{
		Id:               data.Id.ValueString(),
		UserToken:        data.UserToken.ValueString(),
		BoolAttribute:    data.BoolAttribute.ValueBool(),
		Int32Attribute:   data.Int32Attribute.ValueInt32(),
		Int64Attribute:   data.Int64Attribute.ValueInt64(),
		Float64Attribute: data.Float64Attribute.ValueFloat64(),
		Float32Attribute: data.Float32Attribute.ValueFloat32(),
		StringAttribute:  data.StringAttribute.ValueString(),
		NumberAttribute:  data.NumberAttribute.ValueBigFloat(),
	}

	if len(data.MapAttribute.Elements()) > 0 {
		obj.MapAttribute = map[string]string{}
		data.MapAttribute.ElementsAs(ctx, &obj.MapAttribute, false)
	}

	if len(data.ListAttribute.Elements()) > 0 {
		obj.ListAttribute = []string{}
		data.ListAttribute.ElementsAs(ctx, &obj.ListAttribute, false)
	}

	var setAttrs []string
	if len(data.SetAttribute.Elements()) > 0 {
		obj.SetAttribute = map[string]bool{}
		data.SetAttribute.ElementsAs(ctx, &setAttrs, false)
		for _, v := range setAttrs {
			obj.SetAttribute[v] = true
		}
	}

	data.NestedObject.ToGoModel(ctx, &obj.NestedObject)

	if len(data.NestedObjectList) > 0 {
		obj.NestedObjectList = []NestedObject{}
		for _, nestedObjectFromList := range data.NestedObjectList {
			m := NestedObject{}
			nestedObjectFromList.ToGoModel(ctx, &m)
			obj.NestedObjectList = append(obj.NestedObjectList, m)
		}
	}

	if len(data.NestedObjectMap) > 0 {
		obj.NestedObjectMap = map[string]NestedObject{}
		for k, v := range data.NestedObjectMap {
			m := NestedObject{}
			v.ToGoModel(ctx, &m)
			obj.NestedObjectMap[k] = m
		}
	}

	return obj
}

type NestedResourceModel struct {
	StringAttribute    types.String              `tfsdk:"string_attribute"`
	NestedNestedObject NestedNestedResourceModel `tfsdk:"nested_nested_object"`
}

// Fills the target with the values set in the current NestedResourceModel
func (r *NestedResourceModel) ToGoModel(ctx context.Context, target *NestedObject) diag.Diagnostics {
	dgs := diag.Diagnostics{}

	if target == nil {
		dgs.AddError("target cannot be nil", "")
		return dgs
	}
	if r.StringAttribute.IsNull() || r.StringAttribute.IsUnknown() {
		target.StringAttribute = ""
	}
	target.StringAttribute = r.StringAttribute.ValueString()

	diags := r.NestedNestedObject.ToGoModel(ctx, &target.NestedNestedObject)
	dgs.Append(diags...)
	return dgs
}

type NestedNestedResourceModel struct {
	StringAttribute types.String `tfsdk:"string_attribute"`
	BoolAttribute   types.Bool   `tfsdk:"bool_attribute"`
}

// Fills the target with the values set in the current NestedNestedResourceModel
func (r *NestedNestedResourceModel) ToGoModel(ctx context.Context, target *NestedNestedObject) diag.Diagnostics {
	dgs := diag.Diagnostics{}

	if target == nil {
		dgs.AddError("target cannot be nil", "")
		return dgs
	}
	if r.StringAttribute.IsNull() || r.StringAttribute.IsUnknown() {
		target.StringAttribute = ""
	}
	target.StringAttribute = r.StringAttribute.ValueString()
	if r.BoolAttribute.IsNull() || r.BoolAttribute.IsUnknown() {
		target.BoolAttribute = false
	}
	target.BoolAttribute = r.BoolAttribute.ValueBool()

	return dgs
}
