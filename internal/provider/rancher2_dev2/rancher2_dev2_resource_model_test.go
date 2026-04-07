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

var testAPIResponseObject = types.ObjectValueMust(
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

var testAPIResponsesMap = types.MapValueMust(
	types.ObjectType{AttrTypes: apiResponseAttrTypes},
	map[string]attr.Value{
		"create": testAPIResponseObject,
	},
)

var testFullDevResourceModel = Rancher2Dev2ResourceModel{ // Terraform resource model
	ID:           types.StringValue("test_id"),
	APIVersion:   types.StringValue("v1"),
	Kind:         types.StringValue("Rancher2Dev2"),
	Status:       types.StringValue("active"),
	Metadata:     mta.SampleMetadataTypesObject(),
	Spec:         testSpecTypesObject,
	APIResponses: testAPIResponsesMap,
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

var sampledata = `
&rancher2_dev2.Rancher2Dev2ResourceModel{
  ID:basetypes.StringValue{state:0x2, value:"test-id"}, 
  APIVersion:basetypes.StringValue{state:0x2, value:"test-id"}, 
  Kind:basetypes.StringValue{state:0x2, value:"string"}, 
  Metadata:basetypes.ObjectValue{
    attributes:map[string]attr.Value{
      "annotations":basetypes.MapValue{
        elements:map[string]attr.Value(nil), 
        elementType:basetypes.StringType{}, 
        state:0x0
      },
      "creation_timestamp":basetypes.StringValue{state:0x2, value:""}, 
      "deletion_grace_period_seconds":basetypes.Int64Value{state:0x2, value:0}, 
      "deletion_timestamp":basetypes.StringValue{state:0x2, value:""}, 
      "finalizers":basetypes.ListValue{
        elements:[]attr.Value(nil), 
        elementType:basetypes.StringType{}, 
        state:0x0
      }, 
      "generate_name":basetypes.StringValue{state:0x2, value:""}, 
      "generation":basetypes.Int64Value{state:0x2, value:0}, 
      "labels":basetypes.MapValue{
        elements:map[string]attr.Value(nil), 
        elementType:basetypes.StringType{}, 
        state:0x0
      }, 
      "managed_fields":basetypes.StringValue{state:0x2, value:""}, 
      "name":basetypes.StringValue{state:0x2, value:""}, 
      "namespace":basetypes.StringValue{state:0x2, value:""}, 
      "owner_references":basetypes.ListValue{
        elements:[]attr.Value(nil), 
        elementType:basetypes.ObjectType{
          AttrTypes:map[string]attr.Type{
            "api_version":basetypes.StringType{}, 
            "block_owner_deletion":basetypes.BoolType{}, 
            "controller":basetypes.BoolType{}, 
            "kind":basetypes.StringType{}, 
            "name":basetypes.StringType{}, 
            "uid":basetypes.StringType{}
          }
        },
        state:0x0
      }, 
      "resource_version":basetypes.StringValue{state:0x2, value:""}, 
      "self_link":basetypes.StringValue{state:0x2, value:""}, 
      "uid":basetypes.StringValue{state:0x2, value:""}
    },
    attributeTypes:map[string]attr.Type{
      "annotations":basetypes.MapType{
        ElemType:basetypes.StringType{}
      }, 
      "creation_timestamp":basetypes.StringType{}, 
      "deletion_grace_period_seconds":basetypes.Int64Type{}, 
      "deletion_timestamp":basetypes.StringType{}, 
      "finalizers":basetypes.ListType{
        ElemType:basetypes.StringType{}
      }, 
      "generate_name":basetypes.StringType{}, 
      "generation":basetypes.Int64Type{}, 
      "labels":basetypes.MapType{
        ElemType:basetypes.StringType{}
      }, 
      "managed_fields":basetypes.StringType{}, 
      "name":basetypes.StringType{}, 
      "namespace":basetypes.StringType{}, 
      "owner_references":basetypes.ListType{
        ElemType:basetypes.ObjectType{
          AttrTypes:map[string]attr.Type{
            "api_version":basetypes.StringType{}, 
            "block_owner_deletion":basetypes.BoolType{}, 
            "controller":basetypes.BoolType{}, 
            "kind":basetypes.StringType{}, 
            "name":basetypes.StringType{}, 
            "uid":basetypes.StringType{}
          }
        }
      }, 
      "resource_version":basetypes.StringType{}, 
      "self_link":basetypes.StringType{}, 
      "uid":basetypes.StringType{}
    }, 
    state:0x2
  }, 
  Spec:basetypes.ObjectValue{
    attributes:map[string]attr.Value{
      "bool":basetypes.BoolValue{state:0x2, value:false}, 
      "float32":basetypes.Float64Value{state:0x2, value:(*big.Float)(0x2fc0d6d36b0)}, 
      "float64":basetypes.Float64Value{state:0x2, value:(*big.Float)(0x2fc0d6d36e0)}, 
      "int32":basetypes.Int64Value{state:0x2, value:1}, 
      "int64":basetypes.Int64Value{state:0x2, value:1}, 
      "list":basetypes.ListValue{
        elements:[]attr.Value{
          basetypes.StringValue{state:0x2, value:"test"}
        }, elementType:basetypes.StringType{}, 
        state:0x2
      }, 
      "map":basetypes.MapValue{
        elements:map[string]attr.Value{
          "test":basetypes.StringValue{state:0x2, value:"test"}
        }, 
        elementType:basetypes.StringType{}, 
        state:0x2
      }, 
      "number":basetypes.NumberValue{state:0x2, value:(*big.Float)(0x2fc0d6d3680)}, 
      "object":basetypes.ObjectValue{
        attributes:map[string]attr.Value{
          "string_attribute":basetypes.StringValue{state:0x2, value:"test"}
        }, 
        attributeTypes:map[string]attr.Type{
          "string_attribute":basetypes.StringType{}
        }, 
        state:0x2
      }, 
      "object_list":basetypes.ListValue{
        elements:[]attr.Value{
          basetypes.ObjectValue{
            attributes:map[string]attr.Value{
              "string_attribute":basetypes.StringValue{state:0x2, value:"test"}
            }, 
            attributeTypes:map[string]attr.Type{
              "string_attribute":basetypes.StringType{}
            }, 
            state:0x2
          }
        }, 
        elementType:basetypes.ObjectType{
          AttrTypes:map[string]attr.Type{
            "string_attribute":basetypes.StringType{}
          }
        }, 
        state:0x2
      }, 
      "object_map":basetypes.MapValue{
        elements:map[string]attr.Value{
          "test":basetypes.ObjectValue{
            attributes:map[string]attr.Value{
              "string_attribute":basetypes.StringValue{state:0x2, value:"test"}
            }, 
            attributeTypes:map[string]attr.Type{
              "string_attribute":basetypes.StringType{}
            }, 
            state:0x2
          }
        }, 
        elementType:basetypes.ObjectType{
          AttrTypes:map[string]attr.Type{
            "string_attribute":basetypes.StringType{}
          }
        }, 
        state:0x2
      }, 
      "string":basetypes.StringValue{state:0x2, value:"test"}
    }, 
    attributeTypes:map[string]attr.Type{
      "bool":basetypes.BoolType{}, 
      "float32":basetypes.Float64Type{}, 
      "float64":basetypes.Float64Type{}, 
      "int32":basetypes.Int64Type{}, 
      "int64":basetypes.Int64Type{}, 
      "list":basetypes.ListType{
        ElemType:basetypes.StringType{}
      }, 
      "map":basetypes.MapType{
        ElemType:basetypes.StringType{}
      }, 
      "number":basetypes.NumberType{}, 
      "object":basetypes.ObjectType{
        AttrTypes:map[string]attr.Type{
          "string_attribute":basetypes.StringType{}
        }
      }, 
      "object_list":basetypes.ListType{
        ElemType:basetypes.ObjectType{
          AttrTypes:map[string]attr.Type{
            "string_attribute":basetypes.StringType{}
          }
        }
      }, 
      "object_map":basetypes.MapType{
        ElemType:basetypes.ObjectType{
          AttrTypes:map[string]attr.Type{
            "string_attribute":basetypes.StringType{}
          }
        }
      }, 
      "string":basetypes.StringType{}
    }, 
    state:0x2
  }, 
  Status:basetypes.StringValue{
    state:0x2, 
    value:"{\"status\":\"active\"}"
  }, 
  APIResponses:basetypes.MapValue{
    elements:map[string]attr.Value(nil), 
    elementType:basetypes.ObjectType{
      AttrTypes:map[string]attr.Type{
        "body":basetypes.StringType{}, 
        "headers":basetypes.MapType{
          ElemType:basetypes.ListType{
            ElemType:basetypes.StringType{}
          }
        }, 
        "status_code":basetypes.Int64Type{}
      }
    }, 
    state:0x0
  }
}
`
