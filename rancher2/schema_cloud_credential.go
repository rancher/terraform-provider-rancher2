package rancher2

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
	genericCredentialConfig       *genericCredentialConfig
}

const credentialConfigKeySuffix = "credentialConfig"

func isTypedCredentialConfigKey(key string) bool {
	typedDrivers := []string{amazonec2ConfigDriver, azureConfigDriver, digitaloceanConfigDriver, openstackConfigDriver, vmwarevsphereConfigDriver}
	for _, dn := range typedDrivers {
		if key == dn+credentialConfigKeySuffix {
			return true
		}
	}
	return false
}

func (n *CloudCredential) driverName() string {
	switch {
	case n.Amazonec2CredentialConfig != nil:
		return amazonec2ConfigDriver
	case n.AzureCredentialConfig != nil:
		return azureConfigDriver
	case n.DigitaloceanCredentialConfig != nil:
		return digitaloceanConfigDriver
	case n.OpenstackCredentialConfig != nil:
		return openstackConfigDriver
	case n.VmwarevsphereCredentialConfig != nil:
		return vmwarevsphereConfigDriver
	case n.genericCredentialConfig != nil:
		return n.genericCredentialConfig.driverName
	}
	return ""
}

func (n *CloudCredential) UnmarshalJSON(data []byte) error {
	type Alias CloudCredential
	var dest Alias
	if err := json.Unmarshal(data, &dest); err != nil {
		return err
	}

	var rawValues map[string]interface{}
	if err := json.Unmarshal(data, &rawValues); err != nil {
		return err
	}

	for key, value := range rawValues {
		if strings.HasSuffix(key, credentialConfigKeySuffix) && !isTypedCredentialConfigKey(key) {
			driverName := strings.Replace(key, credentialConfigKeySuffix, "", -1)
			if cv, ok := value.(map[string]interface{}); ok {
				dest.genericCredentialConfig = &genericCredentialConfig{
					driverName: driverName,
					config:     cv,
				}
			}
		}
	}

	*n = CloudCredential(dest)
	return nil
}

func (n *CloudCredential) MarshalJSON() ([]byte, error) {
	type Alias CloudCredential
	data, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
	if err != nil {
		return nil, err
	}

	var results map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	if n.genericCredentialConfig != nil {
		configName := n.driverName() + credentialConfigKeySuffix
		results[configName] = n.genericCredentialConfig.config
	}

	return json.Marshal(results)
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
			ConflictsWith: []string{"azure_credential_config", "digitalocean_credential_config", "generic_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialAmazonec2Fields(),
			},
		},
		"azure_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "digitalocean_credential_config", "generic_credential_config", "openstack_credential_config", "vsphere_credential_config"},
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
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "generic_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialDigitaloceanFields(),
			},
		},
		"driver": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"generic_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "digitalocean_credential_config", "openstack_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialGenericFields(),
			},
		},
		"openstack_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "digitalocean_credential_config", "generic_credential_config", "vsphere_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialOpenstackFields(),
			},
		},
		"vsphere_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config", "digitalocean_credential_config", "generic_credential_config", "openstack_credential_config"},
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
