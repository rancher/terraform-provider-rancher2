package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Flatteners

func flattenCloudCredential(d *schema.ResourceData, in *CloudCredential) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	driver := d.Get("driver").(string)
	switch driver {
	case amazonec2ConfigDriver:
		v, ok := d.Get("amazonec2_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("amazonec2_credential_config", flattenCloudCredentialAmazonec2(in.Amazonec2CredentialConfig, v))
		if err != nil {
			return err
		}
	case azureConfigDriver:
		v, ok := d.Get("azure_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("azure_credential_config", flattenCloudCredentialAzure(in.AzureCredentialConfig, v))
		if err != nil {
			return err
		}
	case digitaloceanConfigDriver:
		v, ok := d.Get("digitalocean_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("digitalocean_credential_config", flattenCloudCredentialDigitalocean(in.DigitaloceanCredentialConfig, v))
		if err != nil {
			return err
		}
	case googleConfigDriver:
		v, ok := d.Get("google_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("google_credential_config", flattenCloudCredentialGoogle(in.GoogleCredentialConfig, v))
		if err != nil {
			return err
		}
	case harvesterConfigDriver:
		v, ok := d.Get("harvester_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("harvester_credential_config", flattenCloudCredentialHarvester(in.HarvesterCredentialConfig, v))
		if err != nil {
			return err
		}
	case linodeConfigDriver:
		v, ok := d.Get("linode_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("linode_credential_config", flattenCloudCredentialLinode(in.LinodeCredentialConfig, v))
		if err != nil {
			return err
		}
	case openstackConfigDriver:
		v, ok := d.Get("openstack_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("openstack_credential_config", flattenCloudCredentialOpenstack(in.OpenstackCredentialConfig, v))
		if err != nil {
			return err
		}
	case s3ConfigDriver:
		v, ok := d.Get("s3_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("s3_credential_config", flattenCloudCredentialS3(in.S3CredentialConfig, v))
		if err != nil {
			return err
		}
	case vmwarevsphereConfigDriver:
		v, ok := d.Get("vsphere_credential_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err := d.Set("vsphere_credential_config", flattenCloudCredentialVsphere(in.VmwarevsphereCredentialConfig, v))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("[ERROR] Unsupported driver on cloud credential: %s", driver)
	}

	if in.Annotations != nil {
		err := d.Set("annotations", toMapInterface(in.Annotations))
		if err != nil {
			return err
		}
	}

	if len(in.Labels) > 0 {
		err := d.Set("labels", toMapInterface(in.Labels))
		if err != nil {
			return err
		}
	}

	return nil
}

// Expanders

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

	if v, ok := in.Get("google_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.GoogleCredentialConfig = expandCloudCredentialGoogle(v)
		in.Set("driver", googleConfigDriver)
	}

	if v, ok := in.Get("harvester_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.HarvesterCredentialConfig = expandCloudCredentialHarvester(v)
		in.Set("driver", harvesterConfigDriver)
	}

	if v, ok := in.Get("linode_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.LinodeCredentialConfig = expandCloudCredentialLinode(v)
		in.Set("driver", linodeConfigDriver)
	}

	if v, ok := in.Get("openstack_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.OpenstackCredentialConfig = expandCloudCredentialOpenstack(v)
		in.Set("driver", openstackConfigDriver)
	}

	if v, ok := in.Get("s3_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.S3CredentialConfig = expandCloudCredentialS3(v)
		in.Set("driver", s3ConfigDriver)
	}

	if v, ok := in.Get("vsphere_credential_config").([]interface{}); ok && len(v) > 0 {
		obj.VmwarevsphereCredentialConfig = expandCloudCredentialVsphere(v)
		in.Set("driver", vmwarevsphereConfigDriver)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
