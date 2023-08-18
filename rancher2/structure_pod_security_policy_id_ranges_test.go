package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
