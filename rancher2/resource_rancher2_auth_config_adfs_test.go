package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2AuthConfigADFSType   = "rancher2_auth_config_adfs"
	testAccRancher2AuthConfigADFSConfig = `
resource "` + testAccRancher2AuthConfigADFSType + `" "adfs" {
  display_name_field = "givenname"
  groups_field = "Group"
  uid_field = "upn"
  user_name_field = "name"
  idp_metadata_content = "XXXXXXXX"
  rancher_api_host = "https://RANCHER"
  sp_cert = "XXXXXX"
  sp_key = "XXXXXXXX"
}
`

	testAccRancher2AuthConfigADFSUpdateConfig = `
resource "` + testAccRancher2AuthConfigADFSType + `" "adfs" {
  display_name_field = "givenname"
  groups_field = "Group"
  uid_field = "upn"
  user_name_field = "name-updated"
  idp_metadata_content = "YYYYYYYY"
  rancher_api_host = "https://RANCHER-UPDATED"
  sp_cert = "XXXXXX"
  sp_key = "YYYYYYYY"
}
 `
)

func TestAccRancher2AuthConfigADFS_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigADFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigADFSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "name", AuthConfigADFSName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "user_name_field", "name"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2AuthConfigADFSUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "name", AuthConfigADFSName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "user_name_field", "name-updated"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "rancher_api_host", "https://RANCHER-UPDATED"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "sp_key", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "idp_metadata_content", "YYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2AuthConfigADFSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "name", AuthConfigADFSName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "user_name_field", "name"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigADFS_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigADFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigADFSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigADFSType+"."+AuthConfigADFSName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigADFSType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigADFSDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigADFSType {
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
