package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigCloudProviderAzureConf      *managementClient.AzureCloudProvider
	testClusterRKEConfigCloudProviderAzureInterface []interface{}
)

func init() {
	testClusterRKEConfigCloudProviderAzureConf = &managementClient.AzureCloudProvider{
		AADClientID:                  "XXXXXXXX",
		AADClientSecret:              "XXXXXXXXXXXX",
		SubscriptionID:               "YYYYYYYY",
		TenantID:                     "ZZZZZZZZ",
		AADClientCertPassword:        "password",
		AADClientCertPath:            "/home/user/.ssh",
		Cloud:                        "cloud",
		CloudProviderBackoff:         true,
		CloudProviderBackoffDuration: 30,
		CloudProviderBackoffExponent: 20,
		CloudProviderBackoffJitter:   10,
		CloudProviderBackoffRetries:  5,
		CloudProviderRateLimit:       true,
		CloudProviderRateLimitBucket: 15,
		CloudProviderRateLimitQPS:    100,
		LoadBalancerSku:              cloudProviderAzureLoadBalancerSkuStandard,
		Location:                     "location",
		MaximumLoadBalancerRuleCount: 150,
		PrimaryAvailabilitySetName:   "primary",
		PrimaryScaleSetName:          "primary_scale",
		ResourceGroup:                "resource_group",
		RouteTableName:               "route_table_name",
		SecurityGroupName:            "security_group_name",
		SubnetName:                   "subnet_name",
		UseInstanceMetadata:          true,
		UseManagedIdentityExtension:  true,
		VMType:                       "vm_type",
		VnetName:                     "vnet_name",
		VnetResourceGroup:            "vnet_resource_group",
	}
	testClusterRKEConfigCloudProviderAzureInterface = []interface{}{
		map[string]interface{}{
			"aad_client_id":                    "XXXXXXXX",
			"aad_client_secret":                "XXXXXXXXXXXX",
			"subscription_id":                  "YYYYYYYY",
			"tenant_id":                        "ZZZZZZZZ",
			"aad_client_cert_password":         "password",
			"aad_client_cert_path":             "/home/user/.ssh",
			"cloud":                            "cloud",
			"cloud_provider_backoff":           true,
			"cloud_provider_backoff_duration":  30,
			"cloud_provider_backoff_exponent":  20,
			"cloud_provider_backoff_jitter":    10,
			"cloud_provider_backoff_retries":   5,
			"cloud_provider_rate_limit":        true,
			"cloud_provider_rate_limit_bucket": 15,
			"cloud_provider_rate_limit_qps":    100,
			"load_balancer_sku":                cloudProviderAzureLoadBalancerSkuStandard,
			"location":                         "location",
			"maximum_load_balancer_rule_count": 150,
			"primary_availability_set_name":    "primary",
			"primary_scale_set_name":           "primary_scale",
			"resource_group":                   "resource_group",
			"route_table_name":                 "route_table_name",
			"security_group_name":              "security_group_name",
			"subnet_name":                      "subnet_name",
			"use_instance_metadata":            true,
			"use_managed_identity_extension":   true,
			"vm_type":                          "vm_type",
			"vnet_name":                        "vnet_name",
			"vnet_resource_group":              "vnet_resource_group",
		},
	}
}

func TestFlattenClusterRKEConfigCloudProviderAzure(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AzureCloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderAzureConf,
			testClusterRKEConfigCloudProviderAzureInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderAzure(tc.Input, testClusterRKEConfigCloudProviderAzureInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderAzure(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.AzureCloudProvider
	}{
		{
			testClusterRKEConfigCloudProviderAzureInterface,
			testClusterRKEConfigCloudProviderAzureConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderAzure(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
