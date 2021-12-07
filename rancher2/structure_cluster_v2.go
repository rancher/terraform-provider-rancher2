package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	provisioningV1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
)

const (
	clusterV2Kind             = "Cluster"
	clusterV2APIVersion       = "provisioning.cattle.io/v1"
	clusterV2APIType          = "provisioning.cattle.io.cluster"
	clusterV2ClusterIDsep     = "/"
	clusterV2ActiveCondition  = "Updated"
	clusterV2CreatedCondition = "Created"
)

//Types

type ClusterV2 struct {
	norman.Resource
	provisioningV1.Cluster
}

// Flatteners

func flattenClusterV2(d *schema.ResourceData, in *ClusterV2) error {
	if in == nil {
		return nil
	}

	if len(in.ID) > 0 {
		d.SetId(in.ID)
	}
	d.Set("name", in.ObjectMeta.Name)
	d.Set("fleet_namespace", in.ObjectMeta.Namespace)
	err := d.Set("annotations", toMapInterface(in.ObjectMeta.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.ObjectMeta.Labels))
	if err != nil {
		return err
	}
	d.Set("resource_version", in.ObjectMeta.ResourceVersion)

	if len(in.Spec.KubernetesVersion) > 0 {
		d.Set("kubernetes_version", in.Spec.KubernetesVersion)
	}
	d.Set("local_auth_endpoint", flattenClusterV2LocalAuthEndpoint(in.Spec.LocalClusterAuthEndpoint))
	if in.Spec.RKEConfig != nil {
		d.Set("rke_config", flattenClusterV2RKEConfig(in.Spec.RKEConfig))
	}
	if in.Spec.AgentEnvVars != nil && len(in.Spec.AgentEnvVars) > 0 {
		d.Set("agent_env_vars", flattenEnvVarsV2(in.Spec.AgentEnvVars))
	}
	if len(in.Spec.CloudCredentialSecretName) > 0 {
		d.Set("cloud_credential_secret_name", in.Spec.CloudCredentialSecretName)
	}
	if len(in.Spec.DefaultPodSecurityPolicyTemplateName) > 0 {
		d.Set("default_pod_security_policy_template_name", in.Spec.DefaultPodSecurityPolicyTemplateName)
	}
	if len(in.Spec.DefaultClusterRoleForProjectMembers) > 0 {
		d.Set("default_cluster_role_for_project_members", in.Spec.DefaultClusterRoleForProjectMembers)
	}
	if in.Spec.EnableNetworkPolicy != nil {
		d.Set("enable_network_policy", *in.Spec.EnableNetworkPolicy)
	}
	if len(in.Status.ClusterName) > 0 {
		d.Set("cluster_v1_id", in.Status.ClusterName)
	}

	return nil
}

// Expanders

func expandClusterV2(in *schema.ResourceData) *ClusterV2 {
	if in == nil {
		return nil
	}
	obj := &ClusterV2{}

	if len(in.Id()) > 0 {
		obj.ID = in.Id()
	}
	obj.TypeMeta.Kind = clusterV2Kind
	obj.TypeMeta.APIVersion = clusterV2APIVersion

	obj.ObjectMeta.Name = in.Get("name").(string)
	obj.ObjectMeta.Namespace = in.Get("fleet_namespace").(string)
	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Annotations = toMapString(v)
	}
	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Labels = toMapString(v)
	}
	if v, ok := in.Get("resource_version").(string); ok {
		obj.ObjectMeta.ResourceVersion = v
	}
	if v, ok := in.Get("kubernetes_version").(string); ok && len(v) > 0 {
		obj.Spec.KubernetesVersion = v
	}
	if v, ok := in.Get("local_auth_endpoint").([]interface{}); ok && len(v) > 0 {
		obj.Spec.LocalClusterAuthEndpoint = expandClusterV2LocalAuthEndpoint(v)
	}
	if v, ok := in.Get("rke_config").([]interface{}); ok {
		obj.Spec.RKEConfig = expandClusterV2RKEConfig(v)
	}
	if v, ok := in.Get("agent_env_vars").([]interface{}); ok {
		obj.Spec.AgentEnvVars = expandEnvVarsV2(v)
	}
	if v, ok := in.Get("cloud_credential_secret_name").(string); ok && len(v) > 0 {
		obj.Spec.CloudCredentialSecretName = v
	}
	if v, ok := in.Get("default_pod_security_policy_template_name").(string); ok && len(v) > 0 {
		obj.Spec.DefaultPodSecurityPolicyTemplateName = v
	}
	if v, ok := in.Get("default_cluster_role_for_project_members").(string); ok && len(v) > 0 {
		obj.Spec.DefaultClusterRoleForProjectMembers = v
	}
	if v, ok := in.Get("enable_network_policy").(bool); ok {
		obj.Spec.EnableNetworkPolicy = &v
	}

	return obj
}
