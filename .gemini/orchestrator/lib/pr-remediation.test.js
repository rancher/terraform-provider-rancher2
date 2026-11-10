import assert from "node:assert";
import test from "node:test";
import {
  appendRemediationToPlan,
  buildPRReviewImplementationPrompt,
  fetchUnresolvedPRReviewThreads,
  filterActionableReviewComments,
  formatRemediationPlanSection,
  isWindowsCompatibilityComment,
  normalizeReviewThreads,
  replyToPRReviewThread,
  resolvePRReviewThread,
  sys,
} from "./pr-remediation.js";

test("pr-remediation - isWindowsCompatibilityComment", () => {
  assert.strictEqual(isWindowsCompatibilityComment(""), false);
  assert.strictEqual(isWindowsCompatibilityComment(null), false);
  assert.strictEqual(
    isWindowsCompatibilityComment("Please add unit tests for error handling."),
    false,
  );

  assert.strictEqual(
    isWindowsCompatibilityComment(
      "Consider PowerShell support for this script.",
    ),
    true,
  );
  assert.strictEqual(
    isWindowsCompatibilityComment("Ensure cmd.exe compatibility on Windows."),
    true,
  );
  assert.strictEqual(
    isWindowsCompatibilityComment("Check backslash path separators on Win32."),
    true,
  );
  assert.strictEqual(
    isWindowsCompatibilityComment("What if someone runs .bat or .cmd files?"),
    true,
  );
  assert.strictEqual(
    isWindowsCompatibilityComment("Please provide a setup.ps1 script."),
    true,
  );
  assert.strictEqual(
    isWindowsCompatibilityComment("Handle drive letter like C:\\path"),
    true,
  );
  assert.strictEqual(
    isWindowsCompatibilityComment("Or another drive like D:\\logs"),
    true,
  );
});

test("pr-remediation - normalizeReviewThreads", () => {
  assert.deepStrictEqual(normalizeReviewThreads(null), []);
  assert.deepStrictEqual(normalizeReviewThreads([]), []);

  const rawThreads = [
    {
      id: "PRRT_kwDOPahR0c6y3B6l",
      isResolved: false,
      comments: {
        nodes: [
          {
            databaseId: 101,
            author: { login: "copilot" },
            body: "Consider decoupling repoRoot from cwd.",
            path: "agent-scripts/orchestrator.js",
            line: 42,
            originalLine: 40,
          },
        ],
      },
    },
    {
      id: null,
      isResolved: true,
      comments: {
        nodes: [
          {
            databaseId: 202,
            author: null,
            body: "REST fallback comment",
            path: "README.md",
            line: null,
            originalLine: 5,
          },
        ],
      },
    },
    {
      // Missing ID and databaseId -> should be filtered out
      id: null,
      comments: {
        nodes: [],
      },
    },
  ];

  const normalized = normalizeReviewThreads(rawThreads);
  assert.strictEqual(normalized.length, 2);

  assert.strictEqual(normalized[0].threadId, "PRRT_kwDOPahR0c6y3B6l");
  assert.strictEqual(normalized[0].isResolved, false);
  assert.strictEqual(normalized[0].path, "agent-scripts/orchestrator.js");
  assert.strictEqual(normalized[0].line, 42);
  assert.strictEqual(normalized[0].author, "copilot");
  assert.strictEqual(normalized[0].comments.length, 1);
  assert.strictEqual(normalized[0].comments[0].id, 101);

  assert.strictEqual(normalized[1].threadId, "202");
  assert.strictEqual(normalized[1].isResolved, true);
  assert.strictEqual(normalized[1].line, 5);
  assert.strictEqual(normalized[1].author, "unknown");
  assert.strictEqual(normalized[1].comments[0].author, "ghost");
});

test("pr-remediation - fetchUnresolvedPRReviewThreads", async (t) => {
  await t.test(
    "successfully executes read_comments.js and parses threads",
    async () => {
      let capturedArgs = null;
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        capturedArgs = args;
        return {
          stdout: JSON.stringify({
            prNumber: 452,
            reviewThreads: [
              {
                id: "PRRT_123",
                isResolved: false,
                comments: {
                  nodes: [
                    {
                      databaseId: 999,
                      author: { login: "reviewer1" },
                      body: "Fix logic error here.",
                      path: "main.go",
                      line: 15,
                    },
                  ],
                },
              },
            ],
          }),
        };
      });

      const result = await fetchUnresolvedPRReviewThreads(452, "/mock/repo");
      assert.strictEqual(result.prNumber, 452);
      assert.strictEqual(result.threads.length, 1);
      assert.strictEqual(result.threads[0].threadId, "PRRT_123");
      assert.ok(capturedArgs.includes("--unresolved"));
      assert.ok(capturedArgs.includes("--json"));
    },
  );

  await t.test(
    'falls back to "current" and handles null prTarget safely',
    async () => {
      let capturedArgs = null;
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        capturedArgs = args;
        return {
          stdout: JSON.stringify({
            reviewThreads: [],
          }),
        };
      });

      const result = await fetchUnresolvedPRReviewThreads(null, "/mock/repo");
      assert.strictEqual(result.prNumber, null);
      assert.strictEqual(capturedArgs[1], "current");
    },
  );

  await t.test(
    "parses numeric prTarget safely when json payload omits prNumber",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => ({
        stdout: JSON.stringify({ reviewThreads: [] }),
      }));

      const result = await fetchUnresolvedPRReviewThreads("789", "/mock/repo");
      assert.strictEqual(result.prNumber, 789);
    },
  );

  await t.test(
    "handles non-numeric or edge-case prTarget without yielding 0 or unexpected numbers",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => ({
        stdout: JSON.stringify({ reviewThreads: [] }),
      }));

      const resultWithArray = await fetchUnresolvedPRReviewThreads(
        [],
        "/mock/repo",
      );
      assert.strictEqual(resultWithArray.prNumber, null);

      const resultWithInvalidStr = await fetchUnresolvedPRReviewThreads(
        "abc",
        "/mock/repo",
      );
      assert.strictEqual(resultWithInvalidStr.prNumber, null);
    },
  );

  await t.test(
    "throws formatted SyntaxError when JSON parsing fails",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => ({
        stdout: "not-valid-json {",
      }));

      await assert.rejects(
        async () => {
          await fetchUnresolvedPRReviewThreads("current", "/mock/repo");
        },
        (err) => {
          assert.match(
            err.message,
            /Failed to parse PR review comments JSON\. stdout: not-valid-json \{/,
          );
          assert.ok(err.cause instanceof SyntaxError);
          return true;
        },
      );
    },
  );

  await t.test(
    "propagates error when child process execution fails",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => {
        const err = new Error("Execution failed");
        err.stderr = "API rate limit exceeded";
        throw err;
      });

      await assert.rejects(
        async () => {
          await fetchUnresolvedPRReviewThreads("current");
        },
        {
          message:
            /Failed to fetch unresolved PR review comments: API rate limit exceeded/,
        },
      );
    },
  );
});

test("pr-remediation - filterActionableReviewComments", () => {
  const threads = [
    {
      threadId: "T1",
      body: "Standard bug in error handling.",
      isResolved: false,
    },
    {
      threadId: "T2",
      body: "Please make sure this works with powershell.exe on Windows.",
      isResolved: false,
    },
    {
      threadId: "T3",
      body: "This thread is resolved.",
      isResolved: true,
    },
  ];

  const { actionable, ignored } = filterActionableReviewComments(threads, {
    ignoreWindows: true,
  });
  assert.strictEqual(actionable.length, 1);
  assert.strictEqual(actionable[0].threadId, "T1");

  assert.strictEqual(ignored.length, 2);
  assert.strictEqual(ignored[0].threadId, "T2");
  assert.match(ignored[0].ignoreReason, /Windows compatibility/);
  assert.strictEqual(ignored[1].threadId, "T3");
  assert.match(ignored[1].ignoreReason, /already resolved/);

  // When ignoreWindows is false, Windows comments remain actionable
  const withWindows = filterActionableReviewComments(threads, {
    ignoreWindows: false,
  });
  assert.strictEqual(withWindows.actionable.length, 2);

  // Child comment with Windows compatibility pattern is ignored even if root body is POSIX
  const threadWithWindowsChild = [
    {
      threadId: "T4",
      body: "Improve validation of input argument.",
      isResolved: false,
      comments: [
        {
          id: 1,
          author: "alice",
          body: "Improve validation of input argument.",
        },
        {
          id: 2,
          author: "bob",
          body: "Also please test on powershell on Windows.",
        },
      ],
    },
  ];
  const childFiltered = filterActionableReviewComments(threadWithWindowsChild, {
    ignoreWindows: true,
  });
  assert.strictEqual(childFiltered.actionable.length, 0);
  assert.strictEqual(childFiltered.ignored.length, 1);
  assert.strictEqual(childFiltered.ignored[0].threadId, "T4");
  assert.match(childFiltered.ignored[0].ignoreReason, /Windows compatibility/);

  // Safely handles null threads
  const nullResult = filterActionableReviewComments(null);
  assert.strictEqual(nullResult.actionable.length, 0);
  assert.strictEqual(nullResult.ignored.length, 0);
});

test("pr-remediation - formatRemediationPlanSection", () => {
  const emptySection = formatRemediationPlanSection({
    prNumber: 452,
    reviewThreads: [],
  });
  assert.strictEqual(
    emptySection,
    "## 🔄 PR #452 REVIEW REMEDIATION (Round 1)\n- No unresolved actionable review comments identified.",
  );

  const nullThreadsSection = formatRemediationPlanSection({
    prNumber: 452,
    reviewThreads: null,
  });
  assert.strictEqual(
    nullThreadsSection,
    "## 🔄 PR #452 REVIEW REMEDIATION (Round 1)\n- No unresolved actionable review comments identified.",
  );

  const threads = [
    {
      threadId: "PRRT_abc",
      path: "internal/file.go",
      line: 35,
      author: "alice",
      body: "Improve validation of input argument.\nSecond line explanation.",
    },
    {
      threadId: "PRRT_xyz",
      path: "README.md",
      line: null,
      author: null,
      body: "Fix typo in documentation heading.",
    },
    {
      threadId: "PRRT_zero",
      path: "internal/zero.go",
      line: 0,
      author: "carol",
      body: "Header comment at top of file.",
    },
  ];

  const formatted = formatRemediationPlanSection({
    prNumber: 452,
    reviewThreads: threads,
    round: 2,
  });
  assert.ok(formatted.startsWith("## 🔄 PR #452 REVIEW REMEDIATION (Round 2)"));
  assert.ok(
    formatted.includes(
      '- [ ] Address Review Comment 1 (Thread PRRT_abc): internal/file.go:35 by @alice - "Improve validation of input argument."',
    ),
  );
  assert.ok(
    formatted.includes(
      '- [ ] Address Review Comment 2 (Thread PRRT_xyz): README.md by reviewer - "Fix typo in documentation heading."',
    ),
  );
  assert.ok(
    formatted.includes(
      '- [ ] Address Review Comment 3 (Thread PRRT_zero): internal/zero.go:0 by @carol - "Header comment at top of file."',
    ),
  );
});

test("pr-remediation - appendRemediationToPlan", async (t) => {
  let writtenData = null;

  await t.test("throws when planPath is missing", async () => {
    await assert.rejects(
      async () => appendRemediationToPlan("", "## Section"),
      {
        message: /Plan path must be specified/,
      },
    );
  });

  await t.test("returns appended: false when section is empty", async () => {
    const res = await appendRemediationToPlan("/path/plan.md", "   ");
    assert.strictEqual(res.appended, false);
    assert.strictEqual(res.reason, "Empty remediation section.");
  });

  await t.test("appends section safely to existing plan", async () => {
    t.mock.method(
      sys,
      "readFile",
      async () => "# Existing Plan\n\nSome task.\n",
    );
    t.mock.method(sys, "writeFile", async (filePath, data) => {
      writtenData = data;
    });

    const res = await appendRemediationToPlan(
      "/path/plan.md",
      "## 🔄 PR #452 REVIEW REMEDIATION (Round 1)\n- [ ] Task 1",
    );
    assert.strictEqual(res.appended, true);
    assert.strictEqual(
      writtenData,
      "# Existing Plan\n\nSome task.\n\n## 🔄 PR #452 REVIEW REMEDIATION (Round 1)\n- [ ] Task 1\n",
    );
  });

  await t.test(
    "appends section cleanly to empty plan without leading newlines",
    async () => {
      t.mock.method(sys, "readFile", async () => "   \n\n  ");
      t.mock.method(sys, "writeFile", async (filePath, data) => {
        writtenData = data;
      });

      const res = await appendRemediationToPlan(
        "/path/plan.md",
        "## 🔄 PR #452 REVIEW REMEDIATION (Round 1)\n- [ ] Task 1",
      );
      assert.strictEqual(res.appended, true);
      assert.strictEqual(
        writtenData,
        "## 🔄 PR #452 REVIEW REMEDIATION (Round 1)\n- [ ] Task 1\n",
      );
    },
  );

  await t.test(
    "skips appending if section header is already present",
    async () => {
      t.mock.method(
        sys,
        "readFile",
        async () =>
          "# Existing Plan\n\n## 🔄 PR #452 REVIEW REMEDIATION (Round 1)\n- [ ] Task 1\n",
      );

      const res = await appendRemediationToPlan(
        "/path/plan.md",
        "## 🔄 PR #452 REVIEW REMEDIATION (Round 1)\n- [ ] Duplicate",
      );
      assert.strictEqual(res.appended, false);
      assert.strictEqual(
        res.reason,
        "Remediation section header already exists in plan.",
      );
    },
  );
});

test("pr-remediation - replyToPRReviewThread", async (t) => {
  await t.test("validates required parameters", async () => {
    await assert.rejects(async () => replyToPRReviewThread({ threadId: "" }), {
      message: /Thread ID is required/,
    });
    await assert.rejects(
      async () =>
        replyToPRReviewThread({ threadId: "PRRT_123", replyMessage: " " }),
      {
        message: /Reply message cannot be empty/,
      },
    );
  });

  await t.test("executes reply-comment.js successfully", async () => {
    let capturedArgs = null;
    t.mock.method(sys, "execFileAsync", async (cmd, args) => {
      capturedArgs = args;
      return { stdout: "Reply posted", stderr: "" };
    });

    const res = await replyToPRReviewThread({
      prTarget: 452,
      threadId: "PRRT_123",
      replyMessage: "Fixed in commit abc1234",
      repoRoot: "/mock/repo",
    });

    assert.strictEqual(res.success, true);
    assert.ok(capturedArgs[0].endsWith("reply-comment.js"));
    assert.strictEqual(capturedArgs[1], "PRRT_123");
    assert.strictEqual(capturedArgs[2], "--pr");
    assert.strictEqual(capturedArgs[3], "452");
    assert.strictEqual(capturedArgs[4], "-m");
    assert.strictEqual(capturedArgs[5], "Fixed in commit abc1234");
  });

  await t.test(
    "handles numeric threadId and prTarget without ambiguity",
    async () => {
      let capturedArgs = null;
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        capturedArgs = args;
        return { stdout: "Reply posted", stderr: "" };
      });

      const res = await replyToPRReviewThread({
        prTarget: 1000,
        threadId: 50,
        replyMessage: "Fixed numeric ID",
        repoRoot: "/mock/repo",
      });

      assert.strictEqual(res.success, true);
      assert.strictEqual(capturedArgs[1], "50");
      assert.strictEqual(capturedArgs[2], "--pr");
      assert.strictEqual(capturedArgs[3], "1000");
      assert.strictEqual(capturedArgs[4], "-m");
      assert.strictEqual(capturedArgs[5], "Fixed numeric ID");
    },
  );

  await t.test('falls back to "current" when prTarget is null', async () => {
    let capturedArgs = null;
    t.mock.method(sys, "execFileAsync", async (cmd, args) => {
      capturedArgs = args;
      return { stdout: "Reply posted", stderr: "" };
    });

    const res = await replyToPRReviewThread({
      prTarget: null,
      threadId: "PRRT_123",
      replyMessage: "Fixed",
      repoRoot: "/mock/repo",
    });

    assert.strictEqual(res.success, true);
    assert.strictEqual(capturedArgs[1], "PRRT_123");
    assert.strictEqual(capturedArgs[2], "--pr");
    assert.strictEqual(capturedArgs[3], "current");
  });

  await t.test("propagates error when child process fails", async () => {
    t.mock.method(sys, "execFileAsync", async () => {
      const err = new Error("Process exited with 1");
      err.stderr = "Unauthorized access";
      throw err;
    });

    await assert.rejects(
      async () => {
        await replyToPRReviewThread({
          threadId: "PRRT_123",
          replyMessage: "Fixed",
        });
      },
      {
        message:
          /Failed to reply to PR review thread PRRT_123: Unauthorized access/,
      },
    );
  });

  await t.test(
    "falls back to stdout diagnostics when stderr is empty on failure",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => {
        const err = new Error("Process exited with 1");
        err.stderr = "";
        err.stdout = "Error: GraphQL query failed";
        throw err;
      });

      await assert.rejects(
        async () => {
          await replyToPRReviewThread({
            threadId: "PRRT_123",
            replyMessage: "Fixed",
          });
        },
        {
          message:
            /Failed to reply to PR review thread PRRT_123: Error: GraphQL query failed/,
        },
      );
    },
  );

  await t.test(
    "falls back to stdout diagnostics when stderr is only whitespace",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => {
        const err = new Error("Process exited with 1");
        err.stderr = "   \n\t  ";
        err.stdout = "Error: GraphQL query failed";
        throw err;
      });

      await assert.rejects(
        async () => {
          await replyToPRReviewThread({
            threadId: "PRRT_123",
            replyMessage: "Fixed",
          });
        },
        {
          message:
            /Failed to reply to PR review thread PRRT_123: Error: GraphQL query failed/,
        },
      );
    },
  );
});

test("pr-remediation - resolvePRReviewThread", async (t) => {
  await t.test("validates required parameters", async () => {
    await assert.rejects(async () => resolvePRReviewThread({ threadId: "" }), {
      message: /Thread ID is required/,
    });
  });

  await t.test("executes resolve_comment.js successfully", async () => {
    let capturedArgs = null;
    t.mock.method(sys, "execFileAsync", async (cmd, args) => {
      capturedArgs = args;
      return { stdout: "Thread resolved", stderr: "" };
    });

    const res = await resolvePRReviewThread({
      prTarget: "current",
      threadId: "PRRT_999",
      repoRoot: "/mock/repo",
    });

    assert.strictEqual(res.success, true);
    assert.ok(capturedArgs[0].endsWith("resolve_comment.js"));
    assert.strictEqual(capturedArgs[1], "PRRT_999");
    assert.strictEqual(capturedArgs[2], "--pr");
    assert.strictEqual(capturedArgs[3], "current");
  });

  await t.test(
    "handles numeric threadId and prTarget without ambiguity in resolvePRReviewThread",
    async () => {
      let capturedArgs = null;
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        capturedArgs = args;
        return { stdout: "Thread resolved", stderr: "" };
      });

      const res = await resolvePRReviewThread({
        prTarget: 1000,
        threadId: 50,
        repoRoot: "/mock/repo",
      });

      assert.strictEqual(res.success, true);
      assert.strictEqual(capturedArgs[1], "50");
      assert.strictEqual(capturedArgs[2], "--pr");
      assert.strictEqual(capturedArgs[3], "1000");
    },
  );

  await t.test('falls back to "current" when prTarget is null', async () => {
    let capturedArgs = null;
    t.mock.method(sys, "execFileAsync", async (cmd, args) => {
      capturedArgs = args;
      return { stdout: "Thread resolved", stderr: "" };
    });

    const res = await resolvePRReviewThread({
      prTarget: null,
      threadId: "PRRT_999",
      repoRoot: "/mock/repo",
    });

    assert.strictEqual(res.success, true);
    assert.strictEqual(capturedArgs[1], "PRRT_999");
    assert.strictEqual(capturedArgs[2], "--pr");
    assert.strictEqual(capturedArgs[3], "current");
  });

  await t.test("propagates error when child process fails", async () => {
    t.mock.method(sys, "execFileAsync", async () => {
      const err = new Error("Process exited with 1");
      err.stderr = "Thread not found";
      throw err;
    });

    await assert.rejects(
      async () => {
        await resolvePRReviewThread({
          threadId: "PRRT_999",
        });
      },
      {
        message:
          /Failed to resolve PR review thread PRRT_999: Thread not found/,
      },
    );
  });

  await t.test(
    "falls back to stdout diagnostics when stderr is empty on failure",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => {
        const err = new Error("Process exited with 1");
        err.stderr = "";
        err.stdout = "Error: Mutation resolution failed";
        throw err;
      });

      await assert.rejects(
        async () => {
          await resolvePRReviewThread({
            threadId: "PRRT_999",
          });
        },
        {
          message:
            /Failed to resolve PR review thread PRRT_999: Error: Mutation resolution failed/,
        },
      );
    },
  );

  await t.test(
    "falls back to stdout diagnostics when stderr is only whitespace",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => {
        const err = new Error("Process exited with 1");
        err.stderr = "   \r\n  ";
        err.stdout = "Error: Mutation resolution failed";
        throw err;
      });

      await assert.rejects(
        async () => {
          await resolvePRReviewThread({
            threadId: "PRRT_999",
          });
        },
        {
          message:
            /Failed to resolve PR review thread PRRT_999: Error: Mutation resolution failed/,
        },
      );
    },
  );
});

test("pr-remediation - buildPRReviewImplementationPrompt", () => {
  const prompt = buildPRReviewImplementationPrompt({
    prNumber: 452,
    threads: [
      {
        threadId: "PRRT_001",
        path: "main.go",
        line: 10,
        author: "bob",
        body: "Check return value here.",
      },
    ],
    activePlan: "# Active Plan\n\nTask list.",
  });

  assert.ok(
    prompt.includes('Objective: "Address PR review comments on PR #452"'),
  );
  assert.ok(prompt.includes("Comment 1 (Thread ID: `PRRT_001`)"));
  assert.ok(prompt.includes("main.go:10"));
  assert.ok(prompt.includes("@bob"));
  assert.ok(prompt.includes("Check return value here."));
  assert.ok(prompt.includes("reply-comment.js <thread_id> --pr 452 -m"));
  assert.ok(prompt.includes("resolve_comment.js <thread_id> --pr 452"));
  assert.ok(
    prompt.includes("Windows is strictly out of scope; POSIX target only"),
  );
  assert.ok(
    prompt.includes(
      "<active_plan>\n# Active Plan\n\nTask list.\n</active_plan>",
    ),
  );
});
