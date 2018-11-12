package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2ProjectType   = "rancher2_project"
	testAccRancher2ProjectConfig = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "local"
  description = "Terraform project acceptance test"
}
`

	testAccRancher2ProjectUpdateConfig = `
resource "rancher2_project" "foo" {
  name = "foo-updated"
  cluster_id = "local"
  description = "Terraform project acceptance test - updated"
}
 `

	testAccRancher2ProjectRecreateConfig = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "local"
  description = "Terraform project acceptance test"
}
 `
)

func TestAccRancher2Project_basic(t *testing.T) {
	var project *managementClient.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectExists(testAccRancher2ProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "description", "Terraform project acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "cluster_id", "local"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ProjectUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectExists(testAccRancher2ProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "description", "Terraform project acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "cluster_id", "local"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ProjectRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectExists(testAccRancher2ProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "description", "Terraform project acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "cluster_id", "local"),
				),
			},
		},
	})
}

func TestAccRancher2Project_disappears(t *testing.T) {
	var project *managementClient.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectExists(testAccRancher2ProjectType+".foo", project),
					testAccRancher2ProjectDisappears(project),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ProjectDisappears(pro *managementClient.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ProjectType {
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
				Refresh:    projectStateRefreshFunc(client, pro.ID),
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

func testAccCheckRancher2ProjectExists(n string, pro *managementClient.Project) resource.TestCheckFunc {
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

func testAccCheckRancher2ProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ProjectType {
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
