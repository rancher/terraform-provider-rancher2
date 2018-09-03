package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
	//"gopkg.in/yaml.v2"
)

var (
	clusterKinds = []string{"imported", "eks", "aks", "gke", "rke"}
)

// Schema

func clusterFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"kind": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "rke",
			ValidateFunc: validation.StringInSlice(clusterKinds, true),
		},
		"rke_config": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: rkeConfigFields(),
			},
		},
		"eks_config": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: eksConfigFields(),
			},
		},
		"aks_config": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: aksConfigFields(),
			},
		},
		"gke_config": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: gkeConfigFields(),
			},
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"cluster_registration_token": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRegistationTokenFields(),
			},
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}

// Flatteners

func flattenCluster(d *schema.ResourceData, in *managementClient.Cluster, clusterRegToken *managementClient.ClusterRegistrationToken) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening cluster: Input cluster is nil")
	}

	if clusterRegToken == nil {
		return fmt.Errorf("[ERROR] flattening cluster: Input cluster registration token is nil")
	}

	d.SetId(in.ID)
	err := d.Set("name", in.Name)
	if err != nil {
		return err
	}
	err = d.Set("description", in.Description)
	if err != nil {
		return err
	}
	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}
	regToken, err := flattenClusterRegistationToken(clusterRegToken)
	if err != nil {
		return err
	}
	err = d.Set("cluster_registration_token", regToken)
	if err != nil {
		return err
	}

	switch kind := d.Get("kind").(string); kind {
	case "rke":
		rkeConfig, err := flattenRkeConfig(in.RancherKubernetesEngineConfig)
		if err != nil {
			return err
		}
		err = d.Set("rke_config", rkeConfig)
		if err != nil {
			return err
		}
	case "eks":
		eksConfig, err := flattenEksConfig(in.AmazonElasticContainerServiceConfig)
		if err != nil {
			return err
		}
		d.Set("eks_config", eksConfig)
		if err != nil {
			return err
		}
	case "aks":
		aksConfig, err := flattenAksConfig(in.AzureKubernetesServiceConfig)
		if err != nil {
			return err
		}
		d.Set("aks_config", aksConfig)
		if err != nil {
			return err
		}
	case "gke":
		gkeConfig, err := flattenGkeConfig(in.GoogleKubernetesEngineConfig)
		if err != nil {
			return err
		}
		d.Set("gke_config", gkeConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

// Expanders

func expandCluster(in *schema.ResourceData) (*managementClient.Cluster, error) {
	obj := &managementClient.Cluster{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding cluster: Input ResourceData is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	switch kind := in.Get("kind").(string); kind {
	case "rke":
		rkeConfig, err := expandRkeConfig(in.Get("rke_config").([]interface{}))
		if err != nil {
			return nil, err
		}
		obj.RancherKubernetesEngineConfig = rkeConfig
	case "eks":
		eksConfig, err := expandEksConfig(in.Get("eks_config").([]interface{}))
		if err != nil {
			return nil, err
		}
		obj.AmazonElasticContainerServiceConfig = eksConfig
	case "aks":
		aksConfig, err := expandAksConfig(in.Get("aks_config").([]interface{}))
		if err != nil {
			return nil, err
		}
		obj.AzureKubernetesServiceConfig = aksConfig
	case "gke":
		gkeConfig, err := expandGkeConfig(in.Get("gke_config").([]interface{}))
		if err != nil {
			return nil, err
		}
		obj.GoogleKubernetesEngineConfig = gkeConfig
	}

	return obj, nil
}

func resourceRancher2Cluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterCreate,
		Read:   resourceRancher2ClusterRead,
		Update: resourceRancher2ClusterUpdate,
		Delete: resourceRancher2ClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterImport,
		},
		Schema: clusterFields(),
	}
}

func resourceRancher2ClusterCreate(d *schema.ResourceData, meta interface{}) error {
	cluster, err := expandCluster(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Cluster %s", cluster.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	expectedState := "active"
	kind := d.Get("kind").(string)

	if kind == "imported" {
		expectedState := "pending"
	}

	if kind == "rke" {
		expectedState = "provisioning"
	}

	newCluster, err := client.Cluster.Create(cluster)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{expectedState},
		Refresh:    ClusterStateRefreshFunc(client, newCluster.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster (%s) to be created: %s", newCluster.ID, waitErr)
	}

	clusterReg, err := expandClusterRegistationToken(d.Get("cluster_registration_token").([]interface{}), newCluster.ID)
	if err != nil {
		// If error undo create cluster
		client.Cluster.Delete(cluster)
		return err
	}

	newClusterRegistrationToken, err := client.ClusterRegistrationToken.Create(clusterReg)
	if err != nil {
		// If error undo create cluster
		client.Cluster.Delete(cluster)
		return err
	}

	stateConf = &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    ClusterRegistrationTokenStateRefreshFunc(client, newClusterRegistrationToken.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr = stateConf.WaitForState()
	if waitErr != nil {
		// If error undo create cluster
		client.Cluster.Delete(cluster)
		return fmt.Errorf("[ERROR] waiting for cluster registration token (%s) to be created: %s", newClusterRegistrationToken.ID, waitErr)
	}

	d.SetId(newCluster.ID)
	err = flattenCluster(d, newCluster, newClusterRegistrationToken)
	if err != nil {
		// If error undo create cluster
		client.Cluster.Delete(cluster)
		return err
	}

	return resourceRancher2ClusterRead(d, meta)
}

func resourceRancher2ClusterRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cluster, err := client.Cluster.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	clusterRegistrationToken, err := findClusterRegistrationToken(client, d.Id())
	if err != nil {
		return err
	}

	err = flattenCluster(d, cluster, clusterRegistrationToken)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2ClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Cluster ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cluster, err := client.Cluster.ByID(d.Id())
	if err != nil {
		return err
	}

	rkeConfig, err := expandRkeConfig(d.Get("rke_config").([]interface{}))
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"name":                          d.Get("name").(string),
		"description":                   d.Get("description").(string),
		"rancherKubernetesEngineConfig": rkeConfig,
		"annotations":                   toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                        toMapString(d.Get("labels").(map[string]interface{})),
	}

	newCluster, err := meta.(*Config).UpdateClusterByID(cluster, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "provisioning", "pending"},
		Target:     []string{"active", "provisioning", "pending"},
		Refresh:    ClusterStateRefreshFunc(client, newCluster.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster (%s) to be updated: %s", newCluster.ID, waitErr)
	}

	return resourceRancher2ClusterRead(d, meta)
}

func resourceRancher2ClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cluster, err := client.Cluster.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Cluster.Delete(cluster)
	if err != nil {
		return fmt.Errorf("Error removing Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cluster (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    ClusterStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancher2ClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	cluster, err := client.Cluster.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	clusterRegistrationToken, err := findClusterRegistrationToken(client, d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenCluster(d, cluster, clusterRegistrationToken)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}

// ClusterStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster.
func ClusterStateRefreshFunc(client *managementClient.Client, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clus, err := client.Cluster.ByID(clusterID)
		if err != nil {
			if IsNotFound(err) {
				return clus, "removed", nil
			}
			return nil, "", err
		}

		return clus, clus.State, nil
	}
}

// ClusterRegistrationTokenStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher ClusterRegistrationToken.
func ClusterRegistrationTokenStateRefreshFunc(client *managementClient.Client, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusReg, err := client.ClusterRegistrationToken.ByID(clusterID)
		if err != nil {
			if IsNotFound(err) {
				return clusReg, "removed", nil
			}
			return nil, "", err
		}

		return clusReg, clusReg.State, nil
	}
}

func normalizeAmazonEksConfig(in *managementClient.AmazonElasticContainerServiceConfig) map[string]interface{} {
	out := make(map[string]interface{})

	out["accessKey"] = in.AccessKey
	out["instanceType"] = in.InstanceType
	out["maximumNodes"] = in.MaximumNodes
	out["minimumNodes"] = in.MinimumNodes
	out["region"] = in.Region
	out["secretKey"] = in.SecretKey
	out["securityGroups"] = in.SecurityGroups
	out["serviceRole"] = in.ServiceRole
	out["subnets"] = in.Subnets
	out["virtualNetwork"] = in.VirtualNetwork
	return out
}

func normalizeAzureAksConfig(in *managementClient.AzureKubernetesServiceConfig) map[string]interface{} {
	out := make(map[string]interface{})

	out["adminUsername"] = in.AdminUsername
	out["agentDnsPrefix"] = in.AgentDNSPrefix
	out["agentPoolName"] = in.AgentPoolName
	out["agentVmSize"] = in.AgentVMSize
	out["clientId"] = in.ClientID
	out["clientSecret"] = in.ClientSecret
	out["count"] = in.Count
	out["kubernetesVersion"] = in.KubernetesVersion
	out["location"] = in.Location
	out["masterDnsPrefix"] = in.MasterDNSPrefix
	out["osDiskSizeGb"] = in.OsDiskSizeGB
	out["resourceGroup"] = in.ResourceGroup
	out["sshPublicKeyContents"] = in.SSHPublicKeyContents
	out["subnet"] = in.Subnet
	out["subscriptionId"] = in.SubscriptionID
	out["tags"] = in.Tag
	out["tenantId"] = in.TenantID
	out["virtualNetwork"] = in.VirtualNetwork
	out["virtualNetworkResourceGroup"] = in.VirtualNetworkResourceGroup
	return out
}

func normalizeGoogleGkeConfig(in *managementClient.GoogleKubernetesEngineConfig) map[string]interface{} {
	out := make(map[string]interface{})

	out["clusterIpv4Cidr"] = in.ClusterIpv4Cidr
	out["credential"] = in.Credential
	out["description"] = in.Description
	out["diskSizeGb"] = in.DiskSizeGb
	out["enableAlphaFeature"] = in.EnableAlphaFeature
	out["enableHttpLoadBalancing"] = in.EnableHTTPLoadBalancing
	out["enableHorizontalPodAutoscaling"] = in.EnableHorizontalPodAutoscaling
	out["enableKubernetesDashboard"] = in.EnableKubernetesDashboard
	out["enableLegacyAbac"] = in.EnableLegacyAbac
	out["enableNetworkPolicyConfig"] = in.EnableNetworkPolicyConfig
	out["enableStackdriverLogging"] = in.EnableStackdriverLogging
	out["enableStackdriverMonitoring"] = in.EnableStackdriverMonitoring
	out["imageType"] = in.ImageType
	out["labels"] = in.Labels
	out["locations"] = in.Locations
	out["machineType"] = in.MachineType
	out["maintenanceWindow"] = in.MaintenanceWindow
	out["masterVersion"] = in.MasterVersion
	out["network"] = in.Network
	out["nodeCount"] = in.NodeCount
	out["nodeVersion"] = in.NodeVersion
	out["projectId"] = in.ProjectID
	out["subNetwork"] = in.SubNetwork
	out["zone"] = in.Zone
	return out
}
