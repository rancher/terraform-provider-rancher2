package rancher2_dev

import (
	"bytes"
	"context"
	"encoding/json"

	"math/big"
	"net/http"
	"os"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	h "github.com/rancher/terraform-provider-rancher2/internal/provider/test_helpers"
)

const (
	apiUrl      = "https://rancher.example.com"
	endpoint    = "dev"
	apiEndpoint = apiUrl + "/" + endpoint
	defaultId   = "dev-test"
)

func TestRancherDevResource(t *testing.T) {
	t.Run("Metadata", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherDevResource
			want resource.MetadataResponse
		}{
			{
				"Basic",
				RancherDevResource{},
				resource.MetadataResponse{TypeName: "rancher2_dev_resource"},
			},
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
	t.Run("Schema", func(t *testing.T) {
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
					"bool_attribute",
					"number_attribute",
					"int64_attribute",
					"int32_attribute",
					"float64_attribute",
					"float32_attribute",
					"string_attribute",
					"list_attribute",
					"set_attribute",
					"map_attribute",
					"nested_object",
					"nested_object_list",
					"nested_object_map",
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
	t.Run("Config", func(t *testing.T) {
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
	t.Run("Create", func(t *testing.T) {
		testCases := []struct {
			name          string
			fit           RancherDevResource
			env           map[string]string // a k/v map of environment variables to set
			plan          RancherDevModel   // what is in the plan, translated to struct
			expectedState RancherDevModel   // what should be in the state, translated to struct
			apiRequest    c.Request         // what the request should look like as reported by the client
			apiResponse   c.Response        // the API response to inject into the client
			outcome       string            // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan
				RancherDevModel{
					Id:               defaultId,
					BoolAttribute:    false,
					NumberAttribute:  big.NewFloat(1.0),
					Int64Attribute:   1,
					Int32Attribute:   1,
					Float64Attribute: 1,
					Float32Attribute: 1,
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
				// state expected to match this
				RancherDevModel{
					Id:               defaultId,
					BoolAttribute:    false,
					NumberAttribute:  big.NewFloat(1.0),
					Int64Attribute:   1,
					Int32Attribute:   1,
					Float64Attribute: 1,
					Float32Attribute: 1,
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
				// the client should report this request object
				c.Request{
					Endpoint: apiEndpoint,
					Method:   "POST",
					Body: RancherDevModel{
						Id:               defaultId,
						BoolAttribute:    false,
						NumberAttribute:  big.NewFloat(1.0),
						Int64Attribute:   1,
						Int32Attribute:   1,
						Float64Attribute: 1,
						Float32Attribute: 1,
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
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(RancherDevModel{
						Id:               defaultId,
						BoolAttribute:    false,
						NumberAttribute:  big.NewFloat(1.0),
						Int64Attribute:   1,
						Int32Attribute:   1,
						Float64Attribute: 1,
						Float32Attribute: 1,
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
					}),
				},
				// expected outcome
				"success",
			},
			{
				"API Conflict",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan
				RancherDevModel{
					Id: defaultId,
				},
				// state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// the client should report this request object
				c.Request{
					Endpoint: apiEndpoint,
					Method:   "POST",
					Body: RancherDevModel{
						Id: defaultId,
					},
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusConflict,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(c.ErrorResponse{
						Status:  "409",
						Message: "resource already exists",
					}),
				},
				// expected outcome
				"failure",
			},
			{
				"Server Error",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan
				RancherDevModel{
					Id: defaultId,
				},
				// state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// the client should report this request object
				c.Request{
					Endpoint: apiEndpoint,
					Method:   "POST",
					Body: RancherDevModel{
						Id: "test",
					},
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusInternalServerError,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(c.ErrorResponse{
						Status:  "500",
						Message: "something went wrong",
					}),
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
					t.Fatalf("error configuring resource: %+v", err)
				}
				dgs := diag.Diagnostics{}
				plan := tc.plan.ToResourceModel(ctx, &dgs).ToPlan(ctx, &dgs)
				if dgs.HasError() {
					t.Fatalf("error generating plan: %s", pp.PrettyPrint(dgs))
				}
				req := resource.CreateRequest{Plan: plan}

				// get empty state
				emptyResource := NewRancherDevResource()
				schemaResponseContainer := &resource.SchemaResponse{}
				emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
				state := tfsdk.State{
					Schema: schemaResponseContainer.Schema,
				}

				res := resource.CreateResponse{State: state}
				tc.fit.Create(ctx, req, &res)
				actualState := res

				state = tc.expectedState.ToResourceModel(ctx, &dgs).ToState(ctx, &dgs)
				if dgs.HasError() {
					t.Fatalf("error generating expected state: %s", pp.PrettyPrint(dgs))
				}
				expectedState := resource.CreateResponse{State: state}

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

				if diff := cmp.Diff(expectedApiRequest, actualApiRequest, cmpopts.IgnoreUnexported(big.Float{})); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
				if diff := cmp.Diff(expectedState, actualState, cmpopts.IgnoreUnexported(big.Float{})); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
	t.Run("Read", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherDevResource
			env                map[string]string // a k/v map of environment variables to set
			existingState      RancherDevModel   // this will get injected in the read request
			expectedState      RancherDevModel
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
				RancherDevModel{
					Id: defaultId,
				},
				// resulting state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(RancherDevModel{
						Id: defaultId,
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
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
				RancherDevModel{
					Id: defaultId,
				},
				// resulting state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(RancherDevModel{
						Id: defaultId,
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
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
				RancherDevModel{
					Id: defaultId,
				},
				// resulting state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusInternalServerError,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(c.ErrorResponse{
						Status:  "500",
						Message: "something went wrong",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "GET",
				},
				// expected outcome
				"failure",
			},
			{
				"Unmanaged API data",
				RancherDevResource{},
				// env
				map[string]string{},
				// existing state set in read request
				RancherDevModel{
					Id: defaultId,
				},
				// resulting state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(struct {
						Id                 string `json:"id"`
						UntrackedAttribute string `json:"untracked_attribute,omitempty"`
					}{
						Id:                 defaultId,
						UntrackedAttribute: "untracked",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "GET",
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
					t.Errorf("error configuring resource: %+v", err)
				}
				var dgs diag.Diagnostics
				state := tc.existingState.ToResourceModel(ctx, &dgs).ToState(ctx, &dgs)
				if dgs.HasError() {
					t.Errorf("error generating existing state: %s", pp.PrettyPrint(dgs))
				}
				req := resource.ReadRequest{State: state}

				state = tc.expectedState.ToResourceModel(ctx, &dgs).ToState(ctx, &dgs)
				if dgs.HasError() {
					t.Errorf("error generating expected state: %s", pp.PrettyPrint(dgs))
				}
				expectedState := resource.ReadResponse{State: state}

				// get empty state
				emptyResource := NewRancherDevResource()
				schemaResponseContainer := &resource.SchemaResponse{}
				emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
				state = tfsdk.State{
					Schema: schemaResponseContainer.Schema,
				}
				res := resource.ReadResponse{State: state}

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
	t.Run("Update function", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherDevResource
			env                map[string]string // a k/v map of environment variables to set
			plan               RancherDevModel
			existingState      RancherDevModel
			expectedState      RancherDevModel
			apiResponse        c.Response // this will be injected into the client
			expectedApiRequest c.Request  // the API request expected to be reported from the client
			outcome            string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan
				RancherDevModel{
					Id: defaultId,
				},
				// existing state
				RancherDevModel{
					Id: defaultId,
				},
				// resulting state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(RancherDevModel{
						Id: defaultId,
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "PUT",
					Body: RancherDevModel{
						Id: defaultId,
					},
				},
				// expected outcome
				"success",
			},
			{
				"Update on Deleted Resource",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan set in update request
				RancherDevModel{
					Id: defaultId,
				},
				// existing state set in update request
				RancherDevModel{
					Id: defaultId,
				},
				// resulting state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: http.StatusNotFound,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(c.ErrorResponse{
						Status:  "404",
						Message: "resource not found",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "PUT",
					Body:     RancherDevModel{Id: defaultId},
				},
				// expected outcome
				"failure",
			},
			{
				"Server Error on Update",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan set in update request
				RancherDevModel{
					Id: defaultId,
				},
				// existing state set in update request
				RancherDevModel{
					Id: defaultId,
				},
				// resulting state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: http.StatusInternalServerError,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(c.ErrorResponse{
						Status:  "500",
						Message: "server error",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "PUT",
					Body:     RancherDevModel{Id: defaultId},
				},
				// expected outcome
				"failure",
			},
			{
				"Partial Attribute Update",
				RancherDevResource{},
				// env
				map[string]string{},
				// plan set in update request
				RancherDevModel{
					Id: defaultId,
				},
				// existing state set in update request
				RancherDevModel{
					Id: defaultId,
				},
				// resulting state expected to match this
				RancherDevModel{
					Id: defaultId,
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(RancherDevModel{
						Id: defaultId,
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "PUT",
					Body: RancherDevModel{
						Id: defaultId,
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
				var dgs diag.Diagnostics
				plan := tc.plan.ToResourceModel(ctx, &dgs).ToPlan(ctx, &dgs)
				if dgs.HasError() {
					t.Errorf("error generating plan: %s", pp.PrettyPrint(dgs))
				}
				state := tc.existingState.ToResourceModel(ctx, &dgs).ToState(ctx, &dgs)
				if dgs.HasError() {
					t.Errorf("error generating existing state: %s", pp.PrettyPrint(dgs))
				}
				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				// get empty state
				emptyResource := NewRancherDevResource()
				schemaResponseContainer := &resource.SchemaResponse{}
				emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
				state = tfsdk.State{
					Schema: schemaResponseContainer.Schema,
				}
				res := resource.UpdateResponse{
					State: state,
				}
				tc.fit.Update(ctx, req, &res)

				actualState := res
				state = tc.expectedState.ToResourceModel(ctx, &dgs).ToState(ctx, &dgs)
				if dgs.HasError() {
					t.Errorf("error generating expected state: %s", pp.PrettyPrint(dgs))
				}
				expectedState := resource.UpdateResponse{
					State: state,
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
	t.Run("Delete function", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherDevResource
			env                map[string]string // a k/v map of environment variables to set
			existingState      RancherDevModel
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
				RancherDevModel{
					Id: defaultId,
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: http.StatusNoContent,
					Headers:    map[string][]string{},
					Body:       []byte{},
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "DELETE",
				},
				// expected outcome
				"success",
			},
			{
				"Resource Already Deleted",
				RancherDevResource{},
				// env
				map[string]string{},
				// existing state
				RancherDevModel{
					Id: defaultId,
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: http.StatusNotFound,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(c.ErrorResponse{
						Status:  "404",
						Message: "resource not found",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "DELETE",
				},
				// expected outcome
				"success",
			},
			{
				"Server Error on Delete",
				RancherDevResource{},
				// env
				map[string]string{},
				// existing state
				RancherDevModel{
					Id: defaultId,
				},
				// this response will be injected into the client
				c.Response{
					StatusCode: http.StatusInternalServerError,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: resBodyMarshall(c.ErrorResponse{
						Status:  "500",
						Message: "server error",
					}),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "DELETE",
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

				var dgs diag.Diagnostics
				state := tc.existingState.ToResourceModel(ctx, &dgs).ToState(ctx, &dgs)
				if dgs.HasError() {
					t.Errorf("error generating existing state: %s", pp.PrettyPrint(dgs))
				}
				req := resource.DeleteRequest{
					State: state,
				}
				// get empty state
				emptyResource := NewRancherDevResource()
				schemaResponseContainer := &resource.SchemaResponse{}
				emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
				state = tfsdk.State{
					Schema: schemaResponseContainer.Schema,
				}
				res := resource.DeleteResponse{
					State: state,
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
