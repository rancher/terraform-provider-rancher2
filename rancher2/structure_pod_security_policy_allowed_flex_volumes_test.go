package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testPodSecurityPolicyAllowedFlexVolumesConf           []managementClient.AllowedFlexVolume
	testPodSecurityPolicyAllowedFlexVolumesInterface      []interface{}
	testEmptyPodSecurityPolicyAllowedFlexVolumesConf      []managementClient.AllowedFlexVolume
	testEmptyPodSecurityPolicyAllowedFlexVolumesInterface []interface{}
)

func init() {
	testPodSecurityPolicyAllowedFlexVolumesConf = []managementClient.AllowedFlexVolume{
		{
			Driver: "foo",
		},
		{
			Driver: "bar",
		},
	}
	testPodSecurityPolicyAllowedFlexVolumesInterface = []interface{}{
		map[string]interface{}{
			"driver": "foo",
		},
		map[string]interface{}{
			"driver": "bar",
		},
	}
	testEmptyPodSecurityPolicyAllowedFlexVolumesInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyAllowedFlexVolumes(t *testing.T) {

	cases := []struct {
		Input          []managementClient.AllowedFlexVolume
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyAllowedFlexVolumesConf,
			testPodSecurityPolicyAllowedFlexVolumesInterface,
		},
		{
			testEmptyPodSecurityPolicyAllowedFlexVolumesConf,
			testEmptyPodSecurityPolicyAllowedFlexVolumesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyAllowedFlexVolumes(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandPodSecurityPolicyAllowedFlexVolumes(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.AllowedFlexVolume
	}{
		{
			testPodSecurityPolicyAllowedFlexVolumesInterface,
			testPodSecurityPolicyAllowedFlexVolumesConf,
		},
	}
	for _, tc := range cases {
		output := expandPodSecurityPolicyAllowedFlexVolumes(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
