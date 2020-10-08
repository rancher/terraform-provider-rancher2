package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigServicesKubeControllerConf      *managementClient.KubeControllerService
	testClusterRKEConfigServicesKubeControllerInterface []interface{}
)

func init() {
	testClusterRKEConfigServicesKubeControllerConf = &managementClient.KubeControllerService{
		ClusterCIDR: "10.42.0.0/16",
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds:            []string{"bind_one", "bind_two"},
		ExtraEnv:              []string{"env_one", "env_two"},
		Image:                 "image",
		ServiceClusterIPRange: "10.43.0.0/16",
	}
	testClusterRKEConfigServicesKubeControllerInterface = []interface{}{
		map[string]interface{}{
			"cluster_cidr": "10.42.0.0/16",
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds":              []interface{}{"bind_one", "bind_two"},
			"extra_env":                []interface{}{"env_one", "env_two"},
			"image":                    "image",
			"service_cluster_ip_range": "10.43.0.0/16",
		},
	}
}

func TestFlattenClusterRKEConfigServicesKubeController(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KubeControllerService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesKubeControllerConf,
			testClusterRKEConfigServicesKubeControllerInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesKubeController(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesKubeController(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.KubeControllerService
	}{
		{
			testClusterRKEConfigServicesKubeControllerInterface,
			testClusterRKEConfigServicesKubeControllerConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesKubeController(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
