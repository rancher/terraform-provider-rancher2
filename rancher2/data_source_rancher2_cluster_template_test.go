package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ClusterTemplateDataSource(t *testing.T) {
	testAccCheckRancher2ClusterTemplateDataSourceConfig := testAccRancher2ClusterTemplateConfig + `
data "` + testAccRancher2ClusterTemplateType + `" "foo" {
  name = rancher2_cluster_template.foo.name
}
`
	name := "data." + testAccRancher2ClusterTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform cluster template acceptance test"),
					resource.TestCheckResourceAttr(name, "template_revisions.0.name", "V1"),
					resource.TestCheckResourceAttr(name, "template_revisions.0.default", "true"),
					resource.TestCheckResourceAttr(name, "members.0.access_type", "owner"),
					resource.TestCheckResourceAttr(name, "template_revisions.0.cluster_config.0.rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr(name, "template_revisions.0.cluster_config.0.rke_config.0.services.0.etcd.0.retention", "72h"),
				),
			},
		},
	})
}
