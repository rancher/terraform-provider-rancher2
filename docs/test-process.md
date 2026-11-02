# Rancher2 Provider

### Test Process

To cut an RC, run the `Release` workflow specifying the release branch, the sha to tag, and the tag name.
For example, you might release `abcd1234` from `release/v15` as `v15.0.0-rc.1`
This will tag the sha from the release branch and generate an RC release for it.
The format for the tag must be `v<major>.<minor>.<patch>-rc.<rc number>`

1. Move the issue To Test
2. QA uses the test plan to test the feature or bug fix.
3. If test validation passes, please update the testing issue and close it.
4. If not, please update the testing issue and work with the developer to resolve any issues.

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
- cp the binary to my repo's ./bin directory.  
- Update your terraform rc file to the absolute path, eg. `"/Users/matttrach/go/src/github.com/rancher/terraform-provider-rancher2/bin"` (no file name)
- Make sure your 'TF_CLI_CONFIG_FILE' env variable is pointing to your terraform rc file.  
- Run `make test` to test the binary in `./bin`; do not use `run_tests.sh` for this downloaded RC because it runs `make build` first and overwrites the binary.
