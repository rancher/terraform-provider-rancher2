---
name: git
description: Execute safe Git operations and workflows including branch creation with fork sync, compliant commits, reading gitignored files, status checks, diff reviews, and history. Use whenever agents need to interact with git or inspect repository changes.
---

# Git Management and Safe Operations

## Overview

The `git` skill provides a safe, curated, and standardized set of Git utilities for agents written in JavaScript using `simple-git`. It protects against accidental state corruption, prevents commits directly to `main`/`master`, ensures commit messages comply with Conventional Commits and product boundaries, allows reading gitignored files (unless excluded by `.aiexclude`), and prevents unintended staging.

**Rule:** You should always prefer using the scripts in this skill over running raw, direct `git` commands (such as `git diff`, `git status`, `git commit`, `git checkout`, etc.) to interact with repository state.

## Universal Options & Script Discovery

- **Help Option (`-h`, `--help`)**: Every script in this skill supports `-h` and `--help`. Whenever you need to discover available arguments, positional parameters, or flags for any script, execute it with `--help` instead of reading the source file:
  ```bash
  node ./scripts/<script_name>.js --help
  ```

## Core Tasks

### 1. Check Workspace and Branch Status

- **Working Tree & Branch Status**: View uncommitted changes and list files changed in the current branch compared to the base branch (`main`, `master`, or specified `-b base-branch`):
  ```bash
  node ./scripts/get_files_changed.js [-b base-branch]
  ```
- **List Modified/Untracked Paths**: Get a clean list of explicit file paths that have been modified or are untracked:
  ```bash
  node ./scripts/get_modified_paths.js
  ```
- **Inspect Staged Files in Index**: Inspect what files are currently staged, their stage slot (0 for normal, 1-3 for merge conflicts), and their cached diff status (`A`, `M`, `D`, etc.), or retrieve just file names with `--name-only`:
  ```bash
  node ./scripts/get_staged_files.js [--name-only]
  ```

### 2. Review Detailed Diffs

- **Uncommitted Changes (Staged & Unstaged)**: View a full diff of local changes in the working directory and staging area, optionally filtering by specific files, scoping to `--staged` or `--unstaged`, or requesting a compact summary with `--stat` or `--name-only`:
  ```bash
  node ./scripts/get_diff_to_review.js [--staged|--unstaged] [--stat|--name-only] [file_path...]
  ```
- **Branch Changes**: View all changes introduced in the current branch relative to the base branch (`main`, `master`, or specified `-b base-branch`), optionally filtering by specific files or requesting a compact summary with `--stat` or `--name-only`:
  ```bash
  node ./scripts/get_branch_changes.js [--stat|--name-only] [-b base-branch] [file_path...]
  ```

### 3. Read Files Ignored by .gitignore

To safely read files that are ignored by `.gitignore` (such as plans, reports, or local configurations) while ensuring repository `.aiexclude` privacy boundaries are strictly respected:

```bash
node ./scripts/read-ignored-file.js <file_path> [start_line] [end_line]
```

### 4. Branch Management

- **Create New Branch (with Fork Upstream Sync)**:
  Syncs the local base branch (`main` or `master`) from the upstream remote before creating and checking out a new feature branch:
  ```bash
  node ./scripts/create-branch.js <branch-name> [base-branch]
  ```
- **Switch Branch**:
  Safely switches between existing local or remote branches:
  ```bash
  node ./scripts/switch-branch.js <branch-name>
  ```

### 5. Staging and Unstaging Changes

To stage specific modified files while ensuring files blocked by `.aiexclude` cannot be accidentally staged:

```bash
node ./scripts/stage-files.js <file1> [file2...]
```

To stage all modified tracked files:

```bash
node ./scripts/stage-files.js --all
```

To stage all modified tracked files and untracked files (excluding gitignored files):

```bash
node ./scripts/stage-files.js --all --include-untracked
```

To unstage specific files from the staging index:

```bash
node ./scripts/unstage-files.js <file1> [file2...]
```

To unstage all staged changes:

```bash
node ./scripts/unstage-files.js --all
```

### 6. Create Compliant Commits

To create a commit that blocks commits to `main`/`master`, validates against Conventional Commits, verifies product boundary semver rules (changes outside `internal/` cannot use `feat` or `!`), and runs spell checks:

```bash
node ./scripts/commit.js -m "fix: your commit message"
```

### 7. View Commit Log and File History

- **View Recent Commits**: View a concise, color-stripped oneline commit history on the current branch or scoped to a specific file or branch:
  ```bash
  node ./scripts/get_commit_log.js [-n limit] [-f file] [-b branch]
  ```
- **Retrieve Deleted/Historical File Contents**: Retrieve contents of a deleted or historical file from Git history at the latest commit or a specific revision/commit SHA:
  ```bash
  node ./scripts/get_deleted_file.js <file_path> [commit_or_rev]
  ```
