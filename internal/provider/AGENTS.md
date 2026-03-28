## Resource Development Workflow

When asked to create a new Terraform resource for this provider, you must follow the established conventions found in the `internal/provider` directory. 
The `rancher2_dev2` resource is a perfect template to follow.

Your task is to generate the necessary Go files for a new resource package. A typical resource `rancher2_example_resource` will have the following file structure:

- `internal/provider/rancher2_example_resource/`
- `rancher2_example_resource.go`: The main resource logic (Schema, CRUD functions).
- `rancher2_example_resource_test.go`: Unit tests for the resource's CRUD operations.
- `rancher2_example_resource_model.go`: The Terraform-specific data model (`*ResourceModel` struct with `types.*` fields) and its helper methods (`ToGoModel`, `ToPlan`, `ToState`).
- `rancher2_example_resource_model_test.go`: Unit tests for the resource model conversions.
- `rancher2_example_model.go`: The Go native data model (`*Model` struct with native Go types) and its helper methods (`ToResourceModel`, `ToApiRequestBody`, and various `ToTypesObject` methods for nested object conversions).
- `rancher2_example_model_test.go`: Unit tests for the Go model conversions and API request body construction.

Certain functions MUST exist:

- `rancher2_example_resource.go`
  - `NewRancher2Dev2Resource`
  - `Metadata`
  - `Schema`
  - `Configure`
  - `Create`
  - `Read`
  - `Update`
  - `Delete`
  - `ImportState`
- `rancher2_example_resource_model.go`
  - `ToPlan`
  - `ToState`
  - `ToGoModel`
- `rancher2_example_model.go`
  - `ToApiRequestBody`
  - `ToResourceModel`
  - `ToTypesObject` for each nested struct, used in `ToResourceModel`

**Your process should be:**

1. **Understand the Request:** Clarify the name of the resource (e.g., `rancher2_example_resource`), its API endpoint, and the full schema of attributes (including whether they are required, optional, computed, or sensitive).
2. **Generate Example** Create a new directory in the examples/resources directory for the new resource. Add a resource.tf file to the new directory which outlines the expected outcome for the resource. Add output blocks for all outputs making sure to comment when the outputs are read only. Add all optional arguments and leave a comment specifying that the argument is optional. For mutually exclusive arguments add a commented out argument for one of them specifying that it is mutually exclusive with another argument and which argument it is linked to. Add all arguments which have default values and leave a comment stating that it has a default and what the default is.
3. **Generate Files** Generate the skeleton for the new package using the file structure above.
4. **Implement the Schema** Generate the schema in the resource file along with stud functions for the CRUD operations.
5. **Implement Resource Model** Generate the resource model in the resource model file with the proper tfsdk types.
6. **Implement Go Model** Generate the native Go model in the model file with the proper native Go types.
7. **Implement Basic Unit Tests** Generate tests which validate the expected outputs from all of the required methods listed above. These tests should fail since the functions themselves don't exist. These tests should only validate the most basic functionality.
8. **Implement Required Functions** Implement all functions required above and validate that the tests pass. Use "make test" to run the tests.
9. **Implement Standard Unit Tests** Generate tests which validate standard use cases, such as relying on defaults and standard API errors. The tests should reflect intent, they don't need to pass.
10. **Resolve Standard Unit Test Failures** Resolve any missing functionality in the functions until the tests pass.
11. **Implement Advanced Unit Tests** Generate tests which validate advanced use cases, such as deeply nested structures and difficult API responses. The tests should convey intent, they don't need to pass. Tests should self-document the intended behavior.
12. **Resolve Advanced Unit Test Failures** Resolve any missing functionality in the functions until the tests pass.

By following this structure, you will create new resources that are consistent, maintainable, and testable.

**Principles:**

1. **Metadata Composition:** Metadata is a shared schema which exists in all other resources, it exists in its own package which is imported in resource packages and the functions it provides are used to manage resource metadata.
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
12. **Tests as Documentation** Tests should validate examples, and together they should document the intended behavior of the provider.
13. **Test Driven Development** Tests and examples should be written before behavior because they define the user contract and behavior intent.
