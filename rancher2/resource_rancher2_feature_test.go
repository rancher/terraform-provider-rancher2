package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2FeatureType   = "rancher2_feature"
	testAccRancher2FeatureConfig = `
resource "` + testAccRancher2FeatureType + `" "foo" {
	name = "unsupported-storage-drivers"
	value = true
}
`
	testAccRancher2FeatureUpdateConfig = `
resource "` + testAccRancher2FeatureType + `" "foo" {
	name = "unsupported-storage-drivers"
	value = false
}
 `
)

func TestAccRancher2Feature_basic(t *testing.T) {
	var feature *managementClient.Feature

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2FeatureConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2FeatureExists(testAccRancher2FeatureType+".foo", feature),
					resource.TestCheckResourceAttr(testAccRancher2FeatureType+".foo", "name", "unsupported-storage-drivers"),
					resource.TestCheckResourceAttr(testAccRancher2FeatureType+".foo", "value", "true"),
				),
			},
			{
				Config: testAccRancher2FeatureUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2FeatureExists(testAccRancher2FeatureType+".foo", feature),
					resource.TestCheckResourceAttr(testAccRancher2FeatureType+".foo", "name", "unsupported-storage-drivers"),
					resource.TestCheckResourceAttr(testAccRancher2FeatureType+".foo", "value", "false"),
				),
			},
			{
				Config: testAccRancher2FeatureConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2FeatureExists(testAccRancher2FeatureType+".foo", feature),
					resource.TestCheckResourceAttr(testAccRancher2FeatureType+".foo", "name", "unsupported-storage-drivers"),
					resource.TestCheckResourceAttr(testAccRancher2FeatureType+".foo", "value", "true"),
				),
			},
		},
	})
}

func testAccCheckRancher2FeatureExists(n string, feature *managementClient.Feature) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No feature ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundFeature, err := client.Feature.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Feature not found")
			}
			return err
		}

		feature = foundFeature

		return nil
	}
}
