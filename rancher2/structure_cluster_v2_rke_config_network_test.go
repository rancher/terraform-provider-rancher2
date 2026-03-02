package rancher2

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

func TestFlattenClusterV2RKEConfigNetwork(t *testing.T) {
	cases := []struct {
		Name           string
		Input          *rkev1.Networking
		ExpectedOutput []interface{}
	}{
		{
			Name: "ipv4 stack",
			Input: &rkev1.Networking{
				StackPreference: rkev1.NetworkingStackPreference("ipv4"),
			},
			ExpectedOutput: []interface{}{
				map[string]interface{}{
					"stack_preference": "ipv4",
				},
			},
		},
		{
			Name: "ipv6 stack",
			Input: &rkev1.Networking{
				StackPreference: rkev1.NetworkingStackPreference("ipv6"),
			},
			ExpectedOutput: []interface{}{
				map[string]interface{}{
					"stack_preference": "ipv6",
				},
			},
		},
		{
			Name:           "nil input",
			Input:          nil,
			ExpectedOutput: nil,
		},
		{
			Name: "empty stack_preference",
			Input: &rkev1.Networking{
				StackPreference: "",
			},
			ExpectedOutput: []interface{}{
				map[string]interface{}{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			output := flattenClusterV2Networking(tc.Input)
			assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
		})
	}
}

func TestExpandClusterV2RKEConfigNetwork(t *testing.T) {
	cases := []struct {
		Name           string
		Input          []interface{}
		ExpectedOutput *rkev1.Networking
	}{
		{
			Name: "ipv4 stack",
			Input: []interface{}{
				map[string]interface{}{
					"stack_preference": "ipv4",
				},
			},
			ExpectedOutput: &rkev1.Networking{
				StackPreference: rkev1.NetworkingStackPreference("ipv4"),
			},
		},
		{
			Name: "ipv6 stack",
			Input: []interface{}{
				map[string]interface{}{
					"stack_preference": "ipv6",
				},
			},
			ExpectedOutput: &rkev1.Networking{
				StackPreference: rkev1.NetworkingStackPreference("ipv6"),
			},
		},
		{
			Name: "dualstack",
			Input: []interface{}{
				map[string]interface{}{
					"stack_preference": "dualstack",
				},
			},
			ExpectedOutput: &rkev1.Networking{
				StackPreference: rkev1.NetworkingStackPreference("dualstack"),
			},
		},
		{
			Name:           "nil input",
			Input:          nil,
			ExpectedOutput: nil,
		},
		{
			Name:           "empty slice",
			Input:          []interface{}{},
			ExpectedOutput: nil,
		},
		{
			Name: "empty stack_preference",
			Input: []interface{}{
				map[string]interface{}{
					"stack_preference": "",
				},
			},
			ExpectedOutput: &rkev1.Networking{},
		},
		{
			Name: "invalid type for stack_preference",
			Input: []interface{}{
				map[string]interface{}{
					"stack_preference": 123,
				},
			},
			ExpectedOutput: &rkev1.Networking{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			output := expandClusterV2Networking(tc.Input)
			assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
		})
	}
}

func TestRoundTripClusterV2RKEConfigNetwork(t *testing.T) {
	input := &rkev1.Networking{
		StackPreference: rkev1.NetworkingStackPreference("ipv6"),
	}

	// Flatten then Expand
	flat := flattenClusterV2Networking(input)
	expanded := expandClusterV2Networking(flat)
	assert.Equal(t, input, expanded, "Round-trip flatten to expand failed.")

	// Expand then Flatten
	expandedAgain := expandClusterV2Networking(flat)
	flatAgain := flattenClusterV2Networking(expandedAgain)
	assert.Equal(t, flat, flatAgain, "Round-trip expand to flatten failed.")
}
