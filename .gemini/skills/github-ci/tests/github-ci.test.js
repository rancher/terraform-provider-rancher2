import assert from "node:assert";
import { Buffer } from "node:buffer";
import path from "node:path";
import test from "node:test";
import {
  getCurrentBranch,
  getGhCliToken,
  getRepoDetails,
  sys,
} from "../scripts/auth.js";
import { downloadLogs, parseDownloadArgs } from "../scripts/download_logs.js";
import { findFailedRuns, parseRunArgs } from "../scripts/find_failed_runs.js";
import { getFailedLog, parseFailedLogArgs } from "../scripts/get_failed_log.js";
import { listRuns } from "../scripts/list_runs.js";
import { parseSearchArgs, searchLogs } from "../scripts/search_logs.js";

test("github-ci skill - getRepoDetails unit tests", async (t) => {
  await t.test(
    "resolves owner and repo from GH_REPO env variable",
    async () => {
      const originalGhRepo = process.env.GH_REPO;
      try {
        process.env.GH_REPO = "ci-org/ci-repo";
        const details = await getRepoDetails();
        assert.deepStrictEqual(details, { owner: "ci-org", repo: "ci-repo" });
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
        process.env.GITHUB_REPOSITORY = "ci-upstream/ci-repo";
        const details = await getRepoDetails();
        assert.deepStrictEqual(details, {
          owner: "ci-upstream",
          repo: "ci-repo",
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

  await t.test("resolves owner and repo from git remote fallback", async () => {
    const originalGhRepo = process.env.GH_REPO;
    const originalGithubRepo = process.env.GITHUB_REPOSITORY;
    try {
      delete process.env.GH_REPO;
      delete process.env.GITHUB_REPOSITORY;

      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        if (cmd === "git" && args.includes("upstream")) {
          return { stdout: "git@github.com:remote-owner/remote-repo.git\n" };
        }
        throw new Error("Not found");
      });

      const details = await getRepoDetails();
      assert.deepStrictEqual(details, {
        owner: "remote-owner",
        repo: "remote-repo",
      });
    } finally {
      if (originalGhRepo !== undefined) {
        process.env.GH_REPO = originalGhRepo;
      }
      if (originalGithubRepo !== undefined) {
        process.env.GITHUB_REPOSITORY = originalGithubRepo;
      }
    }
  });
});

test("github-ci skill - getCurrentBranch and getGhCliToken unit tests", async (t) => {
  await t.test(
    "getCurrentBranch resolves current branch from git",
    async () => {
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        if (cmd === "git" && args.includes("--show-current")) {
          return { stdout: "feature/ci-fix\n" };
        }
        return { stdout: "" };
      });

      const branch = await getCurrentBranch();
      assert.strictEqual(branch, "feature/ci-fix");
    },
  );

  await t.test(
    "getCurrentBranch returns empty string on git failure",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => {
        throw new Error("git error");
      });

      const branch = await getCurrentBranch();
      assert.strictEqual(branch, "");
    },
  );

  await t.test("getGhCliToken extracts token from gh auth token", async () => {
    t.mock.method(sys, "execFileAsync", async (cmd, args) => {
      if (cmd === "gh" && args.includes("token")) {
        return { stdout: "gho_secret123\n" };
      }
      return { stdout: "" };
    });

    const token = await getGhCliToken();
    assert.strictEqual(token, "gho_secret123");
  });

  await t.test("getGhCliToken returns empty string on failure", async () => {
    t.mock.method(sys, "execFileAsync", async () => {
      throw new Error("gh not installed");
    });

    const token = await getGhCliToken();
    assert.strictEqual(token, "");
  });
});

test("github-ci skill - listRuns unit tests", async (t) => {
  await t.test(
    "calls octokit.rest.actions.listWorkflowRunsForRepo with options",
    async () => {
      let captured = null;
      const mockOctokit = {
        rest: {
          actions: {
            listWorkflowRunsForRepo: async (params) => {
              captured = params;
              return {
                data: {
                  workflow_runs: [
                    {
                      id: 101,
                      name: "CI",
                      status: "completed",
                      conclusion: "success",
                      head_branch: "main",
                    },
                  ],
                },
              };
            },
          },
        },
      };

      const runs = await listRuns({
        octokit: mockOctokit,
        owner: "rancher",
        repo: "terraform-provider-file",
        branch: "feature/test",
        limit: 5,
        status: "failure",
      });

      assert.strictEqual(runs.length, 1);
      assert.strictEqual(runs[0].id, 101);
      assert.deepStrictEqual(captured, {
        owner: "rancher",
        repo: "terraform-provider-file",
        branch: "feature/test",
        per_page: 5,
        status: "failure",
      });
    },
  );
});

test("github-ci skill - getFailedLog unit tests", async (t) => {
  await t.test("isolates failed jobs and extracts failure lines", async () => {
    const mockOctokit = {
      rest: {
        actions: {
          getWorkflowRun: async () => ({
            data: {
              id: 200,
              name: "Integration",
              status: "completed",
              conclusion: "failure",
            },
          }),
          listJobsForWorkflowRun: async () => ({
            data: {
              jobs: [
                { id: 301, name: "Lint", conclusion: "success" },
                { id: 302, name: "Test", conclusion: "failure" },
              ],
            },
          }),
          downloadJobLogsForWorkflowRun: async ({ job_id }) => {
            if (job_id === 302) {
              return {
                data: "Setting up Go\nRunning tests\nError: TestLocalDirectory failed with exit status 1\nDone",
              };
            }
            return { data: "All good" };
          },
        },
      },
    };

    const { run, failedJobs } = await getFailedLog({
      octokit: mockOctokit,
      owner: "rancher",
      repo: "terraform-provider-file",
      runId: 200,
    });

    assert.strictEqual(run.id, 200);
    assert.strictEqual(failedJobs.length, 1);
    assert.strictEqual(failedJobs[0].job.name, "Test");
    assert.strictEqual(failedJobs[0].errorLines.length, 1);
    assert.match(
      failedJobs[0].errorLines[0],
      /Error: TestLocalDirectory failed with exit status 1/,
    );
  });

  await t.test(
    "uses octokit.paginate to retrieve all jobs across pages",
    async () => {
      let paginateCalled = false;
      const mockOctokit = {
        paginate: async (method, params) => {
          paginateCalled = true;
          assert.strictEqual(params.run_id, 201);
          assert.strictEqual(params.per_page, 100);
          return [
            { id: 401, name: "Job 1", conclusion: "success" },
            { id: 402, name: "Job 2", conclusion: "failure" },
          ];
        },
        rest: {
          actions: {
            getWorkflowRun: async () => ({
              data: {
                id: 201,
                name: "Paginated Run",
                status: "completed",
                conclusion: "failure",
              },
            }),
            downloadJobLogsForWorkflowRun: async () => ({
              data: "Error: Job failed",
            }),
          },
        },
      };

      const { failedJobs } = await getFailedLog({
        octokit: mockOctokit,
        owner: "rancher",
        repo: "terraform-provider-file",
        runId: 201,
      });

      assert.strictEqual(paginateCalled, true);
      assert.strictEqual(failedJobs.length, 1);
      assert.strictEqual(failedJobs[0].job.id, 402);
    },
  );

  await t.test(
    "preserves job.steps containing failed steps in failedJobs output",
    async () => {
      const mockOctokit = {
        rest: {
          actions: {
            getWorkflowRun: async () => ({
              data: {
                id: 202,
                name: "Step Run",
                status: "completed",
                conclusion: "failure",
              },
            }),
            listJobsForWorkflowRun: async () => ({
              data: {
                jobs: [
                  {
                    id: 403,
                    name: "Build and Test",
                    conclusion: "failure",
                    steps: [
                      {
                        number: 1,
                        name: "Setup",
                        status: "completed",
                        conclusion: "success",
                      },
                      {
                        number: 2,
                        name: "Run test suite",
                        status: "completed",
                        conclusion: "failure",
                      },
                    ],
                  },
                ],
              },
            }),
            downloadJobLogsForWorkflowRun: async () => ({
              data: "Error: test suite failed with status 1",
            }),
          },
        },
      };

      const { failedJobs } = await getFailedLog({
        octokit: mockOctokit,
        owner: "rancher",
        repo: "terraform-provider-file",
        runId: 202,
      });

      assert.strictEqual(failedJobs.length, 1);
      assert.strictEqual(failedJobs[0].job.steps.length, 2);
      const failedStep = failedJobs[0].job.steps.find(
        (s) => s.conclusion === "failure",
      );
      assert.strictEqual(failedStep.name, "Run test suite");
      assert.strictEqual(failedStep.number, 2);
    },
  );

  await t.test(
    "preserves multiple failed steps across failure, timed_out, and startup_failure",
    async () => {
      const mockOctokit = {
        rest: {
          actions: {
            getWorkflowRun: async () => ({
              data: {
                id: 203,
                name: "Multi Step Run",
                status: "completed",
                conclusion: "failure",
              },
            }),
            listJobsForWorkflowRun: async () => ({
              data: {
                jobs: [
                  {
                    id: 404,
                    name: "Multi-fail Job",
                    conclusion: "failure",
                    steps: [
                      {
                        number: 1,
                        name: "Setup",
                        status: "completed",
                        conclusion: "success",
                      },
                      {
                        number: 2,
                        name: "Lint",
                        status: "completed",
                        conclusion: "failure",
                      },
                      {
                        number: 3,
                        name: "Test",
                        status: "completed",
                        conclusion: "timed_out",
                      },
                    ],
                  },
                ],
              },
            }),
            downloadJobLogsForWorkflowRun: async () => ({
              data: "Error: lint failed\nError: test timed out",
            }),
          },
        },
      };

      const { failedJobs } = await getFailedLog({
        octokit: mockOctokit,
        owner: "rancher",
        repo: "terraform-provider-file",
        runId: 203,
      });

      assert.strictEqual(failedJobs.length, 1);
      const failedSteps = failedJobs[0].job.steps.filter(
        (s) =>
          s.conclusion === "failure" ||
          s.conclusion === "timed_out" ||
          s.conclusion === "startup_failure",
      );
      assert.strictEqual(failedSteps.length, 2);
      assert.strictEqual(failedSteps[0].name, "Lint");
      assert.strictEqual(failedSteps[1].name, "Test");
    },
  );

  await t.test(
    "captures compilation and runner failures matching expanded regex",
    async () => {
      const mockOctokit = {
        rest: {
          actions: {
            getWorkflowRun: async () => ({
              data: {
                id: 204,
                name: "Compile Run",
                status: "completed",
                conclusion: "failure",
              },
            }),
            listJobsForWorkflowRun: async () => ({
              data: {
                jobs: [{ id: 405, name: "Build", conclusion: "failure" }],
              },
            }),
            downloadJobLogsForWorkflowRun: async () => ({
              data: "Compiling module\n./main.go:15:2: undefined: nonExistentFunc\n./main.go:20:5: cannot find package\n##[error]Process completed with exit code 2\n",
            }),
          },
        },
      };

      const { failedJobs } = await getFailedLog({
        octokit: mockOctokit,
        owner: "rancher",
        repo: "terraform-provider-file",
        runId: 204,
      });

      assert.strictEqual(failedJobs.length, 1);
      assert.strictEqual(failedJobs[0].errorLines.length, 3);
      assert.match(failedJobs[0].errorLines[0], /undefined: nonExistentFunc/);
      assert.match(failedJobs[0].errorLines[1], /cannot find package/);
      assert.match(failedJobs[0].errorLines[2], /##\[error\]Process completed/);
    },
  );

  await t.test(
    "captures tailLines fallback when no explicit regex lines match",
    async () => {
      const mockOctokit = {
        rest: {
          actions: {
            getWorkflowRun: async () => ({
              data: {
                id: 205,
                name: "Implicit Failure Run",
                status: "completed",
                conclusion: "failure",
              },
            }),
            listJobsForWorkflowRun: async () => ({
              data: {
                jobs: [{ id: 406, name: "Custom Tool", conclusion: "failure" }],
              },
            }),
            downloadJobLogsForWorkflowRun: async () => ({
              data: "Step A completed\nStep B starting\nAbnormal termination\n",
            }),
          },
        },
      };

      const { failedJobs } = await getFailedLog({
        octokit: mockOctokit,
        owner: "rancher",
        repo: "terraform-provider-file",
        runId: 205,
      });

      assert.strictEqual(failedJobs.length, 1);
      assert.strictEqual(failedJobs[0].errorLines.length, 0);
      assert.ok(failedJobs[0].tailLines.length >= 2);
      assert.strictEqual(
        failedJobs[0].tailLines[failedJobs[0].tailLines.length - 1],
        "Abnormal termination",
      );
    },
  );

  await t.test(
    "strips ANSI escape sequences and captures colorized error lines",
    async () => {
      const mockOctokit = {
        rest: {
          actions: {
            getWorkflowRun: async () => ({
              data: {
                id: 206,
                name: "Ansi Color Run",
                status: "completed",
                conclusion: "failure",
              },
            }),
            listJobsForWorkflowRun: async () => ({
              data: {
                jobs: [{ id: 407, name: "Compile", conclusion: "failure" }],
              },
            }),
            downloadJobLogsForWorkflowRun: async () => ({
              data: "\x1b[31;1merror: unhandled exception\x1b[0m\n\x1b[0;31m./pkg/mod.go:42:1: fatal: failed to parse\x1b[0m\n",
            }),
          },
        },
      };

      const { failedJobs } = await getFailedLog({
        octokit: mockOctokit,
        owner: "rancher",
        repo: "terraform-provider-file",
        runId: 206,
      });

      assert.strictEqual(failedJobs.length, 1);
      assert.strictEqual(failedJobs[0].errorLines.length, 2);
      assert.strictEqual(
        failedJobs[0].errorLines[0],
        "error: unhandled exception",
      );
      assert.strictEqual(
        failedJobs[0].errorLines[1],
        "./pkg/mod.go:42:1: fatal: failed to parse",
      );
    },
  );

  await t.test("parseFailedLogArgs parses run IDs, branches, and flags", () => {
    assert.deepStrictEqual(parseFailedLogArgs([]), {
      runId: null,
      targetBranch: undefined,
    });
    assert.deepStrictEqual(parseFailedLogArgs(["12345"]), {
      runId: 12345,
      targetBranch: undefined,
    });
    assert.deepStrictEqual(parseFailedLogArgs(["#67890"]), {
      runId: 67890,
      targetBranch: undefined,
    });
    assert.deepStrictEqual(parseFailedLogArgs(["-r", "9999"]), {
      runId: 9999,
      targetBranch: undefined,
    });
    assert.deepStrictEqual(parseFailedLogArgs(["--run=8888"]), {
      runId: 8888,
      targetBranch: undefined,
    });
    assert.deepStrictEqual(parseFailedLogArgs(["-b", "feature/test"]), {
      runId: null,
      targetBranch: "feature/test",
    });
    assert.deepStrictEqual(parseFailedLogArgs(["--branch=main"]), {
      runId: null,
      targetBranch: "main",
    });
    assert.deepStrictEqual(parseFailedLogArgs(["main"]), {
      runId: null,
      targetBranch: "main",
    });
    assert.deepStrictEqual(parseFailedLogArgs(["current", "12345"]), {
      runId: 12345,
      targetBranch: undefined,
    });
  });
});

test("github-ci skill - downloadLogs unit tests", async (t) => {
  await t.test(
    "uses octokit.paginate and downloads logs for all jobs",
    async () => {
      let paginateCalled = false;
      const mockOctokit = {
        paginate: async () => {
          paginateCalled = true;
          return [
            { id: 501, name: "Job A" },
            { id: 502, name: "Job B" },
          ];
        },
        rest: {
          actions: {
            downloadJobLogsForWorkflowRun: async () => ({
              data: "Sample log content",
            }),
          },
        },
      };

      const tempDir = path.join("/tmp", `test-ci-logs-${Date.now()}`);
      try {
        const result = await downloadLogs({
          octokit: mockOctokit,
          owner: "rancher",
          repo: "terraform-provider-file",
          runId: 999,
          outputDir: tempDir,
        });

        assert.strictEqual(paginateCalled, true);
        assert.strictEqual(result.jobsCount, 2);
      } finally {
        await (
          await import("node:fs/promises")
        ).rm(tempDir, { recursive: true, force: true });
      }
    },
  );

  await t.test(
    "filters to failed jobs only when failedOnly is true",
    async () => {
      const mockOctokit = {
        paginate: async () => [
          { id: 501, name: "Job A - Success", conclusion: "success" },
          { id: 502, name: "Job B - Failure", conclusion: "failure" },
          { id: 503, name: "Job C - Timeout", conclusion: "timed_out" },
        ],
        rest: {
          actions: {
            downloadJobLogsForWorkflowRun: async () => ({
              data: "Failure log content",
            }),
          },
        },
      };

      const tempDir = path.join("/tmp", `test-ci-logs-failed-${Date.now()}`);
      try {
        const result = await downloadLogs({
          octokit: mockOctokit,
          owner: "rancher",
          repo: "terraform-provider-file",
          runId: 1000,
          outputDir: tempDir,
          failedOnly: true,
        });

        assert.strictEqual(result.jobsCount, 2);
      } finally {
        await (
          await import("node:fs/promises")
        ).rm(tempDir, { recursive: true, force: true });
      }
    },
  );
});

test("github-ci skill - searchLogs unit tests", async (t) => {
  await t.test("searches directory and matches patterns", async () => {
    // Search the .gemini directory for 'name:'
    const results = await searchLogs({
      logsDir: ".gemini/skills/github-ci",
      pattern: /name:\s*github-ci/,
    });

    assert.ok(results.length >= 1);
    assert.match(results[0].text, /name:\s*github-ci/);
  });

  await t.test(
    "handles invalid/unescaped regex string without throwing SyntaxError",
    async () => {
      const brokenPattern = ["invalid", "_", "syntax", "("].join("");
      const results = await searchLogs({
        logsDir: ".gemini/skills/github-ci",
        pattern: brokenPattern,
      });
      assert.deepStrictEqual(results, []);
    },
  );

  await t.test(
    "matches lines with ANSI escape sequences and cleans output text",
    async () => {
      const tempDir = path.join("/tmp", `test-search-ansi-${Date.now()}`);
      const fsMod = await import("node:fs/promises");
      await fsMod.mkdir(tempDir, { recursive: true });
      try {
        await fsMod.writeFile(
          path.join(tempDir, "sample.log"),
          "\x1b[31merror: something went wrong\x1b[0m\n",
          "utf8",
        );
        const results = await searchLogs({
          logsDir: tempDir,
          pattern: "error: something went wrong",
        });
        assert.strictEqual(results.length, 1);
        assert.strictEqual(results[0].text, "error: something went wrong");
      } finally {
        await fsMod.rm(tempDir, { recursive: true, force: true });
      }
    },
  );

  await t.test(
    "stops searching and limits results when limit is specified",
    async () => {
      const tempDir = path.join("/tmp", `test-search-limit-${Date.now()}`);
      const fsMod = await import("node:fs/promises");
      await fsMod.mkdir(tempDir, { recursive: true });
      try {
        await fsMod.writeFile(
          path.join(tempDir, "many.log"),
          "error: 1\nerror: 2\nerror: 3\nerror: 4\nerror: 5\n",
          "utf8",
        );
        const results = await searchLogs({
          logsDir: tempDir,
          pattern: "error",
          limit: 2,
        });
        assert.strictEqual(results.length, 2);
      } finally {
        await fsMod.rm(tempDir, { recursive: true, force: true });
      }
    },
  );
});

test("github-ci skill - findFailedRuns unit tests", async (t) => {
  await t.test(
    "queries completed runs and filters by failure, timed_out, and startup_failure conclusions",
    async () => {
      let captured = null;
      const mockOctokit = {
        rest: {
          actions: {
            listWorkflowRunsForRepo: async (params) => {
              captured = params;
              return {
                data: {
                  workflow_runs: [
                    {
                      id: 1,
                      name: "Run 1",
                      status: "completed",
                      conclusion: "success",
                    },
                    {
                      id: 2,
                      name: "Run 2",
                      status: "completed",
                      conclusion: "failure",
                    },
                    {
                      id: 3,
                      name: "Run 3",
                      status: "completed",
                      conclusion: "timed_out",
                    },
                    {
                      id: 4,
                      name: "Run 4",
                      status: "completed",
                      conclusion: "startup_failure",
                    },
                    {
                      id: 5,
                      name: "Run 5",
                      status: "completed",
                      conclusion: "cancelled",
                    },
                  ],
                },
              };
            },
          },
        },
      };

      const failed = await findFailedRuns({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        branch: "main",
        limit: 10,
      });

      assert.strictEqual(failed.length, 3);
      assert.strictEqual(failed[0].id, 2);
      assert.strictEqual(failed[1].id, 3);
      assert.strictEqual(failed[2].id, 4);
      assert.strictEqual(captured.status, "completed");
      assert.strictEqual(captured.branch, "main");
    },
  );

  await t.test("supports custom status filter in findFailedRuns", async () => {
    let captured = null;
    const mockOctokit = {
      rest: {
        actions: {
          listWorkflowRunsForRepo: async (params) => {
            captured = params;
            return {
              data: {
                workflow_runs: [
                  {
                    id: 1,
                    name: "Run 1",
                    status: "completed",
                    conclusion: "success",
                  },
                  {
                    id: 5,
                    name: "Run 5",
                    status: "completed",
                    conclusion: "cancelled",
                  },
                ],
              },
            };
          },
        },
      },
    };

    const cancelled = await findFailedRuns({
      octokit: mockOctokit,
      owner: "testowner",
      repo: "testrepo",
      status: "cancelled",
    });

    assert.strictEqual(cancelled.length, 1);
    assert.strictEqual(cancelled[0].id, 5);
    assert.strictEqual(captured.status, "completed");

    await findFailedRuns({
      octokit: mockOctokit,
      owner: "testowner",
      repo: "testrepo",
      status: "in_progress",
    });
    assert.strictEqual(captured.status, "in_progress");
  });

  await t.test(
    "stops paginating at maximum 15 pages to protect rate limits",
    async () => {
      let pagesQueried = 0;
      const mockOctokit = {
        rest: {
          actions: {
            listWorkflowRunsForRepo: async (params) => {
              pagesQueried++;
              return {
                data: {
                  workflow_runs: new Array(params.per_page)
                    .fill(null)
                    .map((_, i) => ({
                      id: pagesQueried * 100 + i,
                      name: "Run",
                      status: "completed",
                      conclusion: "success",
                    })),
                },
              };
            },
          },
        },
      };

      const failed = await findFailedRuns({
        octokit: mockOctokit,
        owner: "testowner",
        repo: "testrepo",
        limit: 10,
      });

      assert.strictEqual(failed.length, 0);
      assert.strictEqual(pagesQueried, 15);
    },
  );
});

test("github-ci skill - timed_out jobs and stateless searchLogs", async (t) => {
  await t.test(
    "getFailedLog captures timed_out jobs and decodes Buffer data",
    async () => {
      const mockOctokit = {
        rest: {
          actions: {
            getWorkflowRun: async () => ({
              data: {
                id: 700,
                name: "Timeout Run",
                status: "completed",
                conclusion: "failure",
              },
            }),
            listJobsForWorkflowRun: async () => ({
              data: {
                jobs: [
                  { id: 801, name: "Build", conclusion: "timed_out" },
                  { id: 802, name: "Lint", conclusion: "startup_failure" },
                ],
              },
            }),
            downloadJobLogsForWorkflowRun: async ({ job_id }) => {
              if (job_id === 801) {
                return {
                  data: Buffer.from(
                    "Job execution timed out after 360m\nError: timed out",
                  ),
                };
              }
              return { data: "Fatal: container failed to initialize" };
            },
          },
        },
      };

      const { failedJobs } = await getFailedLog({
        octokit: mockOctokit,
        owner: "rancher",
        repo: "terraform-provider-file",
        runId: 700,
      });

      assert.strictEqual(failedJobs.length, 2);
      assert.strictEqual(failedJobs[0].job.name, "Build");
      assert.ok(failedJobs[0].errorLines.some((l) => l.includes("timed out")));
      assert.strictEqual(failedJobs[1].job.name, "Lint");
      assert.ok(
        failedJobs[1].errorLines.some((l) =>
          l.includes("Fatal: container failed"),
        ),
      );
    },
  );

  await t.test("searchLogs parses string /pattern/flags format", async () => {
    const results = await searchLogs({
      logsDir: ".gemini/skills/github-ci",
      pattern: "/name:\\s*github-ci/i",
    });

    assert.ok(results.length >= 1);
    assert.match(results[0].text, /name:\s*github-ci/i);
  });
});

test("github-ci skill - parseRunArgs and parseDownloadArgs CLI parsing unit tests", async (t) => {
  await t.test("parseRunArgs handles defaults and various flag formats", () => {
    assert.deepStrictEqual(parseRunArgs([]), {
      branchTarget: "current",
      limit: 10,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["20"]), {
      branchTarget: "current",
      limit: 20,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["-n", "25"]), {
      branchTarget: "current",
      limit: 25,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["--limit", "15"]), {
      branchTarget: "current",
      limit: 15,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["-n=30"]), {
      branchTarget: "current",
      limit: 30,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["--limit=50"]), {
      branchTarget: "current",
      limit: 50,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["main"]), {
      branchTarget: "main",
      limit: 10,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["main", "5"]), {
      branchTarget: "main",
      limit: 5,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["main", "-n", "5"]), {
      branchTarget: "main",
      limit: 5,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["all", "-n", "100"]), {
      branchTarget: "all",
      limit: 100,
      status: undefined,
    });
    assert.deepStrictEqual(parseRunArgs(["current", "50"]), {
      branchTarget: "current",
      limit: 50,
      status: undefined,
    });
    assert.deepStrictEqual(
      parseRunArgs(["-b", "feature-ci", "-s", "in_progress"]),
      {
        branchTarget: "feature-ci",
        limit: 10,
        status: "in_progress",
      },
    );
    assert.deepStrictEqual(
      parseRunArgs(["--branch=main", "--status=completed"]),
      {
        branchTarget: "main",
        limit: 10,
        status: "completed",
      },
    );
  });

  await t.test(
    "parseDownloadArgs handles run IDs, directories, and flags",
    () => {
      assert.deepStrictEqual(parseDownloadArgs([]), {
        runId: null,
        outputDir: null,
        failedOnly: false,
      });
      assert.deepStrictEqual(parseDownloadArgs(["12345"]), {
        runId: 12345,
        outputDir: "/tmp/ci-logs-12345",
        failedOnly: false,
      });
      assert.deepStrictEqual(parseDownloadArgs(["#12345"]), {
        runId: 12345,
        outputDir: "/tmp/ci-logs-12345",
        failedOnly: false,
      });
      assert.deepStrictEqual(parseDownloadArgs(["12345", "/custom/path"]), {
        runId: 12345,
        outputDir: "/custom/path",
        failedOnly: false,
      });
      assert.deepStrictEqual(parseDownloadArgs(["12345", "logs"]), {
        runId: 12345,
        outputDir: "logs",
        failedOnly: false,
      });
      assert.deepStrictEqual(parseDownloadArgs(["current", "/custom/path"]), {
        runId: null,
        outputDir: "/custom/path",
        failedOnly: false,
      });
      assert.deepStrictEqual(parseDownloadArgs(["/custom/path"]), {
        runId: null,
        outputDir: "/custom/path",
        failedOnly: false,
      });
      assert.deepStrictEqual(
        parseDownloadArgs(["-o", "/custom/path", "12345"]),
        {
          runId: 12345,
          outputDir: "/custom/path",
          failedOnly: false,
        },
      );
      assert.deepStrictEqual(parseDownloadArgs(["--output=/custom/path"]), {
        runId: null,
        outputDir: "/custom/path",
        failedOnly: false,
      });
      assert.deepStrictEqual(parseDownloadArgs(["12345", "--failed-only"]), {
        runId: 12345,
        outputDir: "/tmp/ci-logs-12345",
        failedOnly: true,
      });
      assert.deepStrictEqual(parseDownloadArgs(["12345", "-f"]), {
        runId: 12345,
        outputDir: "/tmp/ci-logs-12345",
        failedOnly: true,
      });
    },
  );

  await t.test("parseSearchArgs extracts logsDir, pattern, and limit", () => {
    assert.deepStrictEqual(parseSearchArgs([]), {
      logsDir: undefined,
      pattern: undefined,
      limit: 100,
    });
    assert.deepStrictEqual(parseSearchArgs(["/tmp/logs"]), {
      logsDir: "/tmp/logs",
      pattern: undefined,
      limit: 100,
    });
    assert.deepStrictEqual(parseSearchArgs(["/tmp/logs", "fatal"]), {
      logsDir: "/tmp/logs",
      pattern: "fatal",
      limit: 100,
    });
    assert.deepStrictEqual(
      parseSearchArgs(["-n", "25", "/tmp/logs", "fatal"]),
      {
        logsDir: "/tmp/logs",
        pattern: "fatal",
        limit: 25,
      },
    );
    assert.deepStrictEqual(parseSearchArgs(["--limit=50", "/tmp/logs"]), {
      logsDir: "/tmp/logs",
      pattern: undefined,
      limit: 50,
    });
  });
});
