package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigFreeIpa(d *schema.ResourceData, in *managementClient.LdapConfig) error {
	err := flattenAuthConfigLdap(d, in)
	if err != nil {
		return err
	}

	d.SetId(in.Name)
	d.Set("name", in.Name)
	d.Set("type", managementClient.FreeIpaConfigType)

	return nil
}

// Expanders

func expandAuthConfigFreeIpa(in *schema.ResourceData) (*managementClient.LdapConfig, error) {
	obj, err := expandAuthConfigLdap(in)
	if err != nil {
		return nil, err
	}

	obj.Name = expandAuthConfigName(in, AuthConfigFreeIpaName)
	obj.Type = managementClient.FreeIpaConfigType

	return obj, nil
}
