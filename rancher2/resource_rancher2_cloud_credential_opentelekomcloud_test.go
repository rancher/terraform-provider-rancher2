package rancher2

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/require"
)

type tcloudPublicTestNodeDriver struct {
	managementClient.NodeDriverOperations
	id  string
	err error
}

func (d *tcloudPublicTestNodeDriver) ByID(id string) (*managementClient.NodeDriver, error) {
	d.id = id
	return &managementClient.NodeDriver{State: "active"}, d.err
}

func TestTCloudPublicCredentialCRUDAndImport(t *testing.T) {
	const id = "cattle-global-data:cc-test"
	var mu sync.Mutex
	var stored map[string]interface{}
	var failStatus atomic.Int32
	payloads := make(chan map[string]interface{}, 4)
	var baseURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if status := failStatus.Load(); status != 0 {
			w.WriteHeader(int(status))
			return
		}
		mu.Lock()
		defer mu.Unlock()
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			payloads <- body
			stored = body
		case http.MethodGet:
			if stored == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		case http.MethodDelete:
			stored = nil
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		response := make(map[string]interface{}, len(stored)+3)
		for key, value := range stored {
			response[key] = value
		}
		response["id"] = id
		response["type"] = managementClient.CloudCredentialType
		response["links"] = map[string]string{"self": baseURL + "/v3/cloudCredentials/" + id}
		credential := make(map[string]interface{})
		for key, value := range stored["opentelekomcloudcredentialConfig"].(map[string]interface{}) {
			if key != "password" && key != "secretKey" {
				credential[key] = value
			}
		}
		response["opentelekomcloudcredentialConfig"] = credential
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode mock response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	baseURL = server.URL
	httpClient := server.Client()
	httpClient.Timeout = 5 * time.Second
	schemas := map[string]types.Schema{
		managementClient.CloudCredentialType: {
			Links:             map[string]string{"collection": server.URL + "/v3/cloudCredentials"},
			CollectionMethods: []string{http.MethodPost, http.MethodGet},
			ResourceMethods:   []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		},
	}
	nodeDriver := &tcloudPublicTestNodeDriver{}
	config := &Config{Client: Client{Management: &managementClient.Client{
		APIBaseClient: clientbase.APIBaseClient{Ops: &clientbase.APIOperations{
			Types: schemas, Client: httpClient, Opts: &clientbase.ClientOpts{URL: server.URL},
		}},
		NodeDriver: nodeDriver,
	}}}
	values := map[string]interface{}{"access_key": "test-ak", "secret_key": "test-sk", "region": "eu-de"}
	d := schema.TestResourceDataRaw(t, cloudCredentialFields(), map[string]interface{}{
		"name": "tcloud-public", "tcloud_public_credential_config": []interface{}{values},
	})

	nodeDriver.err = errors.New("test driver unavailable")
	require.ErrorContains(t, resourceRancher2CloudCredentialCreate(d, config), "test driver unavailable")
	require.Empty(t, d.Id())
	require.Empty(t, payloads)
	nodeDriver.err = nil

	require.NoError(t, resourceRancher2CloudCredentialCreate(d, config))
	require.Equal(t, "opentelekomcloud", nodeDriver.id)
	require.Equal(t, id, d.Id())
	created := <-payloads
	require.Contains(t, created, "opentelekomcloudcredentialConfig")
	require.NotContains(t, created, "tcloud_public_credential_config")
	require.Equal(t, "test-sk", d.Get("tcloud_public_credential_config.0.secret_key"))

	values["secret_key"] = "rotated-test-sk"
	require.NoError(t, d.Set("tcloud_public_credential_config", []interface{}{values}))
	require.NoError(t, resourceRancher2CloudCredentialUpdate(d, config))
	updated := <-payloads
	require.Equal(t, "rotated-test-sk", updated["opentelekomcloudcredentialConfig"].(map[string]interface{})["secretKey"])
	require.Equal(t, "rotated-test-sk", d.Get("tcloud_public_credential_config.0.secret_key"))

	passwordValues := map[string]interface{}{"username": "test-user", "password": "test-password", "domain_name": "test-domain"}
	require.NoError(t, d.Set("tcloud_public_credential_config", []interface{}{passwordValues}))
	require.NoError(t, resourceRancher2CloudCredentialUpdate(d, config))
	switched := (<-payloads)["opentelekomcloudcredentialConfig"].(map[string]interface{})
	require.Equal(t, "", switched["accessKey"])
	require.Equal(t, "", switched["secretKey"])
	require.Equal(t, "test-password", switched["password"])
	require.Equal(t, "test-password", d.Get("tcloud_public_credential_config.0.password"))
	require.Equal(t, "", d.Get("tcloud_public_credential_config.0.secret_key"))

	imported := schema.TestResourceDataRaw(t, cloudCredentialFields(), nil)
	imported.SetId(id + ".opentelekomcloud")
	resources, err := resourceRancher2CloudCredentialsImport(imported, config)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	require.Equal(t, id, imported.Id())
	require.Equal(t, "opentelekomcloud", imported.Get("driver"))
	require.Equal(t, "test-user", imported.Get("tcloud_public_credential_config.0.username"))
	require.Equal(t, "", imported.Get("tcloud_public_credential_config.0.password"))

	failStatus.Store(http.StatusBadGateway)
	require.Error(t, resourceRancher2CloudCredentialRead(d, config))
	require.Equal(t, id, d.Id())
	require.Error(t, resourceRancher2CloudCredentialUpdate(d, config))
	failStatus.Store(0)
	require.NoError(t, resourceRancher2CloudCredentialDelete(d, config))
	require.Empty(t, d.Id())
	d.SetId(id)
	require.NoError(t, resourceRancher2CloudCredentialRead(d, config))
	require.Empty(t, d.Id())
}
