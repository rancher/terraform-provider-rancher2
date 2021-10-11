package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccRancher2CatalogV2Type = "rancher2_catalog_v2"

var (
	testAccRancher2CatalogV2             string
	testAccRancher2CatalogV2Update       string
	testAccRancher2CatalogV2Config       string
	testAccRancher2CatalogV2UpdateConfig string
)

func init() {
	testAccRancher2CatalogV2 = `
resource "` + testAccRancher2CatalogV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "foo"
  git_repo = "https://git.rancher.io/charts"
  git_branch = "dev-v2.5"
}
`
	testAccRancher2CatalogV2Update = `
resource "` + testAccRancher2CatalogV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "foo"
  git_repo = "https://git.rancher.io/charts"
  git_branch = "master"
}
 `
	testAccRancher2CatalogV2Config = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2CatalogV2
	testAccRancher2CatalogV2UpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2CatalogV2Update
}

func TestAccRancher2CatalogV2_basic(t *testing.T) {
	var catalog *ClusterRepo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CatalogV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogV2Exists(testAccRancher2CatalogV2Type+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_repo", "https://git.rancher.io/charts"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_branch", "dev-v2.5"),
				),
			},
			{
				Config: testAccRancher2CatalogV2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogV2Exists(testAccRancher2CatalogV2Type+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_repo", "https://git.rancher.io/charts"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_branch", "master"),
				),
			},
			{
				Config: testAccRancher2CatalogV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogV2Exists(testAccRancher2CatalogV2Type+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_repo", "https://git.rancher.io/charts"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_branch", "dev-v2.5"),
				),
			},
		},
	})
}

func TestAccRancher2CatalogV2_disappears(t *testing.T) {
	var catalog *ClusterRepo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CatalogV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogV2Exists(testAccRancher2CatalogV2Type+".foo", catalog),
					testAccRancher2CatalogV2Disappears(catalog),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2CatalogV2Disappears(cat *ClusterRepo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2CatalogV2Type {
				continue
			}
			clusterID := rs.Primary.Attributes["cluster_id"]
			_, rancherID := splitID(rs.Primary.ID)
			catalog, err := getCatalogV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
			if err != nil {
				if IsNotFound(err) || IsForbidden(err) {
					return nil
				}
				return fmt.Errorf("testAccRancher2CatalogV2Disappears-get: %v", err)
			}
			err = deleteCatalogV2(testAccProvider.Meta().(*Config), clusterID, catalog)
			if err != nil {
				return fmt.Errorf("testAccRancher2CatalogV2Disappears-delete: %v", err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{},
				Target:     []string{"removed"},
				Refresh:    catalogV2StateRefreshFunc(testAccProvider.Meta(), clusterID, catalog.ID),
				Timeout:    120 * time.Second,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf("[ERROR] waiting for catalog (%s) to be deleted: %s", catalog.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2CatalogV2Exists(n string, cat *ClusterRepo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No catalog ID is set")
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		foundReg, err := getCatalogV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2CatalogV2Exists: %v", err)
		}

		cat = foundReg

		return nil
	}
}

func testAccCheckRancher2CatalogV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2CatalogV2Type {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		_, err := getCatalogV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2CatalogV2Destroy: %v", err)
		}
		return fmt.Errorf("CatalogV2 still exists")
	}
	return nil
}
