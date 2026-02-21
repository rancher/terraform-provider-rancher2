package rancher2

import (
	"fmt"
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testAuthConfigKeyCloakOIDCConf      *managementClient.KeyCloakOIDCConfig
	testAuthConfigKeyCloakOIDCInterface map[string]any
)

func init() {
	// Upstream Rancher representation.
	testAuthConfigKeyCloakOIDCConf = &managementClient.KeyCloakOIDCConfig{
		Name: AuthConfigKeyCloakOIDCName,
		Type: managementClient.KeyCloakOIDCConfigType,
		Labels: map[string]string{
			"example.com/label": "value",
		},
		Annotations: map[string]string{
			"example.com/annotation": "value",
		},
		ClientSecret:        "top-secret-secret",
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

	// Terraform representation.
	testAuthConfigKeyCloakOIDCInterface = map[string]any{
		"labels": map[string]any{
			"example.com/label": "value",
		},
		"annotations": map[string]any{
			"example.com/annotation": "value",
		},
		"name":                  AuthConfigKeyCloakOIDCName,
		"type":                  managementClient.KeyCloakOIDCConfigType,
		"access_mode":           "access",
		"allowed_principal_ids": []any{"allowed1", "allowed2"},
		"enabled":               true,
		"client_id":             "client_id",
		"client_secret":         "top-secret-secret",
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

func TestFlattenAuthConfigKeyCloakOIDC(t *testing.T) {
	output := schema.TestResourceDataRaw(t, authConfigKeyCloakOIDCFields(), map[string]any{})
	err := flattenAuthConfigKeyCloakOIDC(output, testAuthConfigKeyCloakOIDCConf)
	assert.NoError(t, err, "Error in flattenAuthConfigKeyCloakOIDC")
	actualOutput := map[string]any{}
	for k := range testAuthConfigKeyCloakOIDCInterface {
		if k == "client_secret" {
			continue
		}
		actualOutput[k] = output.Get(k)
	}

	// The example data has a client_secret to test sending it to Rancher.
	// But flattening (i.e. extracting the value from Rancher into Terraform)
	// doesn't pull the value because it won't match the terraform version
	// (which contains the secret).
	wantOutput := maps.Clone(testAuthConfigKeyCloakOIDCInterface)
	delete(wantOutput, "client_secret")

	assert.Equal(t, wantOutput, actualOutput, "Unexpected output from flattenAuthConfigKeyCloakOIDC")
}

func TestExpandAuthConfigKeyCloakOIDC(t *testing.T) {
	inputResourceData := schema.TestResourceDataRaw(t, authConfigKeyCloakOIDCFields(), testAuthConfigKeyCloakOIDCInterface)
	output, err := expandAuthConfigKeyCloakOIDC(inputResourceData)

	assert.NoError(t, err, "Error in expandAuthConfigKeyCloakOIDC")
	assert.Equal(t, testAuthConfigKeyCloakOIDCConf, output, "Unexpected output from expandAuthConfigKeyCloakOIDC")
}

func TestExpandAuthConfigKeyCloakOIDCWithAccessModeRequirements(t *testing.T) {
	tests := []struct {
		accessMode string
		wantErr    string
	}{
		{
			"required",
			"expanding keycloakoidc Auth Config: allowed_principal_ids is required on access_mode required",
		},
		{
			"restricted",
			"expanding keycloakoidc Auth Config: allowed_principal_ids is required on access_mode restricted",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("access_mode %q", tt.accessMode), func(t *testing.T) {
			resourceData := map[string]any{
				"labels": map[string]any{
					"example.com/label": "value",
				},
				"annotations": map[string]any{
					"example.com/annotation": "value",
				},
				"name":        AuthConfigKeyCloakOIDCName,
				"type":        managementClient.KeyCloakOIDCConfigType,
				"access_mode": tt.accessMode,
			}

			inputResourceData := schema.TestResourceDataRaw(t, authConfigKeyCloakOIDCFields(), resourceData)
			_, err := expandAuthConfigKeyCloakOIDC(inputResourceData)

			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}
