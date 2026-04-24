package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlattenOIDCClient(t *testing.T) {
	oidcClient := &managementClient.OIDCClient{
		Name: "testing-client",
		Annotations: map[string]string{
			"example.com/testing": "annotation",
		},
		Labels: map[string]string{
			"example.com/testing": "label",
		},
		RedirectURIs: []string{
			"http://127.0.0.1:5556/auth/rancher/callback",
			"http://127.0.0.1:33418/",
			"https://vscode.dev/redirect",
		},
		Description:                   "Access for Rancher AI Agent",
		TokenExpirationSeconds:        6000,
		RefreshTokenExpirationSeconds: 12000,
	}

	flattened := schema.TestResourceDataRaw(t, oidcClientFields(), nil)
	err := flattenOIDCClient(flattened, oidcClient)
	require.NoError(t, err)

	want := map[string]any{
		"token_expiration_seconds":         6000,
		"refresh_token_expiration_seconds": 12000,
		"redirect_uris": []any{
			"http://127.0.0.1:5556/auth/rancher/callback",
			"http://127.0.0.1:33418/",
			"https://vscode.dev/redirect",
		},
		"annotations": map[string]any{
			"example.com/testing": "annotation",
		},
		"labels": map[string]any{
			"example.com/testing": "label",
		},
	}
	for key, want := range want {
		assert.Equal(t, want, flattened.Get(key), "unexpected output from flattenOIDCClient")
	}
}

func TestExpandOIDCClient(t *testing.T) {
	expandTests := map[string]struct {
		data map[string]any
		want *managementClient.OIDCClient
	}{
		"all fields populated": {
			data: map[string]any{
				"token_expiration_seconds":         6000,
				"refresh_token_expiration_seconds": 12000,
				"redirect_uris": []any{
					"http://127.0.0.1:5556/auth/rancher/callback",
					"http://127.0.0.1:33418/",
					"https://vscode.dev/redirect",
				},
				"description": "Testing OIDC Client",
				"annotations": map[string]any{
					"example.com/testing": "annotation",
				},
				"labels": map[string]any{
					"example.com/testing": "label",
				},
			},
			want: &managementClient.OIDCClient{
				TokenExpirationSeconds:        6000,
				RefreshTokenExpirationSeconds: 12000,
				RedirectURIs: []string{
					"http://127.0.0.1:5556/auth/rancher/callback",
					"http://127.0.0.1:33418/",
					"https://vscode.dev/redirect",
				},
				Description: "Testing OIDC Client",
				Annotations: map[string]string{
					"example.com/testing": "annotation",
				},
				Labels: map[string]string{
					"example.com/testing": "label",
				},
			},
		},
		"only required fields populated": {
			data: map[string]any{
				"redirect_uris": []any{
					"http://127.0.0.1:5556/auth/rancher/callback",
					"http://127.0.0.1:33418/",
					"https://vscode.dev/redirect",
				},
			},
			want: &managementClient.OIDCClient{
				TokenExpirationSeconds: 0,
				RedirectURIs: []string{
					"http://127.0.0.1:5556/auth/rancher/callback",
					"http://127.0.0.1:33418/",
					"https://vscode.dev/redirect",
				},
			},
		},
	}

	for name, tt := range expandTests {
		t.Run(name, func(t *testing.T) {
			inputResourceData := schema.TestResourceDataRaw(t, oidcClientFields(), tt.data)

			expanded, err := expandOIDCClient(inputResourceData)
			assert.NoError(t, err, "Error in expandOIDCClient")

			assert.Equal(t, tt.want, expanded, "Unexpected output from expandOIDCClient")
		})
	}
}
