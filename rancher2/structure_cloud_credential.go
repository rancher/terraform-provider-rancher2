package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

// Flatteners

func flattenCloudCredentialAmazonec2(in *amazonec2CredentialConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	return []interface{}{obj}
}

func flattenCloudCredentialAzure(in *azureCredentialConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.ClientID) > 0 {
		obj["client_id"] = in.ClientID
	}

	if len(in.ClientSecret) > 0 {
		obj["client_secret"] = in.ClientSecret
	}

	if len(in.SubscriptionID) > 0 {
		obj["subscription_id"] = in.SubscriptionID
	}

	return []interface{}{obj}
}

func flattenCloudCredentialDigitalocean(in *digitaloceanCredentialConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessToken) > 0 {
		obj["access_token"] = in.AccessToken
	}

	return []interface{}{obj}
}

func flattenCloudCredential(d *schema.ResourceData, in *CloudCredential) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("name", in.Name)
	if err != nil {
		return err
	}

	if len(in.Description) > 0 {
		err = d.Set("description", in.Description)
		if err != nil {
			return err
		}
	}

	driver := d.Get("driver").(string)

	switch driver {
	case amazonec2ConfigDriver:
		in.Amazonec2CredentialConfig.SecretKey = d.Get("amazonec2_credential_config.0.secret_key").(string)
		d.Set("amazonec2_credential_config", flattenCloudCredentialAmazonec2(in.Amazonec2CredentialConfig))
	case azureConfigDriver:
		in.AzureCredentialConfig.ClientSecret = d.Get("azure_credential_config.0.client_secret").(string)
		d.Set("azure_credential_config", flattenCloudCredentialAzure(in.AzureCredentialConfig))
	case digitaloceanConfigDriver:
		in.DigitaloceanCredentialConfig.AccessToken = d.Get("digitalocean_credential_config.0.access_token").(string)
		d.Set("digitalocean_credential_config", flattenCloudCredentialDigitalocean(in.DigitaloceanCredentialConfig))
	default:
		return fmt.Errorf("[ERROR] Unsupported driver on cloud credential: %s", driver)
	}

	if len(in.Annotations) > 0 {
		err = d.Set("annotations", toMapInterface(in.Annotations))
		if err != nil {
			return err
		}
	}

	if len(in.Labels) > 0 {
		err = d.Set("labels", toMapInterface(in.Labels))
		if err != nil {
			return err
		}
	}

	return nil
}

// Expanders

func expandCloudCredentialAmazonec2(p []interface{}) *amazonec2CredentialConfig {
	obj := &amazonec2CredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	return obj
}

func expandCloudCredentialAzure(p []interface{}) *azureCredentialConfig {
	obj := &azureCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["client_id"].(string); ok && len(v) > 0 {
		obj.ClientID = v
	}

	if v, ok := in["client_secret"].(string); ok && len(v) > 0 {
		obj.ClientSecret = v
	}

	if v, ok := in["subscription_id"].(string); ok && len(v) > 0 {
		obj.SubscriptionID = v
	}

	return obj
}

func expandCloudCredentialDigitalocean(p []interface{}) *digitaloceanCredentialConfig {
	obj := &digitaloceanCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_token"].(string); ok && len(v) > 0 {
		obj.AccessToken = v
	}

	return obj
}

func expandCloudCredential(in *schema.ResourceData) *CloudCredential {
	obj := &CloudCredential{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}
	obj.Name = in.Get("name").(string)

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("amazonec2_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.Amazonec2CredentialConfig = expandCloudCredentialAmazonec2(v)
		in.Set("driver", amazonec2ConfigDriver)
	}

	if v, ok := in.Get("azure_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.AzureCredentialConfig = expandCloudCredentialAzure(v)
		in.Set("driver", azureConfigDriver)
	}

	if v, ok := in.Get("digitalocean_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.DigitaloceanCredentialConfig = expandCloudCredentialDigitalocean(v)
		in.Set("driver", digitaloceanConfigDriver)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
