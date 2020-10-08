package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
