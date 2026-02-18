package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigCognito(d *schema.ResourceData, in *managementClient.CognitoConfig) error {
	d.SetId(AuthConfigCognitoName)
	d.Set("name", AuthConfigCognitoName)
	d.Set("type", managementClient.CognitoConfigType)

	if err := flattenOIDCConfig(d, in); err != nil {
		return fmt.Errorf("flattening AuthConfig for Cognito: %s", err)
	}

	return nil
}

// Expanders

func expandAuthConfigCognito(in *schema.ResourceData) (*managementClient.CognitoConfig, error) {
	obj := &managementClient.CognitoConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding :%s Auth Config: Input ResourceData is nil", AuthConfigCognitoName)
	}

	obj.Name = AuthConfigCognitoName
	obj.Type = managementClient.CognitoConfigType

	return expandOIDCConfig(in, obj)
}
