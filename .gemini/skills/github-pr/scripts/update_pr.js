#!/usr/bin/env node
import { getOctokit, getRepoDetails } from "./auth.js";
import { findPR } from "./find_pr.js";

/**
 * Updates a pull request's title, body, and/or state.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string} params.owner
 * @param {string} params.repo
 * @param {number} params.prNumber
 * @param {string} [params.title]
 * @param {string} [params.body]
 * @param {string} [params.state] - 'open' or 'closed'
 * @returns {Promise<object>}
 */
export async function updatePR({
  octokit,
  owner,
  repo,
  prNumber,
  title,
  body,
  state,
}) {
  const updatePayload = {
    owner,
    repo,
    pull_number: prNumber,
  };

  let hasUpdate = false;
  if (title !== undefined && title !== null && title !== "") {
    updatePayload.title = title;
    hasUpdate = true;
  }
  if (body !== undefined && body !== null && body !== "") {
    updatePayload.body = body;
    hasUpdate = true;
  }
  if (state !== undefined && state !== null && state !== "") {
    updatePayload.state = state;
    hasUpdate = true;
  }

  if (!hasUpdate) {
    throw new Error("No title or body updates provided.");
  }

  const { data } = await octokit.rest.pulls.update(updatePayload);
  return data;
}

/**
 * Parses command-line arguments for update_pr.js flexibly.
 * Supports:
 * - [pr_target, title, body]
 * - [pr_target, title]
 * - [title, body]
 * - [title]
 * - [--title title] [--body body] [--state open|closed] [pr_target]
 * - [-t title] [-b body] [-s open|closed] [pr_target]
 * @param {Array<string>} args
 * @returns {{prTarget?: string, title?: string, body?: string, state?: string}}
 */
export function parseUpdateArgs(args) {
  if (!args || args.length === 0) {
    return {};
  }

  const hasFlags = args.some(
    (a) =>
      a === "-t" ||
      a === "--title" ||
      a === "-b" ||
      a === "--body" ||
      a === "-s" ||
      a === "--state" ||
      a.startsWith("-t=") ||
      a.startsWith("--title=") ||
      a.startsWith("-b=") ||
      a.startsWith("--body=") ||
      a.startsWith("-s=") ||
      a.startsWith("--state="),
  );

  if (hasFlags) {
    let prTarget = undefined;
    let title = undefined;
    let body = undefined;
    let state = undefined;

    for (let i = 0; i < args.length; i++) {
      const a = args[i];
      if ((a === "-t" || a === "--title") && i + 1 < args.length) {
        title = args[++i];
      } else if (a.startsWith("-t=") || a.startsWith("--title=")) {
        title = a.slice(a.indexOf("=") + 1);
      } else if ((a === "-b" || a === "--body") && i + 1 < args.length) {
        body = args[++i];
      } else if (a.startsWith("-b=") || a.startsWith("--body=")) {
        body = a.slice(a.indexOf("=") + 1);
      } else if ((a === "-s" || a === "--state") && i + 1 < args.length) {
        state = args[++i];
      } else if (a.startsWith("-s=") || a.startsWith("--state=")) {
        state = a.slice(a.indexOf("=") + 1);
      } else if (!prTarget && !a.startsWith("-")) {
        prTarget = a;
      }
    }
    const result = { prTarget: prTarget || "current", title, body };
    if (state !== undefined) {
      result.state = state;
    }
    return result;
  }

  const isNum = (v) =>
    /^#?\d+$/.test(String(v || "").trim()) ||
    /\/pull\/(\d+)/.test(String(v || ""));

  if (args.length === 1) {
    return { prTarget: "current", title: args[0], body: undefined };
  }

  if (args.length === 2) {
    if (args[0] === "current" || isNum(args[0])) {
      return { prTarget: args[0], title: args[1], body: undefined };
    }
    return { prTarget: "current", title: args[0], body: args[1] };
  }

  return {
    prTarget: args[0],
    title: args[1],
    body: args.slice(2).join(" "),
  };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: update_pr.js [pr_number_or_current] <title> [body]
   or: update_pr.js [pr_number_or_current] [--title title] [--body body] [--state open|closed]

Arguments:
  pr_number_or_current         PR number, URL, or 'current' (default: 'current')
  title                        New PR title
  body                         New PR body text

Options:
  -t, --title <text>           Set the pull request title
  -b, --body <text>            Set the pull request body
  -s, --state <open|closed>    Update pull request state ('open' or 'closed')
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

  const { prTarget, title, body, state } = parseUpdateArgs(args);

  if (title === undefined && body === undefined && state === undefined) {
    console.error(
      "Usage: update_pr.js [pr_number_or_current] <title_or_empty> [body_or_empty]",
    );
    console.error(
      "   or: update_pr.js [pr_number_or_current] [--title title] [--body body] [--state open|closed]",
    );
    process.exit(1);
  }

  const { owner, repo } = await getRepoDetails();
  const octokit = await getOctokit();

  let prNumber;
  try {
    const pr = await findPR({
      octokit,
      owner,
      repo,
      target: prTarget === "current" ? undefined : prTarget,
    });
    prNumber = pr.number;
  } catch (err) {
    console.error(`❌ Error: ${err.message}`);
    process.exit(1);
  }

  try {
    const updated = await updatePR({
      octokit,
      owner,
      repo,
      prNumber,
      title,
      body,
      state,
    });
    console.log(
      `✅ Successfully updated PR #${updated.number}: "${updated.title}" (State: ${updated.state})`,
    );
  } catch (err) {
    console.error(`❌ Error updating PR: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
