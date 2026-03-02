package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testAuthConfigCognitoConf      *managementClient.CognitoConfig
	testAuthConfigCognitoInterface map[string]any
)

func init() {
	testAuthConfigCognitoConf = &managementClient.CognitoConfig{
		Name:                AuthConfigCognitoName,
		Type:                managementClient.CognitoConfigType,
		AccessMode:          "access",
		AllowedPrincipalIDs: []string{"allowed1", "allowed2"},
		Enabled:             true,
		Annotations:         map[string]string{"example.com/test": "test"},
		Labels:              map[string]string{"example.com/label": "value"},
		ClientID:            "client_id",
		Issuer:              "issuer",
		RancherURL:          "rancher_url",
		AuthEndpoint:        "auth_endpoint",
		TokenEndpoint:       "token_endpoint",
		UserInfoEndpoint:    "userinfo_endpoint",
		JWKSUrl:             "jwks_url",
		Scopes:              "openid profile offline_access",
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
	testAuthConfigCognitoInterface = map[string]any{
		"name":                  AuthConfigCognitoName,
		"type":                  managementClient.CognitoConfigType,
		"access_mode":           "access",
		"allowed_principal_ids": []any{"allowed1", "allowed2"},
		"enabled":               true,
		"client_id":             "client_id",
		"annotations": map[string]any{
			"example.com/test": "test",
		},
		"labels": map[string]any{
			"example.com/label": "value",
		},
		"issuer":               "issuer",
		"rancher_url":          "rancher_url",
		"auth_endpoint":        "auth_endpoint",
		"token_endpoint":       "token_endpoint",
		"userinfo_endpoint":    "userinfo_endpoint",
		"jwks_url":             "jwks_url",
		"scopes":               "openid profile offline_access",
		"group_search_enabled": true,
		"groups_field":         "groups_field",
		"certificate":          "certificate",
		"private_key":          "private_key",
		"name_claim":           "preferred_username",
		"email_claim":          "alt_email",
		"logout_all_enabled":   true,
		"logout_all_forced":    true,
		"end_session_endpoint": "https://example.com/end-session",
	}
}

func TestFlattenAuthConfigCognito(t *testing.T) {
	output := schema.TestResourceDataRaw(t, authConfigCognitoFields(), map[string]any{})
	err := flattenAuthConfigCognito(output, testAuthConfigCognitoConf)
	if err != nil {
		assert.NoError(t, err, "flattenAuthConfigCognito failed")
	}
	expectedOutput := map[string]any{}
	for k := range testAuthConfigCognitoInterface {
		expectedOutput[k] = output.Get(k)
	}
	assert.Equal(t,
		testAuthConfigCognitoInterface,
		expectedOutput, "Unexpected output from flattener")
}

func TestExpandAuthConfigCognito(t *testing.T) {
	inputResourceData := schema.TestResourceDataRaw(t, authConfigCognitoFields(),
		testAuthConfigCognitoInterface)

	output, err := expandAuthConfigCognito(inputResourceData)
	if err != nil {
		assert.NoError(t, err, "expandAuthConfigCognito failed")
	}
	assert.Equal(t, testAuthConfigCognitoConf, output, "Unexpected output from expander")
}
