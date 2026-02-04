package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testAuthConfigGenericOIDCConf      *managementClient.GenericOIDCConfig
	testAuthConfigGenericOIDCInterface map[string]any
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
		NameClaim:           "preferred_username",
		EmailClaim:          "alt_email",
		LogoutAllEnabled:    true,
		LogoutAllForced:     true,
		EndSessionEndpoint:  "https://example.com/end-session",
	}
	testAuthConfigGenericOIDCInterface = map[string]any{
		"name":                  AuthConfigGenericOIDCName,
		"type":                  managementClient.GenericOIDCConfigType,
		"access_mode":           "access",
		"allowed_principal_ids": []any{"allowed1", "allowed2"},
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
		"name_claim":            "preferred_username",
		"email_claim":           "alt_email",
		"logout_all_enabled":    true,
		"logout_all_forced":     true,
		"end_session_endpoint":  "https://example.com/end-session",
	}
}

func TestFlattenAuthConfigGenericOIDC(t *testing.T) {
	output := schema.TestResourceDataRaw(t, authConfigGenericOIDCFields(), map[string]any{})
	err := flattenAuthConfigGenericOIDC(output, testAuthConfigGenericOIDCConf)
	assert.NoError(t, err, "Error in flattenAuthConfigGenericOIDC")
	expectedOutput := map[string]any{}
	for k := range testAuthConfigGenericOIDCInterface {
		expectedOutput[k] = output.Get(k)
	}
	assert.Equal(t, testAuthConfigGenericOIDCInterface, expectedOutput, "Unexpected output from flattenAuthConfigGenericOIDC")
}

func TestExpandAuthConfigGenericOIDC(t *testing.T) {
	inputResourceData := schema.TestResourceDataRaw(t, authConfigGenericOIDCFields(), testAuthConfigGenericOIDCInterface)
	output, err := expandAuthConfigGenericOIDC(inputResourceData)

	assert.NoError(t, err, "Error in expandAuthConfigGenericOIDC")
	assert.Equal(t, testAuthConfigGenericOIDCConf, output, "Unexpected output from expandAuthConfigGenericOIDC")
}
