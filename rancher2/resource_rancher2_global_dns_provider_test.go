package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2GlobalDNSProviderType = "rancher2_global_dns_provider"
)

var (
	testAccRancher2GlobalDNSProviderAlidnsConfig           string
	testAccRancher2GlobalDNSProviderAlidnsUpdateConfig     string
	testAccRancher2GlobalDNSProviderCloudflareConfig       string
	testAccRancher2GlobalDNSProviderCloudflareUpdateConfig string
	testAccRancher2GlobalDNSProviderRoute53Config          string
	testAccRancher2GlobalDNSProviderRoute53UpdateConfig    string
)

func init() {
	testAccRancher2GlobalDNSProviderAlidnsConfig = `
resource "` + testAccRancher2GlobalDNSProviderType + `" "foo-alidns" {
  name = "foo-alidns"
  root_domain = "example.com"
  alidns_config {
	access_key = "YYYYYYYYYYYYYYYYYYYY"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2GlobalDNSProviderAlidnsUpdateConfig = `
resource "` + testAccRancher2GlobalDNSProviderType + `" "foo-alidns" {
  name = "foo-alidns-update"
  root_domain = "update.example.com"
  alidns_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
 `
	testAccRancher2GlobalDNSProviderCloudflareConfig = `
resource "` + testAccRancher2GlobalDNSProviderType + `" "foo-cloudflare" {
  name = "foo-cloudflare"
  root_domain = "example.com"
  cloudflare_config {
	api_email = "test@test.local"
	api_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2GlobalDNSProviderCloudflareUpdateConfig = `
resource "` + testAccRancher2GlobalDNSProviderType + `" "foo-cloudflare" {
  name = "foo-cloudflare-update"
  root_domain = "update.example.com"
  cloudflare_config {
	api_email = "test-update@test.local"
	api_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
 `
	testAccRancher2GlobalDNSProviderRoute53Config = `
resource "` + testAccRancher2GlobalDNSProviderType + `" "foo-route53" {
  name = "foo-route53"
  root_domain = "example.com"
  route53_config {
	access_key = "YYYYYYYYYYYYYYYYYYYY"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	zone_type = "private"
	region = "us-east-1"
  }
}
`
	testAccRancher2GlobalDNSProviderRoute53UpdateConfig = `
resource "` + testAccRancher2GlobalDNSProviderType + `" "foo-route53" {
  name = "foo-route53-update"
  root_domain = "update.example.com"
  route53_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	zone_type = "public"
	region = "us-east-1"
  }
}
 `
}

func TestAccRancher2GlobalDNSProviderAlidns_basic(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDnsProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSProviderAlidnsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-alidns", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "name", "foo-alidns"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "dns_provider", globalDNSProviderAlidnsKind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "root_domain", "example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "alidns_config.0.access_key", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSProviderAlidnsUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-alidns", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "name", "foo-alidns-update"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "dns_provider", globalDNSProviderAlidnsKind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "root_domain", "update.example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "alidns_config.0.access_key", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSProviderAlidnsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-alidns", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "name", "foo-alidns"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "dns_provider", globalDNSProviderAlidnsKind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "root_domain", "example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-alidns", "alidns_config.0.access_key", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
		},
	})
}

func TestAccRancher2GlobalDNSProviderAlidns_disappears(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDnsProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSProviderAlidnsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-alidns", globalDNSProvider),
					testAccRancher2GlobalDNSProviderDisappears(globalDNSProvider),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2GlobalDNSProviderCloudflare_basic(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDnsProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSProviderCloudflareConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "name", "foo-cloudflare"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "dns_provider", globalDNSProviderCloudflareKind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "root_domain", "example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "cloudflare_config.0.api_email", "test@test.local"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSProviderCloudflareUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "name", "foo-cloudflare-update"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "dns_provider", globalDNSProviderCloudflareKind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "root_domain", "update.example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "cloudflare_config.0.api_email", "test-update@test.local"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSProviderCloudflareConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "name", "foo-cloudflare"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "dns_provider", globalDNSProviderCloudflareKind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "root_domain", "example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", "cloudflare_config.0.api_email", "test@test.local"),
				),
			},
		},
	})
}

func TestAccRancher2GlobalDNSProviderCloudflare_disappears(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDnsProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSProviderCloudflareConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-cloudflare", globalDNSProvider),
					testAccRancher2GlobalDNSProviderDisappears(globalDNSProvider),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2GlobalDNSProviderRoute53_basic(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDnsProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSProviderRoute53Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-route53", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "name", "foo-route53"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "dns_provider", globalDNSProviderRoute53Kind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "root_domain", "example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "route53_config.0.access_key", "YYYYYYYYYYYYYYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "route53_config.0.zone_type", "private"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSProviderRoute53UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-route53", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "name", "foo-route53-update"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "dns_provider", globalDNSProviderRoute53Kind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "root_domain", "update.example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "route53_config.0.access_key", "XXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "route53_config.0.zone_type", "public"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSProviderRoute53Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-route53", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "name", "foo-route53"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "dns_provider", globalDNSProviderRoute53Kind),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "root_domain", "example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "route53_config.0.access_key", "YYYYYYYYYYYYYYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo-route53", "route53_config.0.zone_type", "private"),
				),
			},
		},
	})
}

func TestAccRancher2GlobalDNSProviderRoute53_disappears(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDnsProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSProviderRoute53Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo-route53", globalDNSProvider),
					testAccRancher2GlobalDNSProviderDisappears(globalDNSProvider),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2GlobalDNSProviderDisappears(pro *managementClient.GlobalDnsProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2GlobalDNSProviderType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.GlobalDnsProvider.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.GlobalDnsProvider.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Global DNS Provider: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    globalDNSProviderStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Global DNS Provider (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2GlobalDNSProviderExists(n string, pro *managementClient.GlobalDnsProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Global DNS Provider ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.GlobalDnsProvider.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Global DNS Provider not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2GlobalDNSProviderDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2GlobalDNSProviderType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.GlobalDnsProvider.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Global DNS Provider still exists")
	}
	return nil
}
