package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2NotifierDataSourceType = "rancher2_notifier"
)

var (
	testAccCheckRancher2NotifierDataSourceConfig string
)

func init() {
	testAccCheckRancher2NotifierDataSourceConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform notifier acceptance test"
  pagerduty_config {
    service_key = "XXXXXXXX"
    proxy_url = "http://proxy.test.io"
  }
}

data "` + testAccRancher2NotifierDataSourceType + `" "foo" {
  name = "${rancher2_notifier.foo.name}"
  cluster_id = "` + testAccRancher2ClusterID + `"
}
`
}

func TestAccRancher2NotifierDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NotifierDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2NotifierDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NotifierDataSourceType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NotifierDataSourceType+".foo", "pagerduty_config.0.service_key", "XXXXXXXX"),
				),
			},
		},
	})
}
