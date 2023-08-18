# Rancher2 Provider

### Version compatibility matrix

The version matrix specifies the Terraform provider version _recommended_ to use with the associated minor Rancher version that it was released for. When updating the version matrix, add a row for each Terraform version released with a minor Rancher release.

#### Rancher 2.6

| Terraform provider version | Rancher |    Notes    |
|----------------------------------------|:-------:|:-----------:|
| 2.0.0                                  | 2.6.11  | Bug fixes   |

#### Rancher 2.7

| Terraform provider version | Rancher | Notes                                                                                               |
|----------------------------|:-------:|-----------------------------------------------------------------------------------------------------|
| 3.0.0                      |  2.7.2  | Kubernetes 1.25 support, Azure / EKS / Harvester features<br/>and bug fixes                         |
| 3.0.1                      |  2.7.4  | Fix to support old Harvester config                                                                 |
| 3.0.2                      |  2.7.4  | Fix Harvester disk_size default value                                                               |
| 3.1.0                      |  2.7.5  | Cluster Agent customization, PSACT support for 1.25+ clusters,<br/>custom user tokens and bug fixes |
| 3.1.1                      |  2.7.5  | Docs patch                                                                                          |
| 3.2.0                      |  2.7.x  |                                                                                                     |

#### Rancher 2.8

| Terraform provider version | Rancher | Notes |
|----------------------------|:-------:|-------|
| 4.0.0                      |  2.8.x  |       |                                                                                                  |

#### FAQ

**Can I use an earlier Terraform version?** Yes, but Terraform may not support all features and fields supported in your Rancher version so provisioning may be limited.

**Can I use a later Terraform version?** Yes, but you must NOT use any new features and fields that your Rancher version does not support.
