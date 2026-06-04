package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func dataSourceRancher2Cluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterRead,

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
			"generate_kube_config": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Generate a kubeconfig for the cluster. Set to false to avoid creating a new API token on each plan/apply. Default will change to false in a future version.",
			},
			"ca_cert": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"rke_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRKEConfigFieldsData(),
				},
			},
			"rke2_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRKE2ConfigFields(),
				},
			},
			"k3s_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterK3SConfigFields(),
				},
			},
			"eks_config_v2": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterEKSConfigV2Fields(),
				},
			},
			"aks_config_v2": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterAKSConfigV2Fields(),
				},
			},
			"gke_config_v2": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterGKEConfigV2Fields(),
				},
			},
			"oke_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
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
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterAuthEndpoint(),
				},
			},
			"cluster_registration_token": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterRegistrationTokenFields(),
				},
			},
			"cluster_template_answers": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
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
			"default_pod_security_admission_configuration_template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Default pod security admission configuration template name",
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
			"imported_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: clusterImportedConfigFields(),
				},
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

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		cluster := &Cluster{}
		err = client.APIBaseClient.ByID(managementClient.ClusterType, d.Id(), cluster)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		clusterRegistrationToken, err := findClusterRegistrationToken(client, cluster.ID)
		if err != nil && !IsForbidden(err) {
			return resource.NonRetryableError(err)
		}

		defaultProjectID, systemProjectID, err := meta.(*Config).GetClusterSpecialProjectsID(cluster.ID)
		if err != nil && !IsForbidden(err) {
			return resource.NonRetryableError(err)
		}

		var kubeConfig *managementClient.GenerateKubeConfigOutput
		generateKubeConfig := d.Get("generate_kube_config").(bool)
		if generateKubeConfig {
			log.Printf("[WARN] Generating kubeconfig for cluster %s creates a new API token. Set generate_kube_config = false if you don't need kube_config. The default will change to false in a future version.", cluster.ID)
			kubeConfig, err = getClusterKubeconfig(meta.(*Config), cluster.ID, d.Get("kube_config").(string))
			if err != nil && !IsForbidden(err) {
				return resource.NonRetryableError(err)
			}
		}
		if kubeConfig == nil {
			kubeConfig = &managementClient.GenerateKubeConfigOutput{}
		}

		if err = flattenCluster(
			d,
			cluster,
			clusterRegistrationToken,
			kubeConfig,
			defaultProjectID,
			systemProjectID); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}
