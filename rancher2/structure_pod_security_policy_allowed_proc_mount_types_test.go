package rancher2

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

var (
	testPodSecurityPolicyAllowedProcMountTypesConf      []v1.ProcMountType
	testPodSecurityPolicyAllowedProcMountTypesSlice     []string
)

func init() {
	testPodSecurityPolicyAllowedProcMountTypesConf = []v1.ProcMountType{
		"Default",
		"Unmasked",
	}
	testPodSecurityPolicyAllowedProcMountTypesSlice = []string{
		"Default",
		"Unmasked",
	}
}

func TestFlattenPodSecurityAllowedProcMountTypes(t *testing.T) {

	cases := []struct {
		Input          []v1.ProcMountType
		ExpectedOutput []string
	}{
		{
			testPodSecurityPolicyAllowedProcMountTypesConf,
			testPodSecurityPolicyAllowedProcMountTypesSlice,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyAllowedProcMountTypes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityAllowedProcMountTypes(t *testing.T) {

	cases := []struct {
		Input          []string
		ExpectedOutput []v1.ProcMountType
	}{
		{
			testPodSecurityPolicyAllowedProcMountTypesSlice,
			testPodSecurityPolicyAllowedProcMountTypesConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyAllowedProcMountTypes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
