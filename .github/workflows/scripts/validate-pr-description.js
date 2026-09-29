/**
 * Validates that the pull request description contains the required format:
 * "- Addresses: #<issue number>" or "* Addresses: #<issue number>"
 *
 * @param {string} prBody The description of the pull request.
 * @returns {string[]} The matched issue numbers.
 * @throws {Error} If no valid format line is found.
 */
export function validatePrDescription(prBody) {
  if (!prBody) {
    throw new Error('Please add a description line that matches this format: - Addresses: #<issue number>');
  }

  const cleanBody = prBody
    .replace(/<!--[\s\S]*?-->/g, '')
    .replace(/(~{3,}|`{3,})[\s\S]*?\1/g, '');
  const issueNumbers = [];
  const lines = cleanBody.split('\n');
  for (const line of lines) {
    const match = line.trim().match(/^[-*]\s+Addresses:\s+(`?)#([0-9]+)\1\.?$/i);
    if (match) {
      issueNumbers.push(match[2]);
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
 * @param {object} context The GitHub context from github-script.
 * @param {string} issueNumber The issue number to validate.
 * @returns {Promise<boolean>} True if the issue exists.
 * @throws {Error} If the issue does not exist.
 */
export async function validateIssueExists(github, context, issueNumber) {
  const { owner, repo } = context.repo;

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
