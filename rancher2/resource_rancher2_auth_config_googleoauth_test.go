package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2AuthConfigGoogleOauthType   = "rancher2_auth_config_googleoauth"
	testAccRancher2AuthConfigGoogleOauthConfig = `
resource "rancher2_auth_config_googleoauth" "googleoauth" {
  admin_email = "XXXXXX"
  hostname = "XXXXXX"
  oauth_credential = "XXXXXXXX"
  service_account_credential = "XXXXXXXX"
  nested_group_membership_enabled = true
}
`

	testAccRancher2AuthConfigGoogleOauthUpdateConfig = `
resource "rancher2_auth_config_googleoauth" "googleoauth" {
  admin_email = "YYYYYY"
  hostname = "YYYYYY"
  oauth_credential = "YYYYYYYY"
  service_account_credential = "YYYYYYYY"
  nested_group_membership_enabled = false
}
 `

	testAccRancher2AuthConfigGoogleOauthRecreateConfig = `
resource "rancher2_auth_config_googleoauth" "googleoauth" {
  admin_email = "XXXXXX"
  hostname = "XXXXXX"
  oauth_credential = "XXXXXXXX"
  service_account_credential = "XXXXXXXX"
  nested_group_membership_enabled = true
}
 `
)

func TestAccRancher2AuthConfigGoogleOauth_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigGoogleOauthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigGoogleOauthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "name", AuthConfigGoogleOauthName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "admin_email", "XXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "nested_group_membership_enabled", "true"),
				),
			},
			{
				Config: testAccRancher2AuthConfigGoogleOauthUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "name", AuthConfigGoogleOauthName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "admin_email", "YYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "nested_group_membership_enabled", "false"),
				),
			},
			{
				Config: testAccRancher2AuthConfigGoogleOauthRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "name", AuthConfigGoogleOauthName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "admin_email", "XXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, "nested_group_membership_enabled", "true"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigGoogleOauth_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigGoogleOauthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigGoogleOauthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigGoogleOauthType+"."+AuthConfigGoogleOauthName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigGoogleOauthType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigGoogleOauthDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigGoogleOauthType {
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
