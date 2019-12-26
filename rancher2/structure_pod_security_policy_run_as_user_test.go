package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyRunAsUserConf      policyv1.RunAsUserStrategyOptions
	testPodSecurityPolicyRunAsUserInterface []interface{}
)

func init() {
	testPodSecurityPolicyRunAsUserConf = policyv1.RunAsUserStrategyOptions{
		Rule: "RunAsAny",
		Ranges: testPodSecurityPolicyIDRangesConf,
	}
	testPodSecurityPolicyRunAsUserInterface = []interface{}{
		map[string]interface{}{
			"rule": "RunAsAny",
			"ranges": testPodSecurityPolicyIDRangesInterface,
		},
	}
}

func TestFlattenPodSecurityPolicyRunAsUser(t *testing.T) {

	cases := []struct {
		Input          policyv1.RunAsUserStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyRunAsUserConf,
			testPodSecurityPolicyRunAsUserInterface,
		},
	}
	
	for _, tc := range cases {
		t.Logf("Expected: %#v\nGiven:    %#v", tc.Input, tc.ExpectedOutput)
		output := flattenPodSecurityPolicyRunAsUser(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyRunAsUser(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput policyv1.RunAsUserStrategyOptions
	}{
		{
			testPodSecurityPolicyRunAsUserInterface,
			testPodSecurityPolicyRunAsUserConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyRunAsUser(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
