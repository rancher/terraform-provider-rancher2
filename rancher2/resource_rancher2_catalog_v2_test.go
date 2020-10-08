package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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
  git_branch = "v2.5.0"
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
	var catalog interface{}

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
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_branch", "v2.5.0"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			{
				Config: testAccRancher2CatalogV2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogV2Exists(testAccRancher2CatalogV2Type+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_repo", "https://git.rancher.io/charts"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_branch", "master"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			{
				Config: testAccRancher2CatalogV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogV2Exists(testAccRancher2CatalogV2Type+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_repo", "https://git.rancher.io/charts"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "git_branch", "v2.5.0"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogV2Type+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
		},
	})
}

func TestAccRancher2CatalogV2_disappears(t *testing.T) {
	var catalog interface{}

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

func testAccRancher2CatalogV2Disappears(cat interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2CatalogV2Type {
				continue
			}

			clusterID := rs.Primary.Attributes["cluster_id"]
			name := rs.Primary.Attributes["name"]
			client, err := testAccProvider.Meta().(*Config).catalogV2Client(clusterID)
			if err != nil {
				return err
			}
			obj, err := client.Get(name, "", metaV1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) || errors.IsForbidden(err) {
					return nil
				}
				return err
			}
			catalog := obj.(*v1.ClusterRepo)
			err = client.Delete(name, "", nil)
			if err != nil {
				return fmt.Errorf("Error removing Catalog V2 %s: %s", name, err)
			}

			timeout := int64(600)
			listOption := metaV1.ListOptions{
				TypeMeta:        catalog.TypeMeta,
				Watch:           true,
				ResourceVersion: catalog.ObjectMeta.ResourceVersion,
				TimeoutSeconds:  &timeout,
			}
			watcher, err := client.Watch("", listOption)
			for {
				select {
				case event, open := <-watcher.ResultChan():
					if open {
						if event.Type == watch.Deleted {
							watcher.Stop()
							return nil
						}
						continue
					}
					return fmt.Errorf("[ERROR] waiting for catalog V2 (%s) to be deleted", name)
				}
			}
		}
		return nil

	}
}

func testAccCheckRancher2CatalogV2Exists(n string, cat interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No catalog ID is set")
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		client, err := testAccProvider.Meta().(*Config).catalogV2Client(clusterID)
		if err != nil {
			return err
		}
		foundReg, err := client.Get(name, "", metaV1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
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
		name := rs.Primary.Attributes["name"]
		client, err := testAccProvider.Meta().(*Config).catalogV2Client(clusterID)
		if err != nil {
			return err
		}
		_, err = client.Get(name, "", metaV1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("CatalogV2 still exists")
	}
	return nil
}
