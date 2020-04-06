package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenClusterBaseNodePool(t *testing.T) {
	tests := map[string]struct {
		nodePool BaseNodePool
		expected map[string]interface{}
	}{
		"RequiredValues": {
			nodePool: BaseNodePool{Name: "pool-1"},
			expected: map[string]interface{}{
				"add_default_label": false,
				"add_default_taint": false,
				"name":              "pool-1",
			},
		},
		"AddDefaultLabelTrue": {
			nodePool: BaseNodePool{
				AddDefaultLabel: true,
				Name:            "pool-1",
			},
			expected: map[string]interface{}{
				"add_default_label": true,
				"add_default_taint": false,
				"name":              "pool-1",
			},
		},
		"AddDefaultTaintTrue": {
			nodePool: BaseNodePool{
				AddDefaultTaint: true,
				Name:            "pool-1",
			},
			expected: map[string]interface{}{
				"add_default_label": false,
				"add_default_taint": true,
				"name":              "pool-1",
			},
		},
		"AdditionalLabelsSet": {
			nodePool: BaseNodePool{
				AdditionalLabels: map[string]string{
					"label-1": "value-1",
					"label-2": "value-2",
				},
				Name: "pool-1",
			},
			expected: map[string]interface{}{
				"add_default_label": false,
				"add_default_taint": false,
				"additional_labels": map[string]interface{}{
					"label-1": "value-1",
					"label-2": "value-2",
				},
				"name": "pool-1",
			},
		},
		"AdditionalTaintsSet": {
			nodePool: BaseNodePool{
				AdditionalTaints: []K8sTaint{
					{Effect: "NoSchedule", Key: "taint-1", Value: "value-1"},
					{Effect: "NoExecute", Key: "taint-2", Value: "value-2"},
				},
				Name: "pool-1",
			},
			expected: map[string]interface{}{
				"add_default_label": false,
				"add_default_taint": false,
				"additional_taints": []interface{}{
					map[string]interface{}{"effect": "NoSchedule", "key": "taint-1", "value": "value-1"},
					map[string]interface{}{"effect": "NoExecute", "key": "taint-2", "value": "value-2"},
				},
				"name": "pool-1",
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			output := flattenClusterBaseNodePool(input.nodePool)

			assert.Equal(t, output, input.expected)
		})
	}
}

func TestExpandClusterBaseNodePool(t *testing.T) {
	tests := map[string]struct {
		nodePool         map[string]interface{}
		successfulResult BaseNodePool
		failureResult    string
	}{
		"NameNotSet": {
			nodePool:      map[string]interface{}{},
			failureResult: "'name' field must be provided for all pools",
		},
		"RequiredFieldsAreSet": {
			nodePool: map[string]interface{}{
				"name": "a-name",
			},
			successfulResult: BaseNodePool{
				Name: "a-name",
			},
		},
		"DefaultLabelEnabled": {
			nodePool: map[string]interface{}{
				"name":              "a-name",
				"add_default_label": true,
			},
			successfulResult: BaseNodePool{
				Name:            "a-name",
				AddDefaultLabel: true,
			},
		},
		"DefaultTaintEnabled": {
			nodePool: map[string]interface{}{
				"name":              "a-name",
				"add_default_taint": true,
			},
			successfulResult: BaseNodePool{
				Name:            "a-name",
				AddDefaultTaint: true,
			},
		},
		"AdditionalLabelsSet": {
			nodePool: map[string]interface{}{
				"name": "a-name",
				"additional_labels": map[string]interface{}{
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
			successfulResult: BaseNodePool{
				Name: "a-name",
				AdditionalLabels: map[string]string{
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
		},
		"AdditionalTaintsSet": {
			nodePool: map[string]interface{}{
				"name": "a-name",
				"additional_taints": []interface{}{
					map[string]interface{}{"effect": "NoSchedule", "key": "taint-1", "value": "value-1"},
					map[string]interface{}{"effect": "NoExecute", "key": "taint-2", "value": "value-2"},
				},
			},
			successfulResult: BaseNodePool{
				Name: "a-name",
				AdditionalTaints: []K8sTaint{
					{Effect: "NoSchedule", Key: "taint-1", Value: "value-1"},
					{Effect: "NoExecute", Key: "taint-2", Value: "value-2"},
				},
			},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := expandClusterBaseNodePool(input.nodePool)

			if input.failureResult == "" {
				assert.Equal(t, output, input.successfulResult)
			} else {
				assert.EqualError(t, err, input.failureResult)
			}

		})
	}
}
