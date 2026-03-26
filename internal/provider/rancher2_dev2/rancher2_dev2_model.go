package rancher2_dev2

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

type Rancher2Dev2Model struct {
	ID         string   `tfsdk:"id"`
	APIVersion string   `tfsdk:"api_version"`
	Kind       string   `tfsdk:"kind"`
	Metadata   Metadata `tfsdk:"metadata"`
	Spec       Spec     `tfsdk:"spec"`
	Status     string   `tfsdk:"status"`
}

type OwnerReference struct {
	APIVersion         string `tfsdk:"api_version"`
	Kind               string `tfsdk:"kind"`
	Name               string `tfsdk:"name"`
	UID                string `tfsdk:"uid"`
	Controller         bool   `tfsdk:"controller"`
	BlockOwnerDeletion bool   `tfsdk:"block_owner_deletion"`
}

type Metadata struct {
	Name                      string            `tfsdk:"name"`
	Namespace                 string            `tfsdk:"namespace"`
	GenerateName              string            `tfsdk:"generate_name"`
	Annotations               map[string]string `tfsdk:"annotations"`
	Labels                    map[string]string `tfsdk:"labels"`
	Finalizers                []string          `tfsdk:"finalizers"`
	OwnerReferences           []OwnerReference  `tfsdk:"owner_references"`
	UID                       string            `tfsdk:"uid"`
	Generation                int64             `tfsdk:"generation"`
	CreationTimestamp         string            `tfsdk:"creation_timestamp"`
	DeletionGracePeriodSecond int64             `tfsdk:"deletion_grace_period_seconds"`
	DeletionTimestamp         string            `tfsdk:"deletion_timestamp"`
	ManagedFields             string            `tfsdk:"managed_fields"`
	ResourceVersion           string            `tfsdk:"resource_version"`
	SelfLink                  string            `tfsdk:"self_link"`
}

type Spec struct {
	String     string            `tfsdk:"string"`
	Bool       bool              `tfsdk:"bool"`
	Number     float64           `tfsdk:"number"`
	Int32      int32             `tfsdk:"int32"`
	Int64      int64             `tfsdk:"int64"`
	Float32    float32           `tfsdk:"float32"`
	Float64    float64           `tfsdk:"float64"`
	Map        map[string]string `tfsdk:"map"`
	List       []string          `tfsdk:"list"`
	Object     Object            `tfsdk:"object"`
	ObjectList []Object          `tfsdk:"object_list"`
	ObjectMap  map[string]Object `tfsdk:"object_map"`
}

type Object struct {
	StringAttribute string `tfsdk:"string_attribute"`
}

type Rancher2Dev2ResourceModel struct {
	ID         types.String `tfsdk:"id"`
	APIVersion types.String `tfsdk:"api_version"`
	Kind       types.String `tfsdk:"kind"`
	Metadata   types.Object `tfsdk:"metadata"`
	Spec       types.Object `tfsdk:"spec"`
	Status     types.String `tfsdk:"status"`
}

func (obj *Rancher2Dev2Model) ToResourceModel(ctx context.Context, diags *diag.Diagnostics) *Rancher2Dev2ResourceModel {
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2Model to Rancher2Dev2ResourceModel: %+v", pp.PrettyPrint(obj)))
	if diags.HasError() {
		return nil
	}

	var data Rancher2Dev2ResourceModel

	data.ID = types.StringValue(obj.ID)
	data.APIVersion = types.StringValue(obj.APIVersion)
	data.Kind = types.StringValue(obj.Kind)
	data.Status = types.StringValue(obj.Status)

	metadata, d := obj.Metadata.ToResourceModel(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil
	}
	data.Metadata = metadata

	spec, d := obj.Spec.ToResourceModel(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil
	}
	data.Spec = spec

	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2Model to Rancher2Dev2ResourceModel: %+v", pp.PrettyPrint(data)))
	return &data
}

func (m *Metadata) ToResourceModel(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"name":                          types.StringType,
		"namespace":                     types.StringType,
		"generate_name":                 types.StringType,
		"annotations":                   types.MapType{ElemType: types.StringType},
		"labels":                        types.MapType{ElemType: types.StringType},
		"finalizers":                    types.ListType{ElemType: types.StringType},
		"owner_references":              types.ListType{ElemType: types.ObjectType{AttrTypes: ownerReferenceAttrTypes}},
		"uid":                           types.StringType,
		"generation":                    types.Int64Type,
		"creation_timestamp":            types.StringType,
		"deletion_grace_period_seconds": types.Int64Type,
		"deletion_timestamp":            types.StringType,
		"managed_fields":                types.StringType,
		"resource_version":              types.StringType,
		"self_link":                     types.StringType,
	}

	var ownerReferencesValue attr.Value = types.ListNull(types.ObjectType{AttrTypes: ownerReferenceAttrTypes})
	if len(m.OwnerReferences) > 0 {
		var ownerReferences []attr.Value
		for _, or := range m.OwnerReferences {
			orObj, d := or.ToResourceModel(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return types.ObjectNull(attributeTypes), diags
			}
			ownerReferences = append(ownerReferences, orObj)
		}
		var d diag.Diagnostics
		ownerReferencesValue, d = basetypes.NewListValue(types.ObjectType{AttrTypes: ownerReferenceAttrTypes}, ownerReferences)
		diags.Append(d...)
	}

	var annotationsValue attr.Value = types.MapNull(types.StringType)
	if m.Annotations != nil {
		var annotations = make(map[string]attr.Value)
		for k, v := range m.Annotations {
			annotations[k] = basetypes.NewStringValue(v)
		}
		var d diag.Diagnostics
		annotationsValue, d = basetypes.NewMapValue(types.StringType, annotations)
		diags.Append(d...)
	}

	var labelsValue attr.Value = types.MapNull(types.StringType)
	if m.Labels != nil {
		var labels = make(map[string]attr.Value)
		for k, v := range m.Labels {
			labels[k] = basetypes.NewStringValue(v)
		}
		var d diag.Diagnostics
		labelsValue, d = basetypes.NewMapValue(types.StringType, labels)
		diags.Append(d...)
	}

	var finalizersValue attr.Value = types.ListNull(types.StringType)
	if len(m.Finalizers) > 0 {
		var finalizers []attr.Value
		for _, f := range m.Finalizers {
			finalizers = append(finalizers, basetypes.NewStringValue(f))
		}
		var d diag.Diagnostics
		finalizersValue, d = basetypes.NewListValue(types.StringType, finalizers)
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectNull(attributeTypes), diags
	}

	attributes := map[string]attr.Value{
		"name":                          basetypes.NewStringValue(m.Name),
		"namespace":                     basetypes.NewStringValue(m.Namespace),
		"generate_name":                 basetypes.NewStringValue(m.GenerateName),
		"annotations":                   annotationsValue,
		"labels":                        labelsValue,
		"finalizers":                    finalizersValue,
		"owner_references":              ownerReferencesValue,
		"uid":                           basetypes.NewStringValue(m.UID),
		"generation":                    basetypes.NewInt64Value(m.Generation),
		"creation_timestamp":            basetypes.NewStringValue(m.CreationTimestamp),
		"deletion_grace_period_seconds": basetypes.NewInt64Value(m.DeletionGracePeriodSecond),
		"deletion_timestamp":            basetypes.NewStringValue(m.DeletionTimestamp),
		"managed_fields":                basetypes.NewStringValue(m.ManagedFields),
		"resource_version":              basetypes.NewStringValue(m.ResourceVersion),
		"self_link":                     basetypes.NewStringValue(m.SelfLink),
	}

	return basetypes.NewObjectValue(attributeTypes, attributes)
}

var ownerReferenceAttrTypes = map[string]attr.Type{
	"api_version":          types.StringType,
	"kind":                 types.StringType,
	"name":                 types.StringType,
	"uid":                  types.StringType,
	"controller":           types.BoolType,
	"block_owner_deletion": types.BoolType,
}

func (m *OwnerReference) ToResourceModel(ctx context.Context) (types.Object, diag.Diagnostics) {
	return basetypes.NewObjectValueFrom(ctx, ownerReferenceAttrTypes, m)
}

func (m *Spec) ToResourceModel(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"string":      types.StringType,
		"bool":        types.BoolType,
		"number":      types.NumberType,
		"int32":       types.Int64Type,
		"int64":       types.Int64Type,
		"float32":     types.Float64Type,
		"float64":     types.Float64Type,
		"map":         types.MapType{ElemType: types.StringType},
		"list":        types.ListType{ElemType: types.StringType},
		"object":      types.ObjectType{AttrTypes: objectAttrTypes},
		"object_list": types.ListType{ElemType: types.ObjectType{AttrTypes: objectAttrTypes}},
		"object_map":  types.MapType{ElemType: types.ObjectType{AttrTypes: objectAttrTypes}},
	}

	var objectListValue attr.Value = types.ListNull(types.ObjectType{AttrTypes: objectAttrTypes})
	if len(m.ObjectList) > 0 {
		var objectList []attr.Value
		for _, o := range m.ObjectList {
			oObj, d := o.ToResourceModel(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return types.ObjectNull(attributeTypes), diags
			}
			objectList = append(objectList, oObj)
		}
		var d diag.Diagnostics
		objectListValue, d = basetypes.NewListValue(types.ObjectType{AttrTypes: objectAttrTypes}, objectList)
		diags.Append(d...)
	}

	var objectMapValue attr.Value = types.MapNull(types.ObjectType{AttrTypes: objectAttrTypes})
	if len(m.ObjectMap) > 0 {
		var objectMap = make(map[string]attr.Value)
		for k, v := range m.ObjectMap {
			oObj, d := v.ToResourceModel(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return types.ObjectNull(attributeTypes), diags
			}
			objectMap[k] = oObj
		}
		var d diag.Diagnostics
		objectMapValue, d = basetypes.NewMapValue(types.ObjectType{AttrTypes: objectAttrTypes}, objectMap)
		diags.Append(d...)
	}

	objectValue, d := m.Object.ToResourceModel(ctx)
	diags.Append(d...)

	var mapValue attr.Value = types.MapNull(types.StringType)
	if len(m.Map) > 0 {
		var labels = make(map[string]attr.Value)
		for k, v := range m.Map {
			labels[k] = types.StringValue(v)
		}
		var d diag.Diagnostics
		mapValue, d = basetypes.NewMapValue(types.StringType, labels)
		diags.Append(d...)
	}

	var listValue attr.Value = types.ListNull(types.StringType)
	if len(m.List) > 0 {
		var list []attr.Value
		for _, v := range m.List {
			list = append(list, types.StringValue(v))
		}
		var d diag.Diagnostics
		listValue, d = basetypes.NewListValue(types.StringType, list)
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectNull(attributeTypes), diags
	}

	attributes := map[string]attr.Value{
		"string":      types.StringValue(m.String),
		"bool":        types.BoolValue(m.Bool),
		"number":      types.NumberValue(big.NewFloat(m.Number)),
		"int32":       types.Int64Value(int64(m.Int32)),
		"int64":       types.Int64Value(m.Int64),
		"float32":     types.Float64Value(float64(m.Float32)),
		"float64":     types.Float64Value(m.Float64),
		"map":         mapValue,
		"list":        listValue,
		"object":      objectValue,
		"object_list": objectListValue,
		"object_map":  objectMapValue,
	}

	return basetypes.NewObjectValue(attributeTypes, attributes)
}

var objectAttrTypes = map[string]attr.Type{
	"string_attribute": types.StringType,
}

func (m *Object) ToResourceModel(ctx context.Context) (types.Object, diag.Diagnostics) {
	return basetypes.NewObjectValueFrom(ctx, objectAttrTypes, m)
}

func (m *Rancher2Dev2ResourceModel) ToPlan(ctx context.Context, diags *diag.Diagnostics) tfsdk.Plan {
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2ResourceModel to tfsdk.Plan: \n%+v", pp.PrettyPrint(m)))
	if diags.HasError() {
		return tfsdk.Plan{}
	}
	r := NewRancher2Dev2Resource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	plan := tfsdk.Plan{
		Schema: s.Schema,
	}
	if diags.HasError() {
		return plan
	}
	diags.Append(plan.Set(ctx, m)...)
	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2ResourceModel to tfsdk.Plan: \n%+v", pp.PrettyPrint(plan)))
	return plan
}

func (m *Rancher2Dev2ResourceModel) ToState(ctx context.Context, diags *diag.Diagnostics) tfsdk.State {
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2ResourceModel to tfsdk.State: \n%+v", pp.PrettyPrint(m)))
	if diags.HasError() {
		return tfsdk.State{}
	}
	r := NewRancher2Dev2Resource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	state := tfsdk.State{
		Schema: s.Schema,
	}
	diags.Append(state.Set(ctx, m)...)
	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2ResourceModel to tfsdk.State: \n%+v", pp.PrettyPrint(state)))
	return state
}
