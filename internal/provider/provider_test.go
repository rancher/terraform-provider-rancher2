package provider

import (
	"context"
	"slices"
	"testing"

	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProviderMetadata(t *testing.T) {
	testCases := []struct {
		name string
		fit  RancherProvider
		want provider.MetadataResponse
	}{
		{
			"Basic",
			RancherProvider{version: "test"},
			provider.MetadataResponse{TypeName: "rancher", Version: "test"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := provider.MetadataRequest{}
			res := provider.MetadataResponse{}
			tc.fit.Metadata(ctx, req, &res)
			got := res
			if got != tc.want {
				t.Errorf("%#v.Metadata() is %s; want %s", tc.fit, prettyPrint(got), prettyPrint(tc.want))
			}
		})
	}
}

func TestProviderSchema(t *testing.T) {
	testCases := []struct {
		name string
		fit  RancherProvider
		want []string
	}{
		{
			"Basic",
			RancherProvider{version: "test"},
			[]string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := provider.SchemaRequest{}
			res := provider.SchemaResponse{}
			tc.fit.Schema(ctx, req, &res)
			gotKeys := []string{}
			for key := range res.Schema.Attributes {
				gotKeys = append(gotKeys, key)
			}
			for _, wantKey := range tc.want {
				found := slices.Contains(gotKeys, wantKey)
				if !found {
					t.Errorf("%#v.Schema() missing expected key %s", tc.fit, wantKey)
				}
			}
		})
	}
}

// test provider configuration giving some basic values
func TestProviderConfigure(t *testing.T) {
	testCases := []struct {
		name string
		fit  RancherProvider
		have RancherProviderModel
		want string
	}{
		{
			"Basic",
			RancherProvider{version: "test"},
			RancherProviderModel{},
			"succeed",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := provider.ConfigureRequest{Config: tfsdk.Config{
				Raw: tftypes.NewValue(
					getObjectAttributeTypes(),
					getObjectAttributeValues(t, tc.have),
				),
				Schema: getSchema(),
			}}
			res := provider.ConfigureResponse{}
			tc.fit.Configure(ctx, req, &res)
			t.Logf("Configured provider: %s", prettyPrint(res))
			if (tc.want == "succeed") && res.Diagnostics.HasError() {
				t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
			}
			if (tc.want == "fail") && !res.Diagnostics.HasError() {
				t.Errorf("%#v.Configure() did not return expected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
			}
		})
	}
}

// helpers
func getSchema() schema.Schema {
	testProvider := RancherProvider{version: "test"}
	r := provider.SchemaResponse{}
	testProvider.Schema(context.Background(), provider.SchemaRequest{}, &r)
	return r.Schema
}

// getObjectAttributeTypes returns the tftypes.Object for the provider configuration by parsing the schema
func getObjectAttributeTypes() tftypes.Object {
	schema := getSchema()
	attrTypes := map[string]tftypes.Type{}
	for attrName, attr := range schema.GetAttributes() {
		attrTypes[attrName] = attr.GetType().TerraformType(context.Background())
	}
	return tftypes.Object{
		AttributeTypes: attrTypes,
	}
}

var tfsdkTagFieldMap = map[string]string{}

func getStructFieldByTfsdkTag(v reflect.Value, tagName string) (reflect.Value, error) {
	// first check if the tagName is in the map
	if fieldName, ok := tfsdkTagFieldMap[tagName]; ok {
		return v.FieldByName(fieldName), nil
	}
	// otherwise build the map and look again
	for i := 0; i < v.Type().NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("tfsdk")
		tagParts := strings.Split(tag, ",")
		tagValue := tagParts[0]
		tfsdkTagFieldMap[tagValue] = field.Name
	}
	if fieldName, ok := tfsdkTagFieldMap[tagName]; ok {
		return v.FieldByName(fieldName), nil
	}
	// if still not found, return zero Value and error
	return reflect.Value{}, fmt.Errorf("no such field with tfsdk tag %s", tagName)
}

// getObjectAttributeValues converts the RancherProviderModel struct to a map[string]tftypes.Value
// it parses the schema to get the attribute names and types so that it automatically adapts to schema changes
func getObjectAttributeValues(t *testing.T, config RancherProviderModel) map[string]tftypes.Value {
	values := map[string]tftypes.Value{}
	attributeTypes := getObjectAttributeTypes().AttributeTypes
	for attrName, attrType := range attributeTypes {
		var value interface{}
		// use reflect to get the value from the struct based on the attribute name variable
		v := reflect.ValueOf(config)
		fieldVal, err := getStructFieldByTfsdkTag(v, attrName)
		if err != nil {
			t.Fatalf("getObjectAttributeValues: %v", err)
		}
		if !fieldVal.IsValid() {
			t.Fatalf("getObjectAttributeValues: no such field %s in model", attrName)
		}

		// Dynamically call the appropriate Value method based on the tfsdk type.
		// e.g. for types.StringValue, we call ValueString(), for types.BoolValue, we call ValueBool()
		// this works for all simple types, but not for complex types like lists, maps, sets, etc.
		// try to avoid complex types in your schema
		methodName := "Value" + strings.TrimSuffix(fieldVal.Type().Name(), "Value")
		method := fieldVal.MethodByName(methodName)
		if !method.IsValid() {
			t.Fatalf("getObjectAttributeValues: no such method %s for type %s", methodName, fieldVal.Type().Name())
		}
		results := method.Call(nil)
		value = results[0].Interface()

		if value == nil {
			// for nil values, set tftypes.UnknownValue of the appropriate type
			values[attrName] = tftypes.NewValue(attrType, tftypes.UnknownValue)
			continue
		}
		values[attrName] = tftypes.NewValue(attrType, value)
	}
	return values
}

// A helper function to do the pretty-printing.
func prettyPrint(i any) string {
	s, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return "failed to pretty-print object"
	}
	return string(s)
}
