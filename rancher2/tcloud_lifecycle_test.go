package rancher2

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	"github.com/stretchr/testify/require"
)

// tcloudMachineConfigTestHarness is a minimal in-process mock of the Rancher
// catalog-v2 API for the opentelekomcloud machine config API type. It stores
// one object, serves it on GET, accepts POST/PUT, and clears on DELETE.
type tcloudMachineConfigTestHarness struct {
	config *Config
	id     string
	api    string
	mu     sync.Mutex
	stored map[string]interface{}
	rev    int
}

func newTCloudMachineConfigTestHarness(t *testing.T) *tcloudMachineConfigTestHarness {
	t.Helper()
	const (
		id      = "fleet-default:nc-tcloud-acc-test"
		apiType = "rke-machine-config.cattle.io.opentelekomcloudconfig"
	)
	h := &tcloudMachineConfigTestHarness{id: id, api: apiType}
	var baseURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.mu.Lock()
		defer h.mu.Unlock()
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			h.stored = body
			h.rev++
		case http.MethodGet:
			if h.stored == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		case http.MethodDelete:
			h.stored = nil
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		response := make(map[string]interface{}, len(h.stored)+4)
		for key, value := range h.stored {
			response[key] = value
		}
		metadata := make(map[string]interface{})
		for key, value := range h.stored["metadata"].(map[string]interface{}) {
			metadata[key] = value
		}
		metadata["name"] = "nc-tcloud-acc-test"
		metadata["resourceVersion"] = strconv.Itoa(h.rev)
		response["metadata"] = metadata
		response["id"], response["type"] = id, apiType
		response["links"] = map[string]string{"self": baseURL + "/mc"}
		response["actions"] = map[string]string{"validate": baseURL + "/mc?action=validate"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode machine configuration: %v", err)
		}
	}))
	baseURL = server.URL
	t.Cleanup(server.Close)
	httpClient := server.Client()
	httpClient.Timeout = 5 * time.Second
	client := &clientbase.APIBaseClient{Ops: &clientbase.APIOperations{
		Types: map[string]types.Schema{apiType: {
			Links:             map[string]string{"collection": server.URL + "/v1/" + apiType},
			CollectionMethods: []string{http.MethodPost},
			ResourceMethods:   []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		}},
		Client: httpClient, Opts: &clientbase.ClientOpts{URL: server.URL},
	}}
	h.config = &Config{
		Timeout: 2 * time.Second,
		Client:  Client{CatalogV2: map[string]*clientbase.APIBaseClient{"local": client}},
	}
	return h
}

func (h *tcloudMachineConfigTestHarness) newResourceData(t *testing.T, values map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, machineConfigV2Fields(), map[string]interface{}{
		"generate_name": "tcloud-acc-test", "fleet_namespace": "fleet-default",
		"tcloud_public_config": []interface{}{values},
		"labels":               map[string]interface{}{"environment": "test"},
		"annotations":          map[string]interface{}{"ui.rancher/provider": "opentelekomcloud"},
	})
}

func (h *tcloudMachineConfigTestHarness) sharedValues() map[string]interface{} {
	return map[string]interface{}{
		"network_scope": "shared", "skip_default_sg": true, "skip_eip": true,
		"vpc_id": "test-vpc", "subnet_id": "test-subnet", "sec_groups": "test-sg",
		"image_name": "test-image", "tags": "environment.test",
	}
}

// TestTCloudPublicMachineConfigLifecycleRefreshNoOp covers the TC-10
// "refresh no-op / drift" requirement. After a create, two consecutive reads
// (refreshes) against an unchanged server must produce identical state:
// the provider's flatten is idempotent, and a clean refresh must not report
// drift. We verify by reading the same fields twice and asserting equality,
// plus asserting the state still matches the values we created with.
func TestTCloudPublicMachineConfigLifecycleRefreshNoOp(t *testing.T) {
	h := newTCloudMachineConfigTestHarness(t)
	d := h.newResourceData(t, h.sharedValues())

	require.NoError(t, resourceRancher2MachineConfigV2Create(d, h.config))
	require.Equal(t, h.id, d.Id())

	// First refresh after create (create ends with a read, but call again
	// to exercise the read path in isolation).
	require.NoError(t, resourceRancher2MachineConfigV2Read(d, h.config))
	first := map[string]interface{}{
		"vpc_id":          d.Get("tcloud_public_config.0.vpc_id"),
		"subnet_id":       d.Get("tcloud_public_config.0.subnet_id"),
		"sec_groups":      d.Get("tcloud_public_config.0.sec_groups"),
		"image_name":      d.Get("tcloud_public_config.0.image_name"),
		"network_scope":   d.Get("tcloud_public_config.0.network_scope"),
		"skip_default_sg": d.Get("tcloud_public_config.0.skip_default_sg"),
		"name":            d.Get("name"),
		"kind":            d.Get("kind"),
	}

	// Second refresh, server state unchanged.
	require.NoError(t, resourceRancher2MachineConfigV2Read(d, h.config))
	second := map[string]interface{}{
		"vpc_id":          d.Get("tcloud_public_config.0.vpc_id"),
		"subnet_id":       d.Get("tcloud_public_config.0.subnet_id"),
		"sec_groups":      d.Get("tcloud_public_config.0.sec_groups"),
		"image_name":      d.Get("tcloud_public_config.0.image_name"),
		"network_scope":   d.Get("tcloud_public_config.0.network_scope"),
		"skip_default_sg": d.Get("tcloud_public_config.0.skip_default_sg"),
		"name":            d.Get("name"),
		"kind":            d.Get("kind"),
	}

	for key, a := range first {
		require.Equal(t, a, second[key], "refresh must not drift field %q", key)
	}
	// And the state still reflects what we created.
	require.Equal(t, "test-vpc", first["vpc_id"])
	require.Equal(t, "test-subnet", first["subnet_id"])
	require.Equal(t, "test-sg", first["sec_groups"])
	require.Equal(t, "test-image", first["image_name"])
	require.Equal(t, "shared", first["network_scope"])
	require.Equal(t, true, first["skip_default_sg"])
}

// TestTCloudPublicMachineConfigLifecycleDeleteIdempotent covers the TC-10
// "delete idempotency / CR disappearance" requirement. A delete against an
// already-deleted resource must succeed silently (Terraform destroy is
// retried, and the CR can also be removed by the controller when the
// cluster is deleted — the provider must not error on that).
func TestTCloudPublicMachineConfigLifecycleDeleteIdempotent(t *testing.T) {
	h := newTCloudMachineConfigTestHarness(t)
	d := h.newResourceData(t, h.sharedValues())
	require.NoError(t, resourceRancher2MachineConfigV2Create(d, h.config))
	require.Equal(t, h.id, d.Id())

	require.NoError(t, resourceRancher2MachineConfigV2Delete(d, h.config))
	require.Empty(t, d.Id())

	// Simulate a second delete attempt (e.g. user re-runs destroy, or the
	// controller already removed the CR when the cluster went away).
	d.SetId(h.id)
	require.NoError(t, resourceRancher2MachineConfigV2Delete(d, h.config))
	require.Empty(t, d.Id())

	// A read after the resource is gone must also clear the state rather
	// than error (Terraform refresh on a removed resource).
	d.SetId(h.id)
	require.NoError(t, resourceRancher2MachineConfigV2Read(d, h.config))
	require.Empty(t, d.Id())
}

// TestTCloudPublicClusterV2AnnotationRoundTrip covers the TC-10 "annotation
// contract" requirement. The managed-network flow relies on the cluster
// carrying the network-controller annotations (ui.rancher/provider,
// infrastructure.otc.t-systems.com/cluster-network,
// infrastructure.otc.t-systems.com/network-policy). These must round-trip
// through expand/flatten without loss so the controller can bind.
func TestTCloudPublicClusterV2AnnotationRoundTrip(t *testing.T) {
	annotations := map[string]interface{}{
		"ui.rancher/provider":                              "opentelekomcloud",
		"infrastructure.otc.t-systems.com/cluster-network": "tf-tcloud-managed",
		"infrastructure.otc.t-systems.com/network-policy":  "Managed",
	}
	d := schema.TestResourceDataRaw(t, clusterV2Fields(), map[string]interface{}{
		"name": "tcloud-cluster", "fleet_namespace": "fleet-default",
		"kubernetes_version": "v1.36.4+rke2r1", "annotations": annotations,
	})
	obj, err := expandClusterV2(d)
	require.NoError(t, err)
	require.Equal(t, "opentelekomcloud", obj.ObjectMeta.Annotations["ui.rancher/provider"])
	require.Equal(t, "tf-tcloud-managed", obj.ObjectMeta.Annotations["infrastructure.otc.t-systems.com/cluster-network"])
	require.Equal(t, "Managed", obj.ObjectMeta.Annotations["infrastructure.otc.t-systems.com/network-policy"])

	// Flatten back into a fresh resource data and verify the annotations
	// survive the round-trip.
	d2 := schema.TestResourceDataRaw(t, clusterV2Fields(), nil)
	require.NoError(t, flattenClusterV2(d2, obj))
	got := d2.Get("annotations").(map[string]interface{})
	require.Equal(t, annotations["ui.rancher/provider"], got["ui.rancher/provider"])
	require.Equal(t, annotations["infrastructure.otc.t-systems.com/cluster-network"], got["infrastructure.otc.t-systems.com/cluster-network"])
	require.Equal(t, annotations["infrastructure.otc.t-systems.com/network-policy"], got["infrastructure.otc.t-systems.com/network-policy"])
}
