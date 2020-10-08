package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testTaintsConf      []managementClient.Taint
	testTaintsInterface []interface{}
)

func init() {
	testTaintsConf = []managementClient.Taint{
		{
			Key:       "key",
			Value:     "value",
			Effect:    "recipient",
			TimeAdded: "time_added",
		},
	}
	testTaintsInterface = []interface{}{
		map[string]interface{}{
			"key":        "key",
			"value":      "value",
			"effect":     "recipient",
			"time_added": "time_added",
		},
	}
}

func TestFlattenTaints(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Taint
		ExpectedOutput []interface{}
	}{
		{
			testTaintsConf,
			testTaintsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenTaints(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandTaints(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Taint
	}{
		{
			testTaintsInterface,
			testTaintsConf,
		},
	}

	for _, tc := range cases {
		output := expandTaints(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
