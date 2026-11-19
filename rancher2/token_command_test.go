package rancher2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCmdToken = "token-abc:secret"

func TestParseTokenCommandOutput(t *testing.T) {
	cases := []struct {
		name    string
		out     string
		want    string
		wantErr string
	}{
		{"raw token", testCmdToken, testCmdToken, ""},
		{"raw token with newline", testCmdToken + "\n", testCmdToken, ""},
		{"exec credential", `{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1beta1","status":{"token":"` + testCmdToken + `","expirationTimestamp":"2030-01-01T00:00:00Z"}}`, testCmdToken, ""},
		{"empty", "  \n", "", "printed no token"},
		{"json without token", `{"kind":"ExecCredential","status":{}}`, "", "no status.token"},
		{"broken json", `{"kind":`, "", "can't be parsed"},
		{"multiple words", "some log line\n" + testCmdToken, "", "neither a single token"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseTokenCommandOutput([]byte(tc.out))
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				assert.NotContains(t, err.Error(), "secret")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRunTokenCommand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := runTokenCommand([]string{"sh", "-c", "printf '%s\\n' " + testCmdToken}, 5*time.Second)
		require.NoError(t, err)
		assert.Equal(t, testCmdToken, got)
	})

	t.Run("failure never includes command output", func(t *testing.T) {
		_, err := runTokenCommand([]string{"sh", "-c", "echo " + testCmdToken + "; echo login required " + testCmdToken + " >&2; exit 3"}, 5*time.Second)
		require.Error(t, err)
		assert.NotContains(t, err.Error(), testCmdToken)
		assert.NotContains(t, err.Error(), "login required")
		assert.Contains(t, err.Error(), "exit status 3")
		assert.Contains(t, err.Error(), "run the command in a terminal")
	})

	t.Run("timeout", func(t *testing.T) {
		_, err := runTokenCommand([]string{"sh", "-c", "sleep 5"}, 200*time.Millisecond)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timed out")
	})

	t.Run("timeout with descendant holding the output pipes", func(t *testing.T) {
		start := time.Now()
		_, err := runTokenCommand([]string{"sh", "-c", "sleep 30 & wait"}, 200*time.Millisecond)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timed out")
		assert.Less(t, time.Since(start), 10*time.Second)
	})

	t.Run("unbounded output is cut off", func(t *testing.T) {
		start := time.Now()
		_, err := runTokenCommand([]string{"yes"}, 30*time.Second)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "more than")
		assert.NotContains(t, err.Error(), "timed out")
		assert.Less(t, time.Since(start), 10*time.Second)
	})

	t.Run("empty command", func(t *testing.T) {
		_, err := runTokenCommand(nil, time.Second)
		require.Error(t, err)
	})

	t.Run("missing program", func(t *testing.T) {
		_, err := runTokenCommand([]string{"definitely-not-a-real-program-xyz"}, time.Second)
		require.Error(t, err)
	})
}

func TestProviderConfigureTokenCommand(t *testing.T) {
	for _, k := range []string{"RANCHER_TOKEN_KEY", "RANCHER_ACCESS_KEY", "RANCHER_SECRET_KEY", "RANCHER_BOOTSTRAP"} {
		t.Setenv(k, "")
	}
	printToken := []interface{}{"sh", "-c", "printf " + testCmdToken}
	failing := []interface{}{"sh", "-c", "exit 1"}

	cases := []struct {
		name      string
		raw       map[string]interface{}
		wantToken string
		wantErr   string
	}{
		{"token_command used when no static credentials", map[string]interface{}{"api_url": "https://rancher.example.com", "token_command": printToken}, testCmdToken, ""},
		{"partial access_key is an error, not a token_command fallback", map[string]interface{}{"api_url": "https://rancher.example.com", "access_key": "token-static", "token_command": printToken}, "", "No token_key"},
		{"partial secret_key is an error, not a token_command fallback", map[string]interface{}{"api_url": "https://rancher.example.com", "secret_key": "key", "token_command": printToken}, "", "No token_key"},
		{"token_key wins over token_command", map[string]interface{}{"api_url": "https://rancher.example.com", "token_key": "token-static:key", "token_command": failing}, "token-static:key", ""},
		{"access/secret key win over token_command", map[string]interface{}{"api_url": "https://rancher.example.com", "access_key": "token-static", "secret_key": "key", "token_command": failing}, "token-static:key", ""},
		{"token_command error is returned", map[string]interface{}{"api_url": "https://rancher.example.com", "token_command": failing}, "", "token_command"},
		{"no credentials at all", map[string]interface{}{"api_url": "https://rancher.example.com"}, "", "No token_key"},
		{"bootstrap ignores token_command", map[string]interface{}{"api_url": "https://rancher.example.com", "bootstrap": true, "token_command": failing}, providerDefaultEmptyString, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, tc.raw)
			got, err := providerConfigure(d)
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantToken, got.(*Config).TokenKey)
		})
	}
}
