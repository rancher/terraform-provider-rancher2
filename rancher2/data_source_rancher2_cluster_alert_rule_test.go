package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ClusterAlertRuleDataSource(t *testing.T) {
	testAccCheckRancher2ClusterAlertRuleDataSourceConfig := testAccRancher2ClusterAlertRuleConfig + `
data "` + testAccRancher2ClusterAlertRuleType + `" "foo" {
  cluster_id = rancher2_cluster_alert_rule.foo.cluster_id
  name = rancher2_cluster_alert_rule.foo.name
}
`
	name := "data." + testAccRancher2ClusterAlertRuleType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterAlertRuleDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(name, "repeat_interval_seconds", "3600"),
				),
			},
		},
	})
}
