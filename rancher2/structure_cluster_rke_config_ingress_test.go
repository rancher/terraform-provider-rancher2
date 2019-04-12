package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterRKEConfigIngressConf      *managementClient.IngressConfig
	testClusterRKEConfigIngressInterface []interface{}
)

func init() {
	testClusterRKEConfigIngressConf = &managementClient.IngressConfig{
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		NodeSelector: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Provider: "test",
	}
	testClusterRKEConfigIngressInterface = []interface{}{
		map[string]interface{}{
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"node_selector": map[string]interface{}{
				"node_one": "one",
				"node_two": "two",
			},
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"provider": "test",
		},
	}
}

func TestFlattenClusterRKEConfigIngress(t *testing.T) {

	cases := []struct {
		Input          *managementClient.IngressConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigIngressConf,
			testClusterRKEConfigIngressInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigIngress(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigIngress(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.IngressConfig
	}{
		{
			testClusterRKEConfigIngressInterface,
			testClusterRKEConfigIngressConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigIngress(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
