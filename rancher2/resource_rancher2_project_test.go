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
	testAccRancher2ProjectType = "rancher2_project"
)

var (
	testAccRancher2ProjectConfig       string
	testAccRancher2ProjectUpdateConfig string
	testAccRancher2Project             string
	testAccRancher2ProjectUpdate       string
)

func init() {

	testAccRancher2Project = `
resource "` + testAccRancher2ProjectType + `" "foo" {
  name = "foo"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  description = "Terraform project acceptance test"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "500m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
}
`
	testAccRancher2ProjectUpdate = `
resource "` + testAccRancher2ProjectType + `" "foo" {
  name = "foo-updated"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  description = "Terraform project acceptance test - updated"
  resource_quota {
    project_limit {
      limits_cpu = "1000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "700m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
}
 `
}

func TestAccRancher2Project_basic(t *testing.T) {
	var project *managementClient.Project

	testAccRancher2ProjectConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2Project
	testAccRancher2ProjectUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ProjectUpdate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectExists(testAccRancher2ProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "description", "Terraform project acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "wait_for_cluster", "false"),
				),
			},
			{
				Config: testAccRancher2ProjectUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectExists(testAccRancher2ProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "description", "Terraform project acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "wait_for_cluster", "false"),
				),
			},
			{
				Config: testAccRancher2ProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectExists(testAccRancher2ProjectType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "description", "Terraform project acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2ProjectType+".foo", "wait_for_cluster", "false"),
				),
			},
		},
	})
}

func TestAccRancher2Project_disappears(t *testing.T) {
	var project *managementClient.Project

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectDestroy,
		Steps: []resource.TestStep{
			{
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
