#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-all}"

run_terraform() {
  echo "==> Linting Terraform files..."
  terraform fmt -check -recursive -diff || { echo "terraform format error, please run 'terraform fmt -recursive'"; exit 1; }
  tflint --recursive
}

run_actionlint() {
  echo "==> Linting GitHub workflows..."
  actionlint
}

run_shellcheck() {
  echo "==> Running shellcheck..."
  local files
  files=$(grep -Rl -e '^#!' . \
    | grep -v -E "^\./(\.git|\.terraform|\.agent|bin|node_modules)/" \
    | grep -v -E "\.(md|js)$" || true)

  if [[ -z "${files}" ]]; then
    echo "No shell scripts found to check."
    return 0
  fi

  while read -r file; do
    if [[ -f "${file}" ]]; then
      echo "Checking ${file}..."
      shellcheck -x "${file}"
    fi
  done <<< "${files}"
}

run_node_check() {
  echo "==> Running Node Syntax Check..."
  local files
  files=$(find . -type f \( -name "*.js" \) -not -path "*/node_modules/*")
  if [[ -z "${files}" ]]; then
    echo "No Node files found to check."
    return 0
  fi
  while read -r file; do
    if [[ -f "${file}" ]]; then
      echo "Checking ${file}..."
      node --check "${file}"
    fi
  done <<< "${files}"
}

run_node_test() {
  echo "==> Running Node Unit Tests..."
  local files=()
  while IFS= read -r -d '' file; do
    files+=("${file}")
  done < <(find . -type f \( -name "*.test.js" -o -name "*.test.mjs" \) -not -path "*/node_modules/*" -print0)

  if [[ ${#files[@]} -eq 0 ]]; then
    echo "No test files found to run."
    return 0
  fi
  node --test "${files[@]}"
}

run_eslint() {
  echo "==> Running ESLint..."
  npm ci
  eslint .
}

run_gitleaks() {
  echo "==> Scanning for secrets with gitleaks..."
  gitleaks detect --no-banner -v --no-git
  gitleaks detect --no-banner -v
}

case "${MODE}" in
  terraform)
    run_terraform
    ;;
  actionlint)
    run_actionlint
    ;;
  shellcheck)
    run_shellcheck
    ;;
  node-check)
    run_node_check
    ;;
  node-test)
    run_node_test
    ;;
  eslint)
    run_eslint
    ;;
  gitleaks)
    run_gitleaks
    ;;
  all)
    run_terraform
    run_actionlint
    run_shellcheck
    run_node_check
    run_node_test
    run_eslint
    run_gitleaks
    ;;
  *)
    echo "Error: Unknown lint mode: ${MODE}" >&2
    echo "Usage: $0 [terraform|actionlint|shellcheck|node-check|node-test|eslint|gitleaks|all]" >&2
    exit 1
    ;;
esac
