package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ProjectLoggingDataSource(t *testing.T) {
	testAccCheckRancher2ProjectLoggingDataSourceConfig := testAccRancher2ProjectLoggingConfigSyslog + `
data "` + testAccRancher2ProjectLoggingType + `" "foo" {
  project_id = rancher2_project_logging.foo.project_id
}
`
	name := "data." + testAccRancher2ProjectLoggingType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectLoggingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "kind", "syslog"),
					resource.TestCheckResourceAttr(name, "syslog_config.0.endpoint", "192.168.1.1:514"),
				),
			},
		},
	})
}
