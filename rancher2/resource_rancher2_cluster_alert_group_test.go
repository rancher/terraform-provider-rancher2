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
	testAccRancher2ClusterAlertGroupType = "rancher2_cluster_alert_group"
)

var (
	testAccRancher2ClusterAlertGroup             string
	testAccRancher2ClusterAlertGroupUpdate       string
	testAccRancher2ClusterAlertGroupConfig       string
	testAccRancher2ClusterAlertGroupUpdateConfig string
)

func init() {
	testAccRancher2ClusterAlertGroup = `
resource "` + testAccRancher2ClusterAlertGroupType + `" "foo" {
  name = "foo"
  description = "Terraform cluster alert group acceptance test"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
`
	testAccRancher2ClusterAlertGroupUpdate = `
resource "` + testAccRancher2ClusterAlertGroupType + `" "foo" {
  name = "foo"
  description = "Terraform cluster alert group acceptance test - updated"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}
 `
	testAccRancher2ClusterAlertGroupConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ClusterAlertGroup
	testAccRancher2ClusterAlertGroupUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ClusterAlertGroupUpdate
}

func TestAccRancher2ClusterAlertGroup_basic(t *testing.T) {
	var ag *managementClient.ClusterAlertGroup

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterAlertGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterAlertGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterAlertGroupExists(testAccRancher2ClusterAlertGroupType+".foo", ag),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "description", "Terraform cluster alert group acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
			{
				Config: testAccRancher2ClusterAlertGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterAlertGroupExists(testAccRancher2ClusterAlertGroupType+".foo", ag),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "description", "Terraform cluster alert group acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "group_interval_seconds", "600"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "repeat_interval_seconds", "6000"),
				),
			},
			{
				Config: testAccRancher2ClusterAlertGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterAlertGroupExists(testAccRancher2ClusterAlertGroupType+".foo", ag),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "description", "Terraform cluster alert group acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertGroupType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
		},
	})
}

func TestAccRancher2ClusterAlertGroup_disappears(t *testing.T) {
	var ag *managementClient.ClusterAlertGroup

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterAlertGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterAlertGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterAlertGroupExists(testAccRancher2ClusterAlertGroupType+".foo", ag),
					testAccRancher2ClusterAlertGroupDisappears(ag),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterAlertGroupDisappears(ag *managementClient.ClusterAlertGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterAlertGroupType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			ag, err = client.ClusterAlertGroup.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ClusterAlertGroup.Delete(ag)
			if err != nil {
				return fmt.Errorf("Error removing Cluster Alert Group: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    clusterAlertGroupStateRefreshFunc(client, rs.Primary.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for cluster alert group (%s) to be removed: %s", rs.Primary.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ClusterAlertGroupExists(n string, ag *managementClient.ClusterAlertGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster alert group ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundAg, err := client.ClusterAlertGroup.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cluster Alert Group not found")
			}
			return err
		}

		ag = foundAg

		return nil
	}
}

func testAccCheckRancher2ClusterAlertGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterAlertGroupType {
			continue
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ClusterAlertGroup.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cluster Alert Group still exists")
	}
	return nil
}
