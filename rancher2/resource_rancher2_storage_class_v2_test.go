package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccRancher2StorageClassV2Type = "rancher2_storage_class_v2"

var (
	testAccRancher2StorageClassV2             string
	testAccRancher2StorageClassV2Update       string
	testAccRancher2StorageClassV2Config       string
	testAccRancher2StorageClassV2UpdateConfig string
)

func init() {
	testAccRancher2StorageClassV2 = `
resource "` + testAccRancher2StorageClassV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "foo"
  parameters = {
    "param1" = "true"
    "param2" = "40000"
  }
  k8s_provisioner = "tfp.test.io/provisioner"
  reclaim_policy = "Delete"
  volume_binding_mode = "Immediate"
}
`
	testAccRancher2StorageClassV2Update = `
resource "` + testAccRancher2StorageClassV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "foo"
  parameters = {
    "param1" = "false"
    "param2" = "45000"
  }
  k8s_provisioner = "tfp.test.io/provisioner"
  reclaim_policy = "Retain"
  volume_binding_mode = "WaitForFirstConsumer"
}
 `
	testAccRancher2StorageClassV2Config = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2StorageClassV2
	testAccRancher2StorageClassV2UpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2StorageClassV2Update
}

func TestAccRancher2StorageClassV2_basic(t *testing.T) {
	var storageClass *StorageClassV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2StorageClassV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2StorageClassV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2StorageClassV2Exists(testAccRancher2StorageClassV2Type+".foo", storageClass),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "parameters.param1", "true"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "parameters.param2", "40000"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "k8s_provisioner", "tfp.test.io/provisioner"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "reclaim_policy", "Delete"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "volume_binding_mode", "Immediate"),
				),
			},
			{
				Config: testAccRancher2StorageClassV2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2StorageClassV2Exists(testAccRancher2StorageClassV2Type+".foo", storageClass),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "parameters.param1", "false"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "parameters.param2", "45000"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "k8s_provisioner", "tfp.test.io/provisioner"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "reclaim_policy", "Retain"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "volume_binding_mode", "WaitForFirstConsumer"),
				),
			},
			{
				Config: testAccRancher2StorageClassV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2StorageClassV2Exists(testAccRancher2StorageClassV2Type+".foo", storageClass),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "parameters.param1", "true"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "parameters.param2", "40000"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "k8s_provisioner", "tfp.test.io/provisioner"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "reclaim_policy", "Delete"),
					resource.TestCheckResourceAttr(testAccRancher2StorageClassV2Type+".foo", "volume_binding_mode", "Immediate"),
				),
			},
		},
	})
}

func TestAccRancher2StorageClassV2_disappears(t *testing.T) {
	var storageClass *StorageClassV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2StorageClassV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2StorageClassV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2StorageClassV2Exists(testAccRancher2StorageClassV2Type+".foo", storageClass),
					testAccRancher2StorageClassV2Disappears(storageClass),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2StorageClassV2Disappears(cat *StorageClassV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2StorageClassV2Type {
				continue
			}
			clusterID := rs.Primary.Attributes["cluster_id"]
			_, rancherID := splitID(rs.Primary.ID)
			storageClass, err := getStorageClassV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
			if err != nil {
				if IsNotFound(err) || IsForbidden(err) {
					return nil
				}
				return fmt.Errorf("testAccRancher2StorageClassV2Disappears-get: %v", err)
			}
			err = deleteStorageClassV2(testAccProvider.Meta().(*Config), clusterID, storageClass)
			if err != nil {
				return fmt.Errorf("testAccRancher2StorageClassV2Disappears-delete: %v", err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{},
				Target:     []string{"removed"},
				Refresh:    storageClassV2StateRefreshFunc(testAccProvider.Meta(), clusterID, storageClass.ID),
				Timeout:    120 * time.Second,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf("[ERROR] waiting for storageClass (%s) to be deleted: %s", storageClass.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2StorageClassV2Exists(n string, cat *StorageClassV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No storageClass ID is set")
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		foundReg, err := getStorageClassV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2StorageClassV2Exists: %v", err)
		}

		cat = foundReg

		return nil
	}
}

func testAccCheckRancher2StorageClassV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2StorageClassV2Type {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		_, err := getStorageClassV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2StorageClassV2Destroy: %v", err)
		}
		return fmt.Errorf("StorageClassV2 still exists")
	}
	return nil
}
