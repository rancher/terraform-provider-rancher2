package rancher2_dev

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	h "github.com/rancher/terraform-provider-rancher2/internal/provider/test_helpers"
)

const (
	apiUrl      = "https://rancher.example.com"
	endpoint    = "dev"
	apiEndpoint = apiUrl + "/" + endpoint
)

func TestRancherDevResourceMetadata(t *testing.T) {
	t.Run("Metadata function", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherDevResource
			want resource.MetadataResponse
		}{
			{"Basic test", RancherDevResource{}, resource.MetadataResponse{TypeName: "rancher2_dev"}},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				res := resource.MetadataResponse{}
				tc.fit.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "rancher2"}, &res)
				got := res
				if got != tc.want {
					t.Errorf("%+v.Metadata() is %+v; want %+v", tc.fit, got, tc.want)
				}
			})
		}
	})
}

func TestRancherDevResourceSchema(t *testing.T) {
	t.Run("Schema function", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherDevResource
			want []string
		}{
			{
				"Basic",
				RancherDevResource{},
				[]string{
					"id",
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

func TestRancherDevResourceConfig(t *testing.T) {
	t.Run("Config function", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherDevResource
			want RancherDevResource
		}{
			{
				"Basic",
				RancherDevResource{},
				RancherDevResource{},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				req := resource.ConfigureRequest{ProviderData: h.GetTestClient(t, ctx)}
				res := resource.ConfigureResponse{}
				tc.fit.Configure(ctx, req, &res)
			})
		}
	})
}

func TestRancherDevResourceCreate(t *testing.T) {
	t.Run("Create function", func(t *testing.T) {
		testCases := []struct {
			name          string
			fit           RancherDevResource
			env           map[string]string       // a k/v map of environment variables to set
			plan          RancherDevResourceModel // what is in the plan, translated to struct
			expectedState RancherDevResourceModel // what should be in the state, translated to struct
			apiRequest    c.Request               // what the request should look like as reported by the client
			apiResponse   c.Response              // the API response to inject into the client
			outcome       string                  // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// state expected to match this
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// the client should report this request object
				c.Request{
					Endpoint: apiEndpoint,
					Method:   "POST",
					Body: RancherDevModel{
						Id:              "test",
						BoolAttribute:   false,
						NumberAttribute: 1,
						StringAttribute: "test",
					},
				},
				// the response to inject into the client
				c.Response{
					StatusCode: 200,
					Headers:    map[string][]string{},
					Body: resBodyMarshall(RancherDevModel{
						Id:              "test",
						BoolAttribute:   false,
						NumberAttribute: 1,
						StringAttribute: "test",
					}),
				},
				// expected outcome
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				defer func() {
					for k := range tc.env {
						// nolint:usetesting
						os.Unsetenv(k)
					}
				}()
				for k, v := range tc.env {
					// nolint:usetesting
					os.Setenv(k, v)
				}
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10)
				client.SetResponse(tc.apiResponse)

				err := h.GetConfiguredResource(ctx, t, &tc.fit, client)
				if err != nil {
					t.Errorf("Error configuring resource: %+v", err)
				}

				req := resource.CreateRequest{Plan: h.GetPlan(t, &RancherDevResource{}, tc.plan)}
				res := resource.CreateResponse{State: h.GetState(t, &RancherDevResource{}, nil)}
				tc.fit.Create(ctx, req, &res)
				actualState := res
				expectedState := resource.CreateResponse{State: h.GetState(t, &RancherDevResource{}, tc.expectedState)}

				actualApiRequest := client.GetLastRequest()
				expectedApiRequest := tc.apiRequest

				// Always verify outcome before comparing objects
				if tc.outcome == "failure" {
					if !res.Diagnostics.HasError() {
						t.Errorf("%#v.Configure() did not return expected error diagnostics: %+v", tc.fit, pp.PrettyPrint(res.Diagnostics))
					} else {
						return
					}
				}
				if (tc.outcome == "success") && res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, pp.PrettyPrint(res.Diagnostics))
				}
				if diff := cmp.Diff(expectedApiRequest, actualApiRequest); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
				if diff := cmp.Diff(expectedState, actualState); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
}

func TestRancherDevResourceRead(t *testing.T) {
	t.Run("Read function", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherDevResource
			env                map[string]string       // a k/v map of environment variables to set
			existingState      RancherDevResourceModel // this will get injected in the read request
			expectedState      RancherDevResourceModel
			apiResponse        c.Response // this will be injected into the client
			expectedApiRequest c.Request  // the API request expected to be reported from the client
			outcome            string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				// env
				map[string]string{},
				// existing state set in read request
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// resulting state expected to match this
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// the response to inject into the client
				c.Response{
					StatusCode: 200,
					Headers:    map[string][]string{},
					Body: resBodyMarshall(RancherDevModel{
						Id:              "test",
						BoolAttribute:   false,
						NumberAttribute: 1,
						StringAttribute: "test",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/test", // add the id to the path
					Method:   "GET",
				},
				// expected outcome
				"success",
			},
			{
				"Update object",
				RancherDevResource{},
				// env
				map[string]string{},
				// existing state set in read request
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// resulting state expected to match this
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("this is different"),
				},
				// the response to inject into the client
				c.Response{
					StatusCode: 200,
					Headers:    map[string][]string{},
					Body: resBodyMarshall(RancherDevModel{
						Id:              "test",
						BoolAttribute:   false,
						NumberAttribute: 1,
						StringAttribute: "this is different",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/test", // add the id to the path
					Method:   "GET",
				},
				// expected outcome
				"success",
			},
			{
				"Failed Response",
				RancherDevResource{},
				// env
				map[string]string{},
				// existing state set in read request
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// resulting state expected to match this
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("this is different"),
				},
				// the response to inject into the client
				c.Response{
					StatusCode: 500,
					Headers:    map[string][]string{},
					Body:       []byte(""),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/test", // add the id to the path
					Method:   "GET",
				},
				// expected outcome
				"failure",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				defer func() {
					for k := range tc.env {
						// nolint:usetesting
						os.Unsetenv(k)
					}
				}()
				for k, v := range tc.env {
					// nolint:usetesting
					os.Setenv(k, v)
				}
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10)
				client.SetResponse(tc.apiResponse)

				err := h.GetConfiguredResource(ctx, t, &tc.fit, client)
				if err != nil {
					t.Errorf("Error configuring resource: %+v", err)
				}

				req := resource.ReadRequest{State: h.GetState(t, &RancherDevResource{}, tc.existingState)}
				expectedState := resource.ReadResponse{State: h.GetState(t, &RancherDevResource{}, tc.expectedState)}
				res := resource.ReadResponse{State: h.GetState(t, &RancherDevResource{}, nil)}
				tc.fit.Read(context.Background(), req, &res)
				actualState := res
				expectedApiRequest := tc.expectedApiRequest
				actualApiRequest := client.GetLastRequest()

				// Always verify outcome before comparing objects
				if tc.outcome == "failure" {
					if !res.Diagnostics.HasError() {
						t.Errorf("%#v.Configure() did not return expected error diagnostics: %+v", tc.fit, pp.PrettyPrint(res.Diagnostics))
					} else {
						return
					}
				}
				if (tc.outcome == "success") && res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() returned unexpected error diagnostics: %+v", tc.fit, pp.PrettyPrint(res.Diagnostics))
				}
				if diff := cmp.Diff(expectedState, actualState); diff != "" {
					t.Errorf("Create() mismatch (-want +got): %s", diff)
				}
				if diff := cmp.Diff(expectedApiRequest, actualApiRequest); diff != "" {
					t.Errorf("Create() mismatch (-want +got): %s", diff)
				}
			})
		}
	})
}

func TestRancherDevResourceUpdate(t *testing.T) {
	t.Run("Update function", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherDevResource
			env                map[string]string // a k/v map of environment variables to set
			plan               RancherDevResourceModel
			existingState      RancherDevResourceModel
			expectedState      RancherDevResourceModel
			apiResponse        c.Response // this will be injected into the client
			expectedApiRequest c.Request  // the API request expected to be reported from the client
			outcome            string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan set in update request
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// existing state set in update request
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// resulting state expected to match this
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: 200,
					Headers:    map[string][]string{},
					Body: resBodyMarshall(RancherDevModel{
						Id:              "test",
						BoolAttribute:   false,
						NumberAttribute: 1,
						StringAttribute: "test",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/test", // add the id to the path
					Method:   "PUT",
					Body: RancherDevModel{
						Id:              "test",
						BoolAttribute:   false,
						NumberAttribute: 1,
						StringAttribute: "test",
					},
				},
				// expected outcome
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				defer func() {
					for k := range tc.env {
						// nolint:usetesting
						os.Unsetenv(k)
					}
				}()
				for k, v := range tc.env {
					// nolint:usetesting
					os.Setenv(k, v)
				}
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10)
				client.SetResponse(tc.apiResponse)

				err := h.GetConfiguredResource(ctx, t, &tc.fit, client)
				if err != nil {
					t.Errorf("Error configuring resource: %+v", err)
				}

				req := resource.UpdateRequest{
					Plan:  h.GetPlan(t, &RancherDevResource{}, tc.plan),
					State: h.GetState(t, &RancherDevResource{}, tc.existingState),
				}
				res := resource.UpdateResponse{
					State: h.GetState(t, &RancherDevResource{}, nil),
				}
				tc.fit.Update(ctx, req, &res)

				actualState := res
				expectedState := resource.UpdateResponse{
					State: h.GetState(t, &RancherDevResource{}, tc.expectedState),
				}

				actualApiRequest := client.GetLastRequest()
				expectedApiRequest := tc.expectedApiRequest

				// always verify outcome before comparing objects
				if tc.outcome == "failure" {
					if !res.Diagnostics.HasError() {
						t.Errorf("%#v.Configure() did not return expected error diagnostics: %+v", tc.fit, pp.PrettyPrint(res.Diagnostics))
					} else {
						return
					}
				}
				if (tc.outcome == "success") && res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, pp.PrettyPrint(res.Diagnostics))
				}
				if diff := cmp.Diff(expectedState, actualState); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
				if diff := cmp.Diff(expectedApiRequest, actualApiRequest); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
}

func TestRancherDevResourceDelete(t *testing.T) {
	t.Run("Delete function", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherDevResource
			env                map[string]string // a k/v map of environment variables to set
			existingState      RancherDevResourceModel
			apiResponse        c.Response // this will be injected into the client
			expectedApiRequest c.Request  // the API request expected to be reported from the client
			outcome            string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				// env
				map[string]string{},
				// existing state
				RancherDevResourceModel{
					Id:              types.StringValue("test"),
					BoolAttribute:   types.BoolValue(false),
					NumberAttribute: types.Int64Value(1),
					StringAttribute: types.StringValue("test"),
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: 200,
					Headers:    map[string][]string{},
					Body: resBodyMarshall(RancherDevModel{
						Id:              "test",
						BoolAttribute:   false,
						NumberAttribute: 1,
						StringAttribute: "test",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/test", // add the id to the path
					Method:   "DELETE",
				},
				// expected outcome
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				defer func() {
					for k := range tc.env {
						// nolint:usetesting
						os.Unsetenv(k)
					}
				}()
				for k, v := range tc.env {
					// nolint:usetesting
					os.Setenv(k, v)
				}
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting
				ctx := h.GenerateTestContext(t, &buf, nil)

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10)
				client.SetResponse(tc.apiResponse)

				err := h.GetConfiguredResource(ctx, t, &tc.fit, client)
				if err != nil {
					t.Errorf("Error configuring resource: %+v", err)
				}

				req := resource.DeleteRequest{
					State: h.GetState(t, &RancherDevResource{}, tc.existingState),
				}
				res := resource.DeleteResponse{
					State: h.GetState(t, &RancherDevResource{}, nil),
				}
				tc.fit.Delete(ctx, req, &res)

				actualApiRequest := client.GetLastRequest()
				expectedApiRequest := tc.expectedApiRequest

				// always verify outcome before comparing objects
				if tc.outcome == "failure" {
					if !res.Diagnostics.HasError() {
						t.Errorf("%#v.Configure() did not return expected error diagnostics: %+v", tc.fit, pp.PrettyPrint(res.Diagnostics))
					} else {
						return
					}
				}
				if (tc.outcome == "success") && res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, pp.PrettyPrint(res.Diagnostics))
				}
				if diff := cmp.Diff(expectedApiRequest, actualApiRequest); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
}

// helpers.
func resBodyMarshall(obj any) []byte {
	b, _ := json.Marshal(obj)
	return b
}
