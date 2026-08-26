# Developer Architecture & Blueprints Directory (`docs/development/`)

Welcome to the central developer documentation and architecture repository. This directory houses the standing, long-term technical specifications, architectural designs, and domain-specific workflows that govern our development practices and repository-wide automation.

---

## 📖 Key Architectural Concepts

Our engineering culture emphasizes **explicit design**, **cryptographic security**, and **deterministic validation**. To support these principles, all major codebase modifications, developer workflows, and automation components must be thoroughly documented as standing blueprints.

By separating our architectural blueprints (`docs/development/`) from our shared automation tooling boilerplate (`.gemini/`), we achieve:

1. **Generic Portability**: The entire `.gemini/` automation suite contains 0 repository-specific hardcoded values, allowing it to be synced seamlessly across other Rancher and SUSE repositories as a global standard.
2. **Permanent System Memory**: We maintain a comprehensive, standing record of all architectural decisions and design specifications that developers and agents alike can consult, analyze, and expand.

---

## 📂 Blueprint Directory Structure & Nomenclature

To scale our documentation as the project grows, we organize our files into a structured, two-tiered nomenclature pattern:

### 1. Topic Overview (`docs/development/<Topic>.md`)

A **Topic Overview** serves as a high-level domain abstract, indexing the components, files, and architectural goals of a specific domain (such as `AgenticFramework.md` or `ReleaseProcess.md`).

- **Required Sections**:
  - **Purpose**: High-level domain abstract.
  - **Modular Architectural Blueprints Map**: Interactive links to the individual Component Specifications in this domain.

### 2. Component Specification (`docs/development/<Topic>/<Component>.md`)

A **Component Specification** is a detailed, focused technical design document that covers a specific subsystem, utility, or workflow.

- **Required Sections**:
  - **Abstract**: Concise explanation of the component's goal and responsibilities.
  - **Architectural Design / Technical Specification**: Deep dive into the component's mechanics, data flows, and configuration.
  - **Standing Implementation Decisions**: Long-term design choices, trade-offs, and standing parameters that must be preserved.

---

## 🛠️ Domain Index Maps

- **[Agentic Framework](./AgenticFramework.md)**: Secure zero-trust sandboxing, Touch ID gating, PR pipeline automation, and agent developer workflows.
- **[Coding Standards](./CodingStandards.md)**: Repository-wide standards for Go, Terraform, Shell/GitHub scripts, Workflows, and pre-commit reviews.
- **[Release Process](./ReleaseProcess.md)**: Canonical tag, GPG-signing, and release workflows.
- **[Testing](./Testing.md)**: Acceptance and unit-testing standards.
- **[Documentation Changelog](./CHANGELOG.md)**: Revision history of developer blueprints and specifications.
