package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
