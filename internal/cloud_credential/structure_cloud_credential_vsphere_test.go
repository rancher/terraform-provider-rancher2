package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCloudCredentialVsphereConf      *vmwarevsphereCredentialConfig
	testCloudCredentialVsphereInterface []interface{}
)

func init() {
	testCloudCredentialVsphereConf = &vmwarevsphereCredentialConfig{
		Password:    "password",
		Username:    "username",
		Vcenter:     "vcenter",
		VcenterPort: "443",
	}
	testCloudCredentialVsphereInterface = []interface{}{
		map[string]interface{}{
			"password":     "password",
			"username":     "username",
			"vcenter":      "vcenter",
			"vcenter_port": "443",
		},
	}
}

func TestFlattenCloudCredentialVsphere(t *testing.T) {

	cases := []struct {
		Input          *vmwarevsphereCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialVsphereConf,
			testCloudCredentialVsphereInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialVsphere(tc.Input, tc.ExpectedOutput)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandCloudCredentialVsphere(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *vmwarevsphereCredentialConfig
	}{
		{
			testCloudCredentialVsphereInterface,
			testCloudCredentialVsphereConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialVsphere(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
