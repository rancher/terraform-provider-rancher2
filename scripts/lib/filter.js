/**
 * Filtering engine for pull requests targeting EU product changes and release criteria.
 */

export const DEFAULT_PRODUCT_PATHS = ["rancher2/"];

/**
 * Normalizes a file path by removing leading slashes and dot-slashes.
 *
 * @param {string} filePath
 * @returns {string}
 */
export function normalizePath(filePath) {
  if (!filePath || typeof filePath !== "string") {
    return "";
  }
  return filePath.replace(/^(?:\.\/|\/)+/, "").trim();
}

/**
 * Checks whether a given path falls within the allowed product prefixes.
 * Explicitly rejects examples, tests, and non-product paths.
 *
 * @param {string} filePath
 * @param {Array<string>} [allowedPrefixes=DEFAULT_PRODUCT_PATHS]
 * @returns {boolean}
 */
export function isProductPath(
  filePath,
  allowedPrefixes = DEFAULT_PRODUCT_PATHS,
) {
  const norm = normalizePath(filePath);
  if (!norm) {
    return false;
  }

  // Explicit non-product prefixes take precedence
  if (
    norm.startsWith("examples/") ||
    norm.startsWith("test/") ||
    norm.startsWith(".github/") ||
    norm.startsWith(".gemini/")
  ) {
    return false;
  }

  return allowedPrefixes.some((prefix) => {
    const cleanPrefix = normalizePath(prefix);
    return norm.startsWith(cleanPrefix);
  });
}

/**
 * Filters an array of file paths returning only qualifying product paths.
 *
 * @param {Array<string>} filePaths
 * @param {Array<string>} [allowedPrefixes=DEFAULT_PRODUCT_PATHS]
 * @returns {Array<string>}
 */
export function filterProductFiles(
  filePaths,
  allowedPrefixes = DEFAULT_PRODUCT_PATHS,
) {
  if (!Array.isArray(filePaths)) {
    return [];
  }
  return filePaths.filter((p) => isProductPath(p, allowedPrefixes));
}

/**
 * Returns true if at least one file is an EU product change.
 *
 * @param {Array<string>} filePaths
 * @param {Array<string>} [allowedPrefixes=DEFAULT_PRODUCT_PATHS]
 * @returns {boolean}
 */
export function hasProductChanges(
  filePaths,
  allowedPrefixes = DEFAULT_PRODUCT_PATHS,
) {
  if (!Array.isArray(filePaths)) {
    return false;
  }
  return filePaths.some((p) => isProductPath(p, allowedPrefixes));
}

/**
 * Extracts unique directory prefixes or root filenames from an array of file paths.
 * E.g., ["rancher2/resource_cluster.go", "docs/index.md"] -> "rancher2/, docs/"
 *
 * @param {Array<string>} files
 * @returns {string}
 */
export function extractChangesPath(files) {
  if (!Array.isArray(files) || files.length === 0) {
    return "-";
  }

  const prefixes = new Set();
  for (const f of files) {
    const clean = normalizePath(f);
    const slashIdx = clean.indexOf("/");
    if (slashIdx !== -1) {
      prefixes.add(clean.slice(0, slashIdx + 1));
    } else if (clean) {
      prefixes.add(clean);
    }
  }

  const order = ["rancher2/", "rancher/", "docs/", "website/"];
  const sorted = Array.from(prefixes).sort((a, b) => {
    const idxA = order.indexOf(a);
    const idxB = order.indexOf(b);
    if (idxA !== -1 && idxB !== -1) {
      return idxA - idxB;
    }
    if (idxA !== -1) {
      return -1;
    }
    if (idxB !== -1) {
      return 1;
    }
    return a.localeCompare(b);
  });

  return sorted.join(", ") || "-";
}

/**
 * Checks whether an author login matches any excluded username (case-insensitive).
 *
 * @param {string} authorLogin
 * @param {Array<string>} [excludedAuthors=[]]
 * @returns {boolean}
 */
export function isAuthorExcluded(authorLogin, excludedAuthors = []) {
  if (!authorLogin || typeof authorLogin !== "string") {
    return false;
  }
  const cleanAuthor = authorLogin.trim().toLowerCase();
  return excludedAuthors.some(
    (ex) => ex && ex.trim().toLowerCase() === cleanAuthor,
  );
}

/**
 * Validates whether mergedAt timestamp is within since/until boundaries.
 *
 * @param {string|Date} mergedAt
 * @param {object} [options]
 * @param {string|Date} [options.since]
 * @param {string|Date} [options.until]
 * @returns {boolean}
 */
export function isDateInRange(mergedAt, { since, until } = {}) {
  if (!mergedAt) {
    return false;
  }
  const mergeTime = new Date(mergedAt).getTime();
  if (isNaN(mergeTime)) {
    return false;
  }

  if (since) {
    const sinceTime = new Date(since).getTime();
    if (!isNaN(sinceTime) && mergeTime < sinceTime) {
      return false;
    }
  }

  if (until) {
    const untilTime = new Date(until).getTime();
    if (!isNaN(untilTime) && mergeTime > untilTime) {
      return false;
    }
  }

  return true;
}

/**
 * Checks if PR labels match any of the required labels (case-insensitive).
 *
 * @param {Array<{name: string}|string>} prLabels
 * @param {Array<string>} [requiredLabels=[]]
 * @returns {boolean}
 */
export function matchesLabels(prLabels, requiredLabels = []) {
  if (!requiredLabels || requiredLabels.length === 0) {
    return true;
  }
  if (!Array.isArray(prLabels) || prLabels.length === 0) {
    return false;
  }

  const labelNames = prLabels.map((l) =>
    (typeof l === "string" ? l : l.name || "").trim().toLowerCase(),
  );

  return requiredLabels.some((req) =>
    labelNames.includes(req.trim().toLowerCase()),
  );
}

/**
 * Checks if PR milestone matches target milestone (case-insensitive).
 *
 * @param {{title: string}|string|null} prMilestone
 * @param {string} [targetMilestone]
 * @returns {boolean}
 */
export function matchesMilestone(prMilestone, targetMilestone) {
  if (!targetMilestone) {
    return true;
  }
  if (!prMilestone) {
    return false;
  }

  const title =
    typeof prMilestone === "string" ? prMilestone : prMilestone.title || "";
  return title.trim().toLowerCase() === targetMilestone.trim().toLowerCase();
}

/**
 * Comprehensive evaluator for a pull request record.
 *
 * @param {object} pr - GitHub PR object
 * @param {object} [options]
 * @param {Array<string>} [options.filePaths] - List of changed files
 * @param {Array<string>} [options.includePaths] - Allowed path prefixes
 * @param {boolean} [options.allPaths=false] - Bypass path filtering
 * @param {Array<string>} [options.excludedAuthors=[]] - Excluded author logins
 * @param {string|Date} [options.since] - Minimum merge date
 * @param {string|Date} [options.until] - Maximum merge date
 * @param {Array<string>} [options.labels] - Filter by labels
 * @param {string} [options.milestone] - Filter by milestone
 * @returns {{ qualifies: boolean, reason?: string, productFiles: Array<string> }}
 */
export function filterPullRequest(pr, options = {}) {
  const {
    filePaths = [],
    includePaths = DEFAULT_PRODUCT_PATHS,
    allPaths = false,
    excludedAuthors = [],
    since,
    until,
    labels = [],
    milestone,
  } = options;

  if (!pr || !pr.merged_at) {
    return { qualifies: false, reason: "PR is not merged", productFiles: [] };
  }

  const authorLogin = pr.user?.login || pr.author?.login || pr.author || "";
  if (isAuthorExcluded(authorLogin, excludedAuthors)) {
    return {
      qualifies: false,
      reason: `Author '${authorLogin}' is excluded`,
      productFiles: [],
    };
  }

  if (!isDateInRange(pr.merged_at, { since, until })) {
    return {
      qualifies: false,
      reason: "Merge date out of range",
      productFiles: [],
    };
  }

  if (!matchesMilestone(pr.milestone, milestone)) {
    return {
      qualifies: false,
      reason: "Milestone does not match",
      productFiles: [],
    };
  }

  if (!matchesLabels(pr.labels, labels)) {
    return {
      qualifies: false,
      reason: "Labels do not match",
      productFiles: [],
    };
  }

  if (allPaths) {
    return { qualifies: true, productFiles: filePaths };
  }

  const productFiles = filterProductFiles(filePaths, includePaths);
  if (productFiles.length === 0) {
    return {
      qualifies: false,
      reason: "No product files modified",
      productFiles: [],
    };
  }

  return { qualifies: true, productFiles };
}
