package rancher2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
)

func resourceRancher2MachineConfigV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2MachineConfigV2Create,
		Read:   resourceRancher2MachineConfigV2Read,
		Update: resourceRancher2MachineConfigV2Update,
		Delete: resourceRancher2MachineConfigV2Delete,
		Schema: machineConfigV2Fields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2MachineConfigV2Create(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	obj := expandMachineConfigV2(d)

	log.Printf("[INFO] Creating Machine Config V2 %s kind %s", name, obj.TypeMeta.Kind)

	newObj, err := createMachineConfigV2(meta.(*Config), obj)
	if err != nil {
		return err
	}

	d.SetId(newObj.ID)
	d.Set("kind", newObj.TypeMeta.Kind)

	err = waitForMachineConfigV2(d, meta.(*Config), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] waiting for machine config (%s) to be active: %s", newObj.ID, err)
	}

	return flattenMachineConfigV2(d, newObj)
}

func waitForMachineConfigV2(d *schema.ResourceData, config *Config, interval time.Duration) error {
	log.Printf("[INFO] Waiting for state Machine Config V2 %s", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()
	for {
		kind := d.Get("kind").(string)
		_, err := getMachineConfigV2ByID(config, d.Id(), kind)
		if err == nil {
			return nil
		}
		log.Printf("[INFO] Retrying on error Refreshing Machine Config V2 %s: %v", d.Id(), err)
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return fmt.Errorf("Machine Config V2 %s not found: %s", d.Id(), err)
		}
		if IsNotAccessibleByID(err) {
			// Restarting clients to update RBAC
			config.RestartClients()
		}

		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return fmt.Errorf("Timeout waiting for machine config V2 ID %s", d.Id())
		}
	}
}

func resourceRancher2MachineConfigV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Machine Config V2 %s", d.Id())

	kind := d.Get("kind").(string)
	obj, err := getMachineConfigV2ByID(meta.(*Config), d.Id(), kind)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
			log.Printf("[INFO] Machine Config V2 %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	return flattenMachineConfigV2(d, obj)
}

func resourceRancher2MachineConfigV2Update(d *schema.ResourceData, meta interface{}) error {
	obj := expandMachineConfigV2(d)
	log.Printf("[INFO] Updating Machine Config V2 %s", d.Id())

	newObj, err := updateMachineConfigV2(meta.(*Config), obj)
	if err != nil {
		return err
	}
	d.SetId(newObj.ID)
	flattenMachineConfigV2(d, newObj)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    machineConfigV2StateRefreshFunc(meta, newObj.ID, newObj.TypeMeta.Kind),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for machine config (%s) to be active: %s", newObj.ID, waitErr)
	}
	return resourceRancher2MachineConfigV2Read(d, meta)
}

func resourceRancher2MachineConfigV2Delete(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	kind := d.Get("kind").(string)
	log.Printf("[INFO] Deleting Machine Config V2 %s", name)

	obj, err := getMachineConfigV2ByID(meta.(*Config), d.Id(), kind)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	err = deleteMachineConfigV2(meta.(*Config), obj)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    machineConfigV2StateRefreshFunc(meta, obj.ID, obj.TypeMeta.Kind),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for machine config v2 (%s) to be removed: %s", obj.ID, waitErr)
	}
	d.SetId("")
	return nil
}

// machineConfigV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Machine Config v2.
func machineConfigV2StateRefreshFunc(meta interface{}, objID, kind string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := getMachineConfigV2ByID(meta.(*Config), objID, kind)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			// This is required to allow standard user to use this resource
			if !IsNotAccessibleByID(err) {
				return obj, "active", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}

// Rancher2 Machine Config V2 API CRUD functions
func createMachineConfigV2(c *Config, obj *MachineConfigV2) (*MachineConfigV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Creating Machine Config V2: Provider config is nil")
	}
	if obj == nil {
		return nil, fmt.Errorf("Creating Machine Config V2: Machine Config V2 is nil")
	}
	var err error
	out := &MachineConfigV2{}
	kind := obj.TypeMeta.Kind
	switch kind {
	case machineConfigV2Amazonec2Kind:
		resp := &MachineConfigV2Amazonec2{}
		err = c.createObjectV2(rancher2DefaultLocalClusterID, machineConfigV2Amazonec2APIType, obj.Amazonec2Config, resp)
		out.Amazonec2Config = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2AzureKind:
		resp := &MachineConfigV2Azure{}
		err = c.createObjectV2(rancher2DefaultLocalClusterID, machineConfigV2AzureAPIType, obj.AzureConfig, resp)
		out.AzureConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2DigitaloceanKind:
		resp := &MachineConfigV2Digitalocean{}
		err = c.createObjectV2(rancher2DefaultLocalClusterID, machineConfigV2DigitaloceanAPIType, obj.DigitaloceanConfig, resp)
		out.DigitaloceanConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2HarvesterKind:
		resp := &MachineConfigV2Harvester{}
		err = c.createObjectV2(rancher2DefaultLocalClusterID, machineConfigV2HarvesterAPIType, obj.HarvesterConfig, resp)
		out.HarvesterConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2LinodeKind:
		resp := &MachineConfigV2Linode{}
		err = c.createObjectV2(rancher2DefaultLocalClusterID, machineConfigV2LinodeAPIType, obj.LinodeConfig, resp)
		out.LinodeConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2OpenstackKind:
		resp := &MachineConfigV2Openstack{}
		err = c.createObjectV2(rancher2DefaultLocalClusterID, machineConfigV2OpenstackAPIType, obj.OpenstackConfig, resp)
		out.OpenstackConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2VmwarevsphereKind:
		resp := &MachineConfigV2Vmwarevsphere{}
		err = c.createObjectV2(rancher2DefaultLocalClusterID, machineConfigV2VmwarevsphereAPIType, obj.VmwarevsphereConfig, resp)
		out.VmwarevsphereConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	default:
		return nil, fmt.Errorf("[ERROR] Unsupported driver on node template: %s", kind)
	}
	if err != nil {
		return nil, fmt.Errorf("Creating Machine Config V2: %s", err)
	}
	return out, nil
}

func deleteMachineConfigV2(c *Config, obj *MachineConfigV2) error {
	if c == nil {
		return fmt.Errorf("Deleting Machine Config V2: Provider config is nil")
	}
	if obj == nil {
		return fmt.Errorf("Deleting Machine Config V2: Machine Config V2 is nil")
	}
	resource := &norman.Resource{
		ID:      obj.ID,
		Links:   obj.Links,
		Type:    obj.Type,
		Actions: obj.Actions,
	}
	return c.deleteObjectV2(rancher2DefaultLocalClusterID, resource)
}

func getMachineConfigV2ByID(c *Config, id, kind string) (*MachineConfigV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting Machine Config V2: Provider config is nil")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Getting Machine Config V2: Machine Config V2 ID is empty")
	}
	var err error
	out := &MachineConfigV2{}
	switch kind {
	case machineConfigV2Amazonec2Kind:
		resp := &MachineConfigV2Amazonec2{}
		err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, machineConfigV2Amazonec2APIType, resp)
		out.Amazonec2Config = resp
		out.ID = resp.ID
		out.Links = resp.Links
		out.Actions = resp.Actions
		out.Type = resp.Type
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2AzureKind:
		resp := &MachineConfigV2Azure{}
		err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, machineConfigV2AzureAPIType, resp)
		out.AzureConfig = resp
		out.ID = resp.ID
		out.Links = resp.Links
		out.Actions = resp.Actions
		out.Type = resp.Type
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2DigitaloceanKind:
		resp := &MachineConfigV2Digitalocean{}
		err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, machineConfigV2DigitaloceanAPIType, resp)
		out.DigitaloceanConfig = resp
		out.ID = resp.ID
		out.Links = resp.Links
		out.Actions = resp.Actions
		out.Type = resp.Type
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2HarvesterKind:
		resp := &MachineConfigV2Harvester{}
		err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, machineConfigV2HarvesterAPIType, resp)
		out.HarvesterConfig = resp
		out.ID = resp.ID
		out.Links = resp.Links
		out.Actions = resp.Actions
		out.Type = resp.Type
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2LinodeKind:
		resp := &MachineConfigV2Linode{}
		err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, machineConfigV2LinodeAPIType, resp)
		out.LinodeConfig = resp
		out.ID = resp.ID
		out.Links = resp.Links
		out.Actions = resp.Actions
		out.Type = resp.Type
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2OpenstackKind:
		resp := &MachineConfigV2Openstack{}
		err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, machineConfigV2OpenstackAPIType, resp)
		out.OpenstackConfig = resp
		out.ID = resp.ID
		out.Links = resp.Links
		out.Actions = resp.Actions
		out.Type = resp.Type
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2VmwarevsphereKind:
		resp := &MachineConfigV2Vmwarevsphere{}
		err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, machineConfigV2VmwarevsphereAPIType, resp)
		out.VmwarevsphereConfig = resp
		out.ID = resp.ID
		out.Links = resp.Links
		out.Actions = resp.Actions
		out.Type = resp.Type
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	default:
		return nil, fmt.Errorf("[ERROR] Unsupported driver on node template: %s", kind)
	}
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting Machine Config V2: %s", err)
		}
		return nil, err
	}
	return out, nil
}

func updateMachineConfigV2(c *Config, obj *MachineConfigV2) (*MachineConfigV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Updating Machine Config V2: Provider config is nil")
	}
	if obj == nil {
		return nil, fmt.Errorf("Updating Machine Config V2: Machine Config V2 is nil")
	}
	var err error
	out := &MachineConfigV2{}
	kind := obj.TypeMeta.Kind
	switch kind {
	case machineConfigV2Amazonec2Kind:
		resp := &MachineConfigV2Amazonec2{}
		err = c.updateObjectV2(rancher2DefaultLocalClusterID, obj.ID, machineConfigV2Amazonec2APIType, obj.Amazonec2Config, resp)
		out.Amazonec2Config = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2AzureKind:
		resp := &MachineConfigV2Azure{}
		err = c.updateObjectV2(rancher2DefaultLocalClusterID, obj.ID, machineConfigV2AzureAPIType, obj.AzureConfig, resp)
		out.AzureConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2DigitaloceanKind:
		resp := &MachineConfigV2Digitalocean{}
		err = c.updateObjectV2(rancher2DefaultLocalClusterID, obj.ID, machineConfigV2DigitaloceanAPIType, obj.DigitaloceanConfig, resp)
		out.DigitaloceanConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2HarvesterKind:
		resp := &MachineConfigV2Harvester{}
		err = c.updateObjectV2(rancher2DefaultLocalClusterID, obj.ID, machineConfigV2HarvesterAPIType, obj.HarvesterConfig, resp)
		out.HarvesterConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2LinodeKind:
		resp := &MachineConfigV2Linode{}
		err = c.updateObjectV2(rancher2DefaultLocalClusterID, obj.ID, machineConfigV2LinodeAPIType, obj.LinodeConfig, resp)
		out.LinodeConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2OpenstackKind:
		resp := &MachineConfigV2Openstack{}
		err = c.updateObjectV2(rancher2DefaultLocalClusterID, obj.ID, machineConfigV2OpenstackAPIType, obj.OpenstackConfig, resp)
		out.OpenstackConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	case machineConfigV2VmwarevsphereKind:
		resp := &MachineConfigV2Vmwarevsphere{}
		err = c.updateObjectV2(rancher2DefaultLocalClusterID, obj.ID, machineConfigV2VmwarevsphereAPIType, obj.VmwarevsphereConfig, resp)
		out.VmwarevsphereConfig = resp
		out.ID = resp.ID
		out.TypeMeta = resp.TypeMeta
		out.ObjectMeta = resp.ObjectMeta
	default:
		return nil, fmt.Errorf("[ERROR] Unsupported driver on node template: %s", kind)
	}
	if err != nil {
		return nil, fmt.Errorf("Updating Machine Config V2: %s", err)
	}
	return out, err
}
