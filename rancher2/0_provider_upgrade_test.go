package rancher2

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testAccCheckRancher2ClusterSyncTestacc    string
	testAccCheckRancher2NamespaceTestacc      string
	testAccCheckRancher2UpgradeConfig         string
	testAccCheckRancher2UpgradeConfigV23      string
	testAccCheckRancher2UpgradeConfigV24      string
	testAccCheckRancher2UpgradeCluster        string
	testAccCheckRancher2UpgradeVersion        []string
	testAccCheckRancher2UpgradeCatalogV24     string
	testAccCheckRancher2UpgradeCertificateV24 string
)

func init() {
	testAccCheckRancher2ClusterSyncTestacc = `
resource "rancher2_cluster_sync" "testacc" {
  cluster_id =  "` + testAccRancher2ClusterID + `"
}
`
	testAccCheckRancher2NamespaceTestacc = `
resource "rancher2_namespace" "testacc" {
  name = "testacc"
  description = "Terraform namespace acceptance test"
  project_id = rancher2_cluster_sync.testacc.default_project_id
}
`
	testAccCheckRancher2UpgradeVersion = []string{"v2.3.6", "v2.4.8"}
	testAccCheckRancher2UpgradeCluster = os.Getenv("RANCHER_ACC_CLUSTER_NAME")
	testAccCheckRancher2UpgradeCatalogV24 = testAccRancher2CatalogGlobal + testAccRancher2CatalogCluster + testAccRancher2CatalogProject
	testAccCheckRancher2UpgradeCertificateV24 = testAccRancher2Certificate + testAccRancher2CertificateNs

	testAccCheckRancher2UpgradeConfig = `
provider "rancher2" {
  alias = "bootstrap"

  bootstrap = true
  insecure = true
  token_key = "` + providerDefaultEmptyString + `"
}
resource "rancher2_bootstrap" "foo" {
  provider = rancher2.bootstrap

  password = "` + testAccRancher2DefaultAdminPass + `"
  telemetry = true
}
provider "rancher2" {
  api_url = rancher2_bootstrap.foo.url
  token_key = rancher2_bootstrap.foo.token
  insecure = true
}
` + testAccCheckRancher2ClusterSyncTestacc + `
` + testAccCheckRancher2NamespaceTestacc + `
`

	testAccCheckRancher2UpgradeConfigV23 = testAccCheckRancher2UpgradeConfig + `
` + testAccRancher2App + `
` + testAccCheckRancher2UpgradeCatalogV24 + `
` + testAccCheckRancher2UpgradeCertificateV24 + `
` + testAccRancher2CloudCredentialConfigAmazonec2 + `
` + testAccRancher2CloudCredentialConfigAzure + `
` + testAccRancher2CloudCredentialConfigDigitalocean + `
` + testAccRancher2CloudCredentialConfigOpenstack + `
` + testAccRancher2CloudCredentialConfigVsphere + `
` + testAccRancher2ClusterConfigRKE + `
` + testAccRancher2ClusterAlertGroup + `
` + testAccRancher2ClusterAlertRule + `
` + testAccRancher2ClusterDriver + `
` + testAccRancher2ClusterLoggingSyslog + `
` + testAccRancher2User + `
` + testAccRancher2ClusterRoleTemplateBinding + `
` + testAccRancher2ClusterTemplateConfig + `
` + testAccRancher2EtcdBackup + `
` + testAccRancher2GlobalRoleBinding + `
` + testAccRancher2MultiClusterApp + `
` + testAccRancher2Namespace + `
` + testAccRancher2NodeDriver + `
` + testAccRancher2NodePool + `
` + testAccRancher2NodeTemplateAmazonec2 + `
` + testAccRancher2NodeTemplateAzure + `
` + testAccRancher2NodeTemplateDigitalocean + `
` + testAccRancher2NodeTemplateOpennebulaConfig + `
` + testAccRancher2NodeTemplateOpenstack + `
` + testAccRancher2NodeTemplateVsphere + `
` + testAccRancher2NotifierPagerduty + `
` + testAccRancher2NotifierSlack + `
` + testAccRancher2NotifierSMTP + `
` + testAccRancher2NotifierWebhook + `
` + testAccRancher2NotifierWechat + `
` + testAccCheckRancher2PodSecurityPolicyTemplate + `
` + testAccRancher2ProjectAlertGroupConfig + `
` + testAccRancher2ProjectAlertRule + `
` + testAccRancher2ProjectLoggingSyslog + `
` + testAccRancher2ProjectRoleTemplateBinding + `
` + testAccRancher2Project + `
` + testAccRancher2Registry + `
` + testAccRancher2RoleTemplateConfig + `
` + testAccRancher2Secret + `
` + testAccRancher2SecretNs + `
` + testAccRancher2SettingConfig + `
` + testAccRancher2Token + `
` + testAccRancher2TokenCluster + `
`

	testAccCheckRancher2UpgradeConfigV24 = testAccCheckRancher2UpgradeConfig + `
` + testAccRancher2App + `
` + testAccCheckRancher2UpgradeCatalogV24 + `
` + testAccCheckRancher2UpgradeCertificateV24 + `
` + testAccRancher2CloudCredentialConfigAmazonec2 + `
` + testAccRancher2CloudCredentialConfigAzure + `
` + testAccRancher2CloudCredentialConfigDigitalocean + `
` + testAccRancher2CloudCredentialConfigOpenstack + `
` + testAccRancher2CloudCredentialConfigVsphere + `
` + testAccRancher2ClusterConfigRKE + `
` + testAccRancher2ClusterAlertGroup + `
` + testAccRancher2ClusterAlertRule + `
` + testAccRancher2ClusterDriver + `
` + testAccRancher2ClusterLoggingSyslog + `
` + testAccRancher2User + `
` + testAccRancher2ClusterRoleTemplateBinding + `
` + testAccRancher2ClusterTemplateConfig + `
` + testAccRancher2EtcdBackup + `
` + testAccRancher2GlobalRoleBinding + `
` + testAccRancher2MultiClusterApp + `
` + testAccRancher2Namespace + `
` + testAccRancher2NodeDriver + `
` + testAccRancher2NodePool + `
` + testAccRancher2NodeTemplateAmazonec2 + `
` + testAccRancher2NodeTemplateAzure + `
` + testAccRancher2NodeTemplateDigitalocean + `
` + testAccRancher2NodeTemplateOpennebulaConfig + `
` + testAccRancher2NodeTemplateOpenstack + `
` + testAccRancher2NodeTemplateVsphere + `
` + testAccRancher2NotifierPagerduty + `
` + testAccRancher2NotifierSlack + `
` + testAccRancher2NotifierSMTP + `
` + testAccRancher2NotifierWebhook + `
` + testAccRancher2NotifierWechat + `
` + testAccCheckRancher2PodSecurityPolicyTemplate + `
` + testAccRancher2ProjectAlertGroupConfig + `
` + testAccRancher2ProjectAlertRule + `
` + testAccRancher2ProjectLoggingSyslog + `
` + testAccRancher2ProjectRoleTemplateBinding + `
` + testAccRancher2Project + `
` + testAccRancher2Registry + `
` + testAccRancher2RoleTemplateConfig + `
` + testAccRancher2Secret + `
` + testAccRancher2SecretNs + `
` + testAccRancher2SettingConfig + `
` + testAccRancher2Token + `
` + testAccRancher2TokenCluster + `
`
}

func TestAccRancher2Upgrade(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2UpgradeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccRancher2UpgradeVars(),
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", testAccRancher2DefaultAdminPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "true"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", testAccRancher2DefaultAdminPass),
				),
			},
			{
				Config: testAccCheckRancher2UpgradeConfigV23,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", testAccRancher2DefaultAdminPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "true"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", testAccRancher2DefaultAdminPass),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "name", "foo-global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "scope", "global"),
					resource.TestCheckResourceAttr("rancher2_cluster.foo", "name", "foo"),
					testAccRancher2UpgradeRancher(),
				),
			},
			{
				Config: testAccCheckRancher2UpgradeConfigV24,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2BootstrapExists(testAccRancher2BootstrapType+".foo"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "password", testAccRancher2DefaultAdminPass),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "telemetry", "true"),
					resource.TestCheckResourceAttr(testAccRancher2BootstrapType+".foo", "current_password", testAccRancher2DefaultAdminPass),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "name", "foo-global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "scope", "global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "version", "helm_v3"),
					resource.TestCheckResourceAttr("rancher2_cluster.foo", "name", "foo"),
				),
			},
		},
	})
}

func testAccRancher2UpgradeRancher() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		os.Setenv("RANCHER_VERSION", testAccCheckRancher2UpgradeVersion[1])
		cmd := exec.Command("make", "upgrade-rancher")
		cmd.Env = os.Environ()
		path, _ := os.Getwd()
		cmd.Dir = path + "/.."
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("Upgrading rancher: %s\n%v", out, err)
		}
		clusterActive, _, err := testAccProvider.Meta().(*Config).isClusterActive(testAccRancher2ClusterID)
		for retry := 0; retry < 10 && !clusterActive; clusterActive, _, err = testAccProvider.Meta().(*Config).isClusterActive(testAccRancher2ClusterID) {
			fmt.Printf("Waiting for cluster ID %s becomes active", testAccRancher2ClusterID)
			time.Sleep(5 * time.Second)
			retry++
		}

		if err != nil {
			return fmt.Errorf("Getting cluster ID %s state: %s", testAccRancher2ClusterID, err)
		}
		if !clusterActive {
			return fmt.Errorf("Cluster ID %s is not active", testAccRancher2ClusterID)
		}
		return nil
	}
}

func testAccRancher2UpgradeVars() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "rancher2_bootstrap" {
				continue
			}

			token := rs.Primary.Attributes["token"]
			os.Setenv("RANCHER_TOKEN_KEY", token)
			currentPassword := rs.Primary.Attributes["current_password"]
			os.Setenv("RANCHER_ADMIN_PASS", currentPassword)
		}
		return nil

	}
}
