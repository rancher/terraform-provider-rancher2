package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testEtcdBackupConfigS3Conf      *managementClient.S3BackupConfig
	testEtcdBackupConfigS3Interface []interface{}
	testEtcdBackupConfigConf        *managementClient.BackupConfig
	testEtcdBackupConfigInterface   []interface{}
	testEtcdBackupConf              *managementClient.EtcdBackup
	testEtcdBackupInterface         map[string]interface{}
)

func init() {
	testEtcdBackupConfigS3Conf = &managementClient.S3BackupConfig{
		AccessKey:  "access_key",
		BucketName: "bucket_name",
		CustomCA:   "custom_ca",
		Endpoint:   "endpoint",
		Folder:     "folder",
		Region:     "region",
		SecretKey:  "secret",
	}
	testEtcdBackupConfigS3Interface = []interface{}{
		map[string]interface{}{
			"access_key":  "access_key",
			"bucket_name": "bucket_name",
			"custom_ca":   Base64Encode("custom_ca"),
			"endpoint":    "endpoint",
			"folder":      "folder",
			"region":      "region",
			"secret_key":  "secret",
		},
	}
	testEtcdBackupConfigConf = &managementClient.BackupConfig{
		Enabled:        newTrue(),
		IntervalHours:  20,
		Retention:      10,
		S3BackupConfig: testEtcdBackupConfigS3Conf,
		SafeTimestamp:  false,
		Timeout:        500,
	}
	testEtcdBackupConfigInterface = []interface{}{
		map[string]interface{}{
			"enabled":          true,
			"interval_hours":   20,
			"retention":        10,
			"s3_backup_config": testEtcdBackupConfigS3Interface,
			"safe_timestamp":   false,
			"timeout":          500,
		},
	}
	testEtcdBackupConf = &managementClient.EtcdBackup{
		BackupConfig: testEtcdBackupConfigConf,
		ClusterID:    "cluster-test",
		Filename:     "filename",
		Manual:       true,
		Name:         "test",
		NamespaceId:  "namespace_id",
	}
	testEtcdBackupInterface = map[string]interface{}{
		"backup_config": testEtcdBackupConfigInterface,
		"cluster_id":    "cluster-test",
		"filename":      "filename",
		"manual":        true,
		"name":          "test",
		"namespace_id":  "namespace_id",
	}
}

func TestFlattenEtcdBackup(t *testing.T) {

	cases := []struct {
		Input          *managementClient.EtcdBackup
		ExpectedOutput map[string]interface{}
	}{
		{
			testEtcdBackupConf,
			testEtcdBackupInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, etcdBackupFields(), map[string]interface{}{})
		err := flattenEtcdBackup(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandEtcdBackup(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.EtcdBackup
	}{
		{
			testEtcdBackupInterface,
			testEtcdBackupConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, etcdBackupFields(), tc.Input)
		output, err := expandEtcdBackup(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expnader: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
