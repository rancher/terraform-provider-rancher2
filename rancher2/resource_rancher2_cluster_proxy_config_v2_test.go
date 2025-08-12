package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2ClusterProxyConfigV2Type   = "rancher2_cluster_proxy_config_v2"
	testAccRancher2ClusterProxyConfigV2Config = `
resource "` + testAccRancher2ClusterProxyConfigV2Type + `" "foo" {
  cluster_id = rancher2_cluster.foo.id
  enabled = true
}
`
)

func TestAccRancher2ClusterProxyConfigV2_basic(t *testing.T) {
	var clusterProxyConfigV2 *ClusterProxyConfigV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterProxyConfigV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterV2Config + testAccRancher2ClusterProxyConfigV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterProxyConfigV2Exists(testAccRancher2ClusterProxyConfigV2Type+".foo", clusterProxyConfigV2),
					resource.TestCheckResourceAttr(testAccRancher2ClusterProxyConfigV2Type+".foo", "enabled", "true"),
				),
			},
			{
				Config: testAccRancher2ClusterV2Config + testAccRancher2ClusterProxyConfigV2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterProxyConfigV2Exists(testAccRancher2ClusterProxyConfigV2Type+".foo", clusterProxyConfigV2),
					resource.TestCheckResourceAttr(testAccRancher2ClusterProxyConfigV2Type+".foo", "enabled", "false"),
				),
			},
			{
				ResourceName:      testAccRancher2ClusterProxyConfigV2Type + ".foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRancher2ClusterProxyConfigV2ImportStateIdFunc,
			},
		},
	})
}

const testAccRancher2ClusterProxyConfigV2UpdateConfig = `
resource "` + testAccRancher2ClusterProxyConfigV2Type + `" "foo" {
  cluster_id = rancher2_cluster.foo.id
  name = "test-cluster-proxy-config"
  enabled = false
}
`

func testAccRancher2ClusterProxyConfigV2ImportStateIdFunc(s *terraform.State) (string, error) {
	resourceName := testAccRancher2ClusterProxyConfigV2Type + ".foo"
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("Not found: %s", resourceName)
	}
	clusterID := rs.Primary.Attributes["cluster_id"]
	return clusterID + "/" + clusterProxyConfigV2Name, nil
}

func testAccCheckRancher2ClusterProxyConfigV2Exists(n string, obj *ClusterProxyConfigV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ClusterProxyConfig ID is set")
		}

		meta := testAccProvider.Meta().(*Config)
		clusterID := rs.Primary.Attributes["cluster_id"]
		clusterProxyConfigV2Id := clusterID + "/" + clusterProxyConfigV2Name

		foundObj := &ClusterProxyConfigV2{}
		err := meta.getObjectV2ByID(rancher2DefaultLocalClusterID, clusterProxyConfigV2Id, clusterProxyConfigV2ApiType, foundObj)
		if err != nil {
			return err
		}

		obj = foundObj
		return nil
	}
}

func testAccCheckRancher2ClusterProxyConfigV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterProxyConfigV2Type {
			continue
		}

		meta := testAccProvider.Meta().(*Config)
		clusterID := rs.Primary.Attributes["cluster_id"]
		clusterProxyConfigV2Id := clusterID + "/" + clusterProxyConfigV2Name

		obj := &ClusterProxyConfigV2{}
		err := meta.getObjectV2ByID(rancher2DefaultLocalClusterID, clusterProxyConfigV2Id, clusterProxyConfigV2ApiType, obj)
		if err == nil {
			return fmt.Errorf("ClusterProxyConfig still exists")
		}

		if !IsNotFound(err) && !IsForbidden(err) {
			return err
		}
	}

	return nil
}
