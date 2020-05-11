package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2NamespaceDataSource(t *testing.T) {
	testAccCheckRancher2NamespaceDataSourceConfig := testAccRancher2NamespaceConfig + `
data "` + testAccRancher2NamespaceType + `" "foo" {
  name = rancher2_namespace.foo.name
  project_id = rancher2_namespace.foo.project_id
}
`
	name := "data." + testAccRancher2NamespaceType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NamespaceDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform namespace acceptance test"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
