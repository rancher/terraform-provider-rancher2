#!/usr/bin/env bash
#
# Skill: generate-plan-log.sh
# Description: Extracts the Purpose and executed dates of plans in .agent/plans and compiles a Plan Log.
# Usage: .agent/skills/generate-plan-log.sh

set -euo pipefail

# This skill script extracts the "Purpose" and execution dates of all plans
# in the .agent/plans directory (excluding README.md) and outputs them as a Plan Log.

readonly PLANS_DIR=".agent/plans"

extract_date() {
    local file="$1"
    local date_val
    
    date_val=$(awk 'tolower($0) ~ /executed date/ { sub(/^[^:]*:[ \t]*/, ""); gsub(/^[ \t*`_]+|[ \t*`_]+$/, ""); print; exit }' "${file}")
    
    if [[ -z "${date_val}" ]]; then
        echo "Not specified"
    else
        echo "${date_val}"
    fi
}

extract_purpose() {
    local file="$1"
    local purpose_val
    
    purpose_val=$(awk 'tolower($0) ~ /purpose/ { sub(/^[^:]*:[ \t]*/, ""); gsub(/^[ \t*`_]+|[ \t*`_]+$/, ""); print; exit }' "${file}")
    
    if [[ -z "${purpose_val}" ]]; then
        echo "Not specified"
    else
        echo "${purpose_val}"
    fi
}

resolve_executed_date() {
    local file="$1"
    local date_val
    date_val=$(extract_date "${file}")
    
    local lower_date
    lower_date=$(echo "${date_val}" | tr '[:upper:]' '[:lower:]')

    if [[ -z "${date_val}" || "${date_val}" == "Not specified" || "${lower_date}" == *"pending"* ]]; then
        # Try to extract the creation date from git history (oldest commit first)
        local git_date
        git_date=$(git log --follow --format="%ad" --date=short -- "${file}" | tail -n 1 || true)
        
        # Fallback to latest commit if tail is empty
        if [[ -z "${git_date}" ]]; then
            git_date=$(git log -1 --format="%ad" --date=short -- "${file}" || true)
        fi
        
        if [[ -n "${git_date}" ]]; then
            # Update the plan file in-place using a single sed execution with multiple expressions
            # This prevents overwriting the original .bak file and breaking the subsequent cmp check.
            sed -i.bak -E \
                -e "s/(([Ee]xecuted [Dd]ate[ \t]*\**:[ \t]*\**)[ \t]*)[Pp]ending/\1${git_date}/g" \
                -e "s/(([Ee]xecuted [Dd]ate[ \t]*\**:[ \t]*\**)[ \t]*)[Nn]ot [Ss]pecified/\1${git_date}/g" \
                "${file}"
            
            # Check if file was actually modified
            if ! cmp -s "${file}" "${file}.bak"; then
                echo "Self-Healed: Updated '${file}' execution date to '${git_date}' using Git history." >&2
                # Re-extract the date now that it's updated
                date_val=$(extract_date "${file}")
            else
                # Fallback to git_date for the log even if replacement didn't modify file (e.g. read-only fallback)
                date_val="${git_date}"
            fi
            rm -f "${file}.bak"
        fi
    fi
    echo "${date_val}"
}

get_sort_key() {
    local date_val="$1"
    if [[ -z "${date_val}" || "${date_val}" == "Not specified" ]]; then
        echo "0000-00"
        return
    fi
    
    local lower_date
    lower_date=$(echo "${date_val}" | tr '[:upper:]' '[:lower:]')
    
    if [[ "${lower_date}" == *"pending"* ]]; then
        echo "9999-99"
        return
    fi
    
    # If the date is already in YYYY-MM-DD format (or starts with YYYY-MM)
    if [[ "${date_val}" =~ ^([0-9]{4})-([0-9]{2})-([0-9]{2})$ ]]; then
        echo "${BASH_REMATCH[1]}-${BASH_REMATCH[2]}"
        return
    fi

    # Otherwise, parse format like "Month Day, Year" or "Day Month Year"
    # We clean up any commas, brackets, etc.
    local clean_date
    clean_date=$(echo "${date_val}" | tr -d ',')
    
    # Try to find a 4-digit year in the cleaned string
    local year="0000"
    if [[ "${clean_date}" =~ ([0-9]{4}) ]]; then
        year="${BASH_REMATCH[1]}"
    fi

    # Try to find month
    local month_num="00"
    case "${lower_date}" in
        *jan*) month_num="01" ;;
        *feb*) month_num="02" ;;
        *mar*) month_num="03" ;;
        *apr*) month_num="04" ;;
        *may*) month_num="05" ;;
        *jun*) month_num="06" ;;
        *jul*) month_num="07" ;;
        *aug*) month_num="08" ;;
        *sep*) month_num="09" ;;
        *oct*) month_num="10" ;;
        *nov*) month_num="11" ;;
        *dec*) month_num="12" ;;
    esac
    
    echo "${year}-${month_num}"
}

generate_plan_log() {
    echo "# Plan Log"
    echo ""
    
    local has_plans=false
    local tmp_dir
    tmp_dir=$(mktemp -d)
    
    for file in "${PLANS_DIR}"/*.md; do
        # Check if the file exists (in case glob doesn't match anything)
        if [[ -f "${file}" ]]; then
            local filename
            filename=$(basename "${file}")
            
            # Skip README.md
            if [[ "${filename}" == "README.md" ]]; then
                continue
            fi
            
            has_plans=true
            local plan_name="${filename%.md}"
            local date_val
            local purpose_val
            local sort_key
            
            date_val=$(resolve_executed_date "${file}")
            purpose_val=$(extract_purpose "${file}")
            sort_key=$(get_sort_key "${date_val}")
            
            local out_file="${tmp_dir}/${sort_key}_${filename}.txt"
            {
                echo "## ${plan_name}"
                echo "- **Date:** ${date_val}"
                echo "- **Purpose:** ${purpose_val}"
                echo ""
            } > "${out_file}"
        fi
    done
    
    if [[ "${has_plans}" == false ]]; then
        echo "No plans found in ${PLANS_DIR}."
    else
        find "${tmp_dir}" -name "*.txt" | sort -r | while read -r f; do
            cat "$f"
        done
    fi
    
    rm -rf "${tmp_dir}"
}

show_help() {
    cat <<EOF
Usage: generate-plan-log.sh [options]

Extracts the "Purpose" and execution dates of all plans in the .agent/plans directory
(excluding README.md) and outputs them as a sorted, combined Plan Log.

Options:
  -h, --help           Show this help message and exit.

Examples:
  .agent/skills/generate-plan-log.sh
EOF
}

main() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                show_help
                exit 0
                ;;
            -*)
                echo "Error: Unknown option: $1" >&2
                show_help
                exit 1
                ;;
            *)
                echo "Error: Unexpected argument: $1" >&2
                show_help
                exit 1
                ;;
        esac
    done

    if [[ ! -d "${PLANS_DIR}" ]]; then
        echo "Error: Directory ${PLANS_DIR} does not exist." >&2
        exit 1
    fi
    
    generate_plan_log
}

main "$@"
