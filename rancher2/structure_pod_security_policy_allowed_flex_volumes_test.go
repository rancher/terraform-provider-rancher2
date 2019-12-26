package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testPodSecurityPolicyAllowedFlexVolumesConf      []managementClient.AllowedFlexVolume
	testPodSecurityPolicyAllowedFlexVolumesInterface []interface{}
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
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyAllowedFlexVolumes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
