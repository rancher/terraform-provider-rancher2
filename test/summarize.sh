#!/bin/sh
set -e

# summarize.sh - reads report.json and prints a summary

REPORT_FILE="${1:-report.json}"

# Ensure jq is installed
if ! command -v jq >/dev/null; then
  echo "Error: 'jq' is not installed. Please install it to generate the summary."
  exit 1
fi

if [ ! -f "$REPORT_FILE" ]; then
  echo "Error: '$REPORT_FILE' not found." >&2
  exit 1
fi

if [ ! -s "$REPORT_FILE" ]; then
  echo "Warning: Test report '$REPORT_FILE' is empty. No test results to summarize." >&2
  if [ "$REPORT_FILE" = "report.json" ] || [ "$REPORT_FILE" = "./report.json" ]; then
    rm -f "$REPORT_FILE"
  fi
  exit 0
fi

process_tests() {
  action="$1"
  file="${2:-$REPORT_FILE}"
  # slurp is important here, it reads the objects into an array for further processing
  # select all objects with a .Test listed that matches $action and store it as $tests
  # then iterate through tests and save the test name as $prefix
  # then filter out any test names that don't exactly match the current test name, but do have the test name in them followed by a slash
  # this leaves only test names that don't have duplicate prefixes, ie. the actual tests and not the parent blocks
  # eg.
  # changes this:
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults/with_all_fields
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults/with_some_fields
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults/with_no_fields
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults/with_nil_value
  # into this:
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults/with_all_fields
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults/with_some_fields
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults/with_no_fields
  # TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults/with_nil_value
  jq -R -s -r --arg action "$action" \
    '
    split("\n") |
    map(select(length > 0) | fromjson?) |
    (map(select(.Test and .Action == $action) | .Test ) | unique) as $tests |
    $tests[] | . as $prefix |
    select(any($tests[]; . != $prefix and startswith($prefix+"/"))|not)
    ' \
    "$file"
}

printf "\n==================== TEST SUMMARY ===================="

printf "\nPASSED TESTS:"
process_tests "pass" "$REPORT_FILE"

printf "\nFAILED TESTS:"
process_tests "fail" "$REPORT_FILE"

printf "\n======================================================"

if [ "$REPORT_FILE" = "report.json" ] || [ "$REPORT_FILE" = "./report.json" ]; then
  rm -f "$REPORT_FILE"
fi
