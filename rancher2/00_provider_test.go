package rancher2

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	testAccRancher2DefaultClusterID   = "local"
	testAccRancher2DefaultClusterName = "local"
	testAccRancher2DefaultAdminPass   = "admin"
)

var (
	testAccProviderConfig                      *Config
	testAccProviders                           map[string]terraform.ResourceProvider
	testAccProvider                            *schema.Provider
	testAccRancher2ClusterID                   string
	testAccRancher2AdminPass                   string
	testAccRancher2ClusterRKEK8SDefaultVersion string
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"rancher2": testAccProvider,
	}
	testAccRancher2ClusterID = testAccRancher2DefaultClusterID
	testAccRancher2AdminPass = testAccRancher2DefaultAdminPass
	err := testAccCheck()
	if err != nil {
		log.Fatalf("failed check %s", err)
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		assert.FailNow(t, "err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	err := testAccCheck()
	if err != nil {
		assert.FailNow(t, "%v", err)
	}
}

func testAccCheck() error {
	if os.Getenv("TF_ACC") == "1" {
		apiURL := os.Getenv("RANCHER_URL")
		tokenKey := os.Getenv("RANCHER_TOKEN_KEY")
		accessKey := os.Getenv("RANCHER_ACCESS_KEY")
		secretKey := os.Getenv("RANCHER_SECRET_KEY")
		caCerts := os.Getenv("RANCHER_CA_CERTS")
		adminPass := os.Getenv("RANCHER_ADMIN_PASS")
		insecure := os.Getenv("RANCHER_INSECURE") == "true"
		bootstrap := os.Getenv("RANCHER_BOOTSTRAP") == "true"

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
				return fmt.Errorf("failed to test the default cluster name: %w", err)
			}

			testAccRancher2ClusterRKEK8SDefaultVersion, err = testAccProviderConfig.getK8SDefaultVersion()
			if err != nil {
				return fmt.Errorf("failed get the default k8s version: %w", err)
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
