package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestAuthConfigCognitoResourceDefaults(t *testing.T) {
	r := resourceRancher2AuthConfigCognito()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]any{
		"client_id":     "client-id",
		"client_secret": "client-secret",
		"issuer":        "https://issuer.example.com",
		"rancher_url":   "https://rancher.example.com/verify-auth",
	})

	if got := d.Get("name_claim"); got != "cognito:username" {
		t.Fatalf("unexpected default for name_claim: got %v, want %q", got, "cognito:username")
	}

	if got := d.Get("group_search_enabled"); got != false {
		t.Fatalf("unexpected default for group_search_enabled: got %v, want %v", got, false)
	}
}
