package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2GlobalDNSProviderType = "rancher2_global_dns_provider"
)

var (
	testAccRancher2GlobalDNSProviderConfig         string
	testAccRancher2GlobalDNSProviderUpdateConfig   string
	testAccRancher2GlobalDNSProviderRecreateConfig string
)

func init() {
	testAccRancher2GlobalDNSProviderConfig = `
resource "rancher2_global_dns_provider" "foo" {
  name = "foo-test"
  dns_provider = "route53"
  root_domain = "example.com"

  route53_config {
	access_key = "YYYYYYYYYYYYYYYYYYYY"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	zone_type = "private"
	region = "us-east-1"
  }
}
`

	testAccRancher2GlobalDNSProviderUpdateConfig = `
resource "rancher2_global_dns_provider" "foo" {
  name = "foo-test-update"
  dns_provider = "route53"
  root_domain = "update.example.com"

  route53_config {
	access_key = "YYYYYYYYYYYYYYYYYYYY"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	zone_type = "private"
	region = "us-east-1"
  }
}
 `

	testAccRancher2GlobalDNSProviderRecreateConfig = `
resource "rancher2_global_dns_provider" "foo" {
	name = "foo-test-recreate"
	dns_provider = "route53"
	root_domain = "recreate.example.com"
	
	route53_config {
	  access_key = "YYYYYYYYYYYYYYYYYYYY"
	  secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	  zone_type = "private"
	  region = "us-east-1"
	}
}
 `
}

func TestAccRancher2GlobalDNSProvider_basic(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDNSProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSProviderConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo", "name", "foo-test"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo", "root_domain", "example.com"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSProviderUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo", "name", "foo-test-update"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo", "root_domain", "update.example.com"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSProviderRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo", "name", "foo-test-recreate"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSProviderType+".foo", "root_domain", "recreate.example.com"),
				),
			},
		},
	})
}

func TestAccRancher2GlobalDNSProvider_disappears(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDNSProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSProviderConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSProviderExists(testAccRancher2GlobalDNSProviderType+".foo", globalDNSProvider),
					testAccRancher2GlobalDNSProviderDisappears(globalDNSProvider),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2GlobalDNSProviderDisappears(pro *managementClient.GlobalDNSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2GlobalDNSProviderType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.GlobalDNSProvider.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.GlobalDNSProvider.Delete(pro)
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

func testAccCheckRancher2GlobalDNSProviderExists(n string, pro *managementClient.GlobalDNSProvider) resource.TestCheckFunc {
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

		foundPro, err := client.GlobalDNSProvider.ByID(rs.Primary.ID)
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

		_, err = client.GlobalDNSProvider.ByID(rs.Primary.ID)
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
