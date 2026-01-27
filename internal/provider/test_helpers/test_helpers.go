package test_helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
)

func GenerateTestContext(t *testing.T, buf *bytes.Buffer, ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return tflogtest.RootLogger(ctx, buf)
}

var logLevels = map[string]int{
	"TRACE": 0,
	"DEBUG": 1,
	"INFO":  2,
	"WARN":  3,
	"ERROR": 4,
}

func PrintLog(t *testing.T, buf *bytes.Buffer, level ...string) {
	minLevel := -1
	if len(level) > 0 && level[0] != "" {
		var ok bool
		minLevel, ok = logLevels[strings.ToUpper(level[0])]
		if !ok {
			t.Logf("invalid log level: %s", level[0])
			return
		}
	}

	for line := range strings.SplitSeq(buf.String(), "\n") {
		if line == "" {
			continue
		}
		var logEntry map[string]any
		err := json.Unmarshal([]byte(line), &logEntry)
		if err != nil {
			t.Logf("failed to unmarshal log line: %s", err)
			continue
		}

		if minLevel != -1 {
			logLevelStr, ok := logEntry["@level"].(string)
			if !ok {
				continue // Skip logs without a level if filtering
			}

			logLevel, ok := logLevels[strings.ToUpper(logLevelStr)]
			if !ok {
				continue // Skip unknown log levels if filtering
			}

			// fmt.Printf("Log level detected: %d\n", logLevel)
			// fmt.Printf("Min level set: %d\n", minLevel)

			if logLevel < minLevel {
				continue
			}
		}

		if msg, ok := logEntry["@message"]; ok {
			t.Log(msg)
		}
	}
}

// The "typ" should be the resource you want to model, eg. Rancher2<Resource>
// The "model" should be the model struct for the resource, eg. Rancher2<Resource>Model
//   - send "nil" as the model to get a data free struct
func GetPlan(t *testing.T, typ resource.Resource, model any) tfsdk.Plan {
	plan := tfsdk.Plan{}
	getTfsdk(t, typ, model, &plan)
	return plan
}

// The "typ" should be the resource you want to model, eg. Rancher<Resource>
// The "model" should be the model struct for the resource, eg. Rancher<Resource>Model
//   - send "nil" as the model to get a data free struct
func GetState(t *testing.T, typ resource.Resource, model any) tfsdk.State {
	state := tfsdk.State{}
	getTfsdk(t, typ, model, &state)
	return state
}

// The "typ" should be the rprovider you want to model, eg. &RancherProvider{}
// The "model" should be the model struct for the resource, eg. RancherProviderModel{}
//   - send "nil" as the model to get a data free struct
func GetConfig(t *testing.T, typ provider.Provider, model any) tfsdk.Config {
	config := tfsdk.Config{}
	getTfsdk(t, typ, model, &config)
	return config
}

// rsc is a pointer to the resource you want to be configured, eg. *RancherDevResource
// client is the data sent to the resource, this is what the provider configure outputs
// this is where you can inject a custom client for testing
func GetConfiguredResource(ctx context.Context, t *testing.T, rsc resource.ResourceWithConfigure, client c.Client) error {
	req := resource.ConfigureRequest{
		ProviderData: client,
	}
	res := resource.ConfigureResponse{}
	rsc.Configure(ctx, req, &res)
	if res.Diagnostics.HasError() {
		return fmt.Errorf("error configuring resource: %s", res.Diagnostics.Errors())
	}
	return nil
}

func GetTestClient(t *testing.T, ctx context.Context) *c.TestClient {
	return c.NewTestClient(ctx, "https://rancher.example.com", "", false, false, 30, 10)
}

// Internals below here

// Reflect is used a lot in this pkg, mostly for serializing simple structs from and into tfsdk containers

// The "typ" should be the resource you want to model, eg. Rancher2<Resource>
// The "model" should be the model struct for the resource, eg. Rancher2<Resource>Model
//   - send "nil" as the model to get a data free struct
//   - not a pointer, the actual struct
//
// The "container" should be a pointer to a tfsdk.State, tfsdk.Plan, or tfsdk.Config to fill out.
func getTfsdk(t *testing.T, typ any, model any, container any) any {
	v := reflect.ValueOf(container)
	if v.Kind() != reflect.Ptr {
		t.Fatalf("container must be a pointer, got %T", container)
	}
	if v.IsNil() {
		t.Fatalf("container must not be nil")
	}
	if v.Elem().Kind() != reflect.Struct {
		t.Fatalf("container must point to a struct, got %T", container)
	}
	elem := v.Elem()

	// Set Schema
	schema := getSchema(t, typ)
	fSchema := elem.FieldByName("Schema")
	if fSchema.IsValid() && fSchema.CanSet() {
		fSchema.Set(reflect.ValueOf(schema))
	}
	// tfsdk.Container{
	//   Schema: schema.Schema{
	//     Attributes: map[string]schema.Attribute{}...
	//   },
	// }

	if model != nil {
		// m := reflect.ValueOf(model)
		// if m.Kind() != reflect.Ptr {
		// 	t.Fatalf("model must be a pointer, got %T", model)
		// }
		// // At first glance this seems redundant, but it is checking that the pointer doesn't reference a nil object.
		// if m.IsNil() {
		// 	t.Fatalf("model must not be nil")
		// }
		// if m.Elem().Kind() != reflect.Struct {
		// 	t.Fatalf("model must point to a struct, got %T", model)
		// }
		raw := tftypes.NewValue(
			getAttributeTypes(t, typ),
			getObjectAttributeValues(t, model, typ),
		)
		fRaw := elem.FieldByName("Raw")
		if fRaw.IsValid() && fRaw.CanSet() {
			fRaw.Set(reflect.ValueOf(raw))
		}
		// tfsdk.Container{
		//   Schema: schema.Schema{
		//     Attributes: map[string]schema.Attribute{}...
		//   },
		//   Raw: tftypes.Value{
		//     AttributeTypes: map[string]tftypes.Type{}...
		//   }
		// }
	}
	return container
}

func getSchema(t *testing.T, typ any) any {
	switch v := typ.(type) {
	case resource.Resource:
		rsr := resource.SchemaResponse{}
		v.Schema(context.Background(), resource.SchemaRequest{}, &rsr)
		return rsr.Schema
	case provider.Provider:
		psr := provider.SchemaResponse{}
		v.Schema(context.Background(), provider.SchemaRequest{}, &psr)
		return psr.Schema
	default:
		t.Fatalf("getSchema: unsupported type %T", v)
		return nil // unreachable
	}
}

func getAttributeTypes(t *testing.T, typ any) tftypes.Object {
	schema := getSchema(t, typ)

	v := reflect.ValueOf(schema)
	m := v.MethodByName("GetAttributes")
	if !m.IsValid() {
		t.Fatalf("getAttributeTypes: no GetAttributes method on %T", schema)
	}

	result := m.Call(nil)
	if len(result) == 0 {
		t.Fatalf("getAttributeTypes: GetAttributes method returned no value")
	}
	attributes := result[0].Interface()

	return tftypes.Object{
		AttributeTypes: extractAttributeTypes(t, attributes),
	}
}

func extractAttributeTypes(t *testing.T, attributes any) map[string]tftypes.Type {
	attrTypes := map[string]tftypes.Type{}

	v := reflect.ValueOf(attributes)
	if v.Kind() != reflect.Map {
		t.Fatalf("extractAttributeTypes: attributes is not a map, but %s", v.Kind())
	}

	iter := v.MapRange()
	for iter.Next() {
		attrName := iter.Key().String()
		attr := iter.Value() // reflect.Value of an Attribute struct

		getTypeMethod := attr.MethodByName("GetType")
		if !getTypeMethod.IsValid() {
			t.Fatalf("extractAttributeTypes: no GetType method on attribute")
		}

		getTypeResult := getTypeMethod.Call(nil)
		if len(getTypeResult) == 0 {
			t.Fatalf("extractAttributeTypes: GetType method returned no value")
		}
		typeVal := getTypeResult[0] // reflect.Value of a types.Type

		terraformTypeMethod := typeVal.MethodByName("TerraformType")
		if !terraformTypeMethod.IsValid() {
			t.Fatalf("extractAttributeTypes: no TerraformType method on type")
		}

		ctx := reflect.ValueOf(context.Background())
		terraformTypeResult := terraformTypeMethod.Call([]reflect.Value{ctx})
		if len(terraformTypeResult) == 0 {
			t.Fatalf("extractAttributeTypes: TerraformType method returned no value")
		}

		tfType, ok := terraformTypeResult[0].Interface().(tftypes.Type)
		if !ok {
			t.Fatalf("extractAttributeTypes: TerraformType method returned value is not a tftypes.Type")
		}
		attrTypes[attrName] = tfType
	}

	return attrTypes
}

// getObjectAttributeValues converts the struct to a map[string]tftypes.Value.
// it parses the schema to get the attribute names and types so that it automatically adapts to schema changes.
// config can be a struct or a pointer to a struct
func getObjectAttributeValues(t *testing.T, config any, typ any) map[string]tftypes.Value {
	attributeTypes := getAttributeTypes(t, typ).AttributeTypes

	values := map[string]tftypes.Value{}
	for attrName, attrType := range attributeTypes {
		// use reflect to get the value from the struct based on the attribute name variable.
		v := reflect.ValueOf(config)
		fieldVal, err := getStructFieldByTfsdkTag(v, attrName)
		if err != nil {
			t.Fatalf("getObjectAttributeValues: %v", err)
		}
		if !fieldVal.IsValid() {
			t.Fatalf("getObjectAttributeValues: no such field %s in model", attrName)
		}

		// Check for null and unknown values first
		if fieldVal.MethodByName("IsNull").Call(nil)[0].Bool() {
			values[attrName] = tftypes.NewValue(attrType, nil)
			continue
		}
		if fieldVal.MethodByName("IsUnknown").Call(nil)[0].Bool() {
			values[attrName] = tftypes.NewValue(attrType, tftypes.UnknownValue)
			continue
		}

		// Dynamically call the appropriate Value method based on the tfsdk type.
		// e.g. for types.StringValue, we call ValueString(), for types.BoolValue, we call ValueBool().
		// this works for all simple types, but not for complex types like lists, maps, sets, etc.
		// try to avoid complex types in your schema.
		methodName := "Value" + strings.TrimSuffix(fieldVal.Type().Name(), "Value")
		method := fieldVal.MethodByName(methodName)
		if !method.IsValid() {
			t.Fatalf("getObjectAttributeValues: no such method %s for type %s", methodName, fieldVal.Type().Name())
		}
		results := method.Call(nil)
		value := results[0].Interface()

		values[attrName] = tftypes.NewValue(attrType, value)
	}
	return values
}

var tfsdkTagFieldMap = map[string]string{}

// v can be a struct or a pointer to a struct
func getStructFieldByTfsdkTag(v reflect.Value, tagName string) (reflect.Value, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("getStructFieldByTfsdkTag: expected a struct or pointer to a struct, got %s", v.Kind())
	}
	if fieldName, ok := tfsdkTagFieldMap[tagName]; ok {
		return v.FieldByName(fieldName), nil
	}
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
	return reflect.Value{}, fmt.Errorf("no such field with tfsdk tag %s", tagName)
}
