package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ProjectAlertGroupDataSource(t *testing.T) {
	testAccCheckRancher2ProjectAlertGroupDataSourceConfig := testAccRancher2ProjectAlertGroupConfig + `
data "` + testAccRancher2ProjectAlertGroupType + `" "foo" {
  name = rancher2_project_alert_group.foo.name
  project_id = rancher2_project_alert_group.foo.project_id
}
`
	name := "data." + testAccRancher2ProjectAlertGroupType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectAlertGroupDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform project alert group acceptance test"),
					resource.TestCheckResourceAttr(name, "group_interval_seconds", "300"),
				),
			},
		},
	})
}
