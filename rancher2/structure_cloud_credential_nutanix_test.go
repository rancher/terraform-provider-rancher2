package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCloudCredentialNutanixConf      *nutanixCredentialConfig
	testCloudCredentialNutanixInterface []interface{}
)

func init() {
	testCloudCredentialNutanixConf = &nutanixCredentialConfig{
		Endpoint: "pc.example.com",
		Username: "X-ntnx-api-key",
		Password: "secret",
		Port:     "9440",
	}
	testCloudCredentialNutanixInterface = []interface{}{
		map[string]interface{}{
			"endpoint": "pc.example.com",
			"username": "X-ntnx-api-key",
			"password": "secret",
			"port":     "9440",
		},
	}
}

func TestFlattenCloudCredentialNutanix(t *testing.T) {
	cases := []struct {
		Input          *nutanixCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialNutanixConf,
			testCloudCredentialNutanixInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialNutanix(tc.Input, tc.ExpectedOutput)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandCloudCredentialNutanix(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *nutanixCredentialConfig
	}{
		{
			testCloudCredentialNutanixInterface,
			testCloudCredentialNutanixConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialNutanix(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
