package rancher2

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	testAccRancher2DefaultClusterID   = "local"
	testAccRancher2DefaultClusterName = "local"
	testAccRancher2DefaultAdminPass   = "admin"
)

var (
	testAccProviderConfig                      *Config
	testAccProviders                           map[string]func() (*schema.Provider, error)
	testAccProvider                            *schema.Provider
	testAccRancher2ClusterID                   string
	testAccRancher2AdminPass                   string
	testAccRancher2ClusterRKEK8SDefaultVersion string
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"rancher2": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
	testAccRancher2ClusterID = testAccRancher2DefaultClusterID
	testAccRancher2AdminPass = testAccRancher2DefaultAdminPass
	err := testAccCheck()
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		assert.FailNow(t, "err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	err := testAccCheck()
	if err != nil {
		assert.FailNow(t, "%v", err)
	}
}

// TestGetBootstrapEnv This function doesn't test anything, but it must be run before the other tests.
// Originally this was happening as the 1st step of TestAccRancher2Upgrade (1st test), but due to some change on the test
// framework it started generating a flaky test due to the start time of the provider x the env being set.
// I also tried to get the RANCHER_TOKEN_KEY set through the initialization script but that have inconsistent results
// wile being used with the bootstrap initialization of the provider. To use that key it would take a major refactor on TestAccRancher2Upgrade.
func TestGetBootstrapEnv(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
provider "rancher2" {
  alias = "bootstrap"

  bootstrap = true
  insecure = true
  token_key = "` + providerDefaultEmptyString + `"
}
resource "rancher2_bootstrap" "foo" {
  provider = rancher2.bootstrap

  password = "` + testAccRancher2DefaultAdminPass + `"
  telemetry = true
}
provider "rancher2" {
  api_url = rancher2_bootstrap.foo.url
  token_key = rancher2_bootstrap.foo.token
  insecure = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					testAccRancher2UpgradeVars(),
				),
			},
		},
	})
}

func testAccCheck() error {
	if os.Getenv("TF_ACC") == "1" {
		apiURL := os.Getenv("RANCHER_URL")
		tokenKey := os.Getenv("RANCHER_TOKEN_KEY")
		accessKey := os.Getenv("RANCHER_ACCESS_KEY")
		secretKey := os.Getenv("RANCHER_SECRET_KEY")
		caCerts := os.Getenv("RANCHER_CA_CERTS")
		adminPass := os.Getenv("RANCHER_ADMIN_PASS")
		insecure := false
		if os.Getenv("RANCHER_INSECURE") == "true" {
			insecure = true
		}
		bootstrap := false
		if os.Getenv("RANCHER_BOOTSTRAP") == "true" {
			bootstrap = true
		}

		if apiURL == "" {
			return fmt.Errorf("RANCHER_URL must be set for acceptance tests")
		}

		if tokenKey == "" && accessKey != "" && secretKey != "" {
			tokenKey = accessKey + ":" + secretKey
		}

		testAccProviderConfig = &Config{
			URL:       apiURL,
			TokenKey:  tokenKey,
			CACerts:   caCerts,
			Insecure:  insecure,
			Bootstrap: bootstrap,
		}

		if len(adminPass) > 0 {
			testAccRancher2AdminPass = adminPass
		}

		if len(tokenKey) > 5 {
			err := testAccClusterDefaultName(testAccProviderConfig)
			if err != nil {
				return err
			}

			testAccRancher2ClusterRKEK8SDefaultVersion, err = testAccProviderConfig.getK8SDefaultVersion()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func testAccClusterDefaultName(config *Config) error {
	if testAccRancher2ClusterName := os.Getenv("RANCHER_ACC_CLUSTER_NAME"); testAccRancher2ClusterName != testAccRancher2DefaultClusterName {
		clusterID, err := config.GetClusterIDByName(testAccRancher2ClusterName)
		if err != nil {
			return fmt.Errorf("[ERROR] getting cluster id by name: %v", err)
		}
		if len(clusterID) > 0 {
			testAccRancher2ClusterID = clusterID
		}
	}
	return nil
}

func testAccRancher2UpgradeVars() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for k, rs := range s.RootModule().Resources {
			if rs.Type != "rancher2_bootstrap" {
				continue
			}
			currentPassword := rs.Primary.Attributes["current_password"]
			if err := os.Setenv("RANCHER_ADMIN_PASS", currentPassword); err != nil {
				fmt.Printf("Failed to update RANCHER_ADMIN_PASS based on resource %s with err: %s", k, err.Error())
			}
			bootstrapRancherTokenKey := rs.Primary.Attributes["token"]
			if err := os.Setenv("RANCHER_TOKEN_KEY", bootstrapRancherTokenKey); err != nil {
				fmt.Printf("Failed to update RANCHER_TOKEN_KEY based on resource %s with err: %s", k, err.Error())
			}
		}
		return nil
	}
}
