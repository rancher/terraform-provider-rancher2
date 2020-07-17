package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2GlobalDNSProviderDataSourceType = "rancher2_global_dns_provider"
)

var (
	testAccCheckRancher2GlobalDNSProviderDataSourceConfig string
)

func init() {
	testAccCheckRancher2GlobalDNSProviderDataSourceConfig = `
resource "rancher2_global_dns_provider" "dns" {
  name = "foo-test2"
  dns_provider = "route53"
  root_domain = "non.example.com"

  route53_config {
    access_key = "YYYYYYYYYYYYYYYYYYYY"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    zone_type = "private"
    region = "us-east-1"
  }
}

data "` + testAccRancher2GlobalDNSProviderDataSourceType + `" "foo" {
	name = "${rancher2_global_dns_provider.dns.name}"
}
`
}

func TestAccRancher2GlobalDNSProviderDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2GlobalDNSProviderDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalDNSProviderDataSourceType+".foo", "name", "foo-test2"),
				),
			},
		},
	})
}
