#!/usr/bin/env node
import { getCurrentPRNumber, getOctokit, getRepoDetails } from "./auth.js";

/**
 * Fetches pull request details via Octokit.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string} params.owner
 * @param {string} params.repo
 * @param {string|number} [params.target] - PR number or branch name
 * @returns {Promise<object>} PR data
 */
export async function findPR({ octokit, owner, repo, target }) {
  let pullNumber;
  const trimmedTarget =
    target !== undefined && target !== null ? String(target).trim() : "";

  if (
    trimmedTarget &&
    trimmedTarget !== "current" &&
    trimmedTarget !== "HEAD"
  ) {
    const urlMatch = trimmedTarget.match(/\/pull\/(\d+)/);
    if (/^#?\d+$/.test(trimmedTarget)) {
      pullNumber = parseInt(trimmedTarget.replace(/^#/, ""), 10);
    } else if (urlMatch) {
      pullNumber = parseInt(urlMatch[1], 10);
    } else {
      const headSpec = trimmedTarget.includes(":")
        ? trimmedTarget
        : `${owner}:${trimmedTarget}`;
      // Search by head branch, prioritizing open PRs over closed ones
      const { data: pulls } = await octokit.rest.pulls.list({
        owner,
        repo,
        head: headSpec,
        state: "all",
        per_page: 10,
      });
      const match =
        (pulls || []).find((p) => p.state === "open") || (pulls || [])[0];
      if (match) {
        pullNumber = match.number;
      } else {
        // Fallback: search open/all PRs for matching head.ref across forks
        const allPulls = octokit.paginate
          ? await octokit.paginate(octokit.rest.pulls.list, {
            owner,
            repo,
            state: "all",
            per_page: 100,
          })
          : (
            await octokit.rest.pulls.list({
              owner,
              repo,
              state: "all",
              per_page: 100,
            })
          ).data || [];

        const refMatch =
          allPulls.find(
            (p) => p.state === "open" && p.head?.ref === trimmedTarget,
          ) || allPulls.find((p) => p.head?.ref === trimmedTarget);
        if (refMatch) {
          pullNumber = refMatch.number;
        } else {
          throw new Error(
            `No pull request found for branch '${target}' in ${owner}/${repo}.`,
          );
        }
      }
    }
  } else {
    pullNumber = await getCurrentPRNumber(octokit, { owner, repo });
    if (!pullNumber) {
      throw new Error(
        "No pull request target specified and no PR found for current branch.",
      );
    }
  }

  const { data: pr } = await octokit.rest.pulls.get({
    owner,
    repo,
    pull_number: pullNumber,
  });

  return pr;
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: find_pr.js [pr_number_or_branch]

Arguments:
  pr_number_or_branch          PR number, PR URL, or branch name (default: current branch)

Options:
  -h, --help                   Show this help message

Environment Variables:
  GH_REPO / GITHUB_REPOSITORY  Target repository in 'owner/repo' format (overrides git remote detection)
  GITHUB_TOKEN / GH_TOKEN      GitHub authentication token (falls back to 'gh auth token')
`);
}

async function main() {
  const args = process.argv.slice(2);
  if (args.includes("-h") || args.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  const target = args[0];
  const { owner, repo } = await getRepoDetails();
  const octokit = await getOctokit();

  try {
    const pr = await findPR({ octokit, owner, repo, target });
    const stateDesc = pr.merged
      ? "MERGED"
      : pr.draft
        ? "DRAFT"
        : pr.state.toUpperCase();
    console.log(`=== Pull Request #${pr.number}: ${pr.title} ===`);
    console.log(`State:  ${stateDesc}`);
    console.log(`Author: ${pr.user?.login}`);
    console.log(`Branch: ${pr.head.ref} -> ${pr.base.ref}`);
    console.log(`URL:    ${pr.html_url}`);
    console.log(
      `\n--- Description ---\n${pr.body || "(No description provided)"}`,
    );
  } catch (err) {
    console.error(`❌ Error: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
