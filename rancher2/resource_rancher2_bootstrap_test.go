package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2BootstrapType = "rancher2_bootstrap"
	testAccRancher2BootstrapPass = "TestACC123456"
)

var (
	testAccRancher2BootstrapConfig         string
	testAccRancher2BootstrapUpdateConfig   string
	testAccRancher2BootstrapRecreateConfig string
	testAccRancher2ProviderConfig          string
)

func init() {
	testAccRancher2ProviderConfig = `
provider "rancher2" {
  bootstrap = true
  token_key = "` + providerDefaultEmptyString + `"
}
`

	testAccRancher2BootstrapConfig = testAccRancher2ProviderConfig + `
resource "` + testAccRancher2BootstrapType + `" "foo" {
  password = "` + testAccRancher2BootstrapPass + `"
  telemetry = true
}
`

	testAccRancher2BootstrapUpdateConfig = testAccRancher2ProviderConfig + `
resource "` + testAccRancher2BootstrapType + `" "foo" {
  password = "` + testAccRancher2BootstrapPass + `"
  ui_default_landing = "` + bootstrapUILandingExplorer + `"
}
 `
}

func TestAccRancher2Bootstrap_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2BootstrapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2BootstrapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", testAccRancher2BootstrapPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "true"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", testAccRancher2BootstrapPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "ui_default_landing", bootstrapUILandingManager),
				),
			},
			{
				Config: testAccRancher2BootstrapUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", testAccRancher2BootstrapPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "false"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", testAccRancher2BootstrapPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "ui_default_landing", bootstrapUILandingExplorer),
				),
			},
			{
				Config: testAccRancher2BootstrapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", testAccRancher2BootstrapPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "true"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", testAccRancher2BootstrapPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "ui_default_landing", bootstrapUILandingManager),
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
