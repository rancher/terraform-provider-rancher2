package rancher2

import (
	"reflect"
	"testing"
)

var (
	testCloudCredentialOpenstackConf      *openstackCredentialConfig
	testCloudCredentialOpenstackInterface []interface{}
)

func init() {
	testCloudCredentialOpenstackConf = &openstackCredentialConfig{
		Password: "password",
	}
	testCloudCredentialOpenstackInterface = []interface{}{
		map[string]interface{}{
			"password": "password",
		},
	}
}

func TestFlattenCloudCredentialOpenstack(t *testing.T) {

	cases := []struct {
		Input          *openstackCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialOpenstackConf,
			testCloudCredentialOpenstackInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialOpenstack(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialOpenstack(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *openstackCredentialConfig
	}{
		{
			testCloudCredentialOpenstackInterface,
			testCloudCredentialOpenstackConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialOpenstack(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
