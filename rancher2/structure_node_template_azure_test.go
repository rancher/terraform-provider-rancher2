package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"reflect"
	"testing"
)

var (
	testAzureNodeTemplateNodeTaintsConf      []managementClient.Taint
	testAzureNodeTemplateNodeTaintsInterface interface{}
	testNodeTemplateAzureConf                azureConfig
	testAzureNodeTemplateInterface           map[string]interface{}
	testAzureNodeTemplateConf                *NodeTemplate
	testNodeTemplateSquashAzureConfInterface map[string]interface{}
	testNodeTemplateExpandAzureConfInterface map[string]interface{}
)

func init() {
	testAzureNodeTemplateNodeTaintsConf = []managementClient.Taint{
		{
			Key:       "key",
			Value:     "value",
			Effect:    "recipient",
			TimeAdded: "time_added",
		},
	}
	testAzureNodeTemplateNodeTaintsInterface = []interface{}{
		map[string]interface{}{
			"key":        "key",
			"value":      "value",
			"effect":     "recipient",
			"time_added": "time_added",
		},
	}
	testNodeTemplateAzureConf = azureConfig{
		AvailabilitySet:        "docker-machine",
		ClientID:               "test-id",
		ClientSecret:           "test-secret",
		CustomData:             "test-data",
		DiskSize:               "32gb",
		DNS:                    "1.1.1.1",
		Environment:            "prod",
		FaultDomainCount:       "3",
		Image:                  "ubuntu",
		Location:               "us-east-1",
		ManagedDisks:           false,
		NoPublicIP:             false,
		NSG:                    "test-nsg",
		Plan:                   "test-plan",
		OpenPort:               nil,
		PrivateAddressOnly:     false,
		PrivateIPAddress:       "1.1.1.1",
		ResourceGroup:          "test-rg",
		Size:                   "large",
		SSHUser:                "user",
		StaticPublicIP:         false,
		StorageType:            "hdd",
		Subnet:                 "1.1.1.1",
		SubnetPrefix:           "1.1.1.1",
		SubscriptionID:         "sub-id",
		UpdateDomainCount:      "3",
		UsePrivateIP:           false,
		Vnet:                   "test-vnet",
		Tags:                   "key1,value1",
		UsePublicIPStandardSKU: true,
		AvailabilityZone:       "1",
		AcceleratedNetworking:  true,
	}
	testAzureNodeTemplateInterface = map[string]interface{}{
		"availability_set":           "docker-machine",
		"client_id":                  "test-id",
		"client_secret":              "test-secret",
		"custom_data":                "test-data",
		"disk_size":                  "32gb",
		"dns":                        "1.1.1.1",
		"environment":                "prod",
		"fault_domain-count":         "3",
		"image":                      "ubuntu",
		"location":                   "us-east-1",
		"managed_disks":              false,
		"no_public_ip":               false,
		"nsg":                        "test-nsg",
		"plan":                       "test-plan",
		"open_port":                  nil,
		"private_address_only":       false,
		"private_ip_address":         "1.1.1.1",
		"resource_group":             "test-rg",
		"size":                       "large",
		"ssh_user":                   "user",
		"static_public_ip":           false,
		"storage_type":               "hdd",
		"subnet":                     "1.1.1.1",
		"subnet_prefix":              "1.1.1.1",
		"subscription_id":            "sub-id",
		"update_domain_count":        "3",
		"use_private_ip":             false,
		"vnet":                       "test-vnet",
		"tags":                       "key1,value1",
		"use_public_ip_standard_sku": true,
		"availability_zone":          "1",
		"accelerated_networking":     true,
	}
	testAzureNodeTemplateAnnotationsConf := map[string]string{
		"key": "value",
	}
	testAzureNodeTemplateAnnotationsInterface := map[string]interface{}{
		"key": "value",
	}

	useInternalIP := false
	testAzureNodeTemplateConf = &NodeTemplate{
		NodeTemplate: managementClient.NodeTemplate{
			Driver:               "azure",
			UseInternalIPAddress: &useInternalIP,
			Annotations:          testAzureNodeTemplateAnnotationsConf,
			CloudCredentialID:    "abc-test-123",
			NodeTaints:           testAzureNodeTemplateNodeTaintsConf,
			EngineInstallURL:     "http://fake.url",
			Name:                 "test-node-template",
		},
		AzureConfig: &testNodeTemplateAzureConf,
	}

	testNodeTemplateSquashAzureConfInterface = map[string]interface{}{
		"annotations":             testAzureNodeTemplateAnnotationsInterface,
		"driver":                  "azure",
		"cloud_credential_id":     "abc-test-123",
		"use_internal_ip_address": useInternalIP,
		"engine_install_url":      "http://fake.url",
		"name":                    "test-node-template",
	}

	testNodeTemplateExpandAzureConfInterface = map[string]interface{}{
		"annotations":             testAzureNodeTemplateAnnotationsInterface,
		"node_taints":             testAzureNodeTemplateNodeTaintsInterface,
		"driver":                  "azure",
		"cloud_credential_id":     "abc-test-123",
		"use_internal_ip_address": useInternalIP,
		"engine_install_url":      "http://fake.url",
		"name":                    "test-node-template",
		"azure_config":            []interface{}{testAzureNodeTemplateInterface},
	}
}

func TestFlattenAzureNodeTemplate(t *testing.T) {
	cases := []struct {
		Input          *NodeTemplate
		ExpectedOutput map[string]interface{}
	}{
		{
			testAzureNodeTemplateConf,
			testNodeTemplateSquashAzureConfInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, nodeTemplateFields(), map[string]interface{}{})
		err := flattenNodeTemplate(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven: %#v", expectedOutput, tc.ExpectedOutput)
		}
	}
}

func TestExpandAzureNodeTemplate(t *testing.T) {
	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *NodeTemplate
	}{
		{
			Input:          testNodeTemplateExpandAzureConfInterface,
			ExpectedOutput: testAzureNodeTemplateConf,
		},
	}

	for _, tc := range cases {
		inputData := schema.TestResourceDataRaw(t, nodeTemplateFields(), tc.Input)
		output := expandNodeTemplate(inputData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven: %#v", tc.ExpectedOutput, tc.ExpectedOutput)
		}
	}
}
