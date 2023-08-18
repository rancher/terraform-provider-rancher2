package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
