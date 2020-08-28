package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testTargetsConf      []managementClient.Target
	testTargetsInterface []interface{}
)

func init() {
	testTargetsConf = []managementClient.Target{
		{
			ProjectID:   "project_id",
			AppID:       "app_id",
			Healthstate: "health_state",
			State:       "state",
		},
	}
	testTargetsInterface = []interface{}{
		map[string]interface{}{
			"project_id":   "project_id",
			"app_id":       "app_id",
			"health_state": "health_state",
			"state":        "state",
		},
	}
}

func TestFlattenTargets(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Target
		ExpectedOutput []interface{}
	}{
		{
			testTargetsConf,
			testTargetsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenTargets(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandTargets(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Target
	}{
		{
			testTargetsInterface,
			testTargetsConf,
		},
	}

	for _, tc := range cases {
		output := expandTargets(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
