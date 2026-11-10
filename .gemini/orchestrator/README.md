# Programmatic Orchestrator

## Abstract

This directory contains our developer-agent orchestration pipeline. Our workflow is proactively controlled by the central **Orchestrator** (`.gemini/orchestrator/orchestrator.js`), which enforces the Plan -> Implement -> QA -> Commit lifecycle programmatically using the `@google/gemini-cli-sdk` SDK.

---

## The Simplified 4-Phase Developer Loop

Our development lifecycle is steered programmatically by `orchestrator.js`:

1. **Phase 1: Planning (Read-Only)**: The orchestrator restricts the model to read-only actions and generates a plan. The orchestrator pauses and requests explicit developer approval via terminal `readline` before any edits are allowed.
2. **Phase 2: Implementation (Write Access)**: Once approved, the orchestrator initiates the implementation session, enabling the agent to surgically edit code and run local tests.
3. **Phase 3: Automated QA Review Loop**: The orchestrator executes the multi-agent Ralph Review Loop across all workflows. A thinking reviewer evaluates staged code and active plan constraints, a counting model tallies issue severity, and an inner implementation loop surgically resolves critical and medium findings while enforcing test and linter validation backpressure until clean convergence.
4. **Phase 4: Final User Review & Commit**: The orchestrator shows the unified git diff to the developer in the terminal, asks for confirmation, and creates a signed commit.

---

## 1. Core Orchestrator (`.gemini/orchestrator/orchestrator.js`)

The central execution engine using the official `@google/gemini-cli-sdk` SDK. It coordinates all tool invocations (such as file reads/writes, globbing, running commands, and asking for user feedback) and drives the four development phases sequentially.

---

## 2. Supporting Utilities

- **`install-sdk.js`**: Automated builder and installer for `@google/gemini-cli-sdk`.
- **`inspect-prompts.js`**: Inspection utility for orchestrator prompts and agent declarations.
- **`exercise-agents.js`**: Automated multi-model quota check and daily smoke test.
- **`exercise-cron.sh`**: Headless runner for agent health checks.
- **`lib/`**: Deterministic modular libraries supporting agent execution, Ralph review loop, and QA pipeline.
