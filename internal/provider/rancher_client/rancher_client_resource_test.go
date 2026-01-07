package rancher_client

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	// "github.com/google/go-cmp/cmp/cmpopts" .
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
)

const (
	defaultId = "fake123"
)

var rancherClientResourceBooleanFields = []string{}

func TestRancherClientResourceMetadata(t *testing.T) {
	t.Run("Metadata function", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherClientResource
			want resource.MetadataResponse
		}{
			{"Basic test", RancherClientResource{}, resource.MetadataResponse{TypeName: "rancher_client"}},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				res := resource.MetadataResponse{}
				tc.fit.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "rancher"}, &res)
				got := res
				if got != tc.want {
					t.Errorf("%+v.Metadata() is %+v; want %+v", tc.fit, got, tc.want)
				}
			})
		}
	})
}

func TestRancherClientResourceSchema(t *testing.T) {
	t.Run("Schema function", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherClientResource
			want []string
		}{
			{
				"Basic",
				RancherClientResource{},
				[]string{ // want
					"id",
					"api_url",
					"ca_certs",
					"ignore_system_ca",
					"insecure",
					"max_redirects",
					"timeout",
					"access_key",
					"secret_key",
					"token_key",
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				req := resource.SchemaRequest{}
				res := resource.SchemaResponse{}
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
	})
}

func TestRancherClientResourceConfig(t *testing.T) {
	t.Run("Config function", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherClientResource
			want RancherClientResource
		}{
			{
				"Basic",
				RancherClientResource{},
				RancherClientResource{Registry: c.NewRegistry()},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				req := resource.ConfigureRequest{ProviderData: c.NewRegistry()}
				res := resource.ConfigureResponse{}
				tc.fit.Configure(ctx, req, &res)
				if tc.fit.Registry == nil {
					t.Error("Registry was not added to the client resource")
				}
			})
		}
	})
}

func TestRancherClientResourceCreate(t *testing.T) {
	t.Run("Create function", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherClientResource
			have RancherClientResourceModel
			want string // outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherClientResource{Registry: c.NewRegistry()},
				// plan
				RancherClientResourceModel{
					Id:             types.StringValue(defaultId),
					ApiURL:         types.StringValue("https://my-rancher-server.com"),
					CACerts:        types.StringValue(""),
					IgnoreSystemCA: types.BoolValue(false),
					Insecure:       types.BoolValue(false),
					MaxRedirects:   types.Int64Value(3),
					Timeout:        types.StringValue("30s"),
					AccessKey:      types.StringValue("my-access-key"),
					SecretKey:      types.StringValue("my-secret-key"),
					TokenKey:       types.StringValue("bXktYWNjZXNzLWtleTpteS1zZWNyZXQta2V5"), // base64 encoded "my-access-key:my-secret-key"
				},
				// outcome
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw: tftypes.NewValue(
							getObjectAttributeTypes(),
							getObjectAttributeValues(t, tc.have),
						),
						Schema: getSchema(),
					},
				}
				var plannedState RancherClientResourceModel
				if diags := req.Plan.Get(context.Background(), &plannedState); diags.HasError() {
					t.Errorf("Failed to get planned state: %+v", diags)
				}
				plannedId := plannedState.Id.ValueString()
				defer func() { tc.fit.Registry.Delete(plannedId) }()

				expectedState := resource.CreateResponse{
					State: tfsdk.State{
						Raw: tftypes.NewValue(
							getObjectAttributeTypes(),
							getObjectAttributeValues(t, tc.have),
						),
						Schema: getSchema(),
					},
				}
				res := resource.CreateResponse{
					State: tfsdk.State{
						Schema: getSchema(),
					},
				}
				tc.fit.Create(context.Background(), req, &res)
				// t.Logf("Resource State: %s", prettyPrint(res))

				if diff := cmp.Diff(expectedState, res); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}

				if (tc.want == "success") && res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
				}
				if (tc.want == "fail") && !res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() did not return expected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
				}
			})
		}
	})
}

// func TestRancherClientResourceRead(t *testing.T) {
// 	t.Run("Read function", func(t *testing.T) {
// 		testCases := []struct {
// 			name  string
// 			fit   RancherClientResource
// 			have  resource.ReadRequest
// 			want  resource.ReadResponse
// 			setup map[string]string
// 		}{
// 			{
// 				"Basic",
// 				RancherClientResource{client: &memoryBoilerplateClient{}},
// 				// have
// 				getRancherClientResourceReadRequest(t, map[string]string{
// 					"id": defaultId,
// 				}),
// 				// want
// 				getRancherClientResourceReadResponse(t, map[string]string{
// 					"id": defaultId,
// 				}),
// 				map[string]string{
// 					"id": defaultId,
// 				},
// 			},
// 		}
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				if err := tc.fit.client.Create(tc.setup["id"]); err != nil {
// 					t.Errorf("Error setting up: %+v", err)
// 				}
// 				defer func() {
// 					if err := tc.fit.client.Delete(tc.setup["id"]); err != nil {
// 						t.Errorf("Error tearing down: %+v", err)
// 					}
// 				}()
// 				r := getRancherClientResourceReadResponseContainer()
// 				tc.fit.Read(context.Background(), tc.have, &r)
// 				got := r
// 				if diff := cmp.Diff(tc.want, got); diff != "" {
// 					t.Errorf("Read() mismatch (-want +got):\n%+v", diff)
// 				}
// 			})
// 		}
// 	})
// }

// func TestRancherClientResourceUpdate(t *testing.T) {
// 	t.Run("Update function", func(t *testing.T) {
// 		testCases := []struct {
// 			name  string
// 			fit   RancherClientResource
// 			have  resource.UpdateRequest
// 			want  resource.UpdateResponse
// 			setup map[string]string
// 		}{
// 			{
// 				"Basic test",
// 				RancherClientResource{client: &memoryBoilerplateClient{}},
// 				// have
// 				getRancherClientResourceUpdateRequest(t, map[string]map[string]string{
// 					"priorState": {
// 						"id": defaultId,
// 					},
// 					"plan": {
// 						"id": defaultId,
// 					},
// 				}),
// 				// want
// 				getRancherClientResourceUpdateResponse(t, map[string]string{
// 					"id": defaultId,
// 				}),
// 				// setup
// 				map[string]string{
// 					"id": defaultId,
// 				},
// 			},
// 		}
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				if err := tc.fit.client.Create(tc.setup["id"]); err != nil {
// 					t.Errorf("Error setting up: %+v", err)
// 				}
// 				defer func() {
// 					if err := tc.fit.client.Delete(tc.setup["id"]); err != nil {
// 						t.Errorf("Error tearing down: %+v", err)
// 					}
// 				}()
// 				r := getRancherClientResourceUpdateResponseContainer()
// 				tc.fit.Update(context.Background(), tc.have, &r)
// 				got := r
// 				var plannedState RancherClientResourceModel
// 				if diags := tc.have.Plan.Get(context.Background(), &plannedState); diags.HasError() {
// 					t.Errorf("Failed to get planned state: %v", diags)
// 				}
// 				plannedId := plannedState.Id.ValueString()
// 				idAfterUpdate, err := tc.fit.client.Read(plannedState.Id.ValueString())
// 				if err != nil {
// 					t.Errorf("Failed to read boilerplate for update verification: %+v", err)
// 				}
// 				if idAfterUpdate != plannedId {
// 					t.Errorf("Id was not updated correctly. Got %q, want %q", idAfterUpdate, plannedId)
// 				}
// 				if diff := cmp.Diff(tc.want, got); diff != "" {
// 					t.Errorf("Update() mismatch (-want +got):\n%+v", diff)
// 				}
// 			})
// 		}
// 	})
// }

// func TestRancherClientResourceDelete(t *testing.T) {
// 	t.Run("Delete function", func(t *testing.T) {
// 		testCases := []struct {
// 			name  string
// 			fit   RancherClientResource
// 			have  resource.DeleteRequest
// 			want  resource.DeleteResponse
// 			setup map[string]string
// 		}{
// 			{
// 				"Basic test",
// 				RancherClientResource{client: &memoryBoilerplateClient{}},
// 				// have
// 				getRancherClientResourceDeleteRequest(t, map[string]string{
// 					"id": defaultId,
// 				}),
// 				// want
// 				getRancherClientResourceDeleteResponse(),
// 				// setup
// 				map[string]string{
// 					"id": defaultId,
// 				},
// 			},
// 		}
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				if err := tc.fit.client.Create(tc.setup["id"]); err != nil {
// 					t.Errorf("Error setting up: %+v", err)
// 				}
// 				r := getRancherClientResourceDeleteResponseContainer()
// 				tc.fit.Delete(context.Background(), tc.have, &r)
// 				got := r
// 				// Verify the boilerplate was actually deleted.
// 				if id, err := tc.fit.client.Read(tc.setup["id"]); err == nil || err.Error() != "some obj not found" {
// 					if err == nil {
// 						t.Errorf("Expected boilerplate to be deleted, but it still exists. Boilerplate id: %+v", id)
// 					}
// 					t.Errorf("Expected boilerplate to be deleted, but it still exists. Error: %s", err.Error())
// 				}
// 				// Verify that the boilerplate was removed from state.
// 				if diff := cmp.Diff(tc.want, got); diff != "" {
// 					t.Errorf("Update() mismatch (-want +got):\n%+v", diff)
// 				}
// 			})
// 		}
// 	})
// }

// // *** Test Helper Functions *** //
// Create.
func getRancherClientResourceCreateRequest(t *testing.T, data map[string]string) resource.CreateRequest {
	planMap := make(map[string]tftypes.Value)
	for key, value := range data {
		if slices.Contains(rancherClientResourceBooleanFields, key) { // rancherClientResourceBooleanFields is a constant
			if value == "" {
				planMap[key] = tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
			} else {
				v, err := strconv.ParseBool(value)
				if err != nil {
					t.Errorf("Error converting %s to bool %s: ", value, err.Error())
				}
				planMap[key] = tftypes.NewValue(tftypes.Bool, v)
			}
		} else {
			if value == "" {
				planMap[key] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
			} else {
				planMap[key] = tftypes.NewValue(tftypes.String, value)
			}
		}
	}
	planValue := tftypes.NewValue(getRancherClientResourceAttributeTypes(), planMap)
	return resource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw:    planValue,
			Schema: getRancherClientResourceSchema().Schema,
		},
	}
}

func getRancherClientResourceCreateResponseContainer() resource.CreateResponse {
	return resource.CreateResponse{
		State: tfsdk.State{Schema: getRancherClientResourceSchema().Schema},
	}
}

func getRancherClientResourceCreateResponse(t *testing.T, data map[string]string) resource.CreateResponse {
	stateMap := make(map[string]tftypes.Value)
	for key, value := range data {
		if slices.Contains(rancherClientResourceBooleanFields, key) { // rancherClientResourceBooleanFields is a constant
			v, err := strconv.ParseBool(value)
			if err != nil {
				t.Errorf("Error converting %s to bool %s: ", value, err.Error())
			}
			stateMap[key] = tftypes.NewValue(tftypes.Bool, v)
		} else {
			stateMap[key] = tftypes.NewValue(tftypes.String, value)
		}
	}
	stateValue := tftypes.NewValue(getRancherClientResourceAttributeTypes(), stateMap)
	return resource.CreateResponse{
		State: tfsdk.State{
			Raw:    stateValue,
			Schema: getRancherClientResourceSchema().Schema,
		},
	}
}

// // Read.
// func getRancherClientResourceReadRequest(t *testing.T, data map[string]string) resource.ReadRequest {
// 	stateMap := make(map[string]tftypes.Value)
// 	for key, value := range data {
// 		if slices.Contains(rancherClientResourceBooleanFields, key) { // rancherClientResourceBooleanFields is a constant
// 			v, err := strconv.ParseBool(value)
// 			if err != nil {
// 				t.Errorf("Error converting %s to bool %s: ", value, err.Error())
// 			}
// 			stateMap[key] = tftypes.NewValue(tftypes.Bool, v)
// 		} else {
// 			stateMap[key] = tftypes.NewValue(tftypes.String, value)
// 		}
// 	}
// 	stateValue := tftypes.NewValue(getRancherClientResourceAttributeTypes(), stateMap)
// 	return resource.ReadRequest{
// 		State: tfsdk.State{
// 			Raw:    stateValue,
// 			Schema: getRancherClientResourceSchema().Schema,
// 		},
// 	}
// }

// func getRancherClientResourceReadResponseContainer() resource.ReadResponse {
// 	return resource.ReadResponse{
// 		State: tfsdk.State{Schema: getRancherClientResourceSchema().Schema},
// 	}
// }

// func getRancherClientResourceReadResponse(t *testing.T, data map[string]string) resource.ReadResponse {
// 	stateMap := make(map[string]tftypes.Value)
// 	for key, value := range data {
// 		if slices.Contains(rancherClientResourceBooleanFields, key) { // rancherClientResourceBooleanFields is a constant
// 			v, err := strconv.ParseBool(value)
// 			if err != nil {
// 				t.Errorf("Error converting %s to bool %s: ", value, err.Error())
// 			}
// 			stateMap[key] = tftypes.NewValue(tftypes.Bool, v)
// 		} else {
// 			stateMap[key] = tftypes.NewValue(tftypes.String, value)
// 		}
// 	}
// 	stateValue := tftypes.NewValue(getRancherClientResourceAttributeTypes(), stateMap)
// 	return resource.ReadResponse{
// 		State: tfsdk.State{
// 			Raw:    stateValue,
// 			Schema: getRancherClientResourceSchema().Schema,
// 		},
// 	}
// }

// // Update.
// func getRancherClientResourceUpdateRequest(t *testing.T, data map[string]map[string]string) resource.UpdateRequest {
// 	stateMap := make(map[string]tftypes.Value)
// 	for key, value := range data["priorState"] {
// 		if slices.Contains(rancherClientResourceBooleanFields, key) { // rancherClientResourceBooleanFields is a constant
// 			v, err := strconv.ParseBool(value)
// 			if err != nil {
// 				t.Errorf("Error converting %s to bool %s: ", value, err.Error())
// 			}
// 			stateMap[key] = tftypes.NewValue(tftypes.Bool, v)
// 		} else {
// 			stateMap[key] = tftypes.NewValue(tftypes.String, value)
// 		}
// 	}
// 	priorStateValue := tftypes.NewValue(getRancherClientResourceAttributeTypes(), stateMap)

// 	planMap := make(map[string]tftypes.Value)
// 	for key, value := range data["plan"] {
// 		if slices.Contains(rancherClientResourceBooleanFields, key) { // rancherClientResourceBooleanFields is a constant
// 			if value == "" {
// 				planMap[key] = tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
// 			} else {
// 				v, err := strconv.ParseBool(value)
// 				if err != nil {
// 					t.Errorf("Error converting %s to bool %s: ", value, err.Error())
// 				}
// 				planMap[key] = tftypes.NewValue(tftypes.Bool, v)
// 			}
// 		} else {
// 			if value == "" {
// 				planMap[key] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
// 			} else {
// 				planMap[key] = tftypes.NewValue(tftypes.String, value)
// 			}
// 		}
// 	}
// 	planValue := tftypes.NewValue(getRancherClientResourceAttributeTypes(), planMap)

// 	return resource.UpdateRequest{
// 		State: tfsdk.State{
// 			Raw:    priorStateValue,
// 			Schema: getRancherClientResourceSchema().Schema,
// 		},
// 		Plan: tfsdk.Plan{
// 			Raw:    planValue,
// 			Schema: getRancherClientResourceSchema().Schema,
// 		},
// 	}
// }

// func getRancherClientResourceUpdateResponseContainer() resource.UpdateResponse {
// 	return resource.UpdateResponse{
// 		State: tfsdk.State{Schema: getRancherClientResourceSchema().Schema},
// 	}
// }

// func getRancherClientResourceUpdateResponse(t *testing.T, data map[string]string) resource.UpdateResponse {
// 	stateMap := make(map[string]tftypes.Value)
// 	for key, value := range data {
// 		if slices.Contains(rancherClientResourceBooleanFields, key) { // rancherClientResourceBooleanFields is a constant
// 			v, err := strconv.ParseBool(value)
// 			if err != nil {
// 				t.Errorf("Error converting %s to bool %s: ", value, err.Error())
// 			}
// 			stateMap[key] = tftypes.NewValue(tftypes.Bool, v)
// 		} else {
// 			stateMap[key] = tftypes.NewValue(tftypes.String, value)
// 		}
// 	}
// 	stateValue := tftypes.NewValue(getRancherClientResourceAttributeTypes(), stateMap)
// 	return resource.UpdateResponse{
// 		State: tfsdk.State{
// 			Raw:    stateValue,
// 			Schema: getRancherClientResourceSchema().Schema,
// 		},
// 	}
// }

// // Delete.
// func getRancherClientResourceDeleteRequest(t *testing.T, data map[string]string) resource.DeleteRequest {
// 	stateMap := make(map[string]tftypes.Value)
// 	for key, value := range data {
// 		if slices.Contains(rancherClientResourceBooleanFields, key) { // rancherClientResourceBooleanFields is a constant
// 			v, err := strconv.ParseBool(value)
// 			if err != nil {
// 				t.Errorf("Error converting %s to bool %s: ", value, err.Error())
// 			}
// 			stateMap[key] = tftypes.NewValue(tftypes.Bool, v)
// 		} else {
// 			stateMap[key] = tftypes.NewValue(tftypes.String, value)
// 		}
// 	}
// 	stateValue := tftypes.NewValue(getRancherClientResourceAttributeTypes(), stateMap)
// 	return resource.DeleteRequest{
// 		State: tfsdk.State{
// 			Raw:    stateValue,
// 			Schema: getRancherClientResourceSchema().Schema,
// 		},
// 	}
// }

// func getRancherClientResourceDeleteResponseContainer() resource.DeleteResponse {
// 	// A delete response does not need a schema as it results in a null state.
// 	return resource.DeleteResponse{}
// }

// func getRancherClientResourceDeleteResponse() resource.DeleteResponse {
// 	return resource.DeleteResponse{
// 		State: tfsdk.State{
// 			Raw:    tftypes.Value{},
// 			Schema: nil,
// 		},
// 	}
// }

// // The helpers helpers.
func getRancherClientResourceAttributeTypes() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{},
	}
}
func getRancherClientResourceSchema() *resource.SchemaResponse {
	var testResource RancherClientResource
	r := &resource.SchemaResponse{}
	testResource.Schema(context.Background(), resource.SchemaRequest{}, r)
	return r
}
func getSchema() schema.Schema {
	testResource := RancherClientResource{}
	r := resource.SchemaResponse{}
	testResource.Schema(context.Background(), resource.SchemaRequest{}, &r)
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

// getObjectAttributeValues converts the RancherClientResourceModel struct to a map[string]tftypes.Value
// it parses the schema to get the attribute names and types so that it automatically adapts to schema changes
func getObjectAttributeValues(t *testing.T, config RancherClientResourceModel) map[string]tftypes.Value {
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

//// OLD FUNCTIONS BELOW
// // test provider configuration giving some basic values
// func TestProviderConfigure(t *testing.T) {
//   testCases := []struct {
//     name string
//     fit  RancherProvider
//     have RancherProviderModel
//     want string
//   }{
//     {
//       "Basic",
//       RancherProvider{version: "test"},
//       RancherProviderModel{ // have
//         ApiURL:         types.StringValue("https://my-rancher-server.com"),
//         CACerts:        types.StringValue(""),
//         IgnoreSystemCA: types.BoolValue(false),
//         Insecure:       types.BoolValue(false),
//         MaxRedirects:   types.Int64Value(3),
//         Timeout:        types.StringValue("30s"),
//         AccessKey:      types.StringValue("my-access-key"),
//         SecretKey:      types.StringValue("my-secret-key"),
//         TokenKey:       types.StringValue(""),
//       },
//       "succeed",
//     },
//     {
//       "Token",
//       RancherProvider{version: "test"},
//       RancherProviderModel{ // have
//         ApiURL:         types.StringValue("https://my-rancher-server.com"),
//         CACerts:        types.StringValue(""),
//         IgnoreSystemCA: types.BoolValue(false),
//         Insecure:       types.BoolValue(false),
//         MaxRedirects:   types.Int64Value(3),
//         Timeout:        types.StringValue("30s"),
//         AccessKey:      types.StringValue(""),
//         SecretKey:      types.StringValue(""),
//         TokenKey:       types.StringValue("my-token-key"),
//       },
//       "succeed",
//     },
//     {
//       "Missing Credentials",
//       RancherProvider{version: "test"},
//       RancherProviderModel{ // have
//         ApiURL:         types.StringValue("https://my-rancher-server.com"),
//         CACerts:        types.StringValue(""),
//         IgnoreSystemCA: types.BoolValue(false),
//         Insecure:       types.BoolValue(false),
//         MaxRedirects:   types.Int64Value(3),
//         Timeout:        types.StringValue("30s"),
//         AccessKey:      types.StringValue(""),
//         SecretKey:      types.StringValue(""),
//         TokenKey:       types.StringValue(""),
//       },
//       "fail",
//     },
//     {
//       "Invalid ApiURL",
//       RancherProvider{version: "test"},
//       RancherProviderModel{ // have
//         ApiURL:         types.StringValue("my-rancher-server.com"),
//         CACerts:        types.StringValue(""),
//         IgnoreSystemCA: types.BoolValue(false),
//         Insecure:       types.BoolValue(false),
//         MaxRedirects:   types.Int64Value(3),
//         Timeout:        types.StringValue("30s"),
//         AccessKey:      types.StringValue(""),
//         SecretKey:      types.StringValue(""),
//         TokenKey:       types.StringValue("my-token-key"),
//       },
//       "fail",
//     },
//     {
//       "IP address ApiURL",
//       RancherProvider{version: "test"},
//       RancherProviderModel{ // have
//         ApiURL:         types.StringValue("https://192.168.1.1"),
//         CACerts:        types.StringValue(""),
//         IgnoreSystemCA: types.BoolValue(false),
//         Insecure:       types.BoolValue(false),
//         MaxRedirects:   types.Int64Value(3),
//         Timeout:        types.StringValue("30s"),
//         AccessKey:      types.StringValue(""),
//         SecretKey:      types.StringValue(""),
//         TokenKey:       types.StringValue("my-token-key"),
//       },
//       "succeed",
//     },
//     {
//       "Insecure ApiURL without insecure flag",
//       RancherProvider{version: "test"},
//       RancherProviderModel{ // have
//         ApiURL:         types.StringValue("http://192.168.1.1"),
//         CACerts:        types.StringValue(""),
//         IgnoreSystemCA: types.BoolValue(false),
//         Insecure:       types.BoolValue(false),
//         MaxRedirects:   types.Int64Value(3),
//         Timeout:        types.StringValue("30s"),
//         AccessKey:      types.StringValue(""),
//         SecretKey:      types.StringValue(""),
//         TokenKey:       types.StringValue("my-token-key"),
//       },
//       "fail",
//     },
//     {
//       "Secure ApiURL with insecure flag",
//       RancherProvider{version: "test"},
//       RancherProviderModel{ // have
//         ApiURL:         types.StringValue("https://192.168.1.1"),
//         CACerts:        types.StringValue(""),
//         IgnoreSystemCA: types.BoolValue(false),
//         Insecure:       types.BoolValue(true),
//         MaxRedirects:   types.Int64Value(3),
//         Timeout:        types.StringValue("30s"),
//         AccessKey:      types.StringValue(""),
//         SecretKey:      types.StringValue(""),
//         TokenKey:       types.StringValue("my-token-key"),
//       },
//       "succeed", // insecure flag should not cause failure for https URLs
//       // This allows users to set insecure flag for https URLs if they want to skip cert validation
//     },
//   }
//   for _, tc := range testCases {
//     t.Run(tc.name, func(t *testing.T) {
//       ctx := context.Background()
//       req := provider.ConfigureRequest{ Config: tfsdk.Config{
//         Raw: tftypes.NewValue(
//           getObjectAttributeTypes(),
//           getObjectAttributeValues(t, tc.have),
//         ),
//         Schema: getSchema(),
//       }}
//       res := provider.ConfigureResponse{}
//       tc.fit.Configure(ctx, req, &res)
//       t.Logf("Configured provider: %s", prettyPrint(res))
//       if (tc.want == "succeed") && res.Diagnostics.HasError() {
//         t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
//       }
//       if (tc.want == "fail") && !res.Diagnostics.HasError() {
//         t.Errorf("%#v.Configure() did not return expected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
//       }
//     })
//   }
// }
