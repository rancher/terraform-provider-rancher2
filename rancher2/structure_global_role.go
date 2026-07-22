package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenGlobalRole(d *schema.ResourceData, in *managementClient.GlobalRole) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening global role: Input setting is nil")
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("builtin", in.Builtin)
	d.Set("new_user_default", in.NewUserDefault)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	err := d.Set("rules", flattenPolicyRules(in.Rules))
	if err != nil {
		return err
	}

	d.Set("uuid", in.UUID)

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	if len(in.InheritedClusterRoles) > 0 {
		err = d.Set("inherited_cluster_roles", toArrayInterface(in.InheritedClusterRoles))
		if err != nil {
			return err
		}
	}

	err = d.Set("inherited_namespaced_rules", flattenInheritedNamespacedRules(in.InheritedNamespacedRules))
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandGlobalRole(in *schema.ResourceData) *managementClient.GlobalRole {
	obj := &managementClient.GlobalRole{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.NewUserDefault = in.Get("new_user_default").(bool)

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("rules").([]interface{}); ok && len(v) > 0 {
		obj.Rules = expandPolicyRules(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, k := in.Get("inherited_cluster_roles").([]interface{}); k && len(v) > 0 {
		obj.InheritedClusterRoles = toArrayString(v)
	}

	if v, ok := in.Get("inherited_namespaced_rules").(*schema.Set); ok && v.Len() > 0 {
		obj.InheritedNamespacedRules = expandInheritedNamespacedRules(v.List())
	}

	obj.UUID = in.Get("uuid").(string)

	return obj
}

func flattenInheritedNamespacedRules(in map[string][]managementClient.PolicyRule) []interface{} {
	if len(in) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, 0, len(in))
	for namespace, rules := range in {
		out = append(out, map[string]interface{}{
			"namespace": namespace,
			"rules":     flattenPolicyRules(rules),
		})
	}

	return out
}

func expandInheritedNamespacedRules(in []interface{}) map[string][]managementClient.PolicyRule {
	if len(in) == 0 {
		return map[string][]managementClient.PolicyRule{}
	}

	out := make(map[string][]managementClient.PolicyRule, len(in))
	for _, ruleSet := range in {
		ruleSetMap := ruleSet.(map[string]interface{})
		policyRules, _ := ruleSetMap["rules"].([]interface{})
		out[ruleSetMap["namespace"].(string)] = expandPolicyRules(policyRules)
	}

	return out
}
