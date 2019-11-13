package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2CatalogDataSourceType = "rancher2_catalog"
)

var (
	testAccCheckRancher2CatalogGlobalDataSourceConfig  string
	testAccCheckRancher2CatalogClusterDataSourceConfig string
	testAccCheckRancher2CatalogProjectDataSourceConfig string
)

func init() {
	testAccCheckRancher2CatalogGlobalDataSourceConfig = `
data "` + testAccRancher2CatalogDataSourceType + `" "library" {
  name = "library"
}
`
	testAccCheckRancher2CatalogClusterDataSourceConfig = `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  cluster_id = "` + testAccRancher2ClusterID + `"
  scope = "cluster"
}
data "` + testAccRancher2CatalogDataSourceType + `" "library" {
  name = "${rancher2_catalog.foo.name}"
  scope = "cluster"
}
`
	testAccCheckRancher2CatalogProjectDataSourceConfig = `
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
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  project_id = "${rancher2_project.foo.id}"
  scope = "project"
}
data "` + testAccRancher2CatalogDataSourceType + `" "library" {
  name = "${rancher2_catalog.foo.name}"
  scope = "project"
}
`
}

func TestAccRancher2CatalogDataSource_Cluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CatalogClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "scope", "cluster"),
				),
			},
		},
	})
}

func TestAccRancher2CatalogDataSource_Global(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CatalogGlobalDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "name", "library"),
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "url", "https://git.rancher.io/charts"),
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "scope", "global"),
				),
			},
		},
	})
}

func TestAccRancher2CatalogDataSource_Project(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CatalogProjectDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "scope", "project"),
				),
			},
		},
	})
}
