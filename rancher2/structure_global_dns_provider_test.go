package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAccRancher2GlobalDNSProviderAlidnsConf          *managementClient.GlobalDnsProvider
	testAccRancher2GlobalDNSProviderAlidnsInterface     []interface{}
	testAccRancher2GlobalDNSProviderCloudflareConf      *managementClient.GlobalDnsProvider
	testAccRancher2GlobalDNSProviderCloudflareInterface []interface{}
	testAccRancher2GlobalDNSProviderRoute53Conf         *managementClient.GlobalDnsProvider
	testAccRancher2GlobalDNSProviderRoute53Interface    []interface{}
)

func init() {
	testAccRancher2GlobalDNSProviderAlidnsConf = &managementClient.GlobalDnsProvider{
		Name:       "name",
		RootDomain: "root.domain",
		AlidnsProviderConfig: &managementClient.AlidnsProviderConfig{
			AccessKey: "XXXXXXXXXX",
			SecretKey: "YYYYYYYYYY",
		},
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testAccRancher2GlobalDNSProviderAlidnsInterface = []interface{}{
		map[string]interface{}{
			"name":        "name",
			"root_domain": "root.domain",
			"alidns_config": []interface{}{
				map[string]interface{}{
					"access_key": "XXXXXXXXXX",
					"secret_key": "YYYYYYYYYY",
				},
			},
			"annotations": map[string]interface{}{
				"node_one": "one",
				"node_two": "two",
			},
			"labels": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
		},
	}
	testAccRancher2GlobalDNSProviderCloudflareConf = &managementClient.GlobalDnsProvider{
		Name:       "name",
		RootDomain: "root.domain",
		CloudflareProviderConfig: &managementClient.CloudflareProviderConfig{
			APIEmail:     "XXXXXXXXXX",
			APIKey:       "YYYYYYYYYY",
			ProxySetting: newTrue(),
		},
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testAccRancher2GlobalDNSProviderCloudflareInterface = []interface{}{
		map[string]interface{}{
			"name":        "name",
			"root_domain": "root.domain",
			"cloudflare_config": []interface{}{
				map[string]interface{}{
					"api_email":     "XXXXXXXXXX",
					"api_key":       "YYYYYYYYYY",
					"proxy_setting": true,
				},
			},
			"annotations": map[string]interface{}{
				"node_one": "one",
				"node_two": "two",
			},
			"labels": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
		},
	}
	testAccRancher2GlobalDNSProviderRoute53Conf = &managementClient.GlobalDnsProvider{
		Name:       "name",
		RootDomain: "root.domain",
		Route53ProviderConfig: &managementClient.Route53ProviderConfig{
			AccessKey:       "XXXXXXXXXX",
			SecretKey:       "YYYYYYYYYY",
			CredentialsPath: "credPath",
			Region:          "region",
			RoleArn:         "role",
			ZoneType:        "private",
		},
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testAccRancher2GlobalDNSProviderRoute53Interface = []interface{}{
		map[string]interface{}{
			"name":        "name",
			"root_domain": "root.domain",
			"route53_config": []interface{}{
				map[string]interface{}{
					"access_key":       "XXXXXXXXXX",
					"secret_key":       "YYYYYYYYYY",
					"credentials_path": "credPath",
					"region":           "region",
					"role_arn":         "role",
					"zone_type":        "private",
				},
			},
			"annotations": map[string]interface{}{
				"node_one": "one",
				"node_two": "two",
			},
			"labels": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
		},
	}
}

func TestFlattenGlobalDnsProvider(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GlobalDnsProvider
		ExpectedOutput []interface{}
	}{
		{
			testAccRancher2GlobalDNSProviderAlidnsConf,
			testAccRancher2GlobalDNSProviderAlidnsInterface,
		},
		{
			testAccRancher2GlobalDNSProviderCloudflareConf,
			testAccRancher2GlobalDNSProviderCloudflareInterface,
		},
		{
			testAccRancher2GlobalDNSProviderRoute53Conf,
			testAccRancher2GlobalDNSProviderRoute53Interface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, globalDNSProviderFields(), tc.ExpectedOutput[0].(map[string]interface{}))
		err := flattenGlobalDNSProvider(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput[0].(map[string]interface{}) {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual([]interface{}{expectedOutput}, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, []interface{}{expectedOutput})
		}
	}
}

func TestExpandGlobalDnsProvider(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GlobalDnsProvider
	}{
		{
			testAccRancher2GlobalDNSProviderAlidnsInterface,
			testAccRancher2GlobalDNSProviderAlidnsConf,
		},
		{
			testAccRancher2GlobalDNSProviderCloudflareInterface,
			testAccRancher2GlobalDNSProviderCloudflareConf,
		},
		{
			testAccRancher2GlobalDNSProviderRoute53Interface,
			testAccRancher2GlobalDNSProviderRoute53Conf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, globalDNSProviderFields(), tc.Input[0].(map[string]interface{}))
		output := expandGlobalDNSProvider(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
