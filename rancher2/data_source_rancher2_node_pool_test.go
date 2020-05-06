package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2NodePoolDataSource(t *testing.T) {
	testAccCheckRancher2NodePoolDataSourceConfig := testAccRancher2ClusterConfigRKE + testAccRancher2CloudCredentialConfigAmazonec2 + testAccRancher2NodeTemplateAmazonec2 + testAccRancher2NodePool + `
data "` + testAccRancher2NodePoolType + `" "foo" {
  name = rancher2_node_pool.foo.name
  cluster_id = rancher2_node_pool.foo.cluster_id
}
`
	name := "data." + testAccRancher2NodePoolType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NodePoolDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "hostname_prefix", "foo-cluster-0"),
					resource.TestCheckResourceAttr(name, "control_plane", "true"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
