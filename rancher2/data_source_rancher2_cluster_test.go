package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ClusterDataSourceType = "rancher2_cluster"
)

var (
	testAccCheckRancher2ClusterDataSourceConfig string
)

func init() {
	testAccCheckRancher2ClusterDataSourceConfig = `
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
      }
    }
  }
  scheduled_cluster_scan {
    enabled = true
    scan_config {
      cis_scan_config {
        debug_master = true
        debug_worker = true
        override_benchmark_version = "rke-cis-1.5"
      }
    }
    schedule_config {
      cron_schedule = "30 * * * *"
      retention = 5
    }
  }
}
data "` + testAccRancher2ClusterDataSourceType + `" "foo" {
  name = rancher2_cluster.foo.name
}
`
}

func TestAccRancher2ClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "description", "Terraform custom cluster acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "rke_config.0.services.0.etcd.0.creation", "6h"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "scheduled_cluster_scan.0.scan_config.0.cis_scan_config.0.debug_worker", "true"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "scheduled_cluster_scan.0.schedule_config.0.cron_schedule", "30 * * * *"),
				),
			},
		},
	})
}
