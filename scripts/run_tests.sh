#!/bin/bash

get_git_root() {
  git rev-parse --show-toplevel
}

if get_git_root; then cd "$(get_git_root)" || exit; else echo "expects to be run from within the terraform provider git repo"; exit 1; fi
export REPO_ROOT=$PWD

IDENTIFIER=""

rerun_failed=false
specific_test=""
specific_package=""
skip_build=false
cleanup_id=""

while getopts ":r:t:p:s:c:" opt; do
  case $opt in
    r) rerun_failed=true ;;
    t) specific_test="$OPTARG" ;;
    p) specific_package="$OPTARG" ;;
    s) skip_build=true ;;
    c) cleanup_id="$OPTARG" ;;
    \?) cat <<EOT >&2 && exit 1 ;;
Invalid option -$OPTARG, valid options are
  -r to re-run failed tests, only used when running multiple tests
  -t to specify a specific test (eg. TestBase)
  -p to specify a specific test package (eg. base)
  -s to skip building the provider binary and test using a released version
  -c to run clean up only with the given id (eg. abc123)
EOT
  esac
done

if [ -n "$cleanup_id" ]; then
  export IDENTIFIER="$cleanup_id"
fi

run_tests() {
  local rerun=$1
  REPO_ROOT="$(git rev-parse --show-toplevel)"
  cd "$REPO_ROOT" || exit 1
  export REPO_ROOT=$PWD

  if [ "false" = "$SKIP_BUILD" ]; then
    install -d "$REPO_ROOT/.terraform" || true
    touch "$REPO_ROOT/.terraform/terraformrc" || true
    export TF_CLI_CONFIG_FILE="$REPO_ROOT/.terraform/terraformrc"
    cat <<EOF > "$REPO_ROOT/.terraform/terraformrc"
provider_installation {
  dev_overrides {
    "rancher/rancher2" = "$REPO_ROOT/bin"
  }
  direct {
    exclude = []
  }
}
EOF
  fi

  TEST_DIR="$REPO_ROOT/test"

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
  if [ -n "$specific_test" ]; then
    specific_test_flag="-run=$specific_test"
  fi

  local package_pattern=""
  if [ -n "$specific_package" ]; then
    package_pattern="$specific_package"
  else
    package_pattern="/..."
  fi
  # shellcheck disable=SC2086
  gotestsum \
    --format=standard-verbose \
    --jsonfile "/tmp/${IDENTIFIER}_test.log" \
    --post-run-command "sh /tmp/${IDENTIFIER}_test-processor" \
    --packages "$TEST_DIR$package_pattern" \
    -- \
    -parallel=2 \
    -count=1 \
    -failfast=1 \
    -timeout=300m \
    $rerun_flag \
    $specific_test_flag

  if [ "false" = "$SKIP_BUILD" ]; then
    rm -rf "$REPO_ROOT/.terraform"
  fi

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
echo 'if tmp directory is missing, try restarting dev environment'

if [ -z "$cleanup_id" ]; then
  if [ "false" = "$skip_build" ]; then
    echo 'building...'
    if ! "$REPO_ROOT/scripts/gobuild.sh"; then C=$?; echo "failed to compile provider, exit code $C"; exit $C; fi
    export SKIP_BUILD="false"
    echo "provider successfully compiles..."
  else
    echo "skipping build..."
    export SKIP_BUILD="true"
  fi

  # Test if tests can compile
  echo "checking tests for compile errors..."
  cd "$REPO_ROOT/test" || exit
  if ! go mod tidy; then C=$?; echo "failed to tidy, exit code $C"; exit $C; fi

  while read -r file; do
    echo "found $file";
    if ! go test -c "$file" -o "$file.test"; then C=$?; echo "failed to compile $file, exit code $C"; exit $C; fi
    rm -rf "$file.test"
  done <<<"$(find "$REPO_ROOT/test" -name '*.go')"

  echo "compile checks passed..."

  # Run tests initially
  run_tests false

  # Check if we need to rerun failed tests
  if [ "$rerun_failed" = true ] && [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
    echo "Rerunning failed tests..."
    run_tests true
  fi
fi

echo "Clearing leftovers with Id $IDENTIFIER in $AWS_REGION..."
sleep 60

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
fi

#   attempts=0
#   # shellcheck disable=SC2143
#   while [ -n "$(leftovers -d --iaas=aws --aws-region="$AWS_REGION" --type="elbv2-target-group" --filter="tf-" | grep -v 'AccessDenied')" ] && [ $attempts -lt 3 ]; do
#     leftovers --iaas=aws --aws-region="$AWS_REGION" --type="elbv2-target-group" --filter="tf-" --no-confirm | grep -v 'AccessDenied' || true
#     sleep 10
#     attempts=$((attempts + 1))
#   done

#   if [ $attempts -eq 3 ]; then
#     echo "Warning: Failed to clear all EC2 key pairs after 3 attempts."
#   fi
# fi

if [ -f "/tmp/${IDENTIFIER}_failed_tests.txt" ]; then
  echo "done, test failed"
  EXIT=1
else
  echo "done, test passed"
  EXIT=0
fi

exit $EXIT
