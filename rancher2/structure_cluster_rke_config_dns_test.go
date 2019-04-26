package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterRKEConfigDNSConf      *managementClient.DNSConfig
	testClusterRKEConfigDNSInterface []interface{}
)

func init() {
	testClusterRKEConfigDNSConf = &managementClient.DNSConfig{
		NodeSelector: map[string]string{
			"sel1": "value1",
			"sel2": "value2",
		},
		Provider:            "kube-dns",
		ReverseCIDRs:        []string{"rev1", "rev2"},
		UpstreamNameservers: []string{"up1", "up2"},
	}
	testClusterRKEConfigDNSInterface = []interface{}{
		map[string]interface{}{
			"node_selector": map[string]interface{}{
				"sel1": "value1",
				"sel2": "value2",
			},
			"provider":             "kube-dns",
			"reverse_cidrs":        []interface{}{"rev1", "rev2"},
			"upstream_nameservers": []interface{}{"up1", "up2"},
		},
	}
}

func TestFlattenClusterRKEConfigDNS(t *testing.T) {

	cases := []struct {
		Input          *managementClient.DNSConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigDNSConf,
			testClusterRKEConfigDNSInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigDNS(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigDNS(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.DNSConfig
	}{
		{
			testClusterRKEConfigDNSInterface,
			testClusterRKEConfigDNSConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigDNS(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
