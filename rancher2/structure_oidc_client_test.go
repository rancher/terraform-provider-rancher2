package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/types"
	client "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFlattenOIDCClient(t *testing.T) {
	oidcClient := &client.OIDCClient{
		Resource: types.Resource{
			ID: "test-client-id",
		},
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
		Status: client.OIDCClientStatus{
			ClientID: "test-client-id",
		},
	}

	flattened := schema.TestResourceDataRaw(t, oidcClientFields(), nil)
	err := flattenOIDCClient(flattened, oidcClient, nil)
	require.NoError(t, err)

	want := map[string]any{
		"token_expiration_seconds":         6000,
		"refresh_token_expiration_seconds": 12000,
		"redirect_uris": []any{
			"http://127.0.0.1:5556/auth/rancher/callback",
			"http://127.0.0.1:33418/",
			"https://vscode.dev/redirect",
		},
		"client_id": "test-client-id",
		"annotations": map[string]any{
			"example.com/testing": "annotation",
		},
		"labels": map[string]any{
			"example.com/testing": "label",
		},
		"description":   "Access for Rancher AI Agent",
		"client_secret": "",
	}
	assert.Equal(t, "test-client-id", flattened.Id(), "unexpected ID from flattenOIDCClient")
	for key, want := range want {
		assert.Equal(t, want, flattened.Get(key), "unexpected output from flattenOIDCClient")
	}
}

func TestFlattenOIDCClientWithSecretKeyMissing(t *testing.T) {
	oidcClient := &client.OIDCClient{
		Resource: types.Resource{
			ID: "test-client-id",
		},
		Status: client.OIDCClientStatus{
			ClientID: "test-client-id",
			ClientSecrets: map[string]client.OIDCClientSecretStatus{
				"client-secret-1": {
					CreatedAt:          "1785140629",
					LastFiveCharacters: "gttjf",
				},
			},
		},
	}
	// Secret exists but does not contain the expected key.
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-client-id",
			Namespace: oidcSecretsNamespace,
		},
		Data: map[string][]byte{},
	}

	flattened := schema.TestResourceDataRaw(t, oidcClientFields(), nil)
	err := flattenOIDCClient(flattened, oidcClient, secret)
	require.NoError(t, err)

	assert.Equal(t, "", flattened.Get("client_secret"), "client_secret should be empty when secret key is missing")
}

func TestFlattenOIDCClientClearsStaleSecret(t *testing.T) {
	oidcClient := &client.OIDCClient{
		Resource: types.Resource{
			ID: "test-client-id",
		},
		Status: client.OIDCClientStatus{
			ClientID: "test-client-id",
		},
	}

	flattened := schema.TestResourceDataRaw(t, oidcClientFields(), nil)
	// Simulate a stale client_secret value left over from a prior Read.
	flattened.Set("client_secret", "stale-secret-value")
	require.Equal(t, "stale-secret-value", flattened.Get("client_secret"), "precondition: stale value should be set")

	err := flattenOIDCClient(flattened, oidcClient, nil)
	require.NoError(t, err)

	assert.Equal(t, "", flattened.Get("client_secret"), "client_secret should be cleared when no secret exists")
}

func TestFlattenOIDCClientWithSecret(t *testing.T) {
	oidcClient := &client.OIDCClient{
		Resource: types.Resource{
			ID: "test-client-id",
		},
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
		Status: client.OIDCClientStatus{
			ClientID: "test-client-id",
			ClientSecrets: map[string]client.OIDCClientSecretStatus{
				"client-secret-1": client.OIDCClientSecretStatus{
					CreatedAt:          "1785140629",
					LastFiveCharacters: "gttjf",
				},
			},
		},
	}
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-client-id",
			Namespace: oidcSecretsNamespace,
		},
		Data: map[string][]byte{"client-secret-1": []byte("testing-secret-gttjf")},
		Type: "secret",
	}

	flattened := schema.TestResourceDataRaw(t, oidcClientFields(), nil)
	err := flattenOIDCClient(flattened, oidcClient, secret)
	require.NoError(t, err)

	want := map[string]any{
		"token_expiration_seconds":         6000,
		"refresh_token_expiration_seconds": 12000,
		"redirect_uris": []any{
			"http://127.0.0.1:5556/auth/rancher/callback",
			"http://127.0.0.1:33418/",
			"https://vscode.dev/redirect",
		},
		"client_id": "test-client-id",
		"annotations": map[string]any{
			"example.com/testing": "annotation",
		},
		"labels": map[string]any{
			"example.com/testing": "label",
		},
		"description":   "Access for Rancher AI Agent",
		"client_secret": "testing-secret-gttjf",
	}
	assert.Equal(t, "test-client-id", flattened.Id(), "unexpected ID from flattenOIDCClient")
	for key, want := range want {
		assert.Equal(t, want, flattened.Get(key), "unexpected output from flattenOIDCClient")
	}
}

func TestExpandOIDCClient(t *testing.T) {
	expandTests := map[string]struct {
		data map[string]any
		want *client.OIDCClient
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
			want: &client.OIDCClient{
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
			want: &client.OIDCClient{
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
