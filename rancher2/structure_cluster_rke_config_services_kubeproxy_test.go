package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigServicesKubeproxyConf      *managementClient.KubeproxyService
	testClusterRKEConfigServicesKubeproxyInterface []interface{}
)

func init() {
	testClusterRKEConfigServicesKubeproxyConf = &managementClient.KubeproxyService{
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds: []string{"bind_one", "bind_two"},
		ExtraEnv:   []string{"env_one", "env_two"},
		Image:      "image",
	}
	testClusterRKEConfigServicesKubeproxyInterface = []interface{}{
		map[string]interface{}{
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds": []interface{}{"bind_one", "bind_two"},
			"extra_env":   []interface{}{"env_one", "env_two"},
			"image":       "image",
		},
	}
}

func TestFlattenClusterRKEConfigServicesKubeproxy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KubeproxyService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesKubeproxyConf,
			testClusterRKEConfigServicesKubeproxyInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesKubeproxy(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesKubeproxy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.KubeproxyService
	}{
		{
			testClusterRKEConfigServicesKubeproxyInterface,
			testClusterRKEConfigServicesKubeproxyConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesKubeproxy(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
