package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testMonitoringInputConf      *managementClient.MonitoringInput
	testMonitoringInputInterface []interface{}
)

func init() {
	testMonitoringInputConf = &managementClient.MonitoringInput{
		Answers: map[string]string{
			"answer_one": "one",
			"answer_two": "two",
		},
	}
	testMonitoringInputInterface = []interface{}{
		map[string]interface{}{
			"answers": map[string]interface{}{
				"answer_one": "one",
				"answer_two": "two",
			},
		},
	}
}

func TestFlattenMonitoringInput(t *testing.T) {

	cases := []struct {
		Input          *managementClient.MonitoringInput
		ExpectedOutput []interface{}
	}{
		{
			testMonitoringInputConf,
			testMonitoringInputInterface,
		},
	}

	for _, tc := range cases {
		output := flattenMonitoringInput(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandMonitoringInput(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.MonitoringInput
	}{
		{
			testMonitoringInputInterface,
			testMonitoringInputConf,
		},
	}

	for _, tc := range cases {
		output := expandMonitoringInput(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
