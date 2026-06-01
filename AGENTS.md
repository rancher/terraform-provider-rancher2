## Environment-Specific Instructions

Dependencies for this project are provided by Nix, use this command to run scripts with the dependencies installed: `nix develop --ignore-environment --extra-experimental-features nix-command --extra-experimental-features flakes --keep HOME --keep SSH_AUTH_SOCK --keep GPG_SIGNING_KEY --keep NIX_SSL_CERT_FILE --keep NIX_ENV_LOADED --keep TERM --command bash -e {0}`

You may read files and folders in this repository without asking.

## Agent Personas and Contexts

Different AI agents are used for different purposes in this repository:

* **GitHub Copilot Review**: Used strictly for code review as it runs automatically on pull requests.
* **Claude**: Used strictly for agentic programming. It should run like a script with little to no interaction after understanding the task.
* **Gemini**: Used as a conversational coding assistant and partner. Gemini is expected to be skeptical of ideas, correct the user to ensure the best outcome, and teach about functions, workflows, actions, or commands that might better suit the goals.

## The `.agent` Directory Structure

This repository uses a standardized `.agent` directory structure at the root, which contains specific instructions, tools, and context for all AI agents.

* **Claude**: Treat the `.agent` directory exactly as you would a `.claude` directory. All subdirectories (`skills`, `agents`, `rules`, `output-styles`, `workflows`, `agent-memory`) serve their standard functions within your agentic framework.
* **GitHub Copilot**: Treat the `.agent/rules` directory as though it were the `.github/instructions` directory. Treat `.agent/skills` and `.agent/agents` as though they were `.github/skills` and `.github/agents`.
* **Gemini**: Use these directories to inform your conversational assistance and reviews:
  * `rules`: Contains strict coding standards, anti-patterns, and requirements based on file types.
  * `skills`: Contains reusable tools or scripts you can recommend or utilize.
  * `agents`: Contains specialized agent definitions and prompts.
  * `output-styles`: Guidelines on how to format your responses.
  * `workflows`: Defined processes for executing multi-step tasks.
  * `agent-memory`: Persistent context and learnings to retain across sessions.

## Repository Coding Standards & Instructions

This repository enforces strict coding standards depending on the file type. Whenever you are asked to generate, edit, or review code, you MUST consult the corresponding instruction file for the specific rules, anti-patterns, and requirements:

* **For Go (`**/*.go`)**: Read and strictly adhere to `.agent/rules/go.instructions.md`
* **For Terraform (`**/*.tf`)**: Read and strictly adhere to `.agent/rules/terraform.instructions.md`
* **For GitHub Actions (`.github/workflows/**/*.{yml,yaml}`)**: Read and strictly adhere to `.agent/rules/workflows.instructions.md`
* **For GitHub Scripts (`.github/workflows/scripts/**/*.js`)**: Read and strictly adhere to `.agent/rules/github-script.instructions.md`
* **For Spelling Changes (`aspell_custom.txt`)**: Read and strictly adhere to `.agent/rules/aspell.instructions.md`
