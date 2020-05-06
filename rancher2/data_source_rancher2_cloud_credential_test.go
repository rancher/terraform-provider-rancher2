package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2CloudCredentialDataSource(t *testing.T) {
	testAccCheckRancher2CloudCredentialDataSourceConfig := testAccRancher2CloudCredentialConfigAmazonec2 + `
data "` + testAccRancher2CloudCredentialType + `" "foo-aws" {
  name = rancher2_cloud_credential.foo-aws.name
}
`
	name := "data." + testAccRancher2CloudCredentialType + ".foo-aws"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CloudCredentialDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo-aws"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
