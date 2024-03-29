package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		assert.Equal(t, tc.ExpectedOutput, expectedOutput, "Unexpected output from flattener.")
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
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
