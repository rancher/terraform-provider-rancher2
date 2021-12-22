package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2UserType = "rancher2_user"
)

var (
	testAccRancher2User       string
	testAccRancher2UserUpdate string
)

func init() {
	testAccRancher2User = `
resource "` + testAccRancher2UserType + `" "foo" {
  name = "Terraform user acceptance test"
  username = "foo"
  password = "TestACC123456"
  enabled = true
}
`
	testAccRancher2UserUpdate = `
resource "` + testAccRancher2UserType + `" "foo" {
  name = "Terraform user acceptance test - Updated"
  username = "foo"
  password = "TestACC1234567"
  enabled = false
}
 `
}

func TestAccRancher2User_basic(t *testing.T) {
	var user *managementClient.User

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2User,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "TestACC123456"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "true"),
				),
			},
			{
				Config: testAccRancher2UserUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test - Updated"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "TestACC1234567"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "false"),
				),
			},
			{
				Config: testAccRancher2User,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "TestACC123456"),
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
			{
				Config: testAccRancher2User,
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
