package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenRecipients(p []managementClient.Recipient) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		obj["notifier_id"] = in.NotifierID

		if len(in.NotifierType) > 0 {
			obj["notifier_type"] = in.NotifierType
		}

		if len(in.Recipient) > 0 {
			obj["recipient"] = in.Recipient
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandRecipients(p []interface{}) []managementClient.Recipient {
	if len(p) == 0 || p[0] == nil {
		return []managementClient.Recipient{}
	}

	obj := make([]managementClient.Recipient, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		obj[i].NotifierID = in["notifier_id"].(string)
		obj[i].NotifierType = in["notifier_type"].(string)
		obj[i].Recipient = in["recipient"].(string)
	}

	return obj
}
