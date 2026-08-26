# Google Gemini CLI Directory Structure Standard

This document outlines the official, standard directory structure and file layout for Google's Gemini CLI application at both the **User (Global)** and **Workspace (Project)** levels.

---

## 📂 Overview of the `.gemini/` Architecture

Gemini CLI uses a dual-tiered configuration model. Global preferences and personal tools reside in the user's home directory under `~/.gemini/`, while project-specific requirements, team-shared subagents, and automated skills are defined in the workspace root under `.gemini/`.

### 1. User (Global) Directory Layout (`~/.gemini/`)

Located in the user's home directory (`~` on Unix/macOS, `%USERPROFILE%` on Windows), this folder houses personal settings, authentication state, and global extensions:

```text
~/.gemini/
├── settings.json               # Global CLI configuration, UI preferences, and default model routing
├── trustedFolders.json         # Index of workspace directories that have been approved for local execution
├── trusted_hooks.json          # Index of cryptographically signed and trusted project hooks
├── oauth_creds.json            # Google OAuth / Gemini API personal credentials
├── GEMINI.md                   # Global personal memory file (loaded automatically in every session)
├── mcp-oauth-tokens.json       # Secure storage for MCP-associated OAuth access tokens
├── keybindings.json            # Global custom keybindings
│
├── commands/                   # Global custom slash commands
│   └── test.toml               # Registered globally as `/test`
│
├── agents/                     # Global personal custom subagents
│   └── security-auditor.md     # Loaded across all sessions for personal code audits
│
├── skills/                     # Global personal agent skills
│   └── my-skill/
│       ├── SKILL.md            # Skill specification file
│       └── scripts/            # Executable scripts/binaries
│
├── extensions/                 # Installed custom third-party extensions
│
└── tmp/                        # Local session logs, active chats, and system temporary state
    └── <project_hash>/
        ├── chats/              # Active and completed conversation history
        ├── shell_history       # Local interactive shell command execution logs
        └── otel/               # OpenTelemetry trace collector logs
```

---

### 2. Workspace (Project-level) Directory Layout (`your-project/.gemini/`)

Located in the root of your Git repository, this directory contains team-shared automations, strict codebase-level constraints, and pipeline validations. To maintain security, project-level hooks and settings are **untrusted by default** until approved by the developer during session initialization:

```text
your-project/
├── .geminiignore               # Gitignore-style patterns for excluding files from search/RAG indexing
├── .gemini/
│   ├── settings.json           # Workspace settings (overrides user-level settings)
│   ├── .env                    # Workspace-specific, secure environment variables
│   ├── system.md               # Workspace-level system prompt (overrides global CLI instructions)
│   ├── sandbox.Dockerfile      # Custom Docker configuration for fully containerized tool sandboxing
│   ├── sandbox-macos-*.sb      # Custom macOS sandbox profiles for localized secure execution
│   │
│   ├── hooks/                  # Synchronous lifecycle hooks triggered on CLI and tool events
│   │   ├── startup-context.sh  # Runs on SessionStart to load local context
│   │   ├── enforce-planning.js # Runs on BeforeTool to block unauthorized file writes
│   │   └── security.sh         # Runs on BeforeTool (shell) to block unsafe Git remote pushes
│   │
│   ├── agents/                 # Team-shared specialized custom subagents
│   │   └── review_agent.md     # Dedicated local PR review subagent
│   │
│   ├── skills/                 # Team-shared automation capabilities (tools)
│   │   └── commit-push.sh      # Custom command to safely sign and push commits
│   │
│   └── commands/               # Workspace-specific custom slash commands
│       └── changelog.toml      # Registered locally as `/changelog`
```

---

## ⚙️ Core Configuration Files Reference

### A. Settings File (`settings.json`)

The primary configuration engine for adjusting the CLI's behavior. Gemini CLI merges your home directory settings (`~/.gemini/settings.json`) with workspace settings (`.gemini/settings.json`), with workspace keys taking absolute precedence.

**Workspace Example:**

```json
{
  "ui": {
    "compactToolOutput": true,
    "showLineNumbers": true
  },
  "security": {
    "toolSandboxing": true,
    "folderTrust": {
      "enabled": true
    }
  },
  "skills": {
    "enabled": true
  },
  "hooksConfig": {
    "enabled": true,
    "notifications": true
  }
}
```

### B. Environment Variables (`.gemini/.env`)

A specialized file for storing API keys, GPG signing configurations, or other variables isolated to `gemini-cli`'s runtime.

- **Rule:** This file should contain developer-specific environment configurations and is automatically excluded from the AI's search context to prevent accidental leaks.

**Example:**

```bash
GPG_TTY=$(tty)
GEMINI_SYSTEM_MD=1
```

### C. System Prompt (`.gemini/system.md`)

A markdown document that injects repository-wide instructions, team conventions, and platform rules directly into the main model's system context, ensuring the AI behaves consistently with the project's standards.

---

## 🔒 Security & Best Practices

1. **Fingerprinting & Trust Gates**: Because workspace settings and hooks can run arbitrary shell commands, Gemini CLI automatically fingerprints `.gemini/settings.json` and `.gemini/hooks/` files. Any changes made to these files will prompt the developer with a trust gate warning before execution is allowed.
2. **Generic Portability**: To make team automations portable, workspace-level scripts inside `.gemini/` should be repository-agnostic and refer to the standard environment variables:
   - `GEMINI_PROJECT_DIR`: Absolute path to the active workspace.
   - `GEMINI_CWD`: Current working directory of the developer.
   - `GEMINI_SESSION_ID`: Unique conversation session identifier.
