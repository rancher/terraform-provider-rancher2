package rancher2

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Flatteners

func flattenClusterEKSConfig(in *AmazonElasticContainerServiceConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	if len(in.KeyPairName) > 0 {
		obj["key_pair_name"] = in.KeyPairName
	}

	if len(in.KubernetesVersion) > 0 {
		obj["kubernetes_version"] = in.KubernetesVersion
	}

	obj["manage_own_security_groups"] = *in.ManageOwnSecurityGroups

	if len(in.NodeSecurityGroups) > 0 {
		obj["node_security_groups"] = toArrayInterface(in.NodeSecurityGroups)
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

	if len(in.SessionToken) > 0 {
		obj["session_token"] = in.SessionToken
	}

	if len(in.Subnets) > 0 {
		obj["subnets"] = toArrayInterface(in.Subnets)
	}

	if len(in.VirtualNetwork) > 0 {
		obj["virtual_network"] = in.VirtualNetwork
	}

	var workerPoolObjs []interface{}
	for _, workerPoolIn := range in.WorkerPools {
		var workerPool AmazonElasticContainerWorkerPool
		if err := json.Unmarshal([]byte(workerPoolIn), &workerPool); err != nil {
			return nil, err
		}

		workerPoolObj := make(map[string]interface{})

		workerPoolObj["add_default_label"] = workerPool.AddDefaultLabel
		workerPoolObj["add_default_taint"] = workerPool.AddDefaultTaint

		if len(workerPool.AdditionalLabels) > 0 {
			additionalLabelsObj := make(map[string]interface{})
			for key, value := range workerPool.AdditionalLabels {
				additionalLabelsObj[key] = value
			}
			workerPoolObj["additional_labels"] = additionalLabelsObj
		}

		if len(workerPool.AdditionalTaints) > 0 {
			additionalTaintObjs := make([]interface{}, 0, len(workerPool.AdditionalTaints))
			for _, taint := range workerPool.AdditionalTaints {
				additionalTaintObj := make(map[string]interface{})

				if len(taint.Effect) > 0 {
					additionalTaintObj["effect"] = taint.Effect
				}

				if len(taint.Key) > 0 {
					additionalTaintObj["key"] = taint.Key
				}

				if len(taint.Operator) > 0 {
					additionalTaintObj["operator"] = taint.Operator
				}

				if len(taint.Value) > 0 {
					additionalTaintObj["value"] = taint.Value
				}

				additionalTaintObjs = append(additionalTaintObjs, additionalTaintObj)
			}

			workerPoolObj["additional_taints"] = additionalTaintObjs
		}

		if len(workerPool.AMI) > 0 {
			workerPoolObj["ami"] = workerPool.AMI
		}

		workerPoolObj["associate_worker_node_public_ip"] = *workerPool.AssociateWorkerNodePublicIP

		if workerPool.DesiredNodes > 0 {
			workerPoolObj["desired_nodes"] = int(workerPool.DesiredNodes)
		}

		workerPoolObj["ebs_encryption"] = workerPool.EBSEncryption

		if len(workerPool.InstanceType) > 0 {
			workerPoolObj["instance_type"] = workerPool.InstanceType
		}

		if len(workerPool.Name) > 0 {
			workerPoolObj["name"] = workerPool.Name
		}

		if workerPool.MaximumNodes > 0 {
			workerPoolObj["maximum_nodes"] = int(workerPool.MaximumNodes)
		}

		if workerPool.MinimumNodes > 0 {
			workerPoolObj["minimum_nodes"] = int(workerPool.MinimumNodes)
		}

		if workerPool.NodeVolumeSize > 0 {
			workerPoolObj["node_volume_size"] = int(workerPool.NodeVolumeSize)
		}

		if len(workerPool.PlacementGroup) > 0 {
			workerPoolObj["placement_group"] = workerPool.PlacementGroup
		}

		if len(workerPool.UserData) > 0 {
			workerPoolObj["user_data"] = workerPool.UserData
		}

		if len(workerPool.Subnets) > 0 {
			workerPoolObj["subnets"] = toArrayInterface(workerPool.Subnets)
		}

		workerPoolObjs = append(workerPoolObjs, workerPoolObj)
	}

	// when rancher returns details of a cluster that hasn't been migrated to new state model, we fallback to old
	// fields to extract existing worker pool details
	if len(workerPoolObjs) == 0 {
		workerPoolObj := make(map[string]interface{})
		workerPoolObj["name"] = "main-pool"

		if len(in.AMI) > 0 {
			workerPoolObj["ami"] = in.AMI
		}

		workerPoolObj["associate_worker_node_public_ip"] = *in.AssociateWorkerNodePublicIP

		if in.DesiredNodes > 0 {
			workerPoolObj["desired_nodes"] = int(in.DesiredNodes)
		}

		workerPoolObj["ebs_encryption"] = in.EBSEncryption

		if len(in.InstanceType) > 0 {
			workerPoolObj["instance_type"] = in.InstanceType
		}

		if in.MaximumNodes > 0 {
			workerPoolObj["maximum_nodes"] = int(in.MaximumNodes)
		}

		if in.MinimumNodes > 0 {
			workerPoolObj["minimum_nodes"] = int(in.MinimumNodes)
		}

		if in.NodeVolumeSize > 0 {
			workerPoolObj["node_volume_size"] = int(in.NodeVolumeSize)
		}

		if len(in.PlacementGroup) > 0 {
			workerPoolObj["placement_group"] = in.PlacementGroup
		}

		if len(in.UserData) > 0 {
			workerPoolObj["user_data"] = in.UserData
		}

		if len(in.WorkerSubnets) > 0 {
			workerPoolObj["subnets"] = toArrayInterface(in.WorkerSubnets)
		}

		workerPoolObjs = append(workerPoolObjs, workerPoolObj)
	}

	obj["worker_pools"] = workerPoolObjs
	return []interface{}{obj}, nil
}

// Expanders

func expandClusterEKSConfig(obj *AmazonElasticContainerServiceConfig, p []interface{}, name string) (*AmazonElasticContainerServiceConfig, error) {
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	obj.DisplayName = name

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	if v, ok := in["key_pair_name"].(string); ok && len(v) > 0 {
		obj.KeyPairName = v
	}

	if v, ok := in["kubernetes_version"].(string); ok && len(v) > 0 {
		obj.KubernetesVersion = v
	}

	if v, ok := in["manage_own_security_groups"].(bool); ok {
		obj.ManageOwnSecurityGroups = &v
	}

	if v, ok := in["node_security_groups"].([]interface{}); ok && len(v) > 0 {
		obj.NodeSecurityGroups = toArrayString(v)
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

	if v, ok := in["session_token"].(string); ok && len(v) > 0 {
		obj.SessionToken = v
	}

	if v, ok := in["subnets"].([]interface{}); ok && len(v) > 0 {
		obj.Subnets = toArrayString(v)
	}

	if v, ok := in["virtual_network"].(string); ok && len(v) > 0 {
		obj.VirtualNetwork = v
	}

	var workerPoolObjs []string

	if vs, ok := in["worker_pools"]; ok {
		if workerPoolIns, ok := vs.([]interface{}); !ok {
			return nil, errors.New("unexpected content in 'worker_pools'")
		} else if len(workerPoolIns) > 0 {
			for index, v := range workerPoolIns {
				if workerPoolIn, ok := v.(map[string]interface{}); ok {
					workerPoolObj, err := expandWorkerPool(workerPoolIn)
					if err != nil {
						return nil, err
					}

					workerPoolObjs = append(workerPoolObjs, workerPoolObj)
				} else {
					return nil, fmt.Errorf("unexpected content in worker_pool with index %d", index)
				}
			}
		}
	}

	obj.WorkerPools = workerPoolObjs
	return obj, nil
}

func expandWorkerPool(workerPoolIn map[string]interface{}) (string, error) {
	var workerPoolObj AmazonElasticContainerWorkerPool

	if v, ok := workerPoolIn["ami"].(string); ok && len(v) > 0 {
		workerPoolObj.AMI = v
	}

	if v, ok := workerPoolIn["associate_worker_node_public_ip"].(bool); ok {
		workerPoolObj.AssociateWorkerNodePublicIP = &v
	}

	if v, ok := workerPoolIn["desired_nodes"].(int); ok && v > 0 {
		workerPoolObj.DesiredNodes = int64(v)
	}

	if v, ok := workerPoolIn["ebs_encryption"].(bool); ok {
		workerPoolObj.EBSEncryption = v
	}

	if v, ok := workerPoolIn["instance_type"].(string); ok && len(v) > 0 {
		workerPoolObj.InstanceType = v
	}

	if v, ok := workerPoolIn["maximum_nodes"].(int); ok && v > 0 {
		workerPoolObj.MaximumNodes = int64(v)
	}

	if v, ok := workerPoolIn["minimum_nodes"].(int); ok && v > 0 {
		workerPoolObj.MinimumNodes = int64(v)
	}

	if v, ok := workerPoolIn["node_volume_size"].(int); ok && v > 0 {
		workerPoolObj.NodeVolumeSize = int64(v)
	}

	if v, ok := workerPoolIn["placement_group"].(string); ok && len(v) > 0 {
		workerPoolObj.PlacementGroup = v
	}

	if v, ok := workerPoolIn["user_data"].(string); ok && len(v) > 0 {
		workerPoolObj.UserData = v
	}

	if workerPoolObj.DesiredNodes == 0 {
		workerPoolObj.DesiredNodes = workerPoolObj.MinimumNodes
	}

	if v, ok := workerPoolIn["name"].(string); ok && len(v) > 0 {
		workerPoolObj.Name = v
	}

	if v, ok := workerPoolIn["add_default_label"].(bool); ok {
		workerPoolObj.AddDefaultLabel = v
	}

	if v, ok := workerPoolIn["add_default_taint"].(bool); ok {
		workerPoolObj.AddDefaultTaint = v
	}

	if v, ok := workerPoolIn["additional_labels"].(map[string]interface{}); ok && len(v) > 0 {
		workerPoolObj.AdditionalLabels = toMapString(v)
	}

	if v, ok := workerPoolIn["additional_taints"].([]interface{}); ok && len(v) > 0 {
		additionalTaintsObjs, err := expandWorkerPoolAdditionalTaints(v, workerPoolObj.Name)
		if err != nil {
			return "", err
		}

		workerPoolObj.AdditionalTaints = additionalTaintsObjs
	}

	if v, ok := workerPoolIn["subnets"].([]interface{}); ok && len(v) > 0 {
		workerPoolObj.Subnets = toArrayString(v)
	}

	bs, err := json.Marshal(workerPoolObj)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func expandWorkerPoolAdditionalTaints(additionalTaintsIn []interface{}, poolName string) ([]K8sTaint, error) {
	additionalTaintsObjs := make([]K8sTaint, 0, len(additionalTaintsIn))
	for index, additionalTaintIn := range additionalTaintsIn {
		if t, ok := additionalTaintIn.(map[string]interface{}); ok {
			taint := toMapString(t)
			additionalTaintsObj := K8sTaint{}

			if effect, ok := taint["effect"]; ok && len(effect) > 0 {
				additionalTaintsObj.Effect = effect
			}

			if operator, ok := taint["operator"]; ok && len(operator) > 0 {
				additionalTaintsObj.Operator = operator
			}

			if key, ok := taint["key"]; ok && len(key) > 0 {
				additionalTaintsObj.Key = key
			}

			if value, ok := taint["value"]; ok && len(value) > 0 {
				additionalTaintsObj.Value = value
			}

			additionalTaintsObjs = append(additionalTaintsObjs, additionalTaintsObj)
		} else {
			return nil, fmt.Errorf("taint in index %d for worker pool %s contains unexpected content", index, poolName)
		}
	}
	return additionalTaintsObjs, nil
}
