package rancher2

import (
	"testing"

	"github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

func Test_flattenPodDisruptionBudgetV2(t *testing.T) {
	tests := []struct {
		name           string
		input          *v1.PodDisruptionBudgetSpec
		expectedOutput []interface{}
	}{
		{
			name: "min available set",
			input: &v1.PodDisruptionBudgetSpec{
				MinAvailable: "1",
			},
			expectedOutput: []interface{}{
				map[string]interface{}{
					"min_available": "1",
				},
			},
		},
		{
			name: "max unavailable set",
			input: &v1.PodDisruptionBudgetSpec{
				MaxUnavailable: "1",
			},
			expectedOutput: []interface{}{
				map[string]interface{}{
					"max_unavailable": "1",
				},
			},
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			out := flattenPodDisruptionBudgetV2(tt.input)
			assert.Equal(t, tt.expectedOutput, out, "unexpected output from flattener, expected %v, got %v", tt.expectedOutput, out)
		})
	}
}

func Test_expandPodDisruptionBudgetV2(t *testing.T) {
	tests := []struct {
		name           string
		input          []interface{}
		expectedOutput *v1.PodDisruptionBudgetSpec
	}{
		{
			name: "min available set",
			input: []interface{}{
				map[string]interface{}{
					"min_available": "1",
				},
			},
			expectedOutput: &v1.PodDisruptionBudgetSpec{
				MinAvailable: "1",
			},
		},
		{
			name: "max unavailable set",
			input: []interface{}{
				map[string]interface{}{
					"max_unavailable": "1",
				},
			},
			expectedOutput: &v1.PodDisruptionBudgetSpec{
				MaxUnavailable: "1",
			},
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			out := expandPodDisruptionBudgetV2(tt.input)
			assert.Equal(t, tt.expectedOutput, out, "unexpected output from expander, expected %v, got %v", tt.expectedOutput, out)
		})
	}
}
