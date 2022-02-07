package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2Kind       = "MachineConfig"
	machineConfigV2APIVersion = "rke-machine-config.cattle.io/v1"
	machineConfigV2APIType    = "rke-machine-config.cattle.io"
)

//Types

type machineConfigV2 struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata,omitempty"`
	Amazonec2Config     *MachineConfigV2Amazonec2     `json:"amazonec2Config,omitempty" yaml:"amazonec2Config,omitempty"`
	AzureConfig         *MachineConfigV2Azure         `json:"azureConfig,omitempty" yaml:"azureConfig,omitempty"`
	DigitaloceanConfig  *MachineConfigV2Digitalocean  `json:"digitaloceanConfig,omitempty" yaml:"digitaloceanConfig,omitempty"`
	HarvesterConfig     *MachineConfigV2Harvester     `json:"harvesterConfig,omitempty" yaml:"harvesterConfig,omitempty"`
	LinodeConfig        *MachineConfigV2Linode        `json:"linodeConfig,omitempty" yaml:"linodeConfig,omitempty"`
	OpenstackConfig     *MachineConfigV2Openstack     `json:"openstackConfig,omitempty" yaml:"openstackConfig,omitempty"`
	VmwarevsphereConfig *MachineConfigV2Vmwarevsphere `json:"vmwarevsphereConfig,omitempty" yaml:"vmwarevsphereConfig,omitempty"`
}

type MachineConfigV2 struct {
	norman.Resource
	machineConfigV2
}

// Flatteners

func flattenMachineConfigV2(d *schema.ResourceData, in *MachineConfigV2) error {
	if in == nil {
		return nil
	}
	kind := in.TypeMeta.Kind
	d.Set("kind", kind)
	switch kind {
	case machineConfigV2Amazonec2Kind:
		err := d.Set("amazonec2_config", flattenMachineConfigV2Amazonec2(in.Amazonec2Config))
		if err != nil {
			return err
		}
	case machineConfigV2AzureKind:
		err := d.Set("azure_config", flattenMachineConfigV2Azure(in.AzureConfig))
		if err != nil {
			return err
		}
	case machineConfigV2DigitaloceanKind:
		err := d.Set("digitalocean_config", flattenMachineConfigV2Digitalocean(in.DigitaloceanConfig))
		if err != nil {
			return err
		}
	case machineConfigV2HarvesterKind:
		err := d.Set("harvester_config", flattenMachineConfigV2Harvester(in.HarvesterConfig))
		if err != nil {
			return err
		}
	case machineConfigV2LinodeKind:
		err := d.Set("linode_config", flattenMachineConfigV2Linode(in.LinodeConfig))
		if err != nil {
			return err
		}
	case machineConfigV2OpenstackKind:
		err := d.Set("openstack_config", flattenMachineConfigV2Openstack(in.OpenstackConfig))
		if err != nil {
			return err
		}
	case machineConfigV2VmwarevsphereKind:
		err := d.Set("vsphere_config", flattenMachineConfigV2Vmwarevsphere(in.VmwarevsphereConfig))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("[ERROR] Unsupported driver on node template: %s", kind)
	}

	if len(in.ID) > 0 {
		d.SetId(in.ID)
	}
	d.Set("name", in.ObjectMeta.Name)
	d.Set("fleet_namespace", in.ObjectMeta.Namespace)
	err := d.Set("annotations", toMapInterface(in.ObjectMeta.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.ObjectMeta.Labels))
	if err != nil {
		return err
	}
	d.Set("resource_version", in.ObjectMeta.ResourceVersion)

	return nil
}

// Expanders

func expandMachineConfigV2(in *schema.ResourceData) *MachineConfigV2 {
	if in == nil {
		return nil
	}

	obj := &MachineConfigV2{}
	if len(in.Id()) > 0 {
		obj.ID = in.Id()
	}
	obj.ObjectMeta.GenerateName = "nc-" + in.Get("generate_name").(string) + "-"
	obj.ObjectMeta.Namespace = in.Get("fleet_namespace").(string)
	obj.ObjectMeta.Name = in.Get("name").(string)
	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Annotations = toMapString(v)
	}
	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Labels = toMapString(v)
	}
	if v, ok := in.Get("resource_version").(string); ok {
		obj.ObjectMeta.ResourceVersion = v
	}
	if v, ok := in.Get("amazonec2_config").([]interface{}); ok && len(v) > 0 {
		obj.Amazonec2Config = expandMachineConfigV2Amazonec2(v, obj)
	}
	if v, ok := in.Get("azure_config").([]interface{}); ok && len(v) > 0 {
		obj.AzureConfig = expandMachineConfigV2Azure(v, obj)
	}
	if v, ok := in.Get("digitalocean_config").([]interface{}); ok && len(v) > 0 {
		obj.DigitaloceanConfig = expandMachineConfigV2Digitalocean(v, obj)
	}
	if v, ok := in.Get("harvester_config").([]interface{}); ok && len(v) > 0 {
		obj.HarvesterConfig = expandMachineConfigV2Harvester(v, obj)
	}
	if v, ok := in.Get("linode_config").([]interface{}); ok && len(v) > 0 {
		obj.LinodeConfig = expandMachineConfigV2Linode(v, obj)
	}
	if v, ok := in.Get("openstack_config").([]interface{}); ok && len(v) > 0 {
		obj.OpenstackConfig = expandMachineConfigV2Openstack(v, obj)
	}
	if v, ok := in.Get("vsphere_config").([]interface{}); ok && len(v) > 0 {
		obj.VmwarevsphereConfig = expandMachineConfigV2Vmwarevsphere(v, obj)
	}

	return obj
}
