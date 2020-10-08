package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPodSecurityPolicyFSGroupConf              *managementClient.FSGroupStrategyOptions
	testPodSecurityPolicyFSGroupInterface         []interface{}
	testPodSecurityPolicyFSGroupIDRangesConf      []managementClient.IDRange
	testPodSecurityPolicyFSGroupIDRangesInterface []interface{}
	testNilPodSecurityPolicyFSGroupConf           *managementClient.FSGroupStrategyOptions
	testEmptyPodSecurityPolicyFSGroupInterface    []interface{}
)

func init() {
	testPodSecurityPolicyFSGroupIDRangesConf = []managementClient.IDRange{
		{
			Min: int64(1),
			Max: int64(3000),
		},
		{
			Min: int64(0),
			Max: int64(5000),
		},
	}
	testPodSecurityPolicyFSGroupIDRangesInterface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
		map[string]interface{}{
			"min": 0,
			"max": 5000,
		},
	}
	testPodSecurityPolicyFSGroupConf = &managementClient.FSGroupStrategyOptions{
		Rule:   "RunAsAny",
		Ranges: testPodSecurityPolicyFSGroupIDRangesConf,
	}
	testPodSecurityPolicyFSGroupInterface = []interface{}{
		map[string]interface{}{
			"rule":  "RunAsAny",
			"range": testPodSecurityPolicyFSGroupIDRangesInterface,
		},
	}
	testEmptyPodSecurityPolicyFSGroupInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyFSGroup(t *testing.T) {

	cases := []struct {
		Input          *managementClient.FSGroupStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyFSGroupConf,
			testPodSecurityPolicyFSGroupInterface,
		},
		{
			testNilPodSecurityPolicyFSGroupConf,
			testEmptyPodSecurityPolicyFSGroupInterface,
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
		ExpectedOutput *managementClient.FSGroupStrategyOptions
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
