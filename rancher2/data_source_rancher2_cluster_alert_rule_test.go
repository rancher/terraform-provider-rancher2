package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ClusterAlertRuleDataSourceType = "rancher2_cluster_alert_rule"
)

var (
	testAccCheckRancher2ClusterAlertRuleDataSourceConfig string
)

func init() {
	testAccCheckRancher2ClusterAlertRuleDataSourceConfig = `
resource "rancher2_cluster_alert_group" "foo" {
  cluster_id = "` + testAccRancher2ClusterID + `"
  name = "foo"
  description = "Terraform cluster alert rule acceptance test"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}

resource "rancher2_cluster_alert_rule" "foo" {
  cluster_id = "${rancher2_cluster_alert_group.foo.cluster_id}"
  group_id = "${rancher2_cluster_alert_group.foo.id}"
  name = "foo"
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}

data "` + testAccRancher2ClusterAlertRuleDataSourceType + `" "foo" {
  cluster_id = "${rancher2_cluster_alert_rule.foo.cluster_id}"
  name = "${rancher2_cluster_alert_rule.foo.name}"
}
`
}

func TestAccRancher2ClusterAlertRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterAlertRuleDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterAlertRuleDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterAlertRuleDataSourceType+".foo", "group_interval_seconds", "600"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterAlertRuleDataSourceType+".foo", "repeat_interval_seconds", "6000"),
				),
			},
		},
	})
}
