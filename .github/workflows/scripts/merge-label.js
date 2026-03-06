export default async ({ github, context, core}) => {
  // Context for this script
  // https://github.com/actions/github-script?tab=readme-ov-file#this-action
  // https://octokit.github.io/rest.js/v22/#custom-requests replace octokit with github in the examples
  // https://github.com/actions/toolkit/tree/main/packages/core
  // https://docs.github.com/en/actions/reference/workflows-and-actions/contexts

  // this script should find the issue associated with this PR and add a label "internal/merged" to that issue
  const owner = "rancher";
  const repo = "terraform-provider-rancher2";
  const pr = context.payload.pull_request;

  const issueRegex = /#(\d+)/g;
  const matches = pr.body.matchAll(issueRegex);
  const issueNumbers = Array.from(matches, m => parseInt(m[1]));

  core.info(`Found issue numbers in PR body: ${issueNumbers}`);

  for (const issueNumber of issueNumbers) {
    try {
      // Check if it's an issue (not a PR) and has the internal/backport label
      const { data: issueData } = await github.rest.issues.get({
        owner,
        repo,
        issue_number: issueNumber,
      });

      if (!issueData.pull_request && issueData.labels.some(l => l.name === 'internal/backport')) {
        core.info(`Adding 'internal/merged' label to issue #${issueNumber}`);
        await github.rest.issues.addLabels({
          owner,
          repo,
          issue_number: issueNumber,
          labels: ["internal/merged"]
        });
      }
    } catch (error) {
      core.info(`Could not process issue #${issueNumber}: ${error.message}`);
    }
  }
};
