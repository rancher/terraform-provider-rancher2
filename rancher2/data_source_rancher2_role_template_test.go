package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2RoleTemplateDataSource(t *testing.T) {
	testAccCheckRancher2RoleTemplateDataSourceConfig := testAccRancher2RoleTemplateConfig + `
data "` + testAccRancher2RoleTemplateType + `" "foo" {
  name = rancher2_role_template.foo.name
}
`
	name := "data." + testAccRancher2RoleTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2RoleTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "context", roleTemplateContextCluster),
					resource.TestCheckResourceAttr(name, "default_role", "true"),
					resource.TestCheckResourceAttr(name, "description", "Terraform role template acceptance test"),
					resource.TestCheckResourceAttr(name, "rules.0.verbs.0", policyRuleVerbCreate),
				),
			},
		},
	})
}
