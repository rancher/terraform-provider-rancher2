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
	ID           string                 `tfsdk:"id" json:"id"`
	APIVersion   string                 `tfsdk:"api_version" json:"api_version"`
	Kind         string                 `tfsdk:"kind" json:"kind"`
	Metadata     mta.Metadata           `tfsdk:"metadata" json:"metadata"`
	Spec         Spec                   `tfsdk:"spec" json:"spec"`
	Status       any                    `tfsdk:"status" json:"status"` // must be marshalable
	APIResponses map[string]APIResponse `tfsdk:"api_responses" json:"api_responses"`
}

// APIResponses is a map of function to response, eg. create, read, update, delete.
// Dev Only, remove this when using this resource as a template
type APIResponse struct {
	Headers    map[string][]string `tfsdk:"headers" json:"headers"`
	Body       string              `tfsdk:"body" json:"body"`
	StatusCode int64               `tfsdk:"status_code" json:"status_code"`
}

// IsEmpty checks if the APIResponse is empty (has all zero values).
func (ar APIResponse) IsEmpty() bool {
	return len(ar.Headers) == 0 && ar.Body == "" && ar.StatusCode == 0
}

// FromAPIResponseBody converts the raw body response from the API to the Rancher2Dev2Model object.
func (m *Rancher2Dev2Model) FromAPIResponseBody(ctx context.Context, response []byte, diags *diag.Diagnostics) {
	// when chaining commands there won't be time to validate diags
	// so each conversion function should check for diag errors and short circuit
	if diags.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Converting API response body to Rancher2Dev2Model: \n%s\n", pp.PrettyPrint(json.RawMessage(response))))
	var tmpModel struct {
		APIVersion string          `json:"api_version"`
		Kind       string          `json:"kind"`
		Spec       Spec            `json:"spec"`
		Status     json.RawMessage `json:"status"`
	}
	err := json.Unmarshal(response, &tmpModel)
	if err != nil {
		diags.AddError("Error unmarshaling response body:", err.Error())
		return
	}
	m.APIVersion = tmpModel.APIVersion
	m.Kind = tmpModel.Kind
	m.Spec = tmpModel.Spec
	if len(tmpModel.Status) > 0 && string(tmpModel.Status) != "null" {
		m.Status = string(tmpModel.Status)
	} else {
		m.Status = nil
	}

	// Metadata
	var tmpMetadata struct {
		Metadata mta.Metadata `json:"metadata"`
	}
	err = json.Unmarshal(response, &tmpMetadata)
	if err != nil {
		diags.AddError("Error unmarshaling response body metadata:", err.Error())
		return
	}
	m.Metadata = tmpMetadata.Metadata
	tflog.Debug(ctx, fmt.Sprintf("Converted API response body to Rancher2Dev2Model: %+v\n", pp.PrettyPrint(m)))
}

// ToApiRequestBody converts the model to elements that can easily be json marshaled
func (m *Rancher2Dev2Model) ToApiRequestBody(diags *diag.Diagnostics) (jsonMarshalable any) {
	// when chaining commands there won't be time to validate diags
	// so each conversion function should check for diag errors and short circuit
	if diags.HasError() {
		return nil
	}

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
	return requestBody
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

var apiResponseAttrTypes = map[string]attr.Type{
	"headers": types.MapType{
		ElemType: types.ListType{
			ElemType: types.StringType,
		},
	},
	"body":        types.StringType,
	"status_code": types.Int64Type,
}

type Object struct {
	StringAttribute string `tfsdk:"string_attribute"`
}

var objectAttrTypes = map[string]attr.Type{
	"string_attribute": types.StringType,
}

func (goModel *Rancher2Dev2Model) ToResourceModel(ctx context.Context, diags *diag.Diagnostics) (resourceModelPtr *Rancher2Dev2ResourceModel) {
	// when chaining commands there won't be time to validate diags
	// so each conversion function should check for diag errors and short circuit
	if diags.HasError() {
		return nil
	}
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2Model to Rancher2Dev2ResourceModel: \n%+v", pp.PrettyPrint(goModel)))

	resourceModel := Rancher2Dev2ResourceModel{
		ID:         types.StringValue(goModel.ID),
		APIVersion: types.StringValue(goModel.APIVersion),
		Kind:       types.StringValue(goModel.Kind),
		Metadata:   goModel.Metadata.ToTypesObject(ctx, diags),
		Spec:       goModel.Spec.ToTypesObject(ctx, diags),
	}

	if len(goModel.APIResponses) == 0 {
		diags.AddWarning("APIResponses is empty", "The APIResponses map is empty or nil.")
		// return &resourceModel
	}

	apiResponseObjectType := types.ObjectType{
		AttrTypes: apiResponseAttrTypes,
	}

	apiResponseConverter := func(ctx context.Context, diags *diag.Diagnostics, resp APIResponse) attr.Value {
		return resp.ToTypesObject(ctx, diags)
	}

	apiResponsesTypesMap := mapToTypesMap(ctx, diags, goModel.APIResponses, apiResponseObjectType, apiResponseConverter)
	if diags.HasError() {
		return &resourceModel
	}
	if apiResponsesTypesMap.IsNull() {
		diags.AddError("APIResponses types is nil", "The APIResponses types map is nil.")
		return &resourceModel
	}
	resourceModel.APIResponses = apiResponsesTypesMap

	if goModel.Status != nil {
		if status, ok := goModel.Status.(string); ok {
			resourceModel.Status = types.StringValue(status)
		} else {
			statusJson, err := json.Marshal(goModel.Status)
			if err != nil {
				diags.AddError("Error marshaling status", err.Error())
				return &resourceModel
			}
			resourceModel.Status = types.StringValue(string(statusJson))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2Model to Rancher2Dev2ResourceModel: \n%+v", pp.PrettyPrint(resourceModel)))
	return &resourceModel
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

	var objectMapValue attr.Value
	objectMapValue = types.MapNull(types.ObjectType{AttrTypes: objectAttrTypes})
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

func (m *APIResponse) ToTypesObject(ctx context.Context, diags *diag.Diagnostics) types.Object {
	if m == nil {
		return types.ObjectNull(apiResponseAttrTypes)
	}

	var headersValue attr.Value = types.MapNull(types.ListType{ElemType: types.StringType})
	if len(m.Headers) > 0 {
		headersMap := make(map[string]attr.Value)
		for k, v := range m.Headers {
			list, d := types.ListValueFrom(ctx, types.StringType, v)
			diags.Append(d...)
			if diags.HasError() {
				return types.ObjectNull(apiResponseAttrTypes)
			}
			headersMap[k] = list
		}
		var d diag.Diagnostics
		headersValue, d = types.MapValue(types.ListType{ElemType: types.StringType}, headersMap)
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectNull(apiResponseAttrTypes)
	}

	attributes := map[string]attr.Value{
		"headers":     headersValue,
		"body":        types.StringValue(m.Body),
		"status_code": types.Int64Value(m.StatusCode),
	}
	obj, d := types.ObjectValue(apiResponseAttrTypes, attributes)
	diags.Append(d...)
	return obj
}

func (m *Object) ToTypesObject(ctx context.Context, diags *diag.Diagnostics) types.Object {
	attributes := map[string]attr.Value{
		"string_attribute": types.StringValue(m.StringAttribute),
	}

	obj, d := basetypes.NewObjectValue(objectAttrTypes, attributes)
	diags.Append(d...)
	return obj
}

// need a function that converts a map to a types map
// would be good if it can handle both map[string]string and map[string]struct
func mapToTypesMap[T any](ctx context.Context, diags *diag.Diagnostics, m map[string]T, elemType attr.Type, converter func(context.Context, *diag.Diagnostics, T) attr.Value) types.Map {

	elements := make(map[string]attr.Value, len(m))
	for k, v := range m {
		elements[k] = converter(ctx, diags, v)
	}
	mapVal, d := types.MapValue(elemType, elements)
	diags.Append(d...)

	return mapVal
}
