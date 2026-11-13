import assert from "node:assert";
import test from "node:test";
import {
  getCurrentBranch,
  getCurrentPRNumber,
  getRepoDetails,
  sys,
} from "../scripts/auth.js";
import { findPR } from "../scripts/find_pr.js";
import { readComments } from "../scripts/read_comments.js";
import { parseReplyArgs, replyComment } from "../scripts/reply-comment.js";
import { parseResolveArgs, resolveThread } from "../scripts/resolve_comment.js";
import {
  parseRespondArgs,
  respondComment,
} from "../scripts/respond_comment.js";
import { parseUpdateArgs, updatePR } from "../scripts/update_pr.js";

test("github-pr skill - getRepoDetails unit tests", async (t) => {
  await t.test(
    "resolves owner and repo from GH_REPO env variable",
    async () => {
      const originalGhRepo = process.env.GH_REPO;
      try {
        process.env.GH_REPO = "custom-org/custom-repo";
        const details = await getRepoDetails();
        assert.deepStrictEqual(details, {
          owner: "custom-org",
          repo: "custom-repo",
        });
      } finally {
        if (originalGhRepo !== undefined) {
          process.env.GH_REPO = originalGhRepo;
        } else {
          delete process.env.GH_REPO;
        }
      }
    },
  );

  await t.test(
    "resolves owner and repo from GITHUB_REPOSITORY env variable",
    async () => {
      const originalGhRepo = process.env.GH_REPO;
      const originalGithubRepo = process.env.GITHUB_REPOSITORY;
      try {
        delete process.env.GH_REPO;
        process.env.GITHUB_REPOSITORY = "upstream-org/upstream-repo";
        const details = await getRepoDetails();
        assert.deepStrictEqual(details, {
          owner: "upstream-org",
          repo: "upstream-repo",
        });
      } finally {
        if (originalGhRepo !== undefined) {
          process.env.GH_REPO = originalGhRepo;
        } else {
          delete process.env.GH_REPO;
        }
        if (originalGithubRepo !== undefined) {
          process.env.GITHUB_REPOSITORY = originalGithubRepo;
        } else {
          delete process.env.GITHUB_REPOSITORY;
        }
      }
    },
  );
});

test("github-pr skill - getCurrentPRNumber unit tests", async (t) => {
  await t.test(
    "resolves PR number by matching branch in open pulls",
    async () => {
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        if (cmd === "git" && args.includes("--show-current")) {
          return { stdout: "feature/skill-update\n" };
        }
        return { stdout: "" };
      });

      const mockOctokit = {
        rest: {
          pulls: {
            list: async () => ({
              data: [
                { number: 101, head: { ref: "main" } },
                { number: 448, head: { ref: "feature/skill-update" } },
              ],
            }),
          },
        },
      };

      const originalGhRepo = process.env.GH_REPO;
      try {
        process.env.GH_REPO = "testowner/testrepo";
        const num = await getCurrentPRNumber(mockOctokit);
        assert.strictEqual(num, 448);
      } finally {
        if (originalGhRepo !== undefined) {
          process.env.GH_REPO = originalGhRepo;
        } else {
          delete process.env.GH_REPO;
        }
      }
    },
  );

  await t.test("uses provided repoDetails parameter directly", async () => {
    t.mock.method(sys, "execFileAsync", async (cmd, args) => {
      if (cmd === "git" && args.includes("--show-current")) {
        return { stdout: "feature/skill-update\n" };
      }
      return { stdout: "" };
    });

    let capturedParams = null;
    const mockOctokit = {
      rest: {
        pulls: {
          list: async (params) => {
            capturedParams = params;
            return {
              data: [{ number: 777, head: { ref: "feature/skill-update" } }],
            };
          },
        },
      },
    };

    const num = await getCurrentPRNumber(mockOctokit, {
      owner: "explicit-owner",
      repo: "explicit-repo",
    });
    assert.strictEqual(num, 777);
    assert.strictEqual(capturedParams.owner, "explicit-owner");
    assert.strictEqual(capturedParams.repo, "explicit-repo");
  });

  await t.test("getCurrentBranch retrieves branch name from git", async () => {
    t.mock.method(sys, "execFileAsync", async (cmd, args) => {
      if (cmd === "git" && args.includes("--show-current")) {
        return { stdout: "feature/ci-skills\n" };
      }
      return { stdout: "" };
    });

    const branch = await getCurrentBranch();
    assert.strictEqual(branch, "feature/ci-skills");
  });

  await t.test("getCurrentPRNumber uses paginate when available", async () => {
    t.mock.method(sys, "execFileAsync", async (cmd, args) => {
      if (cmd === "git" && args.includes("--show-current")) {
        return { stdout: "feature/paginated-branch\n" };
      }
      return { stdout: "" };
    });

    let paginateCalled = false;
    const mockOctokit = {
      paginate: async (endpoint, params) => {
        paginateCalled = true;
        assert.strictEqual(params.state, "open");
        return [
          { number: 101, head: { ref: "other-branch" } },
          { number: 102, head: { ref: "feature/paginated-branch" } },
        ];
      },
      rest: {
        pulls: {
          list: async () => {
            throw new Error("Should have called octokit.paginate");
          },
        },
      },
    };

    const num = await getCurrentPRNumber(mockOctokit, {
      owner: "org",
      repo: "repo",
    });
    assert.strictEqual(paginateCalled, true);
    assert.strictEqual(num, 102);
  });

  await t.test(
    "constrains match to head owner when fork PR shares branch name",
    async () => {
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        if (cmd === "git" && args.includes("--show-current")) {
          return { stdout: "feature/shared-branch\n" };
        }
        if (
          cmd === "git" &&
          args.includes("get-url") &&
          args.includes("origin")
        ) {
          return {
            stdout: "git@github.com:myfork/terraform-provider-file.git\n",
          };
        }
        return { stdout: "" };
      });

      const mockOctokit = {
        rest: {
          pulls: {
            list: async () => ({
              data: [
                {
                  number: 101,
                  head: {
                    ref: "feature/shared-branch",
                    repo: { owner: { login: "otherfork" } },
                  },
                },
                {
                  number: 202,
                  head: {
                    ref: "feature/shared-branch",
                    repo: { owner: { login: "myfork" } },
                  },
                },
              ],
            }),
          },
        },
      };

      const num = await getCurrentPRNumber(mockOctokit, {
        owner: "upstream-org",
        repo: "upstream-repo",
      });
      assert.strictEqual(num, 202);
    },
  );

  await t.test(
    "paginates pulls.list beyond page 1 when octokit.paginate is unavailable",
    async () => {
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        if (cmd === "git" && args.includes("--show-current")) {
          return { stdout: "feature/page2-branch\n" };
        }
        return { stdout: "" };
      });

      const queriedPages = [];
      const mockOctokit = {
        rest: {
          pulls: {
            list: async (params) => {
              queriedPages.push(params.page);
              if (params.page === 1) {
                const fullPage = Array.from({ length: 100 }, (_, i) => ({
                  number: i + 1,
                  head: { ref: `other-branch-${i}` },
                }));
                return { data: fullPage };
              }
              return {
                data: [{ number: 555, head: { ref: "feature/page2-branch" } }],
              };
            },
          },
        },
      };

      const num = await getCurrentPRNumber(mockOctokit, {
        owner: "org",
        repo: "repo",
      });
      assert.strictEqual(num, 555);
      assert.deepStrictEqual(queriedPages, [1, 2]);
    },
  );
});

test("github-pr skill - replyComment unit tests", async (t) => {
  await t.test("rejects empty commentId or body", async () => {
    await assert.rejects(async () => {
      await replyComment({
        octokit: {},
        owner: "rancher",
        repo: "terraform-provider-file",
        prNumber: 123,
        commentId: null,
        body: "hello",
      });
    }, /Missing comment ID/);

    await assert.rejects(async () => {
      await replyComment({
        octokit: {},
        owner: "rancher",
        repo: "terraform-provider-file",
        prNumber: 123,
        commentId: 456,
        body: "",
      });
    }, /cannot be empty/);
  });

  await t.test("calls createReplyForReviewComment successfully", async () => {
    let captured = null;
    const mockOctokit = {
      rest: {
        pulls: {
          createReplyForReviewComment: async (params) => {
            captured = params;
            return { data: { id: 789, body: params.body } };
          },
        },
      },
    };

    const result = await replyComment({
      octokit: mockOctokit,
      owner: "testowner",
      repo: "testrepo",
      prNumber: 42,
      commentId: 100,
      body: "Thanks for the review!",
    });

    assert.strictEqual(result.type, "review_reply");
    assert.deepStrictEqual(captured, {
      owner: "testowner",
      repo: "testrepo",
      pull_number: 42,
      comment_id: 100,
      body: "Thanks for the review!",
    });
  });

  await t.test(
    "falls back to issue comment when review comment gives 404",
    async () => {
      let capturedIssueComment = null;
      let getCommentCalled = false;
      const mockOctokit = {
        rest: {
          pulls: {
            createReplyForReviewComment: async () => {
              const err = new Error("Not found");
              err.status = 404;
              throw err;
            },
          },
          issues: {
            getComment: async () => {
              getCommentCalled = true;
              return { data: { id: 100 } };
            },
            createComment: async (params) => {
              capturedIssueComment = params;
              return { data: { id: 999 } };
            },
          },
        },
      };

      const result = await replyComment({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 42,
        commentId: 100,
        body: "General conversation reply",
      });

      assert.strictEqual(result.type, "issue_reply");
      assert.strictEqual(getCommentCalled, true);
      assert.match(capturedIssueComment.body, /> In reply to comment #100/);
      assert.match(capturedIssueComment.body, /General conversation reply/);
    },
  );

  await t.test(
    "re-throws 422 validation errors without falling back to issue comment",
    async () => {
      let issueCommentCreated = false;
      const mockOctokit = {
        rest: {
          pulls: {
            createReplyForReviewComment: async () => {
              const err = new Error("Validation Failed");
              err.status = 422;
              throw err;
            },
          },
          issues: {
            getComment: async () => ({ data: { id: 100 } }),
            createComment: async () => {
              issueCommentCreated = true;
            },
          },
        },
      };

      await assert.rejects(
        async () => {
          await replyComment({
            octokit: mockOctokit,
            owner: "testowner",
            repo: "testrepo",
            prNumber: 42,
            commentId: 100,
            body: "invalid body",
          });
        },
        (err) => err.status === 422,
      );

      assert.strictEqual(issueCommentCreated, false);
    },
  );

  await t.test(
    "resolves targetPrNumber from issueComment.issue_url when prNumber omitted",
    async () => {
      let capturedIssueComment = null;
      const mockOctokit = {
        rest: {
          pulls: {
            getReviewComment: async () => {
              const err = new Error("Not Found");
              err.status = 404;
              throw err;
            },
            createReplyForReviewComment: async () => {
              const err = new Error("Not Found");
              err.status = 404;
              throw err;
            },
          },
          issues: {
            getComment: async () => ({
              data: {
                id: 200,
                issue_url:
                  "https://api.github.com/repos/testowner/testrepo/issues/777",
              },
            }),
            createComment: async (params) => {
              capturedIssueComment = params;
              return { data: { id: 999 } };
            },
          },
        },
      };

      const result = await replyComment({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        commentId: 200,
        body: "Replying without known PR number",
      });

      assert.strictEqual(result.type, "issue_reply");
      assert.strictEqual(capturedIssueComment.issue_number, 777);
      assert.match(
        capturedIssueComment.body,
        /Replying without known PR number/,
      );
    },
  );

  await t.test(
    "replies to GraphQL thread ID using addPullRequestReviewThreadReply mutation",
    async () => {
      let capturedQuery = null;
      let capturedVars = null;
      const mockOctokit = {
        graphql: async (query, vars) => {
          capturedQuery = query;
          capturedVars = vars;
          return {
            addPullRequestReviewThreadReply: {
              comment: { id: "PRRC_1", databaseId: 999, body: vars.body },
            },
          };
        },
      };

      const result = await replyComment({
        octokit: mockOctokit,
        commentId: "PRRT_kwDO12345",
        body: "Addressed thread",
      });

      assert.strictEqual(result.type, "thread_reply");
      assert.strictEqual(result.data.databaseId, 999);
      assert.match(capturedQuery, /addPullRequestReviewThreadReply/);
      assert.strictEqual(capturedVars.threadId, "PRRT_kwDO12345");
      assert.strictEqual(capturedVars.body, "Addressed thread");
    },
  );

  await t.test(
    "automatically resolves prNumber from getReviewComment when prNumber is omitted",
    async () => {
      let capturedReply = null;
      const mockOctokit = {
        rest: {
          pulls: {
            getReviewComment: async ({ comment_id }) => {
              assert.strictEqual(comment_id, 100);
              return {
                data: {
                  id: 100,
                  pull_request_url:
                    "https://api.github.com/repos/testowner/testrepo/pulls/88",
                },
              };
            },
            createReplyForReviewComment: async (params) => {
              capturedReply = params;
              return { data: { id: 790, body: params.body } };
            },
          },
        },
      };

      const result = await replyComment({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        commentId: 100,
        body: "Discovered PR reply",
      });

      assert.strictEqual(result.type, "review_reply");
      assert.strictEqual(capturedReply.pull_number, 88);
      assert.strictEqual(capturedReply.comment_id, 100);
    },
  );
});

test("github-pr skill - updatePR unit tests", async (t) => {
  await t.test("throws when neither title nor body is provided", async () => {
    await assert.rejects(async () => {
      await updatePR({
        octokit: {},
        owner: "owner",
        repo: "repo",
        prNumber: 1,
      });
    }, /No title or body updates provided/);
  });

  await t.test(
    "calls octokit.rest.pulls.update with supplied fields",
    async () => {
      let captured = null;
      const mockOctokit = {
        rest: {
          pulls: {
            update: async (params) => {
              captured = params;
              return {
                data: { number: params.pull_number, title: params.title },
              };
            },
          },
        },
      };

      const updated = await updatePR({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 99,
        title: "New PR Title",
      });

      assert.strictEqual(updated.title, "New PR Title");
      assert.strictEqual(captured.title, "New PR Title");
      assert.strictEqual(captured.pull_number, 99);
    },
  );

  await t.test(
    "calls octokit.rest.pulls.update with supplied state",
    async () => {
      let captured = null;
      const mockOctokit = {
        rest: {
          pulls: {
            update: async (params) => {
              captured = params;
              return {
                data: { number: params.pull_number, state: params.state },
              };
            },
          },
        },
      };

      const updated = await updatePR({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 99,
        state: "closed",
      });

      assert.strictEqual(updated.state, "closed");
      assert.strictEqual(captured.state, "closed");
      assert.strictEqual(captured.pull_number, 99);
    },
  );
});

test("github-pr skill - resolveThread unit tests", async (t) => {
  await t.test("calls octokit.graphql with threadId", async () => {
    let capturedQuery = null;
    let capturedVars = null;

    const mockOctokit = {
      graphql: async (query, vars) => {
        capturedQuery = query;
        capturedVars = vars;
        return {
          resolveReviewThread: {
            thread: { id: vars.threadId, isResolved: true },
          },
        };
      },
    };

    const thread = await resolveThread({
      octokit: mockOctokit,
      threadId: "PRRT_kwDOAbc123",
    });

    assert.strictEqual(thread.id, "PRRT_kwDOAbc123");
    assert.strictEqual(thread.isResolved, true);
    assert.match(capturedQuery, /resolveReviewThread/);
    assert.strictEqual(capturedVars.threadId, "PRRT_kwDOAbc123");
  });

  await t.test(
    "resolves thread given numeric comment ID by matching review threads",
    async () => {
      let capturedThreadId = null;
      const mockOctokit = {
        graphql: async (query, vars) => {
          if (query.includes("resolveReviewThread")) {
            capturedThreadId = vars.threadId;
            return {
              resolveReviewThread: {
                thread: { id: vars.threadId, isResolved: true },
              },
            };
          }
          return {
            repository: {
              pullRequest: {
                reviewThreads: {
                  pageInfo: { hasNextPage: false },
                  nodes: [
                    {
                      id: "PRRT_match",
                      comments: {
                        pageInfo: { hasNextPage: false },
                        nodes: [{ databaseId: 1001 }],
                      },
                    },
                  ],
                },
              },
            },
          };
        },
        rest: { issues: { listComments: () => {} } },
        paginate: async () => [],
      };

      const thread = await resolveThread({
        octokit: mockOctokit,
        threadId: 1001,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 42,
      });

      assert.strictEqual(thread.id, "PRRT_match");
      assert.strictEqual(capturedThreadId, "PRRT_match");
    },
  );

  await t.test(
    "throws when numeric comment ID is not found in review threads",
    async () => {
      const mockOctokit = {
        graphql: async () => ({
          repository: {
            pullRequest: {
              reviewThreads: {
                pageInfo: { hasNextPage: false },
                nodes: [],
              },
            },
          },
        }),
        rest: { issues: { listComments: () => {} } },
        paginate: async () => [],
      };

      await assert.rejects(async () => {
        await resolveThread({
          octokit: mockOctokit,
          threadId: 9999,
          owner: "testowner",
          repo: "testrepo",
          prNumber: 42,
        });
      }, /No review thread found containing comment #9999/);
    },
  );

  await t.test(
    "resolves thread by numeric comment ID automatically discovering prNumber via getReviewComment",
    async () => {
      let queriedPr = null;
      const mockOctokit = {
        graphql: async (query, vars) => {
          if (query.includes("resolveReviewThread")) {
            return {
              resolveReviewThread: {
                thread: { id: vars.threadId, isResolved: true },
              },
            };
          }
          queriedPr = vars.prNumber;
          return {
            repository: {
              pullRequest: {
                reviewThreads: {
                  pageInfo: { hasNextPage: false },
                  nodes: [
                    {
                      id: "PRRT_discovered",
                      comments: {
                        pageInfo: { hasNextPage: false },
                        nodes: [{ databaseId: 2002 }],
                      },
                    },
                  ],
                },
              },
            },
          };
        },
        rest: {
          pulls: {
            getReviewComment: async ({ comment_id }) => {
              assert.strictEqual(comment_id, 2002);
              return {
                data: {
                  id: 2002,
                  pull_request_url:
                    "https://api.github.com/repos/testowner/testrepo/pulls/77",
                },
              };
            },
          },
        },
      };

      const thread = await resolveThread({
        octokit: mockOctokit,
        threadId: 2002,
        owner: "testowner",
        repo: "testrepo",
      });

      assert.strictEqual(thread.id, "PRRT_discovered");
      assert.strictEqual(queriedPr, 77);
    },
  );

  await t.test(
    "calls unresolveReviewThread mutation when unresolve is true",
    async () => {
      let capturedQuery = null;
      let capturedVars = null;

      const mockOctokit = {
        graphql: async (query, vars) => {
          capturedQuery = query;
          capturedVars = vars;
          return {
            unresolveReviewThread: {
              thread: { id: vars.threadId, isResolved: false },
            },
          };
        },
      };

      const thread = await resolveThread({
        octokit: mockOctokit,
        threadId: "PRRT_kwDOAbc123",
        unresolve: true,
      });

      assert.strictEqual(thread.id, "PRRT_kwDOAbc123");
      assert.strictEqual(thread.isResolved, false);
      assert.match(capturedQuery, /unresolveReviewThread/);
      assert.strictEqual(capturedVars?.threadId, "PRRT_kwDOAbc123");
    },
  );
});

test("github-pr skill - findPR unit tests", async (t) => {
  await t.test(
    'resolves current branch PR when target is "current" or "HEAD"',
    async () => {
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        if (cmd === "git" && args.includes("--show-current")) {
          return { stdout: "feature/skill-update\n" };
        }
        return { stdout: "" };
      });

      const mockOctokit = {
        rest: {
          pulls: {
            list: async () => ({
              data: [{ number: 555, head: { ref: "feature/skill-update" } }],
            }),
            get: async ({ pull_number }) => ({
              data: { number: pull_number, title: "Current PR" },
            }),
          },
        },
      };

      const originalGhRepo = process.env.GH_REPO;
      try {
        process.env.GH_REPO = "testowner/testrepo";

        const prCurrent = await findPR({
          octokit: mockOctokit,
          owner: "testowner",
          repo: "testrepo",
          target: "current",
        });
        assert.strictEqual(prCurrent.number, 555);

        const prHead = await findPR({
          octokit: mockOctokit,
          owner: "testowner",
          repo: "testrepo",
          target: "HEAD",
        });
        assert.strictEqual(prHead.number, 555);
      } finally {
        if (originalGhRepo !== undefined) {
          process.env.GH_REPO = originalGhRepo;
        } else {
          delete process.env.GH_REPO;
        }
      }
    },
  );
  await t.test("treats decimal digit target as PR number", async () => {
    let captured = null;
    const mockOctokit = {
      rest: {
        pulls: {
          get: async (params) => {
            captured = params;
            return { data: { number: params.pull_number, title: "Test PR" } };
          },
        },
      },
    };

    const pr = await findPR({
      octokit: mockOctokit,
      owner: "testowner",
      repo: "testrepo",
      target: "123",
    });

    assert.strictEqual(pr.number, 123);
    assert.strictEqual(captured.pull_number, 123);
  });

  await t.test(
    "treats numeric-prefix branch like 123-fix as branch name, not PR number",
    async () => {
      let listCaptured = null;
      let getCaptured = null;
      const mockOctokit = {
        rest: {
          pulls: {
            list: async (params) => {
              listCaptured = params;
              return { data: [{ number: 999 }] };
            },
            get: async (params) => {
              getCaptured = params;
              return {
                data: { number: params.pull_number, title: "Branch PR" },
              };
            },
          },
        },
      };

      const pr = await findPR({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        target: "123-fix",
      });

      assert.strictEqual(pr.number, 999);
      assert.strictEqual(listCaptured.head, "testowner:123-fix");
      assert.strictEqual(getCaptured.pull_number, 999);
    },
  );

  await t.test(
    "treats # prefixed number like #448 as numeric PR number",
    async () => {
      let captured = null;
      const mockOctokit = {
        rest: {
          pulls: {
            get: async (params) => {
              captured = params;
              return { data: { number: params.pull_number, title: "Hash PR" } };
            },
          },
        },
      };

      const pr = await findPR({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        target: "#448",
      });

      assert.strictEqual(pr.number, 448);
      assert.strictEqual(captured.pull_number, 448);
    },
  );

  await t.test(
    "preserves explicit user:branch format in head filter",
    async () => {
      let listCaptured = null;
      const mockOctokit = {
        rest: {
          pulls: {
            list: async (params) => {
              listCaptured = params;
              return { data: [{ number: 777 }] };
            },
            get: async (params) => ({ data: { number: params.pull_number } }),
          },
        },
      };

      const pr = await findPR({
        octokit: mockOctokit,
        owner: "upstream-org",
        repo: "provider-repo",
        target: "forkuser:feature-branch",
      });

      assert.strictEqual(pr.number, 777);
      assert.strictEqual(listCaptured.head, "forkuser:feature-branch");
    },
  );

  await t.test(
    "falls back to searching open PRs by head.ref when owner:branch query returns empty",
    async () => {
      let callCount = 0;
      const mockOctokit = {
        rest: {
          pulls: {
            list: async (params) => {
              callCount++;
              if (params.head === "upstream-org:fork-feature") {
                return { data: [] };
              }
              return {
                data: [
                  { number: 888, head: { ref: "other" } },
                  { number: 999, head: { ref: "fork-feature" } },
                ],
              };
            },
            get: async (params) => ({
              data: { number: params.pull_number, title: "Found Fork PR" },
            }),
          },
        },
      };

      const pr = await findPR({
        octokit: mockOctokit,
        owner: "upstream-org",
        repo: "provider-repo",
        target: "fork-feature",
      });

      assert.strictEqual(pr.number, 999);
      assert.strictEqual(callCount, 2);
    },
  );

  await t.test(
    "prioritizes open PRs over closed PRs in head.ref fallback search",
    async () => {
      const mockOctokit = {
        rest: {
          pulls: {
            list: async (params) => {
              if (params.head) {
                return { data: [] };
              }
              return {
                data: [
                  {
                    number: 100,
                    state: "closed",
                    head: { ref: "fork-feature" },
                  },
                  { number: 200, state: "open", head: { ref: "fork-feature" } },
                ],
              };
            },
            get: async (params) => ({
              data: { number: params.pull_number, title: "Found Open Fork PR" },
            }),
          },
        },
      };

      const pr = await findPR({
        octokit: mockOctokit,
        owner: "upstream-org",
        repo: "provider-repo",
        target: "fork-feature",
      });

      assert.strictEqual(pr.number, 200);
    },
  );

  await t.test("resolves PR from GitHub PR URL", async () => {
    let captured = null;
    const mockOctokit = {
      rest: {
        pulls: {
          get: async (params) => {
            captured = params;
            return { data: { number: params.pull_number, title: "From URL" } };
          },
        },
      },
    };

    const pr = await findPR({
      octokit: mockOctokit,
      owner: "testowner",
      repo: "testrepo",
      target: "https://github.com/testowner/testrepo/pull/555",
    });

    assert.strictEqual(pr.number, 555);
    assert.strictEqual(captured.pull_number, 555);
  });

  await t.test(
    "prioritizes open PRs over closed PRs in direct head query",
    async () => {
      const queries = [];
      const mockOctokit = {
        rest: {
          pulls: {
            list: async (params) => {
              queries.push(params);
              return {
                data: [
                  { number: 222, state: "closed", head: { ref: "my-feature" } },
                  { number: 333, state: "open", head: { ref: "my-feature" } },
                ],
              };
            },
            get: async (params) => ({
              data: { number: params.pull_number, title: "Open PR" },
            }),
          },
        },
      };

      const pr = await findPR({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        target: "my-feature",
      });

      assert.strictEqual(pr.number, 333);
      assert.strictEqual(queries[0].head, "testowner:my-feature");
    },
  );
});

test("github-pr skill - readComments unit tests", async (t) => {
  await t.test(
    "paginates GraphQL review threads and thread comments",
    async () => {
      let threadQueryCount = 0;
      const mockOctokit = {
        rest: {
          issues: {
            listComments: () => {},
          },
        },
        paginate: async () => [{ id: 1, body: "Issue comment 1" }],
        graphql: async (query, vars) => {
          if (vars.threadId) {
            // Comment pagination
            return {
              node: {
                comments: {
                  pageInfo: { hasNextPage: false, endCursor: null },
                  nodes: [
                    {
                      databaseId: 1002,
                      body: "Reply comment",
                      author: { login: "dev" },
                    },
                  ],
                },
              },
            };
          }
          threadQueryCount++;
          if (threadQueryCount === 1) {
            return {
              repository: {
                pullRequest: {
                  reviewThreads: {
                    pageInfo: { hasNextPage: true, endCursor: "c1" },
                    nodes: [
                      {
                        id: "PRRT_1",
                        isResolved: false,
                        comments: {
                          pageInfo: { hasNextPage: true, endCursor: "cc1" },
                          nodes: [
                            {
                              databaseId: 1001,
                              body: "First comment",
                              author: { login: "reviewer" },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              },
            };
          }
          return {
            repository: {
              pullRequest: {
                reviewThreads: {
                  pageInfo: { hasNextPage: false, endCursor: null },
                  nodes: [
                    {
                      id: "PRRT_2",
                      isResolved: true,
                      comments: {
                        pageInfo: { hasNextPage: false, endCursor: null },
                        nodes: [
                          {
                            databaseId: 1003,
                            body: "Thread 2 comment",
                            author: { login: "reviewer" },
                          },
                        ],
                      },
                    },
                  ],
                },
              },
            },
          };
        },
      };

      const res = await readComments({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 42,
      });

      assert.strictEqual(res.conversationComments.length, 1);
      assert.strictEqual(res.reviewThreads.length, 2);
      assert.strictEqual(res.reviewThreads[0].id, "PRRT_1");
      assert.strictEqual(res.reviewThreads[0].comments.nodes.length, 2);
      assert.strictEqual(res.reviewThreads[1].id, "PRRT_2");
      assert.strictEqual(res.reviewThreads[1].isResolved, true);
    },
  );

  await t.test(
    "falls back to REST review comments when GraphQL fails and sets id to null",
    async () => {
      const mockOctokit = {
        paginate: async (method) => {
          if (method.name === "listComments") {
            return [{ id: 50, body: "General comment" }];
          }
          return [
            {
              id: 200,
              node_id: "PRRC_abc123",
              user: { login: "reviewer" },
              body: "REST review comment",
              path: "main.go",
              line: 10,
              original_line: 8,
            },
          ];
        },
        graphql: async () => {
          throw new Error("GraphQL forbidden");
        },
        rest: {
          issues: { listComments: { name: "listComments" } },
          pulls: { listReviewComments: { name: "listReviewComments" } },
        },
      };

      const res = await readComments({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 42,
      });

      assert.strictEqual(res.reviewThreads.length, 1);
      assert.strictEqual(res.reviewThreads[0].id, null);
      assert.strictEqual(res.reviewThreads[0].isResolvable, false);
      assert.strictEqual(
        res.reviewThreads[0].comments.nodes[0].databaseId,
        200,
      );
      assert.strictEqual(
        res.reviewThreads[0].comments.nodes[0].originalLine,
        8,
      );
    },
  );

  await t.test(
    "preserves originalLine when line is null in GraphQL query",
    async () => {
      const mockOctokit = {
        graphql: async () => ({
          repository: {
            pullRequest: {
              reviewThreads: {
                pageInfo: { hasNextPage: false, endCursor: null },
                nodes: [
                  {
                    id: "PRRT_orig",
                    isResolved: false,
                    comments: {
                      pageInfo: { hasNextPage: false, endCursor: null },
                      nodes: [
                        {
                          databaseId: 555,
                          body: "Comment on original line",
                          author: { login: "reviewer" },
                          path: "main.go",
                          line: null,
                          originalLine: 42,
                        },
                      ],
                    },
                  },
                ],
              },
            },
          },
        }),
      };

      const res = await readComments({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 99,
      });

      assert.strictEqual(res.reviewThreads.length, 1);
      const comment = res.reviewThreads[0].comments.nodes[0];
      assert.strictEqual(comment.line, null);
      assert.strictEqual(comment.originalLine, 42);
    },
  );

  await t.test(
    'coerces string prNumber like "#42" to integer for GraphQL query',
    async () => {
      let capturedPr = null;
      const mockOctokit = {
        graphql: async (query, vars) => {
          capturedPr = vars.prNumber;
          return {
            repository: {
              pullRequest: {
                reviewThreads: {
                  pageInfo: { hasNextPage: false },
                  nodes: [],
                },
              },
            },
          };
        },
      };

      await readComments({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: "#42",
      });

      assert.strictEqual(typeof capturedPr, "number");
      assert.strictEqual(capturedPr, 42);
    },
  );

  await t.test("resolves PR URL in readComments prNumber", async () => {
    let capturedPr = null;
    const mockOctokit = {
      graphql: async (query, vars) => {
        capturedPr = vars.prNumber;
        return {
          repository: {
            pullRequest: {
              reviewThreads: {
                pageInfo: { hasNextPage: false },
                nodes: [],
              },
            },
          },
        };
      },
    };

    await readComments({
      octokit: mockOctokit,
      owner: "testowner",
      repo: "testrepo",
      prNumber: "https://github.com/testowner/testrepo/pull/999",
    });

    assert.strictEqual(capturedPr, 999);
  });

  await t.test(
    "filters review threads when unresolvedOnly is true",
    async () => {
      const mockOctokit = {
        graphql: async () => ({
          repository: {
            pullRequest: {
              reviewThreads: {
                pageInfo: { hasNextPage: false },
                nodes: [
                  {
                    id: "PRRT_unresolved",
                    isResolved: false,
                    comments: { nodes: [] },
                  },
                  {
                    id: "PRRT_resolved",
                    isResolved: true,
                    comments: { nodes: [] },
                  },
                ],
              },
            },
          },
        }),
      };

      const res = await readComments({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 42,
        unresolvedOnly: true,
      });

      assert.strictEqual(res.reviewThreads.length, 1);
      assert.strictEqual(res.reviewThreads[0].id, "PRRT_unresolved");
      assert.strictEqual(res.reviewThreads[0].isResolved, false);
    },
  );

  await t.test(
    'passes coerced numeric prNumber to REST fallback when prNumber has "#" prefix',
    async () => {
      let capturedParams = null;
      const mockOctokit = {
        graphql: async () => {
          throw new Error("GraphQL error");
        },
        rest: {
          pulls: {
            listReviewComments: async (params) => {
              capturedParams = params;
              return { data: [] };
            },
          },
        },
      };

      await readComments({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: "#77",
      });

      assert.strictEqual(typeof capturedParams.pull_number, "number");
      assert.strictEqual(capturedParams.pull_number, 77);
    },
  );

  await t.test(
    "clears partial GraphQL results when later page fails and falls back cleanly to REST",
    async () => {
      let callCount = 0;
      const mockOctokit = {
        graphql: async () => {
          callCount++;
          if (callCount === 1) {
            return {
              repository: {
                pullRequest: {
                  reviewThreads: {
                    nodes: [
                      {
                        id: "PRRT_partial",
                        isResolved: false,
                        comments: { nodes: [{ id: "1", body: "Partial" }] },
                      },
                    ],
                    pageInfo: { hasNextPage: true, endCursor: "cursor_1" },
                  },
                },
              },
            };
          }
          throw new Error("GraphQL secondary page failure");
        },
        rest: {
          pulls: {
            listReviewComments: async () => ({
              data: [
                {
                  id: 999,
                  body: "REST comment",
                  user: { login: "reviewer" },
                  path: "file.js",
                  line: 10,
                  created_at: "2026-10-20T00:00:00Z",
                },
              ],
            }),
          },
        },
      };

      const res = await readComments({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 448,
      });

      assert.strictEqual(res.reviewThreads.length, 1);
      assert.strictEqual(res.reviewThreads[0].id, null);
      assert.strictEqual(
        res.reviewThreads[0].comments.nodes[0].databaseId,
        999,
      );
      assert.strictEqual(
        res.reviewThreads[0].comments.nodes[0].body,
        "REST comment",
      );
    },
  );
});

test("github-pr skill - parseReplyArgs unit tests", async (t) => {
  await t.test("parses [comment_id, body]", () => {
    const parsed = parseReplyArgs(["123", "Fixed in commit abc"]);
    assert.deepStrictEqual(parsed, {
      prNumber: undefined,
      commentId: 123,
      body: "Fixed in commit abc",
    });
  });

  await t.test("parses unquoted multi-word body [comment_id, ...words]", () => {
    const parsed = parseReplyArgs(["123", "Looks", "good", "to", "me"]);
    assert.deepStrictEqual(parsed, {
      prNumber: undefined,
      commentId: 123,
      body: "Looks good to me",
    });
  });

  await t.test("parses [pr_number, comment_id, ...words]", () => {
    const parsed = parseReplyArgs(["42", "123", "Looks", "great"]);
    assert.deepStrictEqual(parsed, {
      prNumber: 42,
      commentId: 123,
      body: "Looks great",
    });
  });

  await t.test("parses [comment_id, ...words, pr_number]", () => {
    const parsed = parseReplyArgs(["123", "Looks", "great", "42"]);
    assert.deepStrictEqual(parsed, {
      prNumber: 42,
      commentId: 123,
      body: "Looks great",
    });
  });

  await t.test("parses # prefixed comment_id and pr_number", () => {
    const parsed = parseReplyArgs(["#42", "#123", "Looks", "great"]);
    assert.deepStrictEqual(parsed, {
      prNumber: 42,
      commentId: 123,
      body: "Looks great",
    });
  });

  await t.test("parses GraphQL thread ID in parseReplyArgs", () => {
    const parsed1 = parseReplyArgs(["PRRT_kwDOAbc123", "Fixed in commit abc"]);
    assert.deepStrictEqual(parsed1, {
      prNumber: undefined,
      commentId: "PRRT_kwDOAbc123",
      body: "Fixed in commit abc",
    });

    const parsed2 = parseReplyArgs([
      "42",
      "PRRT_kwDOAbc123",
      "Fixed in commit abc",
    ]);
    assert.deepStrictEqual(parsed2, {
      prNumber: 42,
      commentId: "PRRT_kwDOAbc123",
      body: "Fixed in commit abc",
    });
  });

  await t.test('parses "current" keyword in parseReplyArgs', () => {
    const parsed1 = parseReplyArgs(["current", "123", "Looks great"]);
    assert.deepStrictEqual(parsed1, {
      prNumber: undefined,
      commentId: 123,
      body: "Looks great",
    });

    const parsed2 = parseReplyArgs(["123", "Looks great", "current"]);
    assert.deepStrictEqual(parsed2, {
      prNumber: undefined,
      commentId: 123,
      body: "Looks great",
    });
  });

  await t.test(
    "parses [comment_id, pr_number, ...words] with larger comment ID or hash prefix",
    () => {
      const parsed1 = parseReplyArgs([
        "1958291048",
        "448",
        "Fixed",
        "in",
        "commit",
      ]);
      assert.deepStrictEqual(parsed1, {
        prNumber: 448,
        commentId: 1958291048,
        body: "Fixed in commit",
      });

      const parsed2 = parseReplyArgs(["123", "#448", "Addressed", "note"]);
      assert.deepStrictEqual(parsed2, {
        prNumber: 448,
        commentId: 123,
        body: "Addressed note",
      });

      const parsed3 = parseReplyArgs(["123", "current", "Addressed", "note"]);
      assert.deepStrictEqual(parsed3, {
        prNumber: undefined,
        commentId: 123,
        body: "Addressed note",
      });

      const parsed4 = parseReplyArgs([
        "PRRT_kwDO123",
        "448",
        "Addressed",
        "thread",
      ]);
      assert.deepStrictEqual(parsed4, {
        prNumber: 448,
        commentId: "PRRT_kwDO123",
        body: "Addressed thread",
      });

      const parsed5 = parseReplyArgs([
        "PRRT_kwDO123",
        "current",
        "Addressed",
        "thread",
      ]);
      assert.deepStrictEqual(parsed5, {
        prNumber: undefined,
        commentId: "PRRT_kwDO123",
        body: "Addressed thread",
      });

      const parsed6 = parseReplyArgs(["999", "448", "Fixed", "in", "commit"]);
      assert.deepStrictEqual(parsed6, {
        prNumber: 448,
        commentId: 999,
        body: "Fixed in commit",
      });
    },
  );

  await t.test("parses -m and --body flags in parseReplyArgs", () => {
    const res1 = parseReplyArgs(["123", "-m", "Fixed issue in commit"]);
    assert.deepStrictEqual(res1, {
      prNumber: undefined,
      commentId: 123,
      body: "Fixed issue in commit",
    });

    const res2 = parseReplyArgs([
      "448",
      "999",
      "--body",
      "Detailed explanation",
    ]);
    assert.deepStrictEqual(res2, {
      prNumber: 448,
      commentId: 999,
      body: "Detailed explanation",
    });

    const res3 = parseReplyArgs(["PRRT_kwDO123", "-m", "Thread reply"]);
    assert.deepStrictEqual(res3, {
      prNumber: undefined,
      commentId: "PRRT_kwDO123",
      body: "Thread reply",
    });
  });

  await t.test("parses PR URL in parseReplyArgs", () => {
    const res = parseReplyArgs([
      "https://github.com/owner/repo/pull/448",
      "123",
      "URL reply",
    ]);
    assert.deepStrictEqual(res, {
      prNumber: 448,
      commentId: 123,
      body: "URL reply",
    });
  });

  await t.test("parses branch target in parseReplyArgs", () => {
    const res1 = parseReplyArgs(["feature-branch", "123", "Branch reply"]);
    assert.deepStrictEqual(res1, {
      prTarget: "feature-branch",
      prNumber: undefined,
      commentId: 123,
      body: "Branch reply",
    });

    const res2 = parseReplyArgs([
      "feature-branch",
      "123",
      "-m",
      "Explicit message",
    ]);
    assert.deepStrictEqual(res2, {
      prTarget: "feature-branch",
      prNumber: undefined,
      commentId: 123,
      body: "Explicit message",
    });
  });

  await t.test(
    "parses -p and --pr flags unambiguously in parseReplyArgs",
    () => {
      const res1 = parseReplyArgs(["123", "--pr", "448", "-m", "Reply body"]);
      assert.deepStrictEqual(res1, {
        prNumber: 448,
        commentId: 123,
        body: "Reply body",
      });

      const res2 = parseReplyArgs([
        "50",
        "-p",
        "1000",
        "-m",
        "Small comment on large PR",
      ]);
      assert.deepStrictEqual(res2, {
        prNumber: 1000,
        commentId: 50,
        body: "Small comment on large PR",
      });

      const res3 = parseReplyArgs([
        "PRRT_abc",
        "--pr",
        "current",
        "-m",
        "Thread reply",
      ]);
      assert.deepStrictEqual(res3, {
        prNumber: undefined,
        commentId: "PRRT_abc",
        body: "Thread reply",
      });
    },
  );

  await t.test("returns empty object for fewer than 2 arguments", () => {
    assert.deepStrictEqual(parseReplyArgs([]), {});
    assert.deepStrictEqual(parseReplyArgs(["123"]), {});
  });
});

test("github-pr skill - parseResolveArgs unit tests", async (t) => {
  await t.test("returns empty object when args is empty", () => {
    assert.deepStrictEqual(parseResolveArgs([]), { unresolve: false });
    assert.deepStrictEqual(parseResolveArgs(undefined), { unresolve: false });
  });

  await t.test("parses single thread ID or numeric comment ID", () => {
    const res1 = parseResolveArgs(["PRRT_abc123"]);
    assert.deepStrictEqual(res1, {
      threadId: "PRRT_abc123",
      prNumber: undefined,
      prTarget: "current",
      unresolve: false,
    });

    const res2 = parseResolveArgs(["98765"]);
    assert.deepStrictEqual(res2, {
      threadId: "98765",
      prNumber: undefined,
      prTarget: "current",
      unresolve: false,
    });
  });

  await t.test("parses [thread_id, pr_number]", () => {
    const res = parseResolveArgs(["PRRT_abc123", "42"]);
    assert.deepStrictEqual(res, {
      threadId: "PRRT_abc123",
      prNumber: 42,
      prTarget: "42",
      unresolve: false,
    });
  });

  await t.test("parses [pr_number, thread_id] and [current, thread_id]", () => {
    const res1 = parseResolveArgs(["42", "PRRT_abc123"]);
    assert.deepStrictEqual(res1, {
      threadId: "PRRT_abc123",
      prNumber: 42,
      prTarget: "42",
      unresolve: false,
    });

    const res2 = parseResolveArgs(["current", "PRRT_abc123"]);
    assert.deepStrictEqual(res2, {
      threadId: "PRRT_abc123",
      prNumber: undefined,
      prTarget: "current",
      unresolve: false,
    });

    const res3 = parseResolveArgs(["42", "12345"]);
    assert.deepStrictEqual(res3, {
      threadId: 12345,
      prNumber: 42,
      prTarget: "42",
      unresolve: false,
    });

    const res4 = parseResolveArgs(["12345", "42"]);
    assert.deepStrictEqual(res4, {
      threadId: 12345,
      prNumber: 42,
      prTarget: "42",
      unresolve: false,
    });

    const res5 = parseResolveArgs(["#42", "12345"]);
    assert.deepStrictEqual(res5, {
      threadId: 12345,
      prNumber: 42,
      prTarget: "42",
      unresolve: false,
    });

    const res6 = parseResolveArgs(["12345", "#42"]);
    assert.deepStrictEqual(res6, {
      threadId: 12345,
      prNumber: 42,
      prTarget: "42",
      unresolve: false,
    });
  });

  await t.test("parses --unresolve flag", () => {
    const res1 = parseResolveArgs(["PRRT_abc123", "--unresolve"]);
    assert.deepStrictEqual(res1, {
      threadId: "PRRT_abc123",
      prNumber: undefined,
      prTarget: "current",
      unresolve: true,
    });

    const res2 = parseResolveArgs(["--unresolve", "42", "PRRT_abc123"]);
    assert.deepStrictEqual(res2, {
      threadId: "PRRT_abc123",
      prNumber: 42,
      prTarget: "42",
      unresolve: true,
    });
  });

  await t.test("parses PR URL in parseResolveArgs", () => {
    const res = parseResolveArgs([
      "https://github.com/owner/repo/pull/448",
      "9999",
    ]);
    assert.deepStrictEqual(res, {
      threadId: 9999,
      prNumber: 448,
      prTarget: "448",
      unresolve: false,
    });
  });

  await t.test(
    "parses -p and --pr flags unambiguously in parseResolveArgs",
    () => {
      const res1 = parseResolveArgs(["50", "--pr", "1000"]);
      assert.deepStrictEqual(res1, {
        threadId: "50",
        prNumber: 1000,
        prTarget: "1000",
        unresolve: false,
      });

      const res2 = parseResolveArgs(["PRRT_abc", "-p", "448", "--unresolve"]);
      assert.deepStrictEqual(res2, {
        threadId: "PRRT_abc",
        prNumber: 448,
        prTarget: "448",
        unresolve: true,
      });
    },
  );
});

test("github-pr skill - respondComment unit tests", async (t) => {
  await t.test("throws when body is empty or whitespace", async () => {
    await assert.rejects(async () => {
      await respondComment({
        octokit: {},
        owner: "testowner",
        repo: "testrepo",
        prNumber: 10,
        body: "",
      });
    }, /Comment body cannot be empty/);
  });

  await t.test(
    "creates general issue comment when commentId is general or 0",
    async () => {
      let captured = null;
      const mockOctokit = {
        rest: {
          issues: {
            createComment: async (params) => {
              captured = params;
              return { data: { id: 777 } };
            },
          },
        },
      };

      const resGeneral = await respondComment({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 42,
        commentId: "general",
        body: "General feedback note",
      });
      assert.strictEqual(resGeneral.type, "general_comment");
      assert.strictEqual(captured.issue_number, 42);
      assert.strictEqual(captured.body, "General feedback note");

      const resZero = await respondComment({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 42,
        commentId: 0,
        body: "Zero comment note",
      });
      assert.strictEqual(resZero.type, "general_comment");
      assert.strictEqual(captured.body, "Zero comment note");
    },
  );

  await t.test(
    "delegates to replyComment when commentId is numeric",
    async () => {
      let capturedReply = null;
      const mockOctokit = {
        rest: {
          pulls: {
            createReplyForReviewComment: async (params) => {
              capturedReply = params;
              return { data: { id: 888 } };
            },
          },
        },
      };

      const res = await respondComment({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        prNumber: 42,
        commentId: 301,
        body: "Review thread reply",
      });

      assert.strictEqual(res.type, "review_reply");
      assert.strictEqual(capturedReply.comment_id, 301);
      assert.strictEqual(capturedReply.pull_number, 42);
      assert.strictEqual(capturedReply.body, "Review thread reply");
    },
  );
});

test("github-pr skill - parseRespondArgs unit tests", async (t) => {
  await t.test("returns empty object when args is empty", () => {
    assert.deepStrictEqual(parseRespondArgs([]), {});
    assert.deepStrictEqual(parseRespondArgs(), {});
  });

  await t.test("parses general comment on current PR", () => {
    const res = parseRespondArgs(["general", "Great", "job", "team"]);
    assert.deepStrictEqual(res, {
      prTarget: "current",
      commentId: "general",
      body: "Great job team",
    });
  });

  await t.test("parses [current, general, ...words]", () => {
    const res = parseRespondArgs(["current", "general", "Looks", "clean"]);
    assert.deepStrictEqual(res, {
      prTarget: "current",
      commentId: "general",
      body: "Looks clean",
    });
  });

  await t.test("parses [current, comment_id, ...words]", () => {
    const res = parseRespondArgs([
      "current",
      "999",
      "Resolved",
      "in",
      "commit",
    ]);
    assert.deepStrictEqual(res, {
      prTarget: "current",
      commentId: 999,
      body: "Resolved in commit",
    });
  });

  await t.test("parses [pr_number, comment_id, ...words]", () => {
    const res = parseRespondArgs(["448", "999", "Resolved", "in", "commit"]);
    assert.deepStrictEqual(res, {
      prTarget: "448",
      commentId: 999,
      body: "Resolved in commit",
    });
  });

  await t.test(
    "parses [comment_id, pr_number, ...words] with larger comment ID",
    () => {
      const res = parseRespondArgs(["999", "448", "Resolved", "in", "commit"]);
      assert.deepStrictEqual(res, {
        prTarget: "448",
        commentId: 999,
        body: "Resolved in commit",
      });
    },
  );

  await t.test("parses [pr_number, general, ...words]", () => {
    const res = parseRespondArgs(["448", "general", "PR", "level", "note"]);
    assert.deepStrictEqual(res, {
      prTarget: "448",
      commentId: "general",
      body: "PR level note",
    });
  });

  await t.test("parses [pr_number, ...words] with omitted commentId", () => {
    const res = parseRespondArgs([
      "448",
      "Direct",
      "note",
      "without",
      "general",
    ]);
    assert.deepStrictEqual(res, {
      prTarget: "448",
      commentId: "general",
      body: "Direct note without general",
    });
  });

  await t.test("parses bare words targeting current PR", () => {
    const res = parseRespondArgs(["Direct", "comment", "on", "current", "PR"]);
    assert.deepStrictEqual(res, {
      prTarget: "current",
      commentId: "general",
      body: "Direct comment on current PR",
    });
  });

  await t.test("parses [branch_name, comment_id, ...words]", () => {
    const res = parseRespondArgs(["feat-x", "123", "Fixed", "here"]);
    assert.deepStrictEqual(res, {
      prTarget: "feat-x",
      commentId: 123,
      body: "Fixed here",
    });
  });

  await t.test("parses GraphQL thread ID in parseRespondArgs", () => {
    const res1 = parseRespondArgs(["PRRT_kwDOAbc123", "Fixed", "here"]);
    assert.deepStrictEqual(res1, {
      prTarget: "current",
      commentId: "PRRT_kwDOAbc123",
      body: "Fixed here",
    });

    const res2 = parseRespondArgs(["448", "PRRT_kwDOAbc123", "Fixed", "here"]);
    assert.deepStrictEqual(res2, {
      prTarget: "448",
      commentId: "PRRT_kwDOAbc123",
      body: "Fixed here",
    });
  });

  await t.test(
    "parses # prefixed PR number and comment ID in parseRespondArgs",
    () => {
      const res = parseRespondArgs([
        "#448",
        "#999",
        "Resolved",
        "in",
        "commit",
      ]);
      assert.deepStrictEqual(res, {
        prTarget: "#448",
        commentId: 999,
        body: "Resolved in commit",
      });
    },
  );

  await t.test("parses -m and --body flags in parseRespondArgs", () => {
    const res1 = parseRespondArgs(["general", "-m", "Overall LGTM"]);
    assert.deepStrictEqual(res1, {
      prTarget: "current",
      commentId: "general",
      body: "Overall LGTM",
    });

    const res2 = parseRespondArgs([
      "448",
      "999",
      "--body",
      "Addressed in commit",
    ]);
    assert.deepStrictEqual(res2, {
      prTarget: "448",
      commentId: 999,
      body: "Addressed in commit",
    });

    const res3 = parseRespondArgs(["-m", "Quick comment"]);
    assert.deepStrictEqual(res3, {
      prTarget: "current",
      commentId: "general",
      body: "Quick comment",
    });
  });

  await t.test("parses PR URL in parseRespondArgs", () => {
    const res = parseRespondArgs([
      "https://github.com/owner/repo/pull/448",
      "general",
      "URL comment",
    ]);
    assert.deepStrictEqual(res, {
      prTarget: "448",
      commentId: "general",
      body: "URL comment",
    });
  });
});

test("github-pr skill - parseUpdateArgs unit tests", async (t) => {
  await t.test("returns empty object when args is empty", () => {
    assert.deepStrictEqual(parseUpdateArgs([]), {});
    assert.deepStrictEqual(parseUpdateArgs(), {});
  });

  await t.test("parses single title argument for current PR", () => {
    const res = parseUpdateArgs(["fix: typo in provider"]);
    assert.deepStrictEqual(res, {
      prTarget: "current",
      title: "fix: typo in provider",
      body: undefined,
    });
  });

  await t.test("parses [title, body] targeting current PR", () => {
    const res = parseUpdateArgs([
      "fix: typo in provider",
      "Detailed description body",
    ]);
    assert.deepStrictEqual(res, {
      prTarget: "current",
      title: "fix: typo in provider",
      body: "Detailed description body",
    });
  });

  await t.test("parses [current, title]", () => {
    const res = parseUpdateArgs(["current", "fix: updated title"]);
    assert.deepStrictEqual(res, {
      prTarget: "current",
      title: "fix: updated title",
      body: undefined,
    });
  });

  await t.test("parses [pr_number, title]", () => {
    const res = parseUpdateArgs(["448", "fix: updated title"]);
    assert.deepStrictEqual(res, {
      prTarget: "448",
      title: "fix: updated title",
      body: undefined,
    });
  });

  await t.test("parses # prefixed [pr_number, title]", () => {
    const res = parseUpdateArgs(["#448", "fix: updated title"]);
    assert.deepStrictEqual(res, {
      prTarget: "#448",
      title: "fix: updated title",
      body: undefined,
    });
  });

  await t.test("parses [pr_number, title, body]", () => {
    const res = parseUpdateArgs(["448", "fix: title", "Line 1", "Line 2"]);
    assert.deepStrictEqual(res, {
      prTarget: "448",
      title: "fix: title",
      body: "Line 1 Line 2",
    });
  });

  await t.test("parses [branch, title, body]", () => {
    const res = parseUpdateArgs([
      "my-branch",
      "fix: title",
      "Description body",
    ]);
    assert.deepStrictEqual(res, {
      prTarget: "my-branch",
      title: "fix: title",
      body: "Description body",
    });
  });

  await t.test("parses --title and --body flags in parseUpdateArgs", () => {
    const res1 = parseUpdateArgs([
      "--title",
      "New title",
      "--body",
      "New description",
    ]);
    assert.deepStrictEqual(res1, {
      prTarget: "current",
      title: "New title",
      body: "New description",
    });

    const res2 = parseUpdateArgs(["448", "-t", "Updated title"]);
    assert.deepStrictEqual(res2, {
      prTarget: "448",
      title: "Updated title",
      body: undefined,
    });

    const res3 = parseUpdateArgs(["-b=Updated body only", "current"]);
    assert.deepStrictEqual(res3, {
      prTarget: "current",
      title: undefined,
      body: "Updated body only",
    });
  });

  await t.test(
    "parses --state / -s flags and PR URLs in parseUpdateArgs",
    () => {
      const res1 = parseUpdateArgs(["--state", "closed"]);
      assert.deepStrictEqual(res1, {
        prTarget: "current",
        title: undefined,
        body: undefined,
        state: "closed",
      });

      const res2 = parseUpdateArgs([
        "https://github.com/owner/repo/pull/448",
        "-s=open",
        "-t",
        "Reopened PR",
      ]);
      assert.deepStrictEqual(res2, {
        prTarget: "https://github.com/owner/repo/pull/448",
        title: "Reopened PR",
        body: undefined,
        state: "open",
      });
    },
  );
});
