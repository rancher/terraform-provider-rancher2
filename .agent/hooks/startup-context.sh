#!/usr/bin/env bash
set -euo pipefail

# Consume and discard hook input from stdin to prevent broken pipes
cat > /dev/null

# Log diagnostics to stderr to comply with the silence rule on stdout
echo "Loading session-start workspace context..." >&2

combined_context=""
combined_context+=$'###############################################################################\n'
combined_context+=$'#                           CRITICAL AGENT MANDATES                           #\n'
combined_context+=$'#                                                                             #\n'
combined_context+=$'# 1. YOU MUST FOLLOW THE DEVELOPMENT PROCESS IN \'development-process.md\'.     #\n'
combined_context+=$'# 2. YOU MUST NEVER COMMIT OR PUSH DIRECTLY. YOU MUST ALWAYS USE THE CUSTOM   #\n'
combined_context+=$'#    COMMIT-PUSH SKILL: \'.agent/skills/commit-push.sh -m "message"\'.          #\n'
combined_context+=$'# 3. FOR ALL TASKS, YOU MUST DEFINE A SEQUENTIAL IMPLEMENTATION CHECKLIST AT  #\n'
combined_context+=$'#    THE BOTTOM OF YOUR PLAN IN \'.agent/plans/<PlanName>.md\'.                 #\n'
combined_context+=$'#                                                                             #\n'
combined_context+=$'# FAILURE TO COMPLY WILL TRIGGER SECURITY BLOCKS AND PROCESS TERMINATION.     #\n'
combined_context+=$'###############################################################################\n\n'

if [[ -f "AGENTS.md" ]]; then
  combined_context+=$'# Context from AGENTS.md\n\n'
  combined_context+=$(cat AGENTS.md)
  combined_context+=$'\n\n'
  echo "Loaded AGENTS.md" >&2
else
  echo "Warning: AGENTS.md not found" >&2
fi

if [[ -f ".agent/workflows/development-process.md" ]]; then
  combined_context+=$'# Context from .agent/workflows/development-process.md\n\n'
  combined_context+=$(cat .agent/workflows/development-process.md)
  combined_context+=$'\n\n'
  echo "Loaded .agent/workflows/development-process.md" >&2
else
  echo "Warning: .agent/workflows/development-process.md not found" >&2
fi

# Output clean JSON structure to stdout
jq -n --arg ctx "$combined_context" '{
  "hookSpecificOutput": {
    "additionalContext": $ctx
  },
  "systemMessage": "✨ AGENTS.md and development-process.md context injected successfully."
}'
