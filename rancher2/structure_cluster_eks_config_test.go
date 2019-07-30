package rancher2

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	testClusterEKSConfigConf                  *AmazonElasticContainerServiceConfig
	testClusterEKSConfigInterface             []interface{}
	testClusterEKSConfigConfCredsFromEnv      *AmazonElasticContainerServiceConfig
	testClusterEKSConfigInterfaceCredsFromEnv []interface{}
)

func init() {
	testClusterEKSConfigConf = &AmazonElasticContainerServiceConfig{
		AccessKey:                   "XXXXXXXX",
		SecretKey:                   "YYYYYYYY",
		AMI:                         "ami",
		AssociateWorkerNodePublicIP: newTrue(),
		DisplayName:                 "test",
		InstanceType:                "instance",
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
			"instance_type":                   "instance",
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
	testClusterEKSConfigConfCredsFromEnv = &AmazonElasticContainerServiceConfig{
		AccessKey:                   "env_XXXXXXXX",
		SecretKey:                   "env_YYYYYYYY",
		SessionToken:                "env_session_token",
		AssociateWorkerNodePublicIP: newTrue(),
		DisplayName:                 "test",
	}
	testClusterEKSConfigInterfaceCredsFromEnv = []interface{}{
		map[string]interface{}{
			"aws_creds_from_env":              true,
			"associate_worker_node_public_ip": true,
		},
	}
}

func TestFlattenClusterEKSConfig(t *testing.T) {
	cases := []struct {
		Input          *AmazonElasticContainerServiceConfig
		Config         []interface{}
		ExpectedOutput []interface{}
	}{
		{
			testClusterEKSConfigConf,
			[]interface{}{},
			testClusterEKSConfigInterface,
		},
		{
			testClusterEKSConfigConfCredsFromEnv,
			testClusterEKSConfigInterfaceCredsFromEnv,
			testClusterEKSConfigInterfaceCredsFromEnv,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterEKSConfig(tc.Input, tc.Config)
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
		ExtraEnv       map[string]string
		ExpectedOutput *AmazonElasticContainerServiceConfig
		ExpectedError  error
	}{
		{
			testClusterEKSConfigInterface,
			map[string]string{},
			testClusterEKSConfigConf,
			nil,
		},
		{
			testClusterEKSConfigInterfaceCredsFromEnv,
			map[string]string{
				"AWS_ACCESS_KEY_ID":     "env_XXXXXXXX",
				"AWS_SECRET_ACCESS_KEY": "env_YYYYYYYY",
				"AWS_SESSION_TOKEN":     "env_session_token",
			},
			testClusterEKSConfigConfCredsFromEnv,
			nil,
		},
		{
			[]interface{}{
				map[string]interface{}{},
			},
			map[string]string{},
			&AmazonElasticContainerServiceConfig{},
			fmt.Errorf("[ERROR] 'aws_creds_from_env=false' or not set but 'access_key' not set"),
		},
		{
			[]interface{}{
				map[string]interface{}{
					"access_key": "XXXXXXXX",
				},
			},
			map[string]string{},
			&AmazonElasticContainerServiceConfig{},
			fmt.Errorf("[ERROR] 'aws_creds_from_env=false' or not set but 'secret_key' not set"),
		},
		{
			testClusterEKSConfigInterfaceCredsFromEnv,
			map[string]string{},
			testClusterEKSConfigConfCredsFromEnv,
			fmt.Errorf("[ERROR] 'aws_creds_from_env=true' but env var AWS_ACCESS_KEY_ID is not set"),
		},
		{
			testClusterEKSConfigInterfaceCredsFromEnv,
			map[string]string{
				"AWS_ACCESS_KEY_ID": "env_XXXXXXXX",
			},
			testClusterEKSConfigConfCredsFromEnv,
			fmt.Errorf("[ERROR] 'aws_creds_from_env=true' but env var AWS_SECRET_ACCESS_KEY is not set"),
		},
	}

	for _, tc := range cases {
		runWithEnv(tc.ExtraEnv, func() {
			output, err := expandClusterEKSConfig(tc.Input, "test")

			if tc.ExpectedError != nil {
				if !reflect.DeepEqual(err, tc.ExpectedError) {
					t.Fatalf("Unexpected error from expander.\nExpected: %#v\nGiven:    %#v",
						tc.ExpectedError, err)
				} else {
					return
				}
			}

			if err != nil {
				t.Fatalf("[ERROR] on expander: %#v", err)
			}
			if !reflect.DeepEqual(output, tc.ExpectedOutput) {
				t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
					tc.ExpectedOutput, output)
			}
		})
	}
}
