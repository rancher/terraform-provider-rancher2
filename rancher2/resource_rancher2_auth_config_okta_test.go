package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2AuthConfigOKTAType   = "rancher2_auth_config_okta"
	testAccRancher2AuthConfigOKTAConfig = `
resource "` + testAccRancher2AuthConfigOKTAType + `" "okta" {
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
	testAccRancher2AuthConfigOKTAUpdateConfig = `
resource "` + testAccRancher2AuthConfigOKTAType + `" "okta" {
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

func TestAccRancher2AuthConfigOKTA_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigOKTADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigOKTAConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "name", AuthConfigOKTAName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "user_name_field", "sAMAccountName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2AuthConfigOKTAUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "name", AuthConfigOKTAName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "user_name_field", "sAMAccountName-updated"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "rancher_api_host", "https://RANCHER-UPDATED"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "sp_key", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "idp_metadata_content", "YYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2AuthConfigOKTAConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "name", AuthConfigOKTAName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "user_name_field", "sAMAccountName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigOKTA_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigOKTADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigOKTAConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigOKTAType+"."+AuthConfigOKTAName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigOKTAType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigOKTADestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigOKTAType {
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
