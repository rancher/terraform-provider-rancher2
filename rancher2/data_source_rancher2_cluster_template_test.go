package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ClusterTemplateDataSourceType = "rancher2_cluster_template"
)

var (
	testAccCheckRancher2ClusterTemplateDataSourceConfig string
)

func init() {
	testAccCheckRancher2ClusterTemplateDataSourceConfig = `
resource "rancher2_cluster_template" "foo" {
  name = "foo"
  members {
    access_type = "owner"
    user_principal_id = "local://user-XXXXX"
  }
  template_revisions {
    name = "V1"
    cluster_config {
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
          }
        }
        schedule_config {
          cron_schedule = "30 * * * *"
          retention = 5
        }
      }
    }
    default = true
  }
  description = "Terraform cluster template acceptance test"
}
data "` + testAccRancher2ClusterTemplateDataSourceType + `" "foo" {
  name = "${rancher2_cluster_template.foo.name}"
}
`
}

func TestAccRancher2ClusterTemplateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "description", "Terraform cluster template acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "template_revisions.0.name", "V1"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "template_revisions.0.default", "true"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "members.0.access_type", "owner"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.services.0.etcd.0.retention", "24h"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "template_revisions.0.cluster_config.0.scheduled_cluster_scan.0.scan_config.0.cis_scan_config.0.debug_worker", "true"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterTemplateDataSourceType+".foo", "template_revisions.0.cluster_config.0.scheduled_cluster_scan.0.schedule_config.0.cron_schedule", "30 * * * *"),
				),
			},
		},
	})
}
