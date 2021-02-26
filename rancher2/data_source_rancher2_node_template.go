package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func dataSourceRancher2NodeTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2NodeTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_credential_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"driver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_env": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"engine_insecure_registry": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"engine_install_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_label": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"engine_opt": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"engine_registry_mirror": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"engine_storage_driver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_taints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: taintFields(),
				},
			},
			"use_internal_ip_address": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2NodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name": name,
	}
	listOpts := NewListOpts(filters)

	nodeTemplates, err := client.NodeTemplate.List(listOpts)
	if err != nil {
		return err
	}

	count := len(nodeTemplates.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] node template with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d node template with name \"%s\" ", count, name)
	}

	return flattenDataSourceNodeTemplate(d, &nodeTemplates.Data[0])
}

func flattenDataSourceNodeTemplate(d *schema.ResourceData, in *managementClient.NodeTemplate) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("driver", in.Driver)

	if len(in.AuthCertificateAuthority) > 0 {
		d.Set("auth_certificate_authority", in.AuthCertificateAuthority)
	}

	if len(in.AuthKey) > 0 {
		d.Set("auth_key", in.AuthKey)
	}

	if len(in.CloudCredentialID) > 0 {
		d.Set("cloud_credential_id", in.CloudCredentialID)
	}

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.EngineEnv) > 0 {
		err := d.Set("engine_env", toMapInterface(in.EngineEnv))
		if err != nil {
			return err
		}
	}

	if len(in.EngineInsecureRegistry) > 0 {
		err := d.Set("engine_insecure_registry", toArrayInterface(in.EngineInsecureRegistry))
		if err != nil {
			return err
		}
	}

	if len(in.EngineInstallURL) > 0 {
		d.Set("engine_install_url", in.EngineInstallURL)
	}

	if len(in.EngineLabel) > 0 {
		err := d.Set("engine_label", toMapInterface(in.EngineLabel))
		if err != nil {
			return err
		}
	}

	if len(in.EngineOpt) > 0 {
		err := d.Set("engine_opt", toMapInterface(in.EngineOpt))
		if err != nil {
			return err
		}
	}

	if len(in.EngineRegistryMirror) > 0 {
		err := d.Set("engine_registry_mirror", toArrayInterface(in.EngineRegistryMirror))
		if err != nil {
			return err
		}
	}

	if len(in.EngineStorageDriver) > 0 {
		d.Set("engine_storage_driver", in.EngineStorageDriver)
	}

	d.Set("use_internal_ip_address", *in.UseInternalIPAddress)

	if len(in.Annotations) > 0 {
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
