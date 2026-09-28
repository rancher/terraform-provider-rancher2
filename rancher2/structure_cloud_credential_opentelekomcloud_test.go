package rancher2

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestTCloudPublicCredentialRoundTrip(t *testing.T) {
	values := map[string]interface{}{
		"access_key": "test-ak", "secret_key": "test-sk",
		"username": "test-user", "password": "test-password",
		"domain_name": "test-domain", "project_id": "test-project-id",
		"project_name": "test-project", "region": "eu-de",
		"auth_url": "https://identity.example.invalid/v3",
	}
	d := schema.TestResourceDataRaw(t, cloudCredentialFields(), map[string]interface{}{
		"name": "tcloud-public", "tcloud_public_credential_config": []interface{}{values},
	})
	credential := expandCloudCredential(d)
	require.Equal(t, "opentelekomcloud", d.Get("driver"))
	require.Equal(t, &opentelekomcloudCredentialConfig{
		AccessKey: "test-ak", SecretKey: "test-sk", Username: "test-user", Password: "test-password",
		DomainName: "test-domain", ProjectID: "test-project-id", ProjectName: "test-project",
		Region: "eu-de", AuthURL: "https://identity.example.invalid/v3",
	}, credential.OpentelekomcloudCredentialConfig)

	encoded, err := json.Marshal(credential)
	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(encoded, &body))
	require.NotContains(t, body, "tcloud_public_credential_config")
	require.Equal(t, map[string]interface{}{
		"accessKey": "test-ak", "secretKey": "test-sk", "username": "test-user", "password": "test-password",
		"domainName": "test-domain", "projectId": "test-project-id", "projectName": "test-project",
		"region": "eu-de", "authUrl": "https://identity.example.invalid/v3",
	}, body["opentelekomcloudcredentialConfig"])
	require.NoError(t, flattenCloudCredential(d, credential))
	require.Equal(t, []interface{}{values}, d.Get("tcloud_public_credential_config"))
}

func TestTCloudPublicCredentialRedactedSecrets(t *testing.T) {
	prior := []interface{}{map[string]interface{}{
		"password": "prior-password", "secret_key": "prior-secret", "domain_name": "old-domain",
	}}
	response := &opentelekomcloudCredentialConfig{AccessKey: "test-ak", Region: "eu-de"}
	got := flattenCloudCredentialOpentelekomcloud(response, prior)[0].(map[string]interface{})
	require.Equal(t, "prior-password", got["password"])
	require.Equal(t, "prior-secret", got["secret_key"])
	require.Equal(t, "", got["domain_name"])
	require.Equal(t, "old-domain", prior[0].(map[string]interface{})["domain_name"])

	response.Password = "rotated-password"
	response.SecretKey = "rotated-secret"
	got = flattenCloudCredentialOpentelekomcloud(response, prior)[0].(map[string]interface{})
	require.Equal(t, "rotated-password", got["password"])
	require.Equal(t, "rotated-secret", got["secret_key"])
	require.Equal(t, "prior-password", prior[0].(map[string]interface{})["password"])
}

func TestTCloudPublicCredentialEmptyValues(t *testing.T) {
	for _, input := range [][]interface{}{nil, {}, {nil}, {map[string]interface{}{}}} {
		require.Equal(t, &opentelekomcloudCredentialConfig{}, expandCloudCredentialOpentelekomcloud(input))
		require.Empty(t, flattenCloudCredentialOpentelekomcloud(nil, input))
		got := flattenCloudCredentialOpentelekomcloud(&opentelekomcloudCredentialConfig{}, input)
		require.Len(t, got, 1)
		require.Len(t, got[0], 9)
	}
	encoded, err := json.Marshal(&opentelekomcloudCredentialConfig{AccessKey: "test-ak", SecretKey: "test-sk"})
	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(encoded, &body))
	require.Len(t, body, 9)
	require.Equal(t, "", body["password"])
	require.Equal(t, "", body["username"])
	require.Equal(t, "", body["domainName"])
	require.NotContains(t, body, "authMethod")
	require.NotContains(t, body, "domainId")
	require.NotContains(t, body, "token")
	require.Error(t, validateCloudCredentialOpentelekomcloud(nil))
}
