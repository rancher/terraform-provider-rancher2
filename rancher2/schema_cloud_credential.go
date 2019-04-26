package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Types

type amazonec2CredentialConfig struct {
	AccessKey string `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
}

type azureCredentialConfig struct {
	ClientID       string `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	SubscriptionID string `json:"subscriptionId,omitempty" yaml:"subscriptionId,omitempty"`
}

type digitaloceanCredentialConfig struct {
	AccessToken string `json:"accessToken,omitempty" yaml:"accessToken,omitempty"`
}

type CloudCredential struct {
	managementClient.CloudCredential
	Amazonec2CredentialConfig    *amazonec2CredentialConfig    `json:"amazonec2credentialConfig,omitempty" yaml:"amazonec2credentialConfig,omitempty"`
	AzureCredentialConfig        *azureCredentialConfig        `json:"azurecredentialConfig,omitempty" yaml:"azurecredentialConfig,omitempty"`
	DigitaloceanCredentialConfig *digitaloceanCredentialConfig `json:"digitaloceancredentialConfig,omitempty" yaml:"digitaloceancredentialConfig,omitempty"`
}

//Schemas

func cloudCredentialAmazonec2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "AWS Access Key",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "AWS Secret Key",
		},
	}

	return s
}

func cloudCredentialAzureFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"client_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Azure Service Principal Account ID",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Azure Service Principal Account password",
		},
		"subscription_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Azure Subscription ID",
		},
	}

	return s
}

func cloudCredentialDigitaloceanFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_token": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Digital Ocean access token",
		},
	}

	return s
}

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
			ConflictsWith: []string{"azure_credential_config", "digitalocean_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialAmazonec2Fields(),
			},
		},
		"azure_credential_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_credential_config", "digitalocean_credential_config"},
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
			ConflictsWith: []string{"amazonec2_credential_config", "azure_credential_config"},
			Elem: &schema.Resource{
				Schema: cloudCredentialDigitaloceanFields(),
			},
		},
		"driver": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
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
