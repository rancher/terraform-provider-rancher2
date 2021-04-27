package rancher2

import (
	"reflect"
	"testing"
)

var (
	testCloudCredentialGoogleConf      *googleCredentialConfig
	testCloudCredentialGoogleInterface []interface{}
)

func init() {
	testCloudCredentialGoogleConf = &googleCredentialConfig{
		AuthEncodedJSON: "{\"auth_encoded_json\": true}",
	}
	testCloudCredentialGoogleInterface = []interface{}{
		map[string]interface{}{
			"auth_encoded_json": "{\"auth_encoded_json\": true}",
		},
	}
}

func TestFlattenCloudCredentialGoogle(t *testing.T) {

	cases := []struct {
		Input          *googleCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialGoogleConf,
			testCloudCredentialGoogleInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialGoogle(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialGoogle(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *googleCredentialConfig
	}{
		{
			testCloudCredentialGoogleInterface,
			testCloudCredentialGoogleConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialGoogle(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
