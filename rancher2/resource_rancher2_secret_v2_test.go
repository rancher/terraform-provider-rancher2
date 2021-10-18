package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccRancher2SecretV2Type = "rancher2_secret_v2"

var (
	testAccRancher2SecretV2             string
	testAccRancher2SecretV2Update       string
	testAccRancher2SecretV2Config       string
	testAccRancher2SecretV2UpdateConfig string
)

func init() {
	testAccRancher2SecretV2 = `
resource "` + testAccRancher2SecretV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "foo"
  namespace = "default"
  immutable = false
  data = {
  	password = "mypass"
  	username = "test"
  }
  type = "kubernetes.io/basic-auth"
}
`
	testAccRancher2SecretV2Update = `
resource "` + testAccRancher2SecretV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = "foo"
  namespace = "default"
  immutable = true
  data = {
  	password = "mypass-updated"
  	username = "test"
  }
  type = "kubernetes.io/basic-auth"
}
 `
	testAccRancher2SecretV2Config = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2SecretV2
	testAccRancher2SecretV2UpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2SecretV2Update
}

func TestAccRancher2SecretV2_basic(t *testing.T) {
	var secret *SecretV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SecretV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2SecretV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretV2Exists(testAccRancher2SecretV2Type+".foo", secret),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "namespace", "default"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "immutable", "false"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "data.username", "test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "data.password", "mypass"),
				),
			},
			{
				Config: testAccRancher2SecretV2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretV2Exists(testAccRancher2SecretV2Type+".foo", secret),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "namespace", "default"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "immutable", "true"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "data.username", "test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "data.password", "mypass-updated"),
				),
			},
			{
				Config: testAccRancher2SecretV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretV2Exists(testAccRancher2SecretV2Type+".foo", secret),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "namespace", "default"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "immutable", "false"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "data.username", "test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretV2Type+".foo", "data.password", "mypass"),
				),
			},
		},
	})
}

func TestAccRancher2SecretV2_disappears(t *testing.T) {
	var secret *SecretV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SecretV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2SecretV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretV2Exists(testAccRancher2SecretV2Type+".foo", secret),
					testAccRancher2SecretV2Disappears(secret),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2SecretV2Disappears(cat *SecretV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2SecretV2Type {
				continue
			}
			clusterID := rs.Primary.Attributes["cluster_id"]
			_, rancherID := splitID(rs.Primary.ID)
			secret, err := getSecretV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
			if err != nil {
				if IsNotFound(err) || IsForbidden(err) {
					return nil
				}
				return fmt.Errorf("testAccRancher2SecretV2Disappears-get: %v", err)
			}
			err = deleteSecretV2(testAccProvider.Meta().(*Config), clusterID, secret)
			if err != nil {
				return fmt.Errorf("testAccRancher2SecretV2Disappears-delete: %v", err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{},
				Target:     []string{"removed"},
				Refresh:    secretV2StateRefreshFunc(testAccProvider.Meta(), clusterID, secret.ID),
				Timeout:    120 * time.Second,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf("[ERROR] waiting for secret (%s) to be deleted: %s", secret.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2SecretV2Exists(n string, cat *SecretV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No secret ID is set")
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		foundReg, err := getSecretV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2SecretV2Exists: %v", err)
		}

		cat = foundReg

		return nil
	}
}

func testAccCheckRancher2SecretV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2SecretV2Type {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		_, rancherID := splitID(rs.Primary.ID)
		_, err := getSecretV2ByID(testAccProvider.Meta().(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2SecretV2Destroy: %v", err)
		}
		return fmt.Errorf("SecretV2 still exists")
	}
	return nil
}
