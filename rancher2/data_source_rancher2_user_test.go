package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2UserDataSource(t *testing.T) {
	testAccCheckRancher2UserDataSourceConfig := testAccRancher2User + `
data "` + testAccRancher2UserType + `" "foo" {
  username = rancher2_user.foo.username
}
`
	name := "data." + testAccRancher2UserType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2UserDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "username", "foo"),
					resource.TestCheckResourceAttr(name, "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr(name, "enabled", "true"),
				),
			},
		},
	})
}
