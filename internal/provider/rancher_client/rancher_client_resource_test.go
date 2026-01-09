package rancher_client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	// "time"

	"github.com/google/go-cmp/cmp"
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
					"connect_timeout",
					// "access_key",
					// "secret_key",
					// "token_key",
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
			name    string
			fit     RancherClientResource
			have    RancherClientResourceModel // what is in the plan, translated to struct
			env     map[string]string          // a k/v map of environment variables to set
			want    RancherClientResourceModel // what should be in the state, translated to struct
			effect  c.Client                   // the actual client generated
			outcome string                     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherClientResource{Registry: c.NewRegistry()},
				// plan
				RancherClientResourceModel{
					Id:     types.StringValue(defaultId),
					ApiURL: types.StringValue("https://my-rancher-server.com"),
					// AccessKey: types.StringValue("my-access-key"),
					// SecretKey: types.StringValue("my-secret-key"),
				},
				// env
				map[string]string{},
				// state expected to match this
				RancherClientResourceModel{
					Id:             types.StringValue(defaultId),
					ApiURL:         types.StringValue("https://my-rancher-server.com"),
					CACerts:        types.StringValue(""),
					IgnoreSystemCA: types.BoolValue(false),
					Insecure:       types.BoolValue(false),
					MaxRedirects:   types.StringValue("3"),
					ConnectTimeout: types.StringValue("30s"),
					// AccessKey:      types.StringValue("my-access-key"),
					// SecretKey:      types.StringValue("my-secret-key"),
					// TokenKey:       types.StringValue(""), // state shouldn't save default generated token
				},
				// generated client expected to match this
				c.NewHttpClient(
					context.Background(),
					"https://my-rancher-server.com",
					"",
					false,
					false,
					// "my-access-key",
					// "my-secret-key",
					// "bXktYWNjZXNzLWtleTpteS1zZWNyZXQta2V5", // base64 encoded "my-access-key:my-secret-key"
					3,
					"30s",
				),
				// expected outcome
				"success",
			},
			// Please notice that the state doesn't save the data sent in the environment variables, but the client created (saved in memory) does.
			// This means that when using environment passthrough you must always have the variables present.
			// On the next run the state won't have the attributes configured in the environment, so you must supply them.
			{
				"Environment Passthrough",
				RancherClientResource{Registry: c.NewRegistry()},
				// plan
				RancherClientResourceModel{
					Id:     types.StringValue(defaultId),
					ApiURL: types.StringValue("https://my-rancher-server.com"),
				},
				// env
				map[string]string{
					// "RANCHER_ACCESS_KEY": "my-access-key",
					// "RANCHER_SECRET_KEY": "my-secret-key",
				},
				// state expected to match this
				RancherClientResourceModel{
					Id:             types.StringValue(defaultId),
					ApiURL:         types.StringValue("https://my-rancher-server.com"),
					CACerts:        types.StringValue(""),
					IgnoreSystemCA: types.BoolValue(false),
					Insecure:       types.BoolValue(false),
					MaxRedirects:   types.StringValue("3"),
					ConnectTimeout: types.StringValue("30s"),
					// AccessKey:      types.StringValue(""),
					// SecretKey:      types.StringValue(""),
					// TokenKey:       types.StringValue(""), // state shouldn't save default generated token
				},
				// generated client expected to match this
				c.NewHttpClient(
					context.Background(),
					"https://my-rancher-server.com",
					"",
					false,
					false,
					// "my-access-key",
					// "my-secret-key",
					// "bXktYWNjZXNzLWtleTpteS1zZWNyZXQta2V5", // base64 encoded "my-access-key:my-secret-key"
					3,
					"30s",
				),
				// expected outcome
				"success",
			},
			// In this test we override the written config with environment variables.
			// This is tricky and should be avoided in most cases.
			// If the environment variables aren't present you may end up changing the wrong Rancher instance.
			{
				"Environment Override",
				RancherClientResource{Registry: c.NewRegistry()},
				// plan
				RancherClientResourceModel{
					Id:     types.StringValue(defaultId),
					ApiURL: types.StringValue("https://rancher.example.com"),
					// AccessKey: types.StringValue("fake-access-key"),
					// SecretKey: types.StringValue("fake-secret-key"),
				},
				// env
				map[string]string{
					"RANCHER_API_URL":    "https://rancher-staging.example.com",
					"RANCHER_ACCESS_KEY": "real-access-key",
					"RANCHER_SECRET_KEY": "real-secret-key",
				},
				// state expected to match this
				RancherClientResourceModel{
					Id:             types.StringValue(defaultId),
					ApiURL:         types.StringValue("https://rancher.example.com"),
					CACerts:        types.StringValue(""),
					IgnoreSystemCA: types.BoolValue(false),
					Insecure:       types.BoolValue(false),
					MaxRedirects:   types.StringValue("3"),
					ConnectTimeout: types.StringValue("30s"),
					// AccessKey:      types.StringValue("fake-access-key"),
					// SecretKey:      types.StringValue("fake-secret-key"),
					// TokenKey:       types.StringValue(""), // state shouldn't save default generated token
				},
				// generated client expected to match this
				c.NewHttpClient(
					context.Background(),
					"https://rancher-staging.example.com",
					"",
					false,
					false,
					// "real-access-key",
					// "real-secret-key",
					// "cmVhbC1hY2Nlc3Mta2V5OnJlYWwtc2VjcmV0LWtleQ==", // base64 encoded "real-access-key:real-secret-key"
					3,
					"30s",
				),
				// expected outcome
				"success",
			},
			{
				"Environment Explicit Token",
				RancherClientResource{Registry: c.NewRegistry()},
				// plan
				RancherClientResourceModel{
					Id:     types.StringValue(defaultId),
					ApiURL: types.StringValue("https://rancher.example.com"),
				},
				// env
				map[string]string{
					// "RANCHER_TOKEN_KEY": "cmVhbC1hY2Nlc3Mta2V5OnJlYWwtc2VjcmV0LWtleQ==",
				},
				// state expected to match this
				RancherClientResourceModel{
					Id:             types.StringValue(defaultId),
					ApiURL:         types.StringValue("https://rancher.example.com"),
					CACerts:        types.StringValue(""),
					IgnoreSystemCA: types.BoolValue(false),
					Insecure:       types.BoolValue(false),
					MaxRedirects:   types.StringValue("3"),
					ConnectTimeout: types.StringValue("30s"),
					// AccessKey:      types.StringValue(""),
					// SecretKey:      types.StringValue(""),
					// TokenKey:       types.StringValue(""),
				},
				// generated client expected to match this
				c.NewHttpClient(
					context.Background(),
					"https://rancher.example.com",
					"",
					false,
					false,
					// "",
					// "",
					// "cmVhbC1hY2Nlc3Mta2V5OnJlYWwtc2VjcmV0LWtleQ==", // base64 encoded "real-access-key:real-secret-key"
					3,
					"30s",
				),
				// expected outcome
				"success",
			},
			// This test enables the client to generate its own token if possible
			{
				"Auth Undefined",
				RancherClientResource{Registry: c.NewRegistry()},
				// plan
				RancherClientResourceModel{
					Id:     types.StringValue(defaultId),
					ApiURL: types.StringValue("https://rancher.example.com"),
				},
				// env
				map[string]string{},
				// state expected to match this
				RancherClientResourceModel{
					Id:             types.StringValue(defaultId),
					ApiURL:         types.StringValue("https://rancher.example.com"),
					CACerts:        types.StringValue(""),
					IgnoreSystemCA: types.BoolValue(false),
					Insecure:       types.BoolValue(false),
					MaxRedirects:   types.StringValue("3"),
					ConnectTimeout: types.StringValue("30s"),
					// AccessKey:      types.StringValue(""),
					// SecretKey:      types.StringValue(""),
					// TokenKey:       types.StringValue(""),
				},
				// generated client expected to match this
				c.NewHttpClient(
					context.Background(),
					"https://rancher.example.com",
					"",
					false,
					false,
					// "",
					// "",
					// "",
					3,
					"30s",
				),
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
							getObjectAttributeValues(t, tc.want),
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
				expectedClient := tc.effect
				actualClient, err := tc.fit.Registry.LoadOrError(plannedId)
				if err != nil {
					t.Errorf("Error loading client: %+v", err)
				}
				if diff := cmp.Diff(expectedClient, actualClient); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
				// t.Logf("Resource State: %s", prettyPrint(res))
				if (tc.outcome == "failure") && !res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() did not return expected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
				}
				if (tc.outcome == "success") && res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
				}
				if diff := cmp.Diff(expectedState, res); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
}

func TestRancherClientResourceRead(t *testing.T) {
	t.Run("Read function", func(t *testing.T) {
		testCases := []struct {
			name    string
			fit     RancherClientResource
			have    RancherClientResourceModel // what is in the plan, translated to struct
			env     map[string]string          // a k/v map of environment variables to set
			want    RancherClientResourceModel // what should be in the state, translated to struct
			effect  c.Client                   // the actual client generated
			outcome string                     // expected outcome, one of: "success","failure"
		}{
			{
				"Basic",
				RancherClientResource{Registry: c.NewRegistry()},
				// plan
				RancherClientResourceModel{
					Id:     types.StringValue(defaultId),
					ApiURL: types.StringValue("https://my-rancher-server.com"),
					// AccessKey: types.StringValue("my-access-key"),
					// SecretKey: types.StringValue("my-secret-key"),
				},
				// env
				map[string]string{},
				// state expected to match this
				RancherClientResourceModel{
					Id:             types.StringValue(defaultId),
					ApiURL:         types.StringValue("https://my-rancher-server.com"),
					CACerts:        types.StringValue(""),
					IgnoreSystemCA: types.BoolValue(false),
					Insecure:       types.BoolValue(false),
					MaxRedirects:   types.StringValue("3"),
					ConnectTimeout: types.StringValue("30s"),
					// AccessKey:      types.StringValue("my-access-key"),
					// SecretKey:      types.StringValue("my-secret-key"),
					// TokenKey:       types.StringValue(""), // state shouldn't save default generated token
				},
				// generated client expected to match this
				c.NewHttpClient(
					context.Background(),
					"https://my-rancher-server.com",
					"",
					false,
					false,
					// "my-access-key",
					// "my-secret-key",
					// "bXktYWNjZXNzLWtleTpteS1zZWNyZXQta2V5", // base64 encoded "my-access-key:my-secret-key"
					3,
					"30s",
				),
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
							getObjectAttributeValues(t, tc.want),
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
				expectedClient := tc.effect
				actualClient, err := tc.fit.Registry.LoadOrError(plannedId)
				if err != nil {
					t.Errorf("Error loading client: %+v", err)
				}
				if diff := cmp.Diff(expectedClient, actualClient); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
				// t.Logf("Resource State: %s", prettyPrint(res))
				if (tc.outcome == "failure") && !res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() did not return expected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
				}
				if (tc.outcome == "success") && res.Diagnostics.HasError() {
					t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, prettyPrint(res.Diagnostics))
				}
				if diff := cmp.Diff(expectedState, res); diff != "" {
					t.Errorf("Create() mismatch (-want +got):\n%+v", diff)
				}
			})
		}
	})
}

// The helpers.
func getSchema() schema.Schema {
	testResource := RancherClientResource{}
	r := resource.SchemaResponse{}
	testResource.Schema(context.Background(), resource.SchemaRequest{}, &r)
	return r.Schema
}
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

// getObjectAttributeValues converts the RancherClientResourceModel struct to a map[string]tftypes.Value.
// it parses the schema to get the attribute names and types so that it automatically adapts to schema changes.
func getObjectAttributeValues(t *testing.T, config RancherClientResourceModel) map[string]tftypes.Value {
	values := map[string]tftypes.Value{}
	attributeTypes := getObjectAttributeTypes().AttributeTypes
	for attrName, attrType := range attributeTypes {
		var value interface{}
		// use reflect to get the value from the struct based on the attribute name variable.
		v := reflect.ValueOf(config)
		fieldVal, err := getStructFieldByTfsdkTag(v, attrName)
		if err != nil {
			t.Fatalf("getObjectAttributeValues: %v", err)
		}
		if !fieldVal.IsValid() {
			t.Fatalf("getObjectAttributeValues: no such field %s in model", attrName)
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
		value = results[0].Interface()

		if value == nil {
			// for nil values, set tftypes.UnknownValue of the appropriate type.
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
