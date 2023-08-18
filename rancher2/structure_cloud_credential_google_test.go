package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
