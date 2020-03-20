package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	AzureFile             string = "azureFile"
	Flocker               string = "flocker"
	FlexVolume            string = "flexVolume"
	HostPath              string = "hostPath"
	EmptyDir              string = "emptyDir"
	GCEPersistentDisk     string = "gcePersistentDisk"
	AWSElasticBlockStore  string = "awsElasticBlockStore"
	GitRepo               string = "gitRepo"
	Secret                string = "secret"
	NFS                   string = "nfs"
	ISCSI                 string = "iscsi"
	Glusterfs             string = "glusterfs"
	PersistentVolumeClaim string = "persistentVolumeClaim"
	RBD                   string = "rbd"
	Cinder                string = "cinder"
	CephFS                string = "cephFS"
	DownwardAPI           string = "downwardAPI"
	FC                    string = "fc"
	ConfigMap             string = "configMap"
	VsphereVolume         string = "vsphereVolume"
	Quobyte               string = "quobyte"
	AzureDisk             string = "azureDisk"
	PhotonPersistentDisk  string = "photonPersistentDisk"
	StorageOS             string = "storageos"
	Projected             string = "projected"
	PortworxVolume        string = "portworxVolume"
	ScaleIO               string = "scaleIO"
	CSI                   string = "csi"
	All                   string = "*"
)

var (
	fsTypes = []string{
		AzureFile,
		Flocker,
		FlexVolume,
		HostPath,
		EmptyDir,
		GCEPersistentDisk,
		AWSElasticBlockStore,
		GitRepo,
		Secret,
		NFS,
		ISCSI,
		Glusterfs,
		PersistentVolumeClaim,
		RBD,
		Cinder,
		CephFS,
		DownwardAPI,
		FC,
		ConfigMap,
		VsphereVolume,
		Quobyte,
		AzureDisk,
		PhotonPersistentDisk,
		StorageOS,
		Projected,
		PortworxVolume,
		ScaleIO,
		CSI,
		All,
	}
)

func podSecurityPolicyVolumesFields() *schema.Schema {
	s := &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validation.StringInSlice(fsTypes, true),
	}

	return s
}
