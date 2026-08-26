# Agentic Framework: Cryptographic Gating & Approvals

## Abstract

To maintain absolute system integrity and prevent unauthorized code modifications in an autonomous programming workspace, the Agentic Framework implements a secure, **cryptographically chained gating pipeline**. This system mathematically guarantees that no code can be committed or pushed without satisfying sequential, hardware-authorized checkpoints: Planning (Gate 1), Quality Gate (Gate 2), and Commit Gate (Gate 3).

---

## 🔒 Architectural Rationale: Why We Cryptographically Chain Gates

Traditional CI/CD pipelines and pre-commit checks are highly vulnerable to **TOCTOU (Time-of-Check to Time-of-Use)** attacks and race conditions. An autonomous AI agent could theoretically modify files _after_ passing local tests and reviews, but _before_ executing the final commit, thereby injecting unvetted code into the repository.

To eliminate this vulnerability, our framework enforces a **Strict Cryptographic Chaining Rule**:

1. **The Hash Anchor**: Every stage of the development process is anchored to a unique, live SHA-256 hash of the entire active workspace difference (`git diff HEAD`), termed the `diff_hash`.
2. **Immutable Binding**: When a gating step succeeds, its signature (`plan-approval.json` or `review-approval.json`) is cryptographically bound to that exact `diff_hash` and the active `plan_hash`.
3. **Chained Verification**: Subsequent gates—and ultimately the custom `agent-scripts/commit-push-helper.sh` utility—will unconditionally reject operations if any of the following occur:
   - The files on disk change (which changes the current `diff_hash` and causes an immediate mismatch with the signatures).
   - Any signature is missing or altered.
- The signatures are generated out of order (e.g., trying to write a Review approval before a Test approval is signed).
4. **Self-Healing Revocation**: If any check fails or if the agent modifies any source file, the enforcer hooks automatically unlink (delete) the downstream signatures, immediately revoking approvals and halting the pipeline.

---

## 🧬 SSH Key-Backed Gating & Threat Model

In a secure, semi-autonomous engineering environment, the ultimate threat is **Agent Autonomy Escalation**: an AI agent fabricating or signing off on its own changes without real developer oversight.

To prevent this, the framework binds the **Planning Gate (Gate 1)** and the final **Commit Gate (Gate 3)** to a cryptographic key. We use standard SSH key material (located at `~/.gemini/ssh-key` and `~/.gemini/ssh-key.pub`) to sign approvals. This requires human physical validation or authentication (such as Touch ID or SSH Agent key unlocking) to cryptographically sign approvals, ensuring the developer remains the absolute authority over the codebase.

### Cryptographic Decryption Flow

```text
[Gemini CLI]                         [System OS (Enforcer Hook)]                  [SSH Agent / SSH Key]
     |                                           |                                           |
     | -- 1. Trigger exit_plan_mode --------->   |                                           |
     |                                           | -- 2. Generate Challenge token -------->  |
     |                                           |                                           |
     | <--- 3. Prompt: "Approve?" -------------- |                                           |
     |                                           |                                           |
     | -- 4. Click: [Approve Plan] ------------> |                                           |
     |                                           | -- 5. Execute ssh-keygen -Y sign ------>  |
     |                                           |                                           | [SSH Unlock / Touch ID]
     |                                           |                                           | <--- Confirm Signature
     |                                           | <--- 6. Cryptographic Signature --------  |
     |                                           |                                           |
     | <--- 7. Write plan-approval.json.sig ---- |                                           |
```

1. **Challenge Generation**: The enforcer hook (`04-commit-phase.js --after-ask` / `02-plan-phase.js --ask-proof`) intercepts the approved `ask_user` tool execution. It generates a secure, randomized challenge token and writes the payload (`plan-approval.json` or `user-approval.json`).
2. **SSH Key Signing**: The hook invokes `ssh-keygen -Y sign -f ~/.gemini/ssh-key -n gemini` to sign the payload.
3. **Optional Biometric / Agent Unlock**: If the private key is configured to use a passphrase, or if it is integrated with the macOS Secure Enclave/Touch ID (via macOS SSH agent integration), the system natively prompts the developer for authorization (biometrics or passphrase entry).
4. **Signature Verification**: The enforcer hooks verify the signature file using `ssh-keygen -Y verify -f allowed_signers` against the public key `~/.gemini/ssh-key.pub`.

---

## 🛠️ Step-by-Step SSH Key Setup Guide

The framework uses standard SSH keys and `ssh-keygen` to coordinate cryptographic signatures. Both packages are natively available in Darwin, Linux, and our hermetic Nix development shell.

### 1. Generate a Standard SSH Key Pair for Gemini

If you don't already have a dedicated SSH key for Gemini, generate a standard SSH key pair (without a passphrase, or managed via an SSH agent):

```bash
# Generate the key pair and save to the standardized path
ssh-keygen -t ed25519 -f ~/.gemini/ssh-key -N ""
```

### 2. Verify Key Permissions

Ensure the private key file has secure, highly restrictive permissions:

```bash
chmod 600 ~/.gemini/ssh-key
chmod 644 ~/.gemini/ssh-key.pub
```

### 3. (Optional) Back SSH Keys with Apple Secure Enclave (Touch ID / macOS / Darwin)

Alternatively, you can back your local SSH keys using macOS Touch ID to enforce physical biometric confirmation on every signature call. Please refer to standard macOS secure enclave SSH integration guides (e.g. using Touch ID-backed security keys or macOS agent forwarding configurations).

---

## 🤖 Programmatic Subagent Isolation & AfterTool Hook Integration

While Gates 1 and 3 require human SSH Key verification, **Gate 2 (Proactive Review)** is programmatically audited and signed by invoking a specialized local subagent (`review_agent`). Automated tests are enforced programmatically in the pre-review hooks, eliminating the need for a standalone testing subagent.

To guarantee absolute objectivity, the subagent is **fully isolated and sandboxed**:

1. **Read-Only Enforcements**: The `review_agent` is stripped of write capabilities, restricting its toolset strictly to `[read_file]`. It cannot modify code or write approvals.
2. **AfterTool Hook Verification**:
   - When the main agent calls `invoke_agent` targeting `review_agent`, the native `AfterTool` hook (`03-review-phase.js --after-invoke`) intercepts the subagent's execution report.
   - It programmatically parses the report to verify that:
     1. All 4 sequential passes (`Pass 1`, `Pass 2`, `Pass 3`, and `Pass 4`) are checked as complete checklist items (e.g. `- [x] Pass 1`).
     2. Exactly 0 findings are reported, verified by the presence of the clean marker `0 comments/findings` or `0 findings`.
   - If the report is successfully verified as complete and clean, the hook **natively and securely** writes `review-approval.json` to disk, signing it with the active `diff_hash` and `plan_hash`.
   - If any pass is unchecked or if findings are recorded, the hook instantly unlinks (deletes) the signatures, revoking any previous approvals.

This architecture prevents the main agent or subagent from manually writing approvals or manipulating results, enforcing a completely deterministic quality gateway.

---

## 🤖 Hook Timeout & Debugging Configuration

To ensure reliable execution during remote Git pushes and GitHub Pull Request creation, the `commit-phase-after-ask-user` hook (`04-commit-phase.js --after-ask`) is configured with a generous timeout of `60000` ms (60 seconds) inside `.gemini/settings.json`.

The hook isolates the automated commit, push, and PR generation steps into distinct, targeted try-catch blocks. This prevents remote connection timeouts or Git authentication errors from being incorrectly reported as SSH signature failures, ensuring clear, high-signal debugging logs during execution.

---

## 🤖 Modular Controller & Function Decoupling (Agent-Scripts)

The agentic framework decouples core agent logic (such as signing, verification, and Git helper scripts) from the Gemini-specific CLI integration layer. This ensures code reuse, simplifies testing, and supports multi-agent workspaces.

### 1. Root-Level `agent-scripts/` Directory

Generic, reusable agent scripts—both ESM Node.js modules and POSIX Shell helpers—are organized inside the flat root-level `agent-scripts/` directory. These files do not use Gemini-specific CLI primitives, representing pure agent workflows and operations.

### 2. Controller/Shim Architecture

Files in `.gemini/hooks/` and `.gemini/skills/` function as thin **controllers**. Their responsibilities are limited to:

- Parsing Gemini-specific tool execution payloads and arguments.
- Invoking the underlying modular scripts under `agent-scripts/` via clean function calls or command execution.
- Formatting and returning standard output protocols to the Gemini CLI.

### 3. Decoupled Subsystems

The core enforcer and automation layers are decoupled into the following dedicated modules:

- **Planning Logic (`agent-scripts/planning.js`)**: Manages blueprint checks and active plan validations.
- **Cryptographic Gating (`agent-scripts/gating.js`)**: Computes workspace `diff_hash` values and verifies Plan, Test, and Review signatures.
- **Sub-Agent Hook Pipeline (`agent-scripts/after-invoke.js`)**: Evaluates review and testing reports upon tool completion, programmatically signing or revoking gate approvals.
- **Cryptographic Signatures (`agent-scripts/after-ask.js`)**: Conducts SSH signature challenges (via `ssh-keygen -Y sign`) for the Plan and Commit gates.
- **Execution Security (`agent-scripts/security.js`)**: Validates shell commands and paths to prevent command injection or directory bypass.
- **Unified Git Operations (`agent-scripts/git-helpers.js`)**: Consolidates branch switching, remote branch ancestry checks, Conventional commits, and automated PR generation.
- **Automated Quality Verification**: Granular mock-scaffolded unit tests in the `agent-scripts/tests/` directory validate the cryptographic consistency and execution safety of all decoupled helper scripts.
