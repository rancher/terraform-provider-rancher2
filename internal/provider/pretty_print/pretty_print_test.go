package pretty_print

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type testStruct struct {
	Exported   string
	unexported string
	IntVal     int
	BoolVal    bool
}

type nestedStruct struct {
	NestedField string
	Test        *testStruct
}

type rawMessageStruct struct {
	FieldA string
	Raw    json.RawMessage
}

type baseTypes struct {
	StringValue basetypes.StringValue
	IntValue    basetypes.Int64Value
	BoolValue   basetypes.BoolValue
	ObjectValue basetypes.ObjectValue
	ListValue   basetypes.ListValue
	MapValue    basetypes.MapValue
}

// Rancher2Dev2ResourceModel{
//   ID:         basetypes.StringValue{state:0x2, value:""},
//   APIVersion: basetypes.StringValue{state:0x2, value:"test-id"},
//   Kind:basetypes.StringValue{state:0x2, value:"string"},
//   Metadata:basetypes.ObjectValue{attributes:map[string]attr.Value{
//     "annotations":basetypes.MapValue{elements:map[string]attr.Value(nil), elementType:basetypes.StringType{}, state:0x0},
//     "creation_timestamp":basetypes.StringValue{state:0x2, value:""},
//     "deletion_grace_period_seconds":basetypes.Int64Value{state:0x2, value:0},
//     "deletion_timestamp":basetypes.StringValue{state:0x2, value:""},
//     "finalizers":basetypes.ListValue{elements:[]attr.Value(nil), elementType:basetypes.StringType{}, state:0x0},
//     "generate_name":basetypes.StringValue{state:0x2, value:""},
//     "generation":basetypes.Int64Value{state:0x2, value:0},
//     "labels":basetypes.MapValue{elements:map[string]attr.Value(nil), elementType:basetypes.StringType{}, state:0x0},
//     "managed_fields":basetypes.StringValue{state:0x2, value:""},
//     "name":basetypes.StringValue{state:0x2, value:""},
//     "namespace":basetypes.StringValue{state:0x2, value:""},
//     "owner_references":basetypes.ListValue{elements:[]attr.Value(nil), elementType:basetypes.ObjectType{AttrTypes:map[string]attr.Type{
//       "api_version":basetypes.StringType{},
//       "block_owner_deletion":basetypes.BoolType{},
//       "controller":basetypes.BoolType{},
//       "kind":basetypes.StringType{},
//       "name":basetypes.StringType{},
//       "uid":basetypes.StringType{}}
//     }, state:0x0},
//     "resource_version":basetypes.StringValue{state:0x2, value:""},
//     "self_link":basetypes.StringValue{state:0x2, value:""},
//     "uid":basetypes.StringValue{state:0x2, value:""}
//   }, attributeTypes:map[string]attr.Type{"annotations":basetypes.MapType{ElemType:basetypes.StringType{}}, "creation_timestamp":basetypes.StringType{}, "deletion_grace_period_seconds":basetypes.Int64Type{}, "deletion_timestamp":basetypes.StringType{}, "finalizers":basetypes.ListType{ElemType:basetypes.StringType{}}, "generate_name":basetypes.StringType{}, "generation":basetypes.Int64Type{}, "labels":basetypes.MapType{ElemType:basetypes.StringType{}}, "managed_fields":basetypes.StringType{}, "name":basetypes.StringType{}, "namespace":basetypes.StringType{}, "owner_references":basetypes.ListType{ElemType:basetypes.ObjectType{AttrTypes:map[string]attr.Type{"api_version":basetypes.StringType{}, "block_owner_deletion":basetypes.BoolType{}, "controller":basetypes.BoolType{}, "kind":basetypes.StringType{}, "name":basetypes.StringType{}, "uid":basetypes.StringType{}}}}, "resource_version":basetypes.StringType{}, "self_link":basetypes.StringType{}, "uid":basetypes.StringType{}}, state:0x2},
//   Spec:basetypes.ObjectValue{attributes:map[string]attr.Value{
//     "bool":basetypes.BoolValue{state:0x2, value:false},
//     "float32":basetypes.Float64Value{state:0x2, value:(*big.Float)(0x414fd4b74480)},
//     "float64":basetypes.Float64Value{state:0x2, value:(*big.Float)(0x414fd4b744b0)},
//     "int32":basetypes.Int64Value{state:0x2, value:1},
//     "int64":basetypes.Int64Value{state:0x2, value:1},
//     "list":basetypes.ListValue{elements:[]attr.Value{basetypes.StringValue{state:0x2, value:"test"}}, elementType:basetypes.StringType{}, state:0x2},
//     "map":basetypes.MapValue{elements:map[string]attr.Value{
//       "test":basetypes.StringValue{state:0x2, value:"test"}}, elementType:basetypes.StringType{}, state:0x2},
//     "number":basetypes.NumberValue{state:0x2, value:(*big.Float)(0x414fd4b74450)},
//     "object":basetypes.ObjectValue{attributes:map[string]attr.Value{
//       "string_attribute":basetypes.StringValue{state:0x2, value:"test"}
//     }, attributeTypes:map[string]attr.Type{"string_attribute":basetypes.StringType{}}, state:0x2},
//     "object_list":basetypes.ListValue{elements:[]attr.Value{
//       basetypes.ObjectValue{attributes:map[string]attr.Value{
//         "string_attribute":basetypes.StringValue{state:0x2, value:"test"}
//     }, attributeTypes:map[string]attr.Type{"string_attribute":basetypes.StringType{}}, state:0x2}}, elementType:basetypes.ObjectType{AttrTypes:map[string]attr.Type{"string_attribute":basetypes.StringType{}}}, state:0x2}, "object_map":basetypes.MapValue{elements:map[string]attr.Value{"test":basetypes.ObjectValue{attributes:map[string]attr.Value{"string_attribute":basetypes.StringValue{state:0x2, value:"test"}}, attributeTypes:map[string]attr.Type{"string_attribute":basetypes.StringType{}}, state:0x2}}, elementType:basetypes.ObjectType{AttrTypes:map[string]attr.Type{"string_attribute":basetypes.StringType{}}}, state:0x2}, "string":basetypes.StringValue{state:0x2, value:"test"}}, attributeTypes:map[string]attr.Type{"bool":basetypes.BoolType{}, "float32":basetypes.Float64Type{}, "float64":basetypes.Float64Type{}, "int32":basetypes.Int64Type{}, "int64":basetypes.Int64Type{}, "list":basetypes.ListType{ElemType:basetypes.StringType{}}, "map":basetypes.MapType{ElemType:basetypes.StringType{}}, "number":basetypes.NumberType{}, "object":basetypes.ObjectType{AttrTypes:map[string]attr.Type{"string_attribute":basetypes.StringType{}}}, "object_list":basetypes.ListType{ElemType:basetypes.ObjectType{AttrTypes:map[string]attr.Type{"string_attribute":basetypes.StringType{}}}}, "object_map":basetypes.MapType{ElemType:basetypes.ObjectType{AttrTypes:map[string]attr.Type{"string_attribute":basetypes.StringType{}}}}, "string":basetypes.StringType{}}, state:0x2}, Status:basetypes.StringValue{state:0x2, value:"\"[{\\\"status\\\":\\\"active\\\"}]\""}, APIResponses:basetypes.MapValue{elements:map[string]attr.Value(nil), elementType:basetypes.ObjectType{AttrTypes:map[string]attr.Type{"body":basetypes.StringType{}, "headers":basetypes.MapType{ElemType:basetypes.ListType{ElemType:basetypes.StringType{}}}, "status_code":basetypes.Int64Type{}}}, state:0x0}}

func TestToPrintable(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "simple struct",
			input: testStruct{
				Exported:   "hello",
				unexported: "world",
				IntVal:     123,
				BoolVal:    true,
			},
			expected: map[string]any{
				"Exported": "hello",
				"IntVal":   123,
				"BoolVal":  true,
			},
		},
		{
			name: "nested struct",
			input: nestedStruct{
				NestedField: "I am nested",
				Test: &testStruct{
					Exported:   "hello nested",
					unexported: "world nested",
					IntVal:     456,
					BoolVal:    false,
				},
			},
			expected: map[string]any{
				"NestedField": "I am nested",
				"Test": map[string]any{
					"Exported": "hello nested",
					"IntVal":   456,
					"BoolVal":  false,
				},
			},
		},
		{
			name: "nil pointer in struct",
			input: nestedStruct{
				NestedField: "nil pointer test",
				Test:        nil,
			},
			expected: map[string]any{
				"NestedField": "nil pointer test",
				"Test":        nil,
			},
		},
		{
			name:     "nil struct pointer",
			input:    (*testStruct)(nil),
			expected: nil,
		},
		{
			name:     "empty struct",
			input:    struct{}{},
			expected: map[string]any{},
		},
		{
			name:     "string value",
			input:    "a string",
			expected: "a string",
		},
		{
			name:     "integer value",
			input:    42,
			expected: 42,
		},
		{
			name:     "boolean value",
			input:    true,
			expected: true,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
		{
			name:     "slice of strings",
			input:    []string{"a", "b", "c"},
			expected: []any{"a", "b", "c"},
		},
		{
			name:     "map of string to int",
			input:    map[string]int{"one": 1, "two": 2},
			expected: map[string]any{"one": 1, "two": 2},
		},
		{
			name:     "Terraform types.String",
			input:    types.StringValue("tf string"),
			expected: "tf string",
		},
		{
			name:     "Terraform types.Int64",
			input:    types.Int64Value(99),
			expected: int64(99),
		},
		{
			name:     "Terraform types.Bool",
			input:    types.BoolValue(true),
			expected: true,
		},
		{
			name:     "Terraform types.String Unknown",
			input:    types.StringUnknown(),
			expected: "<unknown>",
		},
		{
			name:     "Terraform types.String Null",
			input:    types.StringNull(),
			expected: nil,
		},
		{
			name: "json.RawMessage with valid json",
			input: rawMessageStruct{
				FieldA: "field a",
				Raw:    json.RawMessage(`{"raw_field":"raw_value"}`),
			},
			expected: map[string]any{
				"FieldA": "field a",
				"Raw": map[string]any{
					"raw_field": "raw_value",
				},
			},
		},
		{
			name: "json.RawMessage with empty json",
			input: rawMessageStruct{
				FieldA: "field a",
				Raw:    json.RawMessage(`{}`),
			},
			expected: map[string]any{
				"FieldA": "field a",
				"Raw":    map[string]any{},
			},
		},
		{
			name: "json.RawMessage with string",
			input: rawMessageStruct{
				FieldA: "field a",
				Raw:    json.RawMessage(`"a string"`),
			},
			expected: map[string]any{
				"FieldA": "field a",
				"Raw":    "a string",
			},
		},
		{
			name: "json.RawMessage with invalid json",
			input: rawMessageStruct{
				FieldA: "field a",
				Raw:    json.RawMessage(`invalid-json`),
			},
			expected: map[string]any{
				"FieldA": "field a",
				"Raw":    "invalid-json",
			},
		},
		{
			name: "Terraform basetypes.ObjectValue",
			input: types.ObjectValueMust(
				map[string]attr.Type{"key": types.StringType},
				map[string]attr.Value{"key": types.StringValue("value")},
			),
			expected: map[string]any{
				"key": "value",
			},
		},
		{
			name: "Terraform basetypes.ListValue",
			input: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("item")},
			),
			expected: []any{"item"},
		},
		{
			name: "Terraform basetypes.MapValue",
			input: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{"key": types.StringValue("value")},
			),
			expected: map[string]any{
				"key": "value",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ToPrintable(tc.input)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Logf("have: \n%#v\nwant: \n%#v\n", actual, tc.expected)
				actualJSON, _ := json.MarshalIndent(actual, "", "  ")
				expectedJSON, _ := json.MarshalIndent(tc.expected, "", "  ")
				t.Errorf("have: \n%s\nwant: \n%s\n", string(actualJSON), string(expectedJSON))
			}
		})
	}
}

func TestPrettyPrint(t *testing.T) {
	type rawMessageStruct struct {
		FieldA string
		Raw    json.RawMessage
	}

	testCases := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "json.RawMessage with invalid json",
			input: rawMessageStruct{
				FieldA: "fieldA",
				Raw:    json.RawMessage(`invalid-json`),
			},
			expected: `{
        "FieldA": "fieldA",
        "Raw":    "invalid-json"
      }`,
		},
		{
			name: "json.RawMessage with empty string",
			input: rawMessageStruct{
				FieldA: "field a",
				Raw:    json.RawMessage(``),
			},
			expected: `{
        "FieldA": "field a",
        "Raw":    ""
      }`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := PrettyPrint(tc.input)
			var actualBuf bytes.Buffer
			if err := json.Compact(&actualBuf, []byte(actual)); err != nil {
				t.Fatalf("failed to compact actual: %v", err)
			}
			var expectedBuf bytes.Buffer
			if err := json.Compact(&expectedBuf, []byte(tc.expected.(string))); err != nil {
				t.Fatalf("failed to compact expected: %v", err)
			}
			if !bytes.Equal(actualBuf.Bytes(), expectedBuf.Bytes()) {
				t.Errorf("have: \n%#v\nwant: \n%#v\n", actual, tc.expected)
			}
		})
	}
}
