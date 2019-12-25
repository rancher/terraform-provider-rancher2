package rancher2

import (
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
    AzureFile             FSType = "azureFile"
	Flocker               FSType = "flocker"
	FlexVolume            FSType = "flexVolume"
	HostPath              FSType = "hostPath"
	EmptyDir              FSType = "emptyDir"
	GCEPersistentDisk     FSType = "gcePersistentDisk"
	AWSElasticBlockStore  FSType = "awsElasticBlockStore"
	GitRepo               FSType = "gitRepo"
	Secret                FSType = "secret"
	NFS                   FSType = "nfs"
	ISCSI                 FSType = "iscsi"
	Glusterfs             FSType = "glusterfs"
	PersistentVolumeClaim FSType = "persistentVolumeClaim"
	RBD                   FSType = "rbd"
	Cinder                FSType = "cinder"
	CephFS                FSType = "cephFS"
	DownwardAPI           FSType = "downwardAPI"
	FC                    FSType = "fc"
	ConfigMap             FSType = "configMap"
	VsphereVolume         FSType = "vsphereVolume"
	Quobyte               FSType = "quobyte"
	AzureDisk             FSType = "azureDisk"
	PhotonPersistentDisk  FSType = "photonPersistentDisk"
	StorageOS             FSType = "storageos"
	Projected             FSType = "projected"
	PortworxVolume        FSType = "portworxVolume"
	ScaleIO               FSType = "scaleIO"
	CSI                   FSType = "csi"
	All                   FSType = "*"
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
		Type: schema.TypeString
		ValidateFunc: validation.StringInSlice(fsTypes, true),
	}

	return s
}