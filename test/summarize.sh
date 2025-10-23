#!/bin/sh


# summarize.sh - reads report.json and prints a summary

# Ensure jq is installed
if ! command -v jq > /dev/null; then
    echo "Error: 'jq' is not installed. Please install it to generate the summary."
    exit 1
fi

process_tests() {
  local action=$1
  # slurp is important here, it reads the objects into an array for further processing
  jq --slurp -r --arg action "$action" \
    '
    # select all objects with a .Test listed that matches $acton and store it as $tests
    (map(select(.Test and .Action == $action) | .Test ) | unique) as $tests |
    # iterate through tests and save the test name as $prefix
    $tests[] | . as $prefix |
    # filter out any test names that dont exactly match the current test name, but do have the test name in them followed by a slash
    select(any($tests[]; . != $prefix and startswith($prefix+"/"))|not)
    # this leaves only test names that dont have duplicate prefixes, ie. the actual tests and not the parent blocks
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
