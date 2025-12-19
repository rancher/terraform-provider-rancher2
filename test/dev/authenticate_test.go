package dev

import (
	// "os" .
	// "path/filepath" .
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	// util "github.com/rancher/terraform-provider-rancher2/test" .
	cfg "github.com/rancher/terraform-provider-rancher2/test/config"
)

func TestDevAuthenticate(t *testing.T) {
	t.Parallel()
	var err error

	// generate FIT
	fit := cfg.NewTestConfig(t, "use-cases/dev", nil, nil)
	defer fit.Teardown(t)
	defer fit.GetErrorLogs(t)
	_, err = terraform.InitAndApplyE(t, fit.TerraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		t.Fatalf("Error creating cluster: %s", err)
	}
	fit.CheckReady(t)
	fit.CheckRunning(t)
	output, err := terraform.OutputAllE(t, fit.TerraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		t.Fatalf("Error getting cluster outputs: %s", err)
	}

	// Apply auth config
	authVars := map[string]interface{}{
		"rancher_url": output["address"],
		"identifier":  fit.ID,
		"owner":       fit.Owner,
	}
	auth := cfg.NewTestConfig(t, "use-cases/authenticate", authVars, nil)
	defer auth.Teardown(t)
	_, err = terraform.InitAndApplyE(t, auth.TerraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		t.Fatalf("Error creating cluster: %s", err)
	}

	t.Log("\n\nSleeping for 2 hours to allow for manual inspection...\n\n")
	time.Sleep(2 * time.Hour)

	if t.Failed() {
		t.Log("Test failed...")
	} else {
		t.Log("Test passed...")
	}
}
