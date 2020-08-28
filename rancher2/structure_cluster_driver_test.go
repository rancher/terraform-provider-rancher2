package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterDriverConf      *managementClient.KontainerDriver
	testClusterDriverInterface map[string]interface{}
)

func init() {
	testClusterDriverConf = &managementClient.KontainerDriver{
		Active:           true,
		ActualURL:        "actual_url",
		BuiltIn:          true,
		Checksum:         "XXXXXXXX",
		Name:             "name",
		UIURL:            "ui_url",
		URL:              "url",
		WhitelistDomains: []string{"domain1", "domain2"},
	}
	testClusterDriverInterface = map[string]interface{}{
		"active":            true,
		"actual_url":        "actual_url",
		"builtin":           true,
		"checksum":          "XXXXXXXX",
		"name":              "name",
		"ui_url":            "ui_url",
		"url":               "url",
		"whitelist_domains": []interface{}{"domain1", "domain2"},
	}
}

func TestFlattenClusterDriver(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KontainerDriver
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterDriverConf,
			testClusterDriverInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterDriverFields(), map[string]interface{}{})
		err := flattenClusterDriver(output, tc.Input)
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

func TestExpandClusterDriver(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.KontainerDriver
	}{
		{
			testClusterDriverInterface,
			testClusterDriverConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterDriverFields(), tc.Input)
		output := expandClusterDriver(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
