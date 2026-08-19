# Planning Changes & Long-Lived Architectural Blueprints

All plans MUST be documented in a single, unified plan file located under the `.agent/plans/` directory (e.g., `.agent/plans/<PlanName>.md`). We no longer use separate temporary plans in `.agent/agent-memory/`.

## Plans as Architectural Blueprints

Plans are not merely transient, temporary check-lists; they are foundational, long-lived architectural documents and technical specifications for the repository. Every plan file is structured to represent a distinct architectural domain (e.g., standardizing the release pipeline, implementing a database wrapper, or configuring security scopes).

### Avoid Redundant Plan Sprawl
* **Rule:** You MUST NOT create a brand new plan file if the task fits under an existing architectural domain or plan.
* **Collaboration:** Instead of generating a new plan, you MUST *edit* and *adapt* the existing plan, modifying its top-half architectural blueprint and expanding its bottom-half implementation checklist to encompass the new requirements, modifications, or bug fixes.
* **Lifecycle:** Not every Pull Request requires a new plan file. Updating and expanding an existing blueprint is the preferred mode of development.

## Plan Document Structure

Each plan file MUST follow a strict structure consisting of two main halves:

### Top Half: Detailed Architectural Blueprint & Spec
The top half acts as the authoritative technical specification and blueprint for the domain:
* **Executed Date:** The date the plan was fully executed (formatted as `YYYY-MM-DD` for log sorting), or `"pending"` if execution is ongoing.
* **Purpose:** A high-level, clear abstract explaining the domain's goals, architectural intent, and why these changes/structures are established.
* **Specification Details:** This must contain detailed architectural definitions, sequence/swimlane diagrams, core configuration requirements, code snippets, and structural design rules when necessary (similar to `.agent/plans/ReleaseProcess.md`). It functions as a complete spec for any agent or human to read and follow.

### Bottom Half: Sequential Implementation Checklist
The bottom half is the concrete, repository-specific execution sequence derived from the blueprint:
* **Implementation Checklist:** A section named `## Implementation Checklist` containing a sequential, step-by-step implementation checklist.
* **Agent-Built Checklists:** The agent (or human) reads the top-half architectural specification and dynamically builds/expands the step-by-step implementation checklist to accomplish the target repository changes.
* **Dynamic Expansion:** The checklist is a living document. It MUST be dynamically expanded to add new sub-tasks, verification steps, or specific bug fixes that become necessary during active development, testing, or iteration (e.g., when adding specific sub-tasks to resolve bug workarounds in release workflows).
* **Sequential Work Protocol:** Each step of the checklist MUST be worked strictly in turn. You are NOT allowed to skip steps or run steps in parallel if they depend on one another.
* **Update Checklist Progress:** You MUST update the plan file in place, checking off each step (e.g., `- [ ]` -> `- [x]`) **once it is completed and BEFORE starting the next step**. This ensures the plan file serves as a durable, session-persistent execution state.

### Standard Quality Gate Checklist Integration
To guarantee 100% compliance with repository engineering standards, the sequential implementation checklist of EVERY plan MUST explicitly incorporate the quality gates of our Standard Development Process (Phases 3, 4, 5, and 6) as concrete, checkbox-tracked checklist items:
1. **Implementation & Verification:** Standard build, test execution, and local static analysis/linter verification (`golangci-lint`, `go fmt`, `actionlint`, `shellcheck`, etc.).
2. **Proactive Code Review:** Diff validation against review guidelines to guarantee exactly 0 automated comments.
3. **Upstream Sync & Staging Isolation:** Switch to `main`, execute synchronization (`git-sync.sh` or standard git sync), branch off, isolate the current layer's changes, and keep them **unstaged** in the active workspace to maintain color-coded IDE visibility.
4. **Developer IDE Review Gateway:** Invite the developer to perform their IDE review, obtain explicit manual approval, and perform the authorized conventional commit/push.
5. **PR Generation Gateway:** Create the draft Pull Request on GitHub (using `create-pr.sh` or native `gh` commands).
