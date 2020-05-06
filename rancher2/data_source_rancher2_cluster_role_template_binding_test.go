package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ClusterRoleTemplateBindingDataSource(t *testing.T) {
	testAccCheckRancher2ClusterRoleTemplateBindingDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccRancher2User + testAccRancher2ClusterRoleTemplateBinding + `
data "` + testAccRancher2ClusterRoleTemplateBindingType + `" "foo" {
  name = rancher2_cluster_role_template_binding.foo.name
  cluster_id = rancher2_cluster_role_template_binding.foo.cluster_id
}
`
	name := "data." + testAccRancher2ClusterRoleTemplateBindingType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterRoleTemplateBindingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "role_template_id", "cluster-admin"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
