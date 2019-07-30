package rancher2

import (
	"fmt"
	"os"
)

// Flatteners

func flattenClusterEKSConfig(in *AmazonElasticContainerServiceConfig, p []interface{}) ([]interface{}, error) {
	if in == nil {
		return []interface{}{}, nil
	}

	var obj map[string]interface{}

	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if v, ok := obj["aws_creds_from_env"].(bool); !ok || !v {
		if len(in.AccessKey) > 0 {
			obj["access_key"] = in.AccessKey
		}

		if len(in.SecretKey) > 0 {
			obj["secret_key"] = in.SecretKey
		}

		if len(in.SessionToken) > 0 {
			obj["session_token"] = in.SessionToken
		}
	}

	if len(in.AMI) > 0 {
		obj["ami"] = in.AMI
	}

	obj["associate_worker_node_public_ip"] = *in.AssociateWorkerNodePublicIP

	if len(in.InstanceType) > 0 {
		obj["instance_type"] = in.InstanceType
	}

	if len(in.KubernetesVersion) > 0 {
		obj["kubernetes_version"] = in.KubernetesVersion
	}

	if in.MaximumNodes > 0 {
		obj["maximum_nodes"] = int(in.MaximumNodes)
	}

	if in.MinimumNodes > 0 {
		obj["minimum_nodes"] = int(in.MinimumNodes)
	}

	if in.NodeVolumeSize > 0 {
		obj["node_volume_size"] = int(in.NodeVolumeSize)
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.SecurityGroups) > 0 {
		obj["security_groups"] = toArrayInterface(in.SecurityGroups)
	}

	if len(in.ServiceRole) > 0 {
		obj["service_role"] = in.ServiceRole
	}

	if len(in.Subnets) > 0 {
		obj["subnets"] = toArrayInterface(in.Subnets)
	}

	if len(in.UserData) > 0 {
		obj["user_data"] = in.UserData
	}

	if len(in.VirtualNetwork) > 0 {
		obj["virtual_network"] = in.VirtualNetwork
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterEKSConfig(p []interface{}, name string) (*AmazonElasticContainerServiceConfig, error) {
	obj := &AmazonElasticContainerServiceConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	obj.DisplayName = name

	if v, ok := in["aws_creds_from_env"].(bool); ok && v {
		obj.AccessKey, ok = os.LookupEnv("AWS_ACCESS_KEY_ID")
		if !ok {
			return obj, fmt.Errorf("[ERROR] 'aws_creds_from_env=true' but env var AWS_ACCESS_KEY_ID is not set")
		}
		obj.SecretKey, ok = os.LookupEnv("AWS_SECRET_ACCESS_KEY")
		if !ok {
			return obj, fmt.Errorf("[ERROR] 'aws_creds_from_env=true' but env var AWS_SECRET_ACCESS_KEY is not set")
		}
		obj.SessionToken = os.Getenv("AWS_SESSION_TOKEN")

	} else {
		if v, ok := in["access_key"].(string); ok && len(v) > 0 {
			obj.AccessKey = v
		} else {
			return obj, fmt.Errorf("[ERROR] 'aws_creds_from_env=false' or not set but 'access_key' not set")
		}

		if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
			obj.SecretKey = v
		} else {
			return obj, fmt.Errorf("[ERROR] 'aws_creds_from_env=false' or not set but 'secret_key' not set")
		}

		if v, ok := in["session_token"].(string); ok && len(v) > 0 {
			obj.SessionToken = v
		}
	}

	if v, ok := in["ami"].(string); ok && len(v) > 0 {
		obj.AMI = v
	}

	if v, ok := in["associate_worker_node_public_ip"].(bool); ok {
		obj.AssociateWorkerNodePublicIP = &v
	}

	if v, ok := in["instance_type"].(string); ok && len(v) > 0 {
		obj.InstanceType = v
	}

	if v, ok := in["kubernetes_version"].(string); ok && len(v) > 0 {
		obj.KubernetesVersion = v
	}

	if v, ok := in["maximum_nodes"].(int); ok && v > 0 {
		obj.MaximumNodes = int64(v)
	}

	if v, ok := in["minimum_nodes"].(int); ok && v > 0 {
		obj.MinimumNodes = int64(v)
	}

	if v, ok := in["node_volume_size"].(int); ok && v > 0 {
		obj.NodeVolumeSize = int64(v)
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["security_groups"].([]interface{}); ok && len(v) > 0 {
		obj.SecurityGroups = toArrayString(v)
	}

	if v, ok := in["service_role"].(string); ok && len(v) > 0 {
		obj.ServiceRole = v
	}

	if v, ok := in["subnets"].([]interface{}); ok && len(v) > 0 {
		obj.Subnets = toArrayString(v)
	}

	if v, ok := in["user_data"].(string); ok && len(v) > 0 {
		obj.UserData = v
	}

	if v, ok := in["virtual_network"].(string); ok && len(v) > 0 {
		obj.VirtualNetwork = v
	}

	return obj, nil
}
