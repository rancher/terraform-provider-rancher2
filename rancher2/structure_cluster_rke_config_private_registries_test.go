package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterRKEConfigPrivateRegistriesConf      []managementClient.PrivateRegistry
	testClusterRKEConfigPrivateRegistriesInterface []interface{}
)

func init() {
	testClusterRKEConfigPrivateRegistriesConf = []managementClient.PrivateRegistry{
		managementClient.PrivateRegistry{
			IsDefault: true,
			Password:  "XXXXXXXX",
			URL:       "url.terraform.test",
			User:      "user",
		},
	}
	testClusterRKEConfigPrivateRegistriesInterface = []interface{}{
		map[string]interface{}{
			"is_default": true,
			"password":   "XXXXXXXX",
			"url":        "url.terraform.test",
			"user":       "user",
		},
	}
}

func TestFlattenPrivateRegistries(t *testing.T) {

	cases := []struct {
		Input          []managementClient.PrivateRegistry
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigPrivateRegistriesConf,
			testClusterRKEConfigPrivateRegistriesInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigPrivateRegistries(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPrivateRegistries(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.PrivateRegistry
	}{
		{
			testClusterRKEConfigPrivateRegistriesInterface,
			testClusterRKEConfigPrivateRegistriesConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigPrivateRegistries(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
