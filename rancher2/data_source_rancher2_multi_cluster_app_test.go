package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2MultiClusterAppDataSourceType = "rancher2_multi_cluster_app"
)

var (
	testAccCheckRancher2MultiClusterAppDataSourceConfig string
)

func init() {
	testAccCheckRancher2MultiClusterAppDataSourceConfig = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform project acceptance test"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "500m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
}
resource "rancher2_multi_cluster_app" "foo" {
  catalog_name = "library"
  name = "foo"
  targets {
    project_id = "${rancher2_project.foo.id}"
  }
  template_name = "docker-registry"
  template_version = "1.8.1"
  answers {
    values = {
      "ingress_host" = "test.xip.io"
    }
  }
  roles = ["project-member"]
}
data "` + testAccRancher2MultiClusterAppDataSourceType + `" "foo" {
  name = "${rancher2_multi_cluster_app.foo.name}"
}
`
}

func TestAccRancher2MultiClusterAppDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2MultiClusterAppDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2MultiClusterAppDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2MultiClusterAppDataSourceType+".foo", "template_version_id", "cattle-global-data:library-docker-registry-1.8.1"),
					resource.TestCheckResourceAttr("data."+testAccRancher2MultiClusterAppDataSourceType+".foo", "answers.0.values.ingress_host", "test.xip.io"),
					resource.TestCheckResourceAttr("data."+testAccRancher2MultiClusterAppDataSourceType+".foo", "roles.0", "project-member"),
				),
			},
		},
	})
}
