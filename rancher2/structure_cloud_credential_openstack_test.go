package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
