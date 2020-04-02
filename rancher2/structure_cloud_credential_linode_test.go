package rancher2

import (
	"reflect"
	"testing"
)

var (
	testCloudCredentialLinodeConf      *linodeCredentialConfig
	testCloudCredentialLinodeInterface []interface{}
)

func init() {
	testCloudCredentialLinodeConf = &linodeCredentialConfig{
		Token: "token",
	}
	testCloudCredentialLinodeInterface = []interface{}{
		map[string]interface{}{
			"token": "token",
		},
	}
}

func TestFlattenCloudCredentialLinode(t *testing.T) {

	cases := []struct {
		Input          *linodeCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialLinodeConf,
			testCloudCredentialLinodeInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialLinode(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialLinode(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *linodeCredentialConfig
	}{
		{
			testCloudCredentialLinodeInterface,
			testCloudCredentialLinodeConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialLinode(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
