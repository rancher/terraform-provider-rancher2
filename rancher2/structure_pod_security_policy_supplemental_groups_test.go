package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
