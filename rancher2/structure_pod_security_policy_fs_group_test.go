package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyIDRanges2Conf      []policyv1.IDRange
	testPodSecurityPolicyIDRanges2Interface []interface{}
	testPodSecurityPolicyFSGroupConf      policyv1.FSGroupStrategyOptions
	testPodSecurityPolicyFSGroupInterface []interface{}
)

func init() {
	testPodSecurityPolicyIDRanges2Conf = []policyv1.IDRange{
		{
			Min: int64(1),
			Max: int64(3000),
		},
		{
			Min: int64(0),
			Max: int64(5000),
		},
	}
	testPodSecurityPolicyIDRanges2Interface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
		map[string]interface{}{
			"min": 0,
			"max": 5000,
		},
	}
    testPodSecurityPolicyFSGroupConf = policyv1.FSGroupStrategyOptions{
		Rule: "RunAsAny",
		Ranges: testPodSecurityPolicyIDRanges2Conf,
	}
	testPodSecurityPolicyFSGroupInterface = []interface{}{
		map[string]interface{}{
			"rule": "RunAsAny",
			"ranges": testPodSecurityPolicyIDRanges2Interface,
		},
	}
}

func TestFlattenPodSecurityPolicyFSGroup(t *testing.T) {

	cases := []struct {
		Input          policyv1.FSGroupStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyFSGroupConf,
			testPodSecurityPolicyFSGroupInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyFSGroup(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyFSGroup(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput policyv1.FSGroupStrategyOptions
	}{
		{
			testPodSecurityPolicyFSGroupInterface,
			testPodSecurityPolicyFSGroupConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyFSGroup(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
