package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testAuthConfigGenericOIDCConf      *managementClient.GenericOIDCConfig
	testAuthConfigGenericOIDCInterface map[string]interface{}
	groupSearchEnabled                 = true
)

func init() {
	testAuthConfigGenericOIDCConf = &managementClient.GenericOIDCConfig{
		Name:                AuthConfigGenericOIDCName,
		Type:                managementClient.GenericOIDCConfigType,
		AccessMode:          "access",
		AllowedPrincipalIDs: []string{"allowed1", "allowed2"},
		Enabled:             true,
		ClientID:            "client_id",
		Issuer:              "issuer",
		RancherURL:          "rancher_url",
		AuthEndpoint:        "auth_endpoint",
		TokenEndpoint:       "token_endpoint",
		UserInfoEndpoint:    "userinfo_endpoint",
		JWKSUrl:             "jwks_url",
		Scopes:              "scopes",
		GroupSearchEnabled:  &groupSearchEnabled,
		GroupsClaim:         "groups_field",
		Certificate:         "certificate",
		PrivateKey:          "private_key",
	}
	testAuthConfigGenericOIDCInterface = map[string]interface{}{
		"name":                  AuthConfigGenericOIDCName,
		"type":                  managementClient.GenericOIDCConfigType,
		"access_mode":           "access",
		"allowed_principal_ids": []interface{}{"allowed1", "allowed2"},
		"enabled":               true,
		"client_id":             "client_id",
		"issuer":                "issuer",
		"rancher_url":           "rancher_url",
		"auth_endpoint":         "auth_endpoint",
		"token_endpoint":        "token_endpoint",
		"userinfo_endpoint":     "userinfo_endpoint",
		"jwks_url":              "jwks_url",
		"scopes":                "scopes",
		"group_search_enabled":  true,
		"groups_field":          "groups_field",
		"certificate":           "certificate",
		"private_key":           "private_key",
	}
}

func TestFlattenAuthConfigGenericOIDC(t *testing.T) {
	cases := []struct {
		Input          *managementClient.GenericOIDCConfig
		ExpectedOutput map[string]interface{}
	}{
		{
			testAuthConfigGenericOIDCConf,
			testAuthConfigGenericOIDCInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, authConfigGenericOIDCFields(), map[string]interface{}{})
		err := flattenAuthConfigGenericOIDC(output, tc.Input)
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

func TestExpandAuthConfigGenericOIDC(t *testing.T) {
	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.GenericOIDCConfig
	}{
		{
			testAuthConfigGenericOIDCInterface,
			testAuthConfigGenericOIDCConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, authConfigGenericOIDCFields(), tc.Input)
		output, err := expandAuthConfigGenericOIDC(inputResourceData)
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
