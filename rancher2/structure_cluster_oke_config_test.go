// Copyright 2020 Oracle and/or its affiliates.

package rancher2

import (
	"reflect"
	"testing"
)

var (
	testClusterOKEConfigConf      *OracleKubernetesEngineConfig
	testClusterOKEConfigInterface []interface{}
)

func init() {
	testClusterOKEConfigConf = &OracleKubernetesEngineConfig{
		CompartmentID:               "compartment",
		CustomBootVolumeSize:        0,
		Description:                 "description",
		DisplayName:                 "test",
		DriverName:                  clusterDriverOKE,
		EnableKubernetesDashboard:   true,
		Fingerprint:                 "fingerprint",
		FlexOCPUs:                   0,
		KubernetesVersion:           "version",
		KMSKeyID:                    "ocid1.key.oc1.reg.xxxxxxxxxxxxxxxxxx",
		LimitNodeCount:              0,
		Name:                        "test",
		NodeImage:                   "image",
		NodePoolSubnetDNSDomainName: "nodedns",
		NodePoolSubnetName:          "nodedns",
		NodePublicSSHKeyContents:    "public_key",
		NodeShape:                   "shape",
		PodCidr:                     "10.204.0.0/16",
		PrivateKeyContents:          "private_key_contents",
		PrivateKeyPassphrase:        "",
		PrivateNodes:                false,
		QuantityOfSubnets:           2,
		QuantityPerSubnet:           1,
		Region:                      "region",
		ServiceCidr:                 "10.98.0.0/16",
		ServiceLBSubnet1Name:        "",
		ServiceLBSubnet2Name:        "",
		ServiceSubnetDNSDomainName:  "svcdns",
		SkipVCNDelete:               false,
		TenancyID:                   "tenancy",
		UserOCID:                    "user",
		VCNName:                     "",
		VcnCompartmentID:            "",
		WorkerNodeIngressCidr:       "",
	}
	testClusterOKEConfigInterface = []interface{}{
		map[string]interface{}{
			"compartment_id":               "compartment",
			"custom_boot_volume_size":      0,
			"description":                  "description",
			"enable_private_control_plane": false,
			"enable_kubernetes_dashboard":  true,
			"enable_private_nodes":         false,
			"fingerprint":                  "fingerprint",
			"flex_ocpus":                   0,
			"kms_key_id":                   "ocid1.key.oc1.phx.xxxxxxxxxxxxxxxxxx",
			"kubernetes_version":           "version",
			"limit_node_count":             0,
			"node_image":                   "image",
			"node_pool_dns_domain_name":    "nodedns",
			"node_pool_subnet_name":        "nodedns",
			"node_public_key_contents":     "public_key",
			"node_shape":                   "shape",
			"pod_cidr":                     "10.204.0.0/16",
			"private_key_contents":         "private_key_contents",
			"private_key_passphrase":       "",
			"quantity_of_node_subnets":     2,
			"quantity_per_subnet":          1,
			"region":                       "region",
			"load_balancer_subnet_name_1":  "",
			"load_balancer_subnet_name_2":  "",
			"service_cidr":                 "10.98.0.0/16",
			"service_dns_domain_name":      "svcdns",
			"skip_vcn_delete":              false,
			"tenancy_id":                   "tenancy",
			"user_ocid":                    "user",
			"vcn_compartment_id":           "",
			"vcn_name":                     "",
			"worker_node_ingress_cidr":     "",
		},
	}
}

func TestFlattenClusterOKEConfig(t *testing.T) {

	cases := []struct {
		Input          *OracleKubernetesEngineConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterOKEConfigConf,
			testClusterOKEConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterOKEConfig(tc.Input, testClusterOKEConfigInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterOKEConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *OracleKubernetesEngineConfig
	}{
		{
			testClusterOKEConfigInterface,
			testClusterOKEConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterOKEConfig(tc.Input, "test")
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
