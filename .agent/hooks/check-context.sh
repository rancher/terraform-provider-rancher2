#!/usr/bin/env bash
set -euo pipefail

# Claude Context Limiter Hook
# Reads JSON from stdin, outputs hookSpecificOutput JSON on deny, exits 0 on allow.

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
  
  # Return deny decision via JSON for Claude's PreToolUse event
  jq -n --arg reason "$reason" '{
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
