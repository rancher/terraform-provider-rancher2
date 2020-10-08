package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterAlertRule(d *schema.ResourceData, in *managementClient.ClusterAlertRule) error {
	if in == nil {
		return nil
	}

	if len(in.ID) > 0 {
		d.SetId(in.ID)
	}

	d.Set("name", in.Name)
	d.Set("cluster_id", in.ClusterID)

	if in.EventRule != nil {
		err := d.Set("event_rule", flattenEventRule(in.EventRule))
		if err != nil {
			return err
		}
	}

	d.Set("group_id", in.GroupID)
	d.Set("group_interval_seconds", int(in.GroupIntervalSeconds))
	d.Set("group_wait_seconds", int(in.GroupWaitSeconds))

	if in.Inherited != nil {
		d.Set("inherited", *in.Inherited)
	}

	if in.MetricRule != nil {
		err := d.Set("metric_rule", flattenMetricRule(in.MetricRule))
		if err != nil {
			return err
		}
	}

	if in.NodeRule != nil {
		err := d.Set("node_rule", flattenNodeRule(in.NodeRule))
		if err != nil {
			return err
		}
	}

	d.Set("repeat_interval_seconds", int(in.RepeatIntervalSeconds))

	if len(in.Severity) > 0 {
		d.Set("severity", in.Severity)
	}

	if in.SystemServiceRule != nil {
		err := d.Set("system_service_rule", flattenSystemServiceRule(in.SystemServiceRule))
		if err != nil {
			return err
		}
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

func expandClusterAlertRule(in *schema.ResourceData) *managementClient.ClusterAlertRule {
	obj := &managementClient.ClusterAlertRule{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.ClusterID = in.Get("cluster_id").(string)

	if v, ok := in.Get("event_rule").([]interface{}); ok && len(v) > 0 {
		obj.EventRule = expandEventRule(v)
	}

	obj.GroupID = in.Get("group_id").(string)
	obj.GroupIntervalSeconds = int64(in.Get("group_interval_seconds").(int))
	obj.GroupWaitSeconds = int64(in.Get("group_wait_seconds").(int))

	if v, ok := in.Get("inherited").(bool); ok {
		obj.Inherited = &v
	}

	if v, ok := in.Get("metric_rule").([]interface{}); ok && len(v) > 0 {
		obj.MetricRule = expandMetricRule(v)
	}

	if v, ok := in.Get("node_rule").([]interface{}); ok && len(v) > 0 {
		obj.NodeRule = expandNodeRule(v)
	}

	obj.RepeatIntervalSeconds = int64(in.Get("repeat_interval_seconds").(int))

	if v, ok := in.Get("severity").(string); ok {
		obj.Severity = v
	}

	if v, ok := in.Get("system_service_rule").([]interface{}); ok && len(v) > 0 {
		obj.SystemServiceRule = expandSystemServiceRule(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
