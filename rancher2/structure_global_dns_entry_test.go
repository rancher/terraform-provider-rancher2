package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testGlobalGlobalDNSEntryConf      *managementClient.GlobalDNS
	testGlobalGlobalDNSEntryInterface map[string]interface{}
)

func init() {
	testGlobalGlobalDNSEntryConf = &managementClient.GlobalDNS{
		FQDN:       "test.non.example.com",
		Name:       "test-entry",
		ProviderID: "cattle-global:foo-test2",
		ProjectIDs: []string{"cxs", "bbbb"},
	}
	testGlobalGlobalDNSEntryInterface = map[string]interface{}{
		"name":        "test-entry",
		"provider_id": "cattle-global:foo-test2",
		"fqdn":        "test.non.example.com",
		"project_ids": []interface{}{"cxs", "bbbb"},
	}
}

func TestFlattenGlobalDNSEntry(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GlobalDNS
		ExpectedOutput map[string]interface{}
	}{
		{
			testGlobalGlobalDNSEntryConf,
			testGlobalGlobalDNSEntryInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, GlobalDNSEntryFields(), map[string]interface{}{})
		err := flattenGlobalDNSEntry(output, tc.Input)
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

func TestExpandGlobalDNSEntry(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.GlobalDNS
	}{
		{
			testGlobalGlobalDNSEntryInterface,
			testGlobalGlobalDNSEntryConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, GlobalDNSEntryFields(), tc.Input)
		output, err := expandGlobalDNSEntry(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
