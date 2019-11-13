package rancher2

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2BootstrapType = "rancher2_bootstrap"
)

var (
	testAccRancher2BootstrapConfig         string
	testAccRancher2BootstrapPass           string
	testAccRancher2BootstrapUpdateConfig   string
	testAccRancher2BootstrapRecreateConfig string
	testAccRancher2ProviderConfig          string
)

func init() {
	testAccRancher2BootstrapPass = os.Getenv("RANCHER_ADMIN_PASS")

	testAccRancher2ProviderConfig = `
provider "rancher2" {
  bootstrap = true
  token_key = "` + providerDefaulEmptyString + `"
}
`

	testAccRancher2BootstrapConfig = testAccRancher2ProviderConfig + `
resource "rancher2_bootstrap" "foo" {
  current_password = "` + testAccRancher2BootstrapPass + `"
  password = "TestACC1234"
  telemetry = true
}
`

	testAccRancher2BootstrapUpdateConfig = testAccRancher2ProviderConfig + `
resource "rancher2_bootstrap" "foo" {
  password = "TestACC12345"
}
 `

	testAccRancher2BootstrapRecreateConfig = testAccRancher2ProviderConfig + `
resource "rancher2_bootstrap" "foo" {
  password = "TestACC1234"
  telemetry = true
}
 `
}

func TestAccRancher2Bootstrap_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2BootstrapDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2BootstrapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", "TestACC1234"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "true"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", "TestACC1234"),
				),
				ExpectNonEmptyPlan: true,
			},
			resource.TestStep{
				Config: testAccRancher2BootstrapUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", "TestACC12345"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "false"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", "TestACC12345"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2BootstrapRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", "TestACC1234"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "true"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", "TestACC1234"),
				),
			},
		},
	})
}

func testAccCheckRancher2BootstrapExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		return nil
	}
}

func testAccCheckRancher2BootstrapDestroy(s *terraform.State) error {
	return nil
}
