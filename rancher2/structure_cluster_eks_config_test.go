package rancher2

import (
	"reflect"
	"testing"
)

var (
	testClusterEKSConfigConf      *AmazonElasticContainerServiceConfig
	testClusterEKSConfigInterface []interface{}
)

func init() {
	testClusterEKSConfigConf = &AmazonElasticContainerServiceConfig{
		AccessKey:                   "XXXXXXXX",
		SecretKey:                   "YYYYYYYY",
		AMI:                         "ami",
		AssociateWorkerNodePublicIP: newTrue(),
		DesiredNodes:                4,
		DisplayName:                 "test",
		DriverName:                  clusterDriverEKS,
		EBSEncryption:               true,
		InstanceType:                "instance",
		KeyPairName:                 "key_pair_name",
		KubernetesVersion:           "1.11",
		MaximumNodes:                5,
		MinimumNodes:                3,
		NodeVolumeSize:              40,
		Region:                      "region",
		SecurityGroups:              []string{"sg1", "sg2"},
		ServiceRole:                 "role",
		SessionToken:                "session_token",
		Subnets:                     []string{"subnet1", "subnet2"},
		UserData:                    "user_data",
		VirtualNetwork:              "network",
	}
	testClusterEKSConfigInterface = []interface{}{
		map[string]interface{}{
			"access_key":                      "XXXXXXXX",
			"secret_key":                      "YYYYYYYY",
			"ami":                             "ami",
			"associate_worker_node_public_ip": true,
			"desired_nodes":                   4,
			"ebs_encryption":                  true,
			"instance_type":                   "instance",
			"key_pair_name":                   "key_pair_name",
			"kubernetes_version":              "1.11",
			"maximum_nodes":                   5,
			"minimum_nodes":                   3,
			"node_volume_size":                40,
			"region":                          "region",
			"security_groups":                 []interface{}{"sg1", "sg2"},
			"service_role":                    "role",
			"session_token":                   "session_token",
			"subnets":                         []interface{}{"subnet1", "subnet2"},
			"user_data":                       "user_data",
			"virtual_network":                 "network",
		},
	}
}

func TestFlattenClusterEKSConfig(t *testing.T) {

	cases := []struct {
		Input          *AmazonElasticContainerServiceConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterEKSConfigConf,
			testClusterEKSConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterEKSConfig(tc.Input, testClusterEKSConfigInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterEKSConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *AmazonElasticContainerServiceConfig
	}{
		{
			testClusterEKSConfigInterface,
			testClusterEKSConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterEKSConfig(tc.Input, "test")
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
