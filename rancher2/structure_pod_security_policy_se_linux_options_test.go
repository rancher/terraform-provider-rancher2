package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testPodSecurityPolicySELinuxOptionsConf           *managementClient.SELinuxOptions
	testPodSecurityPolicySELinuxOptionsInterface      []interface{}
	testNilPodSecurityPolicySELinuxOptionsConf        *managementClient.SELinuxOptions
	testEmptyPodSecurityPolicySELinuxOptionsInterface []interface{}
)

func init() {
	testPodSecurityPolicySELinuxOptionsConf = &managementClient.SELinuxOptions{
		User:  "user",
		Role:  "role",
		Type:  "type",
		Level: "level",
	}
	testPodSecurityPolicySELinuxOptionsInterface = []interface{}{
		map[string]interface{}{
			"user":  "user",
			"role":  "role",
			"type":  "type",
			"level": "level",
		},
	}
	testEmptyPodSecurityPolicySELinuxOptionsInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicySELinuxOptions(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SELinuxOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicySELinuxOptionsConf,
			testPodSecurityPolicySELinuxOptionsInterface,
		},
		{
			testNilPodSecurityPolicySELinuxOptionsConf,
			testEmptyPodSecurityPolicySELinuxOptionsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicySELinuxOptions(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandPodSecurityPolicySELinuxOptions(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.SELinuxOptions
	}{
		{
			testPodSecurityPolicySELinuxOptionsInterface,
			testPodSecurityPolicySELinuxOptionsConf,
		},
	}
	for _, tc := range cases {
		output := expandPodSecurityPolicySELinuxOptions(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
