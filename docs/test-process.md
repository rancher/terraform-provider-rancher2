# Rancher2 Provider: Test & QA Process

This document outlines the testing and QA verification process for Release Candidate (RC) and release validation.

For the end-to-end architecture of GitHub workflows and automation triggers, see **[GitHub Workflows & Automation Architecture](workflows.md)**.

### Release Candidate (RC) Flow

Release Candidates are generated automatically by CI:

1. **Automated RC Tagging**: When a backport PR merges into an active release branch (e.g., `release/v15`), automation (`release-please.yml`) automatically tags the commit as the next RC (e.g., `v15.2.0-rc.1`) and dispatches the build and release workflow.
2. **Project Board Progression**: Merging into a release branch automatically transitions the related backport issue on [Rancher2 Terraform Provider](https://github.com/orgs/rancher/projects/73) to **Ready for QA**.
3. **Manual Dispatch (Optional)**: If a specific commit SHA needs to be tagged out-of-band, maintainers can trigger the `Release` workflow manually specifying the target release branch, commit SHA, and tag name (`v<major>.<minor>.<patch>-rc.<N>`).

### QA Testing Lifecycle

1. **Identify Issues to Test**: Filter [Rancher2 Terraform Provider](https://github.com/orgs/rancher/projects/73) for issues in the **Ready for QA** status.
2. **Test Validation**: Execute the test plan using the downloaded RC binary or a locally built binary.
3. **Passing Validation**: When validation passes, comment with test results and close the backport sub-issue.
   - When all sub-issues linked to a tracking issue are closed and a milestone is assigned, automation closes the parent tracking issue and moves it to **Ready For Release**.
4. **Failing Validation**: If testing uncovers defects, update the sub-issue with reproduction details and collaborate with the engineering team to resolve them.

### Configuring Terraform to use a Local Binary

There are two configuration objects you need to configure Terraform to use a local binary.

1. A `terraform.rc` or `terraform.tfrc` or `.terraformrc` file.
2. An environment variable with the full path to that file.

Environment variable: "TF_CLI_CONFIG_FILE"
eg. `export TF_CLI_CONFIG_FILE="full/path/to/terraform.rc"`

This config tells Terraform to get the rancher/rancher2 provider from a path, while directly downloading everything else (exclude nothing) from the Terraform registry:

```hcl2
provider_installation {
  dev_overrides {
    "rancher/rancher2" = "<full path to your repo>/bin"
  }
  direct {
    exclude = []
  }
}
```

The `make build` command places a locally built binary in the repo's directory in ./bin folder (ignored).
This allows the same configuration to work for a binary built from source as well as a downloaded binary from a release.

### Downloading Release Candidate Binary

- Download the proper archive for your hardware/OS from a release candidate release, eg. [v15.2.0-rc.4](https://github.com/rancher/terraform-provider-rancher2/releases/tag/v15.2.0-rc.4).  
  https://github.com/rancher/terraform-provider-rancher2/releases/download/v15.2.0-rc.4/terraform-provider-rancher2_15.2.0-rc.4_darwin_arm64.zip
- Unzip the archive and the binary will be named something like this: `terraform-provider-rancher2_v15.2.0-rc.4`.
- Copy the binary to a trusted directory, make it executable, and rename it to `terraform-provider-rancher2`.

For local testing:

- Copy the binary to the repository's `./bin` directory (`./bin/terraform-provider-rancher2`).
- Update your Terraform rc file to the absolute path, e.g. `"/Users/username/go/src/github.com/rancher/terraform-provider-rancher2/bin"` (directory path without filename).
- Ensure your `TF_CLI_CONFIG_FILE` environment variable points to your Terraform rc file.
- Run `make test` to run unit tests.
- Run `make build` to build the provider binary if testing from source rather than a release-candidate binary.
- To run full e2e acceptance tests using the downloaded RC binary without overwriting it, run `./run_tests.sh -b` (the `-b` flag or `SKIP_BUILD=true` skips `make build`).
  - Ensure your AWS credentials are configured properly before running e2e tests.
  - E2e tests can take multiple hours to complete.
