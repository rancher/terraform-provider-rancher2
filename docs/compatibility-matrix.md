# Rancher2 Provider

## Version compatibility matrix

The version matrix specifies the Terraform provider version _recommended_ to use with the associated minor Rancher version that it was released for.
When updating the version matrix, add a row for each Terraform version released with a minor Rancher release.
This shouldn't be a changelog for every branch, each branch has its own CHANGELOG.md file which should be updated with changes.

| Terraform provider version | Rancher version | Terraform provider branch |
|----------------------------|:---------------:|---------------------------|
| 2.x                        | 2.6.x           | release/v2 |
| 3.x                        | 2.7.x           | release/v3 |
| 4.x                        | 2.8.x           | release/v4 |
| 5.x                        | 2.9.x           | master |

## FAQ

**Can I use an earlier Terraform version?**
Yes, but Terraform may not support all features and fields supported in your Rancher version so provisioning may be limited.

**Can I use a later Terraform version?**
Yes, but you must NOT use any new features and fields that your Rancher version does not support.
