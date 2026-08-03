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
    | grep -v -E "\.md$" || true)

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
    run_eslint
    run_gitleaks
    ;;
  *)
    echo "Error: Unknown lint mode: ${MODE}" >&2
    echo "Usage: $0 [terraform|actionlint|shellcheck|node-check|eslint|gitleaks|all]" >&2
    exit 1
    ;;
esac
