package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2GlobalDNSProviderDataSourceType = "rancher2_global_dns_provider"
)

func TestAccRancher2GlobalDNSProviderDataSource(t *testing.T) {
	testAccCheckRancher2GlobalDNSProviderDataSourceConfig := testAccRancher2GlobalDNSProviderRoute53Config + `
data "` + testAccRancher2GlobalDNSProviderDataSourceType + `" "foo" {
  name = rancher2_global_dns_provider.foo-route53.name
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2GlobalDNSProviderDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalDNSProviderDataSourceType+".foo", "name", "foo-route53"),
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalDNSProviderDataSourceType+".foo", "dns_provider", globalDNSProviderRoute53Kind),
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalDNSProviderDataSourceType+".foo", "root_domain", "example.com"),
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalDNSProviderDataSourceType+".foo", "route53_config.0.access_key", "YYYYYYYYYYYYYYYYYYYY"),
					resource.TestCheckResourceAttr("data."+testAccRancher2GlobalDNSProviderDataSourceType+".foo", "route53_config.0.zone_type", "private"),
				),
			},
		},
	})
}
