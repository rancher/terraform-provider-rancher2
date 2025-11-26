package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2MachineConfigV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	if id == "" {
		return []*schema.ResourceData{}, fmt.Errorf("Importing rancher2_machine_config_v2 requires a non-empty ID")
	}

	cfg, ok := meta.(*Config)
	if !ok || cfg == nil {
		return []*schema.ResourceData{}, fmt.Errorf("Importing rancher2_machine_config_v2: provider config is nil")
	}

	// Available kinds - would be good to have it not hardcoded here....
	kinds := []string{
		machineConfigV2Amazonec2Kind,
		machineConfigV2AzureKind,
		machineConfigV2DigitaloceanKind,
		machineConfigV2HarvesterKind,
		machineConfigV2LinodeKind,
		machineConfigV2OpenstackKind,
		machineConfigV2VmwarevsphereKind,
		machineConfigV2GoogleGCEKind,
	}

	var lastErr error

	for _, kind := range kinds {
		obj, err := getMachineConfigV2ByID(cfg, id, kind)
		if err != nil {
			lastErr = err
			continue
		}

		if err := flattenMachineConfigV2(d, obj); err != nil {
			return []*schema.ResourceData{}, err
		}

		return []*schema.ResourceData{d}, nil
	}

	if lastErr != nil {
		return []*schema.ResourceData{}, lastErr
	}

	return []*schema.ResourceData{}, fmt.Errorf("[ERROR] Unable to find Machine Config V2 %s for any supported driver", id)
}
