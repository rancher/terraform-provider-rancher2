package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ProjectAlertGroupDataSourceType = "rancher2_project_alert_group"
)

var (
	testAccCheckRancher2ProjectAlertGroupDataSourceConfig string
)

func init() {
	testAccCheckRancher2ProjectAlertGroupDataSourceConfig = `
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

resource "rancher2_project_alert_group" "foo" {
  name = "foo"
  description = "Terraform project alert group acceptance test"
  project_id = "${rancher2_project.foo.id}"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}

data "` + testAccRancher2ProjectAlertGroupDataSourceType + `" "foo" {
  name = "${rancher2_project_alert_group.foo.name}"
  project_id = "${rancher2_project_alert_group.foo.project_id}"
}
`
}

func TestAccRancher2ProjectAlertGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectAlertGroupDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectAlertGroupDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectAlertGroupDataSourceType+".foo", "description", "Terraform project alert group acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectAlertGroupDataSourceType+".foo", "group_interval_seconds", "300"),
				),
			},
		},
	})
}
