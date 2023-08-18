package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
