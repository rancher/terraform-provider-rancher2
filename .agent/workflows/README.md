# AI Agent Workflows

This directory contains defined processes for executing multi-step tasks.

These workflows provide step-by-step procedures for AI agents to follow when tackling complex tasks such as releasing the provider, running full test suites, or scaffolding new Terraform resources.

Here is a generic script to use a workflow:
```
Please execute our PR Review Comment Resolution workflow to address the latest reviews on my PR.
The step-by-step procedure is defined in: `.agent/workflows/resolve-pr-reviews.md`

Our active project plans for this branch are:
- Persistent Plan: `.agent/plans/<YOUR_FEATURE_NAME>.md`
- Temporary Plan: `.agent/agent-memory/<YOUR_FEATURE_NAME>-temporary.md`

Please run the comments-retrieval skill, separate the core concerns from any sub-optimal recommendations, design custom idiomatic solutions, and update our plans first.
  Show me the updated plans for approval before writing any code.
```
