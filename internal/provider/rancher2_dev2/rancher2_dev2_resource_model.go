package rancher2_dev2

import (
	"context"
	// "encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	mta "github.com/rancher/terraform-provider-rancher2/internal/provider/rancher2_metadata"
)

type Rancher2Dev2ResourceModel struct {
	ID           types.String `tfsdk:"id"`
	APIVersion   types.String `tfsdk:"api_version"`
	Kind         types.String `tfsdk:"kind"`
	Metadata     types.Object `tfsdk:"metadata"`
	Spec         types.Object `tfsdk:"spec"`
	Status       types.String `tfsdk:"status"` // json string
	APIResponses types.Map    `tfsdk:"api_responses"`
}

func (m *Rancher2Dev2ResourceModel) ToPlan(ctx context.Context, diags *diag.Diagnostics) tfsdk.Plan {
	// when chaining commands there won't be time to validate diags
	// so each conversion function should check for diag errors and short circuit
	if diags.HasError() {
		return tfsdk.Plan{}
	}
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2ResourceModel to tfsdk.Plan: \n%+v", pp.PrettyPrint(m)))
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
	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2ResourceModel to tfsdk.Plan: \n%+v", pp.PrettyPrint(plan.Raw)))
	return plan
}

func (m *Rancher2Dev2ResourceModel) ToState(ctx context.Context, diags *diag.Diagnostics) tfsdk.State {
	// when chaining commands there won't be time to validate diags
	// so each conversion function should check for diag errors and short circuit
	if diags.HasError() {
		return tfsdk.State{}
	}
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2ResourceModel to tfsdk.State: \n%+v", pp.PrettyPrint(m)))
	r := NewRancher2Dev2Resource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	state := tfsdk.State{
		Schema: s.Schema,
	}
	diags.Append(state.Set(ctx, m)...)
	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2ResourceModel to tfsdk.State: \n%+v", pp.PrettyPrint(state.Raw)))
	return state
}

// ToGoModel converts a Rancher2Dev2ResourceModel to a Rancher2Dev2Model.
func (m *Rancher2Dev2ResourceModel) ToGoModel(ctx context.Context, diags *diag.Diagnostics) *Rancher2Dev2Model {
	// when chaining commands there won't be time to validate diags
	// so each conversion function should check for diag errors and short circuit
	if diags.HasError() {
		return nil
	}
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2ResourceModel to Rancher2Dev2Model: \n%+v", pp.PrettyPrint(m)))

	obj := &Rancher2Dev2Model{}

	obj.ID = m.ID.ValueString()
	obj.APIVersion = m.APIVersion.ValueString()
	obj.Kind = m.Kind.ValueString()
	obj.Status = m.Status.ValueString()

	// Metadata
	metadata := mta.ToGoModel(ctx, diags, m.Metadata)
	if diags.HasError() {
		return nil
	}
	obj.Metadata = *metadata

	// Spec
	spec := specToGoModel(ctx, diags, m.Spec)
	if diags.HasError() {
		return nil
	}
	obj.Spec = *spec

	// APIResponse
	if !m.APIResponses.IsNull() && !m.APIResponses.IsUnknown() {
		obj.APIResponses = make(map[string]APIResponse)
		var tmpAPIResponses map[string]types.Object
		diags.Append(m.APIResponses.ElementsAs(ctx, &tmpAPIResponses, false)...)
		if diags.HasError() {
			return nil
		}

		for k, v := range tmpAPIResponses {
			if resp := apiResponseToGoModel(ctx, diags, v); resp != nil {
				obj.APIResponses[k] = *resp
			}
		}
	}

	if diags.HasError() {
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2ResourceModel to Rancher2Dev2Model: \n%+v", pp.PrettyPrint(obj)))
	return obj
}

// specToGoModel converts a types.Object to a Spec struct.
func specToGoModel(ctx context.Context, diags *diag.Diagnostics, specObj types.Object) *Spec {
	// when chaining commands there won't be time to validate diags
	// so each conversion function should check for diag errors and short circuit
	if diags.HasError() {
		return nil
	}
	if specObj.IsNull() || specObj.IsUnknown() {
		return &Spec{}
	}

	type TmpObject struct {
		StringAttribute types.String `tfsdk:"string_attribute"`
	}
	type TmpSpec struct {
		String     types.String  `tfsdk:"string"`
		Bool       types.Bool    `tfsdk:"bool"`
		Number     types.Number  `tfsdk:"number"`
		Int32      types.Int32   `tfsdk:"int32"`
		Int64      types.Int64   `tfsdk:"int64"`
		Float32    types.Float32 `tfsdk:"float32"`
		Float64    types.Float64 `tfsdk:"float64"`
		Map        types.Map     `tfsdk:"map"`
		List       types.List    `tfsdk:"list"`
		Object     types.Object  `tfsdk:"object"`
		ObjectList types.List    `tfsdk:"object_list"`
		ObjectMap  types.Map     `tfsdk:"object_map"`
	}
	var tmpSpec TmpSpec
	diags.Append(specObj.As(ctx, &tmpSpec, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	obj := &Spec{}

	obj.String = tmpSpec.String.ValueString()
	obj.Bool = tmpSpec.Bool.ValueBool()
	f64, _ := tmpSpec.Number.ValueBigFloat().Float64()
	obj.Number = f64
	obj.Int32 = tmpSpec.Int32.ValueInt32()
	obj.Int64 = tmpSpec.Int64.ValueInt64()
	obj.Float32 = tmpSpec.Float32.ValueFloat32()
	obj.Float64 = tmpSpec.Float64.ValueFloat64()
	diags.Append(tmpSpec.Map.ElementsAs(ctx, &obj.Map, false)...)
	diags.Append(tmpSpec.List.ElementsAs(ctx, &obj.List, false)...)

	// Object
	if !tmpSpec.Object.IsNull() && !tmpSpec.Object.IsUnknown() {
		var tmpObject TmpObject
		diags.Append(tmpSpec.Object.As(ctx, &tmpObject, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil
		}
		obj.Object.StringAttribute = tmpObject.StringAttribute.ValueString()
	}

	// ObjectList
	if !tmpSpec.ObjectList.IsNull() && !tmpSpec.ObjectList.IsUnknown() {
		var tmpObjectList []TmpObject
		diags.Append(tmpSpec.ObjectList.ElementsAs(ctx, &tmpObjectList, false)...)
		if diags.HasError() {
			return nil
		}
		for _, tmpObj := range tmpObjectList {
			obj.ObjectList = append(obj.ObjectList, Object{StringAttribute: tmpObj.StringAttribute.ValueString()})
		}
	}

	// ObjectMap
	if !tmpSpec.ObjectMap.IsNull() && !tmpSpec.ObjectMap.IsUnknown() {
		var tmpObjectMap map[string]TmpObject
		diags.Append(tmpSpec.ObjectMap.ElementsAs(ctx, &tmpObjectMap, false)...)
		if diags.HasError() {
			return nil
		}
		obj.ObjectMap = make(map[string]Object)
		for k, tmpObj := range tmpObjectMap {
			obj.ObjectMap[k] = Object{StringAttribute: tmpObj.StringAttribute.ValueString()}
		}
	}

	if diags.HasError() {
		return nil
	}

	return obj
}

// apiResponseToGoModel converts a types.Object to an APIResponse struct.
func apiResponseToGoModel(ctx context.Context, diags *diag.Diagnostics, apiRespObj types.Object) *APIResponse {
	// when chaining commands there won't be time to validate diags
	// so each conversion function should check for diag errors and short circuit
	if diags.HasError() {
		return nil
	}
	if apiRespObj.IsNull() || apiRespObj.IsUnknown() {
		return &APIResponse{}
	}

	type TmpAPIResponse struct {
		Headers    types.Map    `tfsdk:"headers"`
		Body       types.String `tfsdk:"body"`
		StatusCode types.Int64  `tfsdk:"status_code"`
	}
	var tmpAPIResponse TmpAPIResponse
	diags.Append(apiRespObj.As(ctx, &tmpAPIResponse, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	resp := &APIResponse{}
	resp.Body = tmpAPIResponse.Body.ValueString()
	resp.StatusCode = tmpAPIResponse.StatusCode.ValueInt64()

	if !tmpAPIResponse.Headers.IsNull() && !tmpAPIResponse.Headers.IsUnknown() {
		var headers map[string]types.List
		diags.Append(tmpAPIResponse.Headers.ElementsAs(ctx, &headers, false)...)
		if diags.HasError() {
			return nil
		}

		goHeaders := make(map[string][]string)
		if len(headers) != 0 {
			for k, v := range headers {
				var goList []string
				diags.Append(v.ElementsAs(ctx, &goList, false)...)
				if diags.HasError() {
					return nil
				}
				goHeaders[k] = goList
			}
		}
		resp.Headers = goHeaders
	}

	if diags.HasError() {
		return nil
	}

	return resp
}
