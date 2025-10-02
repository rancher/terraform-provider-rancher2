package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	k8sVersion                    string
	testClusterRKEConfigConf      *managementClient.RancherKubernetesEngineConfig
	testClusterRKEConfigInterface []interface{}
)

func testClusterRKEConfig() {
	k8sVersion = "test"
	testClusterRKEConfigConf = &managementClient.RancherKubernetesEngineConfig{
		AddonJobTimeout:     30,
		Addons:              "addons",
		AddonsInclude:       []string{"addon1", "addon2"},
		Authentication:      testClusterRKEConfigAuthenticationConf,
		Authorization:       testClusterRKEConfigAuthorizationConf,
		BastionHost:         testClusterRKEConfigBastionHostConf,
		CloudProvider:       testClusterRKEConfigCloudProviderConf,
		ClusterName:         "test",
		DNS:                 testClusterRKEConfigDNSConf,
		EnableCRIDockerd:    newTrue(),
		IgnoreDockerVersion: newTrue(),
		Ingress:             testClusterRKEConfigIngressConf,
		Monitoring:          testClusterRKEConfigMonitoringConf,
		Network:             testClusterRKEConfigNetworkConfCanal,
		Nodes:               testClusterRKEConfigNodesConf,
		PrefixPath:          "terraform-test",
		WindowsPrefixPath:   "c:\\terraform-test",
		PrivateRegistries:   testClusterRKEConfigPrivateRegistriesConf,
		Services:            testClusterRKEConfigServicesConf,
		SSHAgentAuth:        true,
		SSHKeyPath:          "/home/user/.ssh",
		UpgradeStrategy:     testClusterRKEConfigNodeUpgradeStrategyConf,
		Version:             k8sVersion,
	}
	testClusterRKEConfigInterface = []interface{}{
		map[string]interface{}{
			"addon_job_timeout":     30,
			"addons":                "addons",
			"addons_include":        []interface{}{"addon1", "addon2"},
			"authentication":        testClusterRKEConfigAuthenticationInterface,
			"authorization":         testClusterRKEConfigAuthorizationInterface,
			"bastion_host":          testClusterRKEConfigBastionHostInterface,
			"cloud_provider":        testClusterRKEConfigCloudProviderInterface,
			"dns":                   testClusterRKEConfigDNSInterface,
			"enable_cri_dockerd":    newTrue(),
			"ignore_docker_version": newTrue(),
			"ingress":               testClusterRKEConfigIngressInterface,
			"monitoring":            testClusterRKEConfigMonitoringInterface,
			"network":               testClusterRKEConfigNetworkInterfaceCanal,
			"nodes":                 testClusterRKEConfigNodesInterface,
			"prefix_path":           "terraform-test",
			"win_prefix_path":       "c:\\terraform-test",
			"private_registries":    testClusterRKEConfigPrivateRegistriesInterface,
			"services":              testClusterRKEConfigServicesInterface,
			"ssh_agent_auth":        true,
			"ssh_key_path":          "/home/user/.ssh",
			"upgrade_strategy":      testClusterRKEConfigNodeUpgradeStrategyInterface,
			"kubernetes_version":    k8sVersion,
		},
	}
}

func TestFlattenClusterRKEConfig(t *testing.T) {
	testClusterRKEConfig()
	cases := []struct {
		Input          *managementClient.RancherKubernetesEngineConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigConf,
			testClusterRKEConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfig(tc.Input, testClusterRKEConfigInterface)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterRKEConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.RancherKubernetesEngineConfig
	}{
		{
			testClusterRKEConfigInterface,
			testClusterRKEConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfig(tc.Input, "test")
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
