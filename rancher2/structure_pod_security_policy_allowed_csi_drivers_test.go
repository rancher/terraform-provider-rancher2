package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPodSecurityPolicyAllowedCSIDriversConf           []managementClient.AllowedCSIDriver
	testPodSecurityPolicyAllowedCSIDriversInterface      []interface{}
	testEmptyPodSecurityPolicyAllowedCSIDriversConf      []managementClient.AllowedCSIDriver
	testEmptyPodSecurityPolicyAllowedCSIDriversInterface []interface{}
)

func init() {
	testPodSecurityPolicyAllowedCSIDriversConf = []managementClient.AllowedCSIDriver{
		{
			Name: "foo",
		},
		{
			Name: "bar",
		},
	}
	testPodSecurityPolicyAllowedCSIDriversInterface = []interface{}{
		map[string]interface{}{
			"name": "foo",
		},
		map[string]interface{}{
			"name": "bar",
		},
	}
	testEmptyPodSecurityPolicyAllowedCSIDriversInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyAllowedCSIDrivers(t *testing.T) {

	cases := []struct {
		Input          []managementClient.AllowedCSIDriver
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyAllowedCSIDriversConf,
			testPodSecurityPolicyAllowedCSIDriversInterface,
		},
		{
			testEmptyPodSecurityPolicyAllowedCSIDriversConf,
			testEmptyPodSecurityPolicyAllowedCSIDriversInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyAllowedCSIDrivers(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicyAllowedCSIDrivers(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.AllowedCSIDriver
	}{
		{
			testPodSecurityPolicyAllowedCSIDriversInterface,
			testPodSecurityPolicyAllowedCSIDriversConf,
		},
	}
	for _, tc := range cases {
		output := expandPodSecurityPolicyAllowedCSIDrivers(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
