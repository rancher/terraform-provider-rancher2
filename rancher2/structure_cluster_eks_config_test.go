package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterEKSConfigConf      *managementClient.AmazonElasticContainerServiceConfig
	testClusterEKSConfigInterface []interface{}
)

func init() {
	testClusterEKSConfigConf = &managementClient.AmazonElasticContainerServiceConfig{
		AccessKey:                   "XXXXXXXX",
		SecretKey:                   "YYYYYYYY",
		AMI:                         "ami",
		AssociateWorkerNodePublicIP: newTrue(),
		InstanceType:                "instance",
		MaximumNodes:                5,
		MinimumNodes:                3,
		Region:                      "region",
		SecurityGroups:              []string{"sg1", "sg2"},
		ServiceRole:                 "role",
		Subnets:                     []string{"subnet1", "subnet2"},
		VirtualNetwork:              "network",
	}
	testClusterEKSConfigInterface = []interface{}{
		map[string]interface{}{
			"access_key":                      "XXXXXXXX",
			"secret_key":                      "YYYYYYYY",
			"ami":                             "ami",
			"associate_worker_node_public_ip": true,
			"instance_type":                   "instance",
			"maximum_nodes":                   5,
			"minimum_nodes":                   3,
			"region":                          "region",
			"security_groups":                 []interface{}{"sg1", "sg2"},
			"service_role":                    "role",
			"subnets":                         []interface{}{"subnet1", "subnet2"},
			"virtual_network":                 "network",
		},
	}
}

func TestFlattenClusterEKSConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AmazonElasticContainerServiceConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterEKSConfigConf,
			testClusterEKSConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterEKSConfig(tc.Input)
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
		ExpectedOutput *managementClient.AmazonElasticContainerServiceConfig
	}{
		{
			testClusterEKSConfigInterface,
			testClusterEKSConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterEKSConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
