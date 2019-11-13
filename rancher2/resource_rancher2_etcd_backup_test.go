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
	testAccRancher2EtcdBackupType = "rancher2_etcd_backup"
)

var (
	testAccRancher2EtcdBackupConfigCluster  string
	testAccRancher2EtcdBackupConfig         string
	testAccRancher2EtcdBackupUpdateConfig   string
	testAccRancher2EtcdBackupRecreateConfig string
)

func init() {
	testAccRancher2EtcdBackupConfigCluster = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform custom cluster acceptance test"
  rke_config {
    network {
      plugin = "canal"
    }
    services {
      etcd {
        creation = "6h"
        retention = "24h"
        backup_config {
          enabled = true
          interval_hours = 20
          retention = 10
        }
      }
    }
  }
}
`
	testAccRancher2EtcdBackupConfig = testAccRancher2EtcdBackupConfigCluster + `
resource "rancher2_etcd_backup" "foo" {
  backup_config {
    enabled = true
    interval_hours = 20
    retention = 10
    s3_backup_config {
      access_key = "access_key"
      bucket_name = "bucket_name"
      endpoint = "endpoint"
      folder = "/folder"
      region = "region"
      secret_key = "secret_key"
    }
  }
  cluster_id = "${rancher2_cluster.foo.id}"
  manual = true
  name = "foo"
}
`

	testAccRancher2EtcdBackupUpdateConfig = testAccRancher2EtcdBackupConfigCluster + `
resource "rancher2_etcd_backup" "foo" {
  backup_config {
    enabled = true
    interval_hours = 20
    retention = 10
    s3_backup_config {
      access_key = "access_key"
      bucket_name = "bucket_name"
      endpoint = "endpoint"
      folder = "/folder2"
      region = "region"
      secret_key = "secret_key2"
    }
  }
  cluster_id = "${rancher2_cluster.foo.id}"
  manual = true
  name = "foo"
}
`

	testAccRancher2EtcdBackupRecreateConfig = testAccRancher2EtcdBackupConfigCluster + `
resource "rancher2_etcd_backup" "foo" {
  backup_config {
    enabled = true
    interval_hours = 20
    retention = 10
    s3_backup_config {
      access_key = "access_key"
      bucket_name = "bucket_name"
      endpoint = "endpoint"
      folder = "/folder"
      region = "region"
      secret_key = "secret_key"
    }
  }
  cluster_id = "${rancher2_cluster.foo.id}"
  manual = true
  name = "foo"
}
`
}

func TestAccRancher2EtcdBackup_basic(t *testing.T) {
	var etcdBackup *managementClient.EtcdBackup

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2EtcdBackupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2EtcdBackupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2EtcdBackupExists(testAccRancher2EtcdBackupType+".foo", etcdBackup),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "backup_config.0.s3_backup_config.0.secret_key", "secret_key"),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "backup_config.0.s3_backup_config.0.folder", "/folder"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2EtcdBackupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2EtcdBackupExists(testAccRancher2EtcdBackupType+".foo", etcdBackup),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "backup_config.0.s3_backup_config.0.secret_key", "secret_key2"),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "backup_config.0.s3_backup_config.0.folder", "/folder2"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2EtcdBackupRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2EtcdBackupExists(testAccRancher2EtcdBackupType+".foo", etcdBackup),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "backup_config.0.s3_backup_config.0.secret_key", "secret_key"),
					resource.TestCheckResourceAttr(testAccRancher2EtcdBackupType+".foo", "backup_config.0.s3_backup_config.0.folder", "/folder"),
				),
			},
		},
	})
}

func TestAccRancher2EtcdBackup_disappears(t *testing.T) {
	var etcdBackup *managementClient.EtcdBackup

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2EtcdBackupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2EtcdBackupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2EtcdBackupExists(testAccRancher2EtcdBackupType+".foo", etcdBackup),
					testAccRancher2EtcdBackupDisappears(etcdBackup),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2EtcdBackupDisappears(backup *managementClient.EtcdBackup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2EtcdBackupType {
				continue
			}

			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			backup, err := client.EtcdBackup.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.EtcdBackup.Delete(backup)
			if err != nil {
				return fmt.Errorf("Error removing Etcd Backup: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{},
				Target:     []string{"removed"},
				Refresh:    etcdBackupStateRefreshFunc(client, backup.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Etcd Backup (%s) to be removed: %s", backup.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2EtcdBackupExists(n string, backup *managementClient.EtcdBackup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Etcd Backup ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundBackup, err := client.EtcdBackup.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Etcd Backup not found")
			}
			return err
		}

		backup = foundBackup

		return nil
	}
}

func testAccCheckRancher2EtcdBackupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2EtcdBackupType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		obj, err := client.EtcdBackup.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		if obj.Removed != "" {
			return nil
		}
		return fmt.Errorf("Etcd Backup still exists")
	}
	return nil
}
