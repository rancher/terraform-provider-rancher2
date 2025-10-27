package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rancher/pkg/api/steve/catalog/types"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Flatteners

func flattenAppV2(d *schema.ResourceData, in *AppV2) error {
	if in == nil {
		return nil
	}

	d.Set("name", in.ObjectMeta.Name)
	d.Set("namespace", in.ObjectMeta.Namespace)
	err := d.Set("annotations", toMapInterface(in.ObjectMeta.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.ObjectMeta.Labels))
	if err != nil {
		return err
	}

	if in.Spec.Chart != nil && in.Spec.Chart.Metadata != nil {
		d.Set("chart_name", in.Spec.Chart.Metadata.Name)
		d.Set("chart_version", in.Spec.Chart.Metadata.Version)
	}

	if len(in.Spec.Values) > 0 {
		valuesStr, err := interfaceToGhodssyaml(in.Spec.Values)
		if err != nil {
			return fmt.Errorf("failed to marshal chart values yaml: %v", err)
		}
		d.Set("deployment_values", valuesStr)
		if global, ok := in.Spec.Values["global"].(map[string]interface{}); ok && len(global) > 0 {
			if cattle, ok := global["cattle"].(map[string]interface{}); ok && len(cattle) > 0 {
				if clusterID, ok := cattle["clusterId"].(string); ok && len(clusterID) > 0 {
					d.Set("cluster_id", clusterID)
				}
				if clusterName, ok := cattle["clusterName"].(string); ok && len(clusterName) > 0 {
					d.Set("cluster_name", clusterName)
				}
				if systemDefaultRegistry, ok := cattle["systemDefaultRegistry"].(string); ok && len(systemDefaultRegistry) > 0 {
					d.Set("system_default_registry", systemDefaultRegistry)
				}
			}
		}
	}
	if len(in.ID) > 0 {
		d.SetId(d.Get("cluster_id").(string) + appV2ClusterIDsep + in.ID)
	}

	return nil
}

// Expanders

func expandChartInstallV2(in *schema.ResourceData, chartInfo *types.ChartInfo) (string, []types.ChartInstall, error) {
	if in == nil || chartInfo == nil || chartInfo.Chart == nil {
		return "", nil, nil
	}
	out := make([]types.ChartInstall, 0)
	name := in.Get("name").(string)
	namespace := in.Get("namespace").(string)
	globalInfo := generateGlobalInfoMap(in)
	valuesData := v3.MapStringInterface{}
	if v, ok := in.Get("values").(string); ok {
		values, err := unmarshalValuesContent(v)
		if err != nil {
			return "", nil, err
		}
		mergeGlobalMaps(values, globalInfo)
		valuesData = v3.MapStringInterface(values)
	}
	if chartAnnotations, ok := chartInfo.Chart["annotations"].(map[string]interface{}); ok && len(chartAnnotations) > 0 {
		if autoInstall, ok := chartAnnotations["catalog.cattle.io/auto-install"].(string); ok && len(autoInstall) > 0 {
			chartAuto := splitBySep(autoInstall, "=")
			if len(chartAuto) != 2 {
				return "", nil, fmt.Errorf("wrong format on chart annotation catalog.cattle.io/auto-install: %s", autoInstall)
			}
			chartName := chartAuto[0]
			chartVersion := chartAuto[1]
			if len(chartAuto[1]) == 0 || chartAuto[1] == "match" {
				chartVersion = chartInfo.Chart["version"].(string)
			}
			obj := types.ChartInstall{
				ChartName:   chartName,
				Version:     chartVersion,
				ReleaseName: chartName,
				Values: v3.MapStringInterface{
					"global": valuesData["global"],
				},
			}
			out = append(out, obj)
		}
		// Forcing release name and namespace if rancher certified
		if chartCertified, ok := chartAnnotations["catalog.cattle.io/certified"].(string); ok && chartCertified == "rancher" {
			if chartName, ok := chartAnnotations["catalog.cattle.io/release-name"].(string); ok && len(chartName) > 0 {
				name = chartName
				in.Set("name", name)
			}
			if chartNamespace, ok := chartAnnotations["catalog.cattle.io/namespace"].(string); ok && len(chartNamespace) > 0 {
				namespace = chartNamespace
				in.Set("namespace", namespace)
			}
		}
	}
	obj := types.ChartInstall{
		ChartName:   chartInfo.Chart["name"].(string),
		Version:     chartInfo.Chart["version"].(string),
		ReleaseName: name,
		Values:      valuesData,
	}
	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}
	out = append(out, obj)

	return namespace, out, nil
}

func mergeGlobalMaps(values map[string]interface{}, globalInfo map[string]interface{}) {
	globalInfoCattle := globalInfo["cattle"].(map[string]interface{})

	if global, ok := values["global"].(map[string]interface{}); ok && len(global) > 0 {
		global["systemDefaultRegistry"] = globalInfo["systemDefaultRegistry"]
		if globalCattle, ok := global["cattle"].(map[string]interface{}); ok && len(global) > 0 {
			globalCattle["clusterId"] = globalInfoCattle["clusterId"]
			globalCattle["clusterName"] = globalInfoCattle["clusterName"]
			globalCattle["systemDefaultRegistry"] = globalInfoCattle["systemDefaultRegistry"]

		} else {
			global["cattle"] = globalInfo["cattle"]
		}
	} else {
		values["global"] = globalInfo
	}
}

func generateGlobalInfoMap(in *schema.ResourceData) map[string]interface{} {
	globalInfoCattle := map[string]interface{}{
		"clusterId":             in.Get("cluster_id").(string),
		"clusterName":           in.Get("cluster_name").(string),
		"systemDefaultRegistry": in.Get("system_default_registry").(string),
	}

	globalInfo := map[string]interface{}{
		"systemDefaultRegistry": in.Get("system_default_registry").(string),
		"cattle":                globalInfoCattle,
	}

	return globalInfo
}

func expandChartInstallActionV2(in *schema.ResourceData, chartInfo *types.ChartInfo) (*types.ChartInstallAction, error) {
	if in == nil || chartInfo == nil {
		return nil, nil
	}

	namespace, chartIntalls, err := expandChartInstallV2(in, chartInfo)
	if err != nil {
		return nil, err
	}
	wait := in.Get("wait").(bool)
	if len(chartIntalls) > 1 {
		// Forcing wait = true if chart has dependencies
		wait = true
	}

	timeOut := &metaV1.Duration{}
	timeOut.Duration = in.Timeout(schema.TimeoutCreate)
	obj := &types.ChartInstallAction{
		Timeout:                  timeOut,
		Wait:                     wait,
		DisableHooks:             in.Get("disable_hooks").(bool),
		DisableOpenAPIValidation: in.Get("disable_open_api_validation").(bool),
		Namespace:                namespace,
		Charts:                   chartIntalls,
		ProjectID:                in.Get("project_id").(string),
	}

	return obj, nil
}

func expandChartUpgradeV2(in *schema.ResourceData, chartInfo *types.ChartInfo) (string, []types.ChartUpgrade, error) {
	if in == nil {
		return "", nil, nil
	}

	out := make([]types.ChartUpgrade, 0)
	chartName := in.Get("chart_name").(string)
	chartVersion := in.Get("chart_version").(string)
	name := in.Get("name").(string)
	namespace := in.Get("namespace").(string)
	globalInfo := generateGlobalInfoMap(in)
	valuesData := v3.MapStringInterface{}
	if v, ok := in.Get("values").(string); ok {
		values, err := unmarshalValuesContent(v)
		if err != nil {
			return "", nil, err
		}
		mergeGlobalMaps(values, globalInfo)
		valuesData = v3.MapStringInterface(values)
	}
	if chartAnnotations, ok := chartInfo.Chart["annotations"].(map[string]interface{}); ok && len(chartAnnotations) > 0 {
		if autoInstall, ok := chartAnnotations["catalog.cattle.io/auto-install"].(string); ok && len(autoInstall) > 0 {
			chartAuto := splitBySep(autoInstall, "=")
			if len(chartAuto) != 2 {
				return "", nil, fmt.Errorf("wrong format on chart annotation catalog.cattle.io/auto-install: %s", autoInstall)
			}
			chartName := chartAuto[0]
			chartVersion := chartAuto[1]
			if len(chartAuto[1]) == 0 || chartAuto[1] == "match" {
				chartVersion = chartInfo.Chart["version"].(string)
			}
			obj := types.ChartUpgrade{
				ChartName:   chartName,
				Version:     chartVersion,
				ReleaseName: chartName,
				Values: v3.MapStringInterface{
					"global": valuesData["global"],
				},
				Force: in.Get("force_upgrade").(bool),
			}
			out = append(out, obj)
		}
		// Forcing release name and namespace if rancher certified
		if chartCertified, ok := chartAnnotations["catalog.cattle.io/certified"].(string); ok && chartCertified == "rancher" {
			if chartName, ok := chartAnnotations["catalog.cattle.io/release-name"].(string); ok && len(chartName) > 0 {
				name = chartName
				in.Set("name", name)
			}
			if chartNamespace, ok := chartAnnotations["catalog.cattle.io/namespace"].(string); ok && len(chartNamespace) > 0 {
				namespace = chartNamespace
				in.Set("namespace", namespace)
			}
		}
	}
	obj := types.ChartUpgrade{
		ChartName:   chartName,
		Version:     chartVersion,
		ReleaseName: name,
		Values:      valuesData,
		Force:       in.Get("force_upgrade").(bool),
	}
	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}
	out = append(out, obj)

	return namespace, out, nil
}

func unmarshalValuesContent(v string) (map[string]interface{}, error) {
	values, err := ghodssyamlToMapInterface(v)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to unmarshal chart install values YAML: %#v", err)
	}
	if values == nil {
		values = map[string]interface{}{}
	}

	return values, nil
}

func expandChartUpgradeActionV2(in *schema.ResourceData, chartInfo *types.ChartInfo) (*types.ChartUpgradeAction, error) {
	if in == nil || chartInfo == nil {
		return nil, nil
	}

	namespace, chartUpgrades, err := expandChartUpgradeV2(in, chartInfo)
	if err != nil {
		return nil, err
	}
	wait := in.Get("wait").(bool)
	if len(chartUpgrades) > 1 {
		// Forcing wait = true if chart has dependencies
		wait = true
	}

	timeOut := &metaV1.Duration{}
	timeOut.Duration = in.Timeout(schema.TimeoutUpdate)
	obj := &types.ChartUpgradeAction{
		Timeout:                  timeOut,
		Wait:                     wait,
		DisableHooks:             in.Get("disable_hooks").(bool),
		DisableOpenAPIValidation: in.Get("disable_open_api_validation").(bool),
		Force:                    in.Get("force_upgrade").(bool),
		Namespace:                namespace,
		CleanupOnFail:            in.Get("cleanup_on_fail").(bool),
		Charts:                   chartUpgrades,
	}

	return obj, nil
}
