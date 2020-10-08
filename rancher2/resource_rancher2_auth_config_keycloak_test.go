package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2AuthConfigKeyCloakType   = "rancher2_auth_config_keycloak"
	testAccRancher2AuthConfigKeyCloakConfig = `
resource "` + testAccRancher2AuthConfigKeyCloakType + `" "keycloak" {
  display_name_field = "displayName"
  groups_field = "memberOf"
  uid_field = "distinguishedName"
  user_name_field = "sAMAccountName"
  idp_metadata_content = "XXXXXXXX"
  rancher_api_host = "https://RANCHER"
  sp_cert = "XXXXXX"
  sp_key = "XXXXXXXX"
}
`
	testAccRancher2AuthConfigKeyCloakUpdateConfig = `
resource "` + testAccRancher2AuthConfigKeyCloakType + `" "keycloak" {
  display_name_field = "displayName"
  groups_field = "memberOf"
  uid_field = "distinguishedName"
  user_name_field = "sAMAccountName-updated"
  idp_metadata_content = "YYYYYYYY"
  rancher_api_host = "https://RANCHER-UPDATED"
  sp_cert = "XXXXXX"
  sp_key = "YYYYYYYY"
}
 `
)

func TestAccRancher2AuthConfigKeyCloak_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigKeyCloakDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigKeyCloakConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "name", AuthConfigKeyCloakName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "user_name_field", "sAMAccountName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2AuthConfigKeyCloakUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "name", AuthConfigKeyCloakName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "user_name_field", "sAMAccountName-updated"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "rancher_api_host", "https://RANCHER-UPDATED"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "sp_key", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "idp_metadata_content", "YYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2AuthConfigKeyCloakConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "name", AuthConfigKeyCloakName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "user_name_field", "sAMAccountName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigKeyCloak_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigKeyCloakDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigKeyCloakConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigKeyCloakType+"."+AuthConfigKeyCloakName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigKeyCloakType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigKeyCloakDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigKeyCloakType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		auth, err := client.AuthConfig.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}

		if auth.Enabled == true {
			err = client.Post(auth.Actions["disable"], nil, nil)
			if err != nil {
				return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", auth.ID, err)
			}
		}
		return nil
	}
	return nil
}
