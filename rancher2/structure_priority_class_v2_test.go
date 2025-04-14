package rancher2

import (
	"testing"

	"github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func Test_expandClusterAgentPriorityClassV2(t *testing.T) {
	var preemptionNever = corev1.PreemptionPolicy("Never")

	tests := []struct {
		name           string
		expectedOutput *v1.PriorityClassSpec
		input          []interface{}
	}{
		{
			name: "both fields set",
			expectedOutput: &v1.PriorityClassSpec{
				Value:            123,
				PreemptionPolicy: &preemptionNever,
			},
			input: []interface{}{
				map[string]interface{}{
					"value":             123,
					"preemption_policy": "Never",
				},
			},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			out := expandClusterAgentPriorityClassV2(test.input)
			assert.Equal(t, test.expectedOutput, out, "unexpected output from expander, expected %v, got %v", test.expectedOutput, out)
		})
	}
}

func Test_flattenClusterAgentPriorityClassV2(t *testing.T) {
	var preemptionNever = corev1.PreemptionPolicy("Never")

	tests := []struct {
		name           string
		expectedOutput []interface{}
		input          *v1.PriorityClassSpec
	}{
		{
			name: "both fields set",
			input: &v1.PriorityClassSpec{
				Value:            123,
				PreemptionPolicy: &preemptionNever,
			},
			expectedOutput: []interface{}{
				map[string]interface{}{
					"value":             123,
					"preemption_policy": preemptionNever,
				},
			},
		},
		{
			name: "empty preemption",
			input: &v1.PriorityClassSpec{
				Value: 123,
			},
			expectedOutput: []interface{}{
				map[string]interface{}{
					"value": 123,
				},
			},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			out := flattenClusterAgentPriorityClassV2(test.input)
			assert.Equal(t, test.expectedOutput, out, "unexpected output from flattener, expected %v, got %v", test.expectedOutput, out)
		})
	}
}
