package rancher2_dev2

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	mta "github.com/rancher/terraform-provider-rancher2/internal/provider/rancher2_metadata"
)

type Rancher2Dev2Model struct {
	ID         string       `tfsdk:"id" json:"id"`
	APIVersion string       `tfsdk:"api_version" json:"api_version"`
	Kind       string       `tfsdk:"kind" json:"kind"`
	Metadata   mta.Metadata `tfsdk:"metadata" json:"metadata"`
	Spec       Spec         `tfsdk:"spec" json:"spec"`
	Status     string       `tfsdk:"status" json:"status"`
}

// ToApiRequestBody converts the model to its JSON representation for use as an API request body.
func (m *Rancher2Dev2Model) ToApiRequestBody() ([]byte, error) {

	// apiRequestBody is a temporary struct to hold the data for the API request body.
	// This is used to exclude read-only fields from the JSON output.
	type apiRequestBody struct {
		APIVersion string                 `json:"api_version"`
		Kind       string                 `json:"kind"`
		Metadata   mta.ApiRequestMetadata `json:"metadata"`
		Spec       Spec                   `json:"spec"` // there are no read only attributes in the spec
	}

	requestBody := apiRequestBody{
		APIVersion: m.APIVersion,
		Kind:       m.Kind,
		Metadata: mta.ApiRequestMetadata{
			Name:            m.Metadata.Name,
			Namespace:       m.Metadata.Namespace,
			GenerateName:    m.Metadata.GenerateName,
			Annotations:     m.Metadata.Annotations,
			Labels:          m.Metadata.Labels,
			Finalizers:      m.Metadata.Finalizers,
			OwnerReferences: m.Metadata.OwnerReferences,
		},
		Spec: m.Spec,
	}

	return json.Marshal(requestBody)
}

type Spec struct {
	String     string            `tfsdk:"string" json:"string"`
	Bool       bool              `tfsdk:"bool" json:"bool"`
	Number     float64           `tfsdk:"number" json:"number"`
	Int32      int32             `tfsdk:"int32" json:"int32"`
	Int64      int64             `tfsdk:"int64" json:"int64"`
	Float32    float32           `tfsdk:"float32" json:"float32"`
	Float64    float64           `tfsdk:"float64" json:"float64"`
	Map        map[string]string `tfsdk:"map" json:"map"`
	List       []string          `tfsdk:"list" json:"list"`
	Object     Object            `tfsdk:"object" json:"object"`
	ObjectList []Object          `tfsdk:"object_list" json:"object_list"`
	ObjectMap  map[string]Object `tfsdk:"object_map" json:"object_map"`
}

var specAttrTypes = map[string]attr.Type{
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

type Object struct {
	StringAttribute string `tfsdk:"string_attribute"`
}

var objectAttrTypes = map[string]attr.Type{
	"string_attribute": types.StringType,
}

// Converts Go model Rancher2Dev2Model to a Terraform model Rancher2Dev2ResourceModel.
// Returns a pointer so that you can call further operations such as ToPlan or ToState.
// Never returns nil.
func (obj *Rancher2Dev2Model) ToResourceModel(ctx context.Context, diags *diag.Diagnostics) *Rancher2Dev2ResourceModel {
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2Model to Rancher2Dev2ResourceModel: %+v", pp.PrettyPrint(obj)))
	if diags.HasError() {
		return &Rancher2Dev2ResourceModel{}
	}

	var data Rancher2Dev2ResourceModel

	data.ID = types.StringValue(obj.ID)
	data.APIVersion = types.StringValue(obj.APIVersion)
	data.Kind = types.StringValue(obj.Kind)
	data.Status = types.StringValue(obj.Status)

	metadata := obj.Metadata.ToTypesObject(ctx, diags)
	if diags.HasError() {
		return &Rancher2Dev2ResourceModel{}
	}
	data.Metadata = metadata

	spec := obj.Spec.ToTypesObject(ctx, diags)
	if diags.HasError() {
		return &Rancher2Dev2ResourceModel{}
	}
	data.Spec = spec

	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2Model to Rancher2Dev2ResourceModel: %+v", pp.PrettyPrint(data)))
	return &data
}

func (m *Spec) ToTypesObject(ctx context.Context, diags *diag.Diagnostics) types.Object {
	var objectListValue attr.Value = types.ListNull(types.ObjectType{AttrTypes: objectAttrTypes})
	if len(m.ObjectList) > 0 {
		var objectList []attr.Value
		for _, o := range m.ObjectList {
			oObj := o.ToTypesObject(ctx, diags)
			if diags.HasError() {
				return types.ObjectNull(specAttrTypes)
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
			oObj := v.ToTypesObject(ctx, diags)
			if diags.HasError() {
				return types.ObjectNull(specAttrTypes)
			}
			objectMap[k] = oObj
		}
		var d diag.Diagnostics
		objectMapValue, d = basetypes.NewMapValue(types.ObjectType{AttrTypes: objectAttrTypes}, objectMap)
		diags.Append(d...)
	}

	objectValue := m.Object.ToTypesObject(ctx, diags)

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
		return types.ObjectNull(specAttrTypes)
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
	obj, d := basetypes.NewObjectValue(specAttrTypes, attributes)
	diags.Append(d...)
	return obj
}

func (m *Object) ToTypesObject(ctx context.Context, diags *diag.Diagnostics) types.Object {
	obj, d := basetypes.NewObjectValueFrom(ctx, objectAttrTypes, m)
	diags.Append(d...)
	return obj
}
