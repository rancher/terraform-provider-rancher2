package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPodSecurityPolicySupplementalGroupsConf           *managementClient.SupplementalGroupsStrategyOptions
	testPodSecurityPolicySupplementalGroupsInterface      []interface{}
	testNilPodSecurityPolicySupplementalGroupsConf        *managementClient.SupplementalGroupsStrategyOptions
	testEmptyPodSecurityPolicySupplementalGroupsInterface []interface{}
)

func init() {
	testPodSecurityPolicySupplementalGroupsConf = &managementClient.SupplementalGroupsStrategyOptions{
		Rule:   "RunAsAny",
		Ranges: testPodSecurityPolicyIDRangesConf,
	}
	testPodSecurityPolicySupplementalGroupsInterface = []interface{}{
		map[string]interface{}{
			"rule":  "RunAsAny",
			"range": testPodSecurityPolicyIDRangesInterface,
		},
	}
	testEmptyPodSecurityPolicySupplementalGroupsInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicySupplementalGroups(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SupplementalGroupsStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicySupplementalGroupsConf,
			testPodSecurityPolicySupplementalGroupsInterface,
		},
		{
			testNilPodSecurityPolicySupplementalGroupsConf,
			testEmptyPodSecurityPolicySupplementalGroupsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicySupplementalGroups(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicySupplementalGroups(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.SupplementalGroupsStrategyOptions
	}{
		{
			testPodSecurityPolicySupplementalGroupsInterface,
			testPodSecurityPolicySupplementalGroupsConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicySupplementalGroups(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
