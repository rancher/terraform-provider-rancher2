package rancher2_login

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	apiUrl           = "https://rancher.example.com"
	testSessionToken = "ext/test-token:this1is2a3session4token5it6is7fake"
	testTokenId      = "token-test"
	testTokenKey     = "this1is2a3test4token5it6is7fake"
	testToken        = testTokenId + ":" + testTokenKey
	testUserToken    = "ext/" + testToken
	testInitialToken = "ext/test-token:this1is2the3initial4test5token6it7is8fake"
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
				resource.MetadataResponse{TypeName: "rancher2_login"},
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
			loginRequest  c.Request
			loginResponse c.Response
			apiRequest    c.Request
			apiResponse   c.Response
			expectedState RancherLoginModel
			outcome       string
		}{
			{
				"Basic",
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{ // plan
					Username: "user",
					Password: "password",
				},
				c.Request{ // login
					Endpoint: fmt.Sprintf("%s/%s", apiUrl, loginEndpoint),
					Method:   "POST",
					Headers:  map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"type":         "localProvider",
						"username":     "user",
						"password":     "password",
						"responseType": "json",
					}),
				},
				c.Response{ // login
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"baseType":  "token",
						"expiresAt": "2026-02-13T13:08:30Z",
						"id":        testTokenId,
						"token":     testSessionToken,
						"type":      "token",
					}),
				},
				c.Request{ // token
					Endpoint: fmt.Sprintf("%s/%s", apiUrl, tokenEndpoint),
					Method:   "POST",
					Headers: map[string][]string{
						"Content-Type":  {"application/json"},
						"Authorization": {"Bearer " + testSessionToken},
					},
					Body: rBodyMarshal(map[string]any{
						"apiVersion": apiVersion,
						"kind":       "Token",
						"metadata": map[string]any{
							"name": testTokenId,
						},
						"spec": map[string]any{
							"description": "Terraform login token.",
							"ttl":         7776000000,
						},
					}),
				},
				c.Response{ // token
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"kind":       "Token",
						"apiVersion": apiVersion,
						"metadata": map[string]any{
							"name":              testTokenId,
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
							"description": "Terraform login token.",
							"ttl":         7776000000,
							"enabled":     true,
						},
						"status": map[string]any{
							"value":          testTokenKey,
							"current":        false,
							"expired":        false,
							"expiresAt":      "2026-05-13T16:49:07-05:00",
							"lastUpdateTime": "2026-02-12T15:49:07-06:00",
							"bearerToken":    testUserToken,
						},
					}),
				},
				RancherLoginModel{ // expected state
					Id:                          testTokenId,
					Username:                    "user",
					Password:                    "password",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					UserToken:                   testUserToken,
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserTokenStartDate:          "2026-02-12T21:49:07Z",
					UserTokenEndDate:            "2026-05-13T16:49:07-05:00",
					UserTokenRefreshDate:        "2026-05-03T16:49:07-05:00",
				},
				"success",
			},
			{
				"API Conflict",
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{ // plan
					Username: "user",
					Password: "password",
				},
				c.Request{ // login
					Endpoint: apiUrl + "/" + loginEndpoint,
					Method:   "POST",
					Headers:  map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"type":         "localProvider",
						"username":     "user",
						"password":     "password",
						"responseType": "json",
					}),
				},
				c.Response{ // login
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"baseType":  "token",
						"expiresAt": "2026-02-13T13:08:30Z",
						"id":        testTokenId,
						"token":     testSessionToken,
						"type":      "token",
					}),
				},
				c.Request{ // token
					Endpoint: apiUrl + "/" + tokenEndpoint,
					Method:   "POST",
					Headers: map[string][]string{
						"Content-Type":  {"application/json"},
						"Authorization": {"Bearer " + testSessionToken},
					},
					Body: rBodyMarshal(RancherLoginModel{
						Id:                          testTokenId,
						Username:                    "user",
						Password:                    "password",
						UsernameEnvironmentVariable: "RANCHER_USERNAME",
						PasswordEnvironmentVariable: "RANCHER_PASSWORD",
						TokenTtl:                    "90d",
						RefreshAt:                   "10d",
					}),
				},
				c.Response{ // token
					StatusCode: http.StatusConflict,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "409",
						Message: "resource already exists",
					}),
				},
				RancherLoginModel{}, // expected state
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

				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, &c.TokenStore{})

				loginRequestId := fmt.Sprintf("%s:%s:%s", tc.loginRequest.Endpoint, tc.loginRequest.Method, "")
				client.SetResponse(ctx, loginRequestId, tc.loginResponse)

				tokenRequestId := fmt.Sprintf("%s:%s:%s", tc.apiRequest.Endpoint, tc.apiRequest.Method, testSessionToken)
				client.SetResponse(ctx, tokenRequestId, tc.apiResponse)

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
				actualApiRequest := client.GetRequest(tokenRequestId)
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
			expectedApiRequest c.Request  // the API request expected to be reported from the client
			apiResponse        c.Response // this will be injected into the client
			outcome            string     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic", // read
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{
					Id:                          testTokenId,
					Username:                    "user",
					Password:                    "password",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					UserToken:                   testUserToken,
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserTokenStartDate:          "2026-02-12T21:49:07Z",
					UserTokenEndDate:            "2026-05-13T16:49:07-05:00",
					UserTokenRefreshDate:        "2026-05-03T16:49:07-05:00",
				},
				RancherLoginModel{
					Id:                          testTokenId,
					Username:                    "user",
					Password:                    "password",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					UserToken:                   testUserToken,
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserTokenStartDate:          "2026-02-12T21:49:07Z",
					UserTokenEndDate:            "2026-05-13T16:49:07-05:00",
					UserTokenRefreshDate:        "2026-05-03T16:49:07-05:00",
				},
				c.Request{
					Endpoint: fmt.Sprintf("%s/%s/%s", apiUrl, tokenEndpoint, testTokenId),
					Method:   "GET",
					Headers:  map[string][]string{"Authorization": {"Bearer " + testUserToken}}, // read will use the token in state
				},
				c.Response{
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"kind":       "Token",
						"apiVersion": apiVersion,
						"metadata": map[string]any{
							"name":              testTokenId,
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
							"description": "Terraform login token.",
							"ttl":         7776000000,
							"enabled":     true,
						},
						"status": map[string]any{
							"value":          testTokenKey,
							"current":        false,
							"expired":        false,
							"expiresAt":      "2026-05-13T16:49:07-05:00",
							"lastUpdateTime": "2026-02-12T15:49:07-06:00",
							"bearerToken":    testUserToken,
						},
					}),
				},
				"success",
			},
			{
				"Resource Not Found",
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{Id: testTokenId},
				RancherLoginModel{},
				c.Request{
					Endpoint: fmt.Sprintf("%s/%s/%s", apiUrl, tokenEndpoint, testTokenId),
					Method:   "GET",
					Headers:  map[string][]string{"Authorization": {"Bearer " + testUserToken}}, // read will use the token in state
				},
				c.Response{
					StatusCode: http.StatusNotFound,
				},
				"success",
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
				defer h.PrintLog(t, &buf, "ERROR") // enables printing tflog messages, set to DEBUG when coding
				ctx := h.GenerateTestContext(t, &buf, nil)

				ts := &c.TokenStore{}
				ts.SetToken(testUserToken)
				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, ts)

				apiRequestId := fmt.Sprintf("%s:%s:%s", tc.expectedApiRequest.Endpoint, tc.expectedApiRequest.Method, testUserToken)
				client.SetResponse(ctx, apiRequestId, tc.apiResponse)

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

				tc.fit.Read(ctx, req, &res)
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
				actualApiRequest := client.GetRequest(apiRequestId)

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
			name                 string
			fit                  RancherLoginResource
			env                  map[string]string
			plan                 RancherLoginModel
			existingState        RancherLoginModel
			initalTokenRequest   c.Request
			initialTokenResponse c.Response
			loginRequest         c.Request
			loginResponse        c.Response
			apiRequest           c.Request
			apiResponse          c.Response
			expectedState        RancherLoginModel
			outcome              string
		}{
			{
				"Basic", // update
				RancherLoginResource{},
				map[string]string{},
				RancherLoginModel{ //plan
					Username: "user",
					Password: "password",
				},
				RancherLoginModel{ //current state
					Id:                          testTokenId,
					Username:                    "user",
					Password:                    "password",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					UserToken:                   testUserToken,
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserTokenStartDate:          "2026-02-12T21:49:07Z",
					UserTokenEndDate:            "2026-05-13T16:49:07-05:00",
					UserTokenRefreshDate:        "2026-05-03T16:49:07-05:00",
				},
				c.Request{ //initial token request
					Endpoint: fmt.Sprintf("%s/%s", apiUrl, tokenEndpoint),
					Method:   "GET",
				},
				c.Response{ // initial token response
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"kind":       "Token",
						"apiVersion": apiVersion,
						"metadata": map[string]any{
							"name":              testTokenId,
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
							"description": "Terraform login token.",
							"ttl":         7776000000,
							"enabled":     true,
						},
						"status": map[string]any{
							"value":          testTokenKey,
							"current":        false,
							"expired":        false,
							"expiresAt":      "2026-05-13T16:49:07-05:00",
							"lastUpdateTime": "2026-02-12T15:49:07-06:00",
							"bearerToken":    testUserToken,
						},
					}),
				},
				c.Request{},  // no login request should be made
				c.Response{}, // no login response
				c.Request{},  // no token create should be made
				c.Response{}, // no token create response
				RancherLoginModel{ // expected state should match initial state
					Id:                          testTokenId,
					Username:                    "user",
					Password:                    "password",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					UserToken:                   testUserToken,
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserTokenStartDate:          "2026-02-12T21:49:07Z",
					UserTokenEndDate:            "2026-05-13T16:49:07-05:00",
					UserTokenRefreshDate:        "2026-05-03T16:49:07-05:00",
				},
				"success",
			},
			{
				"Update Token", // update
				RancherLoginResource{},
				map[string]string{
					"RANCHER_USERNAME": "user",
					"RANCHER_PASSWORD": "password",
				},
				RancherLoginModel{}, //plan
				RancherLoginModel{ //existing state
					Id:                          testTokenId,
					Username:                    "user",
					Password:                    "password",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					UserToken:                   testInitialToken,
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserTokenStartDate:          "2026-02-12T21:49:07Z",
					UserTokenEndDate:            "2026-05-13T16:49:07-05:00",
					UserTokenRefreshDate:        "2026-05-03T16:49:07-05:00",
				},
				c.Request{ //initial token request
					Endpoint: fmt.Sprintf("%s/%s", apiUrl, tokenEndpoint),
					Method:   "POST",
					Headers: map[string][]string{
						"Content-Type":  {"application/json"},
						"Authorization": {"Bearer " + testInitialToken},
					},
					Body: rBodyMarshal(map[string]any{
						"apiVersion": apiVersion,
						"kind":       "Token",
						"metadata": map[string]any{
							"name": testTokenId,
						},
						"spec": map[string]any{
							"description": "Terraform login token.",
							"ttl":         7776000000,
						},
					}),
				},
				c.Response{ // initial token response
					StatusCode: http.StatusForbidden,
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
					Body: rBodyMarshal(c.ErrorResponse{
						Status:  "403",
						Message: "Forbidden",
					}),
				},
				c.Request{ // login
					Endpoint: fmt.Sprintf("%s/%s", apiUrl, loginEndpoint),
					Method:   "POST",
					Headers:  map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"type":         "localProvider",
						"username":     "user",
						"password":     "password",
						"responseType": "json",
					}),
				},
				c.Response{ // login
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"baseType":  "token",
						"expiresAt": "2026-02-13T13:08:30Z",
						"id":        testTokenId,
						"token":     testSessionToken,
						"type":      "token",
					}),
				},
				c.Request{ // create token
					Endpoint: fmt.Sprintf("%s/%s", apiUrl, tokenEndpoint),
					Method:   "POST",
					Headers: map[string][]string{
						"Content-Type":  {"application/json"},
						"Authorization": {"Bearer " + testSessionToken},
					},
					Body: rBodyMarshal(map[string]any{
						"apiVersion": apiVersion,
						"kind":       "Token",
						"metadata": map[string]any{
							"name": testTokenId,
						},
						"spec": map[string]any{
							"description": "Terraform login token.",
							"ttl":         7776000000,
						},
					}),
				},
				c.Response{ // create token
					StatusCode: http.StatusOK,
					Headers:    map[string][]string{"Content-Type": {"application/json"}},
					Body: rBodyMarshal(map[string]any{
						"kind":       "Token",
						"apiVersion": apiVersion,
						"metadata": map[string]any{
							"name":              testTokenId,
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
							"description": "Terraform login token.",
							"ttl":         7776000000,
							"enabled":     true,
						},
						"status": map[string]any{
							"value":          testTokenKey,
							"current":        false,
							"expired":        false,
							"expiresAt":      "2026-05-13T16:49:07-05:00",
							"lastUpdateTime": "2026-02-12T15:49:07-06:00",
							"bearerToken":    testUserToken,
						},
					}),
				},
				RancherLoginModel{ //expected state
					Id:                          testTokenId,
					Username:                    "user",
					Password:                    "password",
					UsernameEnvironmentVariable: "RANCHER_USERNAME",
					PasswordEnvironmentVariable: "RANCHER_PASSWORD",
					UserToken:                   testUserToken,
					TokenTtl:                    "90d",
					RefreshAt:                   "10d",
					IgnoreToken:                 false,
					UserTokenStartDate:          "2026-02-12T21:49:07Z",
					UserTokenEndDate:            "2026-05-13T16:49:07-05:00",
					UserTokenRefreshDate:        "2026-05-03T16:49:07-05:00",
				},
				"success",
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

				ts := &c.TokenStore{}
				ts.SetToken(testInitialToken)
				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, ts)
				err := h.GetConfiguredResource(ctx, t, &tc.fit, client)
				if err != nil {
					t.Errorf("Error configuring resource: %+v", err)
				}

				intialRequestId := fmt.Sprintf("%s:%s:%s", tc.initalTokenRequest.Endpoint, tc.initalTokenRequest.Method, testInitialToken)
				client.SetResponse(ctx, intialRequestId, tc.initialTokenResponse)

				loginRequestId := fmt.Sprintf("%s:%s:%s", tc.loginRequest.Endpoint, tc.loginRequest.Method, "")
				client.SetResponse(ctx, loginRequestId, tc.loginResponse)

				tokenRequestId := fmt.Sprintf("%s:%s:%s", tc.apiRequest.Endpoint, tc.apiRequest.Method, testSessionToken)
				client.SetResponse(ctx, tokenRequestId, tc.apiResponse)

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

				state = tc.expectedState.ToResourceModel(ctx, &dgs).ToState(ctx, &dgs)
				if dgs.HasError() {
					t.Errorf("error generating expected state: %s", pp.PrettyPrint(dgs))
				}
				expectedState := resource.UpdateResponse{
					State: state,
				}

				actualState := res

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
			})
		}
	})
	t.Run("Delete", func(t *testing.T) {
		testCases := []struct {
			name               string
			fit                RancherLoginResource
			existingState      RancherLoginModel
			expectedApiRequest c.Request
			apiResponse        c.Response
			outcome            string
		}{
			{
				"Basic", // delete
				RancherLoginResource{},
				RancherLoginModel{Id: testTokenId},
				c.Request{
					Endpoint: fmt.Sprintf("%s/%s/%s", apiUrl, tokenEndpoint, testTokenId),
					Method:   "DELETE",
					Headers:  map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", testUserToken)}},
				},
				c.Response{
					StatusCode: http.StatusNoContent,
				},
				"success",
			},
			{
				"Resource Already Deleted", // delete
				RancherLoginResource{},
				RancherLoginModel{Id: testTokenId},
				c.Request{
					Endpoint: fmt.Sprintf("%s/%s/%s", apiUrl, tokenEndpoint, testTokenId),
					Method:   "DELETE",
					Headers:  map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", testUserToken)}},
				},
				c.Response{
					StatusCode: http.StatusNotFound,
				},
				"success",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				defer h.PrintLog(t, &buf, "ERROR")
				ctx := h.GenerateTestContext(t, &buf, nil)

				ts := &c.TokenStore{}
				ts.SetToken(testUserToken)
				client := c.NewTestClient(ctx, apiUrl, "", false, false, 30, 10, ts)
				apiRequestId := fmt.Sprintf("%s:%s:%s", tc.expectedApiRequest.Endpoint, tc.expectedApiRequest.Method, testUserToken)
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

				actualApiRequest := client.GetRequest(apiRequestId)
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
