#!/bin/bash

specific_test=""
specific_package=""
cleanup_id=""
slow_mode=false

while getopts ":st:p:c:" opt; do
  case $opt in
    t) specific_test="$OPTARG" ;;
    p) specific_package="$OPTARG" ;;
    c) cleanup_id="$OPTARG" ;;
    s) slow_mode=true ;;
    \?) cat <<EOT >&2 && exit 1 ;;
Invalid option -$OPTARG, valid options are
  -s to run tests in slow mode (one at a time to avoid AWS rate limiting)
  -c to run clean up only with the given id (eg. abc123)
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

# shellcheck disable=SC2143
if [ -n "$cleanup_id" ]; then
  export IDENTIFIER="$cleanup_id"
fi

REPO_ROOT="$(git rev-parse --show-toplevel)"

# Find the tests directory
TEST_DIR=""
if [ -d "tests" ]; then
  TEST_DIR="tests"
elif [ -d "test/tests" ]; then
  TEST_DIR="test/tests"
elif [ -d "test" ]; then
  TEST_DIR="test"
else
  echo "Error: Unable to find tests directory" >&2
  exit 1
fi

run_tests() {
  local slow_mode=$1
  local TEST_DIR="$2"

  REPO_ROOT="$(git rev-parse --show-toplevel)"
  cd "$REPO_ROOT" || exit 1

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
  --format "standard-verbose" \
  --jsonfile "/tmp/${IDENTIFIER}_test.log" \
  --post-run-command "sh /tmp/${IDENTIFIER}_test-processor" \
  --packages "$REPO_ROOT/$TEST_DIR/$package_pattern" \
  -- \
  -count=1 \
  -timeout=300m \
  $parallel_packages \
  $parallel_tests \
  $specific_test_flag \
  -failfast
EOT
)
  echo "Running command: $CMD"

  # shellcheck disable=SC2086
  gotestsum \
    --format="standard-verbose" \
    --jsonfile="/tmp/${IDENTIFIER}_test.log" \
    --post-run-command="sh /tmp/${IDENTIFIER}_test-processor" \
    --packages="$REPO_ROOT/$TEST_DIR/$package_pattern" \
    -- \
    -count=1 \
    -timeout=300m \
    $parallel_packages \
    $parallel_tests \
    $specific_test_flag \
    -failfast

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
  echo "checking tests for compile errors..."
  D="$(pwd)"

  cd "$REPO_ROOT/$TEST_DIR" || exit
  if ! go mod tidy; then C=$?; echo "failed to tidy, exit code $C"; exit $C; fi
  echo "completed tidy..."

  while IFS= read -r file; do
    echo "found $file";
    if ! go test -c "$file" -o "${file}.test"; then C=$?; echo "failed to compile $file, exit code $C"; exit $C; fi
    rm -rf "${file}.test"
  done <<< "$(find "$REPO_ROOT/$TEST_DIR" -not \( -path "$REPO_ROOT/$TEST_DIR/data" -prune \) -name '*_test.go')"
  echo "compile checks passed..."

  cd "$D" || exit

  echo "checking terraform configs for errors..."
  if ! tflint --recursive; then C=$?; echo "tflint failed, exit code $C"; exit $C; fi
  echo "terraform configs valid..."

  make build

  run_tests $slow_mode "$TEST_DIR"
  sleep 60
fi

echo "Clearing leftovers with Id $IDENTIFIER in $AWS_REGION..."

# shellcheck disable=SC2143
if [ -n "$IDENTIFIER" ]; then
  attempts=0
  # shellcheck disable=SC2143
  while [ -n "$(leftovers -d --iaas=aws --aws-region="$AWS_REGION" --filter="Id:$IDENTIFIER" | grep -v 'AccessDenied')" ] && [ $attempts -lt 3 ]; do
    leftovers --iaas=aws --aws-region="$AWS_REGION" --filter="Id:$IDENTIFIER" --no-confirm | grep -v 'AccessDenied' || true
    sleep 10
    attempts=$((attempts + 1))
  done

  if [ $attempts -eq 3 ]; then
    echo "Warning: Failed to clear all resources after 3 attempts."
  fi

  # remove key pairs
  attempts=0
  # shellcheck disable=SC2143
  while [ -n "$(leftovers -d --iaas=aws --aws-region="$AWS_REGION" --type="ec2-key-pair" --filter="terraform-ci-$IDENTIFIER" | grep -v 'AccessDenied')" ] && [ $attempts -lt 3 ]; do
    leftovers --iaas=aws --aws-region="$AWS_REGION" --type="ec2-key-pair" --filter="terraform-ci-$IDENTIFIER" --no-confirm | grep -v 'AccessDenied' || true
    sleep 10
    attempts=$((attempts + 1))
  done

  if [ $attempts -eq 3 ]; then
    echo "Warning: Failed to clear all EC2 key pairs after 3 attempts."
  fi

  # remove s3 storage
  attempts=0
  ID="$(aws s3 ls | grep -i "$IDENTIFIER" | awk '{print $3}')"
  # shellcheck disable=SC2143
  while [ -n "$(aws s3 ls | grep -i "$IDENTIFIER")" ] && [ $attempts -lt 3 ]; do
    echo "found s3 bucket $ID, removing..."
    while read -r v; do
      if [ -z "$v" ]; then continue; fi;
      aws s3api delete-object --bucket "$(echo "$ID" | tr '[:upper:]' '[:lower:]')" --key "tfstate" --version-id="$v" | /dev/null;
    done <<<"$(aws s3api list-object-versions --bucket "$(echo "$ID" | tr '[:upper:]' '[:lower:]')" | jq -r '.Versions[]?.VersionId')"

    while read -r v; do
      if [ -z "$v" ]; then continue; fi;
      aws s3api delete-object --bucket "$(echo "$ID" | tr '[:upper:]' '[:lower:]')" --key "tfstate" --version-id="$v";
    done <<<"$(aws s3api list-object-versions --bucket "$(echo "$ID" | tr '[:upper:]' '[:lower:]')" | jq -r '.DeleteMarkers[]?.VersionId')"

    aws s3api delete-bucket --bucket "$(echo "$ID" | tr '[:upper:]' '[:lower:]')";

    sleep 10
    attempts=$((attempts + 1))
  done

  # remove load balancer target groups
  attempts=0
  # shellcheck disable=SC2143
  while [ $attempts -lt 3 ]; do
    while read -r line; do
      if [ -z "$line" ]; then continue; fi
      echo "removing load balancer target group, $line..."
      aws elbv2 delete-target-group --target-group-arn "$line";
    done <<<"$(
      while read -r line; do
        if [ -z "$line" ]; then continue; fi
        aws elbv2 describe-tags --resource-arns "$line" | jq -r --arg id "$IDENTIFIER" '.TagDescriptions[] | select(any(.Tags[]; .Key == "Id" and .Value == $id)) | .ResourceArn // ""';
      done <<<"$(aws elbv2 describe-target-groups | jq -r '.TargetGroups[]?.TargetGroupArn')"
    )"
    sleep 10
    attempts=$((attempts + 1))
  done
fi

if [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
  echo "done, test failed"
  exit 1
else
  echo "done, test passed"
  exit 0
fi
