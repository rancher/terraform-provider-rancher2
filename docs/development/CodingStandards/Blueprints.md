# Architectural Blueprints & Documentation

---

## Abstract

These guidelines describe the layout, structural patterns, and planning procedures used to document repository features, workflows, and automated checks under the `docs/development/` directory.

---

## Nomenclature & Structure

The repository documentation utilizes two main types of declarative architectural specifications, collectively referred to as **Blueprints**:

### 1. Topic Overviews (`docs/development/<Topic>.md`)

An overarching domain document that provides a high-level understanding of a major system area (such as `ReleaseProcess`, `Testing`, or `AgenticFramework`).

- **Abstract/Introduction**: Describes the general concepts of the topic, its purpose, and architectural intent.
- **Architectural Design**: Contains design descriptions, diagrams, or schemas explaining how the system's components fit together.
- **Component Index**: A list of all sub-components that comprise this topic, linking to their respective Component Specifications.

### 2. Component Specifications (`docs/development/<Topic>/<Component>.md`)

A technical design document and actionable specification for a single sub-component under a topic (such as `AgenticFramework/GatingAndApprovals.md`).

- **Abstract**: A clear, technical abstract section named `## Abstract` explaining the component's goals and architectural intent.
- **Specification Details**: Detailed structural design rules, configuration requirements, sequence diagrams, and code patterns describing the system as it currently exists.

---

## Blueprints vs. Plans

To maintain a clean and reliable automated workflow, the repository distinguishes clearly between static architectural specifications and runtime agentic instructions:

- **Architectural Blueprints**: Represent long-lived, declarative documentation files located under the `docs/development/` directory. They describe the system's current state, structure, and design constraints. Blueprints are not modified as part of normal agentic turn cycles unless the overall architecture of a component changes.
- **Imperative Plans**: Represent the step-by-step execution checklists generated during the session's **Plan Mode** (such as files written to `plans/` inside the session workspace). A plan represents the dynamic, imperative contract between the user and the agent for a specific task. To exit Plan Mode and enable code modifications, a signed imperative plan must exist and be cryptographically approved by the user. A Blueprint is not required to exit Plan Mode.

---

## Documentation Life Cycle

- **Avoiding Blueprint Sprawl**: New Topic Overviews or Component Specifications are only created when a new architectural domain is introduced. For changes to existing systems, the corresponding Blueprints are adapted and updated directly to reflect the new state.
- **Declarative Realignment**: Documentation is kept strictly declarative, describing the systems and configurations as they are currently implemented. Historical change logs, implementation checklists, and past milestone histories are omitted to ensure the documentation remains focused on the present architecture.
