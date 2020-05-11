package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ProjectRoleTemplateBindingDataSource(t *testing.T) {
	testAccCheckRancher2ProjectRoleTemplateBindingDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccRancher2User + testAccRancher2ProjectRoleTemplateBinding + `
data "` + testAccRancher2ProjectRoleTemplateBindingType + `" "foo" {
  name = rancher2_project_role_template_binding.foo.name
  project_id = rancher2_project_role_template_binding.foo.project_id
}
`
	name := "data." + testAccRancher2ProjectRoleTemplateBindingType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectRoleTemplateBindingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "role_template_id", "project-member"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
