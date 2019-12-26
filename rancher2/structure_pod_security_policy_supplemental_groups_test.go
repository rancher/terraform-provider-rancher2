package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testPodSecurityPolicyIDRanges3Conf               []managementClient.IDRange
	testPodSecurityPolicyIDRanges3Interface          []interface{}
	testPodSecurityPolicySupplementalGroupsConf      *managementClient.SupplementalGroupsStrategyOptions
	testPodSecurityPolicySupplementalGroupsInterface []interface{}
)

func init() {
	testPodSecurityPolicyIDRanges3Conf = []managementClient.IDRange{
		{
			Min: int64(1),
			Max: int64(3000),
		},
		{
			Min: int64(0),
			Max: int64(5000),
		},
	}
	testPodSecurityPolicyIDRanges3Interface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
		map[string]interface{}{
			"min": 0,
			"max": 5000,
		},
	}
	testPodSecurityPolicySupplementalGroupsConf = &managementClient.SupplementalGroupsStrategyOptions{
		Rule: "RunAsAny",
		Ranges: testPodSecurityPolicyIDRanges3Conf,
	}
	testPodSecurityPolicySupplementalGroupsInterface = []interface{}{
		map[string]interface{}{
			"rule": "RunAsAny",
			"ranges": testPodSecurityPolicyIDRanges3Interface,
		},
	}
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
