package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testPodSecurityPolicyAllowedHostPathsConf           []managementClient.AllowedHostPath
	testPodSecurityPolicyAllowedHostPathsInterface      []interface{}
	testEmptyPodSecurityPolicyAllowedHostPathsConf      []managementClient.AllowedHostPath
	testEmptyPodSecurityPolicyAllowedHostPathsInterface []interface{}
)

func init() {
	testPodSecurityPolicyAllowedHostPathsConf = []managementClient.AllowedHostPath{
		{
			PathPrefix: "/var/lib",
			ReadOnly:   true,
		},
		{
			PathPrefix: "/tmp",
		},
	}
	testPodSecurityPolicyAllowedHostPathsInterface = []interface{}{
		map[string]interface{}{
			"path_prefix": "/var/lib",
			"read_only":   true,
		},
		map[string]interface{}{
			"path_prefix": "/tmp",
			"read_only":   false,
		},
	}
	testEmptyPodSecurityPolicyAllowedHostPathsInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyAllowedHostPaths(t *testing.T) {

	cases := []struct {
		Input          []managementClient.AllowedHostPath
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyAllowedHostPathsConf,
			testPodSecurityPolicyAllowedHostPathsInterface,
		},
		{
			testEmptyPodSecurityPolicyAllowedHostPathsConf,
			testEmptyPodSecurityPolicyAllowedHostPathsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyAllowedHostPaths(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandPodSecurityPolicyAllowedHostPaths(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.AllowedHostPath
	}{
		{
			testPodSecurityPolicyAllowedHostPathsInterface,
			testPodSecurityPolicyAllowedHostPathsConf,
		},
	}
	for _, tc := range cases {
		output := expandPodSecurityPolicyAllowedHostPaths(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
