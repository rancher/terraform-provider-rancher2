package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccRancher2ConfigMapV2Type = "rancher2_config_map_v2"

var (
	testAccRancher2ConfigMapV2             string
	testAccRancher2ConfigMapV2Update       string
	testAccRancher2ConfigMapV2Config       string
	testAccRancher2ConfigMapV2UpdateConfig string
)

func init() {
	testAccRancher2ConfigMapV2 = `
resource "` + testAccRancher2ConfigMapV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "foo"
  namespace = "default"
  data = {
    "param1" = "true"
    "param2" = "40000"
  }
}
`
	testAccRancher2ConfigMapV2Update = `
resource "` + testAccRancher2ConfigMapV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "foo"
  namespace = "default"
  data = {
    "param1" = "false"
    "param2" = "80000"
  }
}
 `
	testAccRancher2ConfigMapV2Config = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ConfigMapV2
	testAccRancher2ConfigMapV2UpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ConfigMapV2Update
}

func TestAccRancher2ConfigMapV2_basic(t *testing.T) {
	var configMap *ConfigMapV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ConfigMapV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ConfigMapV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ConfigMapV2Exists(testAccRancher2ConfigMapV2Type+".foo", configMap),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "namespace", "default"),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "data.param1", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "data.param2", "40000"),
				),
			},
			{
				Config: testAccRancher2ConfigMapV2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ConfigMapV2Exists(testAccRancher2ConfigMapV2Type+".foo", configMap),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "data.param1", "false"),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "data.param2", "80000"),
				),
			},
			{
				Config: testAccRancher2ConfigMapV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ConfigMapV2Exists(testAccRancher2ConfigMapV2Type+".foo", configMap),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "namespace", "default"),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "data.param1", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ConfigMapV2Type+".foo", "data.param2", "40000"),
				),
			},
		},
	})
}

func TestAccRancher2ConfigMapV2_disappears(t *testing.T) {
	var configMap *ConfigMapV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ConfigMapV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ConfigMapV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ConfigMapV2Exists(testAccRancher2ConfigMapV2Type+".foo", configMap),
					testAccRancher2ConfigMapV2Disappears(configMap),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ConfigMapV2Disappears(cat *ConfigMapV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ConfigMapV2Type {
				continue
			}
			clusterID := rs.Primary.Attributes["cluster_id"]
			_, rancherID := splitID(rs.Primary.ID)
			configMap, err := getConfigMapV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
			if err != nil {
				if IsNotFound(err) || IsForbidden(err) {
					return nil
				}
				return fmt.Errorf("testAccRancher2ConfigMapV2Disappears-get: %v", err)
			}
			err = deleteConfigMapV2(testAccProvider.Meta().(*Config), clusterID, configMap)
			if err != nil {
				return fmt.Errorf("testAccRancher2ConfigMapV2Disappears-delete: %v", err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{},
				Target:     []string{"removed"},
				Refresh:    configMapV2StateRefreshFunc(testAccProvider.Meta(), clusterID, configMap.ID),
				Timeout:    120 * time.Second,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf("[ERROR] waiting for configMap (%s) to be deleted: %s", configMap.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ConfigMapV2Exists(n string, cat *ConfigMapV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No configMap ID is set")
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		foundReg, err := getConfigMapV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2ConfigMapV2Exists: %v", err)
		}

		cat = foundReg

		return nil
	}
}

func testAccCheckRancher2ConfigMapV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ConfigMapV2Type {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		_, err := getConfigMapV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2ConfigMapV2Destroy: %v", err)
		}
		return fmt.Errorf("ConfigMapV2 still exists")
	}
	return nil
}
