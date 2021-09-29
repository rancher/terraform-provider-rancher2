package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var (
	testConfigMapV2Conf      *ConfigMapV2
	testConfigMapV2Interface map[string]interface{}
)

func init() {
	testConfigMapV2Conf = &ConfigMapV2{}

	testConfigMapV2Conf.TypeMeta.Kind = configMapV2Kind
	testConfigMapV2Conf.TypeMeta.APIVersion = configMapV2APIVersion

	testConfigMapV2Conf.ObjectMeta.Name = "name"
	testConfigMapV2Conf.ObjectMeta.Namespace = "namespace"
	testConfigMapV2Conf.ObjectMeta.Annotations = map[string]string{
		"value1": "one",
		"value2": "two",
	}
	testConfigMapV2Conf.ObjectMeta.Labels = map[string]string{
		"label1": "one",
		"label2": "two",
	}
	testConfigMapV2Conf.Immutable = newTrue()
	testConfigMapV2Conf.Data = map[string]string{
		"data1": "one",
		"data2": "two",
	}

	testConfigMapV2Interface = map[string]interface{}{
		"name":      "name",
		"namespace": "namespace",
		"immutable": true,
		"data": map[string]interface{}{
			"data1": "one",
			"data2": "two",
		},
		"annotations": map[string]interface{}{
			"value1": "one",
			"value2": "two",
		},
		"labels": map[string]interface{}{
			"label1": "one",
			"label2": "two",
		},
	}
}

func TestFlattenConfigMapV2(t *testing.T) {

	cases := []struct {
		Input          *ConfigMapV2
		ExpectedOutput map[string]interface{}
	}{
		{
			testConfigMapV2Conf,
			testConfigMapV2Interface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, configMapV2Fields(), tc.ExpectedOutput)
		err := flattenConfigMapV2(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandConfigMapV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *ConfigMapV2
	}{
		{
			testConfigMapV2Interface,
			testConfigMapV2Conf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, configMapV2Fields(), tc.Input)
		output := expandConfigMapV2(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
