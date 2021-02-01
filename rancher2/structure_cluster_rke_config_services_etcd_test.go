package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigServicesETCDBackupS3Conf      *managementClient.S3BackupConfig
	testClusterRKEConfigServicesETCDBackupS3Interface []interface{}
	testClusterRKEConfigServicesETCDBackupConf        *managementClient.BackupConfig
	testClusterRKEConfigServicesETCDBackupInterface   []interface{}
	testClusterRKEConfigServicesETCDConf              *managementClient.ETCDService
	testClusterRKEConfigServicesETCDInterface         []interface{}
)

func init() {
	testClusterRKEConfigServicesETCDBackupS3Conf = &managementClient.S3BackupConfig{
		AccessKey:  "access_key",
		BucketName: "bucket_name",
		CustomCA:   "custom_ca",
		Endpoint:   "endpoint",
		Folder:     "folder",
		Region:     "region",
	}
	testClusterRKEConfigServicesETCDBackupS3Interface = []interface{}{
		map[string]interface{}{
			"access_key":  "access_key",
			"bucket_name": "bucket_name",
			"custom_ca":   Base64Encode("custom_ca"),
			"endpoint":    "endpoint",
			"folder":      "folder",
			"region":      "region",
		},
	}
	testClusterRKEConfigServicesETCDBackupConf = &managementClient.BackupConfig{
		Enabled:        newTrue(),
		IntervalHours:  20,
		Retention:      10,
		S3BackupConfig: testClusterRKEConfigServicesETCDBackupS3Conf,
		SafeTimestamp:  true,
		Timeout:        500,
	}
	testClusterRKEConfigServicesETCDBackupInterface = []interface{}{
		map[string]interface{}{
			"enabled":          true,
			"interval_hours":   20,
			"retention":        10,
			"s3_backup_config": testClusterRKEConfigServicesETCDBackupS3Interface,
			"safe_timestamp":   true,
			"timeout":          500,
		},
	}
	testClusterRKEConfigServicesETCDConf = &managementClient.ETCDService{
		BackupConfig: testClusterRKEConfigServicesETCDBackupConf,
		CACert:       "XXXXXXXX",
		Cert:         "YYYYYYYY",
		Creation:     "creation",
		ExternalURLs: []string{"url_one", "url_two"},
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds: []string{"bind_one", "bind_two"},
		ExtraEnv:   []string{"env_one", "env_two"},
		GID:        int64(1001),
		Image:      "image",
		Key:        "ZZZZZZZZ",
		Path:       "/etcd",
		Retention:  "6h",
		Snapshot:   newTrue(),
		UID:        int64(1001),
	}
	testClusterRKEConfigServicesETCDInterface = []interface{}{
		map[string]interface{}{
			"backup_config": testClusterRKEConfigServicesETCDBackupInterface,
			"ca_cert":       "XXXXXXXX",
			"cert":          "YYYYYYYY",
			"creation":      "creation",
			"external_urls": []interface{}{"url_one", "url_two"},
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds": []interface{}{"bind_one", "bind_two"},
			"extra_env":   []interface{}{"env_one", "env_two"},
			"gid":         1001,
			"image":       "image",
			"key":         "ZZZZZZZZ",
			"path":        "/etcd",
			"retention":   "6h",
			"snapshot":    true,
			"uid":         1001,
		},
	}
}

func TestFlattenClusterRKEConfigServicesEtcdBackupConfigS3(t *testing.T) {

	cases := []struct {
		Input          *managementClient.S3BackupConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesETCDBackupS3Conf,
			testClusterRKEConfigServicesETCDBackupS3Interface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKEConfigServicesEtcdBackupConfigS3(tc.Input, testClusterRKEConfigServicesETCDBackupS3Interface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigServicesEtcdBackupConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.BackupConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesETCDBackupConf,
			testClusterRKEConfigServicesETCDBackupInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKEConfigServicesEtcdBackupConfig(tc.Input, testClusterRKEConfigServicesETCDBackupInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigServicesEtcd(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ETCDService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesETCDConf,
			testClusterRKEConfigServicesETCDInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesEtcd(tc.Input, testClusterRKEConfigServicesETCDInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesEtcdBackupConfigS3(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.S3BackupConfig
	}{
		{
			testClusterRKEConfigServicesETCDBackupS3Interface,
			testClusterRKEConfigServicesETCDBackupS3Conf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesEtcdBackupConfigS3(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesEtcdBackupConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.BackupConfig
	}{
		{
			testClusterRKEConfigServicesETCDBackupInterface,
			testClusterRKEConfigServicesETCDBackupConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesEtcdBackupConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesEtcd(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ETCDService
	}{
		{
			testClusterRKEConfigServicesETCDInterface,
			testClusterRKEConfigServicesETCDConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesEtcd(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
