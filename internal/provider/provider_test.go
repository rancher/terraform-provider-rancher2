package provider

import (
	"context"
	"testing"
  // "time"
  "strings"
  "reflect"

	"github.com/hashicorp/terraform-plugin-framework/provider"
  "github.com/hashicorp/terraform-plugin-framework/provider/schema"
  // "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
  "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
  // "github.com/hashicorp/terraform-plugin-log/tflog"
)

func TestProviderMetadata(t *testing.T) {
	testCases := []struct {
		name string
		fit  RancherProvider
		want provider.MetadataResponse
	}{
		{
      "Metadata function basic",
      RancherProvider{version: "test"},
      provider.MetadataResponse{TypeName: "rancher", Version: "test"},
    },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := provider.MetadataRequest{}
			res := provider.MetadataResponse{}
			tc.fit.Metadata(ctx, req, &res)
			got := res
			if got != tc.want {
				t.Errorf("%#v.Metadata() is %v; want %v", tc.fit, got, tc.want)
			}
		})
	}
}

func TestProviderSchema(t *testing.T) {
  testCases := []struct {
    name string
    fit  RancherProvider
    want []string
  }{
    {
      "Schema function basic",
      RancherProvider{version: "test"},
      []string{
        "api_url",
        "ca_certs",
        "insecure",
        "access_key",
        "secret_key",
        "token_key",
        "timeout",
        "bootstrap",
      },
    },
  }
  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      ctx := context.Background()
      req := provider.SchemaRequest{}
      res := provider.SchemaResponse{}
      tc.fit.Schema(ctx, req, &res)
      gotKeys := []string{}
      for key := range res.Schema.Attributes {
        gotKeys = append(gotKeys, key)
      }
      for _, wantKey := range tc.want {
        found := false
        for _, gotKey := range gotKeys {
          if wantKey == gotKey {
            found = true
            break
          }
        }
        if !found {
          t.Errorf("%#v.Schema() missing expected key %s", tc.fit, wantKey)
        }
      }
    })
  }
}

// test provider configuration giving some basic values
func TestProviderConfigure(t *testing.T) {
  t.Run("Configure function basic", func(t *testing.T) {
    testCases := []struct {
      name string
      fit  RancherProvider
      config RancherProviderModel
    }{
      {
        "Configure function basic",
        RancherProvider{version: "test"},
        RancherProviderModel{
          ApiURL:    types.StringValue("https://my-rancher-server.com"),
          CACerts:   types.StringValue(""),
          Insecure:  types.BoolValue(true),
          AccessKey: types.StringValue("my-access-key"),
          SecretKey: types.StringValue("my-secret-key"),
          TokenKey:  types.StringValue(""),
          Timeout:   types.StringValue("30s"),
          Bootstrap: types.BoolValue(false),
        },
      },
    }
    for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
        ctx := context.Background()
        req := provider.ConfigureRequest{ Config: tfsdk.Config{
          Raw: tftypes.NewValue(
            getObjectAttributeTypes(),
            getObjectAttributeValues(t, tc.config),
          ),
        }}
        res := provider.ConfigureResponse{}
        tc.fit.Configure(ctx, req, &res)
        if res.Diagnostics.HasError() {
          t.Errorf("%#v.Configure() returned error diagnostics: %v", tc.fit, res.Diagnostics)
        }
      })
    }
  })
}

func getSchema() schema.Schema {
  testProvider := RancherProvider{version: "test"}
  r := provider.SchemaResponse{}
  testProvider.Schema(context.Background(), provider.SchemaRequest{}, &r)
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

// getObjectAttributeValues converts the RancherProviderModel struct to a map[string]tftypes.Value
// it parses the schema to get the attribute names and types so that it automatically adapts to schema changes
func getObjectAttributeValues(t *testing.T, config RancherProviderModel) map[string]tftypes.Value {
  schema := getSchema()
  values := map[string]tftypes.Value{}
  attributeTypes := getObjectAttributeTypes().AttributeTypes
  for attrName, attr := range schema.GetAttributes() {
    attrType, ok := attributeTypes[attrName]
    if !ok {
      t.Fatalf("getObjectAttributeValues: attribute type for %s not found", attrName)
    }
    var value interface{}
    // use reflect to get the value from the struct based on the attribute name variable
    v := reflect.ValueOf(config)
    fieldVal := v.FieldByNameFunc(func(s string) bool {
      return strings.EqualFold(s, attrName)
    })
    if !fieldVal.IsValid() {
      t.Fatalf("getObjectAttributeValues: no such field %s in RancherProviderModel", attrName)
    }
    // get the underlying value from the types.Value
    switch fv := fieldVal.Interface().(type) {
    case types.String:
      value = fv.ValueString()
    case types.Bool:
      value = fv.ValueBool()
    default:
      t.Fatalf("getObjectAttributeValues: unsupported type %T for field %s", fv, attrName)
    }
    values[attrName] = tftypes.NewValue(attrType, value)
  }
  return values
}
