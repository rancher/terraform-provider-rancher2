#!/bin/bash
set -e
set -o pipefail
# Check commit messages
# This steps enforces https://www.conventionalcommits.org/en/v1.0.0/
# This format enables automatic generation of changelogs and versioning
git version
pwd
git status

if [ -z "${PR_NUMBER:-}" ]; then
  echo "Error: PR_NUMBER environment variable is not set. Please provide it (e.g., PR_NUMBER=123 ./validate-commit-message.sh)"
  exit 1
fi

filter() {
  COMMIT="$1"
  output="$(echo "$COMMIT" | grep -v -e '^fix: ' -e '^feature: ' -e '^feat: ' -e '^refactor!: ' -e '^feature!: ' -e '^feat!: ' -e '^chore(main): ' -e '^Merge branch ' || true)"
  echo "$output"
}
prefix_check() {
  message="$1"
  if [ "" != "$(filter "$message")" ]; then
    cat <<EOF
...Commit message does not start with the required prefix.
Please use one of the following prefixes: "fix:", "feature:", "feat:", "refactor!:", "feature!:", or "feat!:".
This enables release-please to automatically format release notes based on the commit message.
$message
EOF
    exit 1
  else
    echo "...Commit message starts with the required prefix."
  fi
}
empty_check() {
  message="$1"
  if [ "" == "$message" ]; then
    echo "...Empty commit message."
    exit 1
  else
    echo "...Commit message isnt empty."
  fi
}
length_check() {
  message="$1"
  length="$(wc -m <<<"$message")"
  if [ "$length" -gt 100 ]; then
    echo "...Commit message subject line should be less than 100 characters, found $length."
    exit 1
  else
    echo "...Commit message subject line is less than 100 characters."
  fi
}
spell_check() {
  message="$1"
  WORDS="$(aspell list --dont-validate-words <<<"$message")"
  if [ "" != "$WORDS" ]; then
    echo "...Commit message contains spelling errors on: ^$WORDS\$"
    echo "...Also try updating the PR title."
    echo "...If this is a mistake, add your word to the aspell_custom.txt file, it is case insensitive."
    exit 1
  else
    echo "...Commit message doesnt contain spelling errors."
  fi
}

# Fetch the commit messages

COMMIT_MESSAGES="$(gh pr view "$PR_NUMBER" --json commits | jq -r '.commits[].messageHeadline')"
echo "Commit messages found: "
echo "$COMMIT_MESSAGES"

while IFS= read -r message; do
  echo "checking message ^$message\$"
  empty_check "$message"
  prefix_check "$message"
  length_check "$message"
  spell_check "$message"
  echo "message ^$message\$ passed all checks"
done <<<"$COMMIT_MESSAGES"
