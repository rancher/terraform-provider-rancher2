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
	testAccRancher2RoleTemplateType = "rancher2_role_template"
)

var (
	testAccRancher2RoleTemplateConfig       string
	testAccRancher2RoleTemplateUpdateConfig string
)

func init() {
	testAccRancher2RoleTemplateConfig = `
resource "` + testAccRancher2RoleTemplateType + `" "foo" {
  name = "foo"
  context = "` + roleTemplateContextCluster + `"
  default_role = true
  description = "Terraform role template acceptance test"
  rules {
    api_groups = ["*"]
    resources = ["secrets"]
    verbs = ["` + policyRuleVerbCreate + `"]
  }
}
`
	testAccRancher2RoleTemplateUpdateConfig = `
resource "` + testAccRancher2RoleTemplateType + `" "foo" {
  name = "foo-updated"
  context = "` + roleTemplateContextProject + `"
  default_role = false
  description = "Terraform role template acceptance test - updated"
  rules {
    api_groups = ["*"]
    resources = ["secrets"]
    verbs = ["` + policyRuleVerbCreate + `", "` + policyRuleVerbGet + `"]
  }
}
 `
}

func TestAccRancher2RoleTemplate_basic(t *testing.T) {
	var roleTemplate *managementClient.RoleTemplate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RoleTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2RoleTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RoleTemplateExists(testAccRancher2RoleTemplateType+".foo", roleTemplate),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "context", roleTemplateContextCluster),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "default_role", "true"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "description", "Terraform role template acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "rules.0.verbs.0", policyRuleVerbCreate),
				),
			},
			{
				Config: testAccRancher2RoleTemplateUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RoleTemplateExists(testAccRancher2RoleTemplateType+".foo", roleTemplate),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "context", roleTemplateContextProject),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "default_role", "false"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "description", "Terraform role template acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "rules.0.verbs.1", policyRuleVerbGet),
				),
			},
			{
				Config: testAccRancher2RoleTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RoleTemplateExists(testAccRancher2RoleTemplateType+".foo", roleTemplate),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "context", roleTemplateContextCluster),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "default_role", "true"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "description", "Terraform role template acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2RoleTemplateType+".foo", "rules.0.verbs.0", policyRuleVerbCreate),
				),
			},
		},
	})
}

func TestAccRancher2RoleTemplate_disappears(t *testing.T) {
	var roleTemplate *managementClient.RoleTemplate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RoleTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2RoleTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RoleTemplateExists(testAccRancher2RoleTemplateType+".foo", roleTemplate),
					testAccRancher2RoleTemplateDisappears(roleTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2RoleTemplateDisappears(roleTemplate *managementClient.RoleTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2RoleTemplateType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			roleTemplate, err = client.RoleTemplate.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.RoleTemplate.Delete(roleTemplate)
			if err != nil {
				return fmt.Errorf("Error removing role template: %s", err)
			}
		}
		return nil

	}
}

func testAccCheckRancher2RoleTemplateExists(n string, roleTemplate *managementClient.RoleTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No role template ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundRole, err := client.RoleTemplate.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Role Template not found")
			}
			return err
		}

		roleTemplate = foundRole

		return nil
	}
}

func testAccCheckRancher2RoleTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2RoleTemplateType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		for i := 0; i < 5; i++ {
			_, err = client.RoleTemplate.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}
			time.Sleep(2 * time.Second)
		}

		return fmt.Errorf("Role template still exists")
	}
	return nil
}
