package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterRKEConfigCloudProviderConfAzure          *managementClient.CloudProvider
	testClusterRKEConfigCloudProviderInterfaceAzure     []interface{}
	testClusterRKEConfigCloudProviderConfOpenstack      *managementClient.CloudProvider
	testClusterRKEConfigCloudProviderInterfaceOpenstack []interface{}
	testClusterRKEConfigCloudProviderConfVsphere        *managementClient.CloudProvider
	testClusterRKEConfigCloudProviderInterfaceVsphere   []interface{}
	testClusterRKEConfigCloudProviderConf               *managementClient.CloudProvider
	testClusterRKEConfigCloudProviderInterface          []interface{}
)

func init() {
	testClusterRKEConfigCloudProviderConfAzure = &managementClient.CloudProvider{
		AzureCloudProvider: testClusterRKEConfigCloudProviderAzureConf,
		Name:               "azure-test",
	}
	testClusterRKEConfigCloudProviderInterfaceAzure = []interface{}{
		map[string]interface{}{
			"azure_cloud_provider": testClusterRKEConfigCloudProviderAzureInterface,
			"name":                 "azure-test",
		},
	}
	testClusterRKEConfigCloudProviderConfOpenstack = &managementClient.CloudProvider{
		Name:                   "openstack-test",
		OpenstackCloudProvider: testClusterRKEConfigCloudProviderOpenstackConf,
	}
	testClusterRKEConfigCloudProviderInterfaceOpenstack = []interface{}{
		map[string]interface{}{
			"name":                     "openstack-test",
			"openstack_cloud_provider": testClusterRKEConfigCloudProviderOpenstackInterface,
		},
	}
	testClusterRKEConfigCloudProviderConfVsphere = &managementClient.CloudProvider{
		Name:                 "vsphere-test",
		VsphereCloudProvider: testClusterRKEConfigCloudProviderVsphereConf,
	}
	testClusterRKEConfigCloudProviderInterfaceVsphere = []interface{}{
		map[string]interface{}{
			"name":                   "vsphere-test",
			"vsphere_cloud_provider": testClusterRKEConfigCloudProviderVsphereInterface,
		},
	}
	testClusterRKEConfigCloudProviderConf = &managementClient.CloudProvider{
		CustomCloudProvider: "XXXXXXXXXXXX",
		Name:                "test",
	}
	testClusterRKEConfigCloudProviderInterface = []interface{}{
		map[string]interface{}{
			"custom_cloud_provider": "XXXXXXXXXXXX",
			"name":                  "test",
		},
	}
}

func TestFlattenClusterRKEConfigCloudProvider(t *testing.T) {

	cases := []struct {
		Input          *managementClient.CloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderConfAzure,
			testClusterRKEConfigCloudProviderInterfaceAzure,
		},
		{
			testClusterRKEConfigCloudProviderConfOpenstack,
			testClusterRKEConfigCloudProviderInterfaceOpenstack,
		},
		{
			testClusterRKEConfigCloudProviderConfVsphere,
			testClusterRKEConfigCloudProviderInterfaceVsphere,
		},
		{
			testClusterRKEConfigCloudProviderConf,
			testClusterRKEConfigCloudProviderInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProvider(tc.Input, tc.ExpectedOutput)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterRKEConfigCloudProvider(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.CloudProvider
	}{
		{
			testClusterRKEConfigCloudProviderInterfaceAzure,
			testClusterRKEConfigCloudProviderConfAzure,
		},
		{
			testClusterRKEConfigCloudProviderInterfaceOpenstack,
			testClusterRKEConfigCloudProviderConfOpenstack,
		},
		{
			testClusterRKEConfigCloudProviderInterfaceVsphere,
			testClusterRKEConfigCloudProviderConfVsphere,
		},
		{
			testClusterRKEConfigCloudProviderInterface,
			testClusterRKEConfigCloudProviderConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProvider(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
