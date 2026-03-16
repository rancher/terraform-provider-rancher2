package rancher2_dev

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRancherDevModel(t *testing.T) {
	t.Run("ToResource", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherDevModel
			want RancherDevResourceModel
		}{
			{
				"Basic",
				RancherDevModel{
					Id:               "test",
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
				},
				RancherDevResourceModel{
					Id:               types.StringValue("test"),
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
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.fit.ToResourceModel(context.Background(), &diag.Diagnostics{})
				if diff := cmp.Diff(tc.want, *got, cmp.AllowUnexported(tftypes.Value{}, big.Float{})); diff != "" {
					t.Errorf("unexpected diff (-want, +got) = %s", diff)
				}
			})
		}
	})
}

func TestNestedObject(t *testing.T) {
	t.Run("ToResource", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  NestedObject
			want NestedResourceModel
		}{
			{
				"Basic",
				NestedObject{
					StringAttribute: "test",
					NestedNestedObject: NestedNestedObject{
						StringAttribute: "test",
						BoolAttribute:   true,
					},
				},
				NestedResourceModel{
					StringAttribute: types.StringValue("test"),
					NestedNestedObject: NestedNestedResourceModel{
						StringAttribute: types.StringValue("test"),
						BoolAttribute:   types.BoolValue(true),
					},
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got := NestedResourceModel{}
				tc.fit.ToResourceModel(context.Background(), &got)
				if got != tc.want {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})
}

func TestNestedNestedObject(t *testing.T) {
	t.Run("ToResource", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  NestedNestedObject
			want NestedNestedResourceModel
		}{
			{
				"Basic",
				NestedNestedObject{
					StringAttribute: "test",
					BoolAttribute:   true,
				},
				NestedNestedResourceModel{
					StringAttribute: types.StringValue("test"),
					BoolAttribute:   types.BoolValue(true),
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got := NestedNestedResourceModel{}
				tc.fit.ToResourceModel(context.Background(), &got)
				if got != tc.want {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})
}
