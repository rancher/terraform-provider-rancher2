package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterEKSConfigV2NodeGroupConf      []managementClient.NodeGroup
	testClusterEKSConfigV2NodeGroupInterface []interface{}
	testClusterEKSConfigV2Conf               *managementClient.EKSClusterConfigSpec
	testClusterEKSConfigV2Interface          []interface{}
)

func init() {
	size := int64(3)
	testClusterEKSConfigV2NodeGroupConf = []managementClient.NodeGroup{
		{
			NodegroupName: "name",
			InstanceType:  "instance_type",
			DesiredSize:   &size,
			DiskSize:      &size,
			Ec2SshKey:     "ec2_ssh_key",
			Gpu:           newTrue(),
			Labels: map[string]string{
				"label1": "one",
				"label2": "two",
			},
			Tags: map[string]string{
				"tag1": "one",
				"tag2": "two",
			},
			MaxSize: &size,
			MinSize: &size,
			Subnets: []string{"net1", "net2"},
			Version: "kubernetes_version",
		},
	}
	testClusterEKSConfigV2NodeGroupInterface = []interface{}{
		map[string]interface{}{
			"name":          "name",
			"instance_type": "instance_type",
			"desired_size":  3,
			"disk_size":     3,
			"ec2_ssh_key":   "ec2_ssh_key",
			"gpu":           true,
			"labels": map[string]interface{}{
				"label1": "one",
				"label2": "two",
			},
			"tags": map[string]interface{}{
				"tag1": "one",
				"tag2": "two",
			},
			"max_size": 3,
			"min_size": 3,
		},
	}
	testClusterEKSConfigV2Conf = &managementClient.EKSClusterConfigSpec{
		AmazonCredentialSecret: "test",
		DisplayName:            "eksimport",
		LoggingTypes:           []string{"type1", "type2"},
		NodeGroups:             testClusterEKSConfigV2NodeGroupConf,
		Region:                 "test",
		KmsKey:                 "kms_key",
		Imported:               false,
		KubernetesVersion:      "kubernetes_version",
		PrivateAccess:          newTrue(),
		PublicAccess:           newTrue(),
		PublicAccessSources:    []string{"access1", "access2"},
		SecretsEncryption:      newTrue(),
		SecurityGroups:         []string{"sec1", "sec2"},
		ServiceRole:            "service_role",
		Subnets:                []string{"net1", "net2"},
		Tags: map[string]string{
			"value1": "one",
			"value2": "two",
		},
	}
	testClusterEKSConfigV2Interface = []interface{}{
		map[string]interface{}{
			"cloud_credential_id":   "test",
			"kubernetes_version":    "kubernetes_version",
			"imported":              false,
			"kms_key":               "kms_key",
			"logging_types":         []interface{}{"type1", "type2"},
			"node_groups":           testClusterEKSConfigV2NodeGroupInterface,
			"name":                  "eksimport",
			"private_access":        true,
			"public_access":         true,
			"public_access_sources": []interface{}{"access1", "access2"},
			"region":                "test",
			"secrets_encryption":    true,
			"security_groups":       []interface{}{"sec1", "sec2"},
			"service_role":          "service_role",
			"subnets":               []interface{}{"net1", "net2"},
			"tags": map[string]interface{}{
				"value1": "one",
				"value2": "two",
			},
		},
	}
}

func TestFlattenClusterEKSConfigV2NodeGroups(t *testing.T) {

	cases := []struct {
		Input          []managementClient.NodeGroup
		ExpectedOutput []interface{}
	}{
		{
			testClusterEKSConfigV2NodeGroupConf,
			testClusterEKSConfigV2NodeGroupInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterEKSConfigV2NodeGroups(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterEKSConfigV2(t *testing.T) {

	cases := []struct {
		Input          *managementClient.EKSClusterConfigSpec
		ExpectedOutput []interface{}
	}{
		{
			testClusterEKSConfigV2Conf,
			testClusterEKSConfigV2Interface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterEKSConfigV2(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterEKSConfigV2NodeGroups(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.NodeGroup
	}{
		{
			testClusterEKSConfigV2NodeGroupInterface,
			testClusterEKSConfigV2NodeGroupConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterEKSConfigV2NodeGroups(tc.Input, []string{"net1", "net2"}, "kubernetes_version")
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterEKSConfigV2(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.EKSClusterConfigSpec
	}{
		{
			testClusterEKSConfigV2Interface,
			testClusterEKSConfigV2Conf,
		},
	}

	for _, tc := range cases {
		output := expandClusterEKSConfigV2(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
