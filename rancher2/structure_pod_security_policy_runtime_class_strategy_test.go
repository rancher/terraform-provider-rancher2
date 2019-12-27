package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testPodSecurityPolicyRuntimeClassStrategyConf      *managementClient.RuntimeClassStrategyOptions
	testPodSecurityPolicyRuntimeClassStrategyInterface []interface{}
)

func init() {
	testPodSecurityPolicyRuntimeClassStrategyConf = &managementClient.RuntimeClassStrategyOptions{
		AllowedRuntimeClassNames: []string{"foo", "bar"},
		DefaultRuntimeClassName:  "foo",
	}
	testPodSecurityPolicyRuntimeClassStrategyInterface = []interface{}{
		map[string]interface{}{
			"allowed_runtime_class_names": toArrayInterface([]string{"foo", "bar"}),
			"default_runtime_class_name":  "foo",
		},
	}
}

func TestFlattenPodSecurityPolicyRuntimeClassStrategy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RuntimeClassStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyRuntimeClassStrategyConf,
			testPodSecurityPolicyRuntimeClassStrategyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyRuntimeClassStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyRuntimeClassStrategy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.RuntimeClassStrategyOptions
	}{
		{
			testPodSecurityPolicyRuntimeClassStrategyInterface,
			testPodSecurityPolicyRuntimeClassStrategyConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyRuntimeClassStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
