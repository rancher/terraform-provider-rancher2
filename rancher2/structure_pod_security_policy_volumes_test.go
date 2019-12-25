package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyVolumes          []policyv1.FSType
	testPodSecurityPolicyVolumesSlice     []string
)

func init() {
	testPodSecurityPolicyVolumes = []policyv1.FSType{
		"hostPath",
		"emptyDir",
	}
	testPodSecurityPolicyVolumesSlice = []string{
		"hostPath",
		"emptyDir",
	}
}

func TestFlattenPodSecurityPolicyVolumes(t *testing.T) {

	cases := []struct {
		Input          []policyv1.FSType
		ExpectedOutput []string
	}{
		{
			testPodSecurityPolicyVolumes,
			testPodSecurityPolicyVolumesSlice,
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
		Input          []string
		ExpectedOutput []policyv1.FSType
	}{
		{
			testPodSecurityPolicyVolumesSlice,
			testPodSecurityPolicyVolumes,
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
