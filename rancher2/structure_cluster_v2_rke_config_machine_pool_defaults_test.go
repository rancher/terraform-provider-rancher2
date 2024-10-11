package rancher2

import (
	"testing"

	provisionv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterV2RKEConfigMachinePoolDefaultsConf      provisionv1.RKEMachinePoolDefaults
	testClusterV2RKEConfigMachinePoolDefaultsInterface []interface{}
)

func init() {
	testClusterV2RKEConfigMachinePoolDefaultsInterface = []any{
		map[string]any{
			"hostname_length_limit": 32,
		},
	}
	testClusterV2RKEConfigMachinePoolDefaultsConf = provisionv1.RKEMachinePoolDefaults{
		HostnameLengthLimit: 32,
	}
}

func TestFlattenClusterV2RKEConfigMachinePoolDefaults(t *testing.T) {

	cases := []struct {
		Input          provisionv1.RKEMachinePoolDefaults
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigMachinePoolDefaultsConf,
			testClusterV2RKEConfigMachinePoolDefaultsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigMachinePoolDefaults(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterV2RKEConfigMachinePoolDefaults(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput provisionv1.RKEMachinePoolDefaults
	}{
		{
			testClusterV2RKEConfigMachinePoolDefaultsInterface,
			testClusterV2RKEConfigMachinePoolDefaultsConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigMachinePoolDefaults(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
