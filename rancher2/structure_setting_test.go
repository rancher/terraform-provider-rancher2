package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testSettingConf      *managementClient.Setting
	testSettingInterface map[string]interface{}
)

func init() {
	testSettingConf = &managementClient.Setting{
		Name:  "foo",
		Value: "Terraform setting acceptance test",
	}
	testSettingInterface = map[string]interface{}{
		"name":  "foo",
		"value": "Terraform setting acceptance test",
	}
}

func TestFlattenSetting(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Setting
		ExpectedOutput map[string]interface{}
	}{
		{
			testSettingConf,
			testSettingInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, settingFields(), map[string]interface{}{})
		err := flattenSetting(output, tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		assert.Equal(t, tc.ExpectedOutput, expectedOutput, "Unexpected output from flattener.")
	}
}

func TestExpandSetting(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.Setting
	}{
		{
			testSettingInterface,
			testSettingConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, settingFields(), tc.Input)
		output, err := expandSetting(inputResourceData)
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
