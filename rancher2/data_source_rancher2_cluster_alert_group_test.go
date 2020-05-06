package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ClusterAlertGroupDataSource(t *testing.T) {
	testAccCheckRancher2ClusterAlertGroupDataSourceConfig := testAccRancher2ClusterAlertGroupConfig + `
data "` + testAccRancher2ClusterAlertGroupType + `" "foo" {
  cluster_id = rancher2_cluster_alert_group.foo.cluster_id
  name = rancher2_cluster_alert_group.foo.name
}
`
	name := "data." + testAccRancher2ClusterAlertGroupType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterAlertGroupDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform cluster alert group acceptance test"),
					resource.TestCheckResourceAttr(name, "group_interval_seconds", "300"),
				),
			},
		},
	})
}
