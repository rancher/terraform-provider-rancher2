package rancher2

import (
	"reflect"
	"testing"
)

var (
	testCloudCredentialDigitaloceanConf      *digitaloceanCredentialConfig
	testCloudCredentialDigitaloceanInterface []interface{}
)

func init() {
	testCloudCredentialDigitaloceanConf = &digitaloceanCredentialConfig{
		AccessToken: "access_token",
	}
	testCloudCredentialDigitaloceanInterface = []interface{}{
		map[string]interface{}{
			"access_token": "access_token",
		},
	}
}

func TestFlattenCloudCredentialDigitalocean(t *testing.T) {

	cases := []struct {
		Input          *digitaloceanCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialDigitaloceanConf,
			testCloudCredentialDigitaloceanInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialDigitalocean(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialDigitalocean(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *digitaloceanCredentialConfig
	}{
		{
			testCloudCredentialDigitaloceanInterface,
			testCloudCredentialDigitaloceanConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialDigitalocean(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
