package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
    testPodSecurityPolicySupplementalGroupsIDRanges  []policyv1.IDRange
	testPodSecurityPolicySupplementalGroups          *policyv1.SupplementalGroupsStrategyOptions
	testPodSecurityPolicySupplementalGroupsInterface []interface{}
)

func init() {
    testPodSecurityPolicySupplementalGroupsIDRanges = []policyv1.IDRange{
        {
        	Min: int64(0),
			Max: int64(5000),
        },
        {
        	Min: int64(1),
			Max: int64(4000),
        },
    }
	testPodSecurityPolicySupplementalGroups = &policyv1.SupplementalGroupsStrategyOptions{
		Rule: "RunAsAny",
		Ranges: testPodSecurityPolicySupplementalGroupsIDRanges,
	}
	testPodSecurityPolicySupplementalGroupsInterface = []interface{}{
		map[string]interface{}{
			"rule": "RunAsAny",
			"ranges": []interface{}{
				map[string]interface{}{
					"min": 0,
					"max": 5000,
				},
				map[string]interface{}{
					"min": 1,
					"max": 4000,
				},
			},
		},
	}
}

func TestFlattenPodSecurityPolicySupplementalGroups(t *testing.T) {

	cases := []struct {
		Input          *policyv1.SupplementalGroupsStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicySupplementalGroups,
			testPodSecurityPolicySupplementalGroupsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicySupplementalGroup(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicySupplementalGroups(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *policyv1.SupplementalGroupsStrategyOptions
	}{
		{
			testPodSecurityPolicySupplementalGroupsInterface,
			testPodSecurityPolicySupplementalGroups,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicySupplementalGroup(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
