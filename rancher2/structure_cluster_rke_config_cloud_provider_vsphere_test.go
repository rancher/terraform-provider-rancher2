package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigCloudProviderVsphereDiskConf               *managementClient.DiskVsphereOpts
	testClusterRKEConfigCloudProviderVsphereDiskInterface          []interface{}
	testClusterRKEConfigCloudProviderVsphereGlobalConf             *managementClient.GlobalVsphereOpts
	testClusterRKEConfigCloudProviderVsphereGlobalInterface        []interface{}
	testClusterRKEConfigCloudProviderVsphereNetworkConf            *managementClient.NetworkVshpereOpts
	testClusterRKEConfigCloudProviderVsphereNetworkInterface       []interface{}
	testClusterRKEConfigCloudProviderVsphereVirtualCenterConf      map[string]managementClient.VirtualCenterConfig
	testClusterRKEConfigCloudProviderVsphereVirtualCenterInterface []interface{}
	testClusterRKEConfigCloudProviderVsphereWorkspaceConf          *managementClient.WorkspaceVsphereOpts
	testClusterRKEConfigCloudProviderVsphereWorkspaceInterface     []interface{}
	testClusterRKEConfigCloudProviderVsphereConf                   *managementClient.VsphereCloudProvider
	testClusterRKEConfigCloudProviderVsphereInterface              []interface{}
)

func init() {
	testClusterRKEConfigCloudProviderVsphereDiskConf = &managementClient.DiskVsphereOpts{
		SCSIControllerType: "test",
	}
	testClusterRKEConfigCloudProviderVsphereDiskInterface = []interface{}{
		map[string]interface{}{
			"scsi_controller_type": "test",
		},
	}
	testClusterRKEConfigCloudProviderVsphereGlobalConf = &managementClient.GlobalVsphereOpts{
		Datacenters:       "auth.terraform.test",
		InsecureFlag:      true,
		Password:          "YYYYYYYY",
		VCenterPort:       "123",
		User:              "user",
		RoundTripperCount: 10,
	}
	testClusterRKEConfigCloudProviderVsphereGlobalInterface = []interface{}{
		map[string]interface{}{
			"datacenters":          "auth.terraform.test",
			"insecure_flag":        true,
			"password":             "YYYYYYYY",
			"port":                 "123",
			"user":                 "user",
			"soap_roundtrip_count": 10,
		},
	}
	testClusterRKEConfigCloudProviderVsphereNetworkConf = &managementClient.NetworkVshpereOpts{
		PublicNetwork: "test",
	}
	testClusterRKEConfigCloudProviderVsphereNetworkInterface = []interface{}{
		map[string]interface{}{
			"public_network": "test",
		},
	}
	testClusterRKEConfigCloudProviderVsphereVirtualCenterConf = map[string]managementClient.VirtualCenterConfig{
		"test": {
			Datacenters:       "auth.terraform.test",
			Password:          "YYYYYYYY",
			VCenterPort:       "123",
			User:              "user",
			RoundTripperCount: 10,
		},
	}
	testClusterRKEConfigCloudProviderVsphereVirtualCenterInterface = []interface{}{
		map[string]interface{}{
			"name":                 "test",
			"datacenters":          "auth.terraform.test",
			"password":             "YYYYYYYY",
			"port":                 "123",
			"user":                 "user",
			"soap_roundtrip_count": 10,
		},
	}
	testClusterRKEConfigCloudProviderVsphereWorkspaceConf = &managementClient.WorkspaceVsphereOpts{
		Datacenter:       "test",
		Folder:           "folder",
		VCenterIP:        "vcenter",
		DefaultDatastore: "datastore",
		ResourcePoolPath: "resourcepool",
	}
	testClusterRKEConfigCloudProviderVsphereWorkspaceInterface = []interface{}{
		map[string]interface{}{
			"datacenter":        "test",
			"folder":            "folder",
			"server":            "vcenter",
			"default_datastore": "datastore",
			"resourcepool_path": "resourcepool",
		},
	}
	testClusterRKEConfigCloudProviderVsphereConf = &managementClient.VsphereCloudProvider{
		Disk:          testClusterRKEConfigCloudProviderVsphereDiskConf,
		Global:        testClusterRKEConfigCloudProviderVsphereGlobalConf,
		Network:       testClusterRKEConfigCloudProviderVsphereNetworkConf,
		VirtualCenter: testClusterRKEConfigCloudProviderVsphereVirtualCenterConf,
		Workspace:     testClusterRKEConfigCloudProviderVsphereWorkspaceConf,
	}
	testClusterRKEConfigCloudProviderVsphereInterface = []interface{}{
		map[string]interface{}{
			"disk":           testClusterRKEConfigCloudProviderVsphereDiskInterface,
			"global":         testClusterRKEConfigCloudProviderVsphereGlobalInterface,
			"network":        testClusterRKEConfigCloudProviderVsphereNetworkInterface,
			"virtual_center": testClusterRKEConfigCloudProviderVsphereVirtualCenterInterface,
			"workspace":      testClusterRKEConfigCloudProviderVsphereWorkspaceInterface,
		},
	}
}

func TestFlattenClusterRKEConfigCloudProviderVsphereDisk(t *testing.T) {

	cases := []struct {
		Input          *managementClient.DiskVsphereOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderVsphereDiskConf,
			testClusterRKEConfigCloudProviderVsphereDiskInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderVsphereDisk(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderVsphereGlobal(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GlobalVsphereOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderVsphereGlobalConf,
			testClusterRKEConfigCloudProviderVsphereGlobalInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderVsphereGlobal(tc.Input, testClusterRKEConfigCloudProviderVsphereGlobalInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderVsphereNetwork(t *testing.T) {

	cases := []struct {
		Input          *managementClient.NetworkVshpereOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderVsphereNetworkConf,
			testClusterRKEConfigCloudProviderVsphereNetworkInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderVsphereNetwork(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderVsphereVirtualCenter(t *testing.T) {

	cases := []struct {
		Input          map[string]managementClient.VirtualCenterConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderVsphereVirtualCenterConf,
			testClusterRKEConfigCloudProviderVsphereVirtualCenterInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderVsphereVirtualCenter(tc.Input, testClusterRKEConfigCloudProviderVsphereVirtualCenterInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderVsphereWorkspace(t *testing.T) {

	cases := []struct {
		Input          *managementClient.WorkspaceVsphereOpts
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderVsphereWorkspaceConf,
			testClusterRKEConfigCloudProviderVsphereWorkspaceInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderVsphereWorkspace(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigCloudProviderVsphere(t *testing.T) {

	cases := []struct {
		Input          *managementClient.VsphereCloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigCloudProviderVsphereConf,
			testClusterRKEConfigCloudProviderVsphereInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigCloudProviderVsphere(tc.Input, testClusterRKEConfigCloudProviderVsphereInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderVsphereDisk(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.DiskVsphereOpts
	}{
		{
			testClusterRKEConfigCloudProviderVsphereDiskInterface,
			testClusterRKEConfigCloudProviderVsphereDiskConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderVsphereDisk(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderVsphereGlobal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GlobalVsphereOpts
	}{
		{
			testClusterRKEConfigCloudProviderVsphereGlobalInterface,
			testClusterRKEConfigCloudProviderVsphereGlobalConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderVsphereGlobal(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderVsphereNetwork(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.NetworkVshpereOpts
	}{
		{
			testClusterRKEConfigCloudProviderVsphereNetworkInterface,
			testClusterRKEConfigCloudProviderVsphereNetworkConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderVsphereNetwork(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderVsphereVirtualCenter(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput map[string]managementClient.VirtualCenterConfig
	}{
		{
			testClusterRKEConfigCloudProviderVsphereVirtualCenterInterface,
			testClusterRKEConfigCloudProviderVsphereVirtualCenterConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderVsphereVirtualCenter(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderVsphereWorkspace(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.WorkspaceVsphereOpts
	}{
		{
			testClusterRKEConfigCloudProviderVsphereWorkspaceInterface,
			testClusterRKEConfigCloudProviderVsphereWorkspaceConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderVsphereWorkspace(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigCloudProviderVsphere(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.VsphereCloudProvider
	}{
		{
			testClusterRKEConfigCloudProviderVsphereInterface,
			testClusterRKEConfigCloudProviderVsphereConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigCloudProviderVsphere(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
