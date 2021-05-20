package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testEnvVarConf      []managementClient.EnvVar
	testEnvVarInterface []interface{}
)

func init() {
	testEnvVarConf = []managementClient.EnvVar{
		{
			Name:  "name1",
			Value: "value1",
		},
		{
			Name:  "name2",
			Value: "value2",
		},
	}
	testEnvVarInterface = []interface{}{
		map[string]interface{}{
			"name":  "name1",
			"value": "value1",
		},
		map[string]interface{}{
			"name":  "name2",
			"value": "value2",
		},
	}
}

func TestFlattenEnvVars(t *testing.T) {

	cases := []struct {
		Input          []managementClient.EnvVar
		ExpectedOutput []interface{}
	}{
		{
			testEnvVarConf,
			testEnvVarInterface,
		},
	}

	for _, tc := range cases {
		output := flattenEnvVars(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandEnvVars(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.EnvVar
	}{
		{
			testEnvVarInterface,
			testEnvVarConf,
		},
	}

	for _, tc := range cases {
		output := expandEnvVars(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
