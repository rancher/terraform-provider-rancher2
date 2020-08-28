package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPodSecurityPolicyIDRangesConf           []managementClient.IDRange
	testPodSecurityPolicyIDRangesInterface      []interface{}
	testEmptyPodSecurityPolicyIDRangesConf      []managementClient.IDRange
	testEmptyPodSecurityPolicyIDRangesInterface []interface{}
)

func init() {
	testPodSecurityPolicyIDRangesConf = []managementClient.IDRange{
		{
			Min: int64(1),
			Max: int64(3000),
		},
		{
			Min: int64(0),
			Max: int64(5000),
		},
	}
	testPodSecurityPolicyIDRangesInterface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
		map[string]interface{}{
			"min": 0,
			"max": 5000,
		},
	}
	testEmptyPodSecurityPolicyIDRangesInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyIDRanges(t *testing.T) {

	cases := []struct {
		Input          []managementClient.IDRange
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyIDRangesConf,
			testPodSecurityPolicyIDRangesInterface,
		},
		{
			testEmptyPodSecurityPolicyIDRangesConf,
			testEmptyPodSecurityPolicyIDRangesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyIDRanges(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyIDRanges(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.IDRange
	}{
		{
			testPodSecurityPolicyIDRangesInterface,
			testPodSecurityPolicyIDRangesConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyIDRanges(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
