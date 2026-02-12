package rancher2_login

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	h "github.com/rancher/terraform-provider-rancher2/internal/provider/test_helpers"
)

const (
	apiUrl        = "https://rancher.example.com"
	loginEndpoint = "v1-public/login"
	tokenEndpoint = "/apis/ext.cattle.io/v1/tokens"
	defaultId     = "test"
)

func TestRancherLoginResource(t *testing.T) {
	t.Run("Metadata", func(t *testing.T) {
		testCases := []struct {
			name string
			fit  RancherLoginResource
			want resource.MetadataResponse
		}{
			{
				"Basic",
				RancherLoginResource{},
				resource.MetadataResponse{TypeName: "rancher2_login_resource"},
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
			fit  RancherLoginResource
			want []string
		}{
			{
				"Basic",
				RancherLoginResource{},
				[]string{
					"id",
					"username",
					"password",
					"username_environment_variable",
					"password_environment_variable",
					"token_ttl",
					"refresh_at",
					"ignore_token",
					"user_token",
					"session_token",
					"user_token_start_date",
					"user_token_end_date",
					"user_token_refresh_date",
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
			fit  RancherLoginResource
			want RancherLoginResource
		}{
			{
				"Basic",
				RancherLoginResource{},
				RancherLoginResource{},
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
			fit           RancherLoginResource
			env           map[string]string
			plan          RancherLoginModel
			expectedState RancherLoginModel
			loginRequest  c.Request
			loginResponse c.Response
			apiRequest    c.Request
			apiResponse   c.Response
			outcome       string
		}{
			{
				"Basic",
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{
					Username: "user",
					Password: "password",
				},
				RancherLoginModel{
					Id:                          defaultId,
					Username:                    "user",
					Password:                    "password",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					UserToken:                   "a-token",
					SessionToken:                "a-token",
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserTokenStartDate:          "",
					UserTokenEndDate:            "",
					UserTokenRefreshDate:        "",
				},
				// login
				c.Request{
					Endpoint: apiUrl + "/" + loginEndpoint,
					Method:   "POST",
					Body: rBodyMarshal(map[string]any{
						"type":         "localProvider",
						"username":     "user",
						"password":     "password",
						"responseType": "json",
					}),
				},
				// login
				c.Response{
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"baseType":  "token",
						"expiresAt": "2026-02-13T13:08:30Z",
						"id":        "token-r6kwh",
						"token":     "token-r6kwh:bzljw8trcgwq6vltxv74445c2zlspwd5nl828kbbnhwp8jr7bwpq92",
						"type":      "token",
					}),
				},
				// token
				c.Request{
					Endpoint: apiUrl + "/" + tokenEndpoint,
					Method:   "POST",
					Body: rBodyMarshal(map[string]any{
						"apiVersion": "ext.cattle.io/v1",
						"kind":       "Token",
						"spec": map[string]any{
							"description": "terraform-login-49fb2932-8f0c-4f39-ad6b-6e91e40bad78",
							"ttl":         7776000000,
						},
					}),
				},
				// token
				c.Response{
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"kind":       "Token",
						"apiVersion": "ext.cattle.io/v1",
						"metadata": map[string]any{
							"name":              "token-t9d75",
							"generateName":      "token-",
							"uid":               "b9f42296-de97-4150-bd2d-dcc211b9621c",
							"resourceVersion":   "288311",
							"creationTimestamp": "2026-02-12T21:49:07Z",
							"labels": map[string]any{
								"authn.management.cattle.io/kind": "",
								"cattle.io/user-id":               "user-5k5wm",
							},
						},
						"spec": map[string]any{
							"userID": "user-5k5wm",
							"userPrincipal": map[string]any{
								"name":          "local://user-5k5wm",
								"displayName":   "Default Admin",
								"loginName":     "admin",
								"principalType": "user",
								"me":            true,
								"provider":      "local",
							},
							"kind":        "",
							"description": "terraform-login-49fb2932-8f0c-4f39-ad6b-6e91e40bad78",
							"ttl":         7776000000,
							"enabled":     true,
						},
						"status": map[string]any{
							"value":          "ndc7njll8wt4tf2mvzdg4rckgwrkshrckbmj527klt8zs966twpz6h",
							"current":        false,
							"expired":        false,
							"expiresAt":      "2026-05-13T16:49:07-05:00",
							"lastUpdateTime": "2026-02-12T15:49:07-06:00",
							"bearerToken":    "ext/token-t9d75:ndc7njll8wt4tf2mvzdg4rckgwrkshrckbmj527klt8zs966twpz6h",
						},
					}),
				},
				"success",
			},
			{
				"API Conflict",
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{
					Username: "user",
					Password: "password",
				},
				RancherLoginModel{},
				// login
				c.Request{
					Endpoint: apiUrl + "/" + loginEndpoint,
					Method:   "POST",
					Body: rBodyMarshal(map[string]any{
						"type":         "localProvider",
						"username":     "user",
						"password":     "password",
						"responseType": "json",
					}),
				},
				// login
				c.Response{
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"baseType":  "token",
						"expiresAt": "2026-02-13T13:08:30Z",
						"id":        "token-r6kwh",
						"token":     "token-r6kwh:bzljw8trcgwq6vltxv74445c2zlspwd5nl828kbbnhwp8jr7bwpq92",
						"type":      "token",
					}),
				},
				c.Request{
					Endpoint: apiUrl + "/" + tokenEndpoint,
					Method:   "POST",
					Body: rBodyMarshal(RancherLoginModel{
						Id:                          defaultId,
						Username:                    "user",
						Password:                    "password",
						UsernameEnvironmentVariable: "RANCHER_USERNAME",
						PasswordEnvironmentVariable: "RANCHER_PASSWORD",
						TokenTtl:                    "90d",
						RefreshAt:                   "10d",
					}),
				},
				c.Response{
					StatusCode: http.StatusConflict,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "409",
						Message: "resource already exists",
					}),
				},
				"failure",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				defer func() {
					for k := range tc.env {
						os.Unsetenv(k)
					}
				}()
				for k, v := range tc.env {
					os.Setenv(k, v) //nolint:usetesting
				}
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR")
				ctx := h.GenerateTestContext(t, &buf, nil)

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, "")
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

				emptyResource := NewRancherLoginResource()
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
	t.Run("Read", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherLoginResource
			env                map[string]string // a k/v map of environment variables to set
			existingState      RancherLoginModel // this will get injected in the read request
			expectedState      RancherLoginModel
			apiResponse        c.Response // this will be injected into the client
			expectedApiRequest c.Request  // the API request expected to be reported from the client
			outcome            string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{Id: defaultId},
				RancherLoginModel{Id: defaultId},
				c.Response{
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body:       rBodyMarshal(RancherLoginModel{Id: defaultId}),
				},
				c.Request{
					Endpoint: apiUrl + "/" + tokenEndpoint + "/" + defaultId,
					Method:   "GET",
					Body:     rBodyMarshal(nil),
				},
				"success",
			},
			{
				"Resource Not Found",
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{Id: defaultId},
				RancherLoginModel{},
				c.Response{
					StatusCode: http.StatusNotFound,
				},
				c.Request{
					Endpoint: apiUrl + "/" + tokenEndpoint + "/" + defaultId,
					Method:   "GET",
					Body:     rBodyMarshal(nil),
				},
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR")
				ctx := h.GenerateTestContext(t, &buf, nil)

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, "")
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

				emptyResource := NewRancherLoginResource()
				schemaResponseContainer := &resource.SchemaResponse{}
				emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
				resState := tfsdk.State{
					Schema: schemaResponseContainer.Schema,
				}
				res := resource.ReadResponse{State: resState}

				tc.fit.Read(context.Background(), req, &res)
				if tc.name == "Resource Not Found" {
					if !res.State.Raw.IsNull() {
						t.Error("expected state to be removed")
					}
					return
				}

				actualState := res
				state = tc.expectedState.ToResourceModel(ctx, &dgs).ToState(ctx, &dgs)
				if dgs.HasError() {
					t.Errorf("error generating expected state: %s", pp.PrettyPrint(dgs))
				}
				expectedState := resource.ReadResponse{State: state}
				expectedApiRequest := tc.expectedApiRequest
				actualApiRequest := client.GetLastRequest()

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
					t.Errorf("Read() mismatch (-want +got): %s", diff)
				}
				if diff := cmp.Diff(expectedApiRequest, actualApiRequest); diff != "" {
					t.Errorf("Read() mismatch (-want +got): %s", diff)
				}
			})
		}
	})
	t.Run("Update", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherLoginResource
			env                map[string]string
			plan               RancherLoginModel
			existingState      RancherLoginModel
			expectedState      RancherLoginModel
			apiResponse        c.Response
			expectedApiRequest c.Request
			outcome            string
		}{
			{
				"Basic",
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{Id: defaultId, RefreshAt: "5d"},
				RancherLoginModel{Id: defaultId, RefreshAt: "10d"},
				RancherLoginModel{Id: defaultId, RefreshAt: "5d"},
				c.Response{
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body:       rBodyMarshal(RancherLoginModel{Id: defaultId, RefreshAt: "5d"}),
				},
				c.Request{
					Endpoint: apiUrl + "/" + tokenEndpoint + "/" + defaultId,
					Method:   "PUT",
					Body:     rBodyMarshal(RancherLoginModel{Id: defaultId, RefreshAt: "5d"}),
				},
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR")
				ctx := h.GenerateTestContext(t, &buf, nil)

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, "")
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
				emptyResource := NewRancherLoginResource()
				schemaResponseContainer := &resource.SchemaResponse{}
				emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
				resState := tfsdk.State{
					Schema: schemaResponseContainer.Schema,
				}
				res := resource.UpdateResponse{
					State: resState,
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
					t.Errorf("Update() mismatch (-want +got):\n%+v", diff)
				}
				if diff := cmp.Diff(expectedApiRequest, actualApiRequest); diff != "" {
					t.Errorf("Update() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
	t.Run("Delete", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherLoginResource
			existingState      RancherLoginModel
			apiResponse        c.Response
			expectedApiRequest c.Request
			outcome            string
		}{
			{
				"Basic",
				RancherLoginResource{},
				RancherLoginModel{Id: defaultId},
				c.Response{
					StatusCode: http.StatusNoContent,
				},
				c.Request{
					Endpoint: apiUrl + "/" + tokenEndpoint + "/" + defaultId,
					Method:   "DELETE",
					Body:     rBodyMarshal(nil),
				},
				"success",
			},
			{
				"Resource Already Deleted",
				RancherLoginResource{},
				RancherLoginModel{Id: defaultId},
				c.Response{
					StatusCode: http.StatusNotFound,
				},
				c.Request{
					Endpoint: apiUrl + "/" + tokenEndpoint + "/" + defaultId,
					Method:   "DELETE",
					Body:     rBodyMarshal(nil),
				},
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR")
				ctx := h.GenerateTestContext(t, &buf, nil)

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, "")
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
				emptyResource := NewRancherLoginResource()
				schemaResponseContainer := &resource.SchemaResponse{}
				emptyResource.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
				resState := tfsdk.State{
					Schema: schemaResponseContainer.Schema,
				}
				res := resource.DeleteResponse{
					State: resState,
				}
				tc.fit.Delete(ctx, req, &res)

				actualApiRequest := client.GetLastRequest()
				expectedApiRequest := tc.expectedApiRequest

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
					t.Errorf("Delete() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
}

func rBodyMarshal(obj any) []byte {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}
