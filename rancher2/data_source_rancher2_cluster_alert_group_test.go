package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ClusterAlertGroupDataSourceType = "rancher2_cluster_alert_group"
)

var (
	testAccCheckRancher2ClusterAlertGroupDataSourceConfig string
)

func init() {
	testAccCheckRancher2ClusterAlertGroupDataSourceConfig = `
resource "rancher2_cluster_alert_group" "foo" {
  cluster_id = "` + testAccRancher2ClusterID + `"
  name = "foo"
  description = "Terraform cluster alert group acceptance test"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}

data "` + testAccRancher2ClusterAlertGroupDataSourceType + `" "foo" {
  cluster_id = "${rancher2_cluster_alert_group.foo.cluster_id}"
  name = "${rancher2_cluster_alert_group.foo.name}"
}
`
}

func TestAccRancher2ClusterAlertGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterAlertGroupDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterAlertGroupDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterAlertGroupDataSourceType+".foo", "description", "Terraform cluster alert group acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterAlertGroupDataSourceType+".foo", "group_interval_seconds", "300"),
				),
			},
		},
	})
}
