package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPodSecurityPolicyRunAsGroupConf           *managementClient.RunAsGroupStrategyOptions
	testPodSecurityPolicyRunAsGroupInterface      []interface{}
	testNilPodSecurityPolicyRunAsGroupConf        *managementClient.RunAsGroupStrategyOptions
	testEmptyPodSecurityPolicyRunAsGroupInterface []interface{}
)

func init() {
	testPodSecurityPolicyRunAsGroupConf = &managementClient.RunAsGroupStrategyOptions{
		Rule:   "RunAsAny",
		Ranges: testPodSecurityPolicyIDRangesConf,
	}
	testPodSecurityPolicyRunAsGroupInterface = []interface{}{
		map[string]interface{}{
			"rule":  "RunAsAny",
			"range": testPodSecurityPolicyIDRangesInterface,
		},
	}
	testEmptyPodSecurityPolicyRunAsGroupInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyRunAsGroup(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RunAsGroupStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyRunAsGroupConf,
			testPodSecurityPolicyRunAsGroupInterface,
		},
		{
			testNilPodSecurityPolicyRunAsGroupConf,
			testEmptyPodSecurityPolicyRunAsGroupInterface,
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
		ExpectedOutput *managementClient.RunAsGroupStrategyOptions
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
