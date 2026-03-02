package rancher2

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2AuthConfigGenericOIDCName = "genericoidc"
	testAccRancher2AuthConfigGenericOIDCType = "rancher2_auth_config_generic_oidc"
)

var (
	testAccProviders                                 map[string]terraform.ResourceProvider
	testAccProvider                                  *schema.Provider
	testAccRancher2AuthConfigGenericOIDCClientID     = os.Getenv("RANCHER_OIDC_CLIENT_ID")
	testAccRancher2AuthConfigGenericOIDCClientSecret = os.Getenv("RANCHER_OIDC_CLIENT_SECRET")
	testAccRancher2AuthConfigGenericOIDCIssuerURL    = os.Getenv("RANCHER_OIDC_ISSUER_URL")
	testAccRancher2AuthConfigGenericOIDCRancherURL   = os.Getenv("RANCHER_URL")
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"rancher2": testAccProvider,
	}
}

func testAccCheckRancher2AuthConfigGenericOIDCExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Auth Config Generic OIDC ID is set")
		}

		// Reading the auth config can be slow, add a delay.
		time.Sleep(2 * time.Second)

		return nil
	}
}

func testAccCheckRancher2AuthConfigGenericOIDCConfig() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		auth, err := client.AuthConfig.ByID(AuthConfigGenericOIDCName)
		if err != nil {
			return fmt.Errorf("Failed to get Auth Config %s: %s", AuthConfigGenericOIDCName, err)
		}

		if auth.Enabled != true {
			return fmt.Errorf("Auth Config %s is not enabled", AuthConfigGenericOIDCName)
		}

		return nil
	}
}

func testAccCheckRancher2AuthConfigGenericOIDCDisabled(s *terraform.State) error {
	client, err := testAccProvider.Meta().(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigGenericOIDCName)
	if err != nil {
		if IsNotFound(err) {
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		return fmt.Errorf("Auth Config %s is still enabled", AuthConfigGenericOIDCName)
	}

	return nil
}

func testAccRancher2AuthConfigGenericOIDCConfig() string {
	return fmt.Sprintf(`
resource "rancher2_auth_config_generic_oidc" "genericoidc" {
  client_id            = "%s"
  client_secret        = "%s"
  issuer               = "%s"
  rancher_url          = "%s/verify-auth"
  enabled              = true
  scopes               = "openid profile email"
  groups_field         = "groups"
  group_search_enabled = true
}
`, testAccRancher2AuthConfigGenericOIDCClientID, testAccRancher2AuthConfigGenericOIDCClientSecret, testAccRancher2AuthConfigGenericOIDCIssuerURL, testAccRancher2AuthConfigGenericOIDCRancherURL)
}

func testAccRancher2AuthConfigGenericOIDCUpdateConfig() string {
	return fmt.Sprintf(`
resource "rancher2_auth_config_generic_oidc" "genericoidc" {
  client_id            = "%s"
  client_secret        = "%s"
  issuer               = "%s"
  rancher_url          = "%s/verify-auth"
  enabled              = true
  scopes               = "openid profile"
  groups_field         = "group"
  group_search_enabled = false
}
`, testAccRancher2AuthConfigGenericOIDCClientID, testAccRancher2AuthConfigGenericOIDCClientSecret, testAccRancher2AuthConfigGenericOIDCIssuerURL, testAccRancher2AuthConfigGenericOIDCRancherURL)
}
