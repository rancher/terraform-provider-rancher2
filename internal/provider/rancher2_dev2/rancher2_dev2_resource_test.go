package rancher2_dev2

import (
	"encoding/json"
	"fmt"

	// "reflect"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	// "github.com/hashicorp/terraform-plugin-framework/tfsdk"

	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	mta "github.com/rancher/terraform-provider-rancher2/internal/provider/rancher2_metadata"
	h "github.com/rancher/terraform-provider-rancher2/internal/provider/test_helpers"
)

var testToken = "test-token"
var baseEndpoint = "https://rancher.example.com/dev2"
var testID = "test-id"

func TestRancher2Dev2Resource(t *testing.T) {
	t.Run("Metadata", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  Rancher2Dev2Resource
			want resource.MetadataResponse
		}{
			{
				"Basic",
				Rancher2Dev2Resource{},
				resource.MetadataResponse{TypeName: "rancher2_dev2"},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
				defer log()
				res := resource.MetadataResponse{}
				tc.fit.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "rancher2"}, &res)
				got := res
				if got.TypeName != tc.want.TypeName {
					t.Errorf("%+v.Metadata() TypeName is %s; want %s", tc.fit, got.TypeName, tc.want.TypeName)
				}
			})
		}
	})
	t.Run("Schema", func(t *testing.T) {
		testCases := []struct {
			name    string
			fit     Rancher2Dev2Resource
			want    []string
			outcome string
		}{
			{
				"Basic",
				Rancher2Dev2Resource{},
				[]string{
					"id",
					"api_version",
					"kind",
					"metadata",
					"spec",
					"status",
				},
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
				defer log()
				req := resource.SchemaRequest{}
				res := resource.SchemaResponse{}
				tc.fit.Schema(ctx, req, &res)
				if res.Schema.Attributes == nil {
					t.Errorf("%#v.Schema() returned nil attributes", tc.fit)
				}
				if tc.outcome == "failure" && !res.Diagnostics.HasError() {
					t.Errorf("Create() expected error diagnostics, got none")
				}
				gotKeys := []string{}
				for key := range res.Schema.Attributes {
					gotKeys = append(gotKeys, key)
				}
				for _, wantKey := range tc.want {
					if !slices.Contains(gotKeys, wantKey) {
						t.Errorf("%#v.Schema() missing expected key %s", tc.fit, wantKey)
					}
				}
			})
		}
	})
	t.Run("Configure", func(t *testing.T) {
		testCases := []struct {
			name    string
			env     map[string]string
			fit     Rancher2Dev2Resource
			outcome string
		}{
			{
				"Basic",                // configure
				map[string]string{},    // env
				Rancher2Dev2Resource{}, // fit
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
				defer log()
				req := resource.ConfigureRequest{
					ProviderData: h.GetTestClient(t, ctx),
				}
				res := resource.ConfigureResponse{}
				tc.fit.Configure(ctx, req, &res)
				if tc.outcome == "failure" && !res.Diagnostics.HasError() {
					t.Errorf("Create() expected error diagnostics, got none")
				}
				if tc.outcome != "failure" && res.Diagnostics.HasError() {
					t.Errorf("Create() has unexpected error diagnostics:\n%+v\n", pp.PrettyPrint(res.Diagnostics))
				}
			})
		}
	})
	t.Run("Create", func(t *testing.T) {
		testCases := []struct {
			name          string
			env           map[string]string
			fit           Rancher2Dev2Resource
			plan          Rancher2Dev2Model
			apiRequest    c.Request
			apiResponse   c.Response
			expectedState Rancher2Dev2Model
			outcome       string
		}{
			{
				"Basic",                // create
				map[string]string{},    // env
				Rancher2Dev2Resource{}, // fit
				Rancher2Dev2Model{ // Go model to convert to tfsdk.Plan to send to create function
					ID:         testID, // this won't normally be in the plan, but adding it for test consistency
					APIVersion: "string",
					Kind:       "string",
					Metadata:   mta.SampleMetadataGoModel(),
					Spec: Spec{
						String:  "test",
						Bool:    false,
						Number:  1,
						Int32:   1,
						Int64:   1,
						Float32: 1.0,
						Float64: 1.0,
						Map:     map[string]string{"test": "test"},
						List:    []string{"test"},
						Object: Object{
							StringAttribute: "test",
						},
						ObjectList: []Object{
							{
								StringAttribute: "test",
							},
						},
						ObjectMap: map[string]Object{
							"test": {
								StringAttribute: "test",
							},
						},
					},
				},
				c.Request{ // expected API request
					Endpoint: baseEndpoint,
					Method:   "POST",
					Headers: map[string][]string{
						"Content-Type":  {"application/json"},
						"Authorization": {"Bearer " + testToken},
					},
					Body: apiBody(&Rancher2Dev2Model{
						APIVersion: testID,
						Kind:       "string",
						Metadata:   mta.SampleMetadataGoModel(),
						Spec: Spec{
							String:  "test",
							Bool:    false,
							Number:  1,
							Int32:   1,
							Int64:   1,
							Float32: 1.0,
							Float64: 1.0,
							Map:     map[string]string{"test": "test"},
							List:    []string{"test"},
							Object: Object{
								StringAttribute: "test",
							},
							ObjectList: []Object{
								{
									StringAttribute: "test",
								},
							},
							ObjectMap: map[string]Object{
								"test": {
									StringAttribute: "test",
								},
							},
						},
					}),
				},
				c.Response{ // Injected into the client and "served" to the code.
					StatusCode: 201,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: apiResponse(&Rancher2Dev2Model{
						APIVersion: testID,
						Kind:       "string",
						Metadata:   mta.SampleMetadataGoModel(),
						Spec: Spec{
							String:  "test",
							Bool:    false,
							Number:  1,
							Int32:   1,
							Int64:   1,
							Float32: 1.0,
							Float64: 1.0,
							Map:     map[string]string{"test": "test"},
							List:    []string{"test"},
							Object: Object{
								StringAttribute: "test",
							},
							ObjectList: []Object{
								{
									StringAttribute: "test",
								},
							},
							ObjectMap: map[string]Object{
								"test": {
									StringAttribute: "test",
								},
							},
						},
						Status: map[string]any{
							"status": "active",
						},
					}),
				},
				Rancher2Dev2Model{ // Go model to convert to tfsdk.State for comparing against the resulting state
					ID:         testID,
					APIVersion: "string",
					Kind:       "string",
					Metadata:   mta.SampleMetadataGoModel(),
					Spec: Spec{
						String:  "test",
						Bool:    false,
						Number:  1,
						Int32:   1,
						Int64:   1,
						Float32: 1.0,
						Float64: 1.0,
						Map:     map[string]string{"test": "test"},
						List:    []string{"test"},
						Object: Object{
							StringAttribute: "test",
						},
						ObjectList: []Object{
							{
								StringAttribute: "test",
							},
						},
						ObjectMap: map[string]Object{
							"test": {
								StringAttribute: "test",
							},
						},
					},
				},
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
				defer log()
				clnt := h.GetTestClient(t, ctx)
				clnt.SetToken(testToken)

				if tc.apiRequest.Endpoint != "" {
					apiRequestId := fmt.Sprintf("%s:%s:%s", tc.apiRequest.Endpoint, tc.apiRequest.Method, testToken)
					clnt.SetResponse(ctx, apiRequestId, tc.apiResponse)
				}

				err := h.GetConfiguredResource(ctx, t, &tc.fit, clnt)
				if err != nil {
					t.Fatalf("Error getting configured resource: %s", err.Error())
				}
				var dgs diag.Diagnostics

				plan := tc.plan.ToResourceModel(ctx, &dgs).ToPlan(ctx, &dgs)
				if dgs.HasError() {
					t.Fatalf("Error creating plan: %s", pp.PrettyPrint(dgs))
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				// get empty state
				emptyResource := NewRancher2Dev2Resource()
				schemaResponseContainer := &resource.SchemaResponse{}
				emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
				state := tfsdk.State{
					Schema: schemaResponseContainer.Schema,
				}
				res := resource.CreateResponse{State: state}
        // res = resource.CreateResponse{} // This causes a panic!

        tc.fit.Create(ctx, req, &res)
				if tc.outcome == "failure" && !res.Diagnostics.HasError() {
					t.Errorf("Create() expected error diagnostics, got none")
				}
				if tc.outcome != "failure" && res.Diagnostics.HasError() {
					t.Errorf("Create() has unexpected error diagnostics:\n%+v\n", pp.PrettyPrint(res.Diagnostics))
				}
			})
		}
	})
	t.Run("Read", func(t *testing.T) {
		testCases := []struct {
			name    string
			fit     Rancher2Dev2Resource
			outcome string
		}{
			{
				"Basic",
				Rancher2Dev2Resource{},
				"failure",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
				defer log()
				req := resource.ReadRequest{}
				res := resource.ReadResponse{}
				tc.fit.Read(ctx, req, &res)
				if tc.outcome == "failure" && !res.Diagnostics.HasError() {
					t.Errorf("Read() expected error diagnostics, got none")
				}
			})
		}
	})
	t.Run("Update", func(t *testing.T) {
		testCases := []struct {
			name    string
			fit     Rancher2Dev2Resource
			outcome string
		}{
			{
				"Basic",
				Rancher2Dev2Resource{},
				"failure",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
				defer log()
				req := resource.UpdateRequest{}
				res := resource.UpdateResponse{}
				tc.fit.Update(ctx, req, &res)
				if tc.outcome == "failure" && !res.Diagnostics.HasError() {
					t.Errorf("Update() expected error diagnostics, got none")
				}
			})
		}
	})
	t.Run("Delete", func(t *testing.T) {
		testCases := []struct {
			name    string
			fit     Rancher2Dev2Resource
			outcome string
		}{
			{
				"Basic",
				Rancher2Dev2Resource{},
				"failure",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
				defer log()
				req := resource.DeleteRequest{}
				res := resource.DeleteResponse{}
				tc.fit.Delete(ctx, req, &res)
				if tc.outcome == "failure" && !res.Diagnostics.HasError() {
					t.Errorf("Delete() expected error diagnostics, got none")
				}
			})
		}
	})
	t.Run("ImportState", func(t *testing.T) {
		testCases := []struct {
			name    string
			fit     Rancher2Dev2Resource
			outcome string
		}{
			{
				"Basic",
				Rancher2Dev2Resource{},
				"failure",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx, log := h.Cntxt(t, "ERROR") // Change the log level to "DEBUG" here for more logging.
				defer log()
				req := resource.ImportStateRequest{}
				res := resource.ImportStateResponse{}
				tc.fit.ImportState(ctx, req, &res)
				if tc.outcome == "failure" && !res.Diagnostics.HasError() {
					t.Errorf("ImportState() expected error diagnostics, got none")
				}
			})
		}
	})
}

// Helpers.
// apiBody turns a native Go model to a format that can easily be marshaled into json
func apiBody(data *Rancher2Dev2Model) any {
	dgs := diag.Diagnostics{}
	body := data.ToApiRequestBody(&dgs)
	if dgs.HasError() {
		return nil
	}
	return body
}

// apiResponse turns a native Go model to a marshaled byte representation
func apiResponse(data *Rancher2Dev2Model) []byte {
	body, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return body
}
