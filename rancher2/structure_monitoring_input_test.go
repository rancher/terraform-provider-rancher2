package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
