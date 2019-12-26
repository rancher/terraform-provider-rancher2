package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicySELinuxStrategyConf      policyv1.SELinuxStrategyOptions
	testPodSecurityPolicySELinuxStrategyInterface []interface{}
)

func init() {
	testPodSecurityPolicySELinuxStrategyConf = policyv1.SELinuxStrategyOptions{
		Rule: "RunAsAny",
		SELinuxOptions: testPodSecurityPolicySELinuxOptionsConf,
	}
	testPodSecurityPolicySELinuxStrategyInterface = []interface{}{
		map[string]interface{}{
			"rule": "RunAsAny",
			"se_linux_options": testPodSecurityPolicySELinuxOptionsInterface,
		},
	}
}

func TestFlattenPodSecurityPolicySELinuxStrategy(t *testing.T) {

	cases := []struct {
		Input          policyv1.SELinuxStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicySELinuxStrategyConf,
			testPodSecurityPolicySELinuxStrategyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicySELinuxStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicySELinuxStrategy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput policyv1.SELinuxStrategyOptions
	}{
		{
			testPodSecurityPolicySELinuxStrategyInterface,
			testPodSecurityPolicySELinuxStrategyConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicySELinuxStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
