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
	testAccCattleProjectType   = "cattle_project"
	testAccCattleProjectConfig = `
resource "cattle_project" "foo" {
  name = "foo"
  cluster_id = "local"
  description = "Foo project test"
}
`

	testAccCattleProjectUpdateConfig = `
resource "cattle_project" "foo" {
  name = "foo-updated"
  cluster_id = "local"
  description = "Foo project test - updated"
}
 `

	testAccCattleProjectRecreateConfig = `
resource "cattle_project" "foo" {
  name = "foo"
  cluster_id = "local"
  description = "Foo project test"
}
 `
)

func TestAccCattleProject_basic(t *testing.T) {
	var project *managementClient.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCattleProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCattleProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleProjectExists(testAccCattleProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "description", "Foo project test"),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "cluster_id", "local"),
				),
			},
			resource.TestStep{
				Config: testAccCattleProjectUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleProjectExists(testAccCattleProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "description", "Foo project test - updated"),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "cluster_id", "local"),
				),
			},
			resource.TestStep{
				Config: testAccCattleProjectRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleProjectExists(testAccCattleProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "description", "Foo project test"),
					resource.TestCheckResourceAttr(testAccCattleProjectType+".foo", "cluster_id", "local"),
				),
			},
		},
	})
}

func TestAccCattleProject_disappears(t *testing.T) {
	var project *managementClient.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCattleProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCattleProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleProjectExists(testAccCattleProjectType+".foo", project),
					testAccCattleProjectDisappears(project),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCattleProjectDisappears(pro *managementClient.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccCattleProjectType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.Project.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.Project.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Project: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    ProjectStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for project (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckCattleProjectExists(n string, pro *managementClient.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.Project.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Project not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckCattleProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCattleProjectType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.Project.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Project still exists")
	}
	return nil
}
