#!/usr/bin/env bash
#
# Skill: parse-test-logs.sh
# Description: Standard parser for gotestsum / go test -json output files.
#              Extracts and formats passed/failed tests, packages, and elapsed times.
# Usage: .agent/skills/parse-test-logs.sh [options]

set -euo pipefail

# Default options
LOG_FILE=""
NO_COLOR_OPTION=false
FAILED_ONLY=false
PASSED_ONLY=false

show_help() {
  cat <<EOF
Usage: parse-test-logs.sh [options]

Parses gotestsum JSON log files to generate a structured outcome report.

Options:
  -f, --file PATH         Direct path to the JSON log file to parse.
                          If omitted, the newest log file matching /tmp/*_test.log is used.
  --no-color              Disable colored output (respects standard NO_COLOR env var too).
  --failed-only           Only display failed tests and packages.
  --passed-only           Only display passed tests.
  -h, --help              Show this help message and exit.

Examples:
  # Parse the newest test log file
  $ .agent/skills/parse-test-logs.sh

  # Parse a specific log file with color suppressed
  $ .agent/skills/parse-test-logs.sh -f /tmp/my_test.log --no-color

  # Show failures only
  $ .agent/skills/parse-test-logs.sh --failed-only
EOF
}

# Parse command line options
while [[ $# -gt 0 ]]; do
  case "$1" in
    -f|--file)
      if [[ $# -lt 2 ]]; then
        echo "Error: Option $1 requires an argument." >&2
        show_help >&2
        exit 1
      fi
      LOG_FILE="$2"
      shift 2
      ;;
    --no-color)
      NO_COLOR_OPTION=true
      shift
      ;;
    --failed-only)
      FAILED_ONLY=true
      shift
      ;;
    --passed-only)
      PASSED_ONLY=true
      shift
      ;;
    -h|--help)
      show_help
      exit 0
      ;;
    *)
      echo "Error: Unknown option $1" >&2
      show_help >&2
      exit 1
      ;;
  esac
done

# Validate mutual exclusion
if [ "$FAILED_ONLY" = true ] && [ "$PASSED_ONLY" = true ]; then
  echo "Error: --failed-only and --passed-only options are mutually exclusive." >&2
  show_help >&2
  exit 1
fi

# Initialize color settings
if [ "$NO_COLOR_OPTION" = true ] || [ -n "${NO_COLOR:-}" ]; then
  RED=""
  GREEN=""
  YELLOW=""
  BLUE=""
  NC=""
  CHECK_MARK="✓"
  CROSS_MARK="✗"
else
  RED='\033[1;31m'
  GREEN='\033[1;32m'
  YELLOW='\033[1;33m'
  BLUE='\033[1;34m'
  NC='\033[0m'
  CHECK_MARK='\033[1;32m✓\033[0m'
  CROSS_MARK='\033[1;31m✗\033[0m'
fi

# Locate log file if not specified
if [ -z "$LOG_FILE" ]; then
  # Find newest file matching /tmp/*_test.log
  # Avoid using ls in pipe where possible, but here we find newest by time
  LOG_FILE=$(find /tmp -maxdepth 1 -name "*_test.log" -type f -printf "%T@ %p\n" 2>/dev/null | sort -n | tail -n 1 | awk '{print $2}')
fi

if [ -z "$LOG_FILE" ] || [ ! -f "$LOG_FILE" ]; then
  echo -e "${RED}[ERROR]${NC} No valid test log file found or specified." >&2
  exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
  echo -e "${RED}[ERROR]${NC} jq command is required but not installed." >&2
  exit 1
fi

echo -e "${BLUE}=== Parsing Test Log: $LOG_FILE ===${NC}"

# Extract details from JSON
passed_tests=$(jq -r '. | select(.Action == "pass") | select(.Test != null) | "\(.Test)\t\(.Elapsed)"' "$LOG_FILE" 2>/dev/null | sort -u || true)
failed_tests=$(jq -r '. | select(.Action == "fail") | select(.Test != null) | "\(.Test)\t\(.Elapsed)"' "$LOG_FILE" 2>/dev/null | sort -u || true)
failed_pkgs=$(jq -r '. | select(.Action == "fail") | select(.Test == null) | "\(.Package)\t\(.Elapsed)"' "$LOG_FILE" 2>/dev/null | sort -u || true)

# Count summaries
passed_count=0
if [ -n "$passed_tests" ]; then
  passed_count=$(echo "$passed_tests" | wc -l)
fi

failed_test_count=0
if [ -n "$failed_tests" ]; then
  failed_test_count=$(echo "$failed_tests" | wc -l)
fi

failed_pkg_count=0
if [ -n "$failed_pkgs" ]; then
  failed_pkg_count=$(echo "$failed_pkgs" | wc -l)
fi

total_failures=$((failed_test_count + failed_pkg_count))

# Print PASSED section
if [ "$FAILED_ONLY" = false ]; then
  echo ""
  echo -e "${GREEN}PASSED TESTS ($passed_count):${NC}"
  if [ -n "$passed_tests" ]; then
    while IFS=$'\t' read -r name elapsed; do
      if [ -n "$name" ]; then
        printf "  %b %s (%ss)\n" "${CHECK_MARK}" "$name" "$elapsed"
      fi
    done <<< "$passed_tests"
  else
    echo "  None"
  fi
fi

# Print FAILED section
if [ "$PASSED_ONLY" = false ]; then
  echo ""
  echo -e "${RED}FAILED ITEMS ($total_failures):${NC}"
  
  if [ "$failed_test_count" -gt 0 ]; then
    echo -e "  ${YELLOW}Individual Tests:${NC}"
    while IFS=$'\t' read -r name elapsed; do
      if [ -n "$name" ]; then
        printf "    %b %s (%ss)\n" "${CROSS_MARK}" "$name" "$elapsed"
      fi
    done <<< "$failed_tests"
  fi

  if [ "$failed_pkg_count" -gt 0 ]; then
    echo -e "  ${YELLOW}Packages/Builds:${NC}"
    while IFS=$'\t' read -r name elapsed; do
      if [ -n "$name" ]; then
        printf "    %b Package: %s (%ss)\n" "${CROSS_MARK}" "$name" "$elapsed"
      fi
    done <<< "$failed_pkgs"
  fi

  if [ "$total_failures" -eq 0 ]; then
    echo "  None"
  fi
fi

echo ""
echo -e "${BLUE}=======================================${NC}"
if [ "$total_failures" -gt 0 ]; then
  echo -e "${RED}Overall Outcome: FAILED${NC}"
  exit 1
else
  echo -e "${GREEN}Overall Outcome: SUCCESS${NC}"
  exit 0
fi
