package dev

import (
	// "os" .
	// "path/filepath" .
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	// util "github.com/rancher/terraform-provider-rancher2/test" .
	cfg "github.com/rancher/terraform-provider-rancher2/test/config"
)

func TestDevBasic(t *testing.T) {
	t.Parallel()
	config := cfg.NewTestConfig(t, "use-cases/dev", nil, nil)
	defer config.Teardown(t)
	defer config.GetErrorLogs(t)
	_, err := terraform.InitAndApplyE(t, config.TerraformOptions)
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
