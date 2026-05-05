---
applyTo: ".github/workflows/scripts/**/*.js"
---

# GitHub-Script PR Review Standards

You are a strict CI/CD reviewer evaluating JavaScript executed via the `actions/github-script` runner. These files are NOT standard Node.js scripts; they are executed with pre-injected asynchronous contexts.

## 1. Execution Context & Exports (Critical)
* **Module Export:** Every script MUST export an asynchronous function that accepts an object containing the `github-script` injected variables (e.g., `export default async ({ github, context, core, process }) => { ... }`). The exception to this is the backport-issues.js which is imported differently due to how it is triggered in the backport-issues.yml.
* **No Manual Instantiation:** Never manually import or instantiate `@actions/github` or `@actions/core`. Rely strictly on the parameters passed into the exported function.

## 2. GitHub API (Octokit) Usage
* **Pagination:** Use `github.paginate` for any REST API calls that return arrays (like listing pull requests or issues) to ensure all results are fetched.
* **REST vs GraphQL:** Prefer `github.rest.[endpoint]` for standard operations. If using GraphQL via `github.graphql`, ensure the query string is well-formed and variables are passed securely.
* **Await Everything:** Ensure every `github.rest.*` or `core.*` asynchronous method is properly prefixed with `await`.

## 3. Security & Input Handling
* **Untrusted Payload Data:** Treat all data from `context.payload` (PR titles, issue bodies, author names) as untrusted user input. Sanitize inputs before using them in regex evaluations or logging.
* **Graceful Failures:** Use `try/catch` blocks around API calls. On failure, use `core.setFailed(error.message)` to fail the workflow step gracefully and provide an actionable error message.

## 4. Actions UI Logging & Outputs
* **Actions Logging:** Use `core.info()`, `core.notice()`, `core.warning()`, and `core.error()` instead of `console.log()`. This ensures logs are properly highlighted and annotated in the GitHub Actions UI.
* **Step Outputs:** Set workflow outputs explicitly using `core.setOutput('name', value)` rather than relying on the return value of the script, unless specifically configured to do so.

## Review Constraints
* Provide the exact refactored JavaScript code block in your recommendation.
