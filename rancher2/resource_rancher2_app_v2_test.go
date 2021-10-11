package rancher2

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

const testAccRancher2AppV2Type = "rancher2_app_v2"

var (
	testAccRancher2AppV2             string
	testAccRancher2AppV2Update       string
	testAccRancher2AppV2Config       string
	testAccRancher2AppV2UpdateConfig string
)

func init() {
	values := v3.MapStringInterface{
		"alerts": map[string]interface{}{
			"alertmanagerSpec": map[string]interface{}{
				"enabled":  false,
				"severity": "warning",
			},
		},
	}
	valuesUpdated := v3.MapStringInterface{
		"alerts": map[string]interface{}{
			"alertmanagerSpec": map[string]interface{}{
				"enabled":  false,
				"severity": "info",
			},
		},
	}
	valuesStr, err := interfaceToGhodssyaml(values)
	if err != nil {
		log.Fatalf("[ERROR] initializing: %#v", err)
	}
	valuesStrUpdated, err := interfaceToGhodssyaml(valuesUpdated)
	if err != nil {
		log.Fatalf("[ERROR] initializing: %#v", err)
	}
	testAccRancher2AppV2 = `
resource "` + testAccRancher2AppV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "rancher-cis-benchmark"
  namespace = "cis-operator-system"
  repo_name = "rancher-charts"
  chart_name = "rancher-cis-benchmark"
  values = <<EOF
` + valuesStr + `
EOF
}
`
	testAccRancher2AppV2Update = `
resource "` + testAccRancher2AppV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "rancher-cis-benchmark"
  namespace = "cis-operator-system"
  repo_name = "rancher-charts"
  chart_name = "rancher-cis-benchmark"
  values = <<EOF
` + valuesStrUpdated + `
EOF
}
 `
	testAccRancher2AppV2Config = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2AppV2
	testAccRancher2AppV2UpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2AppV2Update
}

func TestAccRancher2AppV2_basic(t *testing.T) {
	var app *AppV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AppV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AppV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AppV2Exists(testAccRancher2AppV2Type+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "name", "rancher-cis-benchmark"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "namespace", "cis-operator-system"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "chart_name", "rancher-cis-benchmark"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			{
				Config: testAccRancher2AppV2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AppV2Exists(testAccRancher2AppV2Type+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "name", "rancher-cis-benchmark"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "namespace", "cis-operator-system"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "chart_name", "rancher-cis-benchmark"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			{
				Config: testAccRancher2AppV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AppV2Exists(testAccRancher2AppV2Type+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "name", "rancher-cis-benchmark"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "namespace", "cis-operator-system"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "chart_name", "rancher-cis-benchmark"),
					resource.TestCheckResourceAttr(testAccRancher2AppV2Type+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
		},
	})
}

func TestAccRancher2AppV2_disappears(t *testing.T) {
	var app *AppV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AppV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AppV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AppV2Exists(testAccRancher2AppV2Type+".foo", app),
					testAccRancher2AppV2Disappears(app),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2AppV2Disappears(cat *AppV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2AppV2Type {
				continue
			}

			clusterID := rs.Primary.Attributes["cluster_id"]
			name := rs.Primary.Attributes["name"]
			_, rancherID := splitID(rs.Primary.ID)
			app, err := getAppV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
			if err != nil {
				if IsNotFound(err) || IsForbidden(err) {
					return nil
				}
				return err
			}
			err = deleteAppV2(testAccProvider.Meta().(*Config), clusterID, app)
			if err != nil {
				return fmt.Errorf("Error removing App V2 %s: %s", name, err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{},
				Target:     []string{"removed"},
				Refresh:    appV2StateRefreshFunc(testAccProvider.Meta(), clusterID, app.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf("[ERROR] waiting for app (%s) to be deleted: %s", app.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2AppV2Exists(n string, cat *AppV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No app ID is set")
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		foundReg, err := getAppV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return err
		}

		cat = foundReg

		return nil
	}
}

func testAccCheckRancher2AppV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AppV2Type {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		_, err := getAppV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("AppV2 still exists")
	}
	return nil
}
