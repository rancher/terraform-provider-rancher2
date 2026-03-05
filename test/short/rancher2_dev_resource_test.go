package short

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	g "github.com/gruntwork-io/terratest/modules/git"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
  fit = providerConfig + `
resource "rancher2_dev" "test" {
	id                = "dev-full-test"
	user_token        = "test"
	bool_attribute    = false
  number_attribute  = 1.1
  int64_attribute   = 1
  int32_attribute   = 1
  float64_attribute = 1.2
  float32_attribute = 1.3
  string_attribute  = "dev-test"
  list_attribute    = ["this", "is", "a", "list"]
  set_attribute     = toset(["this", "is", "a", "set"])
  map_attribute     = {
    "this" = "is"
    "a"    = "map"
  }
  nested_object = {
    string_attribute     = "dev-test"
    nested_nested_object = {
      string_attribute = "tst"
      bool_attribute   = false
    }
  }
  nested_object_list = [
    {
      string_attribute     = "dev-test"
      nested_nested_object = {
        string_attribute = "tst"
        bool_attribute   = false
      }
    },
  ]
  nested_object_map = {
    "first" = {
      string_attribute     = "dev-test"
      nested_nested_object = {
        string_attribute = "tst"
        bool_attribute   = false
      }
    }
  }
}
`
)
func TestAccDevResource(t *testing.T) {
    id := getId()
    testIdPath, err := createTestDirectories(t, id)
    if err != nil {
        t.Fatalf("Error creating test directories: %v", err)
    }
    testDataPath := filepath.Join(testIdPath, "data")
    env := map[string]string{
      "TF_ACC": "1",
      "RANCHER_INSECURE": "true", // using local rancher instance
      "RANCHER_API_URL": "https://127.0.0.1.nip.io", // using local rancher instance
      "TF_LOG": "ERROR", // change to DEBUG for more data
    }
		defer func() {
			for k := range env {
				// nolint:usetesting
				os.Unsetenv(k)
			}
      os.RemoveAll(testIdPath)
		}()
		for k, v := range env {
			// nolint:usetesting
			os.Setenv(k, v)
		}
    resource.Test(t, resource.TestCase{
        WorkingDir: testDataPath,
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read testing
            {
                Config: fit,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify number of items
                    resource.TestCheckResourceAttr("rancher2_dev.test", "list_attribute.#", "4"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "set_attribute.#", "4"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "map_attribute.%", "2"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_list.#", "1"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_map.%", "1"),
                    // Verify attribute
                    resource.TestCheckResourceAttr("rancher2_dev.test", "string_attribute", "dev-test"),
                    // Verify nested object
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object.string_attribute", "dev-test"),
                    // Verify nested attribute
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_list.0.string_attribute", "dev-test"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_list.0.nested_nested_object.string_attribute", "tst"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_list.0.nested_nested_object.bool_attribute", "false"),
                    // Verify map elements
                    resource.TestCheckResourceAttr("rancher2_dev.test", "map_attribute.this", "is"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "map_attribute.a", "map"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_map.first.nested_nested_object.string_attribute", "tst"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_map.first.nested_nested_object.bool_attribute", "false"),
                    // Verify required values are set in the state.
                    resource.TestCheckResourceAttrSet("rancher2_dev.test", "id"),
                    resource.TestCheckResourceAttrSet("rancher2_dev.test", "string_attribute"),
                    resource.TestCheckResourceAttrSet("rancher2_dev.test", "number_attribute"),
                ),
            },
            // ImportState testing
            // {
            //     ResourceName:      "rancher2_dev.test",
            //     ImportState:       true,
            //     ImportStateVerify: true,
            //     // Ignore attributes that don't exist in the remote API
            //     ImportStateVerifyIgnore: []string{
            //       "user_token",
            //     },
            // },
            // Update and Read testing
            {
                // This step's config is altered from the previous, triggering the update lifecycle.
                Config: strings.ReplaceAll(
                  fit,
                  `string_attribute  = "dev-test"`,
                  `string_attribute  = "dev-test2"`,
                ),
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify number of items
                    resource.TestCheckResourceAttr("rancher2_dev.test", "list_attribute.#", "4"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "set_attribute.#", "4"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "map_attribute.%", "2"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_list.#", "1"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_map.%", "1"),
                    // Verify attribute
                    resource.TestCheckResourceAttr("rancher2_dev.test", "string_attribute", "dev-test2"),
                    // Verify nested attribute
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_list.0.string_attribute", "dev-test"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_list.0.nested_nested_object.string_attribute", "tst"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_list.0.nested_nested_object.bool_attribute", "false"),
                    // Verify map elements
                    resource.TestCheckResourceAttr("rancher2_dev.test", "map_attribute.this", "is"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "map_attribute.a", "map"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_map.first.nested_nested_object.string_attribute", "tst"),
                    resource.TestCheckResourceAttr("rancher2_dev.test", "nested_object_map.first.nested_nested_object.bool_attribute", "false"),
                    // Verify required values are set in the state.
                    resource.TestCheckResourceAttrSet("rancher2_dev.test", "id"),
                    resource.TestCheckResourceAttrSet("rancher2_dev.test", "string_attribute"),
                    resource.TestCheckResourceAttrSet("rancher2_dev.test", "number_attribute"),
                ),
            },
            // Delete testing automatically occurs in TestCase
        },
    })
}


// Helpers.
func getRepoRoot(t *testing.T) (string, error) {
  gwd := g.GetRepoRoot(t)
	fwd, err := filepath.Abs(gwd)
	if err != nil {
		return "", err
	}
  return fwd, nil
}

func createTestDirectories(t *testing.T, id string) (testIdPath string, err error) {
  fwd, err := getRepoRoot(t)
	if err != nil {
		return "", err
	}
	paths := []string{
		filepath.Join(fwd, "test", "short", "data"),
		filepath.Join(fwd, "test", "short", "data", id),
		filepath.Join(fwd, "test", "short", "data", id, "data"),
	}
	for _, path := range paths {
		err = os.Mkdir(path, 0755)
		if err != nil && !os.IsExist(err) {
			return "", err
		}
	}
	return filepath.Join(fwd, "test", "short", "data", id), nil
}

func getId() string {
	id := os.Getenv("IDENTIFIER")
	if id == "" {
		id = random.UniqueId()
	}
	id += "-" + random.UniqueId()
	return id
}
