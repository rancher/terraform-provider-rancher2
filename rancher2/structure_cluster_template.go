package rancher2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Sorters

type ByNameClusterTemplateRevisions []managementClient.ClusterTemplateRevision

func (a ByNameClusterTemplateRevisions) Len() int { return len(a) }
func (a ByNameClusterTemplateRevisions) Less(i, j int) bool {
	return strings.Compare(a[i].Name, a[j].Name) == -1
}
func (a ByNameClusterTemplateRevisions) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Flatteners

func flattenQuestions(p []managementClient.Question) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		obj["default"] = in.Default
		obj["required"] = in.Required
		obj["type"] = in.Type
		obj["variable"] = in.Variable

		out[i] = obj
	}

	return out
}

func flattenClusterSpecBase(in *managementClient.ClusterSpecBase, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.LocalClusterAuthEndpoint != nil {
		obj["cluster_auth_endpoint"] = flattenClusterAuthEndpoint(in.LocalClusterAuthEndpoint)
	}

	if len(in.DefaultClusterRoleForProjectMembers) > 0 {
		obj["default_cluster_role_for_project_members"] = in.DefaultClusterRoleForProjectMembers
	}

	if len(in.DefaultPodSecurityPolicyTemplateID) > 0 {
		obj["default_pod_security_policy_template_id"] = in.DefaultPodSecurityPolicyTemplateID
	}

	if len(in.DesiredAgentImage) > 0 {
		obj["desired_agent_image"] = in.DesiredAgentImage
	}

	if len(in.DesiredAuthImage) > 0 {
		obj["desired_auth_image"] = in.DesiredAuthImage
	}

	if len(in.DockerRootDir) > 0 {
		obj["docker_root_dir"] = in.DockerRootDir
	}

	obj["enable_cluster_alerting"] = in.EnableClusterAlerting
	obj["enable_cluster_monitoring"] = in.EnableClusterMonitoring
	obj["enable_network_policy"] = *in.EnableNetworkPolicy

	if in.RancherKubernetesEngineConfig != nil {
		v, ok := obj["rke_config"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		rkeConfig, err := flattenClusterRKEConfig(in.RancherKubernetesEngineConfig, v)
		if err != nil {
			return []interface{}{}, err
		}
		obj["rke_config"] = rkeConfig
	}

	obj["windows_prefered_cluster"] = in.WindowsPreferedCluster

	return []interface{}{obj}, nil
}

func flattenClusterTemplateRevisions(input []managementClient.ClusterTemplateRevision, defaultCtrID string, p []interface{}) ([]interface{}, error) {
	if len(input) == 0 || p == nil {
		return []interface{}{}, nil
	}

	if len(defaultCtrID) == 0 {
		return []interface{}{}, fmt.Errorf("Default Cluster Template Revision ID can't be empty")
	}

	// Sorting input array by data interface
	pIndexID := map[string]int{}
	pIndexName := map[string]int{}
	for i := range p {
		if row, ok := p[i].(map[string]interface{}); ok {
			if v, ok := row["id"].(string); ok && len(v) > 0 {
				pIndexID[v] = i
			}
			if v, ok := row["name"].(string); ok {
				pIndexName[v] = i
			}
		}
	}
	inputLen := len(input)
	pLen := len(p)
	sortedInput := make([]managementClient.ClusterTemplateRevision, inputLen)
	newCTR := []managementClient.ClusterTemplateRevision{}
	lastIndex := 0
	for i := range sortedInput {
		if v, ok := pIndexID[input[i].ID]; ok {
			if v > i && pLen > inputLen {
				v = v - (v - i)
			}
			sortedInput[v] = input[i]
			lastIndex++
			continue
		}
		if v, ok := pIndexName[input[i].Name]; ok {
			if v > i && pLen > inputLen {
				v = v - (v - i)
			}
			sortedInput[v] = input[i]
			lastIndex++
			continue
		}
		newCTR = append(newCTR, input[i])
	}

	for i := range newCTR {
		sortedInput[lastIndex+i] = newCTR[i]
	}

	out := make([]interface{}, len(sortedInput))
	for i, in := range sortedInput {
		var obj map[string]interface{}
		if v, ok := pIndexName[in.Name]; ok {
			if row, ok := p[v].(map[string]interface{}); ok {
				obj = row
			}
		}
		if obj == nil {
			obj = make(map[string]interface{})
		}

		obj["default"] = false
		if len(in.ID) > 0 {
			obj["id"] = in.ID
			if in.ID == defaultCtrID {
				obj["default"] = true
			}
		}

		if in.ClusterConfig != nil {
			v, ok := obj["cluster_config"].([]interface{})
			if !ok {
				v = []interface{}{}
			}
			clusterConfig, err := flattenClusterSpecBase(in.ClusterConfig, v)
			if err != nil {
				return []interface{}{}, err
			}
			obj["cluster_config"] = clusterConfig
		}

		if len(in.ClusterTemplateID) > 0 {
			obj["cluster_template_id"] = in.ClusterTemplateID
		}

		if in.Enabled != nil {
			obj["enabled"] = *in.Enabled
		}

		if len(in.Name) > 0 {
			obj["name"] = in.Name
		}

		if len(in.Questions) > 0 {
			obj["questions"] = flattenQuestions(in.Questions)
		}

		obj["annotations"] = toMapInterface(in.Annotations)
		obj["labels"] = toMapInterface(in.Labels)

		out[i] = obj
	}

	return out, nil
}

func flattenClusterTemplate(d *schema.ResourceData, in *managementClient.ClusterTemplate, revisions []managementClient.ClusterTemplateRevision) error {
	if len(in.ID) > 0 {
		d.SetId(in.ID)
	}

	d.Set("default_revision_id", in.DefaultRevisionID)

	v, ok := d.Get("template_revisions").([]interface{})
	if !ok {
		v = []interface{}{}
	}
	templateRevisions, err := flattenClusterTemplateRevisions(revisions, in.DefaultRevisionID, v)
	if err != nil {
		return err
	}
	d.Set("template_revisions", templateRevisions)
	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	err = d.Set("members", flattenMembers(in.Members))
	if err != nil {
		return err
	}

	d.Set("name", in.Name)

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandQuestions(p []interface{}) []managementClient.Question {
	if len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := make([]managementClient.Question, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["default"].(string); ok && len(v) > 0 {
			obj[i].Default = v
		}

		if v, ok := in["required"].(bool); ok {
			obj[i].Required = v
		}

		if v, ok := in["type"].(string); ok && len(v) > 0 {
			obj[i].Type = v
		}

		if v, ok := in["variable"].(string); ok && len(v) > 0 {
			obj[i].Variable = v
		}
	}

	return obj
}

func expandClusterSpecBase(p []interface{}) (*managementClient.ClusterSpecBase, error) {
	obj := &managementClient.ClusterSpecBase{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_auth_endpoint"].([]interface{}); ok && len(v) > 0 {
		obj.LocalClusterAuthEndpoint = expandClusterAuthEndpoint(v)
	}

	if v, ok := in["default_cluster_role_for_project_members"].(string); ok && len(v) > 0 {
		obj.DefaultClusterRoleForProjectMembers = v
	}

	if v, ok := in["default_pod_security_policy_template_id"].(string); ok && len(v) > 0 {
		obj.DefaultPodSecurityPolicyTemplateID = v
	}

	if v, ok := in["desired_agent_image"].(string); ok && len(v) > 0 {
		obj.DesiredAgentImage = v
	}

	if v, ok := in["desired_auth_image"].(string); ok && len(v) > 0 {
		obj.DesiredAuthImage = v
	}

	if v, ok := in["docker_root_dir"].(string); ok && len(v) > 0 {
		obj.DockerRootDir = v
	}

	if v, ok := in["enable_cluster_alerting"].(bool); ok {
		obj.EnableClusterAlerting = v
	}

	if v, ok := in["enable_cluster_monitoring"].(bool); ok {
		obj.EnableClusterMonitoring = v
	}

	if v, ok := in["enable_network_policy"].(bool); ok {
		obj.EnableNetworkPolicy = &v
	}

	if v, ok := in["rke_config"].([]interface{}); ok && len(v) > 0 {
		rkeConfig, err := expandClusterRKEConfig(v, "")
		if err != nil {
			return nil, err
		}
		obj.RancherKubernetesEngineConfig = rkeConfig
	}

	if v, ok := in["windows_prefered_cluster"].(bool); ok {
		obj.WindowsPreferedCluster = v
	}

	return obj, nil
}

func expandClusterTemplateRevisions(p []interface{}) (int, []managementClient.ClusterTemplateRevision, error) {
	if len(p) == 0 || p[0] == nil {
		return 0, []managementClient.ClusterTemplateRevision{}, nil
	}

	obj := make([]managementClient.ClusterTemplateRevision, len(p))

	indexDefault := 0
	hasDefault := false
	names := map[string]int{}
	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["id"].(string); ok && len(v) > 0 {
			obj[i].ID = v
		}

		if v, ok := in["cluster_config"].([]interface{}); ok && len(v) > 0 {
			var err error
			obj[i].ClusterConfig, err = expandClusterSpecBase(v)
			if err != nil {
				return 0, nil, err
			}
		}

		if v, ok := in["cluster_template_id"].(string); ok && len(v) > 0 {
			obj[i].ClusterTemplateID = v
		}

		if v, ok := in["default"].(bool); ok && v {
			if hasDefault {
				return 0, nil, fmt.Errorf("[ERROR] Expanding cluster template revisions: more than one default defined")
			}
			hasDefault = true
			indexDefault = i
		}

		if v, ok := in["enabled"].(bool); ok {
			obj[i].Enabled = &v
		}

		if v, ok := in["name"].(string); ok {
			obj[i].Name = v
			names[v]++
			if names[v] > 1 {
				return 0, nil, fmt.Errorf("[ERROR] Expanding cluster template revisions: name \"%s\" is repeated", v)
			}
		}

		if v, ok := in["questions"].([]interface{}); ok && len(v) > 0 {
			obj[i].Questions = expandQuestions(v)
		}

		if v, ok := in["annotations"].(map[string]interface{}); ok && len(v) > 0 {
			obj[i].Annotations = toMapString(v)
		}

		if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
			obj[i].Labels = toMapString(v)
		}
	}

	if !hasDefault {
		return 0, nil, fmt.Errorf("[ERROR] Expanding cluster template revisions: NO default defined")
	}

	return indexDefault, obj, nil

}

func expandClusterTemplate(in *schema.ResourceData) (int, *managementClient.ClusterTemplate, []managementClient.ClusterTemplateRevision, error) {
	obj := &managementClient.ClusterTemplate{}
	if in == nil {
		return 0, nil, nil, nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("default_revision_id").(string); ok && len(v) > 0 {
		obj.DefaultRevisionID = v
	}

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("members").([]interface{}); ok && len(v) > 0 {
		obj.Members = expandMembers(v)
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	indexDefault, clusterTemplateRevisions, err := expandClusterTemplateRevisions(in.Get("template_revisions").([]interface{}))

	return indexDefault, obj, clusterTemplateRevisions, err
}
