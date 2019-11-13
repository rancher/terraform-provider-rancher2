package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ProjectAlertRuleDataSourceType = "rancher2_project_alert_rule"
)

var (
	testAccCheckRancher2ProjectAlertRuleDataSourceConfig string
)

func init() {
	testAccCheckRancher2ProjectAlertRuleDataSourceConfig = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform project alert rule acceptance test"
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
  description = "Terraform project alert rule acceptance test"
  project_id = "${rancher2_project.foo.id}"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}

resource "rancher2_project_alert_rule" "foo" {
  project_id = "${rancher2_project_alert_group.foo.project_id}"
  group_id = "${rancher2_project_alert_group.foo.id}"
  name = "foo"
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}

data "` + testAccRancher2ProjectAlertRuleDataSourceType + `" "foo" {
  project_id = "${rancher2_project_alert_rule.foo.project_id}"
  name = "${rancher2_project_alert_rule.foo.name}"
}
`
}

func TestAccRancher2ProjectAlertRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectAlertRuleDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectAlertRuleDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectAlertRuleDataSourceType+".foo", "group_interval_seconds", "600"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectAlertRuleDataSourceType+".foo", "repeat_interval_seconds", "6000"),
				),
			},
		},
	})
}
