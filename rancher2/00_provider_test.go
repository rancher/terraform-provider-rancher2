package rancher2

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2DefaultClusterID = "local"
)

var (
	testAccProviders         map[string]terraform.ResourceProvider
	testAccProvider          *schema.Provider
	testAccRancher2ClusterID string
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"rancher2": testAccProvider,
	}
	err := testAccCheck()
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	err := testAccCheck()
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func testAccCheck() error {
	if os.Getenv("TF_ACC") == "1" {
		apiURL := os.Getenv("RANCHER_URL")
		tokenKey := os.Getenv("RANCHER_TOKEN_KEY")
		accessKey := os.Getenv("RANCHER_ACCESS_KEY")
		secretKey := os.Getenv("RANCHER_SECRET_KEY")
		caCerts := os.Getenv("RANCHER_CA_CERTS")
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

		if tokenKey == "" {
			return fmt.Errorf("RANCHER_TOKEN_KEY or RANCHER_ACCESS_KEY and RANCHER_SECRET_KEY must be set for acceptance tests")
		}

		config := &Config{
			URL:       apiURL,
			TokenKey:  tokenKey,
			CACerts:   caCerts,
			Insecure:  insecure,
			Bootstrap: bootstrap,
		}

		err := testAccClusterDefaultName(config)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	return nil
}

func testAccClusterDefaultName(config *Config) error {
	if testAccRancher2ClusterID == "" {
		testAccRancher2ClusterName := os.Getenv("RANCHER_ACC_CLUSTER_NAME")

		var err error

		testAccRancher2ClusterID, err = config.GetClusterIDByName(testAccRancher2ClusterName)
		if err != nil {
			return fmt.Errorf("[ERROR] getting cluster id by name: %v", err)
		}
		if len(testAccRancher2ClusterID) == 0 {
			testAccRancher2ClusterID = testAccRancher2DefaultClusterID
		}
	}
	return nil
}
