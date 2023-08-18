package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

var (
	testTaintsV2Conf      []corev1.Taint
	testTaintsV2Interface []interface{}
)

func init() {
	testTaintsV2Conf = []corev1.Taint{
		{
			Key:    "key",
			Value:  "value",
			Effect: "recipient",
		},
	}
	testTaintsV2Interface = []interface{}{
		map[string]interface{}{
			"key":    "key",
			"value":  "value",
			"effect": "recipient",
		},
	}
}

func TestFlattenTaintsV2(t *testing.T) {

	cases := []struct {
		Input          []corev1.Taint
		ExpectedOutput []interface{}
	}{
		{
			testTaintsV2Conf,
			testTaintsV2Interface,
		},
	}
	for _, tc := range cases {
		output := flattenTaintsV2(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandTaintsV2(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []corev1.Taint
	}{
		{
			testTaintsV2Interface,
			testTaintsV2Conf,
		},
	}

	for _, tc := range cases {
		output := expandTaintsV2(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
