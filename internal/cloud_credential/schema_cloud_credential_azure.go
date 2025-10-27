package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

type azureCredentialConfig struct {
	ClientID       string `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	SubscriptionID string `json:"subscriptionId,omitempty" yaml:"subscriptionId,omitempty"`
	Environment    string `json:"environment,omitempty" yaml:"environment,omitempty"`
	TenantID       string `json:"tenantId,omitempty" yaml:"tenantId,omitempty"`
}

//Schemas

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
		"environment": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Azure environment (e.g. AzurePublicCloud, AzureChinaCloud)",
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Azure Tenant ID",
		},
	}

	return s
}
