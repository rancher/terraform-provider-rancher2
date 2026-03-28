package rancher2_dev2

import (
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

var testObjectValue = types.ObjectValueMust(objectAttrTypes, map[string]attr.Value{
	"string_attribute": types.StringValue("test_object_string"),
})

var testSpecTypesObject = types.ObjectValueMust(
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
		"object": testObjectValue,
		"object_list": types.ListValueMust(
			types.ObjectType{AttrTypes: objectAttrTypes},
			[]attr.Value{testObjectValue},
		),
		"object_map": types.MapValueMust(
			types.ObjectType{AttrTypes: objectAttrTypes},
			map[string]attr.Value{
				"obj_map_key": testObjectValue,
			},
		),
	},
)

var testFullDevResourceModel = Rancher2Dev2ResourceModel{ // Terraform resource model
	ID:         types.StringValue("test_id"),
	APIVersion: types.StringValue("v1"),
	Kind:       types.StringValue("Rancher2Dev2"),
	Status:     types.StringValue("active"),
	Metadata:   mta.SampleMetadataTypesObject(),
	Spec:       testSpecTypesObject,
}

func TestRancher2Dev2ResourceModel(t *testing.T) {
	t.Run("ToPlan", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  Rancher2Dev2ResourceModel
		}{
			{
				"Basic",
				testFullDevResourceModel,
      },
		}
		for _, tc := range testCases {
			ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
			defer log()

      diags := &diag.Diagnostics{}
      got := tc.fit.ToPlan(ctx, diags)
			if diags.HasError() {
				t.Fatalf("ToPlan had unexpected diagnostics: %v", diags)
			}
      want, err := h.Plan(ctx, &Rancher2Dev2Resource{}, &tc.fit)
      if err != nil {
        t.Fatalf("unexpected error: %v", err)
      }
      if diff := cmp.Diff(want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
    }
  })
	t.Run("ToState", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  Rancher2Dev2ResourceModel
		}{
			{
				"Basic",
				testFullDevResourceModel,
			},
		}
		for _, tc := range testCases {
			ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
			defer log()

      diags := &diag.Diagnostics{}
      got := tc.fit.ToState(ctx, diags)
			if diags.HasError() {
				t.Fatalf("ToState had unexpected diagnostics: %v", diags)
			}
      want, err := h.State(ctx, &Rancher2Dev2Resource{}, &tc.fit)
      if err != nil {
        t.Fatalf("unexpected error: %v", err)
      }
			if diff := cmp.Diff(want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
    }
  })
}

