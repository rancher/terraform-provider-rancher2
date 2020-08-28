package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

var (
	testProjectSecretConf         *projectClient.Secret
	testProjectSecretInterface    map[string]interface{}
	testNamespacedSecretConf      *projectClient.NamespacedSecret
	testNamespacedSecretInterface map[string]interface{}
)

func init() {
	testProjectSecretConf = &projectClient.Secret{
		ProjectID:   "project:test",
		Name:        "name",
		Description: "description",
		Data: map[string]string{
			"data_one": "one",
			"data_two": "two",
		},
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testProjectSecretInterface = map[string]interface{}{
		"project_id":  "project:test",
		"name":        "name",
		"description": "description",
		"data": map[string]interface{}{
			"data_one": "one",
			"data_two": "two",
		},
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNamespacedSecretConf = &projectClient.NamespacedSecret{
		ProjectID:   "project:test",
		Name:        "name",
		Description: "description",
		NamespaceId: "namespace_id",
		Data: map[string]string{
			"data_one": "one",
			"data_two": "two",
		},
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNamespacedSecretInterface = map[string]interface{}{
		"project_id":   "project:test",
		"name":         "name",
		"description":  "description",
		"namespace_id": "namespace_id",
		"data": map[string]interface{}{
			"data_one": "one",
			"data_two": "two",
		},
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
}

func TestFlattenSecret(t *testing.T) {

	cases := []struct {
		Input          interface{}
		ExpectedOutput map[string]interface{}
	}{
		{
			testProjectSecretConf,
			testProjectSecretInterface,
		},
		{
			testNamespacedSecretConf,
			testNamespacedSecretInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, secretFields(), tc.ExpectedOutput)
		err := flattenSecret(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, tc.ExpectedOutput)
		}
	}
}

func TestExpandSecret(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput interface{}
	}{
		{
			testProjectSecretInterface,
			testProjectSecretConf,
		},
		{
			testNamespacedSecretInterface,
			testNamespacedSecretConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, secretFields(), tc.Input)
		output := expandSecret(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
