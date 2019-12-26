package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyVolumesConf      []policyv1.FSType
	testPodSecurityPolicyVolumesInterface []interface{}
)

func init() {
	testPodSecurityPolicyVolumesConf = []policyv1.FSType{
		"hostPath",
		"emptyDir",
	}
	testPodSecurityPolicyVolumesInterface = []interface{}{
		"hostPath",
		"emptyDir",
	}
}

func TestFlattenPodSecurityPolicyVolumes(t *testing.T) {

	cases := []struct {
		Input          []policyv1.FSType
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyVolumesConf,
			testPodSecurityPolicyVolumesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyVolumes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyVolumes(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []policyv1.FSType
	}{
		{
			testPodSecurityPolicyVolumesInterface,
			testPodSecurityPolicyVolumesConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyVolumes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
