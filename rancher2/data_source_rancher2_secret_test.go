package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2SecretDataSource_Project(t *testing.T) {
	testAccCheckRancher2SecretProjectDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccRancher2Secret + `
data "` + testAccRancher2SecretType + `" "foo" {
  name = rancher2_secret.foo.name
  project_id = rancher2_secret.foo.project_id
}
`
	name := "data." + testAccRancher2SecretType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2SecretProjectDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(name, "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(name, "data.username", "dXNlcg=="),
				),
			},
		},
	})
}

func TestAccRancher2SecretDataSource_Namespaced(t *testing.T) {
	testAccCheckRancher2SecretNsDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccCheckRancher2NamespaceTestacc + testAccRancher2SecretNs + `
data "` + testAccRancher2SecretType + `" "foo-ns" {
  name = rancher2_secret.foo-ns.name
  project_id = rancher2_secret.foo-ns.project_id
  namespace_id = rancher2_secret.foo-ns.namespace_id
}
`
	name := "data." + testAccRancher2SecretType + ".foo-ns"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2SecretNsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo-ns"),
					resource.TestCheckResourceAttr(name, "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(name, "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(name, "data.username", "dXNlcg=="),
				),
			},
		},
	})
}
