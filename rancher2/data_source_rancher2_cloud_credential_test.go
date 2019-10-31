package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2CloudCredentialDataSourceType = "rancher2_cloud_credential"
)

var (
	testAccCheckRancher2CloudCredentialDataSourceConfig string
)

func init() {
	testAccCheckRancher2CloudCredentialDataSourceConfig = `
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description = "Terraform cloudCredential acceptance test"
  amazonec2_credential_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}

data "` + testAccRancher2CloudCredentialDataSourceType + `" "foo" {
  name = "${rancher2_cloud_credential.foo.name}"
}
`
}

func TestAccRancher2CloudCredentialDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CloudCredentialDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2CloudCredentialDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2CloudCredentialDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
