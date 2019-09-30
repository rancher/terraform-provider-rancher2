package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRancher2Cluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"driver": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_config": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"rke_config": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRKEConfigFields(),
				},
			},
			"eks_config": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterEKSConfigFields(),
				},
			},
			"aks_config": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterAKSConfigFields(),
				},
			},
			"gke_config": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterGKEConfigFields(),
				},
			},
			"default_project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_auth_endpoint": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterAuthEndpoint(),
				},
			},
			"cluster_monitoring_input": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "Cluster monitoring configuration",
				Elem: &schema.Resource{
					Schema: monitoringInputFields(),
				},
			},
			"cluster_registration_token": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRegistationTokenFields(),
				},
			},
			"default_pod_security_policy_template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default pod security policy template id",
			},
			"enable_cluster_monitoring": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable built-in cluster monitoring",
			},
			"enable_network_policy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable project network isolation",
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name": name,
	}
	listOpts := NewListOpts(filters)

	clusters, err := client.Cluster.List(listOpts)
	if err != nil {
		return err
	}

	count := len(clusters.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] cluster with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d cluster with name \"%s\"", count, name)
	}

	d.SetId(clusters.Data[0].ID)

	return resourceRancher2ClusterRead(d, meta)
}
