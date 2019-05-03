package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

var (
	testCloudCredentialAmazonec2Conf         *amazonec2CredentialConfig
	testCloudCredentialAmazonec2Interface    []interface{}
	testCloudCredentialAzureConf             *azureCredentialConfig
	testCloudCredentialAzureInterface        []interface{}
	testCloudCredentialDigitaloceanConf      *digitaloceanCredentialConfig
	testCloudCredentialDigitaloceanInterface []interface{}
	testCloudCredentialConfAmazonec2         *CloudCredential
	testCloudCredentialInterfaceAmazonec2    map[string]interface{}
	testCloudCredentialConfAzure             *CloudCredential
	testCloudCredentialInterfaceAzure        map[string]interface{}
	testCloudCredentialConfDigitalocean      *CloudCredential
	testCloudCredentialInterfaceDigitalocean map[string]interface{}
)

func init() {
	testCloudCredentialAmazonec2Conf = &amazonec2CredentialConfig{
		AccessKey: "access_key",
		SecretKey: "secret_key",
	}
	testCloudCredentialAmazonec2Interface = []interface{}{
		map[string]interface{}{
			"access_key": "access_key",
			"secret_key": "secret_key",
		},
	}
	testCloudCredentialAzureConf = &azureCredentialConfig{
		ClientID:       "client_id",
		ClientSecret:   "client_secret",
		SubscriptionID: "subscription_id",
	}
	testCloudCredentialAzureInterface = []interface{}{
		map[string]interface{}{
			"client_id":       "client_id",
			"client_secret":   "client_secret",
			"subscription_id": "subscription_id",
		},
	}
	testCloudCredentialDigitaloceanConf = &digitaloceanCredentialConfig{
		AccessToken: "access_token",
	}
	testCloudCredentialDigitaloceanInterface = []interface{}{
		map[string]interface{}{
			"access_token": "access_token",
		},
	}
	testCloudCredentialConfAmazonec2 = &CloudCredential{
		Amazonec2CredentialConfig: testCloudCredentialAmazonec2Conf,
	}
	testCloudCredentialConfAmazonec2.Name = "cloudCredential-test"
	testCloudCredentialConfAmazonec2.Description = "description"
	testCloudCredentialInterfaceAmazonec2 = map[string]interface{}{
		"name":                        "cloudCredential-test",
		"description":                 "description",
		"amazonec2_credential_config": testCloudCredentialAmazonec2Interface,
		"driver":                      amazonec2ConfigDriver,
	}
	testCloudCredentialConfAzure = &CloudCredential{
		AzureCredentialConfig: testCloudCredentialAzureConf,
	}
	testCloudCredentialConfAzure.Name = "cloudCredential-test"
	testCloudCredentialConfAzure.Description = "description"
	testCloudCredentialInterfaceAzure = map[string]interface{}{
		"name":                    "cloudCredential-test",
		"description":             "description",
		"azure_credential_config": testCloudCredentialAzureInterface,
		"driver":                  azureConfigDriver,
	}
	testCloudCredentialConfDigitalocean = &CloudCredential{
		DigitaloceanCredentialConfig: testCloudCredentialDigitaloceanConf,
	}
	testCloudCredentialConfDigitalocean.Name = "cloudCredential-test"
	testCloudCredentialConfDigitalocean.Description = "description"
	testCloudCredentialInterfaceDigitalocean = map[string]interface{}{
		"name":                           "cloudCredential-test",
		"description":                    "description",
		"digitalocean_credential_config": testCloudCredentialDigitaloceanInterface,
		"driver":                         digitaloceanConfigDriver,
	}
}

func TestFlattenCloudCredentialAmazonec2(t *testing.T) {

	cases := []struct {
		Input          *amazonec2CredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialAmazonec2Conf,
			testCloudCredentialAmazonec2Interface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialAmazonec2(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenCloudCredentialAzure(t *testing.T) {

	cases := []struct {
		Input          *azureCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialAzureConf,
			testCloudCredentialAzureInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialAzure(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenCloudCredentialDigitalocean(t *testing.T) {

	cases := []struct {
		Input          *digitaloceanCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialDigitaloceanConf,
			testCloudCredentialDigitaloceanInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialDigitalocean(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenCloudCredential(t *testing.T) {

	cases := []struct {
		Input          *CloudCredential
		ExpectedOutput map[string]interface{}
	}{
		{
			testCloudCredentialConfAmazonec2,
			testCloudCredentialInterfaceAmazonec2,
		},
		{
			testCloudCredentialConfAzure,
			testCloudCredentialInterfaceAzure,
		},
		{
			testCloudCredentialConfDigitalocean,
			testCloudCredentialInterfaceDigitalocean,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, cloudCredentialFields(), tc.ExpectedOutput)
		err := flattenCloudCredential(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandCloudCredentialAmazonec2(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *amazonec2CredentialConfig
	}{
		{
			testCloudCredentialAmazonec2Interface,
			testCloudCredentialAmazonec2Conf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialAmazonec2(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialAzure(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *azureCredentialConfig
	}{
		{
			testCloudCredentialAzureInterface,
			testCloudCredentialAzureConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialAzure(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialDigitalocean(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *digitaloceanCredentialConfig
	}{
		{
			testCloudCredentialDigitaloceanInterface,
			testCloudCredentialDigitaloceanConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialDigitalocean(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredential(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *CloudCredential
	}{
		{
			testCloudCredentialInterfaceAmazonec2,
			testCloudCredentialConfAmazonec2,
		},
		{
			testCloudCredentialInterfaceAzure,
			testCloudCredentialConfAzure,
		},
		{
			testCloudCredentialInterfaceDigitalocean,
			testCloudCredentialConfDigitalocean,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, cloudCredentialFields(), tc.Input)
		output := expandCloudCredential(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
