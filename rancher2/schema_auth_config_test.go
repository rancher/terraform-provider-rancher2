package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAuthConfigFieldsName(t *testing.T) {
	cases := []struct {
		name   string
		fields func() map[string]*schema.Schema
		want   string
	}{
		{"activedirectory", authConfigActiveDirectoryFields, AuthConfigActiveDirectoryName},
		{"adfs", authConfigADFSFields, AuthConfigADFSName},
		{"azuread", authConfigAzureADFields, AuthConfigAzureADName},
		{"cognito", authConfigCognitoFields, AuthConfigCognitoName},
		{"freeipa", authConfigFreeIpaFields, AuthConfigFreeIpaName},
		{"genericoidc", authConfigGenericOIDCFields, AuthConfigGenericOIDCName},
		{"github", authConfigGithubFields, AuthConfigGithubName},
		{"githubapp", authConfigGithubAppFields, AuthConfigGithubAppName},
		{"keycloak", authConfigKeyCloakFields, AuthConfigKeyCloakName},
		{"keycloakoidc", authConfigKeyCloakOIDCFields, AuthConfigKeyCloakOIDCName},
		{"okta", authConfigOKTAFields, AuthConfigOKTAName},
		{"openldap", authConfigOpenLdapFields, AuthConfigOpenLdapName},
		{"ping", authConfigPingFields, AuthConfigPingName},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			name := tc.fields()["name"]
			assert.Equal(t, tc.want, name.Default)
			assert.True(t, name.ForceNew)
		})
	}
}
