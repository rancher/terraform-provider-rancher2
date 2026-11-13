import { execFile } from "node:child_process";
import fs from "node:fs/promises";
import path from "node:path";
import { promisify } from "node:util";
import { getRepoRoot } from "./utils.js";

export const sys = {
  execFileAsync: promisify(execFile),
  readFile: (filePath, encoding) => fs.readFile(filePath, encoding),
  writeFile: (filePath, content, encoding) =>
    fs.writeFile(filePath, content, encoding),
};

const WINDOWS_PATTERNS = [
  /\bpowershell\b/i,
  /\bcmd\.exe\b/i,
  /\bwindows\b/i,
  /\bwin32\b/i,
  /\.(?:bat|cmd|ps1)\b/i,
  /backslash[\s\S]{0,100}path/i,
  /drive\s+letter/i,
  /[a-zA-Z]:\\/i,
];

/**
 * Checks whether a review comment body pertains to Windows-specific compatibility.
 * Per repository standards, POSIX (Nix/Darwin/Linux) is the exclusive target environment.
 *
 * @param {string} body
 * @returns {boolean}
 */
export function isWindowsCompatibilityComment(body) {
  if (typeof body !== "string" || body.trim().length === 0) {
    return false;
  }
  return WINDOWS_PATTERNS.some((pattern) => pattern.test(body));
}

/**
 * Normalizes and extracts review threads from GitHub PR JSON comments payload.
 *
 * @param {Array<object>} rawThreads
 * @returns {Array<object>}
 */
export function normalizeReviewThreads(rawThreads = []) {
  if (!Array.isArray(rawThreads)) {
    return [];
  }

  return rawThreads
    .map((thread) => {
      const comments = Array.isArray(thread?.comments?.nodes)
        ? thread.comments.nodes
        : [];
      const firstComment = comments[0] || {};
      const threadId =
        thread?.id ||
        (firstComment.databaseId ? String(firstComment.databaseId) : null);

      return {
        threadId,
        isResolved: Boolean(thread?.isResolved),
        path: firstComment.path || "unknown",
        line: firstComment.line ?? firstComment.originalLine ?? null,
        author: firstComment.author?.login || "unknown",
        body: firstComment.body || "",
        comments: comments.map((c) => ({
          id: c.databaseId,
          author: c.author?.login || "ghost",
          body: c.body || "",
          path: c.path,
          line: c.line ?? c.originalLine ?? null,
        })),
      };
    })
    .filter((t) => Boolean(t.threadId) && t.comments.length > 0);
}

/**
 * Fetches all unresolved review comment threads for a given PR target.
 *
 * @param {string|number} [prTarget='current']
 * @param {string|null} [repoRoot=null]
 * @returns {Promise<{prNumber: number, threads: Array<object>, rawJson: object}>}
 */
export async function fetchUnresolvedPRReviewThreads(
  prTarget = "current",
  repoRoot = null,
) {
  const target = prTarget || "current";
  const root = repoRoot || (await getRepoRoot());
  const scriptPath = path.join(
    root,
    ".gemini/skills/github-pr/scripts/read_comments.js",
  );

  const args = [scriptPath, String(target), "--unresolved", "--json"];
  let stdoutData = "";
  try {
    const { stdout } = await sys.execFileAsync(process.execPath, args, {
      cwd: root,
      maxBuffer: 10 * 1024 * 1024,
      timeout: 60000,
    });
    stdoutData = stdout;

    const parsed = JSON.parse(stdout);
    const targetNum = parseInt(prTarget, 10);
    const prNumber =
      parsed?.prNumber ||
      (!Number.isNaN(targetNum) && target !== "current" ? targetNum : null);
    const threads = normalizeReviewThreads(parsed?.reviewThreads);

    return {
      prNumber,
      threads,
      rawJson: parsed,
    };
  } catch (err) {
    if (err instanceof SyntaxError) {
      throw new Error(
        `Failed to parse PR review comments JSON. stdout: ${stdoutData.slice(0, 1000)}`,
        { cause: err },
      );
    }
    const stderrMsg = err.stderr || err.message;
    throw new Error(
      `Failed to fetch unresolved PR review comments: ${stderrMsg}`,
      { cause: err },
    );
  }
}

/**
 * Filters review threads into actionable and ignored categories based on project policies.
 *
 * @param {Array<object>} threads
 * @param {object} [options]
 * @param {boolean} [options.ignoreWindows=true] - Exclude Windows compatibility suggestions
 * @returns {{actionable: Array<object>, ignored: Array<object>}}
 */
export function filterActionableReviewComments(threads = [], options = {}) {
  const ignoreWindows = options.ignoreWindows !== false;
  const actionable = [];
  const ignored = [];
  const safeThreads = Array.isArray(threads) ? threads : [];

  for (const thread of safeThreads) {
    if (thread.isResolved) {
      ignored.push({
        ...thread,
        ignoreReason: "Thread is already resolved on GitHub.",
      });
      continue;
    }

    const hasWindowsComment =
      isWindowsCompatibilityComment(thread.body) ||
      (Array.isArray(thread.comments) &&
        thread.comments.some((c) => isWindowsCompatibilityComment(c.body)));

    if (ignoreWindows && hasWindowsComment) {
      ignored.push({
        ...thread,
        ignoreReason:
          "Windows compatibility is strictly out of scope (POSIX target only).",
      });
      continue;
    }

    actionable.push(thread);
  }

  return { actionable, ignored };
}

/**
 * Generates a structured markdown remediation checklist section for the implementation plan.
 *
 * @param {object} params
 * @param {number|string} params.prNumber
 * @param {Array<object>} params.reviewThreads
 * @param {number} [params.round=1]
 * @returns {string} Formatted markdown section
 */
export function formatRemediationPlanSection({
  prNumber,
  reviewThreads = [],
  round = 1,
}) {
  const lines = [
    `## 🔄 PR #${prNumber || "UNKNOWN"} REVIEW REMEDIATION (Round ${round})`,
  ];
  const safeReviewThreads = Array.isArray(reviewThreads) ? reviewThreads : [];

  if (safeReviewThreads.length === 0) {
    lines.push("- No unresolved actionable review comments identified.");
    return lines.join("\n");
  }

  safeReviewThreads.forEach((thread, idx) => {
    const location =
      thread.line !== null && thread.line !== undefined
        ? `${thread.path}:${thread.line}`
        : thread.path;
    const author = thread.author ? `@${thread.author}` : "reviewer";
    const bodyStr = thread.body || "";
    const singleLineSummary = bodyStr
      .trim()
      .split("\n")[0]
      .replace(/^[\s#*>-]+/, "")
      .trim()
      .slice(0, 100);

    lines.push(
      `- [ ] Address Review Comment ${idx + 1} (Thread ${thread.threadId}): ${location} by ${author} - "${singleLineSummary}"`,
    );
  });

  return lines.join("\n");
}

/**
 * Appends a remediation section to the active plan file without overwriting historical context.
 *
 * @param {string} planPath
 * @param {string} remediationSection
 * @returns {Promise<{appended: boolean, reason?: string}>}
 */
export async function appendRemediationToPlan(planPath, remediationSection) {
  if (!planPath) {
    throw new Error("Plan path must be specified.");
  }
  if (!remediationSection || !remediationSection.trim()) {
    return { appended: false, reason: "Empty remediation section." };
  }

  let existingContent;
  try {
    existingContent = await sys.readFile(planPath, "utf8");
  } catch (err) {
    throw new Error(`Could not read plan at ${planPath}: ${err.message}`, {
      cause: err,
    });
  }

  const trimmedSection = remediationSection.trim();
  const headerMatch = trimmedSection.match(/^(##\s+[^\n]+)/);
  if (headerMatch && existingContent.includes(headerMatch[1])) {
    return {
      appended: false,
      reason: "Remediation section header already exists in plan.",
    };
  }

  const updatedContent = existingContent.trim()
    ? existingContent.trimEnd() + "\n\n" + trimmedSection + "\n"
    : trimmedSection + "\n";
  await sys.writeFile(planPath, updatedContent, "utf8");
  return { appended: true };
}

/**
 * Posts a reply to a GitHub review thread using the github-pr skill.
 *
 * @param {object} params
 * @param {string|number} [params.prTarget='current']
 * @param {string|number} params.threadId
 * @param {string} params.replyMessage
 * @param {string|null} [params.repoRoot=null]
 * @returns {Promise<{success: boolean, stdout: string, stderr: string}>}
 */
export async function replyToPRReviewThread({
  prTarget = "current",
  threadId,
  replyMessage,
  repoRoot = null,
}) {
  if (!threadId) {
    throw new Error("Thread ID is required to reply to a review thread.");
  }
  if (!replyMessage || !replyMessage.trim()) {
    throw new Error("Reply message cannot be empty.");
  }

  const target = prTarget || "current";
  const root = repoRoot || (await getRepoRoot());
  const scriptPath = path.join(
    root,
    ".gemini/skills/github-pr/scripts/reply-comment.js",
  );
  const args = [
    scriptPath,
    String(threadId),
    "--pr",
    String(target),
    "-m",
    replyMessage.trim(),
  ];

  try {
    const { stdout, stderr } = await sys.execFileAsync(process.execPath, args, {
      cwd: root,
      maxBuffer: 10 * 1024 * 1024,
      timeout: 60000,
    });
    return { success: true, stdout, stderr };
  } catch (err) {
    const msg = err.stderr?.trim() || err.stdout?.trim() || err.message;
    throw new Error(`Failed to reply to PR review thread ${threadId}: ${msg}`, {
      cause: err,
    });
  }
}

/**
 * Resolves a GitHub review thread using the github-pr skill.
 *
 * @param {object} params
 * @param {string|number} [params.prTarget='current']
 * @param {string|number} params.threadId
 * @param {string|null} [params.repoRoot=null]
 * @returns {Promise<{success: boolean, stdout: string, stderr: string}>}
 */
export async function resolvePRReviewThread({
  prTarget = "current",
  threadId,
  repoRoot = null,
}) {
  if (!threadId) {
    throw new Error("Thread ID is required to resolve a review thread.");
  }

  const target = prTarget || "current";
  const root = repoRoot || (await getRepoRoot());
  const scriptPath = path.join(
    root,
    ".gemini/skills/github-pr/scripts/resolve_comment.js",
  );
  const args = [scriptPath, String(threadId), "--pr", String(target)];

  try {
    const { stdout, stderr } = await sys.execFileAsync(process.execPath, args, {
      cwd: root,
      maxBuffer: 10 * 1024 * 1024,
      timeout: 60000,
    });
    return { success: true, stdout, stderr };
  } catch (err) {
    const msg = err.stderr?.trim() || err.stdout?.trim() || err.message;
    throw new Error(`Failed to resolve PR review thread ${threadId}: ${msg}`, {
      cause: err,
    });
  }
}

/**
 * Builds the implementation prompt for addressing PR review comments.
 *
 * @param {object} params
 * @param {number|string} params.prNumber
 * @param {Array<object>} params.threads - Actionable review threads
 * @param {string} params.activePlan - Full active plan content
 * @returns {string} Formatted implementation prompt
 */
export function buildPRReviewImplementationPrompt({
  prNumber,
  threads = [],
  activePlan = "",
}) {
  const safeThreads = Array.isArray(threads) ? threads : [];
  const formattedComments = safeThreads
    .map((t, idx) => {
      const loc =
        t.line !== null && t.line !== undefined
          ? `${t.path}:${t.line}`
          : t.path;
      return `### Comment ${idx + 1} (Thread ID: \`${t.threadId}\`)\n- File: \`${loc}\`\n- Reviewer: @${t.author}\n\n${t.body}\n`;
    })
    .join("\n");

  const prTargetStr = prNumber || "current";

  return `Objective: "Address PR review comments on PR #${prTargetStr}".

There are ${safeThreads.length} unresolved PR review comment(s) to assess, implement, and resolve:

<pr_review_threads>
${formattedComments}
</pr_review_threads>

Instructions & Protocols:
1. **Assess & Implement Surgically**:
   - Assess each review comment against the active plan in <active_plan>.
   - Work strictly in short iterative development loops (e.g. 2 items at a time).
   - Apply clean, idiomatic code modifications directly addressing the root-cause of the findings.
   - Ignore any Windows-specific compatibility suggestions (Windows is strictly out of scope; POSIX target only).
2. **Deterministic Validation**:
   - Run local validation (e.g. '${process.execPath} .gemini/skills/code-validation/scripts/test.js all' and linters) to verify your changes pass cleanly.
3. **Reply & Resolve Threads via Skills**:
   - After successfully addressing and verifying a comment, post a concise reply explaining the resolution:
     \`${process.execPath} .gemini/skills/github-pr/scripts/reply-comment.js <thread_id> --pr ${prTargetStr} -m "Addressed in latest update: <brief explanation>"\`
   - Resolve the thread on GitHub:
     \`${process.execPath} .gemini/skills/github-pr/scripts/resolve_comment.js <thread_id> --pr ${prTargetStr}\`
4. **Plan Integrity**:
   - Update the checklist items in 'plans/current.md' as you complete them.
   - Do NOT wipe or overwrite 'plans/current.md'.
5. **Git Boundaries**:
   - Do NOT commit or push changes. Staging, QA review loop, and commits are managed by the orchestrator phases.

<active_plan>
${activePlan}
</active_plan>`;
}
