#!/usr/bin/env bash
#
# Skill: update-modules.sh
# Description: Dynamically detects and updates all Terraform Registry module references in all .tf files to their latest registry versions.
# Usage: .gemini/skills/update-modules.sh

set -euo pipefail

show_help() {
  cat <<EOF
Usage: update-modules.sh [options]

Scans all Terraform (.tf) files, dynamically extracts all unique public registry module references
(format: namespace/name/provider), queries the Terraform Registry API with automatic retry/backoff,
and updates their version parameters to the latest available releases.

Options:
  --list-modules       List all modules declared in Terraform files and their current versions.
  -h, --help           Show this help message and exit.

Examples:
  .gemini/skills/update-modules.sh --list-modules
  .gemini/skills/update-modules.sh
EOF
}

# Safely executes a command with retry and exponential backoff
# Usage: run_with_retry cmd args...
run_with_retry() {
  local max_attempts=5
  local base_delay=2
  local attempt=1
  local exit_code=0

  while true; do
    if "$@"; then
      return 0
    else
      exit_code=$?
    fi

    if [[ ${attempt} -ge ${max_attempts} ]]; then
      echo "Error: Command '$*' failed after ${max_attempts} attempts." >&2
      return ${exit_code}
    fi

    local delay
    delay=$(( base_delay * (2 ** (attempt - 1)) ))
    echo "Warning: Command failed (exit code ${exit_code}). Retrying in ${delay} seconds (attempt ${attempt}/${max_attempts})..." >&2
    sleep "${delay}"
    attempt=$((attempt + 1))
  done
}

get_latest_version() {
  local module_name="$1"
  local latest_version
  
  latest_version=$(run_with_retry curl -s "https://registry.terraform.io/v1/modules/${module_name}" 2>/dev/null | jq -r '.version // empty')
  
  if [[ -z "${latest_version}" || "${latest_version}" == "null" ]]; then
    return 1
  fi
  
  echo "${latest_version}"
}

list_current_modules() {
  local search_dir="."
  echo "Scanning Terraform files to list module instances and versions..." >&2

  # We use find to collect all .tf files, and pass them to our portable awk parser
  local modules_data
  # shellcheck disable=SC2016
  modules_data=$(find "${search_dir}" -type d \( -name ".git" -o -name ".terraform" -o -name "tf_plugin_cache" \) -prune -o -type f -name "*.tf" -print0 2>/dev/null \
    | xargs -0 awk '
      # Keep track of the most recently declared module block name
      $1 == "module" {
        split($2, parts, "\"")
        mod_name = parts[2]
      }
      
      # Match only lines declaring a public registry source path (namespace/name/provider)
      $0 ~ /[ \t]*source[ \t]*=[ \t]*"[A-Za-z0-9_-]+\/[A-Za-z0-9_-]+\/[A-Za-z0-9_-]+"/ {
        split($0, parts, "\"")
        source_val = parts[2]
        
        # In Terraform conventional formats, the version line immediately follows the source line
        getline next_line
        
        # Check if this next line indeed defines the module version attribute
        if (next_line ~ /[ \t]*version[ \t]*=[ \t]*/) {
          split(next_line, ver_parts, "\"")
          ver_val = ver_parts[2]
          
          # Clean up path prefix for file name
          filename_clean = FILENAME
          sub(/^\.\//, "", filename_clean)
          
          # Output tab-separated records: label \t source \t version \t filename
          print mod_name "\t" source_val "\t" ver_val "\t" filename_clean
        }
      }
    ' 2>/dev/null || echo "")

  if [[ -z "${modules_data}" ]]; then
    echo "No module references found in Terraform files." >&2
    return 0
  fi

  # Format and print the tab-separated results as an aligned table
  echo "${modules_data}" | sort | awk -F'\t' '
    # Print the table headers and divider line
    BEGIN {
      printf "%-18s %-32s %-18s %s\n", "MODULE NAME", "SOURCE", "CURRENT VERSION", "FILE"
      printf "====================================================================================================\n"
    }
    # Print each aligned row
    {
      printf "%-18s %-32s %-18s %s\n", $1, $2, $3, $4
    }
  '
}

update_all_modules() {
  local search_dir="."
  echo "Scanning all Terraform files to discover public registry modules..." >&2

  # Extract unique public registry module names (pattern: namespace/name/provider)
  local modules
  modules=$(grep -E -o 'source[ \t]*=[ \t]*"[A-Za-z0-9_-]+/[A-Za-z0-9_-]+/[A-Za-z0-9_-]+"' -r --include="*.tf" "${search_dir}" 2>/dev/null \
    | sed -E 's/.*source[ \t]*=[ \t]*"([^"]+)".*/\1/' \
    | sort -u || echo "")

  if [[ -z "${modules}" ]]; then
    echo "No public registry modules found." >&2
    return 0
  fi

  local mod
  for mod in ${modules}; do
    echo "Processing module: ${mod}" >&2
    
    local latest_version
    if ! latest_version=$(get_latest_version "${mod}"); then
      echo "Failed to retrieve version for module ${mod}. Skipping." >&2
      continue
    fi

    # Support prefixing with "v" (common convention in this repo)
    local new_ver="v${latest_version}"
    echo "Latest version of ${mod} is ${new_ver}. Updating references..." >&2

    local file
    while IFS= read -r -d '' file; do
      awk -v new_ver="${new_ver}" -v mod="${mod}" '
      # Match a line specifying the target module source
      # e.g., source = "rancher/access/aws"
      $0 ~ "source[ \t]+=[ \t]+\"" mod "\"" {
        print # Print the source declaration line unchanged
        
        # Read the immediately following line (which should define the module version)
        getline
        
        # Check if this next line defines the version attribute
        if ($0 ~ /version[ \t]+=/) {
          # Substitute the old version string (e.g. "v1.0.0") with the new_ver string
          sub(/"[^"]+"/, "\"" new_ver "\"")
        }
        print # Print the updated (or unchanged) version line
        next
      }
      # Print all other non-matching lines unchanged
      {print}
      ' "${file}" > "${file}.tmp" && mv "${file}.tmp" "${file}"
    done < <(find "${search_dir}" -type d \( -name ".git" -o -name ".terraform" -o -name "tf_plugin_cache" \) -prune -o -type f -name "*.tf" -print0)
  done

  echo "All module updates complete." >&2
}

main() {
  local action="update"

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      --list-modules)
        action="list"
        shift
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

  if [[ "$action" == "list" ]]; then
    list_current_modules
    exit 0
  fi

  update_all_modules
}

main "$@"
