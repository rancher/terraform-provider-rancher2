package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2AuthConfigPingType   = "rancher2_auth_config_ping"
	testAccRancher2AuthConfigPingConfig = `
resource "` + testAccRancher2AuthConfigPingType + `" "ping" {
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
	testAccRancher2AuthConfigPingUpdateConfig = `
resource "` + testAccRancher2AuthConfigPingType + `" "ping" {
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

func TestAccRancher2AuthConfigPing_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigPingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigPingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "name", AuthConfigPingName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "user_name_field", "sAMAccountName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2AuthConfigPingUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "name", AuthConfigPingName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "user_name_field", "sAMAccountName-updated"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "rancher_api_host", "https://RANCHER-UPDATED"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "sp_key", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "idp_metadata_content", "YYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2AuthConfigPingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "name", AuthConfigPingName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "user_name_field", "sAMAccountName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigPing_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigPingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigPingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigPingType+"."+AuthConfigPingName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigPingType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigPingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigPingType {
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
