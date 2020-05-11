package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2NodeTemplateDataSource(t *testing.T) {
	testAccCheckRancher2NodeTemplateDataSourceConfig := testAccRancher2CloudCredentialConfigAmazonec2 + testAccRancher2NodeTemplateAmazonec2 + `
data "` + testAccRancher2NodeTemplateType + `" "foo-aws" {
  name = rancher2_node_template.foo-aws.name
}
`
	name := "data." + testAccRancher2NodeTemplateType + ".foo-aws"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NodeTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo-aws"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver amazonec2 acceptance test"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
