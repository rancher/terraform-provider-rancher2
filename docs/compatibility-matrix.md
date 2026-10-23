# Rancher2 Provider

## Version compatibility matrix

The version matrix specifies the Terraform provider version to use with the associated minor Rancher version that it was released for.
When updating the version matrix, add a row for each Terraform version released with a minor Rancher release.
This shouldn't be a changelog for every branch, each branch has its own CHANGELOG.md file which should be updated with changes.

| Terraform provider version | Rancher version | Terraform provider branch |
|----------------------------|:---------------:|---------------------------|
| 2.x                        | 2.6.x           | release/v2 |
| 3.x                        | 2.7.x           | release/v3 |
| 4.x                        | 2.8.x           | release/v4 |
| 5.x                        | 2.9.x           | release/v5 |
| 6.x                        | 2.10.x          | release/v6 |
| 7.x                        | 2.11.x          | release/v7 |
| 8.x                        | 2.12.x          | release/v8 |
| 13.x                       | 2.13.x          | release/v13 |
| 14.x                       | 2.14.x          | release/v14 |
| 15.x                       | 2.15.x          | release/v15 |
| 16.x                       | 2.16.x          | release/v16 |

## FAQ

**Can I use an earlier or later Rancher provider version?**
Using incorrectly correlated versions isn't supported at this time.
