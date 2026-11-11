/**
 * Native fetch-based GitHub API client with pagination support.
 */

const DEFAULT_USER_AGENT = "rancher-terraform-provider-release-scripts";

/**
 * Builds standard GitHub API headers.
 *
 * @param {string} [token]
 * @returns {Record<string, string>}
 */
export function buildHeaders(token) {
  const headers = {
    Accept: "application/vnd.github+json",
    "User-Agent": DEFAULT_USER_AGENT,
  };
  if (token && token.trim()) {
    headers.Authorization = `Bearer ${token.trim()}`;
  }
  return headers;
}

/**
 * Retrieves closed pull requests targeting a base branch with automatic pagination.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {string} [options.base="main"]
 * @param {number} [options.limit=100]
 * @param {string|Date} [options.since] - Stop fetching if updated_at is older than since
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<Array<object>>}
 */
export async function fetchClosedPRs({
  owner,
  repo,
  base = "main",
  limit = 100,
  since,
  token,
  fetchFn = globalThis.fetch,
}) {
  const prs = [];
  const perPage = Math.min(limit, 100);
  let page = 1;
  const sinceTime = since ? new Date(since).getTime() : null;

  while (prs.length < limit) {
    const url = `https://api.github.com/repos/${owner}/${repo}/pulls?state=closed&base=${encodeURIComponent(base)}&sort=updated&direction=desc&per_page=${perPage}&page=${page}`;
    const res = await fetchFn(url, { headers: buildHeaders(token) });

    if (!res.ok) {
      const errText = await res.text().catch(() => "");
      throw new Error(
        `GitHub API error (${res.status}): ${errText || res.statusText}`,
      );
    }

    const items = await res.json();
    if (!Array.isArray(items) || items.length === 0) {
      break;
    }

    let reachedOlder = false;
    for (const item of items) {
      if (prs.length >= limit) {
        break;
      }
      prs.push(item);

      // Check if PR updated_at is strictly older than since threshold
      if (sinceTime && item.updated_at) {
        const itemUpdated = new Date(item.updated_at).getTime();
        if (!isNaN(itemUpdated) && itemUpdated < sinceTime) {
          reachedOlder = true;
        }
      }
    }

    if (items.length < perPage || reachedOlder) {
      break;
    }
    page++;
  }

  return prs;
}

/**
 * Fetches single pull request details by PR number.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {number} options.prNumber
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<object>}
 */
export async function fetchPullRequest({
  owner,
  repo,
  prNumber,
  token,
  fetchFn = globalThis.fetch,
}) {
  const url = `https://api.github.com/repos/${owner}/${repo}/pulls/${prNumber}`;
  const res = await fetchFn(url, { headers: buildHeaders(token) });

  if (!res.ok) {
    const errText = await res.text().catch(() => "");
    throw new Error(
      `Failed to fetch PR #${prNumber} (${res.status}): ${errText || res.statusText}`,
    );
  }

  return res.json();
}

/**
 * Fetches the changed filenames for a specific pull request.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {number} options.prNumber
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<Array<string>>}
 */
export async function fetchPRFiles({
  owner,
  repo,
  prNumber,
  token,
  fetchFn = globalThis.fetch,
}) {
  const files = [];
  let page = 1;

  while (true) {
    const url = `https://api.github.com/repos/${owner}/${repo}/pulls/${prNumber}/files?per_page=100&page=${page}`;
    const res = await fetchFn(url, { headers: buildHeaders(token) });

    if (!res.ok) {
      const errText = await res.text().catch(() => "");
      throw new Error(
        `Failed to fetch files for PR #${prNumber} (${res.status}): ${errText || res.statusText}`,
      );
    }

    const items = await res.json();
    if (!Array.isArray(items) || items.length === 0) {
      break;
    }

    for (const item of items) {
      if (item && item.filename) {
        files.push(item.filename);
      }
    }

    if (items.length < 100) {
      break;
    }
    page++;
  }

  return files;
}

/**
 * Resolves the timestamp for a release tag (e.g. "v5.2.0").
 * Checks releases first, then tags/commits.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {string} options.tag
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<string|null>} ISO date string or null
 */
export async function fetchTagDate({
  owner,
  repo,
  tag,
  token,
  fetchFn = globalThis.fetch,
}) {
  if (!tag) {
    return null;
  }

  const cleanTag = tag.trim();

  // Try release API
  try {
    const relUrl = `https://api.github.com/repos/${owner}/${repo}/releases/tags/${encodeURIComponent(cleanTag)}`;
    const relRes = await fetchFn(relUrl, { headers: buildHeaders(token) });
    if (relRes.ok) {
      const relData = await relRes.json();
      if (relData && (relData.published_at || relData.created_at)) {
        return relData.published_at || relData.created_at;
      }
    }
  } catch {
    // Fall through to git ref lookup
  }

  // Try git ref tag API
  try {
    const refUrl = `https://api.github.com/repos/${owner}/${repo}/git/ref/tags/${encodeURIComponent(cleanTag)}`;
    const refRes = await fetchFn(refUrl, { headers: buildHeaders(token) });
    if (refRes.ok) {
      const refData = await refRes.json();
      if (refData?.object?.url) {
        const objRes = await fetchFn(refData.object.url, {
          headers: buildHeaders(token),
        });
        if (objRes.ok) {
          const objData = await objRes.json();
          const date =
            objData?.tagger?.date ||
            objData?.author?.date ||
            objData?.committer?.date;
          if (date) {
            return date;
          }
        }
      }
    }
  } catch {
    // Tag not found or network error
  }

  return null;
}

/**
 * Fetches pull requests targeting a specific base branch (e.g. "release/v15").
 * Uses Core REST API with high rate limit (5,000 req/hr) instead of Search API (30 req/min).
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {string} options.baseBranch
 * @param {string} [options.state="all"]
 * @param {number} [options.limit=100]
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<Array<object>>}
 */
export async function fetchBranchPRs({
  owner,
  repo,
  baseBranch,
  state = "all",
  limit = 100,
  token,
  fetchFn = globalThis.fetch,
}) {
  const prs = [];
  let page = 1;
  const perPage = Math.min(limit, 100);

  while (prs.length < limit) {
    const url = `https://api.github.com/repos/${owner}/${repo}/pulls?state=${state}&base=${encodeURIComponent(baseBranch)}&sort=updated&direction=desc&per_page=${perPage}&page=${page}`;
    const res = await fetchFn(url, { headers: buildHeaders(token) });

    if (!res.ok) {
      // Branch might not exist or network issue; return whatever gathered
      break;
    }

    const items = await res.json();
    if (!Array.isArray(items) || items.length === 0) {
      break;
    }

    prs.push(...items);
    if (items.length < perPage) {
      break;
    }
    page++;
  }

  return prs;
}

/**
 * Fetches release branches matching "release/*" sorted descending by version.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<Array<string>>}
 */
export async function fetchReleaseBranches({
  owner,
  repo,
  token,
  fetchFn = globalThis.fetch,
}) {
  const url = `https://api.github.com/repos/${owner}/${repo}/branches?per_page=100`;
  const res = await fetchFn(url, { headers: buildHeaders(token) });

  if (!res.ok) {
    return ["release/v15", "release/v14", "release/v13"];
  }

  const branches = await res.json();
  if (!Array.isArray(branches)) {
    return ["release/v15", "release/v14", "release/v13"];
  }

  return branches
    .map((b) => b.name)
    .filter((name) => name && name.startsWith("release/"))
    .sort((a, b) => b.localeCompare(a, undefined, { numeric: true }));
}

/**
 * Fetches issues with label "internal/tracking" and builds a map of PR number to tracking issue number and state.
 * E.g. issue with body "This is the tracking issue for #2465" maps 2465 -> { number: 2466, state: "open" }.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {number} [options.limit=200]
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<Map<number, { number: number, state: string }>>}
 */
export async function fetchTrackingIssuesMap({
  owner,
  repo,
  token,
  fetchFn = globalThis.fetch,
}) {
  const prToTracking = new Map();
  const issueStates = new Map();
  let page = 1;

  while (true) {
    const url = `https://api.github.com/repos/${owner}/${repo}/issues?labels=internal/tracking&state=all&per_page=100&page=${page}`;
    const res = await fetchFn(url, { headers: buildHeaders(token) });
    if (!res.ok) {
      break;
    }

    const issues = await res.json();
    if (!Array.isArray(issues) || issues.length === 0) {
      break;
    }

    for (const iss of issues) {
      issueStates.set(iss.number, iss.state || "open");

      const match =
        (iss.body || "").match(/tracking issue for (?:PR\s*)?#?(\d+)/i) ||
        (iss.body || "").match(/<!--\s*tracking-pr:\s*#?(\d+)\s*-->/i);
      if (match) {
        const prNum = parseInt(match[1], 10);
        if (!isNaN(prNum) && !prToTracking.has(prNum)) {
          prToTracking.set(prNum, {
            number: iss.number,
            state: iss.state || "open",
          });
        }
      }
    }

    if (issues.length < 100) {
      break;
    }
    page++;
  }

  prToTracking.issueStates = issueStates;
  return prToTracking;
}

/**
 * Fetches the state (open/closed) of a specific issue.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {number} options.issueNumber
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<string>} "open" | "closed"
 */
export async function fetchIssueState({
  owner,
  repo,
  issueNumber,
  token,
  fetchFn = globalThis.fetch,
}) {
  try {
    const url = `https://api.github.com/repos/${owner}/${repo}/issues/${issueNumber}`;
    const res = await fetchFn(url, { headers: buildHeaders(token) });
    if (res.ok) {
      const data = await res.json();
      return data.state || "closed";
    }
  } catch {
    // ignore
  }
  return "closed";
}

/**
 * Fetches details of a specific issue including state, closer, and assignee.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {number} options.issueNumber
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<{ number: number, state: string, person: string, url: string }>}
 */
export async function fetchIssueDetails({
  owner,
  repo,
  issueNumber,
  token,
  fetchFn = globalThis.fetch,
}) {
  try {
    const url = `https://api.github.com/repos/${owner}/${repo}/issues/${issueNumber}`;
    const res = await fetchFn(url, { headers: buildHeaders(token) });
    if (res.ok) {
      const data = await res.json();
      const state = data.state || "open";
      const person =
        state === "closed"
          ? data.closed_by?.login ||
            data.assignee?.login ||
            data.assignees?.[0]?.login ||
            "unknown"
          : data.assignee?.login || data.assignees?.[0]?.login || "unassigned";
      return {
        number: issueNumber,
        state,
        person,
        url:
          data.html_url ||
          `https://github.com/${owner}/${repo}/issues/${issueNumber}`,
      };
    }
  } catch {
    // ignore
  }

  return {
    number: issueNumber,
    state: "unknown",
    person: "unknown",
    url: `https://github.com/${owner}/${repo}/issues/${issueNumber}`,
  };
}

/**
 * Executes a GitHub Search query for pull requests.
 *
 * @param {object} options
 * @param {string} options.query - GitHub search query (e.g. "repo:owner/repo is:pr head:backport-123-")
 * @param {string} [options.token]
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<Array<object>>}
 */
export async function searchPullRequests({
  query,
  token,
  fetchFn = globalThis.fetch,
}) {
  const items = [];
  let page = 1;

  while (true) {
    const url = `https://api.github.com/search/issues?q=${encodeURIComponent(query)}&per_page=100&page=${page}`;
    const res = await fetchFn(url, { headers: buildHeaders(token) });

    if (!res.ok) {
      const errText = await res.text().catch(() => "");
      throw new Error(
        `Search API error (${res.status}): ${errText || res.statusText}`,
      );
    }

    const data = await res.json();
    if (!data || !Array.isArray(data.items) || data.items.length === 0) {
      break;
    }

    items.push(...data.items);

    if (items.length >= data.total_count || data.items.length < 100) {
      break;
    }
    page++;
  }

  return items;
}
