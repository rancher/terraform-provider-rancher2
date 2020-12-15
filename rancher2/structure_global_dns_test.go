package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testGlobalGlobalDNSConf      *managementClient.GlobalDns
	testGlobalGlobalDNSInterface map[string]interface{}
)

func init() {
	testGlobalGlobalDNSConf = &managementClient.GlobalDns{
		FQDN:              "test.non.example.com",
		ProviderID:        "cattle-global:foo-test2",
		MultiClusterAppID: "mca",
		Name:              "test-entry",
		ProjectIDs:        []string{"proj1", "proj2"},
		TTL:               int64(600),
	}
	testGlobalGlobalDNSInterface = map[string]interface{}{
		"fqdn":                 "test.non.example.com",
		"provider_id":          "cattle-global:foo-test2",
		"multi_cluster_app_id": "mca",
		"name":                 "test-entry",
		"project_ids":          []interface{}{"proj1", "proj2"},
		"ttl":                  600,
	}
}

func TestFlattenGlobalDNS(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GlobalDns
		ExpectedOutput map[string]interface{}
	}{
		{
			testGlobalGlobalDNSConf,
			testGlobalGlobalDNSInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, GlobalDNSFields(), map[string]interface{}{})
		err := flattenGlobalDNS(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, output)
		}
	}
}

func TestExpandGlobalDNS(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.GlobalDns
	}{
		{
			testGlobalGlobalDNSInterface,
			testGlobalGlobalDNSConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, GlobalDNSFields(), tc.Input)
		output, err := expandGlobalDNS(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
