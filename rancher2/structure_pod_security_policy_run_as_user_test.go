package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPodSecurityPolicyRunAsUserConf           *managementClient.RunAsUserStrategyOptions
	testPodSecurityPolicyRunAsUserInterface      []interface{}
	testNilPodSecurityPolicyRunAsUserConf        *managementClient.RunAsUserStrategyOptions
	testEmptyPodSecurityPolicyRunAsUserInterface []interface{}
)

func init() {
	testPodSecurityPolicyRunAsUserConf = &managementClient.RunAsUserStrategyOptions{
		Rule:   "RunAsAny",
		Ranges: testPodSecurityPolicyIDRangesConf,
	}
	testPodSecurityPolicyRunAsUserInterface = []interface{}{
		map[string]interface{}{
			"rule":  "RunAsAny",
			"range": testPodSecurityPolicyIDRangesInterface,
		},
	}
	testEmptyPodSecurityPolicyRunAsUserInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyRunAsUser(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RunAsUserStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyRunAsUserConf,
			testPodSecurityPolicyRunAsUserInterface,
		},
		{
			testNilPodSecurityPolicyRunAsUserConf,
			testEmptyPodSecurityPolicyRunAsUserInterface,
		},
	}

	for _, tc := range cases {
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
		ExpectedOutput *managementClient.RunAsUserStrategyOptions
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
