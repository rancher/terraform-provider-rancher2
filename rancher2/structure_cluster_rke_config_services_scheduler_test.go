package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterRKEConfigServicesSchedulerConf      *managementClient.SchedulerService
	testClusterRKEConfigServicesSchedulerInterface []interface{}
)

func init() {
	testClusterRKEConfigServicesSchedulerConf = &managementClient.SchedulerService{
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds: []string{"bind_one", "bind_two"},
		ExtraEnv:   []string{"env_one", "env_two"},
		Image:      "image",
	}
	testClusterRKEConfigServicesSchedulerInterface = []interface{}{
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

func TestFlattenClusterRKEConfigServicesScheduler(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SchedulerService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesSchedulerConf,
			testClusterRKEConfigServicesSchedulerInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesScheduler(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterRKEConfigServicesScheduler(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.SchedulerService
	}{
		{
			testClusterRKEConfigServicesSchedulerInterface,
			testClusterRKEConfigServicesSchedulerConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesScheduler(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
