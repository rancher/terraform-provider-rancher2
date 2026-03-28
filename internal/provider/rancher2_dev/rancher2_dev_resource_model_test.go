package rancher2_dev

import (
	"bytes"
	"context"

	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	h "github.com/rancher/terraform-provider-rancher2/internal/provider/test_helpers"
)

func TestRancher2DevResourceModel(t *testing.T) {
	t.Run("ToGoModel", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  Rancher2DevResourceModel
			want RancherDevModel
		}{
			{
				"Basic",
				getDefaultResourceModel(),
				getDefaultGoModel(),
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)
				got := tc.fit.ToGoModel(ctx)
				if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreUnexported(big.Float{})); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
	t.Run("ToPlan", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  Rancher2DevResourceModel
			want tfsdk.Plan
		}{
			{
				"Basic",
				getDefaultResourceModel(),
				getPlan(context.Background(), getDefaultResourceModel(), &diag.Diagnostics{}),
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)

				model := getDefaultResourceModel()
				got := model.ToPlan(ctx, &diag.Diagnostics{})

				if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{}, big.Float{})); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})
	t.Run("ToState", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  Rancher2DevResourceModel
			want tfsdk.State
		}{
			{
				"Basic",
				getDefaultResourceModel(),
				getState(context.Background(), getDefaultResourceModel(), &diag.Diagnostics{}),
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)

				model := getDefaultResourceModel()
				got := model.ToState(ctx, &diag.Diagnostics{})

				if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{}, big.Float{})); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})

}

func TestNestedResourceModel(t *testing.T) {
	t.Run("ToGoModel", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  NestedResourceModel
			want NestedObject
		}{
			{
				"Basic",
				NestedResourceModel{
					StringAttribute: types.StringValue("test"),
					NestedNestedObject: types.ObjectValueMust(
						nestedNestedObjectAttrTypes,
						map[string]attr.Value{
							"string_attribute": types.StringValue("test"),
							"bool_attribute":   types.BoolValue(true),
						},
					),
				},
				NestedObject{
					StringAttribute: "test",
					NestedNestedObject: NestedNestedObject{
						StringAttribute: "test",
						BoolAttribute:   true,
					},
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)

				got := NestedObject{}
				tc.fit.ToGoModel(ctx, &got)
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})
}

func TestNestedNestedResourceModel(t *testing.T) {
	t.Run("ToGoModel", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  NestedNestedResourceModel
			want NestedNestedObject
		}{
			{
				"Basic",
				NestedNestedResourceModel{
					StringAttribute: types.StringValue("test"),
					BoolAttribute:   types.BoolValue(true),
				},
				NestedNestedObject{
					StringAttribute: "test",
					BoolAttribute:   true,
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)

				got := NestedNestedObject{}
				tc.fit.ToGoModel(ctx, &got)
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})
}

// helpers.
var nestedNestedObjectAttrTypes = map[string]attr.Type{
	"string_attribute": types.StringType,
	"bool_attribute":   types.BoolType,
}
var nestedObjectAttrTypes = map[string]attr.Type{
	"string_attribute": types.StringType,
	"nested_nested_object": types.ObjectType{
		AttrTypes: nestedNestedObjectAttrTypes,
	},
}

func getDefaultResourceModel() Rancher2DevResourceModel {
	nestedNestedObjectValue := types.ObjectValueMust(
		nestedNestedObjectAttrTypes,
		map[string]attr.Value{
			"string_attribute": types.StringValue("test"),
			"bool_attribute":   types.BoolValue(true),
		},
	)

	nestedObjectValue := types.ObjectValueMust(
		nestedObjectAttrTypes,
		map[string]attr.Value{
			"string_attribute":     types.StringValue("test"),
			"nested_nested_object": nestedNestedObjectValue,
		},
	)

	return Rancher2DevResourceModel{
		NumberAttribute:  types.NumberValue(big.NewFloat(1.23)), // required
		StringAttribute:  types.StringValue("test"),             // required
		BoolAttribute:    types.BoolValue(true),                 // default value
		Int32Attribute:   types.Int32Value(123),                 // include read only as well since ToGoModel can be used on state, not just plan and config
		ID:               types.StringValue("test"),             // include read only as well since ToGoModel can be used on state, not just plan and config
		Identifier:       types.StringValue("test"),             // include read only as well since ToGoModel can be used on state, not just plan and config
		Int64Attribute:   types.Int64Value(123),
		Float64Attribute: types.Float64Value(1.23),
		Float32Attribute: types.Float32Value(1.23),
		ListAttribute:    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
		SetAttribute:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
		MapAttribute:     types.MapValueMust(types.StringType, map[string]attr.Value{"test": types.StringValue("test")}),
		NestedObject:     nestedObjectValue,
		NestedObjectList: types.ListValueMust(
			types.ObjectType{AttrTypes: nestedObjectAttrTypes},
			[]attr.Value{nestedObjectValue},
		),
		NestedObjectMap: types.MapValueMust(
			types.ObjectType{AttrTypes: nestedObjectAttrTypes},
			map[string]attr.Value{"test": nestedObjectValue},
		),
	}
}

func getDefaultGoModel() RancherDevModel {
	return RancherDevModel{
		ID:               "test",
		Identifier:       "test",
		StringAttribute:  "test",
		NumberAttribute:  big.NewFloat(1),
		Int32Attribute:   int32(123), // read only
		BoolAttribute:    true,
		Int64Attribute:   int64(123),
		Float64Attribute: 1.23,
		Float32Attribute: float32(1.23),
		ListAttribute:    []string{"test"},
		SetAttribute:     map[string]bool{"test": true},
		MapAttribute:     map[string]string{"test": "test"},
		NestedObject: NestedObject{
			StringAttribute: "test",
			NestedNestedObject: NestedNestedObject{
				StringAttribute: "test",
				BoolAttribute:   true,
			},
		},
		NestedObjectList: []NestedObject{
			{
				StringAttribute: "test",
				NestedNestedObject: NestedNestedObject{
					StringAttribute: "test",
					BoolAttribute:   true,
				},
			},
		},
		NestedObjectMap: map[string]NestedObject{
			"test": {
				StringAttribute: "test",
				NestedNestedObject: NestedNestedObject{
					StringAttribute: "test",
					BoolAttribute:   true,
				},
			},
		},
	}
}

func getPlan(ctx context.Context, m Rancher2DevResourceModel, diags *diag.Diagnostics) tfsdk.Plan {
	plan := tfsdk.Plan{
		Schema: getSchema(ctx),
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

func getState(ctx context.Context, m Rancher2DevResourceModel, diags *diag.Diagnostics) tfsdk.State {
	state := tfsdk.State{
		Schema: getSchema(ctx),
	}
	if diags.HasError() {
		return state
	}

	dgs := state.Set(ctx, m)
	if dgs.HasError() {
		diags.Append(dgs...)
	}
	return state
}

func getSchema(ctx context.Context) schema.Schema {
	emptyResource := NewRancher2DevResource()
	schemaResponseContainer := &resource.SchemaResponse{}
	emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
	return schemaResponseContainer.Schema
}
