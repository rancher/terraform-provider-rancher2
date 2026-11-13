#!/usr/bin/env node
import { getCurrentPRNumber, getOctokit, getRepoDetails } from "./auth.js";
import { findPR } from "./find_pr.js";
import { readComments } from "./read_comments.js";

const RESOLVE_THREAD_MUTATION = `
mutation($threadId: ID!) {
  resolveReviewThread(input: {threadId: $threadId}) {
    thread {
      id
      isResolved
    }
  }
}
`;

const UNRESOLVE_THREAD_MUTATION = `
mutation($threadId: ID!) {
  unresolveReviewThread(input: {threadId: $threadId}) {
    thread {
      id
      isResolved
    }
  }
}
`;

/**
 * Resolves or unresolves an inline review thread on GitHub.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string|number} params.threadId - GraphQL node ID starting with PRRT_ or numeric comment ID
 * @param {string} [params.owner]
 * @param {string} [params.repo]
 * @param {number} [params.prNumber]
 * @param {boolean} [params.unresolve=false] - Whether to unresolve instead of resolve
 * @returns {Promise<{id: string, isResolved: boolean}>}
 */
export async function resolveThread({
  octokit,
  threadId,
  owner,
  repo,
  prNumber,
  unresolve = false,
}) {
  if (!threadId) {
    throw new Error("Missing thread ID.");
  }

  let targetThreadId = String(threadId).trim();

  // If a numeric comment ID was provided (e.g. 1001 or #1001), find its matching review thread
  if (/^#?\d+$/.test(targetThreadId)) {
    const numericCommentId = parseInt(targetThreadId.replace(/^#/, ""), 10);
    let resolvedPrNumber = prNumber;

    if (
      !resolvedPrNumber &&
      owner &&
      repo &&
      octokit.rest?.pulls?.getReviewComment
    ) {
      try {
        const { data: comment } = await octokit.rest.pulls.getReviewComment({
          owner,
          repo,
          comment_id: numericCommentId,
        });
        if (comment.pull_request_url) {
          const parts = comment.pull_request_url.split("/");
          const extractedPr = parseInt(parts[parts.length - 1], 10);
          if (!isNaN(extractedPr)) {
            resolvedPrNumber = extractedPr;
          }
        }
      } catch {
        // Fall back to existing error handling
      }
    }

    if (!owner || !repo || !resolvedPrNumber) {
      throw new Error(
        `Cannot resolve numeric comment ID #${numericCommentId} without owner, repo, and prNumber to locate thread.`,
      );
    }
    const { reviewThreads } = await readComments({
      octokit,
      owner,
      repo,
      prNumber: resolvedPrNumber,
    });
    const matchingThread = reviewThreads.find(
      (t) =>
        t.id &&
        t.comments?.nodes?.some((c) => c.databaseId === numericCommentId),
    );
    if (!matchingThread) {
      throw new Error(
        `No review thread found containing comment #${numericCommentId} on PR #${resolvedPrNumber}.`,
      );
    }
    targetThreadId = matchingThread.id;
  }

  const mutation = unresolve
    ? UNRESOLVE_THREAD_MUTATION
    : RESOLVE_THREAD_MUTATION;
  const response = await octokit.graphql(mutation, {
    threadId: targetThreadId,
  });

  const thread =
    response.unresolveReviewThread?.thread ||
    response.resolveReviewThread?.thread;
  if (!thread) {
    throw new Error(
      `Failed to ${unresolve ? "unresolve" : "resolve"} review thread ${targetThreadId}.`,
    );
  }

  return thread;
}

/**
 * Parses arguments for resolve_comment.js.
 * Supports:
 * - <thread_id_or_comment_id> [pr_number] [--unresolve]
 * - [pr_number_or_current] <thread_id_or_comment_id> [--unresolve]
 * @param {Array<string>} args
 * @returns {{threadId?: string|number, prNumber?: number, prTarget?: string, unresolve: boolean}}
 */
export function parseResolveArgs(args = []) {
  if (!args || args.length === 0) {
    return { unresolve: false };
  }
  const unresolve = args.includes("--unresolve");
  const cleanArgs = args.filter((a) => a !== "--unresolve");

  const isUrl = (v) => /\/pull\/(\d+)/.test(String(v || ""));
  const isNum = (v) => /^#?\d+$/.test(String(v || "").trim()) || isUrl(v);
  const toNum = (v) => {
    const s = String(v).trim();
    const urlMatch = s.match(/\/pull\/(\d+)/);
    if (urlMatch) {
      return parseInt(urlMatch[1], 10);
    }
    return parseInt(s.replace(/^#/, ""), 10);
  };

  let prNumber;
  let prTarget;

  const prIndex = cleanArgs.findIndex((a) => a === "-p" || a === "--pr");
  if (prIndex !== -1 && prIndex + 1 < cleanArgs.length) {
    const rawPr = cleanArgs[prIndex + 1];
    prTarget = rawPr;
    prNumber = isNum(rawPr) ? toNum(rawPr) : undefined;
    cleanArgs.splice(prIndex, 2);
  }

  if (prTarget !== undefined) {
    return {
      threadId: cleanArgs[0],
      prNumber,
      prTarget: prTarget || "current",
      unresolve,
    };
  }

  if (cleanArgs.length === 0) {
    return { unresolve };
  }

  const arg1 = cleanArgs[0];
  const arg2 = cleanArgs[1];

  if (!arg2) {
    return {
      threadId: arg1,
      prNumber: undefined,
      prTarget: "current",
      unresolve,
    };
  }

  if (arg1 === "current") {
    return {
      threadId: arg2,
      prNumber: undefined,
      prTarget: "current",
      unresolve,
    };
  }
  if (arg2 === "current") {
    return {
      threadId: arg1,
      prNumber: undefined,
      prTarget: "current",
      unresolve,
    };
  }

  const isArg1Thread = String(arg1).startsWith("PRRT_");
  const isArg2Thread = String(arg2).startsWith("PRRT_");

  if (isArg1Thread) {
    return {
      threadId: arg1,
      prNumber: isNum(arg2) ? toNum(arg2) : undefined,
      prTarget: String(toNum(arg2)),
      unresolve,
    };
  }
  if (isArg2Thread) {
    return {
      threadId: arg2,
      prNumber: isNum(arg1) ? toNum(arg1) : undefined,
      prTarget: String(toNum(arg1)),
      unresolve,
    };
  }

  // Both are numeric
  if (isNum(arg1) && isNum(arg2)) {
    const num1 = toNum(arg1);
    const num2 = toNum(arg2);

    if (isUrl(arg1) && !isUrl(arg2)) {
      return {
        threadId: num2,
        prNumber: num1,
        prTarget: String(num1),
        unresolve,
      };
    }
    if (isUrl(arg2) && !isUrl(arg1)) {
      return {
        threadId: num1,
        prNumber: num2,
        prTarget: String(num2),
        unresolve,
      };
    }

    const isArg1Hash = String(arg1).startsWith("#");
    const isArg2Hash = String(arg2).startsWith("#");
    if (isArg1Hash && !isArg2Hash) {
      return {
        threadId: num2,
        prNumber: num1,
        prTarget: String(num1),
        unresolve,
      };
    }
    if (isArg2Hash && !isArg1Hash) {
      return {
        threadId: num1,
        prNumber: num2,
        prTarget: String(num2),
        unresolve,
      };
    }

    // In GitHub, global comment IDs are significantly higher than repository PR numbers.
    // If one is greater than the other, the larger is the comment ID and the smaller is the PR number.
    if (num1 > num2) {
      return {
        threadId: num1,
        prNumber: num2,
        prTarget: String(num2),
        unresolve,
      };
    } else if (num2 > num1) {
      return {
        threadId: num2,
        prNumber: num1,
        prTarget: String(num1),
        unresolve,
      };
    }

    return {
      threadId: num1,
      prNumber: num2,
      prTarget: String(num2),
      unresolve,
    };
  }

  return {
    threadId: arg1,
    prNumber: isNum(arg2) ? toNum(arg2) : undefined,
    prTarget: arg2 || "current",
    unresolve,
  };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: resolve_comment.js <thread_id_or_comment_id> [pr_number] [--unresolve]
   or: resolve_comment.js [pr_number_or_current] <thread_id_or_comment_id> [--unresolve]

Arguments:
  thread_id_or_comment_id      GraphQL thread ID (PRRT_...) or numeric comment ID
  pr_number                    PR number or URL (required when resolving by numeric comment ID)

Options:
  -p, --pr <pr_target>         Target PR number, URL, or 'current' (unambiguous)
  --unresolve                  Unresolve the review thread instead of resolving it
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

  const {
    threadId,
    prNumber: parsedPr,
    prTarget,
    unresolve,
  } = parseResolveArgs(args);

  if (!threadId) {
    console.error(
      "Usage: resolve_comment.js <thread_id_or_comment_id> [pr_number] [--unresolve]",
    );
    console.error(
      "   or: resolve_comment.js [pr_number_or_current] <thread_id_or_comment_id> [--unresolve]",
    );
    process.exit(1);
  }

  const octokit = await getOctokit();
  let owner;
  let repo;
  let prNumber = parsedPr;

  if (/^#?\d+$/.test(String(threadId))) {
    const details = await getRepoDetails();
    owner = details.owner;
    repo = details.repo;
    if (!prNumber && prTarget && prTarget !== "current") {
      try {
        const pr = await findPR({ octokit, owner, repo, target: prTarget });
        prNumber = pr?.number;
      } catch {
        // Fall back to current branch or comment lookup
      }
    }
    if (!prNumber) {
      prNumber = await getCurrentPRNumber(octokit, { owner, repo });
    }
    if (!prNumber) {
      const numericCommentId = parseInt(String(threadId).replace(/^#/, ""), 10);
      try {
        const { data: comment } = await octokit.rest.pulls.getReviewComment({
          owner,
          repo,
          comment_id: numericCommentId,
        });
        if (comment.pull_request_url) {
          const parts = comment.pull_request_url.split("/");
          const extractedPr = parseInt(parts[parts.length - 1], 10);
          if (!isNaN(extractedPr)) {
            prNumber = extractedPr;
          }
        }
      } catch {
        // Ignore fallback failure
      }
    }
  }

  try {
    const thread = await resolveThread({
      octokit,
      threadId,
      owner,
      repo,
      prNumber,
      unresolve,
    });
    const actionDesc = unresolve ? "unresolved" : "resolved";
    console.log(
      `✅ Successfully ${actionDesc} review thread ${thread.id} (Status: ${thread.isResolved ? "RESOLVED" : "UNRESOLVED"}).`,
    );
  } catch (err) {
    console.error(`❌ Error resolving thread: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
