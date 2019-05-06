package rancher2

import (
	"reflect"
	"testing"
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
