package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyRunAsGroupConf      *policyv1.RunAsGroupStrategyOptions
	testPodSecurityPolicyRunAsGroupInterface []interface{}
)

func init() {
	testPodSecurityPolicyRunAsGroupConf = &policyv1.RunAsGroupStrategyOptions{
		Rule: "RunAsAny",
		Ranges: testPodSecurityPolicyIDRangesConf,
	}
	testPodSecurityPolicyRunAsGroupInterface = []interface{}{
		map[string]interface{}{
			"rule": "RunAsAny",
			"ranges": testPodSecurityPolicyIDRangesInterface,
		},
	}
}

func TestFlattenPodSecurityPolicyRunAsGroup(t *testing.T) {

	cases := []struct {
		Input          *policyv1.RunAsGroupStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyRunAsGroupConf,
			testPodSecurityPolicyRunAsGroupInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyRunAsGroup(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyRunAsGroup(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *policyv1.RunAsGroupStrategyOptions
	}{
		{
			testPodSecurityPolicyRunAsGroupInterface,
			testPodSecurityPolicyRunAsGroupConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyRunAsGroup(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
