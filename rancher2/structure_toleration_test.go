package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testTolerationsConf      []managementClient.Toleration
	testTolerationsInterface []interface{}
)

func init() {
	seconds := int64(10)
	testTolerationsConf = []managementClient.Toleration{
		{
			Key:               "key",
			Value:             "value",
			Effect:            "recipient",
			Operator:          "operator",
			TolerationSeconds: &seconds,
		},
	}
	testTolerationsInterface = []interface{}{
		map[string]interface{}{
			"key":      "key",
			"value":    "value",
			"effect":   "recipient",
			"operator": "operator",
			"seconds":  10,
		},
	}
}

func TestFlattenTolerations(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Toleration
		ExpectedOutput []interface{}
	}{
		{
			testTolerationsConf,
			testTolerationsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenTolerations(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandTolerations(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Toleration
	}{
		{
			testTolerationsInterface,
			testTolerationsConf,
		},
	}

	for _, tc := range cases {
		output := expandTolerations(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
