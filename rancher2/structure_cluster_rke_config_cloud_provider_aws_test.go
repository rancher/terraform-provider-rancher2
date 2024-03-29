package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterRKEConfigCloudProviderAwsGlobalConf               *managementClient.GlobalAwsOpts
	testClusterRKEConfigCloudProviderAwsGlobalInterface          []interface{}
	testClusterRKEConfigCloudProviderAwsServiceOverrideConf      map[string]managementClient.ServiceOverride
	testClusterRKEConfigCloudProviderAwsServiceOverrideInterface []interface{}
	testClusterRKEConfigCloudProviderAwsConf                     *managementClient.AWSCloudProvider
	testClusterRKEConfigCloudProviderAwsInterface                []interface{}
)

func init() {
	testClusterRKEConfigCloudProviderAwsGlobalConf = &managementClient.GlobalAwsOpts{
		DisableSecurityGroupIngress: true,
		DisableStrictZoneCheck:      true,
		ElbSecurityGroup:            "elb_group",
		KubernetesClusterID:         "k8s_id",
		KubernetesClusterTag:        "k8s_tag",
		RoleARN:                     "role_arn",
		RouteTableID:                "route_table_id",
		SubnetID:                    "subnet_id",
		VPC:                         "vpc",
		Zone:                        "zone",
	}
	testClusterRKEConfigCloudProviderAwsGlobalInterface = []interface{}{
		map[string]interface{}{
			"disable_security_group_ingress": true,
			"disable_strict_zone_check":      true,
			"elb_security_group":             "elb_group",
			"kubernetes_cluster_id":          "k8s_id",
			"kubernetes_cluster_tag":         "k8s_tag",
			"role_arn":                       "role_arn",
			"route_table_id":                 "route_table_id",
			"subnet_id":                      "subnet_id",
			"vpc":                            "vpc",
			"zone":                           "zone",
		},
	}
	testClusterRKEConfigCloudProviderAwsServiceOverrideConf = map[string]managementClient.ServiceOverride{
		"service": {
			Region:        "region",
			Service:       "service",
			SigningMethod: "signing_method",
			SigningName:   "signing_name",
			SigningRegion: "signing_region",
			URL:           "url",
		},
	}
	testClusterRKEConfigCloudProviderAwsServiceOverrideInterface = []interface{}{
		map[string]interface{}{
			"region":         "region",
			"service":        "service",
			"signing_method": "signing_method",
			"signing_name":   "signing_name",
			"signing_region": "signing_region",
			"url":            "url",
		},
	}
	testClusterRKEConfigCloudProviderAwsConf = &managementClient.AWSCloudProvider{
		Global:          testClusterRKEConfigCloudProviderAwsGlobalConf,
		ServiceOverride: testClusterRKEConfigCloudProviderAwsServiceOverrideConf,
	}
	testClusterRKEConfigCloudProviderAwsInterface = []interface{}{
		map[string]interface{}{
			"global":           testClusterRKEConfigCloudProviderAwsGlobalInterface,
			"service_override": testClusterRKEConfigCloudProviderAwsServiceOverrideInterface,
		},
	}
}

func TestFlattenClusterRKEConfigCloudProviderAwsGlobal(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GlobalAwsOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderAwsGlobalConf,
			testClusterRKEConfigCloudProviderAwsGlobalInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKEConfigCloudProviderAwsGlobal(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenClusterRKEConfigCloudProviderAwsServiceOverride(t *testing.T) {

	cases := []struct {
		Input          map[string]managementClient.ServiceOverride
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderAwsServiceOverrideConf,
			testClusterRKEConfigCloudProviderAwsServiceOverrideInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKEConfigCloudProviderAwsServiceOverride(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenClusterRKEConfigCloudProviderAws(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AWSCloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderAwsConf,
			testClusterRKEConfigCloudProviderAwsInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderAws(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterRKEConfigCloudProviderAwsGlobal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GlobalAwsOpts
	}{
		{
			testClusterRKEConfigCloudProviderAwsGlobalInterface,
			testClusterRKEConfigCloudProviderAwsGlobalConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterRKEConfigCloudProviderAwsGlobal(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandClusterRKEConfigCloudProviderAwsServiceOverride(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput map[string]managementClient.ServiceOverride
	}{
		{
			testClusterRKEConfigCloudProviderAwsServiceOverrideInterface,
			testClusterRKEConfigCloudProviderAwsServiceOverrideConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterRKEConfigCloudProviderAwsServiceOverride(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandClusterRKEConfigCloudProviderAws(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.AWSCloudProvider
	}{
		{
			testClusterRKEConfigCloudProviderAwsInterface,
			testClusterRKEConfigCloudProviderAwsConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderAws(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
