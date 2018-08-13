package cattle

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccCattleProjectRoleTemplateBindingType   = "cattle_project_role_template_binding"
	testAccCattleProjectRoleTemplateBindingConfig = `
resource "cattle_project_role_template_binding" "foo" {
  name = "foo"
  project_id = "local:p-2lk7g"
  role_template_id = "project-member"
}
`

	testAccCattleProjectRoleTemplateBindingUpdateConfig = `
resource "cattle_project_role_template_binding" "foo" {
  name = "foo-updated"
  project_id = "local:p-2lk7g"
  role_template_id = "project-member"
  user_id = "u-q2wg7"
}
 `

	testAccCattleProjectRoleTemplateBindingRecreateConfig = `
resource "cattle_project_role_template_binding" "foo" {
  name = "foo"
  project_id = "local:p-2lk7g"
  role_template_id = "project-member"
}
 `
)

func TestAccCattleProjectRoleTemplateBinding_basic(t *testing.T) {
	var projectRole *managementClient.ProjectRoleTemplateBinding

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCattleProjectRoleTemplateBindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCattleProjectRoleTemplateBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleProjectRoleTemplateBindingExists(testAccCattleProjectRoleTemplateBindingType+".foo", projectRole),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "project_id", "local:p-2lk7g"),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "role_template_id", "project-member"),
				),
			},
			resource.TestStep{
				Config: testAccCattleProjectRoleTemplateBindingUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleProjectRoleTemplateBindingExists(testAccCattleProjectRoleTemplateBindingType+".foo", projectRole),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "project_id", "local:p-2lk7g"),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "role_template_id", "project-member"),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "user_id", "local"),
				),
			},
			resource.TestStep{
				Config: testAccCattleProjectRoleTemplateBindingRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleProjectRoleTemplateBindingExists(testAccCattleProjectRoleTemplateBindingType+".foo", projectRole),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "project_id", "Foo project test"),
					resource.TestCheckResourceAttr(testAccCattleProjectRoleTemplateBindingType+".foo", "role_template_id", "project-member"),
				),
			},
		},
	})
}

func TestAccCattleProjectRoleTemplateBinding_disappears(t *testing.T) {
	var projectRole *managementClient.ProjectRoleTemplateBinding

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCattleProjectRoleTemplateBindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCattleProjectRoleTemplateBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleProjectRoleTemplateBindingExists(testAccCattleProjectRoleTemplateBindingType+".foo", projectRole),
					testAccCattleProjectRoleTemplateBindingDisappears(projectRole),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCattleProjectRoleTemplateBindingDisappears(pro *managementClient.ProjectRoleTemplateBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccCattleProjectRoleTemplateBindingType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.ProjectRoleTemplateBinding.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ProjectRoleTemplateBinding.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Project Role Template Binding: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    projectRoleTemplateBindingStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Project Role Template Binding (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckCattleProjectRoleTemplateBindingExists(n string, pro *managementClient.ProjectRoleTemplateBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project Role Template Binding ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.ProjectRoleTemplateBinding.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Project Role Template Binding not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckCattleProjectRoleTemplateBindingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCattleProjectRoleTemplateBindingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ProjectRoleTemplateBinding.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Project Role Template Binding still exists")
	}
	return nil
}
