#!/usr/bin/env bash
#
# Claude Code Stop hook: surfaces token-cost signals for the turn that just finished,
# so a developer can notice expensive patterns in the moment instead of being surprised
# by cumulative spend later. Awareness only — this never blocks or denies anything.
#
# Reads the real per-request `usage` blocks (input/output/cache tokens) that Claude Code
# already writes to the session transcript, and sums them across every model call made
# since the last real user prompt (a "turn" may include several calls if tools were used
# in between). It deliberately does not look at total context size: a large cache_read is
# cheap, reused content, not a cost problem, and conflating "big context" with "expensive
# turn" produces false alarms on long, healthy, well-cached sessions.
#
# Three signals, chosen because they map directly to Anthropic's pricing model rather
# than to a proxy like context size or elapsed time:
#   - cache churn:    cache_creation tokens are billed at a premium (~1.25x input) and
#                      only pay off if reused later via cheap cache_read. A high
#                      cache_creation/(cache_creation+cache_read) ratio means the cached
#                      prefix is being invalidated instead of reused (e.g. a large one-off
#                      tool result sitting where it breaks the cache boundary every turn).
#   - output-heavy:   output tokens cost several times more per token than input.
#   - premium tokens: cache_creation + input (i.e. everything NOT served from cheap
#                      cache_read) summed across every call this turn — the closest
#                      token-only proxy to "what this turn actually cost" without a
#                      hardcoded, easily-stale $-per-token price table.
#
# Thresholds are tunable via environment variables (see defaults below) since "expensive"
# is relative to what a task actually needs.
#
# Gemini's hook system does not currently have a verified equivalent post-turn event (its
# taxonomy has BeforeTool/AfterTool/AfterAgent, but AfterAgent's schema hasn't been tested
# against a live session the way this has), so this hook is Claude Code-only for now and
# is registered solely in .claude/settings.json.

set -euo pipefail

payload=$(cat)
[[ -z "$payload" ]] && exit 0

transcript_path=$(jq -r '.transcript_path // empty' <<<"$payload" 2>/dev/null || true)
[[ -z "$transcript_path" || ! -f "$transcript_path" ]] && exit 0

# Bound how far back we scan. A single turn (even a long tool-calling one) rarely spans
# more than a few hundred transcript lines; this keeps the hook fast on multi-MB sessions.
WINDOW="${EFFICIENCY_SCAN_WINDOW:-500}"

CACHE_WRITE_RATIO_THRESHOLD="${EFFICIENCY_CACHE_WRITE_RATIO_THRESHOLD:-0.15}"
CACHE_WRITE_MIN_TOKENS="${EFFICIENCY_CACHE_WRITE_MIN_TOKENS:-3000}"
OUTPUT_TOKEN_THRESHOLD="${EFFICIENCY_OUTPUT_TOKEN_THRESHOLD:-8000}"
PREMIUM_TOKEN_THRESHOLD="${EFFICIENCY_PREMIUM_TOKEN_THRESHOLD:-15000}"

tail -n "$WINDOW" "$transcript_path" | jq -cs \
  --argjson cache_ratio_threshold "$CACHE_WRITE_RATIO_THRESHOLD" \
  --argjson cache_write_min "$CACHE_WRITE_MIN_TOKENS" \
  --argjson output_threshold "$OUTPUT_TOKEN_THRESHOLD" \
  --argjson premium_threshold "$PREMIUM_TOKEN_THRESHOLD" \
'
  def is_real_user_prompt:
    .type == "user" and (
      (.message.content | type) == "string" or
      ((.message.content | type) == "array" and ((.message.content | map(.type=="tool_result") | any) | not))
    );
  to_entries
  | (map(select(.value | is_real_user_prompt)) | last | .key // 0) as $start
  | .[$start:]
  | map(select(.value.type == "assistant"))
  | map(.value.message.usage) as $u
  | {
      turns: ($u | length),
      output: ($u | map(.output_tokens // 0) | add // 0),
      cache_write: ($u | map(.cache_creation_input_tokens // 0) | add // 0),
      cache_read: ($u | map(.cache_read_input_tokens // 0) | add // 0),
      input: ($u | map(.input_tokens // 0) | add // 0)
    }
  | . + {
      cache_write_ratio: (if (.cache_write + .cache_read) > 0 then (.cache_write / (.cache_write + .cache_read)) else 0 end),
      premium_tokens: (.cache_write + .input)
    }
  | . as $s
  | [
      (if ($s.cache_write > $cache_write_min) and ($s.cache_write_ratio > $cache_ratio_threshold)
       then "Cache churn: \($s.cache_write) tokens were freshly cached this turn (\(($s.cache_write_ratio * 100) | floor)% of cache-write+cache-read), suggesting the cached prefix is being invalidated rather than reused. Check for a large tool result sitting where it breaks the cache boundary."
       else empty end),
      (if $s.output > $output_threshold
       then "Output-heavy: this turn generated \($s.output) output tokens, which cost several times more per token than input. Consider whether the task needed that much generated text."
       else empty end),
      (if $s.premium_tokens > $premium_threshold
       then "High premium-token turn: \($s.premium_tokens) tokens were billed outside of cheap cache reads (cache-write + non-cached input) across \($s.turns) model call(s). This is the closest token-only signal to real cost for this turn."
       else empty end)
    ] as $notes
  | if ($s.turns == 0) or ($notes | length == 0) then empty
    else { systemMessage: ("⚠️ Efficiency notice for this turn:\n" + ($notes | map("- " + .) | join("\n"))) }
    end
' 2>/dev/null || true

exit 0
