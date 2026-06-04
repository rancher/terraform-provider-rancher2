package rancher2

import (
	"log"

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
			"default_pod_security_admission_configuration_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster V2 default pod security admission configuration template name",
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
					Schema: clusterRegistrationTokenFields(),
				},
			},
			"kube_config": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"generate_kube_config": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Generate a kubeconfig for the cluster. Set to false to avoid creating a new API token on each plan/apply. Default will change to false in a future version.",
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

	log.Printf("[INFO] Refreshing Cluster V2 %s", d.Id())

	cluster, err := getClusterV2ByID(meta.(*Config), d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
			log.Printf("[INFO] Cluster V2 %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("cluster_v1_id", cluster.Status.ClusterName)
	generateKubeConfig := d.Get("generate_kube_config").(bool)
	if generateKubeConfig {
		log.Printf("[WARN] Generating kubeconfig for cluster %s creates a new API token. Set generate_kube_config = false if you don't need kube_config. The default will change to false in a future version.", d.Id())
	}
	err = setClusterV2LegacyData(d, meta.(*Config), generateKubeConfig)
	if err != nil {
		return err
	}
	return flattenClusterV2(d, cluster)
}
