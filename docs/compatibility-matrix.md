# Rancher2 Provider

### Version compatibility matrix

The version matrix specifies the Terraform provider version _recommended_ to use with the associated minor Rancher version that it was released for. When updating the version matrix, add a row for each Terraform version released with a minor Rancher release.

#### Rancher 2.6

| Terraform provider version | Rancher |
|----------------------------------------|:-------:|
| 2.0.0                                  | 2.6.11  |
| 2.0.1 (2.1.0 for features)             | 2.6.x   |

#### Rancher 2.7

| Terraform provider version | Rancher |
|----------------------------------------|:-------:|
| 3.0.0                                  | 2.7.2   |
| 3.0.1 (3.1.0 for features)             | 2.7.x   |

**Can I use an earlier Terraform version?** Yes, but Terraform may not support all features and fields supported in your Rancher version so provisioning may be limited.

**Can I use a later Terraform version?** Yes, but you must NOT use any new features and fields that your Rancher version does not support.
