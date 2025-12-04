#!/bin/bash

rerun_failed=false
specific_test=""
specific_package=""
cleanup_id=""
slow_mode=false
dirty_mode=false

while getopts ":rsdt:p:c:" opt; do
  case $opt in
    r) rerun_failed=true ;;
    t) specific_test="$OPTARG" ;;
    p) specific_package="$OPTARG" ;;
    c) cleanup_id="$OPTARG" ;;
    d) dirty_mode=true ;;
    s) slow_mode=true ;;
    \?) cat <<EOT >&2 && exit 1 ;;
Invalid option -$OPTARG, valid options are
  -r to re-run failed tests
  -s to run tests in slow mode (one at a time to avoid AWS rate limiting)
  -c to run clean up only with the given id (eg. abc123)
  -d to skip cleanup (dirty mode)
  -t to specify a specific test (eg. TestBase)
  -p to specify a specific test package (eg. one)
Only one of -c, -t, or -p can be used at a time.
EOT
  esac
done

if [ $slow_mode == true ]; then
  echo "Running in slow mode: tests will be run one at a time to avoid AWS rate limiting."
elif [ $slow_mode == false ]; then
  echo "Running in normal mode: tests will be run in parallel."
fi
if [ $rerun_failed == true ]; then
  echo "Rerun failed tests is enabled."
elif [ $rerun_failed == false ]; then
  echo "Rerun failed tests is disabled."
fi
if [ -n "$specific_test" ]; then
  echo "Specific test to run: $specific_test"
else
  echo "No specific test to run."
fi
if [ -n "$specific_package" ]; then
  echo "Specific package to run: $specific_package"
else
  echo "No specific package to run."
fi
if [ -n "$cleanup_id" ]; then
  echo "Cleanup only mode enabled with id: $cleanup_id"
fi
if [ -n "$cleanup_id" ] && { [ -n "$specific_test" ] || [ -n "$specific_package" ]; }; then
  echo "Error: Only one of -c, -t, or -p can be used at a time." >&2
  exit 1
fi
if [ -n "$specific_test" ] && { [ -n "$specific_package" ] || [ -n "$cleanup_id" ]; }; then
  echo "Error: Only one of -c, -t, or -p can be used at a time." >&2
  exit 1
fi
if [ -n "$specific_package" ] && { [ -n "$specific_test" ] || [ -n "$cleanup_id" ]; }; then
  echo "Error: Only one of -c, -t, or -p can be used at a time." >&2
  exit 1
fi
if [ $dirty_mode == true ]; then
  echo "Running in dirty mode: skipping cleanup."
elif [ $dirty_mode == false ]; then
  echo "Running in normal mode: cleanup will try to remove all resources matching ID."
fi

# shellcheck disable=SC2143
if [ -n "$cleanup_id" ]; then
  export IDENTIFIER="$cleanup_id"
fi

REPO_ROOT="$(git rev-parse --show-toplevel)"

# Find the tests directory
TEST_DIR=""
if [ -d "tests" ]; then
  TEST_DIR="tests"
fi
if [ -d "test" ]; then
  TEST_DIR="test"
fi
if [ -d "test/tests" ]; then
  TEST_DIR="test/tests"
fi
if [ "$TEST_DIR" == "" ]; then
  echo "Error: Unable to find tests directory" >&2
  exit 1
fi

run_tests() {
  local rerun=$1
  local slow_mode=$2
  REPO_ROOT="$(git rev-parse --show-toplevel)"
  cd "$REPO_ROOT" || exit 1

  # Find the tests directory
  TEST_DIR=""
  if [ -d "tests" ]; then
    TEST_DIR="tests"
  fi
  if [ -d "test" ]; then
    TEST_DIR="test"
  fi
  if [ -d "test/tests" ]; then
    TEST_DIR="test/tests"
  fi
  if [ "$TEST_DIR" == "" ]; then
    echo "Error: Unable to find tests directory" >&2
    exit 1
  fi

  echo "" > "/tmp/${IDENTIFIER}_test.log"
  rm -f "/tmp/${IDENTIFIER}_failed_tests.txt"
  cat <<'EOF'> "/tmp/${IDENTIFIER}_test-processor"
echo "Passed: "
export PASS="$(jq -r '. | select(.Action == "pass") | select(.Test != null).Test' "/tmp/${IDENTIFIER}_test.log")"
echo $PASS | tr ' ' '\n'
echo " "
echo "Failed: "
export FAIL="$(jq -r '. | select(.Action == "fail") | select(.Test != null).Test' "/tmp/${IDENTIFIER}_test.log")"
echo $FAIL | tr ' ' '\n'
echo " "
if [ -n "$FAIL" ]; then
  echo $FAIL > "/tmp/${IDENTIFIER}_failed_tests.txt"
  exit 1
fi
exit 0
EOF
  chmod +x "/tmp/${IDENTIFIER}_test-processor"
  export NO_COLOR=1
  echo "starting tests..."
  cd "$TEST_DIR" || return 1;

  local rerun_flag=""
  if [ "$rerun" = true ] && [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
    # shellcheck disable=SC2002
    rerun_flag="-run=$(cat "/tmp/${IDENTIFIER}_failed_tests.txt" | tr '\n' '|')"
  fi

  local specific_test_flag=""
  # shellcheck disable=SC2143
  if [ -n "$specific_test" ]; then
    specific_test_flag="-run=$specific_test"
  fi

  local package_pattern=""
  # shellcheck disable=SC2143
  if [ -n "$specific_package" ]; then
    package_pattern="$specific_package"
  else
    package_pattern="..."
  fi

  # We need both -p and -parallel, as -p sets the number of packages to test in parallel,
  #  and -parallel sets the number of tests to run in parallel.
  # By setting both to 1, we ensure that tests are run sequentially, which can help avoid AWS rate limiting issues.
  # It does increase the runtime significantly though.
  local parallel_packages=""
  local parallel_tests=""
  if [ "$slow_mode" = true ]; then
    echo "Running in slow mode..."
    parallel_packages="-p=1"
    parallel_tests="-parallel=1"
  fi

  CMD=$(cat <<EOT
gotestsum \
  --format=standard-verbose \
  --jsonfile "/tmp/${IDENTIFIER}_test.log" \
  --post-run-command "sh /tmp/${IDENTIFIER}_test-processor" \
  --packages "$REPO_ROOT/$TEST_DIR/$package_pattern" \
  -- \
  -count=1 \
  -timeout=300m \
  -failfast \
  $parallel_packages \
  $parallel_tests \
  $rerun_flag \
  $specific_test_flag
EOT
)
  echo "Running command: $CMD"

  # shellcheck disable=SC2086
  gotestsum \
    --format=standard-verbose \
    --jsonfile "/tmp/${IDENTIFIER}_test.log" \
    --post-run-command "sh /tmp/${IDENTIFIER}_test-processor" \
    --packages "$REPO_ROOT/$TEST_DIR/$package_pattern" \
    -- \
    -count=1 \
    -timeout=300m \
    -failfast \
    $parallel_packages \
    $parallel_tests \
    $rerun_flag \
    $specific_test_flag

  return $?
}

if [ -z "$IDENTIFIER" ]; then
  IDENTIFIER="$(echo a-$RANDOM-d | base64 | tr -d '=')"
  export IDENTIFIER
fi
echo "id is: $IDENTIFIER..."
if [ -z "$GITHUB_TOKEN" ]; then echo "GITHUB_TOKEN isn't set"; else echo "GITHUB_TOKEN is set"; fi
if [ -z "$GITHUB_OWNER" ]; then echo "GITHUB_OWNER isn't set"; else echo "GITHUB_OWNER is set"; fi
if [ -z "$ZONE" ]; then echo "ZONE isn't set"; else echo "ZONE is set"; fi

if [ -z "$cleanup_id" ]; then

  D="$(pwd)"

  echo "tidying..."
  cd "$REPO_ROOT/$TEST_DIR" || exit
  if ! go mod tidy; then C=$?; echo "failed to tidy, exit code $C"; exit $C; fi

  echo "formatting tests..."
  gofmt -s -w -e .
  echo "done formatting"

  echo "checking tests for compile errors..."
  while IFS= read -r file; do
    echo "found $file";
    if ! go test -c "$file" -o "${file}.test"; then C=$?; echo "failed to compile $file, exit code $C"; exit $C; fi
    rm -rf "${file}.test"
  done <<< "$(find "$REPO_ROOT/$TEST_DIR" -not \( -path "$REPO_ROOT/$TEST_DIR/data" -prune \) -name '*.go')"
  echo "compile checks passed..."

  echo "checking tests for go lint errors..."
  if ! golangci-lint run; then echo "lint failed..."; exit 1; fi
  echo "lint errors complete"

  cd "$D" || exit

  echo "checking terraform configs for errors..."
  if ! tflint --recursive; then C=$?; echo "tflint failed, exit code $C"; exit $C; fi
  echo "terraform configs valid..."

  # Run tests initially
  run_tests false "$slow_mode"
  sleep 60

  # Check if we need to rerun failed tests
  if [ "$rerun_failed" = true ] && [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
    echo "Rerunning failed tests..."
    run_tests true "$slow_mode"
    sleep 60
  fi
fi

if [ $dirty_mode == true ]; then
  echo "Running in dirty mode, skipping cleanup..."
else
  echo "Starting cleanup..."
  sh "$REPO_ROOT/cleanup.sh" "$IDENTIFIER"
  C=$?
  if [ $C -ne 0 ]; then
    echo "Cleanup failed with exit code $C"
    exit $C
  fi
  echo "Cleanup completed successfully."
fi

if [ -n "$cleanup_id" ]; then
  # cleanup only mode
  exit 0
fi

if [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
  echo "done, test failed"
  exit 1
else
  echo "done, test passed"
  exit 0
fi
