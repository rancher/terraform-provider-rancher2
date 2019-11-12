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
				),
			},
		},
	})
}
