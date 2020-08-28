package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenProjectAlertRule(d *schema.ResourceData, in *managementClient.ProjectAlertRule) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("project_id", in.ProjectID)
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

	if in.PodRule != nil {
		err := d.Set("pod_rule", flattenPodRule(in.PodRule))
		if err != nil {
			return err
		}
	}

	d.Set("repeat_interval_seconds", int(in.RepeatIntervalSeconds))

	if len(in.Severity) > 0 {
		d.Set("severity", in.Severity)
	}

	if in.WorkloadRule != nil {
		err := d.Set("workload_rule", flattenWorkloadRule(in.WorkloadRule))
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

func expandProjectAlertRule(in *schema.ResourceData) *managementClient.ProjectAlertRule {
	obj := &managementClient.ProjectAlertRule{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.ProjectID = in.Get("project_id").(string)
	obj.GroupID = in.Get("group_id").(string)
	obj.GroupIntervalSeconds = int64(in.Get("group_interval_seconds").(int))
	obj.GroupWaitSeconds = int64(in.Get("group_wait_seconds").(int))

	if v, ok := in.Get("inherited").(bool); ok {
		obj.Inherited = &v
	}

	if v, ok := in.Get("metric_rule").([]interface{}); ok && len(v) > 0 {
		obj.MetricRule = expandMetricRule(v)
	}

	if v, ok := in.Get("pod_rule").([]interface{}); ok && len(v) > 0 {
		obj.PodRule = expandPodRule(v)
	}

	obj.RepeatIntervalSeconds = int64(in.Get("repeat_interval_seconds").(int))

	if v, ok := in.Get("severity").(string); ok {
		obj.Severity = v
	}

	if v, ok := in.Get("workload_rule").([]interface{}); ok && len(v) > 0 {
		obj.WorkloadRule = expandWorkloadRule(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
