#!/bin/bash
# Configuration flags
rerun_failed=false
specific_test=""
specific_package=""
specific_fixture=""
fixture_group=""
cleanup_id=""
wait_time=""
slow_mode=false
dirty_mode=false

# Track whether cleanup has run
cleanup_has_run=false

# Cleanup function that will be called on exit
run_cleanup() {
  # Avoid running cleanup twice
  if [ "$cleanup_has_run" = true ]; then
    return 0
  fi
  cleanup_has_run=true

  # Skip if dirty mode or no identifier
  if [ "$dirty_mode" = true ] || [ -z "$IDENTIFIER" ]; then
    return 0
  fi

  echo ""
  echo "=== Cleanup ==="

  # Wait before cleanup if requested (for investigation)
  if [ -n "$WAIT" ] && [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
    echo "Tests failed. Waiting $WAIT seconds before cleanup for investigation..."
    sleep "$WAIT"
  fi

  # Locate repository root
  REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"

  # Check if cleanup script exists
  if [ -f "$REPO_ROOT/cleanup.sh" ]; then
    echo "Running cleanup script..."
    bash "$REPO_ROOT/cleanup.sh" "$IDENTIFIER"
    cleanup_exit=$?

    if [ $cleanup_exit -ne 0 ]; then
      echo "WARNING: Cleanup failed with exit code $cleanup_exit"
    else
      echo "✓ Cleanup completed successfully"
    fi
  else
    echo "WARNING: cleanup.sh not found, skipping automated cleanup"
    echo "You may need to manually clean up resources with ID: $IDENTIFIER"
  fi
}

# Set trap to run cleanup on exit, error, interrupt, or termination
trap run_cleanup EXIT INT TERM

# Parse command line options
while getopts ":rsdt:p:f:g:c:w:" opt; do
  case $opt in
    r) rerun_failed=true ;;
    t) specific_test="$OPTARG" ;;
    p) specific_package="$OPTARG" ;;
    f) specific_fixture="$OPTARG" ;;
    g) fixture_group="$OPTARG" ;;
    c) cleanup_id="$OPTARG" ;;
    w) wait_time="$OPTARG" ;;
    d) dirty_mode=true ;;
    s) slow_mode=true ;;
    \?) cat <<EOT >&2 && exit 1
Invalid option: -$OPTARG

Usage: $0 [OPTIONS]

Options:
  -r              Re-run failed tests, requires dirty mode.
  -s              Run tests in slow mode (sequential, one at a time)
  -d              Skip cleanup (dirty mode)
  -t TEST         Run specific test (eg. TestMatrix)
  -p PACKAGE      Run specific test package (eg. one)
  -f FIXTURE      Run specific fixture combination (eg. "sle-micro-61-canal-stable-one-rpm-ipv4")
  -g GROUP        Run specific fixture group (eg. "necessary" or "extended")
  -c ID           Cleanup-only mode with the given identifier
  -w SECONDS      Wait time in seconds before cleanup on test failure (for investigation)

Notes:
  - Only one of -c, -t, -p, -f, or -g can be used at a time
  - The -f option sets the COMBO environment variable for fixture selection
  - The -g option sets the GROUP environment variable for fixture group selection
  - The -w option sets the WAIT environment variable for error investigation
EOT
  esac
done

# Validate mutually exclusive options
exclusive_count=0
[ -n "$cleanup_id" ] && ((exclusive_count++))
[ -n "$specific_test" ] && ((exclusive_count++))
[ -n "$specific_package" ] && ((exclusive_count++))
[ -n "$specific_fixture" ] && ((exclusive_count++))
[ -n "$fixture_group" ] && ((exclusive_count++))

if [ $exclusive_count -gt 1 ]; then
  echo "Error: Only one of -c, -t, -p, -f, or -g can be used at a time." >&2
  exit 1
fi

if [ "$rerun_failed" = true ] && [ "$dirty_mode" = false ]; then
  echo "Error: Rerun failed tests requires dirty mode." >&2
  exit 1
fi

# Display configuration
echo "=== Test Configuration ==="
if [ "$slow_mode" = true ]; then
  echo "Mode: Slow (sequential execution to avoid AWS rate limiting)"
else
  echo "Mode: Normal (parallel execution)"
fi

if [ "$rerun_failed" = true ]; then
  echo "Rerun failed tests: Enabled"
fi

if [ "$dirty_mode" = true ]; then
  echo "Cleanup: Disabled (dirty mode)"
else
  echo "Cleanup: Enabled"
fi

if [ -n "$specific_test" ]; then
  echo "Specific test: $specific_test"
fi

if [ -n "$specific_package" ]; then
  echo "Specific package: $specific_package"
fi

if [ -n "$specific_fixture" ]; then
  echo "Specific fixture: $specific_fixture"
fi

if [ -n "$fixture_group" ]; then
  echo "Fixture group: $fixture_group"
fi

if [ -n "$cleanup_id" ]; then
  echo "Cleanup-only mode: $cleanup_id"
fi

if [ -n "$wait_time" ]; then
  if ! [[ "$wait_time" =~ ^[0-9]+$ ]]; then
    echo "Error: -w must be an integer number of seconds, got: $wait_time" >&2
    exit 1
  fi
  echo "Wait time on failure: $wait_time seconds"
fi

echo "=========================="
echo ""

# Set cleanup ID if provided
if [ -n "$cleanup_id" ]; then
  export IDENTIFIER="$cleanup_id"
fi

# Set COMBO environment variable for fixture selection
export COMBO="$specific_fixture"
if [ -n "$COMBO" ]; then
  echo "COMBO environment variable set to: $COMBO"
fi

# Set GROUP environment variable for fixture group selection
export GROUP="$fixture_group"
if [ -n "$GROUP" ]; then
  echo "GROUP environment variable set to: $GROUP"
fi

# Set WAIT environment variable for error investigation
export WAIT="$wait_time"
if [ -n "$WAIT" ]; then
  echo "WAIT environment variable set to: $WAIT seconds"
fi

# Locate repository root
REPO_ROOT="$(git rev-parse --show-toplevel)"

# Find the tests directory
find_test_dir() {
  local test_dir=""
  if [ -d "$REPO_ROOT/test/tests" ]; then
    test_dir="test/tests"
  elif [ -d "$REPO_ROOT/tests" ]; then
    test_dir="tests"
  elif [ -d "$REPO_ROOT/test" ]; then
    test_dir="test"
  else
    echo "Error: Unable to find tests directory" >&2
    exit 1
  fi
  echo "$test_dir"
}

TEST_DIR="$(find_test_dir)"
echo "Using test directory: $TEST_DIR"
echo ""

# Run tests function
run_tests() {
  local rerun=$1
  local slow_mode=$2

  cd "$REPO_ROOT" || exit 1

  # Setup test log and processor
  echo "" > "/tmp/${IDENTIFIER}_test.log"

  cat <<'EOF' > "/tmp/${IDENTIFIER}_test-processor"
echo "Passed: "
export PASS="$(jq -r '. | select(.Action == "pass") | select(.Test != null).Test' "/tmp/${IDENTIFIER}_test.log")"
echo "$PASS" | tr ' ' '\n'
echo " "
echo "Failed: "
export FAIL="$(jq -r '. | select(.Action == "fail") | select(.Test != null).Test' "/tmp/${IDENTIFIER}_test.log")"
echo "$FAIL" | tr ' ' '\n'
echo " "
if [ -n "$FAIL" ]; then
  echo "$FAIL" > "/tmp/${IDENTIFIER}_failed_tests.txt"
  exit 1
fi
exit 0
EOF
  chmod +x "/tmp/${IDENTIFIER}_test-processor"

  export NO_COLOR=1
  echo "Starting tests..."
  cd "$TEST_DIR" || exit 1

  # Build rerun flag
  local rerun_flag=""
  if [ "$rerun" = true ] && [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
    # shellcheck disable=SC2002
    rerun_flag="-run=$(cat "/tmp/${IDENTIFIER}_failed_tests.txt" | tr '\n' '|' | sed 's/|$//')"
    echo "Rerunning failed tests: $rerun_flag"
    rm -f "/tmp/${IDENTIFIER}_failed_tests.txt"
  fi

  # Build specific test flag
  local specific_test_flag=""
  if [ -n "$specific_test" ] && [ "$rerun" != true ]; then
    specific_test_flag="-run=$specific_test"
    echo "Running specific test: $specific_test"
  fi

  # Build package pattern
  local package_pattern=""
  if [ -n "$specific_package" ]; then
    package_pattern="$specific_package"
    echo "Running specific package: $specific_package"
  else
    package_pattern="..."
  fi

  # Build parallel flags for slow mode
  local parallel_packages=""
  local parallel_tests=""
  if [ "$slow_mode" = true ]; then
    echo "Slow mode: Running tests sequentially"
    parallel_packages="-p=1"
    parallel_tests="-parallel=1"
  else
    parallel_tests="-parallel=10"
  fi

  # Display the command that will be run
  echo ""
  echo "Test command:"
  echo "  gotestsum --format=standard-verbose \\"
  echo "    --jsonfile /tmp/${IDENTIFIER}_test.log \\"
  echo "    --post-run-command 'sh /tmp/${IDENTIFIER}_test-processor' \\"
  echo "    --packages $REPO_ROOT/$TEST_DIR/$package_pattern \\"
  echo "    -- -count=1 -timeout=300m -failfast \\"
  echo "      $parallel_packages $parallel_tests \\"
  echo "      $rerun_flag $specific_test_flag"
  echo ""

  # Run tests
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

# Generate and export identifier
if [ -z "$IDENTIFIER" ]; then
  IDENTIFIER="tfci-$(printf "%04x%04x" $RANDOM $RANDOM)"
  export IDENTIFIER
fi

echo "Test identifier: $IDENTIFIER"
echo ""

# Check required environment variables
echo "=== Environment Check ==="
if [ -z "$GITHUB_TOKEN" ]; then
  echo "WARNING: GITHUB_TOKEN is not set"
else
  echo "GITHUB_TOKEN: Set"
fi

if [ -z "$GITHUB_OWNER" ]; then
  echo "WARNING: GITHUB_OWNER is not set"
else
  echo "GITHUB_OWNER: Set ($GITHUB_OWNER)"
fi

if [ -z "$ZONE" ]; then
  echo "WARNING: ZONE is not set"
else
  echo "ZONE: Set"
fi
echo "========================="
echo ""

# If cleanup-only mode, skip tests and run cleanup directly
if [ -n "$cleanup_id" ]; then
  echo "Cleanup-only mode enabled, skipping tests..."
  # In cleanup-only mode, we want to run cleanup immediately
  run_cleanup
  echo "Cleanup-only mode completed"
  exit 0
else
  # Pre-test validation
  current_dir="$(pwd)"

  echo "=== Pre-Test Validation ==="

  echo "Running go mod tidy..."
  cd "$REPO_ROOT/$TEST_DIR" || exit 1
  if ! go mod tidy; then
    echo "ERROR: go mod tidy failed"
    exit 1
  fi
  echo "go mod tidy passed"

  echo "Formatting tests..."
  gofmt -s -w -e .
  echo "Formatting complete"

  echo "Checking for compile errors..."
  if ! go list ./... | grep -v '/data' | xargs -r go test -run='^$'; then
    echo "ERROR: Compile checks failed"
    exit 1
  fi
  echo "Compile checks passed"

  echo "Running go lint..."
  if ! golangci-lint run; then
    echo "ERROR: Linting failed"
    exit 1
  fi
  echo "✓ Lint passed"

  cd "$current_dir" || exit 1

  echo "Checking terraform configs..."
  if ! tflint --recursive; then
    echo "ERROR: tflint failed"
    exit 1
  fi
  echo "Terraform configs valid"
  make build
  echo "============================"
  echo ""

  # Clear failed tests before initial run
  rm -f "/tmp/${IDENTIFIER}_failed_tests.txt"

  # Run tests initially
  echo "=== Running Tests ==="
  run_tests false "$slow_mode"
  test_exit_code=$?

  if [ $test_exit_code -ne 0 ]; then
    echo "Tests failed with exit code: $test_exit_code"
  else
    echo "Tests passed"
  fi

  # Brief pause between test runs
  sleep 5

  # Check if we need to rerun failed tests
  if [ "$rerun_failed" = true ] && [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
    echo ""
    echo "=== Rerunning Failed Tests ==="
    run_tests true "$slow_mode"
    test_exit_code=$?

    if [ $test_exit_code -ne 0 ]; then
      echo "Rerun failed with exit code: $test_exit_code"
    else
      echo "All tests passed on rerun"
    fi

    sleep 5
  fi
fi

echo ""
echo "=== Test Summary ==="

# Exit with appropriate code based on test results
if [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
  echo "Tests FAILED"
  echo "Failed tests logged to: /tmp/${IDENTIFIER}_failed_tests.txt"
  exit 1
else
  echo "All tests PASSED"
  exit 0
fi
