package rancher2

import (
	"errors"
	"net/http"
	"testing"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

func TestIsFeatureEnabled(t *testing.T) {
	lookupErr := errors.New("lookup failed")

	cases := []struct {
		name           string
		feature        *managementClient.Feature
		lookupErr      error
		wantEnabled    bool
		wantErr        bool
		wantWrappedErr bool
	}{
		{
			name:        "enabled",
			feature:     &managementClient.Feature{Value: new(true)},
			wantEnabled: true,
		},
		{
			name:    "disabled",
			feature: &managementClient.Feature{Value: new(false)},
		},
		{
			name:    "nil value",
			feature: &managementClient.Feature{},
		},
		{
			name:      "not found",
			lookupErr: &clientbase.APIError{StatusCode: http.StatusNotFound},
		},
		{
			name:      "forbidden",
			lookupErr: &clientbase.APIError{StatusCode: http.StatusForbidden},
		},
		{
			name:           "other error",
			lookupErr:      lookupErr,
			wantErr:        true,
			wantWrappedErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			features := &fakeFeatureOperations{feature: tc.feature, err: tc.lookupErr}
			gotEnabled, err := isFeatureEnabled(features, "test-feature")

			assert.Equal(t, tc.wantEnabled, gotEnabled)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "getting feature test-feature")
				if tc.wantWrappedErr {
					assert.ErrorIs(t, err, lookupErr)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, "test-feature", features.requestedID)
		})
	}
}

type fakeFeatureOperations struct {
	managementClient.FeatureOperations
	feature     *managementClient.Feature
	err         error
	requestedID string
}

func (f *fakeFeatureOperations) ByID(id string) (*managementClient.Feature, error) {
	f.requestedID = id
	return f.feature, f.err
}

func TestCheckAuthConfigCanBeEnabled(t *testing.T) {
	lookupErr := errors.New("lookup failed")
	alreadyEnabled := "enabling Auth Config okta: Auth Config github is already enabled, and the multiple-auth-providers feature is not enabled in Rancher"

	cases := []struct {
		name      string
		feature   *managementClient.Feature
		lookupErr error
		auths     []managementClient.AuthConfig
		listErr   error
		wantErr   string
	}{
		{
			name:    "feature enabled",
			feature: &managementClient.Feature{Value: new(true)},
			auths:   []managementClient.AuthConfig{{Resource: types.Resource{ID: "github"}, Enabled: true}},
		},
		{
			name:    "feature disabled and only local enabled",
			feature: &managementClient.Feature{Value: new(false)},
			auths: []managementClient.AuthConfig{
				{Resource: types.Resource{ID: "local"}, Enabled: true},
				{Resource: types.Resource{ID: "github"}},
			},
		},
		{
			name:    "feature disabled and same config enabled",
			feature: &managementClient.Feature{Value: new(false)},
			auths:   []managementClient.AuthConfig{{Resource: types.Resource{ID: "okta"}, Enabled: true}},
		},
		{
			name:    "feature disabled and another config enabled",
			feature: &managementClient.Feature{Value: new(false)},
			auths: []managementClient.AuthConfig{
				{Resource: types.Resource{ID: "local"}, Enabled: true},
				{Resource: types.Resource{ID: "github"}, Enabled: true},
			},
			wantErr: alreadyEnabled,
		},
		{
			name:      "feature not found and another config enabled",
			lookupErr: &clientbase.APIError{StatusCode: http.StatusNotFound},
			auths:     []managementClient.AuthConfig{{Resource: types.Resource{ID: "github"}, Enabled: true}},
			wantErr:   alreadyEnabled,
		},
		{
			name:      "feature lookup error",
			lookupErr: lookupErr,
			wantErr:   "getting feature multiple-auth-providers: lookup failed",
		},
		{
			name:    "list error",
			feature: &managementClient.Feature{Value: new(false)},
			listErr: lookupErr,
			wantErr: "listing Auth Configs: lookup failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			features := &fakeFeatureOperations{feature: tc.feature, err: tc.lookupErr}
			authConfigs := &fakeAuthConfigOperations{auths: tc.auths, err: tc.listErr}

			err := checkAuthConfigCanBeEnabled(features, authConfigs, "okta")

			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErr)
			}
			assert.Equal(t, multipleAuthProvidersFeature, features.requestedID)
		})
	}
}

type fakeAuthConfigOperations struct {
	managementClient.AuthConfigOperations
	auths []managementClient.AuthConfig
	err   error
}

func (f *fakeAuthConfigOperations) List(opts *types.ListOpts) (*managementClient.AuthConfigCollection, error) {
	if f.err != nil {
		return nil, f.err
	}

	return &managementClient.AuthConfigCollection{Data: f.auths}, nil
}
