#!/usr/bin/env bash
set -euo pipefail

# Context Limiter Hook (PreToolUse on Claude Code, BeforeTool on Gemini).
# Reads JSON from stdin, exits 0 on allow. On deny, emits a JSON object with both
# top-level decision/reason (Gemini) and hookSpecificOutput.permissionDecision
# (Claude Code) — each host reads only the fields it recognizes.

# Read standard input into a variable
payload=$(cat)

# Return allow immediately if payload is empty to prevent jq parsing errors
if [[ -z "$payload" ]]; then
  exit 0
fi

token_usage=$(printf "%s" "$payload" | jq -r '.tokens // 0')
transcript_path=$(printf "%s" "$payload" | jq -r '.transcript_path // ""')

MAX_TOKENS=200000

# Approximate Claude tokens from transcript if .tokens not provided
if [[ "$token_usage" -eq 0 ]] && [[ -n "$transcript_path" ]] && [[ "$transcript_path" != "null" ]] && [[ -f "$transcript_path" ]]; then
  file_size=$(wc -c < "$transcript_path" | tr -d ' ')
  token_usage=$((file_size / 4))
fi

if [[ "$token_usage" -gt "$MAX_TOKENS" ]]; then
  reason="Context limit of $MAX_TOKENS tokens reached. You must halt operations, update plans, and prompt the user for a new session."

  jq -n --arg reason "$reason" '{
    decision: "deny",
    reason: $reason,
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "deny",
      permissionDecisionReason: $reason
    }
  }'
  exit 0
fi

# Exit 0 with empty stdout to allow normal flow
exit 0
