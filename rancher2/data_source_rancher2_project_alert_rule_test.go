package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ProjectAlertRuleDataSource(t *testing.T) {
	testAccCheckRancher2ProjectAlertRuleDataSourceConfig := testAccRancher2ProjectAlertRuleConfig + `
data "` + testAccRancher2ProjectAlertRuleType + `" "foo" {
  project_id = rancher2_project_alert_rule.foo.project_id
  name = rancher2_project_alert_rule.foo.name
}
`
	name := "data." + testAccRancher2ProjectAlertRuleType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectAlertRuleDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(name, "repeat_interval_seconds", "3600"),
				),
			},
		},
	})
}
