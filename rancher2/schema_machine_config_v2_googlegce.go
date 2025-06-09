package rancher2

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func machineConfigV2GoogleGCEFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "GCE Instance External IP",
		},
		"auth_encoded_json": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "GCE service account auth json file path",
		},
		"disk_size": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "GCE Instance Disk Size (in GB)",
		},
		"disk_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "GCE Instance Disk Type",
		},
		"external_firewall_rule_prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A prefix to be added to firewall rules created when exposing ports publicly. Required if exposing ports publicly.",
		},
		"internal_firewall_rule_prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A prefix to be added to an internal firewall rule created to ensure virtual machines can communicate with one another.",
		},
		"labels": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A set of labels to be added to each VM, in the format of 'key1,value1,key2,value2'",
		},
		"machine_image": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "GCE instance image absolute URL",
		},
		"machine_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "GCE instance type",
		},
		"network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The network to provision virtual machines within",
		},
		"open_port": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "A list of ports to be opened publicly. 'external_firewall_rule_prefix' must also be set",
		},
		"preemptable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates if the virtual machine can be preempted",
		},
		"project": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The GCP project to create virtual machines within",
		},
		"scopes": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Access scopes to be set on the virtual machine",
		},
		"sub_network": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The subnetwork to provision virtual machines within",
		},
		"tags": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A set of network tags to be added to each VM, in the format of 'tag1,tag2'",
		},
		"use_existing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates if an existing VM should be used. This is not currently support in Rancher.",
		},
		"use_internal_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates if the virtual machines should use an internal IP",
		},
		"use_internal_ip_only": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates if the virtual machines should use an internal IP only and not be assigned a public IP",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "GCE user-data file path",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The username to be set when logging into the virtual machines",
		},
		"zone": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The region and zone to create virtual machines within (e.g. us-east1-b)",
		},
	}
}
