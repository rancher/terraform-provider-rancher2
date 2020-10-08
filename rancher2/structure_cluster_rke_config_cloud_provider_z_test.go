package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
