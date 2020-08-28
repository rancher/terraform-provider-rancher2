package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven: %#v", tc.ExpectedOutput, output)
		}
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
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven: %#v", tc.ExpectedOutput, output)
		}
	}
}
