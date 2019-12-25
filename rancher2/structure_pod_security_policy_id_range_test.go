package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyIDRange          []policyv1.IDRange
	testPodSecurityPolicyIDRangeInterface []interface{}
)

func init() {
	testPodSecurityPolicyIDRange = []policyv1.IDRange{
		{
			Min: int64(1),
			Max: int64(3000),
		},
		{
			Min: int64(0),
			Max: int64(5000),
		},
	}
	testPodSecurityPolicyIDRangeInterface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
		map[string]interface{}{
			"min": 0,
			"max": 5000,
		},
	}
}

func TestFlattenPodSecurityPolicyIPRanges(t *testing.T) {

	cases := []struct {
		Input          []policyv1.IDRange
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyIDRange,
			testPodSecurityPolicyIDRangeInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyIDRanges(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyIPRanges(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []policyv1.IDRange
	}{
		{
			testPodSecurityPolicyIDRangeInterface,
			testPodSecurityPolicyIDRange,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyIDRanges(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
