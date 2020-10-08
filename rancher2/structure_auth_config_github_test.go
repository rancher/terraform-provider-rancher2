package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAuthConfigGithubConf      *managementClient.GithubConfig
	testAuthConfigGithubInterface map[string]interface{}
)

func init() {
	testAuthConfigGithubConf = &managementClient.GithubConfig{
		Name:                AuthConfigGithubName,
		Type:                managementClient.GithubConfigType,
		AccessMode:          "access",
		AllowedPrincipalIDs: []string{"allowed1", "allowed2"},
		Enabled:             true,
		ClientID:            "client_id",
		Hostname:            "hostname",
		TLS:                 true,
	}
	testAuthConfigGithubInterface = map[string]interface{}{
		"name":                  AuthConfigGithubName,
		"type":                  managementClient.GithubConfigType,
		"access_mode":           "access",
		"allowed_principal_ids": []interface{}{"allowed1", "allowed2"},
		"enabled":               true,
		"client_id":             "client_id",
		"hostname":              "hostname",
		"tls":                   true,
	}
}

func TestFlattenAuthConfigGithub(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GithubConfig
		ExpectedOutput map[string]interface{}
	}{
		{
			testAuthConfigGithubConf,
			testAuthConfigGithubInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, authConfigGithubFields(), map[string]interface{}{})
		err := flattenAuthConfigGithub(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandAuthConfigGithub(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.GithubConfig
	}{
		{
			testAuthConfigGithubInterface,
			testAuthConfigGithubConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, authConfigGithubFields(), tc.Input)
		output, err := expandAuthConfigGithub(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
