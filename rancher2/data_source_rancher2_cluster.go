package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2Cluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2ClusterRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_env_vars": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Optional Agent Env Vars for Rancher agent",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"driver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_cert": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"rke_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRKEConfigFieldsData(),
				},
			},
			"rke2_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRKE2ConfigFields(),
				},
			},
			"k3s_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterK3SConfigFields(),
				},
			},
			"eks_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterEKSConfigFields(),
				},
			},
			"eks_config_v2": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterEKSConfigV2Fields(),
				},
			},
			"aks_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterAKSConfigFields(),
				},
			},
			"aks_config_v2": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterAKSConfigV2Fields(),
				},
			},
			"gke_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterGKEConfigFields(),
				},
			},
			"gke_config_v2": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterGKEConfigV2Fields(),
				},
			},
			"oke_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterOKEConfigFields(),
				},
			},
			"default_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_auth_endpoint": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterAuthEndpoint(),
				},
			},
			"cluster_monitoring_input": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster monitoring configuration",
				Elem: &schema.Resource{
					Schema: monitoringInputFields(),
				},
			},
			"cluster_registration_token": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRegistrationTokenFields(),
				},
			},
			"cluster_template_answers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster template answers",
				Elem: &schema.Resource{
					Schema: answerFields(),
				},
			},
			"cluster_template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster template ID",
			},
			"cluster_template_questions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster template questions",
				Elem: &schema.Resource{
					Schema: questionFields(),
				},
			},
			"cluster_template_revision_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster template revision ID",
			},
			"default_pod_security_policy_template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default pod security policy template ID",
			},
			"default_pod_security_admission_configuration_template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Default pod security admission configuration template name",
			},
			"enable_cluster_alerting": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable built-in cluster alerting",
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
			"fleet_workspace_name": {
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

			"cluster_agent_deployment_customization": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional customization for cluster agent",
				Elem: &schema.Resource{
					Schema: agentDeploymentCustomizationFields(),
				},
			},
			"fleet_agent_deployment_customization": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional customization for fleet agent",
				Elem: &schema.Resource{
					Schema: agentDeploymentCustomizationFields(),
				},
			},
			"desired_agent_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"desired_auth_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"docker_root_dir": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"istio_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"windows_prefered_cluster": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name": name,
	}
	listOpts := NewListOpts(filters)

	clusters, err := client.Cluster.List(listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	count := len(clusters.Data)
	if count <= 0 {
		return diag.Errorf("[ERROR] cluster with name \"%s\" not found", name)
	}
	if count > 1 {
		return diag.Errorf("[ERROR] found %d cluster with name \"%s\"", count, name)
	}

	d.SetId(clusters.Data[0].ID)

	return resourceRancher2ClusterRead(ctx, d, meta)
}
