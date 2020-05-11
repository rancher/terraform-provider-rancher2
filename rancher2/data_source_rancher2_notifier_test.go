package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2NotifierDataSource(t *testing.T) {
	testAccCheckRancher2NotifierDataSourceConfig := testAccRancher2NotifierPagerdutyConfig + `
data "` + testAccRancher2NotifierType + `" "foo" {
  name = rancher2_notifier.foo-pagerduty.name
  cluster_id = rancher2_notifier.foo-pagerduty.cluster_id
}
`
	name := "data." + testAccRancher2NotifierType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NotifierDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo-pagerduty"),
					resource.TestCheckResourceAttr(name, "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(name, "pagerduty_config.0.service_key", "XXXXXXXX"),
				),
			},
		},
	})
}
