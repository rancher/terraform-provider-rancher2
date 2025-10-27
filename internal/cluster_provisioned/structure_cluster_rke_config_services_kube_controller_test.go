package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
