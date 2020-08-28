package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPodSecurityPolicySELinuxStrategyConf           *managementClient.SELinuxStrategyOptions
	testPodSecurityPolicySELinuxStrategyInterface      []interface{}
	testNilPodSecurityPolicySELinuxStrategyConf        *managementClient.SELinuxStrategyOptions
	testEmptyPodSecurityPolicySELinuxStrategyInterface []interface{}
)

func init() {
	testPodSecurityPolicySELinuxStrategyConf = &managementClient.SELinuxStrategyOptions{
		Rule:           "RunAsAny",
		SELinuxOptions: testPodSecurityPolicySELinuxOptionsConf,
	}
	testPodSecurityPolicySELinuxStrategyInterface = []interface{}{
		map[string]interface{}{
			"rule":            "RunAsAny",
			"se_linux_option": testPodSecurityPolicySELinuxOptionsInterface,
		},
	}
	testEmptyPodSecurityPolicySELinuxStrategyInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicySELinuxStrategy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SELinuxStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicySELinuxStrategyConf,
			testPodSecurityPolicySELinuxStrategyInterface,
		},
		{
			testNilPodSecurityPolicySELinuxStrategyConf,
			testEmptyPodSecurityPolicySELinuxStrategyInterface,
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
		ExpectedOutput *managementClient.SELinuxStrategyOptions
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
