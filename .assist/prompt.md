## PROMPTS FOR AI ASSISTANTS
- Dependencies for this project are provided by Nix, use this command to run scripts with the dependencies installed: `nix develop --ignore-environment --extra-experimental-features nix-command --extra-experimental-features flakes --keep HOME --keep SSH_AUTH_SOCK --keep GPG_SIGNING_KEY --keep NIX_SSL_CERT_FILE --keep NIX_ENV_LOADED --keep TERM --command bash -e {0}`
- Scripts in the `.github/workflows/scripts` directory are actually actions/github-script scripts, not just javascript
- The following links have information on github-script scripts:
  - https://github.com/actions/github-script?tab=readme-ov-file#this-action
  - https://octokit.github.io/rest.js/v22/#custom-requests replace octokit with github in the examples
- You may read files and folders in this repository without asking.
