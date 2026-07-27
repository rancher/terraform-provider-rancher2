#!/usr/bin/env bash
set -euo pipefail

workflows=$(find .github/workflows -type f \( -name '*.yml' -o -name '*.yaml' \) 2>/dev/null || echo "")

if [[ -z "${workflows}" ]]; then
    echo "No workflow files found."
    exit 0
fi

for workflow in ${workflows}; do
    echo "Processing ${workflow}..."
    
    # Extract unique repos from the workflow file
    # Format expected: # https://github.com/owner/repo/releases
    repos=$(awk '/^[[:space:]]*#[[:space:]]*https:\/\/github\.com\/[^\/]+\/[^\/]+\/releases/ {
        str = $0
        sub(/^[[:space:]]*#[[:space:]]*https:\/\/github\.com\//, "", str)
        sub(/\/releases.*/, "", str)
        print str
    }' "${workflow}" | sort -u || true)

    for repo in ${repos}; do
        echo "Found release link for ${repo}"
        
        # Get latest release tag
        tag=$(gh api "/repos/${repo}/releases/latest" --jq '.tag_name' 2>/dev/null || true)
        
        if [[ -z "${tag}" ]]; then
            echo "Failed to get latest tag for ${repo}"
            continue
        fi

        # Get commit sha for tag
        sha=$(gh api "/repos/${repo}/commits/${tag}" --jq '.sha' 2>/dev/null || true)
        
        if [[ -z "${sha}" ]]; then
            echo "Failed to get commit SHA for ${repo} tag ${tag}"
            continue
        fi

        echo "Latest version for ${repo} is ${tag} (${sha})"
        
        # Use a temporary file for replacement
        tmp_file=$(mktemp)
        
        awk -v repo="${repo}" -v new_sha="${sha}" -v new_tag="${tag}" '
        {
            if ($0 ~ "^[[:space:]]*#[[:space:]]*https://github.com/" repo "/releases") {
                print $0
                getline
                if ($0 ~ "^[[:space:]]*- uses: " repo "@") {
                    idx = index($0, "- uses: ")
                    indent = substr($0, 1, idx - 1)
                    print indent "- uses: " repo "@" new_sha " # " new_tag
                } else {
                    print $0
                }
                next
            }
            print $0
        }
        ' "${workflow}" > "${tmp_file}"
        
        mv "${tmp_file}" "${workflow}"
    done
done

echo "Done updating action versions."
