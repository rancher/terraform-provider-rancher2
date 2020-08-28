package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigAuthenticationConf      *managementClient.AuthnConfig
	testClusterRKEConfigAuthenticationInterface []interface{}
)

func init() {
	testClusterRKEConfigAuthenticationConf = &managementClient.AuthnConfig{
		SANs:     []string{"sans1", "sans2"},
		Strategy: "strategy",
	}
	testClusterRKEConfigAuthenticationInterface = []interface{}{
		map[string]interface{}{
			"sans":     []interface{}{"sans1", "sans2"},
			"strategy": "strategy",
		},
	}
}

func TestFlattenClusterRKEConfigAuthentication(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AuthnConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigAuthenticationConf,
			testClusterRKEConfigAuthenticationInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigAuthentication(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigAuthentication(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.AuthnConfig
	}{
		{
			testClusterRKEConfigAuthenticationInterface,
			testClusterRKEConfigAuthenticationConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigAuthentication(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
