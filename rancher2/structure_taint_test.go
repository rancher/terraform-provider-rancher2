package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testTaintsConf      []managementClient.Taint
	testTaintsInterface []interface{}
)

func init() {
	testTaintsConf = []managementClient.Taint{
		{
			Key:       "key",
			Value:     "value",
			Effect:    "recipient",
			TimeAdded: "time_added",
		},
	}
	testTaintsInterface = []interface{}{
		map[string]interface{}{
			"key":        "key",
			"value":      "value",
			"effect":     "recipient",
			"time_added": "time_added",
		},
	}
}

func TestFlattenTaints(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Taint
		ExpectedOutput []interface{}
	}{
		{
			testTaintsConf,
			testTaintsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenTaints(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandTaints(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Taint
	}{
		{
			testTaintsInterface,
			testTaintsConf,
		},
	}

	for _, tc := range cases {
		output := expandTaints(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
