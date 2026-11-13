/**
 * Parsing utilities for pull request titles and descriptions.
 */

/**
 * Parses conventional commit information from a pull request title.
 * E.g., "feat(cluster): add proxy support" -> { type: "feat", scope: "cluster", breaking: false, description: "add proxy support" }
 *
 * @param {string} title
 * @returns {{ type: string | null, scope: string | null, breaking: boolean, description: string }}
 */
export function parseConventionalTitle(title) {
  if (!title || typeof title !== "string") {
    return { type: null, scope: null, breaking: false, description: "" };
  }

  const cleanTitle = title.trim();
  const match = cleanTitle.match(
    /^([a-zA-Z0-9_-]+)(?:\(([^)]+)\))?(!)?:\s*(.+)$/,
  );
  if (!match) {
    return {
      type: null,
      scope: null,
      breaking: false,
      description: cleanTitle,
    };
  }

  return {
    type: match[1].toLowerCase(),
    scope: match[2] ? match[2].trim() : null,
    breaking: Boolean(match[3]),
    description: match[4].trim(),
  };
}

/**
 * Extracts addressed issue numbers from a pull request description.
 * Matches:
 * - "- Addresses: #123"
 * - "Addresses #123"
 * - "Fixes: #123" / "Fixes #123"
 * - "Closes: #123" / "Closes #123"
 * - "Resolves: #123" / "Resolves #123"
 * - "- Addresses: https://github.com/rancher/terraform-provider-rancher2/issues/123"
 *
 * Scopes full GitHub issue URLs and owner/repo prefixes to the target repository,
 * ignoring links to foreign repositories.
 *
 * @param {string} body
 * @param {string|{owner?: string, repo?: string}} [targetRepo="rancher/terraform-provider-rancher2"]
 * @returns {Array<number>} Deduplicated list of issue numbers
 */
export function extractAddressedIssues(
  body,
  targetRepo = "rancher/terraform-provider-rancher2",
) {
  if (typeof body !== "string") {
    return [];
  }

  const cleanBody = body
    .replace(/<!--[\s\S]*?-->/g, "")
    .replace(/```[\s\S]*?```/g, "")
    .replace(/`[^`]+`/g, "");

  const cleanTargetRepo = (
    typeof targetRepo === "string"
      ? targetRepo
      : targetRepo?.owner && targetRepo?.repo
        ? `${targetRepo.owner}/${targetRepo.repo}`
        : "rancher/terraform-provider-rancher2"
  )
    .toLowerCase()
    .trim();

  const issueSet = new Set();

  // Pattern: Keywords followed by issue references or issue URLs
  // Added support for case-insensitivity, robust URL handling, and multiline matching
  const keywordRegex =
    /(?:Addresses|Fixes|Closes|Resolves):?\s*([\s\S]+?)(?=\n(?:Addresses|Fixes|Closes|Resolves):|\n\n|$)/gi;
  let match;
  while ((match = keywordRegex.exec(cleanBody)) !== null) {
    const segment = match[1].replace(/(github\.com\/[^\s/]+\/[^\s/]+\/issues\/\d+)[?#][^\s]*/gi, "$1");

    // 1. Full GitHub issue URLs: require matching target repository
    // Supports http/https, optional www, optional trailing slashes or anchors
    const urlRegex =
      /(?:https?:\/\/)?(?:www\.)?github\.com\/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)\/issues\/(\d+)\b/gi;
    let urlMatch;
    while ((urlMatch = urlRegex.exec(segment)) !== null) {
      const urlRepo = `${urlMatch[1]}/${urlMatch[2]}`.toLowerCase();
      if (urlRepo === cleanTargetRepo) {
        const num = parseInt(urlMatch[3], 10);
        if (!isNaN(num)) {
          issueSet.add(num);
        }
      }
    }

    // 2. Hash issue references (#123 or owner/repo#123)
    // Uses \b to ensure word boundary
    const hashRefRegex =
      /(?:^|[\s,;([{\-:*\`"'])(?:([a-zA-Z0-9_.-]+\/[a-zA-Z0-9_.-]+))?#(\d+)\b/g;
    let hashMatch;
    while ((hashMatch = hashRefRegex.exec(segment)) !== null) {
      const prefix = hashMatch[1];
      if (prefix) {
        if (prefix.toLowerCase() === cleanTargetRepo) {
          const num = parseInt(hashMatch[2], 10);
          if (!isNaN(num)) {
            issueSet.add(num);
          }
        }
      } else {
        const num = parseInt(hashMatch[2], 10);
        if (!isNaN(num)) {
          issueSet.add(num);
        }
      }
    }
  }

  return Array.from(issueSet).sort((a, b) => a - b);
}
