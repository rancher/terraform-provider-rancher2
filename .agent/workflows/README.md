# AI Agent Workflows

This directory contains defined, step-by-step procedures that AI agents must follow when executing complex, multi-step tasks in this repository. 

Using these workflows ensures maximum consistency, rigorous quality control, and clean engineering practices.

## Available Workflows

### 1. [Standard Development Process](development-process.md)
* **Purpose:** Outlines the lifecycle for developing new features, applying bug fixes, and performing refactoring.
* **Key Steps:** Exploration, bug reproduction, plan-creation, surgical edits, format/compiles, unit testing, schema documentation updates, and lint validations.

### 2. [Troubleshooting CI/CD Workflows](troubleshoot-workflows.md)
* **Purpose:** Explains how to diagnose, triage, and repair broken GitHub Actions or release workflows.
* **Key Steps:** Log retrieval, error isolation, script/YAML auditing, secret token sanitization, and verification with `actionlint` and `shellcheck`.

### 3. [PR Review Comment Resolution](resolve-pr-reviews.md)
* **Purpose:** Defines the repeatable, high-standard process for analyzing, planning, and executing resolutions for Pull Request review comments.
* **Key Steps:** Review retrieval, separation of concerns (discernment phase), plan updating, responding on GitHub, surgical refactoring, and programmatic thread resolution.
