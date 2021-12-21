package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigPrivateRegistriesECRCredentialsConf      *managementClient.ECRCredentialPlugin
	testClusterRKEConfigPrivateRegistriesECRCredentialsInterface []interface{}
	testClusterRKEConfigPrivateRegistriesConf                    []managementClient.PrivateRegistry
	testClusterRKEConfigPrivateRegistriesInterface               []interface{}
)

func init() {
	testClusterRKEConfigPrivateRegistriesECRCredentialsConf = &managementClient.ECRCredentialPlugin{
		AwsAccessKeyID:     "aws_access_key_id",
		AwsSecretAccessKey: "aws_secret_access_key",
		AwsSessionToken:    "aws_session_token",
	}
	testClusterRKEConfigPrivateRegistriesECRCredentialsInterface = []interface{}{
		map[string]interface{}{
			"aws_access_key_id":     "aws_access_key_id",
			"aws_secret_access_key": "aws_secret_access_key",
			"aws_session_token":     "aws_session_token",
		},
	}
	testClusterRKEConfigPrivateRegistriesConf = []managementClient.PrivateRegistry{
		{
			ECRCredentialPlugin: testClusterRKEConfigPrivateRegistriesECRCredentialsConf,
			IsDefault:           true,
			Password:            "XXXXXXXX",
			URL:                 "url.terraform.test",
			User:                "user",
		},
	}
	testClusterRKEConfigPrivateRegistriesInterface = []interface{}{
		map[string]interface{}{
			"ecr_credential_plugin": testClusterRKEConfigPrivateRegistriesECRCredentialsInterface,
			"is_default":            true,
			"password":              "XXXXXXXX",
			"url":                   "url.terraform.test",
			"user":                  "user",
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
		output, err := flattenClusterRKEConfigPrivateRegistries(tc.Input, tc.ExpectedOutput)
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
