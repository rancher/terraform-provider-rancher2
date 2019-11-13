package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2RoleTemplateDataSourceType = "rancher2_role_template"
)

var (
	testAccCheckRancher2RoleTemplateDataSourceConfig string
)

func init() {
	testAccCheckRancher2RoleTemplateDataSourceConfig = `
resource "rancher2_role_template" "foo" {
  name = "foo"
  context = "` + roleTemplateContextCluster + `"
  default_role = true
  description = "Terraform role template acceptance test"
  rules {
    api_groups = ["*"]
    resources = ["secrets"]
    verbs = ["` + policyRuleVerbCreate + `"]
  }
}

data "` + testAccRancher2RoleTemplateDataSourceType + `" "foo" {
  name = "${rancher2_role_template.foo.name}"
}
`
}

func TestAccRancher2RoleTemplateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2RoleTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2RoleTemplateDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RoleTemplateDataSourceType+".foo", "context", roleTemplateContextCluster),
					resource.TestCheckResourceAttr("data."+testAccRancher2RoleTemplateDataSourceType+".foo", "default_role", "true"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RoleTemplateDataSourceType+".foo", "description", "Terraform role template acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RoleTemplateDataSourceType+".foo", "rules.0.verbs.0", policyRuleVerbCreate),
				),
			},
		},
	})
}
