package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyIDRangesConf      []policyv1.IDRange
	testPodSecurityPolicyIDRangesInterface []interface{}
)

func init() {
	testPodSecurityPolicyIDRangesConf = []policyv1.IDRange{
		{
			Min: int64(1),
			Max: int64(3000),
		},
		{
			Min: int64(0),
			Max: int64(5000),
		},
	}
	testPodSecurityPolicyIDRangesInterface = []interface{}{
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

func TestFlattenPodSecurityPolicyIDRanges(t *testing.T) {

	cases := []struct {
		Input          []policyv1.IDRange
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyIDRangesConf,
			testPodSecurityPolicyIDRangesInterface,
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

func TestExpandPodSecurityPolicyIDRanges(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []policyv1.IDRange
	}{
		{
			testPodSecurityPolicyIDRangesInterface,
			testPodSecurityPolicyIDRangesConf,
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
