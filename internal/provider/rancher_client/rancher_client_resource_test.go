package rancher_client

import (
	"context"
	// "slices".
	// "strconv".
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	// "github.com/hashicorp/terraform-plugin-framework/tfsdk".
	// "github.com/hashicorp/terraform-plugin-go/tftypes".
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
			want resource.SchemaResponse
		}{
			{"Basic test", RancherClientResource{}, *getRancherClientResourceSchema()},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				r := resource.SchemaResponse{}
				tc.fit.Schema(context.Background(), resource.SchemaRequest{}, &r)
				got := r
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("Schema() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
}

// func TestRancherClientResourceCreate(t *testing.T) {
// 	t.Run("Create function", func(t *testing.T) {
// 		testCases := []struct {
// 			name string
// 			fit  RancherClientResource
// 			have resource.CreateRequest
// 			want resource.CreateResponse
// 		}{
// 			{
// 				"Basic",
// 				RancherClientResource{client: &memoryBoilerplateClient{}},
// 				// have
// 				getRancherClientResourceCreateRequest(t, map[string]string{
// 					"id": defaultId,
// 				}),
// 				// want
// 				getRancherClientResourceCreateResponse(t, map[string]string{
// 					"id": "fake123",
// 				}),
// 			},
// 		}
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				var plannedState RancherClientResourceModel
// 				if diags := tc.have.Plan.Get(context.Background(), &plannedState); diags.HasError() {
// 					t.Errorf("Failed to get planned state: %+v", diags)
// 				}
// 				plannedId := plannedState.Id.ValueString()
// 				r := getRancherClientResourceCreateResponseContainer()
// 				// run the resource's Create command
// 				tc.fit.Create(context.Background(), tc.have, &r)
// 				defer func() {
// 					// run the client's Delete function
// 					if err := tc.fit.client.Delete(plannedId); err != nil {
// 						t.Errorf("Error cleaning up: %+v", err)
// 					}
// 				}()
// 				got := r
// 				if diff := cmp.Diff(tc.want, got); diff != "" {
// 					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
// 				}
// 			})
// 		}
// 	})
// }
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
// // Create.
// func getRancherClientResourceCreateRequest(t *testing.T, data map[string]string) resource.CreateRequest {
// 	planMap := make(map[string]tftypes.Value)
// 	for key, value := range data {
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
// 	return resource.CreateRequest{
// 		Plan: tfsdk.Plan{
// 			Raw:    planValue,
// 			Schema: getRancherClientResourceSchema().Schema,
// 		},
// 	}
// }

// func getRancherClientResourceCreateResponseContainer() resource.CreateResponse {
// 	return resource.CreateResponse{
// 		State: tfsdk.State{Schema: getRancherClientResourceSchema().Schema},
// 	}
// }

// func getRancherClientResourceCreateResponse(t *testing.T, data map[string]string) resource.CreateResponse {
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
// 	return resource.CreateResponse{
// 		State: tfsdk.State{
// 			Raw:    stateValue,
// 			Schema: getRancherClientResourceSchema().Schema,
// 		},
// 	}
// }

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
// func getRancherClientResourceAttributeTypes() tftypes.Object {
// 	return tftypes.Object{
// 		AttributeTypes: map[string]tftypes.Type{
// 			"id": tftypes.String,
// 		},
// 	}
// }

func getRancherClientResourceSchema() *resource.SchemaResponse {
	var testResource RancherClientResource
	r := &resource.SchemaResponse{}
	testResource.Schema(context.Background(), resource.SchemaRequest{}, r)
	return r
}
