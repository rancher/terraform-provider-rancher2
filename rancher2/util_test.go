package rancher2

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type loginInput struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	TTL         int64  `json:"ttl"`
	Description string `json:"description"`
}

func TestDoUserLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tokenID := "token-xqmfl"
		tokenValue := tokenID + ":fq4nsrx7tcqhbn6bnqlf74mlp8hn8q947nlxbcgv2hsctxnvncprmk"

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/v1-public/login", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var reqBody loginInput
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			require.NoError(t, err)
			assert.Equal(t, "localProvider", reqBody.Type)
			assert.Equal(t, bootstrapDefaultUser, reqBody.Username)
			assert.Equal(t, bootstrapDefaultPassword, reqBody.Password)
			assert.Equal(t, int64(60000), reqBody.TTL)
			assert.Equal(t, bootstrapDefaultSessionDesc, reqBody.Description)

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"token": tokenValue,
			})
		}))
		defer srv.Close()

		id, token, err := DoUserLogin(srv.URL, bootstrapDefaultUser, bootstrapDefaultPassword, bootstrapDefaultTTL, bootstrapDefaultSessionDesc, "", true)
		require.NoError(t, err)
		assert.Equal(t, tokenID, id)
		assert.Equal(t, tokenValue, token)
	})

	t.Run("success with ext prefix", func(t *testing.T) {
		tokenID := "ext/saml-user-abc123"
		tokenValue := tokenID + ":secrettoken"

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"token": tokenValue,
			})
		}))
		defer srv.Close()

		id, token, err := DoUserLogin(srv.URL, bootstrapDefaultUser, bootstrapDefaultPassword, bootstrapDefaultTTL, bootstrapDefaultSessionDesc, "", true)
		require.NoError(t, err)
		assert.Equal(t, "saml-user-abc123", id, "ext/ prefix should be stripped from ID")
		assert.Equal(t, tokenValue, token)
	})

	t.Run("missing token in response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"type": "error",
				"code": "Unauthorized",
			})
		}))
		defer srv.Close()

		_, _, err := DoUserLogin(srv.URL, bootstrapDefaultUser, "wrongpass", bootstrapDefaultTTL, bootstrapDefaultSessionDesc, "", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unauthorized")
	})

	t.Run("invalid token format", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"token": "invalid-token-without-colon",
			})
		}))
		defer srv.Close()

		_, _, err := DoUserLogin(srv.URL, bootstrapDefaultUser, bootstrapDefaultPassword, bootstrapDefaultTTL, bootstrapDefaultSessionDesc, "", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token format")
	})

	t.Run("server error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"type":    "error",
				"code":    "ServerError",
				"message": "internal server error",
			})
		}))
		defer srv.Close()

		_, _, err := DoUserLogin(srv.URL, bootstrapDefaultUser, bootstrapDefaultPassword, bootstrapDefaultTTL, bootstrapDefaultSessionDesc, "", true)
		require.Error(t, err)
	})

	t.Run("fallback to v3-public on 404", func(t *testing.T) {
		tokenID := "token-v3fallback"
		tokenValue := tokenID + ":v3secrettoken"
		callCount := 0

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")

			if r.URL.Path == "/v1-public/login" {
				// Simulate v1 endpoint not available
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"type": "error",
					"code": "NotFound",
				})
				return
			}

			if r.URL.Path == "/v3-public/localProviders/local" {
				assert.Equal(t, "login", r.URL.Query().Get("action"))

				var reqBody loginInput
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				require.NoError(t, err)
				assert.Equal(t, int64(60000), reqBody.TTL)

				_ = json.NewEncoder(w).Encode(map[string]any{
					"token": tokenValue,
				})
				return
			}

			t.Errorf("unexpected request path: %s", r.URL.Path)
		}))
		defer srv.Close()

		id, token, err := DoUserLogin(srv.URL, bootstrapDefaultUser, bootstrapDefaultPassword, bootstrapDefaultTTL, bootstrapDefaultSessionDesc, "", true)
		require.NoError(t, err)
		assert.Equal(t, tokenID, id)
		assert.Equal(t, tokenValue, token)
		assert.Equal(t, 2, callCount, "should have made 2 requests (v1 then v3)")
	})

	t.Run("fallback to v3-public fails", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if r.URL.Path == "/v1-public/login" {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"type": "error",
					"code": "NotFound",
				})
				return
			}

			// v3 endpoint also fails with auth error
			_ = json.NewEncoder(w).Encode(map[string]any{
				"type": "error",
				"code": "Unauthorized",
			})
		}))
		defer srv.Close()

		_, _, err := DoUserLogin(srv.URL, bootstrapDefaultUser, "wrongpass", bootstrapDefaultTTL, bootstrapDefaultSessionDesc, "", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unauthorized")
	})

	t.Run("invalid ttl value - not a number", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("should not make request with invalid TTL")
		}))
		defer srv.Close()

		_, _, err := DoUserLogin(srv.URL, bootstrapDefaultUser, bootstrapDefaultPassword, "not-a-number", bootstrapDefaultSessionDesc, "", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid ttl value")
	})

	t.Run("invalid ttl value - negative number", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("should not make request with negative TTL")
		}))
		defer srv.Close()

		_, _, err := DoUserLogin(srv.URL, bootstrapDefaultUser, bootstrapDefaultPassword, "-1000", bootstrapDefaultSessionDesc, "", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid ttl value")
	})

	t.Run("valid ttl value - zero", func(t *testing.T) {
		tokenID := "token-zero-ttl"
		tokenValue := tokenID + ":secrettoken"

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var reqBody loginInput
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			require.NoError(t, err)
			assert.Equal(t, int64(0), reqBody.TTL)

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"token": tokenValue,
			})
		}))
		defer srv.Close()

		id, token, err := DoUserLogin(srv.URL, bootstrapDefaultUser, bootstrapDefaultPassword, "0", bootstrapDefaultSessionDesc, "", true)
		require.NoError(t, err)
		assert.Equal(t, tokenID, id)
		assert.Equal(t, tokenValue, token)
	})
}
