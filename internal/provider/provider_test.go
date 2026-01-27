package provider

import (
	"bytes"
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	h "github.com/rancher/terraform-provider-rancher2/internal/provider/test_helpers"
)

func TestProviderMetadata(t *testing.T) {
	testCases := []struct {
		name string
		fit  RancherProvider
		want provider.MetadataResponse
	}{
		{
			"Basic",
			RancherProvider{version: "test"},
			provider.MetadataResponse{TypeName: "rancher2", Version: "test"},
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
				t.Errorf("%#v.Metadata() is %s; want %s", tc.fit, pp.PrettyPrint(got), pp.PrettyPrint(tc.want))
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
			"Basic",
			RancherProvider{version: "test"},
			[]string{
				"api_url",
				"ca_cert",
				"insecure",
				"ignore_system_ca",
				"timeout",
				"max_redirects",
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
				found := slices.Contains(gotKeys, wantKey)
				if !found {
					t.Errorf("%#v.Schema() missing expected key %s", tc.fit, wantKey)
				}
			}
		})
	}
}

func TestProviderConfigure(t *testing.T) {
	testCases := []struct {
		name    string
		fit     RancherProvider
		have    RancherProviderModel
		outcome string
	}{
		{
			"Basic",
			RancherProvider{version: "test"},
			RancherProviderModel{
				ApiUrl: types.StringValue("https://test-rancher.example.com"),
			},
			"succeed",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := h.GenerateTestContext(t, &buf, nil)
			defer h.PrintLog(t, &buf, "ERROR") // change to debug when troubleshooting

			req := provider.ConfigureRequest{Config: h.GetConfig(t, &tc.fit, tc.have)}
			res := provider.ConfigureResponse{}
			tc.fit.Configure(ctx, req, &res)
			if (tc.outcome == "succeed") && res.Diagnostics.HasError() {
				t.Errorf("%#v.Configure() returned unexpected error diagnostics: %s", tc.fit, pp.PrettyPrint(res.Diagnostics))
			}
			if (tc.outcome == "fail") && !res.Diagnostics.HasError() {
				t.Errorf("%#v.Configure() did not return expected error diagnostics: %s", tc.fit, pp.PrettyPrint(res.Diagnostics))
			}
		})
	}
}
