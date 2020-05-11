package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2MultiClusterAppDataSource(t *testing.T) {
	testAccCheckRancher2MultiClusterAppDataSourceConfig := testAccRancher2MultiClusterAppConfig + `
data "` + testAccRancher2MultiClusterAppType + `" "foo" {
  name = rancher2_multi_cluster_app.foo.name
}
`
	name := "data." + testAccRancher2MultiClusterAppType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2MultiClusterAppDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "template_version_id", "cattle-global-data:library-docker-registry-1.8.1"),
					resource.TestCheckResourceAttr(name, "answers.0.values.ingress_host", "test.xip.io"),
					resource.TestCheckResourceAttr(name, "roles.0", "project-member"),
				),
			},
		},
	})
}
