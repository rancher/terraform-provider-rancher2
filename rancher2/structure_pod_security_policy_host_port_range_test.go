package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyHostPortRangesConf      []policyv1.HostPortRange
	testPodSecurityPolicyHostPortRangesInterface []interface{}
)

func init() {
	testPodSecurityPolicyHostPortRangesConf = []policyv1.HostPortRange{
		{
			Min: 1,
			Max: 3000,
		},
		{
			Min: 2,
			Max: 4000,
		},
	}
	testPodSecurityPolicyHostPortRangesInterface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
		map[string]interface{}{
			"min": 2,
			"max": 4000,
		},
	}
}

func TestFlattenPodSecurityPolicyHostPortRanges(t *testing.T) {

	cases := []struct {
		Input          []policyv1.HostPortRange
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyHostPortRangesConf,
			testPodSecurityPolicyHostPortRangesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyHostPortRanges(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyHostPortRanges(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []policyv1.HostPortRange
	}{
		{
			testPodSecurityPolicyHostPortRangesInterface,
			testPodSecurityPolicyHostPortRangesConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyHostPortRanges(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
