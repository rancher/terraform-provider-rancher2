package short

import (
	"os"
	"path/filepath"
	// "strings"
	"testing"

	g "github.com/gruntwork-io/terratest/modules/git"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevResource(t *testing.T) {
	id := getId()
	testIdPath, err := createTestDirectories(t, id)
	if err != nil {
		t.Fatalf("Error creating test directories: %v", err)
	}
	testDataPath := filepath.Join(testIdPath, "data")
	env := map[string]string{
		"TF_ACC":           "1",
		"RANCHER_INSECURE": "true",                     // using local rancher instance
		"RANCHER_API_URL":  "https://127.0.0.1.nip.io", // using local rancher instance
		"TF_LOG":           "DEBUG",                    // change to DEBUG for more data
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
		WorkingDir:               testDataPath,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				ConfigFile: config.StaticFile("../../examples/resources/rancher2_dev/resource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("rancher2_dev.full", "list_attribute.#", "4"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "set_attribute.#", "4"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "map_attribute.%", "2"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_list.#", "1"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_map.%", "1"),
					// Verify attribute
					resource.TestCheckResourceAttr("rancher2_dev.full", "string_attribute", "dev-test"),
					// Verify nested object
					resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object.string_attribute", "test"),
					// Verify nested attribute
					resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_list.0.string_attribute", "test"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_list.0.nested_nested_object.string_attribute", "tst"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_list.0.nested_nested_object.bool_attribute", "false"),
					// Verify map elements
					resource.TestCheckResourceAttr("rancher2_dev.full", "map_attribute.this", "is"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "map_attribute.a", "map"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_map.first.nested_nested_object.string_attribute", "tst"),
					resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_map.first.nested_nested_object.bool_attribute", "false"),
					// Verify required values are set in the state.
					resource.TestCheckResourceAttrSet("rancher2_dev.full", "id"),
					resource.TestCheckResourceAttrSet("rancher2_dev.full", "string_attribute"),
					resource.TestCheckResourceAttrSet("rancher2_dev.full", "number_attribute"),
				),
			},
			// ImportState testing
			// {
			//     ResourceName:      "rancher2_dev.full",
			//     ImportState:       true,
			//     ImportStateVerify: true,
			//     // Ignore attributes that don't exist in the remote API
			//     ImportStateVerifyIgnore: []string{
			//       "user_token",
			//     },
			// },
			// // Update and Read testing
			// {
			// 	// This step's config is altered from the previous, triggering the update lifecycle.
			// 	Config: strings.ReplaceAll(
			// 		// read the file into a string and replace the value of string_attribute so that an update is necessary
			// 		func() string {
			// 			b, err := os.ReadFile("../../examples/resources/rancher2_dev/resource.tf")
			// 			if err != nil {
			// 				t.Fatalf("Error reading resource.tf: %v", err)
			// 			}
			// 			return string(b)
			// 		}(),
			// 		`string_attribute  = "dev-test"`,
			// 		`string_attribute  = "dev-test2"`,
			// 	),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		// Verify number of items
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "list_attribute.#", "4"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "set_attribute.#", "4"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "map_attribute.%", "2"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_list.#", "1"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_map.%", "1"),
			// 		// Verify attribute
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "string_attribute", "dev-test2"),
			// 		// Verify nested attribute
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_list.0.string_attribute", "test"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_list.0.nested_nested_object.string_attribute", "tst"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_list.0.nested_nested_object.bool_attribute", "false"),
			// 		// Verify map elements
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "map_attribute.this", "is"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "map_attribute.a", "map"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_map.first.nested_nested_object.string_attribute", "tst"),
			// 		resource.TestCheckResourceAttr("rancher2_dev.full", "nested_object_map.first.nested_nested_object.bool_attribute", "false"),
			// 		// Verify required values are set in the state.
			// 		resource.TestCheckResourceAttrSet("rancher2_dev.full", "id"),
			// 		resource.TestCheckResourceAttrSet("rancher2_dev.full", "string_attribute"),
			// 		resource.TestCheckResourceAttrSet("rancher2_dev.full", "number_attribute"),
			// 	),
			// },
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
