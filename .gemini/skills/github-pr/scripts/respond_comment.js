#!/usr/bin/env node
import { getOctokit, getRepoDetails } from "./auth.js";
import { findPR } from "./find_pr.js";
import { replyComment } from "./reply-comment.js";

/**
 * Adds a general comment or replies to an inline comment on a pull request.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string} params.owner
 * @param {string} params.repo
 * @param {number} params.prNumber
 * @param {string|number} [params.commentId] - numeric ID or 'general' / '0'
 * @param {string} params.body
 * @returns {Promise<object>}
 */
export async function respondComment({
  octokit,
  owner,
  repo,
  prNumber,
  commentId,
  body,
}) {
  if (!body || !body.trim()) {
    throw new Error("Comment body cannot be empty.");
  }

  const isGeneral =
    !commentId ||
    commentId === "general" ||
    commentId === "0" ||
    commentId === 0;

  if (isGeneral) {
    const { data } = await octokit.rest.issues.createComment({
      owner,
      repo,
      issue_number: prNumber,
      body,
    });
    return { type: "general_comment", data };
  }

  return replyComment({
    octokit,
    owner,
    repo,
    prNumber,
    commentId:
      typeof commentId === "string" && commentId.startsWith("PRRT_")
        ? commentId
        : parseInt(commentId, 10),
    body,
  });
}

/**
 * Parses command-line arguments for respond_comment.js flexibly.
 * Supports:
 * - [pr_target, comment_id, ...body]
 * - [pr_target, 'general', ...body]
 * - ['general', ...body]
 * - [pr_target, ...body]
 * - [...body]
 * @param {Array<string>} args
 * @returns {{prTarget?: string, commentId?: string|number, body?: string}}
 */
export function parseRespondArgs(args) {
  if (!args || args.length === 0) {
    return {};
  }

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
  const isRef = (v) => isNum(v) || String(v || "").startsWith("PRRT_");
  const toRef = (v) =>
    String(v || "").startsWith("PRRT_") ? String(v).trim() : toNum(v);

  let prTarget = "current";
  let commentId = "general";
  let bodyTokens;

  let cleanArgs = [...args];
  let explicitBody = null;
  const mIndex = cleanArgs.findIndex((a) => a === "-m" || a === "--body");
  if (mIndex !== -1 && mIndex + 1 < cleanArgs.length) {
    explicitBody = cleanArgs.slice(mIndex + 1).join(" ");
    cleanArgs = cleanArgs.slice(0, mIndex);
  }

  if (cleanArgs.length === 0) {
    return { prTarget, commentId, body: explicitBody || "" };
  }

  const first = cleanArgs[0];

  if (first === "general" || first === "0") {
    commentId = "general";
    bodyTokens = cleanArgs.slice(1);
  } else if (String(first || "").startsWith("PRRT_")) {
    commentId = String(first).trim();
    bodyTokens = cleanArgs.slice(1);
  } else if (first === "current") {
    prTarget = "current";
    if (
      cleanArgs.length >= 2 &&
      (cleanArgs[1] === "general" || isRef(cleanArgs[1]))
    ) {
      commentId = isRef(cleanArgs[1]) ? toRef(cleanArgs[1]) : "general";
      bodyTokens = cleanArgs.slice(2);
    } else {
      commentId = "general";
      bodyTokens = cleanArgs.slice(1);
    }
  } else if (isNum(first)) {
    if (
      cleanArgs.length >= 2 &&
      (cleanArgs[1] === "general" || isRef(cleanArgs[1]))
    ) {
      const second = cleanArgs[1];
      if (isUrl(first) && !isUrl(second)) {
        prTarget = String(toNum(first));
        commentId = isRef(second) ? toRef(second) : "general";
      } else if (isUrl(second) && !isUrl(first)) {
        prTarget = String(toNum(second));
        commentId = toNum(first);
      } else if (isNum(second) && toNum(first) > toNum(second)) {
        prTarget = String(toNum(second));
        commentId = toNum(first);
      } else {
        prTarget = isUrl(first) ? String(toNum(first)) : first;
        commentId = isRef(second) ? toRef(second) : "general";
      }
      bodyTokens = cleanArgs.slice(2);
    } else {
      prTarget = isUrl(first) ? String(toNum(first)) : first;
      commentId = "general";
      bodyTokens = cleanArgs.slice(1);
    }
  } else {
    if (
      cleanArgs.length >= 2 &&
      (cleanArgs[1] === "general" || isRef(cleanArgs[1]))
    ) {
      prTarget = first;
      commentId = isRef(cleanArgs[1]) ? toRef(cleanArgs[1]) : "general";
      bodyTokens = cleanArgs.slice(2);
    } else {
      prTarget = "current";
      commentId = "general";
      bodyTokens = cleanArgs;
    }
  }

  const body = (explicitBody || bodyTokens.join(" ")).trim();
  return { prTarget, commentId, body };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: respond_comment.js [pr_number_or_current] [comment_id_or_general] "body"
   or: respond_comment.js [pr_number_or_current] [comment_id_or_general] -m "body"

Arguments:
  pr_number_or_current         PR number, URL, or 'current' (default: 'current')
  comment_id_or_general        Inline comment ID, GraphQL thread ID (PRRT_...), or 'general' (default: 'general')
  body                         Comment body text

Options:
  -m, --body <text>            Specify comment body explicitly
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

  const { prTarget, commentId, body } = parseRespondArgs(args);

  if (!body) {
    console.error(
      'Usage: respond_comment.js [pr_number_or_current] [comment_id_or_general] "body"',
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
    if (!String(commentId).startsWith("PRRT_")) {
      console.error(`❌ Error: ${err.message}`);
      process.exit(1);
    }
  }

  try {
    const result = await respondComment({
      octokit,
      owner,
      repo,
      prNumber,
      commentId,
      body,
    });
    const targetDesc = prNumber
      ? `on PR #${prNumber}`
      : `to thread ${commentId}`;
    console.log(`✅ Successfully responded ${targetDesc} (${result.type}).`);
  } catch (err) {
    console.error(`❌ Error responding to PR comment: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
