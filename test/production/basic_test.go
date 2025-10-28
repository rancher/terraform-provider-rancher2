package production

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	util "github.com/rancher/terraform-provider-rancher2/test"
)

func TestProductionBasic(t *testing.T) {
	t.Parallel()
	config := util.NewTestConfig(t, "use-cases/production")

	config.AddVars(map[string]interface{}{
		"aws_access_key_id":     util.GetAwsAccessKey(),
		"aws_secret_access_key": util.GetAwsSecretKey(),
		"aws_session_token":     util.GetAwsSessionToken(),
		"aws_region":            util.GetRegion(),
	})

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
