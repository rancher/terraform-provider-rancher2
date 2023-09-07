package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
