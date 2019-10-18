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
	testAccRancher2ProjectAlertGroupType = "rancher2_project_alert_group"
)

var (
	testAccRancher2ProjectAlertGroupProject        string
	testAccRancher2ProjectAlertGroupConfig         string
	testAccRancher2ProjectAlertGroupUpdateConfig   string
	testAccRancher2ProjectAlertGroupRecreateConfig string
)

func init() {
	testAccRancher2ProjectAlertGroupProject = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform project alert group acceptance test"
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

	testAccRancher2ProjectAlertGroupConfig = testAccRancher2ProjectAlertGroupProject + `
resource "rancher2_project_alert_group" "foo" {
  name = "foo"
  description = "Terraform project alert group acceptance test"
  project_id = "${rancher2_project.foo.id}"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
`

	testAccRancher2ProjectAlertGroupUpdateConfig = testAccRancher2ProjectAlertGroupProject + `
resource "rancher2_project_alert_group" "foo" {
  name = "foo"
  description = "Terraform project alert group acceptance test - updated"
  project_id = "${rancher2_project.foo.id}"
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}
 `

	testAccRancher2ProjectAlertGroupRecreateConfig = testAccRancher2ProjectAlertGroupProject + `
resource "rancher2_project_alert_group" "foo" {
  name = "foo"
  description = "Terraform project alert group acceptance test"
  project_id = "${rancher2_project.foo.id}"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
 `
}

func TestAccRancher2ProjectAlertGroup_basic(t *testing.T) {
	var ag *managementClient.ProjectAlertGroup

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectAlertGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ProjectAlertGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertGroupExists(testAccRancher2ProjectAlertGroupType+".foo", ag),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "description", "Terraform project alert group acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ProjectAlertGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertGroupExists(testAccRancher2ProjectAlertGroupType+".foo", ag),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "description", "Terraform project alert group acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "group_interval_seconds", "600"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "repeat_interval_seconds", "6000"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ProjectAlertGroupRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertGroupExists(testAccRancher2ProjectAlertGroupType+".foo", ag),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "description", "Terraform project alert group acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertGroupType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
		},
	})
}

func TestAccRancher2ProjectAlertGroup_disappears(t *testing.T) {
	var ag *managementClient.ProjectAlertGroup

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectAlertGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ProjectAlertGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertGroupExists(testAccRancher2ProjectAlertGroupType+".foo", ag),
					testAccRancher2ProjectAlertGroupDisappears(ag),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ProjectAlertGroupDisappears(ag *managementClient.ProjectAlertGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ProjectAlertGroupType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			ag, err = client.ProjectAlertGroup.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ProjectAlertGroup.Delete(ag)
			if err != nil {
				return fmt.Errorf("Error removing Project Alert Group: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    projectAlertGroupStateRefreshFunc(client, rs.Primary.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for project alert group (%s) to be removed: %s", rs.Primary.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ProjectAlertGroupExists(n string, ag *managementClient.ProjectAlertGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project alert group ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundAg, err := client.ProjectAlertGroup.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Project Alert Group not found")
			}
			return err
		}

		ag = foundAg

		return nil
	}
}

func testAccCheckRancher2ProjectAlertGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ProjectAlertGroupType {
			continue
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		obj, err := client.ProjectAlertGroup.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		if obj.Removed != "" {
			return nil
		}
		return fmt.Errorf("Project Alert Group still exists")
	}
	return nil
}
