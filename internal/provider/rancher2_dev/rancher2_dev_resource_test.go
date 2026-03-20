package rancher2_dev

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

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
	apiUrl        = "https://rancher.example.com"
	endpoint      = "dev"
	apiEndpoint   = apiUrl + "/" + endpoint
	defaultId     = "dev-test"
	testTokenId   = "token-test"
	testTokenKey  = "this1is2a3test4token5it6is7fake"
	testToken     = testTokenId + ":" + testTokenKey
	testUserToken = "ext/" + testToken
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
				resource.MetadataResponse{TypeName: "rancher2_dev"},
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
			apiRequest    c.Request         // what the request should look like as reported by the client
			apiResponse   c.Response        // the API response to inject into the client
			expectedState RancherDevModel   // what should be in the state, translated to struct
			outcome       string            // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{ // plan
					Id:              defaultId,       // required
					StringAttribute: "dev-test",      // required
					NumberAttribute: big.NewFloat(1), // required
					// Int32Attribute:   1,      // this attribute is read only so it shouldn't be in the plan
					BoolAttribute:    false,
					Int64Attribute:   int64(1),
					Float64Attribute: 1.0,
					Float32Attribute: float32(1.0),
					ListAttribute:    []string{"test"},
					SetAttribute:     map[string]bool{"test": true},
					MapAttribute:     map[string]string{"test": "test"},
					NestedObject:     NestedObject{
						StringAttribute:    "test",
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
				c.Request{
					Endpoint: apiEndpoint,
					Method:   "POST",
					Headers: map[string][]string{
						"Content-Type":  {"application/json"},
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(RancherDevModel{
						Id:               defaultId,
						BoolAttribute:    false,
						NumberAttribute:  big.NewFloat(1),
						// Int32Attribute:   1, // This is a read only attribute, so it shouldn't be in the request
						Int64Attribute:   int64(1),
						Float64Attribute: 1.0,
						Float32Attribute: float32(1.0),
						StringAttribute:  "dev-test",
						ListAttribute:    []string{"test"},
						SetAttribute:     map[string]bool{"test": true},
						MapAttribute:     map[string]string{"test": "test"},
						NestedObject:     NestedObject{
							StringAttribute:    "test",
							NestedNestedObject: NestedNestedObject{
								StringAttribute: "test",
								BoolAttribute:   true,
							},
						},
						NestedObjectList: []NestedObject{
							{
								StringAttribute:    "test",
								NestedNestedObject: NestedNestedObject{
									StringAttribute: "test",
									BoolAttribute:   true,
								},
							},
						},
						NestedObjectMap: map[string]NestedObject{
							"test": {
								StringAttribute:    "test",
								NestedNestedObject: NestedNestedObject{
									StringAttribute: "test",
									BoolAttribute:   true,
								},
							},
						},
					}),
				},
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(RancherDevModel{
						Id:               defaultId,
						NumberAttribute:  big.NewFloat(1),
						StringAttribute:  "dev-test",
						Int32Attribute:   int32(1),  // respond with the read-only value
						BoolAttribute:    false,
						Int64Attribute:   int64(1),
						Float64Attribute: 1.0,
						Float32Attribute: float32(1.0),
						ListAttribute:    []string{"test"},
						SetAttribute:     map[string]bool{"test": true},
						MapAttribute:     map[string]string{"test": "test"},
						NestedObject:     NestedObject{
							StringAttribute:    "test",
							NestedNestedObject: NestedNestedObject{
								StringAttribute: "test",
								BoolAttribute:   true,
							},
						},
						NestedObjectList: []NestedObject{
							{
								StringAttribute:    "test",
								NestedNestedObject: NestedNestedObject{
									StringAttribute: "test",
									BoolAttribute:   true,
								},
							},
						},
						NestedObjectMap: map[string]NestedObject{
							"test": {
								StringAttribute:    "test",
								NestedNestedObject: NestedNestedObject{
									StringAttribute: "test",
									BoolAttribute:   true,
								},
							},
						},
					}),
				},
				RancherDevModel{ // expected state
					Id:               defaultId,        // required
					StringAttribute:  "dev-test",       // required
					NumberAttribute:  big.NewFloat(1),  // required
					Int32Attribute:   int32(1),  // the read only attribute should be available for query
					BoolAttribute:    false,
					Int64Attribute:   int64(1),
					Float64Attribute: 1.0,
					Float32Attribute: float32(1.0),
					ListAttribute:    []string{"test"},
					SetAttribute:     map[string]bool{"test": true},
					MapAttribute:     map[string]string{"test": "test"},
					NestedObject:     NestedObject{
						StringAttribute:    "test",
						NestedNestedObject: NestedNestedObject{
							StringAttribute: "test",
							BoolAttribute:   true,
						},
					},
					NestedObjectList: []NestedObject{
						{
							StringAttribute:    "test",
							NestedNestedObject: NestedNestedObject{
								StringAttribute: "test",
								BoolAttribute:   true,
							},
						},
					},
					NestedObjectMap: map[string]NestedObject{
						"test": {
							StringAttribute:    "test",
							NestedNestedObject: NestedNestedObject{
								StringAttribute: "test",
								BoolAttribute:   true,
							},
						},
					},
				},
				"success", // expected outcome
			},
			{
				"API Conflict",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{ // plan MUST include required attributes
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				c.Request{
					Endpoint: apiEndpoint,
					Method:   "POST",
					Headers: map[string][]string{
						"Content-Type":  {"application/json"},
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(RancherDevModel{ // request should include the required attributes as well
						Id:              defaultId,
						StringAttribute: "dev-test",
						NumberAttribute: big.NewFloat(1),
					}),
				},
				c.Response{
					StatusCode: http.StatusConflict,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "409",
						Message: "resource already exists",
					}),
				},
				RancherDevModel{ // expected state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				"failure", // expected outcome
			},
			{
				"Server Error",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{ // plan
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				c.Request{
					Endpoint: apiEndpoint,
					Method:   "POST",
					Headers: map[string][]string{
						"Content-Type":  {"application/json"},
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(RancherDevModel{
						Id: "test",
					}),
				},
				c.Response{
					StatusCode: http.StatusInternalServerError,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "500",
						Message: "something went wrong",
					}),
				},
				RancherDevModel{ // expected state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				"failure", // expected outcome
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
				var err error
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR") // this enables tflog.Debug, change to DEBUG when troubleshooting

        ctx := h.GenerateTestContext(t, &buf, nil)
				ts := &c.TokenStore{}
				ts.SetToken(testUserToken)
				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, ts)

        apiRequestId := fmt.Sprintf("%s:%s:%s", tc.apiRequest.Endpoint, tc.apiRequest.Method, testUserToken)
				client.SetResponse(ctx, apiRequestId, tc.apiResponse)

        err = h.GetConfiguredResource(ctx, t, &tc.fit, client)
				if err != nil {
					t.Fatalf("error configuring resource: %+v", err)
				}

        t.Logf("Fit after configure: %#v", tc.fit)

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

				actualApiRequest := client.GetRequest(apiRequestId)
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

				if diff := cmp.Diff(expectedApiRequest, actualApiRequest,
					cmpopts.IgnoreUnexported(big.Float{}),
					cmp.Comparer(func(x, y c.Request) bool {
						if x.Method != y.Method || x.Endpoint != y.Endpoint {
							return false
						}
						xBody, xOk := x.Body.([]byte)
						if !xOk {
							t.Logf("failed to assert x.Body to []byte")
							return false
						}
						yBody, yOk := y.Body.([]byte)
						if !yOk {
							t.Logf("failed to assert y.Body to []byte")
							return false
						}
						if len(xBody) == 0 && len(yBody) == 0 {
							return true
						}
						if len(xBody) == 0 || len(yBody) == 0 {
							return false
						}
						var want, got RancherDevModel
						if err := json.Unmarshal(xBody, &want); err != nil {
							t.Logf("failed to unmarshal want body: %v", err)
							return false
						}
						if err := json.Unmarshal(yBody, &got); err != nil {
							t.Logf("failed to unmarshal got body: %v", err)
							return false
						}
						return cmp.Equal(want, got, cmpopts.IgnoreUnexported(big.Float{}))
					}),
				); diff != "" {
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
			name          string
			fit           RancherDevResource
			env           map[string]string // a k/v map of environment variables to set
			existingState RancherDevModel   // this will get injected in the read request
			expectedState RancherDevModel
			apiRequest    c.Request  // the API request expected to be reported from the client
			apiResponse   c.Response // this will be injected into the client
			outcome       string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				// env
				map[string]string{},
				// existing state set in read request
				RancherDevModel{
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				// resulting state expected to match this
				RancherDevModel{
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "GET",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(nil),
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(RancherDevModel{
						Id:              defaultId,
						StringAttribute: "dev-test",
						NumberAttribute: big.NewFloat(1),
					}),
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
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				// resulting state expected to match this
				RancherDevModel{
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "GET",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(nil),
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(RancherDevModel{
						Id:              defaultId,
						StringAttribute: "dev-test",
						NumberAttribute: big.NewFloat(1),
					}),
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
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				// resulting state expected to match this
				RancherDevModel{
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "GET",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(nil),
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusInternalServerError,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "500",
						Message: "something went wrong",
					}),
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
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				// resulting state expected to match this
				RancherDevModel{
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				// the API request expected to be reported
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "GET",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(nil),
				},
				// the response to inject into the client
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(struct {
						Id                 string     `json:"id"`
						StringAttribute    string     `json:"string_attribute,omitempty"`
						NumberAttribute    *big.Float `json:"number_attribute,omitempty"`
						UntrackedAttribute string     `json:"untracked_attribute,omitempty"`
					}{
						Id:                 defaultId,
						StringAttribute:    "dev-test",
						NumberAttribute:    big.NewFloat(1),
						UntrackedAttribute: "untracked",
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

				ts := &c.TokenStore{}
				ts.SetToken(testUserToken)
				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, ts)
				apiRequestId := fmt.Sprintf("%s:%s:%s", tc.apiRequest.Endpoint, tc.apiRequest.Method, testUserToken)
				client.SetResponse(ctx, apiRequestId, tc.apiResponse)
				_, err := client.GetResponse(apiRequestId)
				if err != nil {
					t.Errorf("error getting response: %+v", err)
				}

				err = h.GetConfiguredResource(ctx, t, &tc.fit, client)
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

				tc.fit.Read(ctx, req, &res)
				actualState := res
				expectedApiRequest := tc.apiRequest
				actualApiRequest := client.GetRequest(apiRequestId)

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
				if diff := cmp.Diff(expectedApiRequest, actualApiRequest,
					cmpopts.IgnoreUnexported(big.Float{}),
					cmp.Comparer(func(x, y c.Request) bool {
						if x.Method != y.Method || x.Endpoint != y.Endpoint {
							return false
						}
						xBody, xOk := x.Body.([]byte)
						if !xOk {
							t.Logf("failed to assert x.Body to []byte")
							return false
						}
						yBody, yOk := y.Body.([]byte)
						if !yOk {
							t.Logf("failed to assert y.Body to []byte")
							return false
						}
						if len(xBody) == 0 && len(yBody) == 0 {
							return true
						}
						if len(xBody) == 0 || len(yBody) == 0 {
							return false
						}
						var want, got RancherDevModel
						if err := json.Unmarshal(xBody, &want); err != nil {
							t.Logf("failed to unmarshal want body: %v", err)
							return false
						}
						if err := json.Unmarshal(yBody, &got); err != nil {
							t.Logf("failed to unmarshal got body: %v", err)
							return false
						}
						return cmp.Equal(want, got, cmpopts.IgnoreUnexported(big.Float{}))
					}),
				); diff != "" {
					t.Errorf("Create() mismatch (-want +got): %s", diff)
				}
			})
		}
	})
	t.Run("Update", func(t *testing.T) {
		testCases := []struct {
			name          string
			fit           RancherDevResource
			env           map[string]string // a k/v map of environment variables to set
			plan          RancherDevModel
			existingState RancherDevModel
			apiRequest    c.Request  // the API request expected to be reported from the client
			apiResponse   c.Response // this will be injected into the client
			expectedState RancherDevModel
			outcome       string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{ // plan
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted, (defaulted attributes need to be included in the plan's model)
        },
				RancherDevModel{ // existing state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "PUT",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(RancherDevModel{
						Id:              defaultId,
						StringAttribute: "dev-test",
						NumberAttribute: big.NewFloat(1),
						BoolAttribute:   true,
					}),
				},
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(map[string]any{
						"Id":              defaultId,
						"StringAttribute": "dev-test",
						"NumberAttribute": big.NewFloat(1),
					}),
				},
				RancherDevModel{ // expected state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				"success",
			},
			{
				"Update on Deleted Resource",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{ // plan
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				RancherDevModel{ // existing state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "PUT",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(RancherDevModel{Id: defaultId}),
				},
				c.Response{
					StatusCode: http.StatusNotFound,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "404",
						Message: "resource not found",
					}),
				},
				RancherDevModel{ // expected state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
        "failure",
			},
			{
				"Server Error on Update",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{ // plan
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				RancherDevModel{ // existing state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted attributes end up looking like required ones on the other side of the plan phase
				},
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "PUT",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(RancherDevModel{Id: defaultId}),
				},
				c.Response{
					StatusCode: http.StatusInternalServerError,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "500",
						Message: "server error",
					}),
				},
				RancherDevModel{ // expected state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				"failure",
			},
			{
				"Partial Attribute Update",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{ // plan
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				RancherDevModel{ // existing state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "PUT",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(RancherDevModel{
						Id:              defaultId,
						StringAttribute: "dev-test",
						NumberAttribute: big.NewFloat(1),
						BoolAttribute:   true, // default
					}),
				},
				c.Response{
					StatusCode: http.StatusOK,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(RancherDevModel{
						Id:              defaultId,
						StringAttribute: "dev-test",
						NumberAttribute: big.NewFloat(1),
            BoolAttribute:   true, // defaulted
					}),
				},
				RancherDevModel{ // expected state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
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

				ts := &c.TokenStore{}
				ts.SetToken(testUserToken)
				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, ts)
				apiRequestId := fmt.Sprintf("%s:%s:%s", tc.apiRequest.Endpoint, tc.apiRequest.Method, testUserToken)
				client.SetResponse(ctx, apiRequestId, tc.apiResponse)

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

				actualApiRequest := client.GetRequest(apiRequestId)
				expectedApiRequest := tc.apiRequest

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
	t.Run("Delete", func(t *testing.T) {
		testCases := []struct {
			name          string
			fit           RancherDevResource
			env           map[string]string // a k/v map of environment variables to set
			existingState RancherDevModel
			apiRequest    c.Request  // the API request expected to be reported from the client
			apiResponse   c.Response // this will be injected into the client
			outcome       string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{ // existing state
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
          BoolAttribute:   true, // defaulted
				},
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId, // add the id to the path
					Method:   "DELETE",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(nil),
				},
				c.Response{
					StatusCode: http.StatusNoContent,
					Headers:    map[string][]string{},
					Body:       []byte{},
				},
				"success",
			},
			{
				"Resource Already Deleted",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "DELETE",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(nil),
				},
				c.Response{
					StatusCode: http.StatusNotFound,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "404",
						Message: "resource not found",
					}),
				},
				"success",
			},
			{
				"Server Error on Delete",
				RancherDevResource{},
				map[string]string{},
				RancherDevModel{
					Id:              defaultId,
					StringAttribute: "dev-test",
					NumberAttribute: big.NewFloat(1),
				},
				c.Request{
					Endpoint: apiEndpoint + "/" + defaultId,
					Method:   "DELETE",
					Headers: map[string][]string{
						"Authorization": {"Bearer " + testUserToken},
					},
					Body: rBodyMarshal(nil),
				},
				c.Response{
					StatusCode: http.StatusInternalServerError,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "500",
						Message: "server error",
					}),
				},
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

				ts := &c.TokenStore{}
				ts.SetToken(testUserToken)
				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, ts)
				apiRequestId := fmt.Sprintf("%s:%s:%s", tc.apiRequest.Endpoint, tc.apiRequest.Method, testUserToken)
				client.SetResponse(ctx, apiRequestId, tc.apiResponse)

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

				actualApiRequest := client.GetRequest(apiRequestId)
				expectedApiRequest := tc.apiRequest

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
				if diff := cmp.Diff(expectedApiRequest, actualApiRequest,
					cmpopts.IgnoreUnexported(big.Float{}),
					cmp.Comparer(func(x, y c.Request) bool {
						if x.Method != y.Method || x.Endpoint != y.Endpoint {
							return false
						}
						xBody, xOk := x.Body.([]byte)
						if !xOk {
							t.Logf("failed to assert x.Body to []byte")
							return false
						}
						yBody, yOk := y.Body.([]byte)
						if !yOk {
							t.Logf("failed to assert y.Body to []byte")
							return false
						}
						if len(xBody) == 0 && len(yBody) == 0 {
							return true
						}
						if len(xBody) == 0 || len(yBody) == 0 {
							return false
						}
						var want, got RancherDevModel
						if err := json.Unmarshal(xBody, &want); err != nil {
							t.Logf("failed to unmarshal want body: %v", err)
							return false
						}
						if err := json.Unmarshal(yBody, &got); err != nil {
							t.Logf("failed to unmarshal got body: %v", err)
							return false
						}
						return cmp.Equal(want, got, cmpopts.IgnoreUnexported(big.Float{}))
					}),
				); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
}

// helpers.
func rBodyMarshal(obj any) []byte {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}
