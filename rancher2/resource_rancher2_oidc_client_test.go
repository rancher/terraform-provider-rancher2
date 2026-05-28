package rancher2

import (
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceRancher2OIDCClientRead(t *testing.T) {
	newGetter := func(ops managementClient.OIDCClientOperations) *fakeManagementClientGetter {
		return &fakeManagementClientGetter{
			client: &managementClient.Client{OIDCClient: ops},
		}
	}

	newResourceData := func() *schema.ResourceData {
		d := schema.TestResourceDataRaw(t, oidcClientFields(), map[string]any{
			"redirect_uris": []any{},
		})
		d.SetId("oidc-client-1")
		return d
	}

	t.Run("populates state when resource exists", func(t *testing.T) {
		existing := &managementClient.OIDCClient{
			Resource: types.Resource{
				ID: "oidc-client-1",
			},
			Description:                   "Test OIDC Client",
			TokenExpirationSeconds:        3600,
			RefreshTokenExpirationSeconds: 7200,
			RedirectURIs:                  []string{"http://127.0.0.1:5556/auth/rancher/callback"},
			Status: managementClient.OIDCClientStatus{
				ClientID: "my-client-id",
			},
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientRead(d, newGetter(&fakeOIDCClientOperations{oidcClient: existing}))

		require.NoError(t, err)
		assert.Equal(t, "oidc-client-1", d.Id())
		assert.Equal(t, "Test OIDC Client", d.Get("description"))
		assert.Equal(t, 3600, d.Get("token_expiration_seconds"))
		assert.Equal(t, 7200, d.Get("refresh_token_expiration_seconds"))
		assert.Equal(t, "my-client-id", d.Get("client_id"))
	})

	t.Run("clears state when resource does not exist or is not accessible by ID", func(t *testing.T) {

		clearIDTests := map[string]struct {
			err error
		}{
			"not found clears resource ID": {
				err: &clientbase.APIError{StatusCode: http.StatusNotFound},
			},
			"forbidden clears resource ID": {
				err: &clientbase.APIError{StatusCode: http.StatusForbidden},
			},
			"not accessible by ID clears resource ID": {
				err: errors.New("can not be looked up by ID"),
			},
		}

		for name, tc := range clearIDTests {
			t.Run(name, func(t *testing.T) {
				d := newResourceData()
				err := resourceRancher2OIDCClientRead(d, newGetter(&fakeOIDCClientOperations{err: tc.err}))

				require.NoError(t, err)
				assert.Empty(t, d.Id())
			})
		}
	})

	t.Run("returns error on unexpected API error", func(t *testing.T) {
		d := newResourceData()
		err := resourceRancher2OIDCClientRead(d, newGetter(&fakeOIDCClientOperations{err: errors.New("internal server error")}))

		require.Error(t, err)
		assert.Equal(t, "oidc-client-1", d.Id())
	})
}

func TestResourceRancher2OIDCClientCreate(t *testing.T) {
	newGetter := func(ops managementClient.OIDCClientOperations) *fakeManagementClientGetter {
		return &fakeManagementClientGetter{
			client: &managementClient.Client{OIDCClient: ops},
		}
	}

	newResourceData := func() *schema.ResourceData {
		return schema.TestResourceDataRaw(t, oidcClientFields(), map[string]any{
			"description": "Test OIDC Client",
			"redirect_uris": []any{
				"http://127.0.0.1:5556/auth/rancher/callback",
			},
			"token_expiration_seconds":         3600,
			"refresh_token_expiration_seconds": 7200,
		})
	}

	t.Run("sets ID and populates state on success", func(t *testing.T) {
		created := &managementClient.OIDCClient{
			Resource: types.Resource{
				ID: "oidc-client-new",
			},
			Description:                   "Test OIDC Client",
			TokenExpirationSeconds:        3600,
			RefreshTokenExpirationSeconds: 7200,
			RedirectURIs:                  []string{"http://127.0.0.1:5556/auth/rancher/callback"},
		}

		// ByID is called by the Read that follows Create.
		ops := &fakeOIDCClientOperations{
			createResult: created,
			oidcClient:   created,
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientCreate(d, newGetter(ops))

		require.NoError(t, err)
		assert.Equal(t, "oidc-client-new", d.Id())
		assert.Equal(t, "Test OIDC Client", d.Get("description"))
		assert.Equal(t, 3600, d.Get("token_expiration_seconds"))
		assert.Equal(t, 7200, d.Get("refresh_token_expiration_seconds"))
	})

	t.Run("returns error when Create API call fails", func(t *testing.T) {
		ops := &fakeOIDCClientOperations{
			createErr: errors.New("server unavailable"),
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientCreate(d, newGetter(ops))

		require.Error(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("clears ID when Read after Create returns not found", func(t *testing.T) {
		created := &managementClient.OIDCClient{
			Resource: types.Resource{
				ID: "oidc-client-new",
			},
		}

		ops := &fakeOIDCClientOperations{
			createResult: created,
			err:          &clientbase.APIError{StatusCode: http.StatusNotFound},
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientCreate(d, newGetter(ops))

		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestResourceRancher2OIDCClientDelete(t *testing.T) {
	newGetter := func(ops managementClient.OIDCClientOperations) *fakeManagementClientGetter {
		return &fakeManagementClientGetter{
			client: &managementClient.Client{OIDCClient: ops},
		}
	}

	newResourceData := func() *schema.ResourceData {
		d := schema.TestResourceDataRaw(t, oidcClientFields(), map[string]any{
			"redirect_uris": []any{},
		})
		d.SetId("oidc-client-1")
		return d
	}

	t.Run("deletes resource successfully", func(t *testing.T) {
		existing := &managementClient.OIDCClient{}
		existing.ID = "oidc-client-1"

		ops := &fakeOIDCClientOperations{oidcClient: existing}

		d := newResourceData()
		err := resourceRancher2OIDCClientDelete(d, newGetter(ops))

		require.NoError(t, err)
	})

	t.Run("ByID fails with not found or forbidden error and Delete is not called, resulting in cleared ID", func(t *testing.T) {
		clearIDTests := map[string]struct {
			err error
		}{
			"not found on ByID clears resource ID": {
				err: &clientbase.APIError{StatusCode: http.StatusNotFound},
			},
			"forbidden on ByID clears resource ID": {
				err: &clientbase.APIError{StatusCode: http.StatusForbidden},
			},
			// This comes from Norman code (clientbase/ops.go)
			"not accessible by ID on ByID clears resource ID": {
				err: errors.New("Resource Type oidcClient can not be looked up by ID"),
			},
		}

		for name, tc := range clearIDTests {
			t.Run(name, func(t *testing.T) {
				ops := &fakeOIDCClientOperations{err: tc.err}

				d := newResourceData()
				err := resourceRancher2OIDCClientDelete(d, newGetter(ops))

				require.NoError(t, err)
				assert.Empty(t, d.Id())
			})
		}
	})

	t.Run("returns error on unexpected ByID error", func(t *testing.T) {
		ops := &fakeOIDCClientOperations{err: errors.New("connection refused")}

		d := newResourceData()
		err := resourceRancher2OIDCClientDelete(d, newGetter(ops))

		require.ErrorContains(t, err, "connection refused")
	})

	t.Run("returns error when Delete API call fails", func(t *testing.T) {
		existing := &managementClient.OIDCClient{
			Resource: types.Resource{
				ID: "oidc-client-1",
			},
		}

		ops := &fakeOIDCClientOperations{
			oidcClient: existing,
			deleteErr:  errors.New("delete rejected"),
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientDelete(d, newGetter(ops))

		require.ErrorContains(t, err, "delete rejected")
	})

	t.Run("succeeds when Delete returns not found", func(t *testing.T) {
		existing := &managementClient.OIDCClient{
			Resource: types.Resource{
				ID: "oidc-client-1",
			},
		}

		ops := &fakeOIDCClientOperations{
			oidcClient: existing,
			deleteErr:  &clientbase.APIError{StatusCode: http.StatusNotFound},
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientDelete(d, newGetter(ops))

		require.NoError(t, err)
		assert.Empty(t, d.Id(), "Resource ID should be cleared when Delete returns not found")
	})
}

func TestResourceRancher2OIDCClientUpdate(t *testing.T) {
	newGetter := func(ops managementClient.OIDCClientOperations) *fakeManagementClientGetter {
		return &fakeManagementClientGetter{
			client: &managementClient.Client{OIDCClient: ops},
		}
	}

	newResourceData := func() *schema.ResourceData {
		d := schema.TestResourceDataRaw(t, oidcClientFields(), map[string]any{
			"description": "Updated description",
			"redirect_uris": []any{
				"http://127.0.0.1:5556/auth/rancher/callback",
			},
			"token_expiration_seconds":         3600,
			"refresh_token_expiration_seconds": 7200,
		})
		d.SetId("oidc-client-1")
		return d
	}

	t.Run("sends correct payload to Update", func(t *testing.T) {
		existing := &managementClient.OIDCClient{
			Resource: types.Resource{
				ID: "oidc-client-1",
			},
		}

		ops := &fakeOIDCClientOperations{
			oidcClient:   existing,
			updateResult: existing,
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientUpdate(d, newGetter(ops))

		require.NoError(t, err)
		assert.Same(t, existing, ops.updateCalledOn)
		wantPayload := map[string]any{
			"labels":      map[string]string{},
			"annotations": map[string]string{},
			"description": "Updated description",
			"redirectURIs": []string{
				"http://127.0.0.1:5556/auth/rancher/callback",
			},
			"tokenExpirationSeconds":        3600,
			"refreshTokenExpirationSeconds": 7200,
		}
		assert.Equal(t, wantPayload, ops.updateCalledWith)
	})

	t.Run("sends labels and annotations in update payload", func(t *testing.T) {
		existing := &managementClient.OIDCClient{
			Resource: types.Resource{
				ID: "oidc-client-1",
			},
		}

		ops := &fakeOIDCClientOperations{
			oidcClient:   existing,
			updateResult: existing,
		}

		d := schema.TestResourceDataRaw(t, oidcClientFields(), map[string]any{
			"description": "Updated description",
			"redirect_uris": []any{
				"http://127.0.0.1:5556/auth/rancher/callback",
			},
			"token_expiration_seconds":         3600,
			"refresh_token_expiration_seconds": 7200,
			"labels": map[string]any{
				"example.com/env": "production",
			},
			"annotations": map[string]any{
				"example.com/owner": "team-a",
			},
		})
		d.SetId("oidc-client-1")

		err := resourceRancher2OIDCClientUpdate(d, newGetter(ops))

		require.NoError(t, err)
		assert.Same(t, existing, ops.updateCalledOn)
		wantPayload := map[string]any{
			"labels": map[string]string{
				"example.com/env": "production",
			},
			"annotations": map[string]string{
				"example.com/owner": "team-a",
			},
			"description": "Updated description",
			"redirectURIs": []string{
				"http://127.0.0.1:5556/auth/rancher/callback",
			},
			"tokenExpirationSeconds":        3600,
			"refreshTokenExpirationSeconds": 7200,
		}
		assert.Equal(t, wantPayload, ops.updateCalledWith)
	})

	t.Run("updates resource and refreshes state on success", func(t *testing.T) {
		existing := &managementClient.OIDCClient{
			Resource: types.Resource{
				ID: "oidc-client-1",
			},
		}

		updated := &managementClient.OIDCClient{
			Description:                   "Updated description",
			TokenExpirationSeconds:        3600,
			RefreshTokenExpirationSeconds: 7200,
			RedirectURIs:                  []string{"http://127.0.0.1:5556/auth/rancher/callback"},
		}
		updated.ID = "oidc-client-1"

		ops := &fakeOIDCClientOperations{
			oidcClient:   existing,
			updateResult: updated,
		}
		// After Update, Read calls ByID again — return the updated client.
		// We reuse oidcClient for the Read; override with updated after first call.
		ops.oidcClient = updated

		d := newResourceData()
		err := resourceRancher2OIDCClientUpdate(d, newGetter(ops))

		require.NoError(t, err)
		assert.Equal(t, "oidc-client-1", d.Id())
		assert.Equal(t, "Updated description", d.Get("description"))
		assert.Equal(t, 3600, d.Get("token_expiration_seconds"))
		assert.Equal(t, 7200, d.Get("refresh_token_expiration_seconds"))
	})

	t.Run("returns error when ByID fails during update", func(t *testing.T) {
		ops := &fakeOIDCClientOperations{
			err: errors.New("connection refused"),
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientUpdate(d, newGetter(ops))

		require.Error(t, err)
	})

	t.Run("returns error when Update API call fails", func(t *testing.T) {
		existing := &managementClient.OIDCClient{}
		existing.ID = "oidc-client-1"

		ops := &fakeOIDCClientOperations{
			oidcClient: existing,
			updateErr:  errors.New("update rejected"),
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientUpdate(d, newGetter(ops))

		require.Error(t, err)
	})

	t.Run("clears ID when Read after Update returns not found", func(t *testing.T) {
		existing := &managementClient.OIDCClient{}
		existing.ID = "oidc-client-1"

		ops := &fakeOIDCClientOperations{
			updateResult: existing,
			// First ByID (fetch-before-update): success.
			// Second ByID (inside trailing Read): not found → ID cleared, no error.
			byIDQueue: byIDQueue{
				{client: existing},
				{err: &clientbase.APIError{StatusCode: http.StatusNotFound}},
			},
		}

		d := newResourceData()
		err := resourceRancher2OIDCClientUpdate(d, newGetter(ops))

		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

// byIDResponse holds a single canned response for a ByID call.
type byIDResponse struct {
	client *managementClient.OIDCClient
	err    error
}

// byIDQueue is a FIFO queue of byIDResponse values.
type byIDQueue []byIDResponse

func (q *byIDQueue) pop() (byIDResponse, bool) {
	if len(*q) == 0 {
		return byIDResponse{}, false
	}
	front := (*q)[0]
	*q = (*q)[1:]
	return front, true
}

// fakeOIDCClientOperations is an implementation of managementClient.OIDCClientOperations.
type fakeOIDCClientOperations struct {
	// ByID response — used when byIDQueue is empty.
	oidcClient *managementClient.OIDCClient
	err        error
	// byIDQueue lets tests supply different responses for successive ByID calls.
	// Each call dequeues the first entry; falls back to oidcClient/err when empty.
	byIDQueue byIDQueue

	// Create response
	createResult *managementClient.OIDCClient
	createErr    error

	// Update response
	updateResult *managementClient.OIDCClient
	updateErr    error
	// Captured Update arguments for assertion in tests.
	updateCalledWith interface{}
	updateCalledOn   *managementClient.OIDCClient

	// Delete response
	deleteErr error
}

func (f *fakeOIDCClientOperations) List(_ *types.ListOpts) (*managementClient.OIDCClientCollection, error) {
	return nil, nil
}

func (f *fakeOIDCClientOperations) ListAll(_ *types.ListOpts) (*managementClient.OIDCClientCollection, error) {
	return nil, nil
}

func (f *fakeOIDCClientOperations) Create(_ *managementClient.OIDCClient) (*managementClient.OIDCClient, error) {
	return f.createResult, f.createErr
}

func (f *fakeOIDCClientOperations) Update(existing *managementClient.OIDCClient, updates interface{}) (*managementClient.OIDCClient, error) {
	f.updateCalledOn = existing
	f.updateCalledWith = updates
	return f.updateResult, f.updateErr
}

func (f *fakeOIDCClientOperations) Replace(_ *managementClient.OIDCClient) (*managementClient.OIDCClient, error) {
	return nil, nil
}

func (f *fakeOIDCClientOperations) ByID(_ string) (*managementClient.OIDCClient, error) {
	if resp, ok := f.byIDQueue.pop(); ok {
		return resp.client, resp.err
	}
	return f.oidcClient, f.err
}

func (f *fakeOIDCClientOperations) Delete(_ *managementClient.OIDCClient) error {
	return f.deleteErr
}

// fakeManagementClientGetter satisfies the managementClientGetter interface using
// a pre-built *managementClient.Client with stubbed operations.
type fakeManagementClientGetter struct {
	client *managementClient.Client
}

func (f *fakeManagementClientGetter) ManagementClient() (*managementClient.Client, error) {
	return f.client, nil
}
