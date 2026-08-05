#!/usr/bin/env bash
set -euo pipefail

update_module_version() {
    local module_name="rancher/access/aws"
    local search_dir="."
    
    echo "Fetching latest version of ${module_name} from Terraform registry..."
    local latest_version
    latest_version=$(curl -s "https://registry.terraform.io/v1/modules/${module_name}" | jq -r '.version')

    if [[ -z "${latest_version}" || "${latest_version}" == "null" ]]; then
        echo "Failed to fetch the latest version from the Terraform registry." >&2
        return 1
    fi

    local new_ver="v${latest_version}"
    echo "Latest version is ${new_ver}. Updating instances..."

    while IFS= read -r -d '' file; do
        awk -v new_ver="${new_ver}" -v mod="${module_name}" '
        $0 ~ "source[ \t]+=[ \t]+\"" mod "\"" {
            print
            getline
            if ($0 ~ /version[ \t]+=/) {
                sub(/"[^"]+"/, "\"" new_ver "\"")
            }
            print
            next
        }
        {print}
        ' "${file}" > "${file}.tmp" && mv "${file}.tmp" "${file}"
    done < <(find "${search_dir}" -type d \( -name ".git" -o -name ".terraform" -o -name "tf_plugin_cache" \) -prune -o -type f -name "*.tf" -print0)

    echo "Update complete."
}

main() {
    update_module_version
}

main "$@"