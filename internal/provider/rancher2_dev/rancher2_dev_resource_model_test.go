package rancher2_dev

import (
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
)

func TestRancherDevResourceModel(t *testing.T) {
	t.Run("ToGoModel", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherDevResourceModel
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
				got := tc.fit.ToGoModel(context.Background())
				if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreUnexported(big.Float{})); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
	t.Run("ToPlan", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherDevResourceModel
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
				model := getDefaultResourceModel()
				got := model.ToPlan(context.Background(), &diag.Diagnostics{})

				if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{}, big.Float{})); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})
	t.Run("ToState", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherDevResourceModel
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
				model := getDefaultResourceModel()
				got := model.ToState(context.Background(), &diag.Diagnostics{})

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
					NestedNestedObject: NestedNestedResourceModel{
						StringAttribute: types.StringValue("test"),
						BoolAttribute:   types.BoolValue(true),
					},
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
				got := NestedObject{}
				tc.fit.ToGoModel(context.Background(), &got)
				if got != tc.want {
					t.Errorf("got %v, want %v", got, tc.want)
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
				got := NestedNestedObject{}
				tc.fit.ToGoModel(context.Background(), &got)
				if got != tc.want {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})
}

// helpers.
func getDefaultResourceModel() RancherDevResourceModel {
	return RancherDevResourceModel{
		Id:               types.StringValue("test"),
		UserToken:        types.StringValue(""),
		BoolAttribute:    types.BoolValue(true),
		NumberAttribute:  types.NumberValue(big.NewFloat(1.23)),
		Int64Attribute:   types.Int64Value(123),
		Int32Attribute:   types.Int32Value(123),
		Float64Attribute: types.Float64Value(1.23),
		Float32Attribute: types.Float32Value(1.23),
		StringAttribute:  types.StringValue("test"),
		ListAttribute:    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
		SetAttribute:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
		MapAttribute:     types.MapValueMust(types.StringType, map[string]attr.Value{"test": types.StringValue("test")}),
		NestedObject: NestedResourceModel{
			StringAttribute: types.StringValue("test"),
			NestedNestedObject: NestedNestedResourceModel{
				StringAttribute: types.StringValue("test"),
				BoolAttribute:   types.BoolValue(true),
			},
		},
		NestedObjectList: []NestedResourceModel{
			{
				StringAttribute: types.StringValue("test"),
				NestedNestedObject: NestedNestedResourceModel{
					StringAttribute: types.StringValue("test"),
					BoolAttribute:   types.BoolValue(true),
				},
			},
		},
		NestedObjectMap: map[string]NestedResourceModel{
			"test": {
				StringAttribute: types.StringValue("test"),
				NestedNestedObject: NestedNestedResourceModel{
					StringAttribute: types.StringValue("test"),
					BoolAttribute:   types.BoolValue(true),
				},
			},
		},
	}
}

func getDefaultGoModel() RancherDevModel {
	return RancherDevModel{
		Id:               "test",
		UserToken:        "",
		BoolAttribute:    true,
		NumberAttribute:  big.NewFloat(1.23),
		Int64Attribute:   123,
		Int32Attribute:   123,
		Float64Attribute: 1.23,
		Float32Attribute: 1.23,
		StringAttribute:  "test",
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

func getPlan(ctx context.Context, m RancherDevResourceModel, diags *diag.Diagnostics) tfsdk.Plan {
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

func getState(ctx context.Context, m RancherDevResourceModel, diags *diag.Diagnostics) tfsdk.State {
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
	emptyResource := NewRancherDevResource()
	schemaResponseContainer := &resource.SchemaResponse{}
	emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
	return schemaResponseContainer.Schema
}
