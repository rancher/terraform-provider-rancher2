# Rancher2 Provider

_For Terraform Maintainers Only_

QA has a new process to test Terraform using RCs which will enable us to better test this tool before a release is published.

### Test Process

1. Cut a Terraform RC.

Terraform RCs enable QA to test features and bug fixes _before_ the next Terraform version is released. This process is also how Rancher is tested and released. Terraform is expected to mirror Rancher's functionality and use of third party packages because it uses the Rancher API, so it's ideal and best practice to test the provider in the same way.

To cut an RC, tag the latest commit. Bump the minor version for new features and the patch version for bug fixes only.

For example, if the latest tag on `master` is `v3.0.0` and the commits since that tag contain a new feature,

```sh
$ git tag v3.1.0-rc1
$ git push upstream v3.1.0-rc1
```

If all commits since the latest tag are only bug fixes,

```sh
$ git tag v3.0.1-rc1
$ git push upstream v3.0.1-rc1
```

If the latest tag is already an RC, say `v3.0.1-rc1`,

```sh
$ git tag v3.0.1-rc2
$ git push upstream v3.0.1-rc2
```

2. Move the issue To Test
3. QA uses the test plan and the Terraform RC to test the feature or bug fix using [this script](https://github.com/rancher/terraform-provider-rancher2/blob/master/setup-provider.sh). How to instructions are included in the script.

QA will test using a downloaded binary from the RC. Using an RC asset enables us to quickly test features and bugs without waiting for Hashicorp to publish a pre-release. This script sets up the correct binary using a defined `<provider> <version>` to test updates locally.

```sh
$ ./setup-provider.sh <provider> <version>
```

 Example
```sh
$ ./setup-provider.sh rancher2 v3.0.0-rc1
```

There is also a [windows script](https://github.com/rancher/terraform-provider-rancher2/blob/master/setup-provider-windows.ps1) for cross-platform support.

```powershell
PS /> ./setup-provider-windows.ps1 <provider> <version>
```

Example
```powershell
PS /> ./setup-provider-windows.ps1 rancher2 v3.0.0-rc1
```

4. If test validation passes, issue is closed. If not, it is reopened.
