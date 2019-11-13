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
	OpenstackCredentialConfig     *openstackCredentialConfig     `json:"openstackcredentialConfig,omitempty" yaml:"openstackcredentialConfig,omitempty"`
	VmwarevsphereCredentialConfig *vmwarevsphereCredentialConfig `json:"vmwarevspherecredentialConfig,omitempty" yaml:"vmwarevspherecredentialConfig,omitempty"`
}

//Schemas

func cloudCredentialFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"amazonec2_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"azure_credential_config", "digitalocean_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialAmazonec2Fields(),
			},
		},
		"azure_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "digitalocean_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialAzureFields(),
			},
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"digitalocean_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialDigitaloceanFields(),
			},
		},
		"driver": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"openstack_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "digitalocean_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialOpenstackFields(),
			},
		},
		"vsphere_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "digitalocean_credential_config", "openstack_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialVsphereFields(),
			},
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
