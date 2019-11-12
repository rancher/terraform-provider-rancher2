package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2EtcdBackupDataSourceType = "rancher2_etcd_backup"
)

var (
	testAccCheckRancher2EtcdBackupDataSourceConfig string
)

func init() {
	testAccCheckRancher2EtcdBackupDataSourceConfig = `
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
resource "rancher2_etcd_backup" "foo" {
  backup_config {
    enabled = true
    interval_hours = 20
    retention = 10
    s3_backup_config {
      access_key = "access_key"
      bucket_name = "bucket_name"
      endpoint = "endpoint"
      region = "region"
      secret_key = "secret_key"
    }
  }
  cluster_id = "${rancher2_cluster.foo.id}"
  manual = true
  name = "foo"
}
data "` + testAccRancher2EtcdBackupDataSourceType + `" "foo" {
  name = "${rancher2_etcd_backup.foo.name}"
  cluster_id = "${rancher2_etcd_backup.foo.cluster_id}"
}
`
}

func TestAccRancher2EtcdBackupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2EtcdBackupDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2EtcdBackupDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2EtcdBackupDataSourceType+".foo", "manual", "true"),
					resource.TestCheckResourceAttr("data."+testAccRancher2EtcdBackupDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
