package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ClusterRoleTemplateBindingDataSourceType = "rancher2_cluster_role_template_binding"
)

var (
	testAccCheckRancher2ClusterRoleTemplateBindingDataSourceConfig string
)

func init() {
	testAccCheckRancher2ClusterRoleTemplateBindingDataSourceConfig = `
resource "rancher2_cluster_role_template_binding" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  role_template_id = "cluster-member"
}

data "` + testAccRancher2ClusterRoleTemplateBindingDataSourceType + `" "foo" {
  name = "${rancher2_cluster_role_template_binding.foo.name}"
  cluster_id = "${rancher2_cluster_role_template_binding.foo.cluster_id}"
}
`
}

func TestAccRancher2ClusterRoleTemplateBindingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterRoleTemplateBindingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterRoleTemplateBindingDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterRoleTemplateBindingDataSourceType+".foo", "role_template_id", "cluster-member"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterRoleTemplateBindingDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
