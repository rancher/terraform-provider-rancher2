package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2AuthConfigActiveDirectoryType   = "rancher2_auth_config_activedirectory"
	testAccRancher2AuthConfigActiveDirectoryConfig = `
resource "` + testAccRancher2AuthConfigActiveDirectoryType + `" "activedirectory" {
  servers = ["ad.test.local"]
  service_account_username = "XXXXXX"
  service_account_password = "XXXXXXXXX"
  user_search_base = "dc=test,dc=local"
  port = 389
  default_login_domain = "test"
  enabled = false
  test_username = "test"
  test_password = "test"
}
`

	testAccRancher2AuthConfigActiveDirectoryUpdateConfig = `
resource "` + testAccRancher2AuthConfigActiveDirectoryType + `" "activedirectory" {
  servers = ["ad.test.local"]
  service_account_username = "XXXXXX"
  service_account_password = "XXXXXXXXX"
  user_search_base = "dc=users,dc=test,dc=local"
  port = 389
  default_login_domain = "test-updated"
  enabled = false
  test_username = "test"
  test_password = "test"
}
 `
)

func TestAccRancher2AuthConfigActiveDirectory_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigActiveDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigActiveDirectoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "name", AuthConfigActiveDirectoryName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "user_search_base", "dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "default_login_domain", "test"),
				),
			},
			{
				Config: testAccRancher2AuthConfigActiveDirectoryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "name", AuthConfigActiveDirectoryName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "user_search_base", "dc=users,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "default_login_domain", "test-updated"),
				),
			},
			{
				Config: testAccRancher2AuthConfigActiveDirectoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "name", AuthConfigActiveDirectoryName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "user_search_base", "dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, "default_login_domain", "test"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigActiveDirectory_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigActiveDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AuthConfigActiveDirectoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigActiveDirectoryType+"."+AuthConfigActiveDirectoryName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigActiveDirectoryType),
				),
			},
		},
	})
}

func testAccRancher2AuthConfigDisappears(auth *managementClient.AuthConfig, objType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != objType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			auth, err = client.AuthConfig.ByID(rs.Primary.ID)
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
}

func testAccCheckRancher2AuthConfigExists(n string, auth *managementClient.AuthConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Auth Config ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundReg, err := client.AuthConfig.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Auth Config %s not found", AuthConfigActiveDirectoryName)
			}
			return err
		}

		auth = foundReg

		return nil
	}
}

func testAccCheckRancher2AuthConfigActiveDirectoryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigActiveDirectoryType {
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
