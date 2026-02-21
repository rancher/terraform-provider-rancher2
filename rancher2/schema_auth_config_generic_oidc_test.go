package rancher2

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAuthConfigGenericOIDCResourceValidateNoData(t *testing.T) {
	r := resourceRancher2AuthConfigGenericOIDC()
	d := terraform.NewResourceConfigRaw(map[string]any{})
	warns, errs := r.Validate(d)

	assert.Empty(t, warns)

	assert.Len(t, errs, 4)
	joined := errors.Join(errs...)
	assert.ErrorContains(t, joined, "\"issuer\": required field is not set")
	assert.ErrorContains(t, joined, "\"client_id\": required field is not set")
	assert.ErrorContains(t, joined, "\"client_secret\": required field is not set")
	assert.ErrorContains(t, joined, "\"rancher_url\": required field is not set")
}

func TestAuthConfigGenericOIDCResourceEndSessionEndpointValidation(t *testing.T) {
	r := resourceRancher2AuthConfigGenericOIDC()
	d := terraform.NewResourceConfigRaw(map[string]any{
		"client_id":            "client-id",
		"client_secret":        "client-secret",
		"issuer":               "https://issuer.example.com",
		"rancher_url":          "https://rancher.example.com/verify-auth",
		"end_session_endpoint": "test",
	})
	warns, errs := r.Validate(d)

	assert.Empty(t, warns)
	assert.ErrorContains(t, errors.Join(errs...), `expected "end_session_endpoint" to have a host`)
}
