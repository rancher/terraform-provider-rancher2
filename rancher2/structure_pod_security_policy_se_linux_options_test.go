package rancher2

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

var (
	testPodSecurityPolicySELinuxOptionsConf      *v1.SELinuxOptions
	testPodSecurityPolicySELinuxOptionsInterface []interface{}
)

func init() {
	testPodSecurityPolicySELinuxOptionsConf = &v1.SELinuxOptions{
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
}

func TestFlattenPodSecurityPolicySELinuxOptions(t *testing.T) {

	cases := []struct {
		Input          *v1.SELinuxOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicySELinuxOptionsConf,
			testPodSecurityPolicySELinuxOptionsInterface,
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
		ExpectedOutput *v1.SELinuxOptions
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
