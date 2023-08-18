package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
