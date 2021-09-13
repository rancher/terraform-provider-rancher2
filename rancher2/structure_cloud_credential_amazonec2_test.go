package rancher2

import (
	"reflect"
	"testing"
)

var (
	testCloudCredentialAmazonec2Conf      *amazonec2CredentialConfig
	testCloudCredentialAmazonec2Interface []interface{}
)

func init() {
	testCloudCredentialAmazonec2Conf = &amazonec2CredentialConfig{
		AccessKey:     "access_key",
		SecretKey:     "secret_key",
		DefaultRegion: "default_region",
	}
	testCloudCredentialAmazonec2Interface = []interface{}{
		map[string]interface{}{
			"access_key":     "access_key",
			"secret_key":     "secret_key",
			"default_region": "default_region",
		},
	}
}

func TestFlattenCloudCredentialAmazonec2(t *testing.T) {

	cases := []struct {
		Input          *amazonec2CredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialAmazonec2Conf,
			testCloudCredentialAmazonec2Interface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialAmazonec2(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialAmazonec2(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *amazonec2CredentialConfig
	}{
		{
			testCloudCredentialAmazonec2Interface,
			testCloudCredentialAmazonec2Conf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialAmazonec2(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
