package rancher2

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// Empty fields are sent explicitly so updates can clear an old authentication method.
type opentelekomcloudCredentialConfig struct {
	AccessKey   string `json:"accessKey" yaml:"accessKey"`
	AuthURL     string `json:"authUrl" yaml:"authUrl"`
	DomainName  string `json:"domainName" yaml:"domainName"`
	Password    string `json:"password" yaml:"password"`
	ProjectID   string `json:"projectId" yaml:"projectId"`
	ProjectName string `json:"projectName" yaml:"projectName"`
	Region      string `json:"region" yaml:"region"`
	SecretKey   string `json:"secretKey" yaml:"secretKey"`
	Username    string `json:"username" yaml:"username"`
}

func cloudCredentialOpentelekomcloudFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "T-Cloud Public access key. Requires secret_key.",
		},
		"auth_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "T-Cloud Public identity service endpoint.",
		},
		"domain_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "T-Cloud Public domain name for username/password authentication.",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "T-Cloud Public password. Requires username and domain_name.",
		},
		"project_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "T-Cloud Public project ID.",
		},
		"project_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "T-Cloud Public project name.",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "T-Cloud Public region.",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "T-Cloud Public secret key. Requires access_key.",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "T-Cloud Public username. Requires password and domain_name.",
		},
	}
}
