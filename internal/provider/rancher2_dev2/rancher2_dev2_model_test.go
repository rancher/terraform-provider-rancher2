package rancher2_dev2

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	h "github.com/rancher/terraform-provider-rancher2/internal/provider/test_helpers"
)

var metadataAttrTypes = map[string]attr.Type{
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

func TestRancher2Dev2ModelToResourceModel(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {

		ownerReferenceValue := types.ObjectValueMust(
			ownerReferenceAttrTypes,
			map[string]attr.Value{
				"api_version":          types.StringValue("v1"),
				"kind":                 types.StringValue("some_kind"),
				"name":                 types.StringValue("owner"),
				"uid":                  types.StringValue("some_uid"),
				"controller":           types.BoolValue(true),
				"block_owner_deletion": types.BoolValue(true),
			},
		)

		metadataValue := types.ObjectValueMust(
			metadataAttrTypes,
			map[string]attr.Value{
				"name":                          types.StringValue("test_metadata"),
				"namespace":                     types.StringValue("test_namespace"),
				"generate_name":                 types.StringValue("test_generate_name"),
				"annotations":                   types.MapValueMust(types.StringType, map[string]attr.Value{"ann_key": types.StringValue("ann_value")}),
				"labels":                        types.MapValueMust(types.StringType, map[string]attr.Value{"label_key": types.StringValue("label_value")}),
				"finalizers":                    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("finalizer_a")}),
				"owner_references":              types.ListValueMust(types.ObjectType{AttrTypes: ownerReferenceAttrTypes}, []attr.Value{ownerReferenceValue}),
				"uid":                           types.StringValue("test_uid"),
				"generation":                    types.Int64Value(1),
				"creation_timestamp":            types.StringValue("2023-01-01T00:00:00Z"),
				"deletion_grace_period_seconds": types.Int64Value(30),
				"deletion_timestamp":            types.StringValue("2023-01-01T01:00:00Z"),
				"managed_fields":                types.StringValue("test_managed_fields"),
				"resource_version":              types.StringValue("test_resource_version"),
				"self_link":                     types.StringValue("/api/v1/namespaces/default/rancher2_dev2s/test"),
			},
		)

		objectValue := types.ObjectValueMust(
			objectAttrTypes,
			map[string]attr.Value{
				"string_attribute": types.StringValue("test_object_string"),
			},
		)

		specValue := types.ObjectValueMust(
			specAttrTypes,
			map[string]attr.Value{
				"string":      types.StringValue("test_spec_string"),
				"bool":        types.BoolValue(true),
				"number":      types.NumberValue(big.NewFloat(1.25)),
				"int32":       types.Int64Value(123),
				"int64":       types.Int64Value(456),
				"float32":     types.Float64Value(1.25),
				"float64":     types.Float64Value(4.50),
				"map":         types.MapValueMust(types.StringType, map[string]attr.Value{"map_key": types.StringValue("map_value")}),
				"list":        types.ListValueMust(types.StringType, []attr.Value{types.StringValue("list_value")}),
				"object":      objectValue,
				"object_list": types.ListValueMust(types.ObjectType{AttrTypes: objectAttrTypes}, []attr.Value{objectValue}),
				"object_map":  types.MapValueMust(types.ObjectType{AttrTypes: objectAttrTypes}, map[string]attr.Value{"obj_map_key": objectValue}),
			},
		)

		fit := Rancher2Dev2Model{
			ID:         "test_id",
			APIVersion: "v1",
			Kind:       "Rancher2Dev2",
			Status:     "active",
			Metadata: Metadata{
				Name:         "test_metadata",
				Namespace:    "test_namespace",
				GenerateName: "test_generate_name",
				Annotations:  map[string]string{"ann_key": "ann_value"},
				Labels:       map[string]string{"label_key": "label_value"},
				Finalizers:   []string{"finalizer_a"},
				OwnerReferences: []OwnerReference{
					{
						APIVersion:         "v1",
						Kind:               "some_kind",
						Name:               "owner",
						UID:                "some_uid",
						Controller:         true,
						BlockOwnerDeletion: true,
					},
				},
				UID:                       "test_uid",
				Generation:                1,
				CreationTimestamp:         "2023-01-01T00:00:00Z",
				DeletionGracePeriodSecond: 30,
				DeletionTimestamp:         "2023-01-01T01:00:00Z",
				ManagedFields:             "test_managed_fields",
				ResourceVersion:           "test_resource_version",
				SelfLink:                  "/api/v1/namespaces/default/rancher2_dev2s/test",
			},
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
		}

		want := &Rancher2Dev2ResourceModel{
			ID:         types.StringValue("test_id"),
			APIVersion: types.StringValue("v1"),
			Kind:       types.StringValue("Rancher2Dev2"),
			Status:     types.StringValue("active"),
			Metadata:   metadataValue,
			Spec:       specValue,
		}

		var buf bytes.Buffer
		defer h.PrintLog(t, &buf, "ERROR")
		ctx := h.GenerateTestContext(t, &buf, nil)

		got := fit.ToResourceModel(ctx, &diag.Diagnostics{})

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
			t.Errorf("unexpected diff (-want, +got) = %s", diff)
		}
	})
}

func TestMetadataToResourceModel(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		ownerReferenceValue := types.ObjectValueMust(
			ownerReferenceAttrTypes,
			map[string]attr.Value{
				"api_version":          types.StringValue("v1"),
				"kind":                 types.StringValue("some_kind"),
				"name":                 types.StringValue("owner"),
				"uid":                  types.StringValue("some_uid"),
				"controller":           types.BoolValue(true),
				"block_owner_deletion": types.BoolValue(true),
			},
		)

		fit := Metadata{
			Name:         "test_metadata",
			Namespace:    "test_namespace",
			GenerateName: "test_generate_name",
			Annotations:  map[string]string{"ann_key": "ann_value"},
			Labels:       map[string]string{"label_key": "label_value"},
			Finalizers:   []string{"finalizer_a"},
			OwnerReferences: []OwnerReference{
				{
					APIVersion:         "v1",
					Kind:               "some_kind",
					Name:               "owner",
					UID:                "some_uid",
					Controller:         true,
					BlockOwnerDeletion: true,
				},
			},
			UID:                       "test_uid",
			Generation:                1,
			CreationTimestamp:         "2023-01-01T00:00:00Z",
			DeletionGracePeriodSecond: 30,
			DeletionTimestamp:         "2023-01-01T01:00:00Z",
			ManagedFields:             "test_managed_fields",
			ResourceVersion:           "test_resource_version",
			SelfLink:                  "/api/v1/namespaces/default/rancher2_dev2s/test",
		}

		want := types.ObjectValueMust(
			metadataAttrTypes,
			map[string]attr.Value{
				"name":                          types.StringValue("test_metadata"),
				"namespace":                     types.StringValue("test_namespace"),
				"generate_name":                 types.StringValue("test_generate_name"),
				"annotations":                   types.MapValueMust(types.StringType, map[string]attr.Value{"ann_key": types.StringValue("ann_value")}),
				"labels":                        types.MapValueMust(types.StringType, map[string]attr.Value{"label_key": types.StringValue("label_value")}),
				"finalizers":                    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("finalizer_a")}),
				"owner_references":              types.ListValueMust(types.ObjectType{AttrTypes: ownerReferenceAttrTypes}, []attr.Value{ownerReferenceValue}),
				"uid":                           types.StringValue("test_uid"),
				"generation":                    types.Int64Value(1),
				"creation_timestamp":            types.StringValue("2023-01-01T00:00:00Z"),
				"deletion_grace_period_seconds": types.Int64Value(30),
				"deletion_timestamp":            types.StringValue("2023-01-01T01:00:00Z"),
				"managed_fields":                types.StringValue("test_managed_fields"),
				"resource_version":              types.StringValue("test_resource_version"),
				"self_link":                     types.StringValue("/api/v1/namespaces/default/rancher2_dev2s/test"),
			},
		)

		var buf bytes.Buffer
		defer h.PrintLog(t, &buf, "ERROR")
		ctx := h.GenerateTestContext(t, &buf, nil)

		got, diags := fit.ToResourceModel(ctx)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
			t.Errorf("unexpected diff (-want, +got) = %s", diff)
		}
	})
}

func TestOwnerReferenceToResourceModel(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		fit := OwnerReference{
			APIVersion:         "v1",
			Kind:               "some_kind",
			Name:               "owner",
			UID:                "some_uid",
			Controller:         true,
			BlockOwnerDeletion: true,
		}

		want := types.ObjectValueMust(
			ownerReferenceAttrTypes,
			map[string]attr.Value{
				"api_version":          types.StringValue("v1"),
				"kind":                 types.StringValue("some_kind"),
				"name":                 types.StringValue("owner"),
				"uid":                  types.StringValue("some_uid"),
				"controller":           types.BoolValue(true),
				"block_owner_deletion": types.BoolValue(true),
			},
		)

		var buf bytes.Buffer
		defer h.PrintLog(t, &buf, "ERROR")
		ctx := h.GenerateTestContext(t, &buf, nil)

		got, diags := fit.ToResourceModel(ctx)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
			t.Errorf("unexpected diff (-want, +got) = %s", diff)
		}
	})
}

func TestSpecToResourceModel(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		objectValue := types.ObjectValueMust(
			objectAttrTypes,
			map[string]attr.Value{
				"string_attribute": types.StringValue("test_object_string"),
			},
		)

		fit := Spec{
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
		}

		want := types.ObjectValueMust(
			specAttrTypes,
			map[string]attr.Value{
				"string":      types.StringValue("test_spec_string"),
				"bool":        types.BoolValue(true),
				"number":      types.NumberValue(big.NewFloat(1.25)),
				"int32":       types.Int64Value(123),
				"int64":       types.Int64Value(456),
				"float32":     types.Float64Value(1.25),
				"float64":     types.Float64Value(4.50),
				"map":         types.MapValueMust(types.StringType, map[string]attr.Value{"map_key": types.StringValue("map_value")}),
				"list":        types.ListValueMust(types.StringType, []attr.Value{types.StringValue("list_value")}),
				"object":      objectValue,
				"object_list": types.ListValueMust(types.ObjectType{AttrTypes: objectAttrTypes}, []attr.Value{objectValue}),
				"object_map":  types.MapValueMust(types.ObjectType{AttrTypes: objectAttrTypes}, map[string]attr.Value{"obj_map_key": objectValue}),
			},
		)

		var buf bytes.Buffer
		defer h.PrintLog(t, &buf, "ERROR")
		ctx := h.GenerateTestContext(t, &buf, nil)

		got, diags := fit.ToResourceModel(ctx)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
			t.Errorf("unexpected diff (-want, +got) = %s", diff)
		}
	})
}

func TestObjectToResourceModel(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		fit := Object{
			StringAttribute: "test_object_string",
		}

		want := types.ObjectValueMust(
			objectAttrTypes,
			map[string]attr.Value{
				"string_attribute": types.StringValue("test_object_string"),
			},
		)

		var buf bytes.Buffer
		defer h.PrintLog(t, &buf, "ERROR")
		ctx := h.GenerateTestContext(t, &buf, nil)

		got, diags := fit.ToResourceModel(ctx)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
			t.Errorf("unexpected diff (-want, +got) = %s", diff)
		}
	})
}
