package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
