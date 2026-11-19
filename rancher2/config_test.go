package rancher2

import (
	"testing"

	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

func TestFindUserByPrincipalID(t *testing.T) {
	users := []managementClient.User{
		{
			Resource:     types.Resource{ID: "user-local"},
			PrincipalIDs: []string{"local://local:user-local"},
		},
		{
			Resource:     types.Resource{ID: "user-remote"},
			PrincipalIDs: []string{"saml://remote-user", "local://user-remote"},
		},
	}

	tests := []struct {
		name        string
		principalID string
		expected    *managementClient.User
	}{
		{
			name:        "local user principal",
			principalID: "local://local:user-local",
			expected:    &users[0],
		},
		{
			name:        "external user principal",
			principalID: "saml://remote-user",
			expected:    &users[1],
		},
		{
			name:        "secondary local principal",
			principalID: "local://user-remote",
			expected:    &users[1],
		},
		{
			name:        "not found",
			principalID: "unknown://nobody",
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findUserByPrincipalID(users, tt.principalID)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAuthConfigUpdateBody(t *testing.T) {
	t.Run("without keepEmpty returns the object unchanged", func(t *testing.T) {
		obj := &managementClient.GithubConfig{Name: "github"}

		body, err := authConfigUpdateBody(obj)

		assert.NoError(t, err)
		assert.Same(t, obj, body)
	})

	t.Run("sends empty keepEmpty fields", func(t *testing.T) {
		obj := &managementClient.GithubConfig{Name: "github", ClientID: "client-id"}

		body, err := authConfigUpdateBody(obj, managementClient.AuthConfigFieldDisplayName)

		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{
			"name":        "github",
			"clientId":    "client-id",
			"displayName": "",
		}, body)
	})

	t.Run("keeps set keepEmpty fields", func(t *testing.T) {
		obj := &managementClient.GithubConfig{Name: "github", DisplayName: "Login via GitHub"}

		body, err := authConfigUpdateBody(obj, managementClient.AuthConfigFieldDisplayName)

		assert.NoError(t, err)
		assert.Equal(t, "Login via GitHub", body.(map[string]interface{})["displayName"])
	})
}

func TestCheckAuthConfigType(t *testing.T) {
	existing := &managementClient.AuthConfig{Name: "github", Type: managementClient.GithubConfigType}

	t.Run("same type", func(t *testing.T) {
		obj := &managementClient.GithubConfig{Name: "github", Type: managementClient.GithubConfigType}

		assert.NoError(t, checkAuthConfigType("github", existing, obj))
	})

	t.Run("different type", func(t *testing.T) {
		obj := &managementClient.GithubAppConfig{Name: "github", Type: managementClient.GithubAppConfigType}

		err := checkAuthConfigType("github", existing, obj)

		assert.EqualError(t, err, "Auth Config github already exists with type githubConfig, not githubAppConfig")
	})

	t.Run("missing type", func(t *testing.T) {
		obj := &managementClient.GithubConfig{Name: "github"}

		err := checkAuthConfigType("github", existing, obj)

		assert.EqualError(t, err, "Auth Config github has no type")
	})
}

func TestCheckAuthConfigCanBeEnabledSkipsDisabled(t *testing.T) {
	// A disabled AuthConfig never needs the management client.
	obj := &managementClient.OKTAConfig{Name: "okta", Type: managementClient.OKTAConfigType}

	assert.NoError(t, (&Config{}).CheckAuthConfigCanBeEnabled("okta", obj))
}

func TestAuthConfigCreateBody(t *testing.T) {
	t.Run("adds the cleanup annotation", func(t *testing.T) {
		obj := &managementClient.OKTAConfig{Name: "okta", Type: managementClient.OKTAConfigType}

		body, err := authConfigCreateBody(obj)

		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{
			"name": "okta",
			"type": managementClient.OKTAConfigType,
			"annotations": map[string]interface{}{
				"management.cattle.io/auth-provider-cleanup": "rancher-locked",
			},
		}, body)
	})

	t.Run("keeps existing annotations", func(t *testing.T) {
		obj := &managementClient.OKTAConfig{
			Name:        "okta",
			Type:        managementClient.OKTAConfigType,
			Annotations: map[string]string{"example.com/owner": "team-a"},
		}

		body, err := authConfigCreateBody(obj)

		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{
			"example.com/owner":                          "team-a",
			"management.cattle.io/auth-provider-cleanup": "rancher-locked",
		}, body["annotations"])
	})
}
