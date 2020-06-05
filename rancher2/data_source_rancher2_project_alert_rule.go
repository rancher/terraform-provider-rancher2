package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ProjectAlertRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ProjectAlertRuleRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert rule project ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert rule name",
			},
			"metric_rule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: metricRuleFields(),
				},
				Description: "Alert metric rule",
			},
			"pod_rule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: podRuleFields(),
				},
				Description: "Alert pod rule",
			},
			"workload_rule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: workloadRuleFields(),
				},
				Description: "Alert workload rule",
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alert rule group ID",
			},
			"group_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert rule interval seconds",
			},
			"group_wait_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert rule wait seconds",
			},
			"inherited": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Alert rule inherited",
			},
			"repeat_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert rule repeat interval seconds",
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alert rule severity",
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func dataSourceRancher2ProjectAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"projectId": projectID,
		"name":      name,
	}
	listOpts := NewListOpts(filters)

	alertRules, err := client.ProjectAlertRule.List(listOpts)
	if err != nil {
		return err
	}

	count := len(alertRules.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] project alert rule with name \"%s\" on project ID \"%s\" not found", name, projectID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d project alert rule with name \"%s\" on project ID \"%s\"", count, name, projectID)
	}

	return flattenProjectAlertRule(d, &alertRules.Data[0])
}
