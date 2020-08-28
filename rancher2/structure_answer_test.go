package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAnswersConf      []managementClient.Answer
	testAnswersInterface []interface{}
)

func init() {
	testAnswersConf = []managementClient.Answer{
		{
			ClusterID: "cluster_id",
			ProjectID: "project_id",
			Values: map[string]string{
				"value1": "one",
				"value2": "two",
			},
		},
	}
	testAnswersInterface = []interface{}{
		map[string]interface{}{
			"cluster_id": "cluster_id",
			"project_id": "project_id",
			"values": map[string]interface{}{
				"value1": "one",
				"value2": "two",
			},
		},
	}
}

func TestFlattenAnswers(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Answer
		ExpectedOutput []interface{}
	}{
		{
			testAnswersConf,
			testAnswersInterface,
		},
	}

	for _, tc := range cases {
		output := flattenAnswers(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandAnswers(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Answer
	}{
		{
			testAnswersInterface,
			testAnswersConf,
		},
	}

	for _, tc := range cases {
		output := expandAnswers(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
