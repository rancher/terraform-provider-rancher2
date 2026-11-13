---
name: github-ci
description: Inspect, troubleshoot, and resolve GitHub Action workflow failures. Use this skill whenever a CI check fails, when you need to list workflow runs, check failed runs, download job/step logs, or search/grep through downloaded log files to identify compilation or testing errors.
---

# GitHub CI / Actions Troubleshooting

## Overview

The `github-ci` skill provides tools to investigate and resolve failed GitHub Action workflow runs using pure JavaScript and the Octokit SDK. It helps you quickly list workflow statuses, isolate failures, parse failed step logs directly, download run logs locally, and search through them for compiler errors, test failures, or syntax issues without leaving the CLI. All scripts automatically provide fallback authentication (`GITHUB_TOKEN` -> default `gh` auth).

**Rule:** You should always use the scripts in this skill instead of running direct `gh` CLI commands when investigating CI failures or fetching workflow logs.

## Universal Options & Environment Variables

- **Help Option (`-h`, `--help`)**: Every script supports `-h` and `--help`. Run any script with `--help` to view parameter syntax, flags, and usage without reading the script file:
  ```bash
  node ./scripts/<script_name>.js --help
  ```
- **Environment Variables**:
  - `GH_REPO` / `GITHUB_REPOSITORY`: Explicit target repository in `owner/repo` format (e.g. `rancher/terraform-provider-file`). Overrides git remote detection.
  - `GITHUB_TOKEN` / `GH_TOKEN`: GitHub personal access token or app installation token. When omitted, scripts automatically fall back to the token stored in GitHub CLI (`gh auth token`).
  - Verify authentication and repository detection anytime by running:
    ```bash
    node ./scripts/auth.js --help
    ```

## Core Tasks

### 1. List Actions / Workflow Runs

To view the status of recent workflow runs on the current branch:

```bash
node ./scripts/list_runs.js [-n limit] [--status completed|in_progress|queued]
```

To check runs on a specific branch or across all branches, pass them as arguments:

```bash
node ./scripts/list_runs.js [branch_name_or_all_or_current] [-n limit] [--status status]
```

### 2. Isolate Failed Runs

To list only failed workflow runs on the current branch (to find their run IDs):

```bash
node ./scripts/find_failed_runs.js [-n limit] [--status completed|in_progress|queued]
```

To check failed runs on a specific branch or across all branches:

```bash
node ./scripts/find_failed_runs.js [branch_name_or_all_or_current] [-n limit] [--status status]
```

### 3. Fetch and Parse Failed Step Logs

To quickly extract and parse the exact failure details and error messages from a failed run without downloading and wading through full logs:

```bash
node ./scripts/get_failed_log.js [-r run_id] [-b branch] [run_id_or_branch]
```

If `run_id_or_branch` is omitted or `current`, it automatically targets the most recent failed workflow run on the current branch. If a branch name is supplied (e.g. `main` or `-b main`), it targets the latest failed run on that branch.

### 4. Download Job and Step Logs

For in-depth analysis when full log dumps are required, download and extract all job and step logs into `/tmp/ci-logs-<run_id>`:

```bash
node ./scripts/download_logs.js [run_id] [--failed-only]
```

If `run_id` is omitted or set to `current`, it automatically targets the most recent failed workflow run on the current branch. You can also specify a custom output directory or download only failed job logs with `--failed-only`:

```bash
node ./scripts/download_logs.js [run_id] [/path/to/directory] [--failed-only]
```

### 5. Search and Grep Logs for Errors

Use the search tool to scan through a downloaded logs folder for error messages, test failures, or exceptions. It uses regex patterns to catch common failure keywords (e.g., error, fail, exit status, panic, exception):

```bash
node ./scripts/search_logs.js <logs_directory> [regex_pattern] [-n limit]
```

To search for a specific error message or component failure, pass a custom pattern:

```bash
node ./scripts/search_logs.js <logs_directory> "cannot find package"
```
