package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCloudCredentialIonoscloudConf      *ionoscloudCredentialConfig
	testCloudCredentialIonoscloudInterface []interface{}
)

func init() {
	testCloudCredentialIonoscloudConf = &ionoscloudCredentialConfig{
		Token:    "123",
		Username: "Pavels",
		Password: "1337",
		Endpoint: "https://api.ionos.com",
	}
	testCloudCredentialIonoscloudInterface = []interface{}{
		map[string]interface{}{
			"token":    "123",
			"username": "Pavels",
			"password": "1337",
			"endpoint": "https://api.ionos.com",
		},
	}
}

func TestFlattenCloudCredentialIonoscloud(t *testing.T) {

	cases := []struct {
		Input          *ionoscloudCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialIonoscloudConf,
			testCloudCredentialIonoscloudInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialIonoscloud(tc.Input, tc.ExpectedOutput)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandCloudCredentialIonoscloud(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *ionoscloudCredentialConfig
	}{
		{
			testCloudCredentialIonoscloudInterface,
			testCloudCredentialIonoscloudConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialIonoscloud(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
