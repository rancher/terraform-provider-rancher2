package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2GlobalRoleDataSource(t *testing.T) {
	testAccCheckRancher2GlobalRoleDataSourceConfig := testAccRancher2GlobalRoleConfig + `
data "` + testAccRancher2GlobalRoleType + `" "foo" {
  name = ` + testAccRancher2GlobalRoleType + `.foo.name
}
`
	name := "data." + testAccRancher2GlobalRoleType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2GlobalRoleDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "new_user_default", "true"),
					resource.TestCheckResourceAttr(name, "description", "Terraform global role acceptance test"),
					resource.TestCheckResourceAttr(name, "rules.0.verbs.0", policyRuleVerbCreate),
				),
			},
		},
	})
}
