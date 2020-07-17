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
	testAccRancher2GlobalDNSEntryType = "rancher2_global_dns_entry"
)

var (
	testAccRancher2GlobalDNSEntryConfig         string
	testAccRancher2GlobalDNSEntryUpdateConfig   string
	testAccRancher2GlobalDNSEntryRecreateConfig string
)

func init() {
	testAccRancher2GlobalDNSEntryConfig = `
data "rancher2_project" "proj" {
	cluster_id = ` + testAccRancher2ClusterID + `
	name = "Default"	
}

resource "rancher2_global_dns_provider" "entry" {
	name = "foo-test-entry-provider"
	dns_provider = "route53"
	root_domain = "test-entry-example.com"
	
	route53_config {
		access_key = "YYYYYYYYYYYYYYYYYYYY"
		secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
		zone_type = "private"
		region = "us-east-1"
	}
}

resource "rancher2_global_dns_entry" "entry" {
	name = "test-entry"
	fqdn = "test-entry.example.com"
	provider_id = "${rancher2_global_dns_provider.entry.id}"
	project_ids = ["${data.rancher2_project.proj.id}"]
}

`

	testAccRancher2GlobalDNSEntryUpdateConfig = `
resource "rancher2_global_dns_entry" "entry-update" {
	name = "test-entry"
	fqdn = "test-entry-update.example.com"
	provider_id = "${rancher2_global_dns_provider.entry.id}"
	project_ids = ["${data.rancher2_project.proj.id}"]
}
 `

	testAccRancher2GlobalDNSEntryRecreateConfig = `
resource "rancher2_global_dns_entry" "entry" {
	name = "test-entry"
	fqdn = "test-entry-rc.example.com"
	provider_id = "${rancher2_global_dns_provider.entry.id}"
	project_ids = ["${data.rancher2_project.proj.id}"]
}
 `
}

func TestAccRancher2GlobalDNSEntry_basic(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDNSProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSEntryExists(testAccRancher2GlobalDNSEntryType+".foo", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSEntryType+".foo", "name", "test-entry"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSEntryType+".foo", "fqdn", "test-entry.example.com"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSEntryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSEntryExists(testAccRancher2GlobalDNSEntryType+".foo", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSEntryType+".foo", "name", "test-entry"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSEntryType+".foo", "fqdn", "test-entry-update.example.com"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSEntryRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSEntryExists(testAccRancher2GlobalDNSEntryType+".foo", globalDNSProvider),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSEntryType+".foo", "name", "test-entry"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSEntryType+".foo", "fqdn", "test-entry-rc.example.com"),
				),
			},
		},
	})
}

func TestAccRancher2GlobalDNSEntry_disappears(t *testing.T) {
	var globalDNSProvider *managementClient.GlobalDNSProvider

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSEntryExists(testAccRancher2GlobalDNSEntryType+".foo", globalDNSProvider),
					testAccRancher2GlobalDNSEntryDisappears(globalDNSProvider),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2GlobalDNSEntryDisappears(pro *managementClient.GlobalDNSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2GlobalDNSEntryType {
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

func testAccCheckRancher2GlobalDNSEntryExists(n string, pro *managementClient.GlobalDNSProvider) resource.TestCheckFunc {
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

func testAccCheckRancher2GlobalDNSEntryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2GlobalDNSEntryType {
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
		return fmt.Errorf("Global DNS Entry still exists")
	}
	return nil
}
