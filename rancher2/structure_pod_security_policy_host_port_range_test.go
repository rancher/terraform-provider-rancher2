package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPodSecurityPolicyHostPortRangesConf           []managementClient.HostPortRange
	testPodSecurityPolicyHostPortRangesInterface      []interface{}
	testEmptyPodSecurityPolicyHostPortRangesConf      []managementClient.HostPortRange
	testEmptyPodSecurityPolicyHostPortRangesInterface []interface{}
)

func init() {
	testPodSecurityPolicyHostPortRangesConf = []managementClient.HostPortRange{
		{
			Min: 1,
			Max: 3000,
		},
		{
			Min: 2,
			Max: 4000,
		},
	}
	testPodSecurityPolicyHostPortRangesInterface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
		map[string]interface{}{
			"min": 2,
			"max": 4000,
		},
	}
	testEmptyPodSecurityPolicyHostPortRangesInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyHostPortRanges(t *testing.T) {

	cases := []struct {
		Input          []managementClient.HostPortRange
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyHostPortRangesConf,
			testPodSecurityPolicyHostPortRangesInterface,
		},
		{
			testEmptyPodSecurityPolicyHostPortRangesConf,
			testEmptyPodSecurityPolicyHostPortRangesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyHostPortRanges(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyHostPortRanges(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.HostPortRange
	}{
		{
			testPodSecurityPolicyHostPortRangesInterface,
			testPodSecurityPolicyHostPortRangesConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyHostPortRanges(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
