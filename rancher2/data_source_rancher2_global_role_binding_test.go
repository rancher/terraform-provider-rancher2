package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2GlobalRoleBindingDataSource(t *testing.T) {
	testAccCheckRancher2GlobalRoleBindingDataSourceConfig := testAccRancher2User + testAccRancher2GlobalRoleBinding + `
data "` + testAccRancher2GlobalRoleBindingType + `" "foo" {
  name = rancher2_global_role_binding.foo.name
  global_role_id = rancher2_global_role_binding.foo.global_role_id
}
`
	name := "data." + testAccRancher2GlobalRoleBindingType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2GlobalRoleBindingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo-test"),
					resource.TestCheckResourceAttr(name, "global_role_id", "user-base"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
