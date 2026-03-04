package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigKeyCloakOIDC(d *schema.ResourceData, in *managementClient.KeyCloakOIDCConfig) error {
	d.SetId(AuthConfigKeyCloakOIDCName)
	d.Set("name", AuthConfigKeyCloakOIDCName)
	d.Set("type", managementClient.KeyCloakOIDCConfigType)

	if err := flattenOIDCConfig(d, in); err != nil {
		return fmt.Errorf("flattening AuthConfig for KeyCloakOIDC: %s", err)
	}

	return nil
}

// Expanders

func expandAuthConfigKeyCloakOIDC(in *schema.ResourceData) (*managementClient.KeyCloakOIDCConfig, error) {
	obj := &managementClient.KeyCloakOIDCConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", AuthConfigGenericOIDCName)
	}

	obj.Name = AuthConfigKeyCloakOIDCName
	obj.Type = managementClient.KeyCloakOIDCConfigType

	return expandOIDCConfig(in, obj)
}
