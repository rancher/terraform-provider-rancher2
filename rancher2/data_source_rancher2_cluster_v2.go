package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterV2Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster V2 name",
			},
			"fleet_namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "fleet-default",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster V2 kubernetes version",
			},
			"rke_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "Cluster V2 rke config",
				Elem: &schema.Resource{
					Schema: clusterV2RKEConfigFields(),
				},
			},
			"agent_env_vars": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster V2 default agent env vars",
				Elem: &schema.Resource{
					Schema: envVarFields(),
				},
			},
			"cloud_credential_secret_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster V2 cloud credential secret name",
			},
			"default_pod_security_policy_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster V2 default pod security policy template name",
			},
			"default_cluster_role_for_project_members": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster V2 default cluster role for project members",
			},
			"enable_network_policy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable k8s network policy",
			},
			"cluster_registration_token": {
				Type:      schema.TypeList,
				MaxItems:  1,
				Computed:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: clusterRegistationTokenFields(),
				},
			},
			"kube_config": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"cluster_v1_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterV2Read(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	namespace := d.Get("fleet_namespace").(string)
	d.SetId(namespace + clusterV2ClusterIDsep + name)

	return resourceRancher2ClusterV2Read(d, meta)
}
