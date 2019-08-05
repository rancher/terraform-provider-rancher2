package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	projectClient "github.com/rancher/types/client/project/v3"
)

var (
	testAppConf      *projectClient.App
	testAppInterface map[string]interface{}
)

func init() {
	testAppConf = &projectClient.App{
		ProjectID:   "project:test",
		Name:        "name",
		Description: "description",
		Answers: map[string]string{
			"answers1": "one",
			"answers2": "two",
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
	testAppInterface = map[string]interface{}{
		"project_id":  "project:test",
		"name":        "name",
		"description": "description",
		"answers": map[string]interface{}{
			"answers1": "one",
			"answers2": "two",
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

func TestFlattenApp(t *testing.T) {

	cases := []struct {
		Input          *projectClient.App
		ExpectedOutput map[string]interface{}
	}{
		{
			testAppConf,
			testAppInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, appFields(), tc.ExpectedOutput)
		err := flattenApp(output, tc.Input)
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

func TestExpandApp(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput interface{}
	}{
		{
			testAppInterface,
			testAppConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, appFields(), tc.Input)
		output := expandApp(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
