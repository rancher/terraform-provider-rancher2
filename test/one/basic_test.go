package one

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	util "github.com/rancher/terraform-provider-rancher2/test"
)

func TestOneBasic(t *testing.T) {
	t.Parallel()
	config := util.NewTestConfig(t, "use-cases/one")

	defer config.Teardown(t)
	defer config.GetErrorLogs(t)
	_, err := terraform.InitAndApplyE(t, config.TerraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		t.Fatalf("Error creating cluster: %s", err)
	}
	config.CheckReady(t)
	config.CheckRunning(t)

	// Validate that remote state can be picked up from a clean system and everything rebuilds
	os.RemoveAll(config.TestDir)
	err = util.CreateTestDirectories(t, config.ID)
	if err != nil {
		t.Log("Test failed, tearing down...")
		t.Fatalf("Error creating cluster: %s", err)
	}

	// Running the apply again should re-create everything from state in S3
	// This should only recreate the files, the AWS and Rancher resources should be untouched
	err = os.WriteFile(filepath.Join(config.TestDir, "id_rsa"), []byte(config.KeyPair.KeyPair.PrivateKey), 0600)
	if err != nil {
		t.Log("Test failed, tearing down...")
		t.Fatalf("Error creating cluster: %s", err)
	}

	err = util.WriteTerraformRc(t, config.TerraformRcPath)
	if err != nil {
		t.Fatalf("Error writing Terraform RC file: %s", err)
	}

	_, err = terraform.InitAndApplyE(t, config.TerraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		t.Fatalf("Error creating cluster: %s", err)
	}
	config.CheckReady(t)
	config.CheckRunning(t)

	// Running the apply again should not change anything
	_, err = terraform.InitAndApplyE(t, config.TerraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		t.Fatalf("Error creating cluster: %s", err)
	}
	config.CheckReady(t)
	config.CheckRunning(t)

	if t.Failed() {
		t.Log("Test failed...")
	} else {
		t.Log("Test passed...")
	}
}
