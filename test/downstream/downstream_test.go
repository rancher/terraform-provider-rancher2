package downstream

import (
	"os"
	"path/filepath"
	"testing"

	aws "github.com/gruntwork-io/terratest/modules/aws"
	g "github.com/gruntwork-io/terratest/modules/git"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	util "github.com/rancher/terraform-provider-rancher2/test"
)

func TestDownstreamBasic(t *testing.T) {
	t.Parallel()
	id := util.GetId()
	region := util.GetRegion()
	accessKey := util.GetAwsAccessKey()
	secretKey := util.GetAwsSecretKey()
	sessionToken := util.GetAwsSessionToken()
	directory := "deploy_rke2"
	owner := "terraform-ci@suse.com"
	util.SetAcmeServer()

	repoRoot, err := filepath.Abs(g.GetRepoRoot(t))
	if err != nil {
		t.Fatalf("Error getting git root directory: %v", err)
	}

	exampleDir := repoRoot + "/examples/" + directory
	testDir := repoRoot + "/test/data/" + id

	err = util.CreateTestDirectories(t, id)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test data directories: %s", err)
	}
	keyPair, err := util.CreateKeypair(t, region, owner, id)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test key pair: %s", err)
	}
	err = os.WriteFile(testDir+"/id_rsa", []byte(keyPair.KeyPair.PrivateKey), 0600)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test key pair: %s", err)
	}
	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
	t.Logf("Key %s created and added to agent", keyPair.Name)

	// use oldest RKE2, remember it releases much more than Rancher
	_, _, rke2Version, err := util.GetRke2Releases()
	if err != nil {
		os.RemoveAll(testDir)
		aws.DeleteEC2KeyPair(t, keyPair)
		sshAgent.Stop()
		t.Fatalf("Error getting Rke2 release version: %s", err)
	}

	rancherVersion := os.Getenv("RANCHER_VERSION")
	if rancherVersion == "" {
		// use stable version if not specified
		// using stable prevents problems where the Rancher provider hasn't released to fit the latest Rancher
		_, rancherVersion, _, err = util.GetRancherReleases()
	}
	if err != nil {
		os.RemoveAll(testDir)
		aws.DeleteEC2KeyPair(t, keyPair)
		sshAgent.Stop()
		t.Fatalf("Error getting Rancher release version: %s", err)
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: exampleDir,
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"identifier":            id,
			"owner":                 owner,
			"key_name":              keyPair.Name,
			"key":                   keyPair.KeyPair.PublicKey,
			"zone":                  os.Getenv("ZONE"),
			"rke2_version":          rke2Version,
			"rancher_version":       rancherVersion,
			"file_path":             testDir,
			"aws_access_key_id":     accessKey,
			"aws_secret_access_key": secretKey,
			"aws_session_token":     sessionToken,
			"aws_region":            region,
		},
		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": region,
			"AWS_REGION":         region,
			"TF_DATA_DIR":        testDir,
			"TF_IN_AUTOMATION":   "1",
			"KUBECONFIG":         testDir + "/kubeconfig",
			"KUBE_CONFIG_PATH":   testDir,
			// Adding these as environment variables so that they are not displayed in the log
			"TF_VAR_aws_access_key_id":     accessKey,
			"TF_VAR_aws_secret_access_key": secretKey,
			"TF_VAR_aws_session_token":     sessionToken,
			"TF_CLI_ARGS_plan":             "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_apply":            "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_destroy":          "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_output":           "-no-color -state=" + testDir + "/tfstate",
		},
		RetryableTerraformErrors: util.GetRetryableTerraformErrors(),
		NoColor:                  true,
		SshAgent:                 sshAgent,
		Upgrade:                  true,
	})

	_, err = terraform.InitAndApplyE(t, terraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		util.GetErrorLogs(t, testDir+"/kubeconfig")
		util.Teardown(t, testDir, terraformOptions, keyPair)
		os.Remove(exampleDir + ".terraform.lock.hcl")
		sshAgent.Stop()
		t.Fatalf("Error creating cluster: %s", err)
	}
	util.CheckReady(t, testDir+"/kubeconfig")
	util.CheckRunning(t, testDir+"/kubeconfig")
	if t.Failed() {
		t.Log("Test failed...")
	} else {
		t.Log("Test passed...")
	}
	util.Teardown(t, testDir, terraformOptions, keyPair)
	os.Remove(exampleDir + "/.terraform.lock.hcl")
	sshAgent.Stop()
}

func TestDownstreamRoles(t *testing.T) {
	t.Parallel()
	id := util.GetId()
	region := util.GetRegion()
	accessKey := util.GetAwsAccessKey()
	secretKey := util.GetAwsSecretKey()
	sessionToken := util.GetAwsSessionToken()
	directory := "deploy_rke2_multiple_pools"
	owner := "terraform-ci@suse.com"
	util.SetAcmeServer()

	repoRoot, err := filepath.Abs(g.GetRepoRoot(t))
	if err != nil {
		t.Fatalf("Error getting git root directory: %v", err)
	}

	exampleDir := repoRoot + "/examples/" + directory
	testDir := repoRoot + "/test/data/" + id

	err = util.CreateTestDirectories(t, id)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test data directories: %s", err)
	}
	keyPair, err := util.CreateKeypair(t, region, owner, id)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test key pair: %s", err)
	}

	err = os.WriteFile(testDir+"/id_rsa", []byte(keyPair.KeyPair.PrivateKey), 0600)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test key pair: %s", err)
	}

	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
	t.Logf("Key %s created and added to agent", keyPair.Name)

	// use oldest RKE2, remember it releases much more than Rancher
	_, _, rke2Version, err := util.GetRke2Releases()
	if err != nil {
		os.RemoveAll(testDir)
		aws.DeleteEC2KeyPair(t, keyPair)
		sshAgent.Stop()
		t.Fatalf("Error getting Rke2 release version: %s", err)
	}

	rancherVersion := os.Getenv("RANCHER_VERSION")
	if rancherVersion == "" {
		// use stable version if not specified
		// using stable prevents problems where the Rancher provider hasn't released to fit the latest Rancher
		_, rancherVersion, _, err = util.GetRancherReleases()
	}
	if err != nil {
		os.RemoveAll(testDir)
		aws.DeleteEC2KeyPair(t, keyPair)
		sshAgent.Stop()
		t.Fatalf("Error getting Rancher release version: %s", err)
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: exampleDir,
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"identifier":            id,
			"owner":                 owner,
			"key_name":              keyPair.Name,
			"key":                   keyPair.KeyPair.PublicKey,
			"zone":                  os.Getenv("ZONE"),
			"rke2_version":          rke2Version,
			"rancher_version":       rancherVersion,
			"file_path":             testDir,
			"aws_access_key_id":     accessKey,
			"aws_secret_access_key": secretKey,
			"aws_session_token":     sessionToken,
			"aws_region":            region,
		},
		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION":  region,
			"AWS_REGION":          region,
			"TF_DATA_DIR":         testDir,
			"TF_IN_AUTOMATION":    "1",
			"KUBECONFIG":          testDir + "/kubeconfig",
			"KUBE_CONFIG_PATH":    testDir,
			"TF_CLI_ARGS_plan":    "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_apply":   "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_destroy": "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_output":  "-no-color -state=" + testDir + "/tfstate",
		},
		RetryableTerraformErrors: util.GetRetryableTerraformErrors(),
		NoColor:                  true,
		SshAgent:                 sshAgent,
		Upgrade:                  true,
	})

	_, err = terraform.InitAndApplyE(t, terraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		util.GetErrorLogs(t, testDir+"/kubeconfig")
		util.Teardown(t, testDir, terraformOptions, keyPair)
		os.Remove(exampleDir + ".terraform.lock.hcl")
		sshAgent.Stop()
		t.Fatalf("Error creating cluster: %s", err)
	}
	util.CheckReady(t, testDir+"/kubeconfig")
	util.CheckRunning(t, testDir+"/kubeconfig")
	if t.Failed() {
		t.Log("Test failed...")
	} else {
		t.Log("Test passed...")
	}
	util.Teardown(t, testDir, terraformOptions, keyPair)
	os.Remove(exampleDir + "/.terraform.lock.hcl")
	sshAgent.Stop()
}

func TestDownstreamImport(t *testing.T) {
	t.Parallel()
	id := util.GetId()
	region := util.GetRegion()
	accessKey := util.GetAwsAccessKey()
	secretKey := util.GetAwsSecretKey()
	sessionToken := util.GetAwsSessionToken()
	directory := "state_import"
	owner := "terraform-ci@suse.com"
	util.SetAcmeServer()

	repoRoot, err := filepath.Abs(g.GetRepoRoot(t))
	if err != nil {
		t.Fatalf("Error getting git root directory: %v", err)
	}

	exampleDir := repoRoot + "/examples/" + directory
	testDir := repoRoot + "/test/data/" + id

	err = util.CreateTestDirectories(t, id)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test data directories: %s", err)
	}
	keyPair, err := util.CreateKeypair(t, region, owner, id)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test key pair: %s", err)
	}
	err = os.WriteFile(testDir+"/id_rsa", []byte(keyPair.KeyPair.PrivateKey), 0600)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("Error creating test key pair: %s", err)
	}
	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
	t.Logf("Key %s created and added to agent", keyPair.Name)

	// use oldest RKE2, remember it releases much more than Rancher
	_, _, rke2Version, err := util.GetRke2Releases()
	if err != nil {
		os.RemoveAll(testDir)
		aws.DeleteEC2KeyPair(t, keyPair)
		sshAgent.Stop()
		t.Fatalf("Error getting Rke2 release version: %s", err)
	}

	rancherVersion := os.Getenv("RANCHER_VERSION")
	if rancherVersion == "" {
		// use stable version if not specified
		// using stable prevents problems where the Rancher provider hasn't released to fit the latest Rancher
		_, rancherVersion, _, err = util.GetRancherReleases()
	}
	if err != nil {
		os.RemoveAll(testDir)
		aws.DeleteEC2KeyPair(t, keyPair)
		sshAgent.Stop()
		t.Fatalf("Error getting Rancher release version: %s", err)
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: exampleDir,
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"identifier":            id,
			"owner":                 owner,
			"key_name":              keyPair.Name,
			"key":                   keyPair.KeyPair.PublicKey,
			"zone":                  os.Getenv("ZONE"),
			"rke2_version":          rke2Version,
			"rancher_version":       rancherVersion,
			"file_path":             testDir,
			"aws_access_key_id":     accessKey,
			"aws_secret_access_key": secretKey,
			"aws_session_token":     sessionToken,
			"aws_region":            region,
			"data_path":             testDir,
		},
		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": region,
			"AWS_REGION":         region,
			"TF_DATA_DIR":        testDir,
			"TF_IN_AUTOMATION":   "1",
			"KUBECONFIG":         testDir + "/kubeconfig",
			"KUBE_CONFIG_PATH":   testDir,
			// Adding these as environment variables so that they are not displayed in the log
			"TF_VAR_aws_access_key_id":     accessKey,
			"TF_VAR_aws_secret_access_key": secretKey,
			"TF_VAR_aws_session_token":     sessionToken,
			"TF_CLI_ARGS_init":             "-upgrade",
			"TF_CLI_ARGS_plan":             "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_apply":            "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_destroy":          "-no-color -state=" + testDir + "/tfstate",
			"TF_CLI_ARGS_output":           "-no-color -state=" + testDir + "/tfstate",
		},
		RetryableTerraformErrors: util.GetRetryableTerraformErrors(),
		NoColor:                  true,
		SshAgent:                 sshAgent,
		Upgrade:                  true,
	})

	terraformOptions.Vars["project_mismatch"] = "true"
	_, err = terraform.InitAndApplyE(t, terraformOptions)
	if err == nil {
		t.Log("Test failed, tearing down...")
		util.GetErrorLogs(t, testDir+"/kubeconfig")
		util.Teardown(t, testDir, terraformOptions, keyPair)
		os.Remove(exampleDir + ".terraform.lock.hcl")
		sshAgent.Stop()
		t.Fatalf("Error creating cluster: %s", err)
	}

	terraformOptions.Vars["project_mismatch"] = "false"
	_, err = terraform.ApplyE(t, terraformOptions)
	if err != nil {
		t.Log("Test failed, tearing down...")
		util.GetErrorLogs(t, testDir+"/kubeconfig")
		util.Teardown(t, testDir, terraformOptions, keyPair)
		os.Remove(exampleDir + ".terraform.lock.hcl")
		sshAgent.Stop()
		t.Fatalf("Error creating cluster: %s", err)
	}
	util.CheckReady(t, testDir+"/kubeconfig")
	util.CheckRunning(t, testDir+"/kubeconfig")

	if t.Failed() {
		t.Log("Test failed...")
	} else {
		t.Log("Test passed...")
	}

	util.Teardown(t, testDir, terraformOptions, keyPair)
	os.Remove(exampleDir + "/.terraform.lock.hcl")
	os.Remove(exampleDir + "/TF_DATA_DIR.env")
	sshAgent.Stop()
}
