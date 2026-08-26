# Agentic Framework: Boilerplate Sync Skill

## Abstract

To maintain cross-project consistency across multiple repositories without the heavy administrative overhead of Git submodules, the Agentic Framework implements a lightweight, **manifest-driven file synchronization skill** (`sync-boilerplate.sh`). This tool allows developers and agents to dynamically track, compare, pull, and push common configuration files (such as linters, CI/CD scripts, and enforcer hooks) to and from a centralized master template repository, keeping all systems fully aligned.

---

## Technical Specification

### 1. The Sync Manifest (`.boilerplate-sync.json`)

Each child repository declares its tracking files, local destinations, and master template repository source inside a `.boilerplate-sync.json` configuration file located in the repository's root directory:

```json
{
  "template_repo": "git@github.com:your-organization/your-boilerplate-template.git",
  "files": [
    { "local": ".prettierrc", "remote": "shared-configs/.prettierrc" },
    { "local": ".golangci.yml", "remote": "shared-configs/.golangci.yml" },
    { "local": "cspell.json", "remote": "shared-configs/cspell.json" },
    { "local": ".github/workflows/scripts/lint.sh", "remote": "scripts/lint.sh" }
  ]
}
```

### 2. Operational Logic & Sequence

The sync utility `.gemini/skills/sync-boilerplate.sh` executes the following sequence to compare or copy files:

```text
[Local Repo]                     [System OS / TMP]                     [Remote Repo]
     │                                   │                                   │
     │ ── 1. Read Manifest JSON ───────► │                                   │
     │                                   │ ── 2. git clone --depth 1 ──────► │
     │                                   │ ◄── 3. Shallow Local Checkout ─── │
     │                                   │                                   │
     │ ◄── 4. Execute 'diff' or 'cp' ─── │                                   │
     │                                   │                                   │
     │                                   │ ── 5. Trap: rm -rf /tmp/clone ──► │
```

1. **Manifest Parsing**: Reads and validates the JSON fields `.template_repo` and `.files` array using `jq`.
2. **Hermetic Sandbox Prep**: Establishes a temporary workspace directory under `/tmp/boilerplate-sync-XXXXXX` using `mktemp -d`.
3. **Repository Fetching**: Clones the remote master repository off the parent branch using `git clone --depth 1 --no-checkout <repo_url> <tmp_dir>`, then checks out strictly the tracked files to minimize network and disk footprints.
4. **Operations Modes**:
   - **Diff Mode (`--diff`)**: Runs `git diff --no-index` or standard `diff -u` between the local file and its remote template file counterpart.
   - **Sync/Pull Mode (`--pull`)**: Overwrites the local file by copying the template file into place, creating any missing parent directories natively.
   - **Sync/Push Mode (`--push`)**: Copies local files that differ back to the remote template clone, commits them conventionally, and pushes the updates to the template repository using native developer credentials.
5. **Secure Workspace Cleanup**: Registers an exit trap (`trap 'cleanup' EXIT`) that mathematically guarantees the temporary directories are fully destroyed on exit, preventing `/tmp` clutter or memory leak vectors.

---

## Standing Implementation Decisions

1. **Environment Variable Override**: If the `CENTRAL_FILE_REPO` environment variable is defined, the utility automatically overrides the manifest's `.template_repo` with it, allowing safe targeting of private central repositories without hardcoding sensitive URLs in version-controlled JSON manifests.
2. **Standard Commits on Push**: When pushing updates back to the centralized template repository, changes are always committed using standard, non-bumping conventional commits (e.g. `build(sync): sync configurations`) to maintain pristine release-please semantics in the master repository.
3. **No-Checkout Optimization**: `git clone` always uses `--no-checkout` to avoid pulling unnecessary files, checking out only the exact files mapped in `.boilerplate-sync.json`.
