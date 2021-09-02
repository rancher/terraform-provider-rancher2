package rancher2

import (
	"reflect"
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

var (
	testClusterV2RKEConfigETCDSnapshotS3Conf      *rkev1.ETCDSnapshotS3
	testClusterV2RKEConfigETCDSnapshotS3Interface []interface{}
	testClusterV2RKEConfigETCDConf                *rkev1.ETCD
	testClusterV2RKEConfigETCDInterface           []interface{}
)

func init() {
	testClusterV2RKEConfigETCDSnapshotS3Conf = &rkev1.ETCDSnapshotS3{
		Bucket:              "bucket",
		CloudCredentialName: "cloud_credential_name",
		Endpoint:            "endpoint",
		EndpointCA:          "endpoint_ca",
		Folder:              "folder",
		Region:              "region",
		SkipSSLVerify:       true,
	}

	testClusterV2RKEConfigETCDSnapshotS3Interface = []interface{}{
		map[string]interface{}{
			"bucket":                "bucket",
			"cloud_credential_name": "cloud_credential_name",
			"endpoint":              "endpoint",
			"endpoint_ca":           "endpoint_ca",
			"folder":                "folder",
			"region":                "region",
			"skip_ssl_verify":       true,
		},
	}
	testClusterV2RKEConfigETCDConf = &rkev1.ETCD{
		DisableSnapshots:     true,
		SnapshotScheduleCron: "snapshot_schedule_cron",
		SnapshotRetention:    10,
		S3:                   testClusterV2RKEConfigETCDSnapshotS3Conf,
	}

	testClusterV2RKEConfigETCDInterface = []interface{}{
		map[string]interface{}{
			"disable_snapshots":      true,
			"snapshot_schedule_cron": "snapshot_schedule_cron",
			"snapshot_retention":     10,
			"s3_config":              testClusterV2RKEConfigETCDSnapshotS3Interface,
		},
	}
}

func TestFlattenClusterV2RKEConfigETCDSnapshotS3(t *testing.T) {

	cases := []struct {
		Input          *rkev1.ETCDSnapshotS3
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigETCDSnapshotS3Conf,
			testClusterV2RKEConfigETCDSnapshotS3Interface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigETCDSnapshotS3(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterV2RKEConfigETCD(t *testing.T) {

	cases := []struct {
		Input          *rkev1.ETCD
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigETCDConf,
			testClusterV2RKEConfigETCDInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigETCD(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigETCDSnapshotS3(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rkev1.ETCDSnapshotS3
	}{
		{
			testClusterV2RKEConfigETCDSnapshotS3Interface,
			testClusterV2RKEConfigETCDSnapshotS3Conf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigETCDSnapshotS3(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigETCD(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rkev1.ETCD
	}{
		{
			testClusterV2RKEConfigETCDInterface,
			testClusterV2RKEConfigETCDConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigETCD(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
