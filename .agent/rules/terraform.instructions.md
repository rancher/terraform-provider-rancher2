---
applyTo: "**/*.tf"
---

# Terraform PR Review Standards

As a strict Infrastructure-as-Code reviewer, enforce these standards on all Terraform changes. Flag violations with a concise explanation and provide the refactored code block.

## 1. Validations (Critical)
* **Simple Validations:** Use `variable` validations for simple validations (e.g. regex matching on a single variable).
* **Complex Validations (Check Blocks):** Use `check` blocks for complex validations. Validations which include inputs from other variables MUST go here.
* **Complex Validations (Preconditions):** Use `terraform_data` preconditions for complex validations that cannot be evaluated in `check` blocks. Validations which include outputs from other resources MUST go here.

## Review Constraints
* Assume the codebase uses `terraform fmt`. DO NOT comment on spacing or alignment.
* Provide the exact refactored Terraform code block in your recommendation.
