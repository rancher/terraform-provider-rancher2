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
	testAccRancher2GlobalRoleType = "rancher2_global_role"
)

var (
	testAccRancher2GlobalRoleConfig       string
	testAccRancher2GlobalRoleUpdateConfig string
)

func init() {
	testAccRancher2GlobalRoleConfig = `
resource "` + testAccRancher2GlobalRoleType + `" "foo" {
  name = "foo"
  new_user_default = true
  description = "Terraform global role acceptance test"
  rules {
    api_groups = ["*"]
    resources = ["secrets"]
    verbs = ["` + policyRuleVerbCreate + `"]
  }
}
`
	testAccRancher2GlobalRoleUpdateConfig = `
resource "` + testAccRancher2GlobalRoleType + `" "foo" {
  name = "foo-updated"
  new_user_default = false
  description = "Terraform global role acceptance test - updated"
  rules {
    api_groups = ["*"]
    resources = ["secrets"]
    verbs = ["` + policyRuleVerbCreate + `", "` + policyRuleVerbGet + `"]
  }
}
 `
}

func TestAccRancher2GlobalRole_basic(t *testing.T) {
	var globalRole *managementClient.GlobalRole

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalRoleExists(testAccRancher2GlobalRoleType+".foo", globalRole),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "new_user_default", "true"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "description", "Terraform global role acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "rules.0.verbs.0", policyRuleVerbCreate),
				),
			},
			{
				Config: testAccRancher2GlobalRoleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalRoleExists(testAccRancher2GlobalRoleType+".foo", globalRole),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "new_user_default", "false"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "description", "Terraform global role acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "rules.0.verbs.1", policyRuleVerbGet),
				),
			},
			{
				Config: testAccRancher2GlobalRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalRoleExists(testAccRancher2GlobalRoleType+".foo", globalRole),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "new_user_default", "true"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "description", "Terraform global role acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleType+".foo", "rules.0.verbs.0", policyRuleVerbCreate),
				),
			},
		},
	})
}

func TestAccRancher2GlobalRole_disappears(t *testing.T) {
	var globalRole *managementClient.GlobalRole

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalRoleExists(testAccRancher2GlobalRoleType+".foo", globalRole),
					testAccRancher2GlobalRoleDisappears(globalRole),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2GlobalRoleDisappears(globalRole *managementClient.GlobalRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2GlobalRoleType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			globalRole, err = client.GlobalRole.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.GlobalRole.Delete(globalRole)
			if err != nil {
				return fmt.Errorf("Error removing global role: %s", err)
			}
		}
		return nil

	}
}

func testAccCheckRancher2GlobalRoleExists(n string, globalRole *managementClient.GlobalRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No global role ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundRole, err := client.GlobalRole.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Global role not found")
			}
			return err
		}

		globalRole = foundRole

		return nil
	}
}

func testAccCheckRancher2GlobalRoleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2GlobalRoleType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		for i := 0; i < 5; i++ {
			_, err = client.GlobalRole.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}
			time.Sleep(2 * time.Second)
		}

		return fmt.Errorf("Global role still exists")
	}
	return nil
}
