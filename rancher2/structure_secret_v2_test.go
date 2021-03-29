package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var (
	testSecretV2Conf      *SecretV2
	testSecretV2Interface map[string]interface{}
)

func init() {
	testSecretV2Conf = &SecretV2{}

	testSecretV2Conf.TypeMeta.Kind = secretV2Kind
	testSecretV2Conf.TypeMeta.APIVersion = secretV2APIVersion

	testSecretV2Conf.ObjectMeta.Name = "name"
	testSecretV2Conf.ObjectMeta.Namespace = "namespace"
	testSecretV2Conf.ObjectMeta.Annotations = map[string]string{
		"value1": "one",
		"value2": "two",
	}
	testSecretV2Conf.ObjectMeta.Labels = map[string]string{
		"label1": "one",
		"label2": "two",
	}
	testSecretV2Conf.Immutable = newTrue()
	testSecretV2Conf.Resource.Type = "type"
	testSecretV2Conf.K8SType = "type"
	testSecretV2Conf.StringData = map[string]string{
		"data1": "one",
		"data2": "two",
	}

	testSecretV2Interface = map[string]interface{}{
		"name":      "name",
		"namespace": "namespace",
		"immutable": true,
		"type":      "type",
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

func TestFlattenSecretV2(t *testing.T) {

	cases := []struct {
		Input          *SecretV2
		ExpectedOutput map[string]interface{}
	}{
		{
			testSecretV2Conf,
			testSecretV2Interface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, secretV2Fields(), tc.ExpectedOutput)
		err := flattenSecretV2(output, tc.Input)
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

func TestExpandSecretV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *SecretV2
	}{
		{
			testSecretV2Interface,
			testSecretV2Conf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, secretV2Fields(), tc.Input)
		output := expandSecretV2(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
