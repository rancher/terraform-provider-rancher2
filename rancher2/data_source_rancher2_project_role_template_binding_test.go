package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ProjectRoleTemplateBindingDataSourceType = "rancher2_project_role_template_binding"
)

var (
	testAccCheckRancher2ProjectRoleTemplateBindingDataSourceConfig string
)

func init() {
	testAccCheckRancher2ProjectRoleTemplateBindingDataSourceConfig = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform project role template binding acceptance test"
}

resource "rancher2_project_role_template_binding" "foo" {
  name = "foo"
  project_id = "${rancher2_project.foo.id}"
  role_template_id = "project-member"
}

data "` + testAccRancher2ProjectRoleTemplateBindingDataSourceType + `" "foo" {
  name = "${rancher2_project_role_template_binding.foo.name}"
  project_id = "${rancher2_project_role_template_binding.foo.project_id}"
}
`
}

func TestAccRancher2ProjectRoleTemplateBindingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectRoleTemplateBindingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectRoleTemplateBindingDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectRoleTemplateBindingDataSourceType+".foo", "role_template_id", "project-member"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectRoleTemplateBindingDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
