package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testTolerationsConf      []managementClient.Toleration
	testTolerationsInterface []interface{}
)

func init() {
	seconds := int64(10)
	testTolerationsConf = []managementClient.Toleration{
		{
			Key:               "key",
			Value:             "value",
			Effect:            "NoSchedule",
			Operator:          "Equal",
			TolerationSeconds: &seconds,
		},
	}
	testTolerationsInterface = []interface{}{
		map[string]interface{}{
			"key":      "key",
			"value":    "value",
			"effect":   "NoSchedule",
			"operator": "Equal",
			"seconds":  10,
		},
	}
}

func TestFlattenTolerations(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Toleration
		ExpectedOutput []interface{}
	}{
		{
			testTolerationsConf,
			testTolerationsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenTolerations(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandTolerations(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Toleration
	}{
		{
			testTolerationsInterface,
			testTolerationsConf,
		},
	}

	for _, tc := range cases {
		output := expandTolerations(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
