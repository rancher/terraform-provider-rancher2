/**
 * Validates that the pull request description contains the required format, e.g.:
 * - Addresses: #<issue number>
 * * Addresses: #<issue number>
 * - Addresses: <org>/<repo>#<issue number>
 *
 * @param {string} owner The owner to be assumed for issue references that don't specify one
 * @param {string} repo The repo to be assumed for issue references that don't specify one
 * @param {string} prBody The description of the pull request.
 * @returns {obj[]} Array of objects of the matched owner, repo and issue number
 * @throws {Error} If no valid format line is found.
 */
export function validatePrDescription(owner, repo, prBody) {
  if (!prBody) {
    throw new Error('Please add a description line that matches this format: - Addresses: #<issue number>');
  }

  const cleanBody = prBody
    .replace(/<!--[\s\S]*?-->/g, '')
    .replace(/(~{3,}|`{3,})[\s\S]*?\1/g, '');
  const issueNumbers = [];
  const lines = cleanBody.split('\n');
  for (const line of lines) {
    const match = line.trim().match(/^\s*[-*]\s+Addresses:\s+(`?)(([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+))?#([0-9]+)\1(\.?)$/im);
    if (match) {
      issueNumbers.push({
        owner: match[3] ? match[3] : owner,
        repo: match[4] ? match[4] : repo,
        issue: match[5],
      });
    }
  }

  if (issueNumbers.length === 0) {
    throw new Error('Please add a description line that matches this format: - Addresses: #<issue number>');
  }

  return issueNumbers;
}

/**
 * Validates that the issue exists in the repository.
 *
 * @param {object} github The Octokit GitHub client from github-script.
 * @param {string} owner The GitHub owner of the issue
 * @param {string} repo The GitHub repo of the issue
 * @param {string} issueNumber The issue number to validate.
 * @returns {Promise<boolean>} True if the issue exists.
 * @throws {Error} If the issue does not exist.
 */
export async function validateIssueExists(github, owner, repo, issueNumber) {
  let issue;
  try {
    ({ data: issue } = await github.rest.issues.get({
      owner,
      repo,
      issue_number: parseInt(issueNumber, 10),
    }));
  } catch (error) {
    if (error.status === 404) {
      throw new Error(`Issue #${issueNumber} not found in ${owner}/${repo} repo.`);
    }
    throw error;
  }

  if (issue.pull_request) {
    throw new Error(`Issue #${issueNumber} is a pull request, not an issue.`);
  }

  if (issue.state !== 'open') {
    throw new Error(`Issue #${issueNumber} is already closed.`);
  }

  return true;
}

/**
 * Main entrypoint for GitHub Actions execution.
 *
 * @param {object} params The github-script parameters.
 * @param {object} params.github The Octokit GitHub client from github-script.
 * @param {object} params.context The GitHub context from github-script.
 * @param {object} params.core The GitHub Actions core toolkit.
 */
export default async function run({ github, context, core }) {
  try {
    const prBody = context.payload?.pull_request?.body ?? '';
    const {owner, repo } = context.repo ?? {owner: 'rancher', repo: 'terraform-provider-rancher2'};
    const issues = validatePrDescription(owner, repo, prBody);
    const issueMap = new Map(
      issues.map((issue) => { return [JSON.stringify(issue), issue]; }),
    );
    const uniqueIssues = [...issueMap.values()];
    const results = await Promise.allSettled(
      uniqueIssues.map(({owner, repo, issue}) => validateIssueExists(github, owner, repo, issue))
    );
    const errors = results
      .filter((r) => r.status === 'rejected')
      .map((r) => r.reason?.message ?? String(r.reason));
    if (errors.length > 0) {
      throw new Error(errors.join('\n'));
    }
    core?.info?.(`Successfully validated issue reference(s): ${uniqueIssues.map((n) => `#${n}`).join(', ')}`);
  } catch (error) {
    core.setFailed(error?.message ?? String(error));
  }
}
