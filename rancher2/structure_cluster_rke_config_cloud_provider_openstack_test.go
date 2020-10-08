package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigCloudProviderOpenstackBlockStorageConf      *managementClient.BlockStorageOpenstackOpts
	testClusterRKEConfigCloudProviderOpenstackBlockStorageInterface []interface{}
	testClusterRKEConfigCloudProviderOpenstackGlobalConf            *managementClient.GlobalOpenstackOpts
	testClusterRKEConfigCloudProviderOpenstackGlobalInterface       []interface{}
	testClusterRKEConfigCloudProviderOpenstackLoadBalancerConf      *managementClient.LoadBalancerOpenstackOpts
	testClusterRKEConfigCloudProviderOpenstackLoadBalancerInterface []interface{}
	testClusterRKEConfigCloudProviderOpenstackMetadataConf          *managementClient.MetadataOpenstackOpts
	testClusterRKEConfigCloudProviderOpenstackMetadataInterface     []interface{}
	testClusterRKEConfigCloudProviderOpenstackRouteConf             *managementClient.RouteOpenstackOpts
	testClusterRKEConfigCloudProviderOpenstackRouteInterface        []interface{}
	testClusterRKEConfigCloudProviderOpenstackConf                  *managementClient.OpenstackCloudProvider
	testClusterRKEConfigCloudProviderOpenstackInterface             []interface{}
)

func init() {
	testClusterRKEConfigCloudProviderOpenstackBlockStorageConf = &managementClient.BlockStorageOpenstackOpts{
		BSVersion:       "test",
		IgnoreVolumeAZ:  true,
		TrustDevicePath: true,
	}
	testClusterRKEConfigCloudProviderOpenstackBlockStorageInterface = []interface{}{
		map[string]interface{}{
			"bs_version":        "test",
			"ignore_volume_az":  true,
			"trust_device_path": true,
		},
	}
	testClusterRKEConfigCloudProviderOpenstackGlobalConf = &managementClient.GlobalOpenstackOpts{
		AuthURL:    "auth.terraform.test",
		Password:   "XXXXXXXX",
		TenantID:   "YYYYYYYY",
		Username:   "user",
		CAFile:     "ca_file",
		DomainID:   "domain_id",
		DomainName: "domain_name",
		Region:     "region",
		TenantName: "tenant",
		TrustID:    "VVVVVVVV",
	}
	testClusterRKEConfigCloudProviderOpenstackGlobalInterface = []interface{}{
		map[string]interface{}{
			"auth_url":    "auth.terraform.test",
			"password":    "XXXXXXXX",
			"tenant_id":   "YYYYYYYY",
			"username":    "user",
			"ca_file":     "ca_file",
			"domain_id":   "domain_id",
			"domain_name": "domain_name",
			"region":      "region",
			"tenant_name": "tenant",
			"trust_id":    "VVVVVVVV",
		},
	}
	testClusterRKEConfigCloudProviderOpenstackLoadBalancerConf = &managementClient.LoadBalancerOpenstackOpts{
		CreateMonitor:        true,
		FloatingNetworkID:    "test",
		LBMethod:             "method",
		LBProvider:           "provider",
		LBVersion:            "version",
		ManageSecurityGroups: true,
		MonitorDelay:         "30s",
		MonitorMaxRetries:    5,
		MonitorTimeout:       "10s",
		SubnetID:             "subnet",
		UseOctavia:           true,
	}
	testClusterRKEConfigCloudProviderOpenstackLoadBalancerInterface = []interface{}{
		map[string]interface{}{
			"create_monitor":         true,
			"floating_network_id":    "test",
			"lb_method":              "method",
			"lb_provider":            "provider",
			"lb_version":             "version",
			"manage_security_groups": true,
			"monitor_delay":          "30s",
			"monitor_max_retries":    5,
			"monitor_timeout":        "10s",
			"subnet_id":              "subnet",
			"use_octavia":            true,
		},
	}
	testClusterRKEConfigCloudProviderOpenstackMetadataConf = &managementClient.MetadataOpenstackOpts{
		RequestTimeout: 30,
		SearchOrder:    "order",
	}
	testClusterRKEConfigCloudProviderOpenstackMetadataInterface = []interface{}{
		map[string]interface{}{
			"request_timeout": 30,
			"search_order":    "order",
		},
	}
	testClusterRKEConfigCloudProviderOpenstackRouteConf = &managementClient.RouteOpenstackOpts{
		RouterID: "test",
	}
	testClusterRKEConfigCloudProviderOpenstackRouteInterface = []interface{}{
		map[string]interface{}{
			"router_id": "test",
		},
	}
	testClusterRKEConfigCloudProviderOpenstackConf = &managementClient.OpenstackCloudProvider{
		BlockStorage: testClusterRKEConfigCloudProviderOpenstackBlockStorageConf,
		Global:       testClusterRKEConfigCloudProviderOpenstackGlobalConf,
		LoadBalancer: testClusterRKEConfigCloudProviderOpenstackLoadBalancerConf,
		Metadata:     testClusterRKEConfigCloudProviderOpenstackMetadataConf,
		Route:        testClusterRKEConfigCloudProviderOpenstackRouteConf,
	}
	testClusterRKEConfigCloudProviderOpenstackInterface = []interface{}{
		map[string]interface{}{
			"block_storage": testClusterRKEConfigCloudProviderOpenstackBlockStorageInterface,
			"global":        testClusterRKEConfigCloudProviderOpenstackGlobalInterface,
			"load_balancer": testClusterRKEConfigCloudProviderOpenstackLoadBalancerInterface,
			"metadata":      testClusterRKEConfigCloudProviderOpenstackMetadataInterface,
			"route":         testClusterRKEConfigCloudProviderOpenstackRouteInterface,
		},
	}
}

func TestFlattenClusterRKEConfigCloudProviderOpenstackBlockStorage(t *testing.T) {

	cases := []struct {
		Input          *managementClient.BlockStorageOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackBlockStorageConf,
			testClusterRKEConfigCloudProviderOpenstackBlockStorageInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderOpenstackBlockStorage(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderOpenstackGlobal(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GlobalOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackGlobalConf,
			testClusterRKEConfigCloudProviderOpenstackGlobalInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderOpenstackGlobal(tc.Input, testClusterRKEConfigCloudProviderOpenstackGlobalInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderOpenstackLoadBalancer(t *testing.T) {

	cases := []struct {
		Input          *managementClient.LoadBalancerOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackLoadBalancerConf,
			testClusterRKEConfigCloudProviderOpenstackLoadBalancerInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderOpenstackLoadBalancer(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderOpenstackMetadata(t *testing.T) {

	cases := []struct {
		Input          *managementClient.MetadataOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackMetadataConf,
			testClusterRKEConfigCloudProviderOpenstackMetadataInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderOpenstackMetadata(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderOpenstackRoute(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RouteOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackRouteConf,
			testClusterRKEConfigCloudProviderOpenstackRouteInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderOpenstackRoute(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderOpenstack(t *testing.T) {

	cases := []struct {
		Input          *managementClient.OpenstackCloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackConf,
			testClusterRKEConfigCloudProviderOpenstackInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderOpenstack(tc.Input, testClusterRKEConfigCloudProviderOpenstackInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderOpenstackBlockStorage(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.BlockStorageOpenstackOpts
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackBlockStorageInterface,
			testClusterRKEConfigCloudProviderOpenstackBlockStorageConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderOpenstackBlockStorage(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderOpenstackGlobal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GlobalOpenstackOpts
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackGlobalInterface,
			testClusterRKEConfigCloudProviderOpenstackGlobalConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderOpenstackGlobal(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderOpenstackLoadBalancer(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.LoadBalancerOpenstackOpts
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackLoadBalancerInterface,
			testClusterRKEConfigCloudProviderOpenstackLoadBalancerConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderOpenstackLoadBalancer(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderOpenstackMetadata(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.MetadataOpenstackOpts
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackMetadataInterface,
			testClusterRKEConfigCloudProviderOpenstackMetadataConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderOpenstackMetadata(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderOpenstackRoute(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.RouteOpenstackOpts
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackRouteInterface,
			testClusterRKEConfigCloudProviderOpenstackRouteConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderOpenstackRoute(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderOpenstack(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.OpenstackCloudProvider
	}{
		{
			testClusterRKEConfigCloudProviderOpenstackInterface,
			testClusterRKEConfigCloudProviderOpenstackConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderOpenstack(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
