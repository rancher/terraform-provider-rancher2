package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2GlobalRoleBindingDataSourceType = "rancher2_global_role_binding"
)

var (
	testAccCheckRancher2GlobalRoleBindingDataSourceConfig string
)

func init() {
	testAccCheckRancher2GlobalRoleBindingDataSourceConfig = `
resource "rancher2_user" "foo" {
  name = "Terraform user acceptance test"
  username = "foo"
  password = "changeme"
}
resource "rancher2_global_role_binding" "foo" {
  name = "foo"
  global_role_id = "user-base"
  user_id = "${rancher2_user.foo.id}"
}

data "` + testAccRancher2GlobalRoleBindingDataSourceType + `" "foo" {
  name = "${rancher2_global_role_binding.foo.name}"
  global_role_id = "${rancher2_global_role_binding.foo.global_role_id}"
}
`
}

func TestAccRancher2GlobalRoleBindingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2GlobalRoleBindingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalRoleBindingDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalRoleBindingDataSourceType+".foo", "global_role_id", "user-base"),
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalRoleBindingDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
