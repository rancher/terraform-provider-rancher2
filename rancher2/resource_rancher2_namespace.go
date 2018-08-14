package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

func init() {
	descriptions = map[string]string{
		"name": "Name of the k8s namespace managed by rancher v2",

		"cluster_id": "Cluster ID where k8s namespace belongs",

		"project_name": "Project name where k8s namespace belongs",

		"description": "Description of the k8s namespace managed by rancher v2",

		"resource_quota_template_id": "Resource quota template id to apply on k8s namespace",

		"annotations": "Annotations of the k8s namespace managed by rancher v2",

		"labels": "Labels of the k8s namespace managed by rancher v2",
	}
}

func resourceCattleNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceCattleNamespaceCreate,
		Read:   resourceCattleNamespaceRead,
		Update: resourceCattleNamespaceUpdate,
		Delete: resourceCattleNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCattleNamespaceImport,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: descriptions["cluster_id"],
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: descriptions["name"],
			},
			"project_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_name"],
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["description"],
			},
			"resource_quota_template_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["resource_quota_template_id"],
			},
			"annotations": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["annotations"],
			},
			"labels": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["labels"],
			},
		},
	}
}

func resourceCattleNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	clusterID := d.Get("cluster_id").(string)
	projectName := d.Get("project_name").(string)
	projectID, err := meta.(*Config).GetProjectIDByName(projectName, clusterID)
	if err != nil {
		return err
	}

	//annotations := d.Get("annotations").(map[string]string)
	//labels := d.Get("labels").(map[string]string)

	log.Printf("[INFO] Creating Namespace %s", name)

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return err
	}

	ns := &clusterClient.Namespace{
		Name:                    name,
		ProjectID:               projectID,
		Description:             d.Get("description").(string),
		ResourceQuotaTemplateID: d.Get("resource_quota_template_id").(string),
		//	Annotations: annotations,
		//	Labels:      labels,
	}

	newNs, err := client.Namespace.Create(ns)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"activating"},
		Target:     []string{"active"},
		Refresh:    NamespaceStateRefreshFunc(client, newNs.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for namespace (%s) to be created: %s", newNs.ID, waitErr)
	}

	d.SetId(newNs.ID)
	d.Set("project_id", projectID)

	return resourceCattleNamespaceRead(d, meta)
}

func resourceCattleNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	log.Printf("[INFO] Refreshing Namespace ID %s", d.Id())

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return err
	}

	ns, err := client.Namespace.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Namespace ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	projectName, err := meta.(*Config).GetProjectNameByID(ns.ProjectID)
	if err != nil {
		return err
	}

	d.Set("name", ns.Name)
	d.Set("project_id", ns.ProjectID)
	d.Set("project_name", projectName)
	d.Set("description", ns.Description)
	d.Set("resource_quota_template_id", ns.ResourceQuotaTemplateID)
	//d.Set("annotations", ns.Annotations)
	//d.Set("labels", ns.Labels)

	return nil
}

func resourceCattleNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	projectName := d.Get("project_name").(string)
	projectID, err := meta.(*Config).GetProjectIDByName(projectName, clusterID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Namespace ID %s", d.Id())

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return err
	}

	ns, err := client.Namespace.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]string{
		"projectId":               projectID,
		"description":             d.Get("description").(string),
		"resourceQuotaTemplateId": d.Get("resource_quota_template_id").(string),
		//"annotations": d.Get("annotations").(map[string]string),
		//"labels":      d.Get("labels").(map[string]string),
	}

	newNs, err := client.Namespace.Update(ns, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    NamespaceStateRefreshFunc(client, newNs.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for namespace (%s) to be updated: %s", newNs.ID, waitErr)
	}

	d.Set("project_id", projectID)

	return resourceCattleNamespaceRead(d, meta)
}

func resourceCattleNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Namespace ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ClusterClient(d.Get("cluster_id").(string))
	if err != nil {
		return err
	}

	ns, err := client.Namespace.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Namespace ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Namespace.Delete(ns)
	if err != nil {
		return fmt.Errorf("Error removing Namespace: %s", err)
	}

	log.Printf("[DEBUG] Waiting for namespace (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    NamespaceStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for namespace (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceCattleNamespaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID, resourceID := splitID(d.Id())

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	ns, err := client.Namespace.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	projectName, err := meta.(*Config).GetProjectNameByID(ns.ProjectID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(resourceID)
	d.Set("name", ns.Name)
	d.Set("cluster_id", clusterID)
	d.Set("project_id", ns.ProjectID)
	d.Set("project_name", projectName)
	d.Set("description", ns.Description)
	d.Set("resource_quota_template_id", ns.ResourceQuotaTemplateID)
	//d.Set("annotations", ns.Annotations)
	//d.Set("labels", ns.Labels)

	return []*schema.ResourceData{d}, nil
}

// NamespaceStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Namespace.
func NamespaceStateRefreshFunc(client *clusterClient.Client, nsID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ns, err := client.Namespace.ByID(nsID)
		if err != nil {
			if IsNotFound(err) {
				return ns, "removed", nil
			}
			return nil, "", err
		}

		return ns, ns.State, nil
	}
}
