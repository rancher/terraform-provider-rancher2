package rancher2

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"rancher2": testAccProvider,
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
	url := os.Getenv("RANCHER_URL")
	token := os.Getenv("RANCHER_TOKEN_KEY")
	accessKey := os.Getenv("RANCHER_ACCESS_KEY")
	secretKey := os.Getenv("RANCHER_SECRET_KEY")

	if url == "" {
		t.Fatal("RANCHER_URL must be set for acceptance tests")
	}

	if token == "" && (accessKey == "" || secretKey == "") {
		t.Fatal("RANCHER_TOKEN_KEY or RANCHER_ACCESS_KEY and RANCHER_SECRET_KEY must be set for acceptance tests")
	}
}
