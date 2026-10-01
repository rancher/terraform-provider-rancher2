package rancher2

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	"github.com/stretchr/testify/require"
)

func TestTCloudPublicMachineConfigCRUD(t *testing.T) {
	const (
		id         = "fleet-default:nc-tcloud-test-abc"
		kind       = "OpentelekomcloudConfig"
		apiType    = "rke-machine-config.cattle.io.opentelekomcloudconfig"
		collection = "/v1/" + apiType
		selfPath   = "/machine-config-object"
	)
	var mu sync.Mutex
	var stored map[string]interface{}
	var revision int
	var failMethod atomic.Value
	failMethod.Store("")
	payloads := make(chan map[string]interface{}, 2)
	var baseURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		expectedPath := selfPath
		switch r.Method {
		case http.MethodPost:
			expectedPath = collection
		case http.MethodGet:
			expectedPath = collection + "/" + id
		}
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected API route: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method == failMethod.Load().(string) {
			w.WriteHeader(http.StatusBadRequest)
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
			revision++
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
		response := make(map[string]interface{}, len(stored)+4)
		for key, value := range stored {
			response[key] = value
		}
		metadata := make(map[string]interface{})
		for key, value := range stored["metadata"].(map[string]interface{}) {
			metadata[key] = value
		}
		metadata["name"] = "nc-tcloud-test-abc"
		metadata["resourceVersion"] = strconv.Itoa(revision)
		response["metadata"] = metadata
		response["id"], response["type"] = id, apiType
		response["links"] = map[string]string{"self": baseURL + selfPath}
		response["actions"] = map[string]string{"validate": baseURL + selfPath + "?action=validate"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode machine configuration: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	baseURL = server.URL
	httpClient := server.Client()
	httpClient.Timeout = 5 * time.Second
	client := &clientbase.APIBaseClient{Ops: &clientbase.APIOperations{
		Types: map[string]types.Schema{apiType: {
			Links:             map[string]string{"collection": server.URL + collection},
			CollectionMethods: []string{http.MethodPost},
			ResourceMethods:   []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		}},
		Client: httpClient, Opts: &clientbase.ClientOpts{URL: server.URL},
	}}
	config := &Config{
		Timeout: 2 * time.Second,
		Client:  Client{CatalogV2: map[string]*clientbase.APIBaseClient{"local": client}},
	}
	values := map[string]interface{}{
		"network_scope": "shared", "skip_default_sg": true, "skip_eip": true,
		"vpc_id": "test-vpc", "subnet_id": "test-subnet", "sec_groups": "test-sg",
		"tags": "environment.test", "image_name": "test-image",
	}
	d := schema.TestResourceDataRaw(t, machineConfigV2Fields(), map[string]interface{}{
		"generate_name": "tcloud-test", "tcloud_public_config": []interface{}{values},
		"labels":      map[string]interface{}{"environment": "test"},
		"annotations": map[string]interface{}{"ui.rancher/provider": "opentelekomcloud"},
	})
	failMethod.Store(http.MethodPost)
	require.Error(t, resourceRancher2MachineConfigV2Create(d, config))
	require.Empty(t, d.Id())
	failMethod.Store("")
	require.NoError(t, resourceRancher2MachineConfigV2Create(d, config))
	created := <-payloads
	require.Equal(t, kind, created["kind"])
	require.Equal(t, "rke-machine-config.cattle.io/v1", created["apiVersion"])
	require.NotContains(t, created, "spec")
	require.NotContains(t, created, "opentelekomcloudConfig")
	require.NotContains(t, created, "tcloud_public_config")
	require.Equal(t, "shared", created["networkScope"])
	require.Equal(t, true, created["skipDefaultSg"])
	require.Equal(t, id, d.Id())
	require.Equal(t, kind, d.Get("kind"))
	require.Equal(t, "nc-tcloud-test-abc", d.Get("name"))
	require.Equal(t, "1", d.Get("resource_version"))

	obj, err := getMachineConfigV2ByID(config, id, kind)
	require.NoError(t, err)
	require.Equal(t, obj.Resource, obj.OpentelekomcloudConfig.Resource)
	require.Equal(t, apiType, obj.Type)
	require.Equal(t, server.URL+selfPath, obj.Links["self"])
	require.Equal(t, server.URL+selfPath+"?action=validate", obj.Actions["validate"])
	require.Equal(t, "test-vpc", obj.OpentelekomcloudConfig.VPCID)

	failMethod.Store(http.MethodGet)
	require.Error(t, resourceRancher2MachineConfigV2Read(d, config))
	require.Equal(t, id, d.Id())
	require.Error(t, resourceRancher2MachineConfigV2Update(d, config))
	require.Error(t, resourceRancher2MachineConfigV2Delete(d, config))
	failMethod.Store(http.MethodPut)
	require.Error(t, resourceRancher2MachineConfigV2Update(d, config))
	require.Equal(t, id, d.Id())
	failMethod.Store("")
	values["network_scope"] = "machine"
	values["skip_default_sg"], values["skip_eip"] = false, false
	values["tags"], values["image_name"], values["vpc_id"] = "", "", ""
	values["root_volume_size"] = "0"
	require.NoError(t, d.Set("tcloud_public_config", []interface{}{values}))
	require.NoError(t, resourceRancher2MachineConfigV2Update(d, config))
	updated := <-payloads
	for key, expected := range map[string]interface{}{
		"skipDefaultSg": false, "skipEip": false, "tags": "", "imageName": "", "vpcId": "", "rootVolumeSize": "0",
	} {
		require.Contains(t, updated, key)
		require.Equal(t, expected, updated[key], key)
	}
	metadata := updated["metadata"].(map[string]interface{})
	require.Equal(t, "1", metadata["resourceVersion"])
	require.Equal(t, "fleet-default", metadata["namespace"])
	require.Equal(t, "opentelekomcloud", metadata["annotations"].(map[string]interface{})["ui.rancher/provider"])
	require.Equal(t, "test", metadata["labels"].(map[string]interface{})["environment"])
	require.Equal(t, "2", d.Get("resource_version"))
	require.Equal(t, false, d.Get("tcloud_public_config.0.skip_default_sg"))
	require.Equal(t, "", d.Get("tcloud_public_config.0.image_name"))
	require.Equal(t, "", d.Get("tcloud_public_config.0.tags"))

	failMethod.Store(http.MethodDelete)
	require.Error(t, resourceRancher2MachineConfigV2Delete(d, config))
	require.Equal(t, id, d.Id())
	failMethod.Store("")
	require.NoError(t, resourceRancher2MachineConfigV2Delete(d, config))
	require.Empty(t, d.Id())
	d.SetId(id)
	require.NoError(t, resourceRancher2MachineConfigV2Read(d, config))
	require.Empty(t, d.Id())
	d.SetId(id)
	require.NoError(t, resourceRancher2MachineConfigV2Delete(d, config))
	require.Empty(t, d.Id())
}
