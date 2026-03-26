## Environment-Specific Instructions

Dependencies for this project are provided by Nix, use this command to run scripts with the dependencies installed: `nix develop --ignore-environment --extra-experimental-features nix-command --extra-experimental-features flakes --keep HOME --keep SSH_AUTH_SOCK --keep GPG_SIGNING_KEY --keep NIX_SSL_CERT_FILE --keep NIX_ENV_LOADED --keep TERM --command bash -e {0}`

The GNUmakefile present in this repository gives a few convienience methods for common necessities:
- use `make build` to build the project and check for compilation errors
- use `make test` to run unit tests

You may read files and folders in this repository without asking.

## Repository Coding Standards & Instructions

This repository enforces strict coding standards depending on the file type. Whenever you are asked to generate, edit, or review code, you MUST consult the corresponding instruction file for the specific rules, anti-patterns, and requirements:

* **For Go (`**/*.go`)**: Read and strictly adhere to `.github/instructions/go.instructions.md`
* **For Terraform (`**/*.tf`)**: Read and strictly adhere to `.github/instructions/terraform.instructions.md`
* **For GitHub Actions (`.github/workflows/**/*.{yml,yaml}`)**: Read and strictly adhere to `.github/instructions/workflows.instructions.md`
* **For GitHub Scripts (`.github/workflows/scripts/**/*.js`)**: Read and strictly adhere to `.github/instructions/github-script.instructions.md`
* **For Spelling Changes (`aspell_custom.txt`)**: Read and strictly adhere to `.github/instructions/aspell.instructions.md`

## Resource Development Workflow

When asked to create a new Terraform resource for this provider, you must follow the established conventions found in the `internal/provider` directory. The `rancher2_dev` resource is a perfect template to follow.

Your task is to generate the necessary Go files for a new resource package. A typical resource `rancher2_example_resource` will have the following file structure:

- `internal/provider/rancher2_example_resource/`
  - `rancher2_example_resource.go`: The main resource logic (Schema, CRUD functions).
  - `rancher2_example_resource_model.go`: The Terraform-specific data model (`*ResourceModel` struct with `types.*` fields) and its helper methods (`ToGoModel`, `ToPlan`, `ToState`).
  - `rancher2_example_model.go`: The API-specific data model (`*Model` struct with native Go types) and its helper method (`ToResourceModel`).
  - `rancher2_example_resource_test.go`: Unit tests for the resource's CRUD operations.
  - `rancher2_example_resource_model_test.go`: Unit tests for the resource model conversions.
  - `rancher2_example_model_test.go`: Unit tests for the API model conversions.

**Your process should be:**

1.  **Understand the Request:** Clarify the name of the resource (e.g., `rancher2_example_resource`), its API endpoint, and the full schema of attributes (including whether they are required, optional, computed, or sensitive).
2.  **Generate the Models:**
    *   Create the `..._model.go` file with the API model struct (`ExampleModel`). Use native Go types and `json` tags that match the Rancher API payload.
    *   Create the `..._resource_model.go` file with the Terraform model struct (`ExampleResourceModel`). Use `types.*` fields corresponding to the API model.
3.  **Implement Model Conversions:**
    *   Implement the `ToGoModel()` method on `*ExampleResourceModel` to convert from Terraform types to Go types. This is for preparing the API request.
    *   Implement the `ToResourceModel()` method on `*ExampleModel` to convert from Go types to Terraform types. This is for processing the API response.
4.  **Implement the Resource:**
    *   In `..._resource.go`, define the `Schema` using the attributes identified in step 1.
    *   Implement the `Create`, `Read`, `Update`, and `Delete` methods. Use the injected `client` to make API calls.
    *   Use the model conversion functions to translate data between the Terraform plan/state and the API request/response bodies.
5.  **Generate Tests:**
    *   Create comprehensive unit tests for all new files. Mock the API client and responses to test the CRUD functions under various conditions (success, failure, conflict, etc.).
    *   Write tests for the model conversion functions to ensure data integrity.

By following this structure, you will create new resources that are consistent, maintainable, and testable.

**Principles:**
1. **Metadata Composition:** Use shared `Metadata` blocks; do not redefine `uid`, `labels`, etc.
2. **Sticky Logic:** Use `PlanModifiers` on `labels`/`annotations`. HCL keys reconcile; API-only keys are adopted.
3. **OwnerRefs:** Match by `uid`. Merge API-side refs into the Plan; do not delete them.
4. **Finalizer Override:** Merge by default. If HCL is an explicit empty list `[]`, overwrite State to clear blockers.
5. **Read-Only Attributes:** `uid`, `resource_version`, `creation_timestamp` are `Computed` & `ReadOnly`.
6. **Read/Map Logic:** `Read` must handle 404s (RemoveResource). `MapFromApi` must preserve adopted/sticky keys in State.
7. **Provider-Owned ID:** The top-level `id` attribute is reserved for Terraform's state management and MUST be a provider-generated UUID. Store the API's unique identifier in a separate, `Computed` attribute (e.g., `cluster_id`, `token_id`).
8. **Sanitize API Requests:** Before `Create` or `Update` calls, construct the request payload from the plan but explicitly nullify all `Computed`-only attributes. This prevents the API from rejecting requests containing read-only, server-generated fields.
9. **Idempotent Deletion:** The `Delete` function MUST treat a `404 Not Found` API response as a success, ensuring that destroying an already-deleted resource does not fail.
10. **Type Conversion Models** Resources MUST be convertable to tfsdk.plan, tfsdk.state, native Go model, and Terraform Resource model. Models need to accomodate this requirement, for instance Resource models need to use types.Object for nested structures rather than pointers to another struct. Each package MUST include functions to facilitate the conversions such as `ToGoModel`, `ToResourceModel`, `ToPlan`, and `ToState`.
11. **Thin Controller Model** Resources MUST reveal the API structure they are managing, they MUST NOT use flatteners and expanders to change the structure of the data before sending calls to the API.
