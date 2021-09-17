package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ClusterV2DataSource_Cluster(t *testing.T) {
	testAccCheckRancher2ClusterV2DataSourceConfig := testAccRancher2ClusterV2 + `
data "` + testAccRancher2ClusterV2Type + `" "foo" {
  name = rancher2_cluster_v2.foo.name
  fleet_namespace = rancher2_cluster_v2.foo.fleet_namespace
}
`
	name := "data." + testAccRancher2ClusterV2Type + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterV2DataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "fleet_namespace", "fleet-default"),
					resource.TestCheckResourceAttr(name, "kubernetes_version", "v1.21.4+k3s1"),
					resource.TestCheckResourceAttr(name, "enable_network_policy", "true"),
					resource.TestCheckResourceAttr(name, "default_cluster_role_for_project_members", "user"),
				),
			},
		},
	})
}
