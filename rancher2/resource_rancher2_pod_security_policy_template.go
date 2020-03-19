package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
	"log"
	"time"
)

func resourceRancher2PodSecurityPolicyTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2PodSecurityPolicyTemplateCreate,
		Read:   resourceRancher2PodSecurityPolicyTemplateRead,
		Update: resourceRancher2PodSecurityPolicyTemplateUpdate,
		Delete: resourceRancher2PodSecurityPolicyTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2PodSecurityPolicyTemplateImport,
		},

		Schema: podSecurityPolicyTemplateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2PodSecurityPolicyTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	podSecurityPolicyTemplate := expandPodSecurityPolicyTemplate(d)

	log.Printf("[INFO] Creating PodSecurityPolicyTemplate %s", podSecurityPolicyTemplate.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newPodSecurityPolicyTemplate, err := client.PodSecurityPolicyTemplate.Create(podSecurityPolicyTemplate)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{newPodSecurityPolicyTemplate.ID},
		Refresh:    podSecurityPolicyTemplateStateRefreshFunc(client, newPodSecurityPolicyTemplate.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for podSecurityPolicyTemplate (%s) to be created: %s", newPodSecurityPolicyTemplate.ID, waitErr)
	}

	err = flattenPodSecurityPolicyTemplate(d, newPodSecurityPolicyTemplate)
	if err != nil {
		return err
	}

	return resourceRancher2PodSecurityPolicyTemplateRead(d, meta)
}

func resourceRancher2PodSecurityPolicyTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing PodSecurityPolicyTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	pspt, err := client.PodSecurityPolicyTemplate.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] PodSecurityPolicyTemplate with ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenPodSecurityPolicyTemplate(d, pspt)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2PodSecurityPolicyTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating PodSecurityPolicyTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	pspt, err := client.PodSecurityPolicyTemplate.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"annotations":  toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":       toMapString(d.Get("labels").(map[string]interface{})),
		"allow_privilege_escalation": d.Get("allow_privilege_escalation").(bool),
		"allowed_capabilities": d.Get("allowed_capabilities").([]interface{}),
		"allowed_proc_mount_types": d.Get("allowed_proc_mount_types").([]interface{}),
		"allowed_unsafe_sysctls": d.Get("allowed_unsafe_sysctls").([]interface{}),
		"default_add_capabilities": d.Get("default_add_capabilities").([]interface{}),
		"default_allow_privilege_escalation": d.Get("default_allow_privilege_escalation").(bool),
		"forbidden_sysctls": d.Get("forbidden_sysctls").([]interface{}),
		"host_ipc": d.Get("host_ipc").(bool),
		"host_network": d.Get("host_network").(bool),
		"host_pid": d.Get("host_pid").(bool),
		"privileged": d.Get("privileged").(bool),
		"read_only_root_filesystem": d.Get("read_only_root_filesystem").(bool),
		"required_drop_capabilities": d.Get("required_drop_capabilities").([]interface{}),
		"volumes": d.Get("volumes").([]interface{}),
	}

	if pspt.AllowedCSIDrivers != nil && d.HasChange("allowed_csi_driver") {
		update["allowed_csi_driver"] = expandPodSecurityPolicyAllowedCSIDrivers(d.Get("allowed_csi_driver").([]interface{}))
	}

	if pspt.AllowedFlexVolumes != nil && d.HasChange("allowed_flex_volume") {
		update["allowed_flex_volume"] = expandPodSecurityPolicyAllowedFlexVolumes(d.Get("allowed_flex_volume").([]interface{}))
	}

	if pspt.FSGroup != nil && d.HasChange("fs_group") {
		update["fs_group"] = expandPodSecurityPolicyFSGroup(d.Get("fs_group").([]interface{}))
	}

	if pspt.HostPorts != nil && d.HasChange("host_port") {
		update["host_port"] = expandPodSecurityPolicyHostPortRanges(d.Get("host_port").([]interface{}))
	}

	if pspt.RunAsUser != nil && d.HasChange("run_as_user") {
		update["run_as_user"] = expandPodSecurityPolicyRunAsUser(d.Get("run_as_user").([]interface{}))
	}

	if pspt.RunAsGroup != nil && d.HasChange("run_as_group") {
		update["run_as_group"] = expandPodSecurityPolicyRunAsGroup(d.Get("run_as_group").([]interface{}))
	}

	if pspt.RuntimeClass != nil && d.HasChange("runtime_class") {
		update["runtime_class"] = expandPodSecurityPolicyRuntimeClassStrategy(d.Get("runtime_class").([]interface{}))
	}

	if pspt.SELinux != nil && d.HasChange("se_linux") {
		update["se_linux"] = expandPodSecurityPolicySELinuxStrategy(d.Get("runtime_class").([]interface{}))
	}

	if pspt.SupplementalGroups != nil && d.HasChange("supplemental_group") {
		update["supplemental_group"] = expandPodSecurityPolicySupplementalGroups(d.Get("supplemental_group").([]interface{}))
	}

	newPspt, err := client.PodSecurityPolicyTemplate.Update(pspt, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{pspt.ID},
		Target:     []string{newPspt.ID},
		Refresh:    podSecurityPolicyTemplateStateRefreshFunc(client, newPspt.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for PodSecurityPolicyTemplate (%s) to be updated: %s", newPspt.ID, waitErr)
	}

	return resourceRancher2PodSecurityPolicyTemplateRead(d, meta)
}

func resourceRancher2PodSecurityPolicyTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting PodSecurityPolicyTemplate with ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	pspt, err := client.PodSecurityPolicyTemplate.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] PodSecurityPolicyTemplate with ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.PodSecurityPolicyTemplate.Delete(pspt)
	if err != nil {
		return fmt.Errorf("Error removing PodSecurityPolicyTemplate: %s", err)
	}

	log.Printf("[DEBUG] Waiting for PodSecurityPolicyTemplate (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    podSecurityPolicyTemplateStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for PodSecurityPolicyTemplate (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// podSecurityPolicyTemplateStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher PodSecurityPolicyTemplate.
func podSecurityPolicyTemplateStateRefreshFunc(client *managementClient.Client, psptId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.PodSecurityPolicyTemplate.ByID(psptId)
		if err != nil {
			if IsNotFound(err) {
				return obj, "not found", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, obj.Created, nil
	}
}

