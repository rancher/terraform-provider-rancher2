package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2UserDataSourceType = "rancher2_user"
)

var (
	testAccCheckRancher2UserDataSourceConfig string
)

func init() {
	testAccCheckRancher2UserDataSourceConfig = `
resource "rancher2_user" "foo" {
  name = "Terraform user acceptance test"
  username = "foo"
  password = "changeme"
  enabled = "true"
}

data "` + testAccRancher2UserDataSourceType + `" "foo" {
  username = "${rancher2_user.foo.username}"
}
`
}

func TestAccRancher2UserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2UserDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2UserDataSourceType+".foo", "username", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2UserDataSourceType+".foo", "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2UserDataSourceType+".foo", "enabled", "true"),
				),
			},
		},
	})
}
