package short

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDev2Resource(t *testing.T) {
	id := id()
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
		WorkingDir: testDataPath,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ConfigFile:               config.StaticFile("../../examples/resources/rancher2_dev2/resource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("rancher2_dev2.full", "api_version", "v1"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "kind", "Rancher2Dev2"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "metadata.name", "test-rancher2-dev2"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "metadata.namespace", "test-namespace"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.string", "test_string"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.bool", "true"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.number", "1.5"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.int32", "32"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.int64", "64"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.float32", "1.25"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.float64", "4.5"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.map.foo", "bar"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.list.#", "2"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.list.0", "item1"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.list.1", "item2"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.object.string_attribute", "obj_string"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.object_list.#", "1"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.object_list.0.string_attribute", "list_obj_string"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.object_map.%", "1"),
					resource.TestCheckResourceAttr("rancher2_dev2.full", "spec.object_map.first.string_attribute", "map_obj_string"),
					resource.TestCheckResourceAttrSet("rancher2_dev2.full", "id"),
				),
			},
		},
	})
}
