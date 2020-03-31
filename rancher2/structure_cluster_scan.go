package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterScanCisConfig(in *managementClient.CisScanConfig) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})
	obj["debug_master"] = in.DebugMaster
	obj["debug_worker"] = in.DebugWorker
	if len(in.OverrideBenchmarkVersion) > 0 {
		obj["override_benchmark_version"] = in.OverrideBenchmarkVersion
	}
	if in.OverrideSkip != nil && len(in.OverrideSkip) > 0 {
		obj["override_skip"] = toArrayInterface(in.OverrideSkip)
	}
	if len(in.Profile) > 0 {
		obj["profile"] = in.Profile
	}

	return []interface{}{obj}
}

func flattenClusterScanConfig(in *managementClient.ClusterScanConfig) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})
	obj["cis_scan_config"] = flattenClusterScanCisConfig(in.CisScanConfig)

	return []interface{}{obj}
}

func flattenClusterScan(d *schema.ResourceData, in *managementClient.ClusterScan) error {
	if in == nil {
		return nil
	}

	d.Set("cluster_id", in.ClusterID)
	d.Set("name", in.Name)
	if len(in.RunType) > 0 {
		d.Set("run_type", in.RunType)
	}
	d.Set("scan_config", flattenClusterScanConfig(in.ScanConfig))

	if len(in.ScanType) > 0 {
		d.Set("scan_type", in.ScanType)
	}

	if in.Status != nil {
		str, _ := interfaceToYAML(in.Status)
		d.Set("status", str)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandClusterScanCisConfig(p []interface{}) *managementClient.CisScanConfig {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	obj := &managementClient.CisScanConfig{}
	if v, ok := in["debug_master"].(bool); ok {
		obj.DebugMaster = v
	}

	if v, ok := in["debug_worker"].(bool); ok {
		obj.DebugWorker = v
	}

	if v, ok := in["override_benchmark_version"].(string); ok && len(v) > 0 {
		obj.OverrideBenchmarkVersion = v
	}

	if v, ok := in["override_skip"].([]interface{}); ok && len(v) > 0 {
		obj.OverrideSkip = toArrayString(v)
	}

	if v, ok := in["profile"].(string); ok && len(v) > 0 {
		obj.Profile = v
	}

	return obj
}

func expandClusterScanConfig(p []interface{}) *managementClient.ClusterScanConfig {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})
	obj := &managementClient.ClusterScanConfig{}
	if v, ok := in["cis_scan_config"].([]interface{}); ok && len(v) > 0 {
		obj.CisScanConfig = expandClusterScanCisConfig(v)
	}

	return obj
}

func expandClusterScan(in *schema.ResourceData) *managementClient.ClusterScan {
	if in == nil {
		return nil
	}
	obj := &managementClient.ClusterScan{}
	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("cluster_id").(string); ok && len(v) > 0 {
		obj.ClusterID = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("run_type").(string); ok && len(v) > 0 {
		obj.RunType = v
	}

	if v, ok := in.Get("scan_config").([]interface{}); ok && len(v) > 0 {
		obj.ScanConfig = expandClusterScanConfig(v)
	}

	if v, ok := in.Get("scan_type").(string); ok && len(v) > 0 {
		obj.ScanType = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
