#!/usr/bin/env node
import { getOctokit, getRepoDetails } from "./auth.js";
import { findPR } from "./find_pr.js";

const REVIEW_THREADS_QUERY = `
query($owner: String!, $repo: String!, $prNumber: Int!, $threadCursor: String) {
  repository(owner: $owner, name: $repo) {
    pullRequest(number: $prNumber) {
      reviewThreads(first: 50, after: $threadCursor) {
        pageInfo {
          hasNextPage
          endCursor
        }
        nodes {
          id
          isResolved
          comments(first: 50) {
            pageInfo {
              hasNextPage
              endCursor
            }
            nodes {
              databaseId
              author {
                login
              }
              body
              path
              line
              originalLine
            }
          }
        }
      }
    }
  }
}
`;

const THREAD_COMMENTS_QUERY = `
query($threadId: ID!, $commentCursor: String) {
  node(id: $threadId) {
    ... on PullRequestReviewThread {
      comments(first: 50, after: $commentCursor) {
        pageInfo {
          hasNextPage
          endCursor
        }
        nodes {
          databaseId
          author {
            login
          }
          body
          path
          line
          originalLine
        }
      }
    }
  }
}
`;

/**
 * Reads all conversation and review comments for a PR.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string} params.owner
 * @param {string} params.repo
 * @param {number|string} params.prNumber
 * @param {boolean} [params.unresolvedOnly=false]
 * @returns {Promise<{conversationComments: Array, reviewThreads: Array}>}
 */
export async function readComments({
  octokit,
  owner,
  repo,
  prNumber,
  unresolvedOnly = false,
}) {
  let numericPrNumber = prNumber;
  if (typeof prNumber === "string") {
    const urlMatch = prNumber.match(/\/pull\/(\d+)/);
    if (urlMatch) {
      numericPrNumber = parseInt(urlMatch[1], 10);
    } else {
      numericPrNumber = parseInt(prNumber.replace(/^#/, ""), 10);
    }
  }

  // 1. Fetch conversation comments (paginated)
  let conversationComments = [];
  if (octokit.paginate && octokit.rest?.issues?.listComments) {
    conversationComments = await octokit.paginate(
      octokit.rest.issues.listComments,
      {
        owner,
        repo,
        issue_number: numericPrNumber,
        per_page: 100,
      },
    );
  } else if (octokit.rest?.issues?.listComments) {
    const res = await octokit.rest.issues.listComments({
      owner,
      repo,
      issue_number: numericPrNumber,
      per_page: 100,
    });
    conversationComments = res.data || [];
  }

  // 2. Fetch review threads via GraphQL with cursor pagination
  const reviewThreads = [];
  try {
    let threadCursor = null;
    let hasNextThreadPage = true;
    const gqlThreads = [];

    while (hasNextThreadPage) {
      const gqlRes = await octokit.graphql(REVIEW_THREADS_QUERY, {
        owner,
        repo,
        prNumber: numericPrNumber,
        threadCursor,
      });

      const threadConnection = gqlRes.repository?.pullRequest?.reviewThreads;
      const threads = threadConnection?.nodes || [];

      for (const thread of threads) {
        let commentCursor = thread.comments?.pageInfo?.endCursor;
        let hasNextCommentPage = thread.comments?.pageInfo?.hasNextPage;

        while (hasNextCommentPage) {
          const commentRes = await octokit.graphql(THREAD_COMMENTS_QUERY, {
            threadId: thread.id,
            commentCursor,
          });
          const commentConnection = commentRes.node?.comments;
          const moreComments = commentConnection?.nodes || [];
          if (moreComments.length > 0) {
            thread.comments.nodes.push(...moreComments);
          }
          hasNextCommentPage =
            commentConnection?.pageInfo?.hasNextPage || false;
          commentCursor = commentConnection?.pageInfo?.endCursor;
        }

        gqlThreads.push(thread);
      }

      hasNextThreadPage = threadConnection?.pageInfo?.hasNextPage || false;
      threadCursor = threadConnection?.pageInfo?.endCursor;
    }

    reviewThreads.push(...gqlThreads);
  } catch {
    // If GraphQL query fails, reset collection and fallback to REST review comments
    reviewThreads.length = 0;
    const restReviewComments = octokit.paginate
      ? await octokit.paginate(octokit.rest.pulls.listReviewComments, {
        owner,
        repo,
        pull_number: numericPrNumber,
        per_page: 100,
      })
      : (
        await octokit.rest.pulls.listReviewComments({
          owner,
          repo,
          pull_number: numericPrNumber,
          per_page: 100,
        })
      ).data || [];

    for (const c of restReviewComments) {
      reviewThreads.push({
        id: null,
        isResolved: false,
        isResolvable: false,
        comments: {
          nodes: [
            {
              databaseId: c.id,
              author: { login: c.user?.login || "ghost" },
              body: c.body,
              path: c.path,
              line: c.line || c.original_line,
              originalLine: c.original_line,
            },
          ],
        },
      });
    }
  }

  const filteredThreads = unresolvedOnly
    ? reviewThreads.filter((t) => !t.isResolved)
    : reviewThreads;
  return { conversationComments, reviewThreads: filteredThreads };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: read_comments.js [pr_number_or_branch] [--unresolved] [--json]

Options:
  --unresolved  Filter to only unresolved review comment threads
  --json        Output the results as JSON
  -h, --help    Show this help message

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

  const unresolvedOnly = args.includes("--unresolved");
  const jsonOutput = args.includes("--json");
  const targetArgs = args.filter((a) => !a.startsWith("-"));
  const argTarget = targetArgs[0];

  const { owner, repo } = await getRepoDetails();
  const octokit = await getOctokit();

  let prNumber;
  try {
    const target = argTarget ? argTarget.replace(/^#/, "") : undefined;
    const pr = await findPR({
      octokit,
      owner,
      repo,
      target,
    });
    prNumber = pr.number;
  } catch (err) {
    const targetDesc = argTarget ? `'${argTarget}'` : "for current branch";
    console.error(
      `❌ Error finding PR ${targetDesc} in repository ${owner}/${repo}: ${err.message}`,
    );
    if (err.status === 404) {
      console.error(
        `💡 Tip: Verify that PR ${targetDesc} exists in ${owner}/${repo}. ` +
          `If working from a fork, set an 'upstream' git remote or export GH_REPO="${owner}/${repo}".`,
      );
    }
    process.exit(1);
  }

  try {
    const { conversationComments, reviewThreads } = await readComments({
      octokit,
      owner,
      repo,
      prNumber,
      unresolvedOnly,
    });

    if (jsonOutput) {
      console.log(
        JSON.stringify(
          { prNumber, conversationComments, reviewThreads },
          null,
          2,
        ),
      );
      return;
    }

    console.log(`=== Comments for PR #${prNumber} in ${owner}/${repo} ===\n`);

    console.log(
      `--- Conversation Comments (${conversationComments.length}) ---`,
    );
    if (conversationComments.length === 0) {
      console.log("No general conversation comments found.");
    } else {
      for (const comment of conversationComments) {
        console.log(
          `\n[Comment #${comment.id}] @${comment.user?.login || "ghost"}:`,
        );
        console.log(comment.body);
      }
    }

    console.log(
      `\n--- Review Comment Threads (${reviewThreads.length}${unresolvedOnly ? " unresolved" : ""}) ---`,
    );
    if (reviewThreads.length === 0) {
      console.log("No review threads found.");
    } else {
      for (const thread of reviewThreads) {
        if (thread.id) {
          const status = thread.isResolved ? "RESOLVED" : "UNRESOLVED";
          console.log(`\n[Thread ID: ${thread.id}] Status: ${status}`);
        } else {
          console.log(
            "\n[Review Comment (REST Fallback - Not resolvable via GraphQL)]",
          );
        }
        for (const c of thread.comments.nodes || []) {
          const lineNum = c.line || c.originalLine || "?";
          console.log(
            `  Comment #${c.databaseId} by @${c.author?.login || "ghost"} on ${c.path}:${lineNum}:`,
          );
          console.log(`    ${c.body.replace(/\n/g, "\n    ")}`);
        }
      }
    }
  } catch (err) {
    console.error(
      `❌ Error fetching comments for PR #${prNumber} in repository ${owner}/${repo}: ${err.message}`,
    );
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
