package rancher2

import (
	"encoding/json"
	"reflect"
	"testing"
)

var (
	testClusterEKSConfigConf      *AmazonElasticContainerServiceConfig
	testClusterEKSConfigInterface []interface{}
)

func init() {
	workerPoolBytes, _ := json.Marshal(
		AmazonElasticContainerWorkerPool{
			AddDefaultLabel: false,
			AddDefaultTaint: false,
			AdditionalLabels: map[string]string{
				"pool-name": "main-pool",
			},
			AdditionalTaints: []K8sTaint{
				{
					Key:      "pool-name",
					Operator: "Equal",
					Value:    "main-pool",
					Effect:   "NoSchedule",
				},
			},
			AMI:                         "ami",
			AssociateWorkerNodePublicIP: newTrue(),
			DesiredNodes:                4,
			InstanceType:                "instance",
			MaximumNodes:                5,
			MinimumNodes:                3,
			Name:                        "main-pool",
			NodeVolumeSize:              40,
			PlacementGroup:              "placement_group",
			UserData:                    "user_data",
			Subnets:                     []string{"worker1", "worker2"},
		},
	)

	testClusterEKSConfigConf = &AmazonElasticContainerServiceConfig{
		AccessKey:               "XXXXXXXX",
		SecretKey:               "YYYYYYYY",
		DisplayName:             "test",
		KeyPairName:             "key_pair_name",
		KubernetesVersion:       "1.11",
		ManageOwnSecurityGroups: newTrue(),
		NodeSecurityGroups:      []string{"node-sg1", "node-sg2"},
		Region:                  "region",
		SecurityGroups:          []string{"sg1", "sg2"},
		ServiceRole:             "role",
		SessionToken:            "session_token",
		Subnets:                 []string{"subnet1", "subnet2"},
		VirtualNetwork:          "network",
		WorkerPools: []string{
			string(workerPoolBytes),
		},
	}
	testClusterEKSConfigInterface = []interface{}{
		map[string]interface{}{
			"access_key":                 "XXXXXXXX",
			"secret_key":                 "YYYYYYYY",
			"key_pair_name":              "key_pair_name",
			"kubernetes_version":         "1.11",
			"manage_own_security_groups": true,
			"node_security_groups":       []interface{}{"node-sg1", "node-sg2"},
			"region":                     "region",
			"security_groups":            []interface{}{"sg1", "sg2"},
			"service_role":               "role",
			"session_token":              "session_token",
			"subnets":                    []interface{}{"subnet1", "subnet2"},
			"virtual_network":            "network",
			"worker_pools": []interface{}{
				map[string]interface{}{
					"add_default_label": false,
					"add_default_taint": false,
					"additional_labels": map[string]interface{}{
						"pool-name": "main-pool",
					},
					"additional_taints": []interface{}{
						map[string]interface{}{
							"key":      "pool-name",
							"operator": "Equal",
							"value":    "main-pool",
							"effect":   "NoSchedule",
						},
					},
					"ami":                             "ami",
					"associate_worker_node_public_ip": true,
					"desired_nodes":                   4,
					"ebs_encryption":                  false,
					"instance_type":                   "instance",
					"maximum_nodes":                   5,
					"minimum_nodes":                   3,
					"name":                            "main-pool",
					"node_volume_size":                40,
					"placement_group":                 "placement_group",
					"user_data":                       "user_data",
					"subnets":                         []interface{}{"worker1", "worker2"},
				},
			},
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
		ExpectedOutput *AmazonElasticContainerServiceConfig
	}{
		{
			testClusterEKSConfigInterface,
			testClusterEKSConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterEKSConfig(&AmazonElasticContainerServiceConfig{}, tc.Input, "test")
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
