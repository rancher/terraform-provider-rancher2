package rancher2

import (
	norman "github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigV2NutanixKind       = "NutanixConfig"
	machineConfigV2NutanixAPIVersion = "rke-machine-config.cattle.io/v1"
	machineConfigV2NutanixAPIType    = "rke-machine-config.cattle.io.nutanixconfig"
)

type machineConfigV2Nutanix struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Endpoint          string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	Username          string   `json:"username,omitempty" yaml:"username,omitempty"`
	Password          string   `json:"password,omitempty" yaml:"password,omitempty"`
	Port              string   `json:"port,omitempty" yaml:"port,omitempty"`
	Insecure          bool     `json:"insecure,omitempty" yaml:"insecure,omitempty"`
	Cluster           string   `json:"cluster,omitempty" yaml:"cluster,omitempty"`
	VMMem             string   `json:"vmMem,omitempty" yaml:"vmMem,omitempty"`
	VMCPUs            string   `json:"vmCpus,omitempty" yaml:"vmCpus,omitempty"`
	VMCores           string   `json:"vmCores,omitempty" yaml:"vmCores,omitempty"`
	VMCPUPassthrough  bool     `json:"vmCpuPassthrough,omitempty" yaml:"vmCpuPassthrough,omitempty"`
	VMNetwork         []string `json:"vmNetwork,omitempty" yaml:"vmNetwork,omitempty"`
	VMImage           string   `json:"vmImage,omitempty" yaml:"vmImage,omitempty"`
	VMImageSize       string   `json:"vmImageSize,omitempty" yaml:"vmImageSize,omitempty"`
	VMCategories      []string `json:"vmCategories,omitempty" yaml:"vmCategories,omitempty"`
	StorageContainer  string   `json:"storageContainer,omitempty" yaml:"storageContainer,omitempty"`
	DiskSize          string   `json:"diskSize,omitempty" yaml:"diskSize,omitempty"`
	CloudInit         string   `json:"cloudInit,omitempty" yaml:"cloudInit,omitempty"`
	VMSerialPort      bool     `json:"vmSerialPort,omitempty" yaml:"vmSerialPort,omitempty"`
	Project           string   `json:"project,omitempty" yaml:"project,omitempty"`
	BootType          string   `json:"bootType,omitempty" yaml:"bootType,omitempty"`
	Timeout           string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	VMGPU             []string `json:"vmGpu,omitempty" yaml:"vmGpu,omitempty"`
}

type MachineConfigV2Nutanix struct {
	norman.Resource
	machineConfigV2Nutanix
}

func flattenMachineConfigV2Nutanix(in *MachineConfigV2Nutanix) []any {
	if in == nil {
		return nil
	}

	obj := make(map[string]any)

	if len(in.Endpoint) > 0 {
		obj["endpoint"] = in.Endpoint
	}
	if len(in.Username) > 0 {
		obj["username"] = in.Username
	}
	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}
	if len(in.Port) > 0 {
		obj["port"] = in.Port
	}
	obj["insecure"] = in.Insecure
	if len(in.Cluster) > 0 {
		obj["cluster"] = in.Cluster
	}
	if len(in.VMMem) > 0 {
		obj["vm_mem"] = in.VMMem
	}
	if len(in.VMCPUs) > 0 {
		obj["vm_cpus"] = in.VMCPUs
	}
	if len(in.VMCores) > 0 {
		obj["vm_cores"] = in.VMCores
	}
	obj["vm_cpu_passthrough"] = in.VMCPUPassthrough
	if len(in.VMNetwork) > 0 {
		obj["vm_network"] = toArrayInterface(in.VMNetwork)
	}
	if len(in.VMImage) > 0 {
		obj["vm_image"] = in.VMImage
	}
	if len(in.VMImageSize) > 0 {
		obj["vm_image_size"] = in.VMImageSize
	}
	if len(in.VMCategories) > 0 {
		obj["vm_categories"] = toArrayInterface(in.VMCategories)
	}
	if len(in.StorageContainer) > 0 {
		obj["storage_container"] = in.StorageContainer
	}
	if len(in.DiskSize) > 0 {
		obj["disk_size"] = in.DiskSize
	}
	if len(in.CloudInit) > 0 {
		obj["cloud_init"] = in.CloudInit
	}
	obj["vm_serial_port"] = in.VMSerialPort
	if len(in.Project) > 0 {
		obj["project"] = in.Project
	}
	if len(in.BootType) > 0 {
		obj["boot_type"] = in.BootType
	}
	if len(in.Timeout) > 0 {
		obj["timeout"] = in.Timeout
	}
	if len(in.VMGPU) > 0 {
		obj["vm_gpu"] = toArrayInterface(in.VMGPU)
	}

	return []any{obj}
}

func expandMachineConfigV2Nutanix(p []any, source *MachineConfigV2) *MachineConfigV2Nutanix {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &MachineConfigV2Nutanix{}
	if len(source.ID) > 0 {
		obj.ID = source.ID
	}

	in := p[0].(map[string]any)

	obj.TypeMeta.Kind = machineConfigV2NutanixKind
	obj.TypeMeta.APIVersion = machineConfigV2NutanixAPIVersion
	source.TypeMeta = obj.TypeMeta
	obj.ObjectMeta = source.ObjectMeta

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}
	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}
	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in["port"].(string); ok && len(v) > 0 {
		obj.Port = v
	}
	if v, ok := in["insecure"].(bool); ok {
		obj.Insecure = v
	}
	if v, ok := in["cluster"].(string); ok && len(v) > 0 {
		obj.Cluster = v
	}
	if v, ok := in["vm_mem"].(string); ok && len(v) > 0 {
		obj.VMMem = v
	}
	if v, ok := in["vm_cpus"].(string); ok && len(v) > 0 {
		obj.VMCPUs = v
	}
	if v, ok := in["vm_cores"].(string); ok && len(v) > 0 {
		obj.VMCores = v
	}
	if v, ok := in["vm_cpu_passthrough"].(bool); ok {
		obj.VMCPUPassthrough = v
	}
	if v, ok := in["vm_network"].([]any); ok && len(v) > 0 {
		obj.VMNetwork = toArrayString(v)
	}
	if v, ok := in["vm_image"].(string); ok && len(v) > 0 {
		obj.VMImage = v
	}
	if v, ok := in["vm_image_size"].(string); ok && len(v) > 0 {
		obj.VMImageSize = v
	}
	if v, ok := in["vm_categories"].([]any); ok && len(v) > 0 {
		obj.VMCategories = toArrayString(v)
	}
	if v, ok := in["storage_container"].(string); ok && len(v) > 0 {
		obj.StorageContainer = v
	}
	if v, ok := in["disk_size"].(string); ok && len(v) > 0 {
		obj.DiskSize = v
	}
	if v, ok := in["cloud_init"].(string); ok && len(v) > 0 {
		obj.CloudInit = v
	}
	if v, ok := in["vm_serial_port"].(bool); ok {
		obj.VMSerialPort = v
	}
	if v, ok := in["project"].(string); ok && len(v) > 0 {
		obj.Project = v
	}
	if v, ok := in["boot_type"].(string); ok && len(v) > 0 {
		obj.BootType = v
	}
	if v, ok := in["timeout"].(string); ok && len(v) > 0 {
		obj.Timeout = v
	}
	if v, ok := in["vm_gpu"].([]any); ok && len(v) > 0 {
		obj.VMGPU = toArrayString(v)
	}

	return obj
}
