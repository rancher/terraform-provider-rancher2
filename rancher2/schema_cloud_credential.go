package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Types

type CloudCredential struct {
	managementClient.CloudCredential
	Amazonec2CredentialConfig     *amazonec2CredentialConfig     `json:"amazonec2credentialConfig,omitempty" yaml:"amazonec2credentialConfig,omitempty"`
	AzureCredentialConfig         *azureCredentialConfig         `json:"azurecredentialConfig,omitempty" yaml:"azurecredentialConfig,omitempty"`
	DigitaloceanCredentialConfig  *digitaloceanCredentialConfig  `json:"digitaloceancredentialConfig,omitempty" yaml:"digitaloceancredentialConfig,omitempty"`
	LinodeCredentialConfig        *linodeCredentialConfig        `json:"linodecredentialConfig,omitempty" yaml:"linodecredentialConfig,omitempty"`
	OpenstackCredentialConfig     *openstackCredentialConfig     `json:"openstackcredentialConfig,omitempty" yaml:"openstackcredentialConfig,omitempty"`
	VmwarevsphereCredentialConfig *vmwarevsphereCredentialConfig `json:"vmwarevspherecredentialConfig,omitempty" yaml:"vmwarevspherecredentialConfig,omitempty"`
}

//Schemas

func cloudCredentialFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"amazonec2_credential_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"azure_credential_config", "digitalocean_credential_config", "linode_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialAmazonec2Fields(),
			},
		},
		"azure_credential_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "digitalocean_credential_config", "linode_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialAzureFields(),
			},
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"digitalocean_credential_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "linode_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialDigitaloceanFields(),
			},
		},
		"driver": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"linode_credential_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "digitalocean_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialLinodeFields(),
			},
		},
		"openstack_credential_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "digitalocean_credential_config", "linode_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialOpenstackFields(),
			},
		},
		"vsphere_credential_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "digitalocean_credential_config", "linode_credential_config", "openstack_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialVsphereFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
