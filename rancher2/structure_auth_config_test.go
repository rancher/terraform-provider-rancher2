package rancher2

import (
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type authConfigStructureTestCase struct {
	name    string
	fields  func() map[string]*schema.Schema
	input   map[string]interface{}
	expand  func(*schema.ResourceData) (string, string, error)
	flatten func(d *schema.ResourceData, name, displayName string) error
}

func authConfigStructureTestCases() []authConfigStructureTestCase {
	return []authConfigStructureTestCase{
		{
			name:   "activedirectory",
			fields: authConfigActiveDirectoryFields,
			input:  testAuthConfigActiveDirectoryInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigActiveDirectory(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigActiveDirectoryConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigActiveDirectory(d, &in)
			},
		},
		{
			name:   "adfs",
			fields: authConfigADFSFields,
			input:  testAuthConfigADFSInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigADFS(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigADFSConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigADFS(d, &in)
			},
		},
		{
			name:   "azuread",
			fields: authConfigAzureADFields,
			input:  testAuthConfigAzureADInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigAzureAD(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigAzureADConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigAzureAD(d, &in)
			},
		},
		{
			name:   "cognito",
			fields: authConfigCognitoFields,
			input:  testAuthConfigCognitoInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigCognito(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigCognitoConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigCognito(d, &in)
			},
		},
		{
			name:   "freeipa",
			fields: authConfigFreeIpaFields,
			input:  testAuthConfigFreeIpaInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigFreeIpa(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigFreeIpaConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigFreeIpa(d, &in)
			},
		},
		{
			name:   "genericoidc",
			fields: authConfigGenericOIDCFields,
			input:  testAuthConfigGenericOIDCInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigGenericOIDC(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigGenericOIDCConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigGenericOIDC(d, &in)
			},
		},
		{
			name:   "github",
			fields: authConfigGithubFields,
			input:  testAuthConfigGithubInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigGithub(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigGithubConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigGithub(d, &in)
			},
		},
		{
			name:   "githubapp",
			fields: authConfigGithubAppFields,
			input:  testAuthConfigGithubAppInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigGithubApp(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigGithubAppConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigGithubApp(d, &in)
			},
		},
		{
			name:   "keycloak",
			fields: authConfigKeyCloakFields,
			input:  testAuthConfigKeyCloakInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigKeyCloak(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigKeyCloakConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigKeyCloak(d, &in)
			},
		},
		{
			name:   "keycloakoidc",
			fields: authConfigKeyCloakOIDCFields,
			input:  testAuthConfigKeyCloakOIDCInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigKeyCloakOIDC(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigKeyCloakOIDCConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigKeyCloakOIDC(d, &in)
			},
		},
		{
			name:   "okta",
			fields: authConfigOKTAFields,
			input:  testAuthConfigOKTAInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigOKTA(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigOKTAConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigOKTA(d, &in)
			},
		},
		{
			name:   "openldap",
			fields: authConfigOpenLdapFields,
			input:  testAuthConfigOpenLdapInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigOpenLdap(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigOpenLdapConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigOpenLdap(d, &in)
			},
		},
		{
			name:   "ping",
			fields: authConfigPingFields,
			input:  testAuthConfigPingInterface,
			expand: func(d *schema.ResourceData) (string, string, error) {
				obj, err := expandAuthConfigPing(d)
				if err != nil {
					return "", "", err
				}
				return obj.Name, obj.DisplayName, nil
			},
			flatten: func(d *schema.ResourceData, name, displayName string) error {
				in := *testAuthConfigPingConf
				in.Name = name
				in.DisplayName = displayName
				return flattenAuthConfigPing(d, &in)
			},
		},
	}
}

func TestExpandAuthConfigCustomName(t *testing.T) {
	for _, tc := range authConfigStructureTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			input := maps.Clone(tc.input)
			input["name"] = "custom-" + tc.name

			name, _, err := tc.expand(schema.TestResourceDataRaw(t, tc.fields(), input))

			require.NoError(t, err)
			assert.Equal(t, "custom-"+tc.name, name)
		})
	}
}

func TestExpandAuthConfigDefaultName(t *testing.T) {
	for _, tc := range authConfigStructureTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			input := maps.Clone(tc.input)
			delete(input, "name")

			name, _, err := tc.expand(schema.TestResourceDataRaw(t, tc.fields(), input))

			require.NoError(t, err)
			assert.Equal(t, tc.fields()["name"].Default, name)
		})
	}
}

func TestExpandAuthConfigDisplayName(t *testing.T) {
	for _, tc := range authConfigStructureTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			input := maps.Clone(tc.input)
			input["display_name"] = "Login via " + tc.name

			_, displayName, err := tc.expand(schema.TestResourceDataRaw(t, tc.fields(), input))

			require.NoError(t, err)
			assert.Equal(t, "Login via "+tc.name, displayName)
		})
	}
}

func TestFlattenAuthConfigNameAndDisplayName(t *testing.T) {
	for _, tc := range authConfigStructureTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, tc.fields(), map[string]interface{}{})

			err := tc.flatten(d, "custom-"+tc.name, "Login via "+tc.name)

			require.NoError(t, err)
			assert.Equal(t, "custom-"+tc.name, d.Id())
			assert.Equal(t, "custom-"+tc.name, d.Get("name"))
			assert.Equal(t, "Login via "+tc.name, d.Get("display_name"))
		})
	}
}
