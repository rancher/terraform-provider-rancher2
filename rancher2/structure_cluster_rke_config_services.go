package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterRKEConfigServicesKubeproxy(in *managementClient.KubeproxyService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigServicesKubelet(in *managementClient.KubeletService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ClusterDNSServer) > 0 {
		obj["cluster_dns_server"] = in.ClusterDNSServer
	}

	if len(in.ClusterDomain) > 0 {
		obj["cluster_domain"] = in.ClusterDomain
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	obj["fail_swap_on"] = in.FailSwapOn

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.InfraContainerImage) > 0 {
		obj["infra_container_image"] = in.InfraContainerImage
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigServicesKubeController(in *managementClient.KubeControllerService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ClusterCIDR) > 0 {
		obj["cluster_cidr"] = in.ClusterCIDR
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.ServiceClusterIPRange) > 0 {
		obj["service_cluster_ip_range"] = in.ServiceClusterIPRange
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigServicesKubeAPI(in *managementClient.KubeAPIService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	obj["pod_security_policy"] = in.PodSecurityPolicy

	if len(in.ServiceClusterIPRange) > 0 {
		obj["service_cluster_ip_range"] = in.ServiceClusterIPRange
	}

	if len(in.ServiceNodePortRange) > 0 {
		obj["service_node_port_range"] = in.ServiceNodePortRange
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigServicesEtcdBackupConfigS3(in *managementClient.S3BackupConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	obj["access_key"] = in.AccessKey
	obj["bucket_name"] = in.BucketName
	obj["endpoint"] = in.Endpoint
	obj["region"] = in.Region

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	return []interface{}{obj}
}

func flattenClusterRKEConfigServicesEtcdBackupConfig(in *managementClient.BackupConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	obj["enabled"] = *in.Enabled

	if in.IntervalHours > 0 {
		obj["interval_hours"] = int(in.IntervalHours)
	}

	if in.Retention > 0 {
		obj["retention"] = int(in.Retention)
	}

	if in.S3BackupConfig != nil {
		v, ok := obj["s3_backup_config"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		obj["s3_backup_config"] = flattenClusterRKEConfigServicesEtcdBackupConfigS3(in.S3BackupConfig, v)
	}

	return []interface{}{obj}
}

func flattenClusterRKEConfigServicesEtcd(in *managementClient.ETCDService, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.BackupConfig != nil {
		v, ok := obj["backup_config"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		obj["backup_config"] = flattenClusterRKEConfigServicesEtcdBackupConfig(in.BackupConfig, v)
	}

	if len(in.CACert) > 0 {
		obj["ca_cert"] = in.CACert
	}

	if len(in.Cert) > 0 {
		obj["cert"] = in.Cert
	}

	if len(in.Creation) > 0 {
		obj["creation"] = in.Creation
	}

	if len(in.ExternalURLs) > 0 {
		obj["external_urls"] = toArrayInterface(in.ExternalURLs)
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.Key) > 0 {
		obj["key"] = in.Key
	}

	if len(in.Path) > 0 {
		obj["path"] = in.Path
	}

	if len(in.Retention) > 0 {
		obj["retention"] = in.Retention
	}

	obj["snapshot"] = *in.Snapshot

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigServices(in *managementClient.RKEConfigServices, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.Etcd != nil {
		v, ok := obj["etcd"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		etcd, err := flattenClusterRKEConfigServicesEtcd(in.Etcd, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["etcd"] = etcd
	}

	if in.KubeAPI != nil {
		kubeAPI, err := flattenClusterRKEConfigServicesKubeAPI(in.KubeAPI)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kube_api"] = kubeAPI
	}

	if in.KubeController != nil {
		kubeController, err := flattenClusterRKEConfigServicesKubeController(in.KubeController)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kube_controller"] = kubeController
	}

	if in.Kubelet != nil {
		kubelet, err := flattenClusterRKEConfigServicesKubelet(in.Kubelet)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kubelet"] = kubelet
	}

	if in.Kubeproxy != nil {
		kubeproxy, err := flattenClusterRKEConfigServicesKubeproxy(in.Kubeproxy)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["kubeproxy"] = kubeproxy
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigServicesKubeproxy(p []interface{}) (*managementClient.KubeproxyService, error) {
	obj := &managementClient.KubeproxyService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	return obj, nil
}

func expandClusterRKEConfigServicesKubelet(p []interface{}) (*managementClient.KubeletService, error) {
	obj := &managementClient.KubeletService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_dns_server"].(string); ok && len(v) > 0 {
		obj.ClusterDNSServer = v
	}

	if v, ok := in["cluster_domain"].(string); ok && len(v) > 0 {
		obj.ClusterDomain = v
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["fail_swap_on"].(bool); ok {
		obj.FailSwapOn = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["infra_container_image"].(string); ok && len(v) > 0 {
		obj.InfraContainerImage = v
	}

	return obj, nil
}

func expandClusterRKEConfigServicesKubeController(p []interface{}) (*managementClient.KubeControllerService, error) {
	obj := &managementClient.KubeControllerService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_cidr"].(string); ok && len(v) > 0 {
		obj.ClusterCIDR = v
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["service_cluster_ip_range"].(string); ok && len(v) > 0 {
		obj.ServiceClusterIPRange = v
	}

	return obj, nil
}

func expandClusterRKEConfigServicesKubeAPI(p []interface{}) (*managementClient.KubeAPIService, error) {
	obj := &managementClient.KubeAPIService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["pod_security_policy"].(bool); ok {
		obj.PodSecurityPolicy = v
	}

	if v, ok := in["service_cluster_ip_range"].(string); ok && len(v) > 0 {
		obj.ServiceClusterIPRange = v
	}

	if v, ok := in["service_node_port_range"].(string); ok && len(v) > 0 {
		obj.ServiceNodePortRange = v
	}

	return obj, nil
}

func expandClusterRKEConfigServicesEtcdBackupConfigS3(p []interface{}) *managementClient.S3BackupConfig {
	obj := &managementClient.S3BackupConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["bucket_name"].(string); ok && len(v) > 0 {
		obj.BucketName = v
	}

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	return obj
}

func expandClusterRKEConfigServicesEtcdBackupConfig(p []interface{}) *managementClient.BackupConfig {
	obj := &managementClient.BackupConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = &v
	}

	if v, ok := in["interval_hours"].(int); ok && v > 0 {
		obj.IntervalHours = int64(v)
	}

	if v, ok := in["retention"].(int); ok && v > 0 {
		obj.Retention = int64(v)
	}

	if v, ok := in["s3_backup_config"].([]interface{}); ok && len(v) > 0 {
		obj.S3BackupConfig = expandClusterRKEConfigServicesEtcdBackupConfigS3(v)
	}

	return obj
}

func expandClusterRKEConfigServicesEtcd(p []interface{}) (*managementClient.ETCDService, error) {
	obj := &managementClient.ETCDService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["backup_config"].([]interface{}); ok && len(v) > 0 {
		obj.BackupConfig = expandClusterRKEConfigServicesEtcdBackupConfig(v)
	}

	if v, ok := in["ca_cert"].(string); ok && len(v) > 0 {
		obj.CACert = v
	}

	if v, ok := in["cert"].(string); ok && len(v) > 0 {
		obj.Cert = v
	}

	if v, ok := in["creation"].(string); ok && len(v) > 0 {
		obj.Creation = v
	}

	if v, ok := in["external_urls"].([]interface{}); ok && len(v) > 0 {
		obj.ExternalURLs = toArrayString(v)
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["key"].(string); ok && len(v) > 0 {
		obj.Key = v
	}

	if v, ok := in["path"].(string); ok && len(v) > 0 {
		obj.Path = v
	}

	if v, ok := in["retention"].(string); ok && len(v) > 0 {
		obj.Retention = v
	}

	if v, ok := in["snapshot"].(bool); ok {
		obj.Snapshot = &v
	}

	return obj, nil
}

func expandClusterRKEConfigServices(p []interface{}) (*managementClient.RKEConfigServices, error) {
	obj := &managementClient.RKEConfigServices{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["etcd"].([]interface{}); ok && len(v) > 0 {
		etcd, err := expandClusterRKEConfigServicesEtcd(v)
		if err != nil {
			return obj, err
		}
		obj.Etcd = etcd
	}

	if v, ok := in["kube_api"].([]interface{}); ok && len(v) > 0 {
		kubeAPI, err := expandClusterRKEConfigServicesKubeAPI(v)
		if err != nil {
			return obj, err
		}
		obj.KubeAPI = kubeAPI
	}

	if v, ok := in["kube_controller"].([]interface{}); ok && len(v) > 0 {
		kubeController, err := expandClusterRKEConfigServicesKubeController(v)
		if err != nil {
			return obj, err
		}
		obj.KubeController = kubeController
	}

	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0 {
		kubelet, err := expandClusterRKEConfigServicesKubelet(v)
		if err != nil {
			return obj, err
		}
		obj.Kubelet = kubelet
	}

	if v, ok := in["kubeproxy"].([]interface{}); ok && len(v) > 0 {
		kubeproxy, err := expandClusterRKEConfigServicesKubeproxy(v)
		if err != nil {
			return obj, err
		}
		obj.Kubeproxy = kubeproxy
	}

	return obj, nil
}
