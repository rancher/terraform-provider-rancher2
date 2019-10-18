package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2UserType = "rancher2_user"
)

var (
	testAccRancher2UserConfig         string
	testAccRancher2UserUpdateConfig   string
	testAccRancher2UserRecreateConfig string
)

func init() {
	testAccRancher2UserConfig = `
resource "rancher2_user" "foo" {
  name = "Terraform user acceptance test"
  username = "foo"
  password = "changeme"
  enabled = true
}
`

	testAccRancher2UserUpdateConfig = `
resource "rancher2_user" "foo" {
  name = "Terraform user acceptance test - Updated"
  username = "foo"
  password = "changeme2"
  enabled = false
}
 `

	testAccRancher2UserRecreateConfig = `
resource "rancher2_user" "foo" {
  name = "Terraform user acceptance test"
  username = "foo"
  password = "changeme"
  enabled = true
}
 `
}

func TestAccRancher2User_basic(t *testing.T) {
	var user *managementClient.User

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2UserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2UserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "changeme"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "true"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2UserUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test - Updated"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "changeme2"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "false"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2UserRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "changeme"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccRancher2User_disappears(t *testing.T) {
	var user *managementClient.User

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2UserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2UserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					testAccRancher2UserDisappears(user),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2UserDisappears(pro *managementClient.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2UserType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.User.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.User.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing User: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    userStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for User (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2UserExists(n string, pro *managementClient.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.User.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("User not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2UserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2UserType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.User.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("User still exists")
	}
	return nil
}
