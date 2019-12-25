package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyHostPortRange          *policyv1.HostPortRange
	testPodSecurityPolicyHostPortRangeInterface []interface{}
)

func init() {
	testPodSecurityPolicyHostPortRange = &policyv1.HostPortRange{
		Min: 1,
		Max: 3000,
	}
	testPodSecurityPolicyHostPortRangeInterface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
	}
}

func TestFlattenPodSecurityPolicyHostPortRange(t *testing.T) {

	cases := []struct {
		Input          *policyv1.HostPortRange
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyHostPortRange,
			testPodSecurityPolicyHostPortRangeInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyHostPortRange(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyHostPortRange(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *policyv1.HostPortRange
	}{
		{
			testPodSecurityPolicyHostPortRangeInterface,
			testPodSecurityPolicyHostPortRange,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyHostPortRange(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
