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
	testAccRancher2GlobalDNSType = "rancher2_global_dns"
)

var (
	testAccRancher2GlobalDNS             string
	testAccRancher2GlobalDNSUpdate       string
	testAccRancher2GlobalDNSConfig       string
	testAccRancher2GlobalDNSUpdateConfig string
)

func init() {
	testAccRancher2GlobalDNS = `data "rancher2_project" "default" {
	cluster_id = "` + testAccRancher2ClusterID + `"
	name = "Default"	
}
resource "` + testAccRancher2GlobalDNSType + `" "foo" {
	name = "foo"
	fqdn = "foo.example.com"
	provider_id = rancher2_global_dns_provider.foo-route53.id
	project_ids = [data.rancher2_project.default.id]
	ttl = 500
}
`

	testAccRancher2GlobalDNSUpdate = `data "rancher2_project" "default" {
	cluster_id = "` + testAccRancher2ClusterID + `"
	name = "Default"	
}
data "rancher2_project" "system" {
	cluster_id = "` + testAccRancher2ClusterID + `"
	name = "System"	
}
resource "` + testAccRancher2GlobalDNSType + `" "foo" {
	name = "foo-update"
	fqdn = "foo-update.example.com"
	provider_id = rancher2_global_dns_provider.foo-route53.id
	project_ids = [data.rancher2_project.default.id,data.rancher2_project.system.id]
	ttl = 600
}
`

	testAccRancher2GlobalDNSConfig = testAccRancher2GlobalDNSProviderRoute53Config + testAccRancher2GlobalDNS
	testAccRancher2GlobalDNSUpdateConfig = testAccRancher2GlobalDNSProviderRoute53Config + testAccRancher2GlobalDNSUpdate
}

func TestAccRancher2GlobalDNS_basic(t *testing.T) {
	var globalDNS *managementClient.GlobalDns

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSExists(testAccRancher2GlobalDNSType+".foo", globalDNS),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "fqdn", "foo.example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "ttl", "500"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSExists(testAccRancher2GlobalDNSType+".foo", globalDNS),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "name", "foo-update"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "fqdn", "foo-update.example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "ttl", "600"),
				),
			},
			{
				Config: testAccRancher2GlobalDNSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSExists(testAccRancher2GlobalDNSType+".foo", globalDNS),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "fqdn", "foo.example.com"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalDNSType+".foo", "ttl", "500"),
				),
			},
		},
	})
}

func TestAccRancher2GlobalDNS_disappears(t *testing.T) {
	var globalDNS *managementClient.GlobalDns

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalDNSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalDNSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalDNSExists(testAccRancher2GlobalDNSType+".foo", globalDNS),
					testAccRancher2GlobalDNSDisappears(globalDNS),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2GlobalDNSDisappears(pro *managementClient.GlobalDns) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2GlobalDNSType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.GlobalDns.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.GlobalDns.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Global DNS registry: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    globalDNSStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Global DNS registry (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2GlobalDNSExists(n string, pro *managementClient.GlobalDns) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Global DNS registry ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.GlobalDns.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Global DNS registry not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2GlobalDNSDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2GlobalDNSType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.GlobalDns.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Global DNS registry still exists")
	}
	return nil
}
