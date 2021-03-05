package rancher2

import (
	"fmt"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

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

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	obj["bucket_name"] = in.BucketName
	obj["endpoint"] = in.Endpoint

	if len(in.Folder) > 0 {
		obj["folder"] = in.Folder
	}

	obj["region"] = in.Region

	if len(in.CustomCA) > 0 {
		obj["custom_ca"] = Base64Encode(in.CustomCA)
	}

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

	obj["safe_timestamp"] = in.SafeTimestamp

	if in.Timeout > 0 {
		obj["timeout"] = int(in.Timeout)
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

	if in.GID >= 0 {
		obj["gid"] = int(in.GID)
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

	if in.UID >= 0 {
		obj["uid"] = int(in.UID)
	}

	obj["snapshot"] = *in.Snapshot

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigServicesEtcdBackupConfigS3(p []interface{}) (*managementClient.S3BackupConfig, error) {
	obj := &managementClient.S3BackupConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["bucket_name"].(string); ok && len(v) > 0 {
		obj.BucketName = v
	}

	if v, ok := in["custom_ca"].(string); ok && len(v) > 0 {
		customCA, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding etcd backup S3 Config: custom_ca is not base64 encoded: %s", v)
		}
		obj.CustomCA = customCA
	}

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["folder"].(string); ok && len(v) > 0 {
		obj.Folder = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	return obj, nil
}

func expandClusterRKEConfigServicesEtcdBackupConfig(p []interface{}) (*managementClient.BackupConfig, error) {
	obj := &managementClient.BackupConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
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
		s3BackupConfig, err := expandClusterRKEConfigServicesEtcdBackupConfigS3(v)
		if err != nil {
			return nil, err
		}
		obj.S3BackupConfig = s3BackupConfig
	}

	if v, ok := in["safe_timestamp"].(bool); ok {
		obj.SafeTimestamp = v
	}

	if v, ok := in["timeout"].(int); ok && v > 0 {
		obj.Timeout = int64(v)
	}

	return obj, nil
}

func expandClusterRKEConfigServicesEtcd(p []interface{}) (*managementClient.ETCDService, error) {
	obj := &managementClient.ETCDService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["backup_config"].([]interface{}); ok && len(v) > 0 {
		backupConfig, err := expandClusterRKEConfigServicesEtcdBackupConfig(v)
		if err != nil {
			return nil, err
		}
		obj.BackupConfig = backupConfig
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

	if v, ok := in["gid"].(int); ok && v >= 0 {
		obj.GID = int64(v)
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

	if v, ok := in["uid"].(int); ok && v >= 0 {
		obj.UID = int64(v)
	}

	return obj, nil
}
