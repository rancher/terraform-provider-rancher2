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
	testAccRancher2ClusterLoggingType = "rancher2_cluster_logging"
)

var (
	testAccRancher2ClusterLoggingConfigSyslog         string
	testAccRancher2ClusterLoggingUpdateConfigSyslog   string
	testAccRancher2ClusterLoggingRecreateConfigSyslog string
)

func init() {
	testAccRancher2ClusterLoggingConfigSyslog = `
resource "rancher2_cluster_logging" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}
`

	testAccRancher2ClusterLoggingUpdateConfigSyslog = `
resource "rancher2_cluster_logging" "foo" {
  name = "foo-updated"
  cluster_id = "` + testAccRancher2ClusterID + `"
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}
 `

	testAccRancher2ClusterLoggingRecreateConfigSyslog = `
resource "rancher2_cluster_logging" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}
 `
}

func TestAccRancher2ClusterLogging_basic_syslog(t *testing.T) {
	var cluster *managementClient.ClusterLogging

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterLoggingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterLoggingConfigSyslog,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterLoggingExists(testAccRancher2ClusterLoggingType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterLoggingUpdateConfigSyslog,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterLoggingExists(testAccRancher2ClusterLoggingType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterLoggingRecreateConfigSyslog,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterLoggingExists(testAccRancher2ClusterLoggingType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
		},
	})
}

func TestAccRancher2ClusterLogging_disappears_syslog(t *testing.T) {
	var cluster *managementClient.ClusterLogging

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterLoggingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterLoggingConfigSyslog,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterLoggingExists(testAccRancher2ClusterLoggingType+".foo", cluster),
					testAccRancher2ClusterLoggingDisappears(cluster),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterLoggingDisappears(clu *managementClient.ClusterLogging) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterLoggingType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			clu, err = client.ClusterLogging.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ClusterLogging.Delete(clu)
			if err != nil {
				return fmt.Errorf("Error removing Cluster Logging: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    clusterLoggingStateRefreshFunc(client, clu.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for cluster logging (%s) to be removed: %s", clu.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ClusterLoggingExists(n string, clu *managementClient.ClusterLogging) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster Logging ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundClu, err := client.ClusterLogging.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cluster Logging not found")
			}
			return err
		}

		clu = foundClu

		return nil
	}
}

func testAccCheckRancher2ClusterLoggingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterLoggingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ClusterLogging.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cluster Logging still exists")
	}
	return nil
}
