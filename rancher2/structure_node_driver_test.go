package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testNodeDriverConf      *managementClient.NodeDriver
	testNodeDriverInterface map[string]interface{}
)

func init() {
	testNodeDriverConf = &managementClient.NodeDriver{
		Active:             true,
		AddCloudCredential: true,
		Builtin:            true,
		Checksum:           "XXXXXXXX",
		Description:        "description",
		ExternalID:         "external",
		Name:               "name",
		UIURL:              "ui_url",
		URL:                "url",
		WhitelistDomains:   []string{"domain1", "domain2"},
	}
	testNodeDriverInterface = map[string]interface{}{
		"active":               true,
		"add_cloud_credential": true,
		"builtin":              true,
		"checksum":             "XXXXXXXX",
		"description":          "description",
		"external_id":          "external",
		"name":                 "name",
		"ui_url":               "ui_url",
		"url":                  "url",
		"whitelist_domains":    []interface{}{"domain1", "domain2"},
	}
}

func TestFlattenNodeDriver(t *testing.T) {

	cases := []struct {
		Input          *managementClient.NodeDriver
		ExpectedOutput map[string]interface{}
	}{
		{
			testNodeDriverConf,
			testNodeDriverInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, nodeDriverFields(), map[string]interface{}{})
		err := flattenNodeDriver(output, tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			assert.FailNow(t, "Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, output)
		}
	}
}

func TestExpandNodeDriver(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.NodeDriver
	}{
		{
			testNodeDriverInterface,
			testNodeDriverConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, nodeDriverFields(), tc.Input)
		output := expandNodeDriver(inputResourceData)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
