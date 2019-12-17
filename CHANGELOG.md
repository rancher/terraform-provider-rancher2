## 1.7.2 (Unreleased)

FEATURES:



ENHANCEMENTS:

* Added `refresh` argument to `rancher2_catalog` resource
* Added `name` and `is_external` arguments to `rancher2_user` datasource
* Added `delete_not_ready_after_secs` and `node_taints` arguments to `node_pool` resource
* Added `delete_not_ready_after_secs` and `node_taints` arguments to `rancher2_node_pool` resource
* Updated `github.com/rancher/types` and `github.com/rancher/norman` go modules and vendor files to support rancher v2.3.3
* Splitted schema, structure and test `cluster_rke_config_services` files for every rke service 
* Added `ssh_cert_path` argument to `rke_config` argument on `rancher2_cluster` resource
* Added `audit_log`, `event_rate_limit` and `secrets_encryption_config` arguments to `rke_config.services.kube_api` argument on `rancher2_cluster` resource
* Added `generate_serving_certificate` argument to `rke_config.services.kubelet` argument on `rancher2_cluster` resource

BUG FIXES:



## 1.7.1 (December 04, 2019)

FEATURES:



ENHANCEMENTS:

* Added GetRancherVersion function to provider config
* Updated `vsphere_config` argument schema on `rancher2_node_template` resource to support Rancher v2.3.3 features
* Updated rancher to v2.3.3 and k3s to v0.10.2 on acceptance tests

BUG FIXES:

* Set `annotations` argument as computed on `rancher2_node_template` resource
* Added `rancher2_node_template` resource workaround on docs when upgrade Rancher to v2.3.3

## 1.7.0 (November 20, 2019)

FEATURES:

* **New Resource:** `rancher2_token`

ENHANCEMENTS:

* Added `always_pull_images` argument on `kube_api` argument on `rke_config` argument for `rancher2_clusters` resource
* Added resource deletion if not getting active state on creation for `rancher2_catalog` resource
* Updated rancher to v2.3.2 and k3s to v0.10.1 on acceptance tests
* Added `desired nodes` support on `eks_config` argument on `rancher2_cluster` resource
* Added `managed disk` support on `azure_config` argument on `rancher2_node_template` resource
* Migrated provider to use `terraform-plugin-sdk`
* Updated `rancher2_etcd_backup` documentation

BUG FIXES:

* Fix `password` argument update for `rancher2_catalog` resource
* Fix `rancher2_app` update issue on Rancher v2.3.2
* Fix: set `key` argument as sensitive on `rancher2_certificate` resource. 
* Fix continuous diff issues on `rancher2_project` resource
* Fix `pod_security_policy_template_id` update on `rancher2_project` resource
* Fix continuous diff issues on `rancher2_namespace` resource

## 1.6.0 (October 08, 2019)

FEATURES:

* **New Data Source:** `rancher2_cluster_alert_group`
* **New Data Source:** `rancher2_cluster_alert_rule`
* **New Data Source:** `rancher2_cluster_template`
* **New Data Source:** `rancher2_notifier`
* **New Data Source:** `rancher2_project_alert_group`
* **New Data Source:** `rancher2_project_alert_rule`
* **New Data Source:** `rancher2_role_template`
* **New Resource:** `rancher2_auth_config_keycloak`
* **New Resource:** `rancher2_auth_config_okta`
* **New Resource:** `rancher2_cluster_alert_group`
* **New Resource:** `rancher2_cluster_alert_rule`
* **New Resource:** `rancher2_cluster_sync`
* **New Resource:** `rancher2_cluster_template`
* **New Resource:** `rancher2_notifier`
* **New Resource:** `rancher2_project_alert_group`
* **New Resource:** `rancher2_project_alert_rule`
* **New Resource:** `rancher2_role_template`

ENHANCEMENTS:

* Added `monitoring_input` argument to define monitoring config for `rancher2_cluster` and `rancher2_project`
* Improved capitalization/spelling/grammar/etc in docs

BUG FIXES:

* Fix `expandAppExternalID` function on `rancher2_app` resource. Function was generating a wrong `ExternalID` catalog URL, on `cluster` and `project` scope
* Fix `flattenMultiClusterApp` function on `rancher2_multi-cluster_app` resource. Function wasn't updating fine `catalog_name`, `template_name` and/or `template_version` arguments, when contains char `-`
* Fix: set `value_yaml` multiline argument as base64 encoded
* Fix: removed `restricted` and `unrestricted` values checking for `default_pod_security_policy_template_id` argument on `rancher2_cluster` resource

## 1.5.0 (September 06, 2019)

FEATURES:

* **New Data Source:** `rancher2_app`
* **New Data Source:** `rancher2_certificate`
* **New Data Source:** `rancher2_multi_cluster_app`
* **New Data Source:** `rancher2_node_template`
* **New Data Source:** `rancher2_secret`
* **New Resource:** `rancher2_certificate`
* **New Resource:** `rancher2_app`
* **New Resource:** `rancher2_multi_cluster_app`
* **New Resource:** `rancher2_secret`

ENHANCEMENTS:

* Updated default image to `canonical:UbuntuServer:18.04-LTS:latest` on Azure node template
* Added `folder` argument on `s3_backup_config`
* Updated `github.com/rancher/types` and `github.com/rancher/norman` go modules and vendor files to support rancher v2.2.8
* Updated rancher to v2.2.8 and k3s to v0.8.0 on acceptance tests
* Added `key_pair_name` argument on `eks_config` argument on `rancher2_cluster` resource
* Set `kubernetes_version` argument as required on `eks_config` argument on `rancher2_cluster` resource
* Set `quantity` argument as optional with default value `1` on `rancher2_node_pool` resource. Added validation that value >= 1 

BUG FIXES:

* Fix: `container_resource_limit` argument update issue on `rancher2_namespace` and `rancher2_project` resources update
* Fix: `sidebar_current` definition on datasources docs
* Fix: set `access_key` and `secret_key` arguments as optional on `s3_backup_config`
* Fix: crash `rancher2_cluster`  datasource and resource if `enableNetworkPolicy` doesn't exist
* Fix: don't delete builtin cluster nor node drivers from rancher on tf destroy
* Fix: wrong updates on not changed sensitive arguments on `rancher2_cluster_logging` and `rancher2_project_logging` resources

## 1.4.1 (August 16, 2019)

FEATURES:

ENHANCEMENTS:

BUG FIXES:

* Fix: auth issue when using `access_key` and `secret_key`

## 1.4.0 (August 15, 2019)

FEATURES:

* **New Data Source:** `rancher2_catalog`
* **New Data Source:** `rancher2_cloud_credential`
* **New Data Source:** `rancher2_cluster`
* **New Data Source:** `rancher2_cluster_driver`
* **New Data Source:** `rancher2_cluster_logging`
* **New Data Source:** `rancher2_cluster_role_template_binding`
* **New Data Source:** `rancher2_etcd_backup`
* **New Data Source:** `rancher2_global_role_binding`
* **New Data Source:** `rancher2_namespace`
* **New Data Source:** `rancher2_node_driver`
* **New Data Source:** `rancher2_node_pool`
* **New Data Source:** `rancher2_project_logging`
* **New Data Source:** `rancher2_project_role_template_binding`
* **New Data Source:** `rancher2_registry`
* **New Data Source:** `rancher2_user`
* **New Resource:** `rancher2_global_role_binding`
* **New Resource:** `rancher2_registry`
* **New Resource:** `rancher2_user`

ENHANCEMENTS:

* Set `session_token` argument as sensitive on `eks_config` argument on `rancher2_cluster` resource
* Added `wait_for_cluster` argument on `rancher2_namespace` and `rancher2_project` resources
* Set default value to `engine_install_url` argument on `rancher2_node_template` resource
* Added `enable_cluster_monitoring` argument to `rancher2_cluster` resource and datasource
* Added `enable_project_monitoring` argument to `rancher2_project` resource and datasource
* Added `token` argument on `cluster_registration_token` argument to rancher2_cluster resource and datasource
* Set default value to `engine_install_url` argument on `rancher2_node_template` resource
* Added `custom_ca` argument on etcd `s3_backup_config` on `rancher2_cluster` and `rancher2_etcd_backup` resources
* Updated `github.com/rancher/types` and `github.com/rancher/norman` go modules and vendor files to support rancher v2.2.6
* Updated rancher to v2.2.6 and k3s to v0.7.0 on acceptance tests
* Added cluster and project scope support on `rancher2_catalog` resource and datasource
* Updated `provider` config validation to enable bootstrap and resource creation at same run
* Added `container_resource_limit` argument on `rancher2_namespace` and `rancher2_project` resources and datasources
* Added `pod_security_policy_template_id` on `rancher2_project` resource

BUG FIXES:

* Fix: `toArrayString` and `toMapString` functions to check `nil` values
* Fix: Set `kubernetes_version` argument as required on `aks_config` argument on `rancher2_cluster` resource
* Fix: Set `security_groups`, `service_role`, `subnets` and `virtual_network` arguments as optional to `eks_config` argument on `rancher2_cluster` resource
* Fix: Removed `docker_version` argument from `rancher2_node_template` resource

## 1.3.0 (June 26, 2019)

FEATURES:

ENHANCEMENTS:

* Added `scheduler` argument to `services`-`rke_config` argument on `rancher2_cluster` resource

BUG FIXES:

* Fix: index out of range issue on `vsphere_cloud_provider`-`cloud_provider`-`rke_config` argument on `rancher2_cluster` resource 

## 1.2.0 (June 12, 2019)

FEATURES:

* **New Data Source:** `rancher2_project`

ENHANCEMENTS:

* Added `cluster_auth_endpoint` argument to `rancher2_cluster` resource
* Added `default_pod_security_policy_template_id` argument to `rancher2_cluster` resource
* Added `enable_network_policy` argument to `rancher2_cluster` resource
* Updated acceptance tests
  * k3s version updated to v0.5.0
  * Rancher version updated to v2.2.4

BUG FIXES:

* Fix: set default value to `true` on `ignore_docker_version`-`rke_config` argument on `rancher2_cluster` resource
* Fix: set default value to `false` on `pod_security_policy`-`services`-`rke_config` argument on `rancher2_cluster` resource
* Fix: typo on `boot2docker_url`-`vsphere_config` argument name on `rancher2_node_template` resource docs
* Fix: set `monitor_delay` and `monitor_timeout` fields as string type for openstack load_balancer config on `cloud_provider`-`rke_config` argument on `rancher2_cluster` resource
* Fix: Updated `rancher2_etcd_backup` resource to work on rancher v2.2.4

## 1.1.0 (May 29, 2019)

FEATURES:

ENHANCEMENTS:

* Added `default_project_id` & `system_project_id` attributes to `rancher2_cluster` resource
* Added support to move `rancher2_namespace` resource to a rancher project when import
* Added support to terraform 0.12

BUG FIXES:

* Fix: Updated `flattenNamespace` function on `rancher2_namespace` resource to avoid no empty plan if `resource_quota` is not specified
* Fix: Updated `rke_config` argument for openstack cloud_provider on `rancher2_cluster` resource:
  * Removed `used_id` field on global argument in favour of `username` following [k8s openstack cloud provider docs](https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/provider-configuration.md#global-required-parameters)
  * Set computed=true on optional field to avoid no empty plan if not specified

## 1.0.0 (May 14, 2019)

* Initial Terraform Ecosystem Release


## v0.2.0-rc5 (Unreleased)

FEATURES:

ENHANCEMENTS:

* Updated `rancher2_cluster` `rke_config` argument to support `aws_cloud_provider` config
* Updated k3s version to v0.4.0 to run acceptance tests
* Added support to openstack and vsphere drivers on `rancher2_cloud_credential` resource
* Added support to openstack and vsphere drivers on `rancher2_node_template` resource

BUG FIXES:

* Fix: Updated `rancher2_cluster` resource to save correctly S3 and cloud providers passwords on `rke_config`
* Updated `rancher2_cloud_credential` resource to save correctly S3 password
* Updated `rancher2_etcd_backup` resource to save correctly S3 password

## v0.2.0-rc4 (Unreleased)

FEATURES:

* **New Resource:** `rancher2_bootstrap`
* **New Resource:** `rancher2_cloud_credential`
* **New Resource:** `rancher2_cluster_driver`
* **New Resource:** `rancher2_etcd_backup`

ENHANCEMENTS:

* Added `.drone.yml` file to also support run rancher pipeline
* Added `rancher2_node_pool` resource tests
* Added `rancher2_auth_config_*` resource tests
* Updated and reviewed docs format
* Added support to rancher v2.2.x
* Updated `rancher2_cluster` `rke_config` argument to support:
  * etcd service `backup_config` with local and S3 storage backends
  * `dns` config
  * `weave` network provider
* Splitted resources into own schema, structure and import files.
* Added support to amazonec2, azure and digitalocean drivers on `rancher2_cloud_credential` resource
* Added support to local and S3 storage backends on `rancher2_etcd_backup` resource

BUG FIXES:

* Fix: drone build image to golang:1.12.3 to fix go fmt issues
* Fix: removed test on apply for `rancher2_auth_config_*` resources
* Fix: updated `api_url` field as required on provider.go
* Fix: updated `rancher2_namespace` move to a project after import it from k8s cluster

## v0.2.0-rc3 (Unreleased)

FEATURES:

ENHANCEMENTS:

* Added `Sensitive: true` option to fields with sensible data

BUG FIXES:

* Fix: set rke cluster `cloud_provider_vsphere` disk and network as optional and computed fields

## v0.2.0-rc2 (Unreleased)

FEATURES:

ENHANCEMENTS:

* Added `Sensitive: true` option to fields with sensible data
* Added `kube_config` computed field on cluster resources
* Added `ami` and `associate_worker_node_public_ip` fields for `eks_config` on cluster resources
* Added all available fields for rke_config on cluster resources
* Added `manifest_url` and `windows_node_command` fields for `cluster_registration_token` on cluster resources
* Added `creation` argument on `etcd` service for rke_config on cluster resources

BUG FIXES:

* Fix: added updating pending state on cluster resource update
* Fix: checking if `cluster_registration_token` exists on cluster resource creation
* Fix: typo on `gke_config` credential field on cluster resource
* Fix: Updated auth resources to avoid permission denied error

## 0.1.0-rc1 (Unreleased)

FEATURES:

* **New Resource:** `rancher2_auth_config_activedirectory`
* **New Resource:** `rancher2_auth_config_adfs`
* **New Resource:** `rancher2_auth_config_azuread`
* **New Resource:** `rancher2_auth_config_freeipa`
* **New Resource:** `rancher2_auth_config_github`
* **New Resource:** `rancher2_auth_config_openldap`
* **New Resource:** `rancher2_auth_config_ping`
* **New Resource:** `rancher2_catalog`
* **New Resource:** `rancher2_cluster`
* **New Resource:** `rancher2_cluster_logging`
* **New Resource:** `rancher2_cluster_role_template_binding`
* **New Resource:** `rancher2_namespace`
* **New Resource:** `rancher2_node_driver`
* **New Resource:** `rancher2_node_pool`
* **New Resource:** `rancher2_node_template`
* **New Resource:** `rancher2_project`
* **New Resource:** `rancher2_project_logging`
* **New Resource:** `rancher2_project_role_template_binding`
* **New Resource:** `rancher2_setting`

ENHANCEMENTS:

* First release candidate of the rancher2 provider.
* resource/rancher2_cluster: support for providers:
  * Amazon EKS
  * Azure AKS
  * Google GKE
  * Imported
  * RKE
    * Cloud providers adding node pools
    * Custom
* resource/rancher2_cluster_logging: support for providers:
  * Elasticsearch
  * Fluentd
  * Kafka
  * Splunk
  * Syslog
* resource/rancher2_namespace: quota limits support on Rancher v2.1.x or higher
  * Amazon EC2
  * Azure
  * Digitalocean
* resource/rancher2_project: quota limits support on Rancher v2.1.x or higher
* resource/rancher2_project_logging: support for providers:
  * Elasticsearch
  * Fluentd
  * Kafka
  * Splunk
  * Syslog
* resource/rancher2_node_template: support for providers:

BUG FIXES:
