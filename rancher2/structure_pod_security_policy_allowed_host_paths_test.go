package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
