#!/usr/bin/env bash
set -euo pipefail

# This skill script extracts the "Purpose" and execution dates of all plans
# in the .agent/plans directory (excluding README.md) and outputs them as a Plan Log.

readonly PLANS_DIR=".agent/plans"

extract_date() {
    local file="$1"
    local date_val
    
    date_val=$(awk 'tolower($0) ~ /^\**executed date:\**/ { sub(/^\**[Ee]xecuted [Dd]ate:\**[ \t]*/, ""); print; exit }' "${file}")
    
    if [[ -z "${date_val}" ]]; then
        echo "Not specified"
    else
        echo "${date_val}"
    fi
}

extract_purpose() {
    local file="$1"
    local purpose_val
    
    purpose_val=$(awk 'tolower($0) ~ /^\**purpose:\**/ { sub(/^\**[Pp]urpose:\**[ \t]*/, ""); print; exit }' "${file}")
    
    if [[ -z "${purpose_val}" ]]; then
        echo "Not specified"
    else
        echo "${purpose_val}"
    fi
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
            
            date_val=$(extract_date "${file}")
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

main() {
    if [[ ! -d "${PLANS_DIR}" ]]; then
        echo "Error: Directory ${PLANS_DIR} does not exist." >&2
        exit 1
    fi
    
    generate_plan_log
}

main "$@"
