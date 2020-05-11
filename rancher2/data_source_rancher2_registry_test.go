package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2RegistryDataSource_Project(t *testing.T) {
	testAccCheckRancher2RegistryProjectDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccRancher2Registry + `
data "` + testAccRancher2RegistryType + `" "foo" {
  name = rancher2_registry.foo.name
  project_id = rancher2_registry.foo.project_id
}
`
	name := "data." + testAccRancher2RegistryType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2RegistryProjectDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr(name, "registries.0.address", "test.io"),
					resource.TestCheckResourceAttr(name, "registries.0.username", "user"),
				),
			},
		},
	})
}

func TestAccRancher2RegistryDataSource_Namespaced(t *testing.T) {
	testAccCheckRancher2RegistryNsDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccCheckRancher2NamespaceTestacc + testAccRancher2RegistryNs + `
data "` + testAccRancher2RegistryType + `" "foo-ns" {
  name = rancher2_registry.foo-ns.name
  project_id = rancher2_registry.foo-ns.project_id
  namespace_id = rancher2_registry.foo-ns.namespace_id
}
`
	name := "data." + testAccRancher2RegistryType + ".foo-ns"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2RegistryNsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo-ns"),
					resource.TestCheckResourceAttr(name, "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr(name, "registries.0.address", "test.io"),
					resource.TestCheckResourceAttr(name, "registries.0.username", "user"),
				),
			},
		},
	})
}
