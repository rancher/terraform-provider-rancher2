package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyAllowedHostPathsConf      []policyv1.AllowedHostPath
	testPodSecurityPolicyAllowedHostPathsInterface []interface{}
)

func init() {
	testPodSecurityPolicyAllowedHostPathsConf = []policyv1.AllowedHostPath{
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
			"read_only": true,
		},
		map[string]interface{}{
			"path_prefix": "/tmp",
			"read_only": false,
		},
	}
}

func TestFlattenPodSecurityPolicyAllowedHostPaths(t *testing.T) {

	cases := []struct {
		Input          []policyv1.AllowedHostPath
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyAllowedHostPathsConf,
			testPodSecurityPolicyAllowedHostPathsInterface,
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
		ExpectedOutput []policyv1.AllowedHostPath
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
