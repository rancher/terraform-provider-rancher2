package rancher2_dev2

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	mta "github.com/rancher/terraform-provider-rancher2/internal/provider/rancher2_metadata"
	h "github.com/rancher/terraform-provider-rancher2/internal/provider/test_helpers"
)

var specTypesObject = types.ObjectValueMust(
	specAttrTypes,
	map[string]attr.Value{
		"string":  types.StringValue("test_spec_string"),
		"bool":    types.BoolValue(true),
		"number":  types.NumberValue(big.NewFloat(1.25)),
		"int32":   types.Int64Value(123),
		"int64":   types.Int64Value(456),
		"float32": types.Float64Value(1.25),
		"float64": types.Float64Value(4.50),
		"map": types.MapValueMust(types.StringType, map[string]attr.Value{
			"map_key": types.StringValue("map_value"),
		}),
		"list": types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("list_value"),
		}),
		"object": types.ObjectValueMust(objectAttrTypes, map[string]attr.Value{
			"string_attribute": types.StringValue("test_object_string"),
		}),
		"object_list": types.ListValueMust(types.ObjectType{AttrTypes: objectAttrTypes}, []attr.Value{
			types.ObjectValueMust(objectAttrTypes, map[string]attr.Value{
				"string_attribute": types.StringValue("test_object_string"),
			},
			)}),
		"object_map": types.MapValueMust(types.ObjectType{AttrTypes: objectAttrTypes}, map[string]attr.Value{
			"obj_map_key": types.ObjectValueMust(objectAttrTypes, map[string]attr.Value{
				"string_attribute": types.StringValue("test_object_string"),
			},
			)}),
	},
)

var sampleAPIResponse = APIResponse{
	Headers:    map[string][]string{"Content-Type": {"application/json"}},
	Body:       `{"id":"test","type":"dev2"}`,
	StatusCode: 200,
}

var sampleAPIResponses = map[string]APIResponse{
	"create": sampleAPIResponse,
}

var sampleAPIResponseObject = types.ObjectValueMust(
	apiResponseAttrTypes,
	map[string]attr.Value{
		"headers": types.MapValueMust(
			types.ListType{ElemType: types.StringType},
			map[string]attr.Value{
				"Content-Type": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("application/json")}),
			},
		),
		"body":        types.StringValue(`{"id":"test","type":"dev2"}`),
		"status_code": types.Int64Value(200),
	},
)

var sampleAPIResponsesMap = types.MapValueMust(
	types.ObjectType{AttrTypes: apiResponseAttrTypes},
	map[string]attr.Value{
		"create": sampleAPIResponseObject,
	},
)

var fullDevModel = Rancher2Dev2Model{ // native Go model which we will convert to a Terraform resource model
	ID:         "test_id",
	APIVersion: "v1",
	Kind:       "Rancher2Dev2",
	Status:     "active",
	Metadata:   mta.SampleMetadataGoModel(),
	Spec: Spec{
		String:  "test_spec_string",
		Bool:    true,
		Number:  1.25,
		Int32:   123,
		Int64:   456,
		Float32: 1.25,
		Float64: 4.50,
		Map:     map[string]string{"map_key": "map_value"},
		List:    []string{"list_value"},
		Object: Object{
			StringAttribute: "test_object_string",
		},
		ObjectList: []Object{
			{StringAttribute: "test_object_string"},
		},
		ObjectMap: map[string]Object{
			"obj_map_key": {StringAttribute: "test_object_string"},
		},
	},
	APIResponses: sampleAPIResponses,
}

var fullDevResourceModel = Rancher2Dev2ResourceModel{ // Terraform resource model
	ID:           types.StringValue("test_id"),
	APIVersion:   types.StringValue("v1"),
	Kind:         types.StringValue("Rancher2Dev2"),
	Status:       types.StringValue("active"),
	Metadata:     mta.SampleMetadataTypesObject(),
	Spec:         specTypesObject,
	APIResponses: sampleAPIResponsesMap,
}

var correctedFullDevResourceModel = Rancher2Dev2ResourceModel{ // Terraform resource model
	ID:         types.StringValue("test_id"),
	APIVersion: types.StringValue("v1"),
	Kind:       types.StringValue("Rancher2Dev2"),
	Status:     types.StringValue("active"),
	Metadata: types.ObjectValueMust(
		mta.MetadataAttrTypes,
		map[string]attr.Value{
			"name":          types.StringValue("test_metadata"),
			"namespace":     types.StringValue("test_namespace"),
			"generate_name": types.StringValue(""), // mutually exclusive to the name attribute
			"annotations": types.MapValueMust(types.StringType, map[string]attr.Value{
				"ann_key": types.StringValue("ann_value"),
			}),
			"labels": types.MapValueMust(types.StringType, map[string]attr.Value{
				"label_key": types.StringValue("label_value"),
			}),
			"finalizers": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("finalizer_a"),
			}),
			"owner_references": types.ListValueMust(types.ObjectType{AttrTypes: mta.OwnerReferenceAttrTypes}, []attr.Value{
				types.ObjectValueMust(mta.OwnerReferenceAttrTypes, map[string]attr.Value{
					"api_version":          types.StringValue("v1"),
					"kind":                 types.StringValue("some_kind"),
					"name":                 types.StringValue("owner"),
					"uid":                  types.StringValue("some_uid"),
					"controller":           types.BoolValue(true),
					"block_owner_deletion": types.BoolValue(true),
				}),
			}),
			"uid":                           types.StringValue("test_uid"),
			"generation":                    types.Int64Value(1),
			"creation_timestamp":            types.StringValue("2023-01-01T00:00:00Z"),
			"deletion_grace_period_seconds": types.Int64Value(30),
			"deletion_timestamp":            types.StringValue("2023-01-01T01:00:00Z"),
			"managed_fields":                types.StringValue("{\"field\": \"test_managed_fields\"}"),
			"resource_version":              types.StringValue("test_resource_version"),
			"self_link":                     types.StringValue("/api/v1/namespaces/default/rancher2_dev2s/test"),
		},
	),
	Spec: types.ObjectValueMust(
		specAttrTypes,
		map[string]attr.Value{
			"string":  types.StringValue("test_spec_string"),
			"bool":    types.BoolValue(true),
			"number":  types.NumberValue(big.NewFloat(1.25)),
			"int32":   types.Int64Value(123),
			"int64":   types.Int64Value(456),
			"float32": types.Float64Value(1.25),
			"float64": types.Float64Value(4.50),
			"map": types.MapValueMust(types.StringType, map[string]attr.Value{
				"map_key": types.StringValue("map_value"),
			}),
			"list": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("list_value"),
			}),
			"object": types.ObjectValueMust(objectAttrTypes, map[string]attr.Value{
				"string_attribute": types.StringValue("test_object_string"),
			}),
			"object_list": types.ListValueMust(types.ObjectType{AttrTypes: objectAttrTypes}, []attr.Value{
				types.ObjectValueMust(objectAttrTypes, map[string]attr.Value{
					"string_attribute": types.StringValue("test_object_string"),
				}),
			}),
			"object_map": types.MapValueMust(types.ObjectType{AttrTypes: objectAttrTypes}, map[string]attr.Value{
				"obj_map_key": types.ObjectValueMust(objectAttrTypes, map[string]attr.Value{
					"string_attribute": types.StringValue("test_object_string"),
				}),
			}),
		},
	),
	APIResponses: sampleAPIResponsesMap,
}

func TestRancher2Dev2ModelToResourceModel(t *testing.T) {
	testCases := []struct {
		name string
		fit  Rancher2Dev2Model
		want *Rancher2Dev2ResourceModel
	}{
		{
			"Basic",
			fullDevModel,
			&fullDevResourceModel,
		},
		{
			"CorrectedNestedObjects",
			fullDevModel,
			&correctedFullDevResourceModel,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			got := tc.fit.ToResourceModel(ctx, &diag.Diagnostics{})

			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
		})
	}
}

func TestSpecToResourceModel(t *testing.T) {
	testCases := []struct {
		name string
		fit  Spec
		want types.Object
	}{
		{
			"Basic",
			Spec{
				String:  "test_spec_string",
				Bool:    true,
				Number:  1.25,
				Int32:   123,
				Int64:   456,
				Float32: 1.25,
				Float64: 4.50,
				Map:     map[string]string{"map_key": "map_value"},
				List:    []string{"list_value"},
				Object: Object{
					StringAttribute: "test_object_string",
				},
				ObjectList: []Object{
					{StringAttribute: "test_object_string"},
				},
				ObjectMap: map[string]Object{
					"obj_map_key": {StringAttribute: "test_object_string"},
				},
			},
			specTypesObject,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
			defer log()
			diags := &diag.Diagnostics{}
			got := tc.fit.ToTypesObject(ctx, diags)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
		})
	}
}

func TestObjectToResourceModel(t *testing.T) {
	testCases := []struct {
		name string
		fit  Object
		want types.Object
	}{
		{
			"Basic",
			Object{
				StringAttribute: "test_object_string",
			},
			types.ObjectValueMust(
				objectAttrTypes,
				map[string]attr.Value{
					"string_attribute": types.StringValue("test_object_string"),
				},
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
			defer log()
			diags := &diag.Diagnostics{}
			got := tc.fit.ToTypesObject(ctx, diags)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
		})
	}
}

func TestModelToApiRequestBody(t *testing.T) {
	testCases := []struct {
		name string
		fit  Rancher2Dev2Model
		want []byte
	}{
		{
			"Basic",
			fullDevModel,
			[]byte(`{
        "api_version": "v1",
        "kind": "Rancher2Dev2",` +
				mta.SampleMetadataApiRequestJson() +
				`"spec": {
          "string": "test_spec_string",
          "bool": true,
          "number": 1.25,
          "int32": 123,
          "int64": 456,
          "float32": 1.25,
          "float64": 4.5,
          "map": {
            "map_key": "map_value"
          },
          "list": [
            "list_value"
          ],
          "object": {
            "StringAttribute": "test_object_string"
          },
          "object_list": [
            {
              "StringAttribute": "test_object_string"
            }
          ],
          "object_map": {
            "obj_map_key": {
              "StringAttribute": "test_object_string"
            }
          }
        }
      }`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// compact both sides to avoid diffs in spacing
			diags := &diag.Diagnostics{}
			g, err := json.Marshal(tc.fit.ToApiRequestBody(diags))
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var got bytes.Buffer
			err = json.Compact(&got, g)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var want bytes.Buffer
			err = json.Compact(&want, tc.want)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(want.String(), got.String()); diff != "" {
				t.Log(want.String())
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
		})
	}
}
