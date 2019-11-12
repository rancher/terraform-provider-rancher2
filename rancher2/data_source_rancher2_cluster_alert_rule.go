package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterAlertRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterAlertRuleRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert rule cluster ID",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert rule name",
			},
			"event_rule": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: eventRuleFields(),
				},
				Description: "Alert event rule",
			},
			"metric_rule": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: metricRuleFields(),
				},
				Description: "Alert metric rule",
			},
			"node_rule": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: nodeRuleFields(),
				},
				Description: "Alert node rule",
			},
			"system_service_rule": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: systemServiceRuleFields(),
				},
				Description: "Alert system service rule",
			},
			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alert rule group ID",
			},
			"group_interval_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert rule interval seconds",
			},
			"group_wait_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert rule wait seconds",
			},
			"inherited": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Alert rule inherited",
			},
			"repeat_interval_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert rule repeat interval seconds",
			},
			"severity": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alert rule severity",
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func dataSourceRancher2ClusterAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
		"name":      name,
	}
	listOpts := NewListOpts(filters)

	alertRules, err := client.ClusterAlertRule.List(listOpts)
	if err != nil {
		return err
	}

	count := len(alertRules.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] cluster alert rule with name \"%s\" on cluster ID \"%s\" not found", name, clusterID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d cluster alert rule with name \"%s\" on cluster ID \"%s\"", count, name, clusterID)
	}

	return flattenClusterAlertRule(d, &alertRules.Data[0])
}
