---
name: github-pr
description: Manage, review, and address GitHub pull request comments and metadata. Use this skill whenever you need to find a PR, fetch conversation or inline review comments (including thread IDs), update PR titles/descriptions, respond to comments, or resolve review threads.
---

# GitHub PR Management

## Overview

The `github-pr` skill enables seamless local management of GitHub pull requests and review comment threads. It uses pure JavaScript and the Octokit SDK with automatic credential fallback (`GITHUB_TOKEN` -> default `gh` auth) to view, update, respond to, and resolve comments directly from the command line.

## Universal Options & Environment Variables

- **Help Option (`-h`, `--help`)**: Every script in this skill supports `-h` and `--help`. Run any script with `--help` to inspect its parameters, flag options, and usage examples without reading the source file:
  ```bash
  node ./scripts/<script_name>.js --help
  ```
- **Environment Variables**:
  - `GH_REPO` / `GITHUB_REPOSITORY`: Explicit target repository in `owner/repo` format (e.g. `rancher/terraform-provider-file`). Overrides git remote detection.
  - `GITHUB_TOKEN` / `GH_TOKEN`: GitHub personal access token. When unset, scripts automatically fall back to `gh auth token`.
- **Auth & Repository Inspection**:
  - Check your active credentials, detected branch, and detected PR number by running:
    ```bash
    node ./scripts/auth.js
    ```

## Core Tasks

### 1. View and Find a PR

To view the pull request for the current branch:

```bash
node ./scripts/find_pr.js
```

To view a specific PR (e.g. PR #123) or branch:

```bash
node ./scripts/find_pr.js 123
```

### 2. Read Comments and Review Threads

To read conversation comments and inline code review comment threads (which will also output thread IDs required for resolving):

```bash
node ./scripts/read_comments.js
```

To read comments for a specific PR, pass the PR number:

```bash
node ./scripts/read_comments.js 123
```

To filter and display only unresolved review threads (excluding already resolved discussions):

```bash
node ./scripts/read_comments.js [pr_number_or_current] --unresolved
```

To output full comments and review thread data in structured JSON format:

```bash
node ./scripts/read_comments.js [pr_number_or_current] --json
```

### 3. Respond to Comments and Reply to Reviews

- **Reply directly to a code review comment or thread**:
  To post a threaded reply directly to an inline review comment using its comment ID or GraphQL Thread ID (`PRRT_...`, from `read_comments.js`):

  ```bash
  node ./scripts/reply-comment.js <comment_id_or_thread_id> "Your reply text"
  ```

  You can also explicitly pass the PR number (or `current`), or use `-m` / `--body`:

  ```bash
  node ./scripts/reply-comment.js [pr_number_or_current] <comment_id_or_thread_id> [-m "Your reply text"]
  ```

- **Add general PR comment or reply via respond_comment**:
  - Add general conversation comment:
    ```bash
    node ./scripts/respond_comment.js [pr_number_or_current] general [-m "Your comment text"]
    ```
  - Reply to comment or thread:
    ```bash
    node ./scripts/respond_comment.js [pr_number_or_current] <comment_id_or_thread_id> [-m "Your reply text"]
    ```

### 4. Resolve and Unresolve Review Threads

To resolve an inline review thread on GitHub, use either the alphanumeric GraphQL Thread ID (starts with `PRRT_`, found via `read_comments.js`) or any numeric review comment ID within that thread:

```bash
node ./scripts/resolve_comment.js <thread_id_or_comment_id> [pr_number]
```

You can also pass `[pr_number_or_current]` first:

```bash
node ./scripts/resolve_comment.js [pr_number_or_current] <thread_id_or_comment_id>
```

To unresolve a previously resolved thread, pass `--unresolve`:

```bash
node ./scripts/resolve_comment.js <thread_id_or_comment_id> [pr_number] --unresolve
```

### 5. Update PR Title, Description, and State

To edit the title and/or description of a pull request:

```bash
node ./scripts/update_pr.js [pr_number_or_current] "New Title" "New Description Body"
```

You can also use `--title` / `-t`, `--body` / `-b`, and `--state` / `-s` flags to update title, body, or state (`open` or `closed`):

```bash
node ./scripts/update_pr.js [pr_number_or_current] [--title "New Title"] [--body "New Description Body"] [--state open|closed]
```

_Note: All commands accepting a `pr_number` also accept full GitHub PR URLs (e.g. `https://github.com/owner/repo/pull/123`)._
