package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2CatalogDataSource_Cluster(t *testing.T) {
	testAccCheckRancher2CatalogClusterDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccRancher2CatalogCluster + `
data "` + testAccRancher2CatalogType + `" "library" {
  name = rancher2_catalog.foo-cluster.name
  scope = rancher2_catalog.foo-cluster.scope
}
`
	name := "data." + testAccRancher2CatalogType + ".library"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CatalogClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo-cluster"),
					resource.TestCheckResourceAttr(name, "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(name, "scope", "cluster"),
					resource.TestCheckResourceAttr(name, "version", "helm_v2"),
				),
			},
		},
	})
}

func TestAccRancher2CatalogDataSource_Global(t *testing.T) {
	testAccCheckRancher2CatalogGlobalDataSourceConfig := testAccRancher2CatalogGlobal + `
data "` + testAccRancher2CatalogType + `" "library" {
  name = "library"
}
`
	name := "data." + testAccRancher2CatalogType + ".library"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CatalogGlobalDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "library"),
					resource.TestCheckResourceAttr(name, "url", "https://git.rancher.io/charts"),
					resource.TestCheckResourceAttr(name, "scope", "global"),
				),
			},
		},
	})
}

func TestAccRancher2CatalogDataSource_Project(t *testing.T) {
	testAccCheckRancher2CatalogProjectDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccRancher2CatalogProject + `
data "` + testAccRancher2CatalogType + `" "library" {
  name = rancher2_catalog.foo-project.name
  scope = rancher2_catalog.foo-project.scope
}
`
	name := "data." + testAccRancher2CatalogType + ".library"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CatalogProjectDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo-project"),
					resource.TestCheckResourceAttr(name, "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(name, "scope", "project"),
					resource.TestCheckResourceAttr(name, "version", "helm_v2"),
				),
			},
		},
	})
}
