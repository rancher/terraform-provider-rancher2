package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyAllowedFlexVolumesConf      []policyv1.AllowedFlexVolume
	testPodSecurityPolicyAllowedFlexVolumesInterface []interface{}
)

func init() {
	testPodSecurityPolicyAllowedFlexVolumesConf = []policyv1.AllowedFlexVolume{
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
		Input          []policyv1.AllowedFlexVolume
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
		ExpectedOutput []policyv1.AllowedFlexVolume
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
