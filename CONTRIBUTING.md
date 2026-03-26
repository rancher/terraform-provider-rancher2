### How to Add/Update a Resource

1. **Context:** Provide the LLM with `internal/prompts/system_prompt.md`.
2. **Input:** Provide the API/CRD JSON or YAML schema.
3. **Prompt:** *"Using the System Prompt, generate the `rancher2_[NAME]` resource for this schema: [SCHEMA]"*
4. **Audit:** Ask the AI: *"Verify this follows the Sticky Metadata and Finalizer Override patterns."*
5. **Deploy:** Register the resource in `provider.go` and run Acceptance Tests.

