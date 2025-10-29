#!/bin/sh
set -e

# summarize.sh - reads report.json and prints a summary

# Ensure jq is installed
if ! command -v jq > /dev/null; then
    echo "Error: 'jq' is not installed. Please install it to generate the summary."
    exit 1
fi

process_tests() {
  local action=$1
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
  jq --slurp -r --arg action "$action" \
    '
    (map(select(.Test and .Action == $action) | .Test ) | unique) as $tests |
    $tests[] | . as $prefix |
    select(any($tests[]; . != $prefix and startswith($prefix+"/"))|not)
    '\
    report.json
}

echo "\n==================== TEST SUMMARY ===================="

echo "\nPASSED TESTS:"
process_tests "pass"

echo "\nFAILED TESTS:"
process_tests "fail"

echo "\n======================================================"

rm -rf report.json
