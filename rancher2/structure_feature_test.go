package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testFeatureConf      *managementClient.Feature
	testFeatureInterface map[string]interface{}
)

func init() {
	testFeatureConf = &managementClient.Feature{
		Name:  "foo",
		Value: newTrue(),
	}
	testFeatureInterface = map[string]interface{}{
		"name":  "foo",
		"value": true,
	}
}

func TestFlattenFeature(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Feature
		ExpectedOutput map[string]interface{}
	}{
		{
			testFeatureConf,
			testFeatureInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, featureFields(), map[string]interface{}{})
		err := flattenFeature(output, tc.Input)
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

func TestExpandFeature(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.Feature
	}{
		{
			testFeatureInterface,
			testFeatureConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, featureFields(), tc.Input)
		output, err := expandFeature(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven: %#v", tc.ExpectedOutput, output)
		}
	}
}
