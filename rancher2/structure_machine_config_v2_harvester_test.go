package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenMachineConfigV2HarvesterCPU(t *testing.T) {
	testcases := []struct {
		Input       *MachineConfigV2HarvesterCPU
		Expectation []any
	}{
		{
			// All default values
			Input: &MachineConfigV2HarvesterCPU{},
			Expectation: []any{map[string]interface{}{
				"count":                 2,
				"pinning":               false,
				"isolateEmulatorThread": false,
			}},
		},
		{
			Input: &MachineConfigV2HarvesterCPU{
				Count:                 3,
				Pinning:               true,
				IsolateEmulatorThread: true,
			},
			Expectation: []any{map[string]interface{}{
				"count":                 3,
				"pinning":               true,
				"isolateEmulatorThread": true,
			}},
		},
	}

	for _, tc := range testcases {
		output := flattenMachineConfigV2HarvesterCPU(tc.Input)
		assert.Equal(t, tc.Expectation, output, "unexpected output from flattenMachineConfigV2HarvesterCPU")
	}
}

func TestExpandMachineConfigV2HarvesterCPU(t *testing.T) {
	testcases := []struct {
		Input       []any
		Expectation *MachineConfigV2HarvesterCPU
	}{
		{
			// All default values
			Input: []any{map[string]interface{}{
				"count":                 2,
				"pinning":               false,
				"isolateEmulatorThread": false,
			}},
			Expectation: &MachineConfigV2HarvesterCPU{
				Count: 2,
			},
		},
		{
			Input: []any{map[string]interface{}{
				"count":                 3,
				"pinning":               true,
				"isolateEmulatorThread": true,
			}},
			Expectation: &MachineConfigV2HarvesterCPU{
				Count:                 3,
				Pinning:               true,
				IsolateEmulatorThread: true,
			},
		},
	}

	for _, tc := range testcases {
		output := expandMachineConfigV2HarvesterCPU(tc.Input)
		assert.Equal(t, tc.Expectation, output, "unexpected output from flattenMachineConfigV2HarvesterCPU")
	}
}
