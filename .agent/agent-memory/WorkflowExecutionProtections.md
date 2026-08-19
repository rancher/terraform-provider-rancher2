# GitHub Actions Policies & Workflow Execution Protections

* **Source Documentation:** 
  - [About Actions Policies](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/actions-policies/about-actions-policies)
  - [Workflow Execution Protections](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/actions-policies/workflow-execution-protections)
* **Status:** Public Preview (As of June/July 2026)

---

## 1. Context & Architectural Overview

GitHub Actions Policies and Workflow Execution Protections introduce a centralized, platform-level governance layer to GitHub Actions. This closes a critical security vector and provides administrators with a declarative way to restrict workflow execution.

Previously, workflow execution was dictated strictly by the YAML definitions found within the triggering commit/ref. A malicious actor with write access (or via an unreviewed fork PR) could manipulate a workflow file to run arbitrary code with privileged tokens. This new policies framework shifts the trust boundary from the distributed, easily compromised codebase up to a centralized repository/organization ruleset.

---

## 2. Core Concepts

### Actions Policies
A dedicated administration interface separate from the existing *Actions > General* settings. It is designed to host multiple governance policies over time. Currently, **Workflow Execution Protections** is the primary policy type available.

### Workflow Execution Protections
An allowlist mechanism configured via GitHub Rulesets that acts as a gatekeeper *prior* to a workflow run starting. When an event occurs, GitHub evaluates the configured ruleset. If the event or the actor who initiated it is not permitted, the workflow execution is blocked before any runner is provisioned or code is executed.

---

## 3. Ruleset Features & Configurations

Because protections are built on top of the GitHub Rulesets framework, they inherit several advanced enterprise-grade controls:

### Targeting & Scoping
* **Multi-Level Enforcement:** Policies can be established and enforced at the **Enterprise**, **Organization**, or individual **Repository** levels.
* **Custom Property Scoping:** Organizations can target rulesets to groups of repositories dynamically using **Repository Custom Properties** (e.g., applying strict protections only to repositories marked as `public` or `tier-1`).

### Deployment Modes
* **Active:** Enforces the policies immediately, blocking unauthorized events/actors.
* **Evaluate Mode:** Runs the policies in a passive "shadow" mode. Blocked workflows are logged in Policy Insights and audit logs, but execution is not halted. This allows administrators to safely test and audit policies before active enforcement.

---

## 4. Supported Policy Rule Options

> 💡 **Maintainer Note:** Because this feature is currently in public preview, these details and options are subject to change by GitHub. Always refer to the [authoritative GitHub documentation](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/actions-policies/workflow-execution-protections) for the most up-to-date information.

As of **August 2026**, the Actions Policies user interface exposes **three primary options/rules** for configuring protections:

### A. Restrict Actors (Who can trigger workflows)
Controls which identities are allowed to trigger workflow runs. This separates the privilege of pushing code from the privilege of executing CI, preventing unreviewed/untrusted contributions from executing actions.
* **Supported Actor Types (as of August 2026):**
  * Individual GitHub Users
  * Repository Roles (e.g., `Read`, `Maintain`, `Admin`, `Triage`)
  * Installed GitHub Apps (including the system-level **GitHub Actions** App for GITHUB_TOKEN pushes)
  * GitHub Copilot
  * Dependabot (`dependabot[bot]`)

### B. Restrict Events (Which events can trigger workflows)
Declares exactly which GitHub Actions webhook events are permitted to trigger runs on the repository.
* **Supported Events (as of August 2026):**
  * `push`
  * `pull_request`
  * `pull_request_target`
  * `workflow_dispatch`

### C. Require Lockfile
Enforces strict locking requirements on configuration files (e.g., locking down workflow configurations or package dependencies) to prevent accidental or malicious modifications during execution.

---

## 5. Security & Threat Mitigations

Workflow execution protections are highly effective at disrupting several common pipeline attacks:

| Attack Vector | Threat Description | Ruleset Mitigation Strategy |
| :--- | :--- | :--- |
| **Poisoned Pipeline Execution (PPE)** | Attackers submit a PR from a fork modifying workflow files or exploiting `pull_request_target` to execute malicious code in the privileged base context. | Use **Restrict Events** to prohibit or limit `pull_request_target` runs. |
| **Manual-Trigger Abuse** | Non-maintainer contributors with write access invoke sensitive deployment/release pipelines manually. | Use **Restrict Actors** to limit `workflow_dispatch` to specific roles (e.g. `Maintain` or `Admin`). |
| **Untrusted Actor Runs** | External contributors or compromised integration tokens trigger costly automated runs. | Use **Restrict Actors** to block low-trust identities or unauthorized bots from initiating any workflow execution. |
| **Workflow File Tampering** | Attackers bypass file-level branch protection rules (like CODEOWNERS) to force-execute compromised workflows. | Centralized repository policies override any configuration declared in the local workspace YAML files. |

---

## 6. Setup and Administration Instructions

To configure these protections in GitHub:

1. Navigate to your repository on GitHub.
2. Under the repository name, click **Settings**.
3. In the left sidebar, locate the **Actions** section and click **Policies**.
4. Create a new ruleset under the Policies screen (note: this applies globally to the repository; there is **no branch selector**).
5. Add and configure the desired rules:
   * **Restrict Actors** (e.g., add the **GitHub Actions** App or write-level roles to allow GITHUB_TOKEN pushes to trigger downstream release workflows).
   * **Restrict Events** (e.g., permit **`push`** to enable release pipelines, while restricting dangerous events).
   * **Require Lockfile** (to secure execution state).
6. Select your enforcement status (**Evaluate** to test, **Active** to enforce).
7. Click **Create** or **Save changes**.
