# AI Agent Hooks

This directory contains lifecycle hooks that gate or observe agent behavior at the
process level, so process rules are enforced by the tool rather than relying on the
agent choosing to follow a written instruction. Hooks are only active once registered
in a host's settings file — `.claude/settings.json` for Claude Code, `.gemini/settings.json`
for Gemini. Dropping a script in this directory alone does nothing.

Claude Code and Gemini use different tool-name vocabularies (`Write`/`Edit`/`Bash` vs.
`write_file`/`replace`/`run_shell_command`) and different decision-output schemas for the
same underlying event (`PreToolUse` vs. `BeforeTool`). Where a hook supports both hosts,
it checks for either tool-name vocabulary and emits a JSON object carrying both decision
schemas at once (top-level `decision`/`reason` for Gemini, `hookSpecificOutput.permissionDecision`
for Claude Code) — each host reads only the fields it recognizes.

## Hooks

- **`startup-context.sh`** (`SessionStart`, both hosts) — injects `AGENTS.md` and
  `development-process.md` as additional context at the start of a session.
- **`check-context.sh`** (`PreToolUse`/`BeforeTool`, both hosts) — denies tool calls once
  the session transcript grows past a token-count estimate, prompting a fresh session.
- **`enforce-planning.js`** (`PreToolUse`/`BeforeTool`, both hosts) — denies file writes
  outside `.agent/`, `.gemini/`, `.claude/`, `AGENTS.md`, and `GEMINI.md` unless an active
  plan exists in `.agent/plans/` per `git status`. Also unconditionally denies any write
  to `review-approval.json`, which must only be produced by the review agent.
- **`block-rancher-git.js`** (`PreToolUse`/`BeforeTool`, both hosts) — denies direct
  `git commit`/`git push`, denies any git remote operation referencing `rancher`, denies
  branch switches while the current branch's PR is still a draft, and denies manipulating
  `review-approval.json` via shell commands.
- **`notice-inefficiency.sh`** (`Stop`, Claude Code only) — surfaces token-cost signals
  (cache churn, output-heavy turns, high premium-token turns) for the turn that just
  finished, using the real per-request `usage` data Claude Code already logs to the
  session transcript. Awareness only; it never blocks anything. Thresholds are tunable
  via environment variables documented in the script's header. Gemini is not wired up
  yet — its closest equivalent event (`AfterAgent`) hasn't been verified against a live
  session.

## Adding a new hook

1. Write the script to read the host's stdin JSON and emit the decision schema(s) it
   needs to support — verify field names and event names against a real payload rather
   than assuming they match between hosts (Gemini and Claude Code diverge in exactly the
   ways described above).
2. Register it in `.claude/settings.json` and/or `.gemini/settings.json` under the
   correct event name and matcher. An unregistered hook script does nothing.
3. Document it in this file.
