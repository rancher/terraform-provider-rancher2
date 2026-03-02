# Rancher2 Provider

_For Terraform Maintainers Only_

QA has a new process to test Terraform using RCs which will enable us to better test this tool before a release is published.

### Test Process

1. Cut a Terraform RC.

Terraform RCs enable QA to test features and bug fixes _before_ the next Terraform version is released. This process is also how Rancher is tested and released. Terraform is expected to mirror Rancher's functionality and use of third party packages because it uses the Rancher API, so it's ideal and best practice to test the provider in the same way.

To cut an RC, run the `Manually Create RC Release` workflow specifying the release branch, the sha to tag, and the tag name.
For example, you might release `abcd1234` from `release/v13` as `v13.0.0-rc.1`
This will tag the sha from the release branch and generate an RC release for it.
The format for the tag must be `v<major>.<minor>.<patch>-rc.<rc number>`

2. Move the issue To Test

3. QA uses the test plan and the Terraform RC to test the feature or bug fix.

QA will test using a downloaded binary from the RC. Using an RC asset enables us to quickly test features and bugs without waiting for Hashicorp to publish a pre-release.

4. If test validation passes, please update the testing issue and close it. If not, please update the testing issue and work with the developer to resolve any issues.
