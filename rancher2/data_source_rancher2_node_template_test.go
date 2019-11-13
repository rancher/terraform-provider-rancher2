package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2NodeTemplateDataSourceType = "rancher2_node_template"
)

var (
	testAccCheckRancher2NodeTemplateDataSourceConfig string
)

func init() {
	testAccCheckRancher2NodeTemplateDataSourceConfig = `
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description= "Terraform cloudCredential acceptance test"
  amazonec2_credential_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}

resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node pool acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  amazonec2_config {
	ami =  "ami-XXXXXXXXXXXXXXX"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-XXXXXXXX"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}

data "` + testAccRancher2NodeTemplateDataSourceType + `" "foo" {
  name = "${rancher2_node_template.foo.name}"
}
`
}

func TestAccRancher2NodeTemplateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NodeTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2NodeTemplateDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NodeTemplateDataSourceType+".foo", "description", "Terraform node pool acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NodeTemplateDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
