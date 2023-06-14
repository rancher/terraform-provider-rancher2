package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testResourceRequirementsConf      *managementClient.ResourceRequirements
	testResourceRequirementsInterface []interface{}
)

func init() {
	testVal := "500"
	testResourceRequirementsConf = &managementClient.ResourceRequirements{
		Limits: map[string]string{
			"cpu":    testVal,
			"memory": testVal,
		},
		Requests: map[string]string{
			"cpu":    testVal,
			"memory": testVal,
		},
	}
	testResourceRequirementsInterface = []interface{}{
		map[string]interface{}{
			"cpu_limit":      testVal,
			"memory_limit":   testVal,
			"cpu_request":    testVal,
			"memory_request": testVal,
		},
	}
}

func TestFlattenResourceRequirements(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ResourceRequirements
		ExpectedOutput []interface{}
	}{
		{
			testResourceRequirementsConf,
			testResourceRequirementsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenResourceRequirements(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandResourceRequirements(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ResourceRequirements
	}{
		{
			testResourceRequirementsInterface,
			testResourceRequirementsConf,
		},
	}

	for _, tc := range cases {
		output := expandResourceRequirements(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
