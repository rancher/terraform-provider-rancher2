package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFlattenMachineConfigV2Harvester_CPU(t *testing.T) {
	testcases := []struct {
		Input       *MachineConfigV2Harvester
		Expectation []any
	}{
		{
			// All default values
			Input:       &MachineConfigV2Harvester{},
			Expectation: []any{map[string]any{}},
		},
		{
			Input: &MachineConfigV2Harvester{
				machineConfigV2Harvester: machineConfigV2Harvester{
					CPUCount:              "3",
					CPUPinning:            true,
					IsolateEmulatorThread: true,
				},
			},
			Expectation: []any{map[string]any{
				"cpu_count":               "3",
				"cpu_pinning":             true,
				"isolate_emulator_thread": true,
			}},
		},
	}

	for _, tc := range testcases {
		output := flattenMachineConfigV2Harvester(tc.Input)
		assert.Equal(t, tc.Expectation, output, "unexpected output from flattenMachineConfigV2Harvester")
	}
}

func TestExpandMachineConfigV2Harvester_CPU(t *testing.T) {
	testcases := []struct {
		Input       []any
		Expectation *MachineConfigV2Harvester
	}{
		{
			// All default values
			Input: []any{map[string]any{}},
			Expectation: &MachineConfigV2Harvester{
				machineConfigV2Harvester: machineConfigV2Harvester{
					TypeMeta: metav1.TypeMeta{
						Kind:       machineConfigV2HarvesterKind,
						APIVersion: machineConfigV2HarvesterAPIVersion,
					},
				},
			},
		},
		{
			Input: []any{map[string]any{
				"cpu_count":               "3",
				"cpu_pinning":             true,
				"isolate_emulator_thread": true,
			}},
			Expectation: &MachineConfigV2Harvester{
				machineConfigV2Harvester: machineConfigV2Harvester{
					TypeMeta: metav1.TypeMeta{
						Kind:       machineConfigV2HarvesterKind,
						APIVersion: machineConfigV2HarvesterAPIVersion,
					},
					CPUCount:              "3",
					CPUPinning:            true,
					IsolateEmulatorThread: true,
				},
			},
		},
	}

	for _, tc := range testcases {
		output := expandMachineConfigV2Harvester(tc.Input, &MachineConfigV2{})
		assert.Equal(t, tc.Expectation, output, "unexpected output from expandMachineConfigV2Harvester")
	}
}
