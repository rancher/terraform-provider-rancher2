package rancher2

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

var (
	testPodSecurityPolicyCapabilitiesConf      []v1.Capability
	testPodSecurityPolicyCapabilitiesSlice     []interface{}
)

func init() {
	testPodSecurityPolicyCapabilitiesConf = []v1.Capability{
		"foo",
		"bar",
	}
	testPodSecurityPolicyCapabilitiesSlice = []interface{}{
		"foo",
		"bar",
	}
}

func TestFlattenPodSecurityPolicyCapabilities(t *testing.T) {

	cases := []struct {
		Input          []v1.Capability
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyCapabilitiesConf,
			testPodSecurityPolicyCapabilitiesSlice,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyCapabilities(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyCapabilities(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []v1.Capability
	}{
		{
			testPodSecurityPolicyCapabilitiesSlice,
			testPodSecurityPolicyCapabilitiesConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyCapabilities(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
