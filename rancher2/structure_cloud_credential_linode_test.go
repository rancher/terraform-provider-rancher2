package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
