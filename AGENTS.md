## Environmental Specific Instructions

Dependencies for this project are provided by Nix, use this command to run scripts with the dependencies installed: `nix develop --ignore-environment --extra-experimental-features nix-command --extra-experimental-features flakes --keep HOME --keep SSH_AUTH_SOCK --keep GPG_SIGNING_KEY --keep NIX_SSL_CERT_FILE --keep NIX_ENV_LOADED --keep TERM --command bash -e {0}`

You may read files and folders in this repository without asking.

## Repository Coding Standards & Instructions

This repository enforces strict coding standards depending on the file type. Whenever you are asked to generate, edit, or review code, you MUST consult the corresponding instruction file for the specific rules, anti-patterns, and requirements:

* **For Go (`**/*.go`)**: Read and strictly adhere to `.github/instructions/go.instructions.md`
* **For Terraform (`**/*.tf`)**: Read and strictly adhere to `.github/instructions/terraform.instructions.md`
* **For GitHub Actions (`.github/workflows/**/*.yml`)**: Read and strictly adhere to `.github/instructions/workflows.instructions.md`
* **For GitHub Scripts (`.github/workflows/scripts/**/*.js`)**: Read and strictly adhere to `.github/instructions/github-script.instructions.md`
* **For Spelling Changes (`aspell_custom.txt`)**: Read and strictly adhere to `.github/instructions/aspell.instructions.md`
