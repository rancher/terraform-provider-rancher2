import assert from "node:assert";
import test from "node:test";

import {
  buildImplementationPrompt,
  buildReviewerPrompt,
  buildReviewPrompt,
  buildBatchSuggestionImplementationPrompt,
  buildTallyPrompt,
  calculateScaledTurns,
  captureGitState,
  detectGitStateChange,
  extractCompletedSuggestionIds,
  formatValidationCommands,
  formatValidationSpillover,
  GIT_EMPTY_TREE_HASH,
  parseIndividualSuggestions,
  parseSuggestionCounts,
  runImplementationRalphLoop,
  runRalphLoop,
  runReviewRalphLoop,
  tallySuggestionsWithFlashLite,
} from "./ralph-loop.js";

test("detectGitStateChange state comparison logic", async (t) => {
  await t.test("detects when diff hash changes", () => {
    const before = { diffHash: "hash-1", status: "", plan: "" };
    const after = { diffHash: "hash-2", status: "", plan: "" };
    const res = detectGitStateChange(before, after);
    assert.strictEqual(res.hasChanged, true);
    assert.strictEqual(res.diffChanged, true);
    assert.strictEqual(res.statusChanged, false);
    assert.strictEqual(res.planChanged, false);
  });

  await t.test("detects when git porcelain status changes", () => {
    const before = {
      diffHash: "hash-1",
      status: " M file.go",
      plan: "plan text",
    };
    const after = {
      diffHash: "hash-1",
      status: "M  file.go",
      plan: "plan text",
    };
    const res = detectGitStateChange(before, after);
    assert.strictEqual(res.hasChanged, true);
    assert.strictEqual(res.diffChanged, false);
    assert.strictEqual(res.statusChanged, true);
    assert.strictEqual(res.planChanged, false);
  });

  await t.test("detects when plan file content changes", () => {
    const before = { diffHash: "hash-1", status: "", plan: "old plan" };
    const after = {
      diffHash: "hash-1",
      status: "",
      plan: "new plan with decisions",
    };
    const res = detectGitStateChange(before, after);
    assert.strictEqual(res.hasChanged, true);
    assert.strictEqual(res.diffChanged, false);
    assert.strictEqual(res.statusChanged, false);
    assert.strictEqual(res.planChanged, true);
  });

  await t.test(
    "returns hasChanged: false when state snapshots are identical",
    () => {
      const before = {
        diffHash: "hash-1",
        status: "M main.go\n",
        plan: "same plan",
      };
      const after = {
        diffHash: "hash-1",
        status: "M main.go",
        plan: "same plan",
      };
      const res = detectGitStateChange(before, after);
      assert.strictEqual(res.hasChanged, false);
      assert.strictEqual(res.diffChanged, false);
      assert.strictEqual(res.statusChanged, false);
      assert.strictEqual(res.planChanged, false);
    },
  );

  await t.test("handles missing or null snapshots gracefully", () => {
    const res1 = detectGitStateChange(null, { diffHash: "1" });
    assert.strictEqual(res1.hasChanged, true);

    const res2 = detectGitStateChange({ diffHash: "1" }, null);
    assert.strictEqual(res2.hasChanged, true);
  });
});

test("buildReviewPrompt formatting and feedback injection", async (t) => {
  await t.test(
    "formats standard review prompt with plan, diff, and modified files",
    () => {
      const prompt = buildReviewPrompt({
        iteration: 1,
        maxIterations: 10,
        stateBefore: {
          diff: "diff --git a/main.go b/main.go",
          modifiedFilesContext: '<file path="main.go">package main</file>',
          plan: "## Objective: Review and improve",
        },
      });

      assert.ok(prompt.includes("Review iteration 1/10."));
      assert.ok(
        prompt.includes(
          "<active_plan>\n## Objective: Review and improve\n</active_plan>",
        ),
      );
      assert.ok(
        prompt.includes(
          "<git_diff>\ndiff --git a/main.go b/main.go\n</git_diff>",
        ),
      );
      assert.ok(
        prompt.includes(
          '<modified_files>\n<file path="main.go">package main</file>\n</modified_files>',
        ),
      );
      assert.ok(!prompt.includes("<validation_failure>"));
    },
  );

  await t.test("injects negative feedback when lastValidation failed", () => {
    const prompt = buildReviewPrompt({
      iteration: 2,
      maxIterations: 10,
      stateBefore: {
        diff: "some diff",
        modifiedFilesContext: "",
        plan: "plan",
      },
      lastValidation: {
        success: false,
        stdout: "golangci-lint found 1 error",
        stderr: "main.go:10:5: undefined variable x",
      },
    });

    assert.ok(
      prompt.includes(
        "⚠️ CRITICAL VALIDATION FEEDBACK (Negative Backpressure):",
      ),
    );
    assert.ok(prompt.includes("<validation_failure>"));
    assert.ok(prompt.includes("STDOUT:\ngolangci-lint found 1 error"));
    assert.ok(
      prompt.includes("STDERR / ERRORS:\nmain.go:10:5: undefined variable x"),
    );
  });

  await t.test(
    "injects validation status notice when lastValidation succeeded",
    () => {
      const prompt = buildReviewPrompt({
        iteration: 2,
        maxIterations: 10,
        stateBefore: { diff: "diff", modifiedFilesContext: "", plan: "plan" },
        lastValidation: { success: true },
      });

      assert.ok(prompt.includes("<validation_status>"));
      assert.ok(
        prompt.includes(
          "🟢 All workspace tests and linters passed cleanly in the previous iteration.",
        ),
      );
      assert.ok(!prompt.includes("<validation_failure>"));
    },
  );
});

test("captureGitState mock execution", async (t) => {
  await t.test(
    "captures status, diff, hash, and files using custom exec function",
    async () => {
      const executedCommands = [];
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        executedCommands.push(fullCmd);
        if (fullCmd.includes("status --porcelain")) {
          return { stdout: "M  internal/provider.go\n" };
        }
        if (
          fullCmd.includes("diff --no-ext-diff") &&
          fullCmd.includes("HEAD") &&
          !fullCmd.includes("--name-only")
        ) {
          return {
            stdout:
              "diff --git a/internal/provider.go b/internal/provider.go\nindex abc..def 100644\n",
          };
        }
        if (fullCmd.includes("--name-only") && fullCmd.includes("HEAD")) {
          return { stdout: "internal/provider.go\n" };
        }
        return { stdout: "" };
      };

      const state = await captureGitState({
        repoRoot: process.cwd(),
        planPath: "/nonexistent/plan.md",
        stageBeforeIteration: false,
        execFn: mockExec,
      });

      assert.strictEqual(state.status, "M  internal/provider.go\n");
      assert.ok(state.diff.includes("diff --git"));
      assert.ok(
        typeof state.diffHash === "string" && state.diffHash.length === 64,
      );
      assert.deepStrictEqual(state.changedFiles, ["internal/provider.go"]);
      assert.ok(
        executedCommands.some((c) => c.includes("core.quotepath=false")),
      );
      assert.ok(executedCommands.some((c) => c.includes("status --porcelain")));
      assert.ok(
        !executedCommands.some((c) =>
          c.includes("status --porcelain --no-color"),
        ),
      );
      assert.ok(executedCommands.some((c) => c.includes("--no-color")));
    },
  );

  await t.test(
    "falls back to --cached when HEAD is unborn (initial commit / empty branch)",
    async () => {
      const executedCommands = [];
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        executedCommands.push(fullCmd);
        if (fullCmd.includes("status --porcelain")) {
          return { stdout: "A  newfile.go\n" };
        }
        // HEAD diff throws because HEAD is unborn
        if (fullCmd.includes("HEAD")) {
          throw new Error(
            'fatal: ambiguous argument "HEAD": unknown revision or path not in the working tree.',
          );
        }
        if (fullCmd.includes("--cached") && !fullCmd.includes("--name-only")) {
          return {
            stdout:
              "diff --git a/newfile.go b/newfile.go\nnew file mode 100644\n",
          };
        }
        if (fullCmd.includes("--cached") && fullCmd.includes("--name-only")) {
          return { stdout: "newfile.go\n" };
        }
        return { stdout: "" };
      };

      const state = await captureGitState({
        repoRoot: process.cwd(),
        planPath: "/nonexistent/plan.md",
        stageBeforeIteration: false,
        execFn: mockExec,
      });

      assert.strictEqual(state.status, "A  newfile.go\n");
      assert.ok(state.diff.includes("new file mode 100644"));
      assert.deepStrictEqual(state.changedFiles, ["newfile.go"]);
      assert.ok(executedCommands.some((c) => c.includes("--cached")));
    },
  );

  await t.test(
    "falls back to empty tree diff with --cached when HEAD is unborn and simple diff fails",
    async () => {
      const executedCommands = [];
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        executedCommands.push(fullCmd);
        if (fullCmd.includes("status --porcelain")) {
          return { stdout: "A  newfile.go\n" };
        }
        if (fullCmd.includes("HEAD")) {
          throw new Error('fatal: ambiguous argument "HEAD"');
        }
        if (
          fullCmd.includes("--cached") &&
          !fullCmd.includes(GIT_EMPTY_TREE_HASH)
        ) {
          throw new Error("fatal: cached diff error");
        }
        if (fullCmd.includes(GIT_EMPTY_TREE_HASH)) {
          return {
            stdout:
              "diff --git a/newfile.go b/newfile.go\nnew file mode 100644\n",
          };
        }
        return { stdout: "" };
      };

      const state = await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: false,
        execFn: mockExec,
      });

      assert.ok(state.diff.includes("new file mode 100644"));
      assert.ok(
        executedCommands.some(
          (c) => c.includes("--cached") && c.includes(GIT_EMPTY_TREE_HASH),
        ),
      );
    },
  );

  await t.test(
    "omits content of files exceeding size limit to protect context window",
    async () => {
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        if (fullCmd.includes("--name-only")) {
          return { stdout: "package.json\n" };
        }
        return { stdout: "" };
      };

      const state = await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: false,
        maxFileSizeBytes: 50,
        execFn: mockExec,
      });

      assert.deepStrictEqual(state.changedFiles, ["package.json"]);
      assert.ok(
        state.modifiedFilesContext.includes("[File content omitted: size ("),
      );
      assert.ok(
        state.modifiedFilesContext.includes("exceeds limit of 50 bytes"),
      );
      assert.ok(
        state.modifiedFilesContext.includes(
          "Inspect 'package.json' with read_file",
        ),
      );
    },
  );

  await t.test(
    "enforces maxTotalContextBytes across multiple modified files",
    async () => {
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        if (fullCmd.includes("--name-only")) {
          return {
            stdout:
              ".gemini/orchestrator/lib/utils.js\n.gemini/orchestrator/lib/ralph-loop.js\n",
          };
        }
        return { stdout: "" };
      };

      // Set maxTotalContextBytes small enough that the second file gets omitted
      const state = await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: false,
        maxFileSizeBytes: 100 * 1024,
        maxTotalContextBytes: 15 * 1024, // 15KB
        execFn: mockExec,
      });

      assert.deepStrictEqual(state.changedFiles, [
        ".gemini/orchestrator/lib/utils.js",
        ".gemini/orchestrator/lib/ralph-loop.js",
      ]);
      assert.ok(
        state.modifiedFilesContext.includes(
          '<file path=".gemini/orchestrator/lib/utils.js">',
        ),
      );
      assert.ok(
        state.modifiedFilesContext.includes(
          "[File content omitted: total modified context limit (15360 bytes) reached",
        ) ||
          state.modifiedFilesContext.includes(
            "[File content omitted: total modified context limit (15360 bytes) exceeded",
          ),
      );
      assert.ok(
        state.modifiedFilesContext.includes(
          "Inspect '.gemini/orchestrator/lib/ralph-loop.js' with read_file",
        ),
      );
    },
  );

  await t.test(
    "proactively excludes git pack and index binary files (.pack, .idx)",
    async () => {
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        if (fullCmd.includes("--name-only")) {
          return { stdout: "objects/pack/test.pack\nobjects/pack/test.idx\n" };
        }
        return { stdout: "" };
      };

      // Use a custom repoRoot where dummy test.pack and test.idx files exist or mock fs via captureGitState
      const state = await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: false,
        execFn: mockExec,
      });

      assert.deepStrictEqual(state.changedFiles, [
        "objects/pack/test.pack",
        "objects/pack/test.idx",
      ]);
    },
  );

  await t.test(
    "caps omitted file placeholders at 50 and outputs summary for remaining files",
    async () => {
      // Generate 75 files that point to an existing file so stat/exists succeed
      const fakeFiles = Array.from({ length: 75 }, () => "package.json");
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        if (fullCmd.includes("--name-only")) {
          return { stdout: fakeFiles.join("\n") + "\n" };
        }
        return { stdout: "" };
      };

      const state = await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: false,
        maxFileSizeBytes: 100 * 1024,
        maxTotalContextBytes: 100, // Very low limit so files get omitted
        execFn: mockExec,
      });

      assert.ok(
        state.modifiedFilesContext.includes(
          "[File content omitted: total modified context limit (100 bytes) exceeded",
        ),
      );
      assert.ok(
        state.modifiedFilesContext.includes(
          "additional modified file(s) omitted due to total context limit",
        ),
      );
    },
  );

  await t.test(
    "caps omitted file placeholders globally when many individual files are oversized",
    async () => {
      // Generate 75 files where each file size exceeds maxFileSizeBytes
      const fakeFiles = Array.from({ length: 75 }, () => "package.json");
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        if (fullCmd.includes("--name-only")) {
          return { stdout: fakeFiles.join("\n") + "\n" };
        }
        return { stdout: "" };
      };

      const state = await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: false,
        maxFileSizeBytes: 10, // Exceeds limit on individual file basis
        maxTotalContextBytes: 10 * 1024 * 1024, // High total limit
        execFn: mockExec,
      });

      assert.ok(
        state.modifiedFilesContext.includes("[File content omitted: size ("),
      );
      assert.ok(
        state.modifiedFilesContext.includes(
          "additional omitted file placeholder(s) omitted due to placeholder limit (50)",
        ),
      );
    },
  );

  await t.test(
    "skips paths escaping repository context (path traversal hardening)",
    async () => {
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        if (fullCmd.includes("--name-only")) {
          return { stdout: "../escaping.go\n/etc/passwd\npackage.json\n" };
        }
        return { stdout: "" };
      };

      const state = await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: false,
        execFn: mockExec,
      });

      assert.ok(
        !state.modifiedFilesContext.includes('<file path="../escaping.go">'),
      );
      assert.ok(
        !state.modifiedFilesContext.includes('<file path="/etc/passwd">'),
      );
      assert.ok(
        state.modifiedFilesContext.includes('<file path="package.json">'),
      );
    },
  );

  await t.test(
    "logs debug notice when falling back to empty tree diff",
    async () => {
      const originalDebug = console.debug;
      const debugLogs = [];
      console.debug = (...args) => debugLogs.push(args.join(" "));

      try {
        const mockExec = async (cmd, args) => {
          const fullCmd = `${cmd} ${args.join(" ")}`;
          if (fullCmd.includes("status --porcelain")) {
            return { stdout: "A  newfile.go\n" };
          }
          if (fullCmd.includes("HEAD")) {
            throw new Error('fatal: ambiguous argument "HEAD"');
          }
          if (
            fullCmd.includes("--cached") &&
            !fullCmd.includes(GIT_EMPTY_TREE_HASH)
          ) {
            throw new Error("fatal: cached diff error");
          }
          if (fullCmd.includes(GIT_EMPTY_TREE_HASH)) {
            return {
              stdout:
                "diff --git a/newfile.go b/newfile.go\nnew file mode 100644\n",
            };
          }
          return { stdout: "" };
        };

        await captureGitState({
          repoRoot: process.cwd(),
          stageBeforeIteration: false,
          execFn: mockExec,
        });

        assert.ok(
          debugLogs.some(
            (log) =>
              log.includes("Comparing staged index against empty tree hash") &&
              log.includes(GIT_EMPTY_TREE_HASH),
          ),
        );
      } finally {
        console.debug = originalDebug;
      }
    },
  );

  await t.test(
    "spills over massive diff to .gemini/tmp/staged-review.diff and embeds preview",
    async () => {
      const massiveDiff =
        "diff --git a/big.go b/big.go\n" + "+ line of code\n".repeat(5000); // > 50KB
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        if (fullCmd.includes("status --porcelain")) {
          return { stdout: "M big.go\n" };
        }
        if (fullCmd.includes("--name-only")) {
          return { stdout: "big.go\n" };
        }
        return { stdout: massiveDiff };
      };

      const state = await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: false,
        execFn: mockExec,
      });

      assert.ok(state.diff.includes("Git diff truncated"));
      assert.ok(state.diff.includes(".gemini/tmp/staged-review.diff"));
      assert.strictEqual(state.rawDiff, massiveDiff);
    },
  );

  await t.test("throws an error when stageBeforeIteration fails", async () => {
    const mockExec = async (cmd, args) => {
      const fullCmd = `${cmd} ${args.join(" ")}`;
      if (fullCmd.includes("stage-files.js")) {
        throw new Error("stage script execution failed: index locked");
      }
      return { stdout: "" };
    };

    await assert.rejects(async () => {
      await captureGitState({
        repoRoot: process.cwd(),
        stageBeforeIteration: true,
        execFn: mockExec,
      });
    }, /stage script execution failed/);
  });

  await t.test(
    "throws an error when all git diff strategies fail",
    async () => {
      const mockExec = async (cmd, args) => {
        const fullCmd = `${cmd} ${args.join(" ")}`;
        if (fullCmd.includes("status --porcelain")) {
          return { stdout: "M file.go\n" };
        }
        if (fullCmd.includes("diff")) {
          throw new Error("fatal: repository corrupted");
        }
        return { stdout: "" };
      };

      await assert.rejects(async () => {
        await captureGitState({
          repoRoot: process.cwd(),
          stageBeforeIteration: false,
          execFn: mockExec,
        });
      }, /All git diff strategies failed to capture repository diff/);
    },
  );
});

test("formatValidationSpillover spillover behavior", async (t) => {
  await t.test(
    "spills over large stdout and stderr to disk artifacts and provides previews",
    async () => {
      const massiveStdout = "stdout log line\n".repeat(500); // > 7KB
      const massiveStderr = "stderr log line\n".repeat(500); // > 7KB

      const validationResult = {
        success: false,
        stdout: massiveStdout,
        stderr: massiveStderr,
      };

      const formatted = await formatValidationSpillover(
        validationResult,
        process.cwd(),
        500,
      );

      assert.ok(formatted.stdout.includes("Output truncated"));
      assert.ok(formatted.stdout.includes(".gemini/tmp/validation-stdout.log"));
      assert.ok(formatted.stderr.includes("Output truncated"));
      assert.ok(formatted.stderr.includes(".gemini/tmp/validation-stderr.log"));
      assert.strictEqual(formatted.rawStdout, massiveStdout);
      assert.strictEqual(formatted.rawStderr, massiveStderr);
    },
  );

  await t.test(
    "passes through small outputs without spillover or modification",
    async () => {
      const validationResult = {
        success: false,
        stdout: "small stdout",
        stderr: "small stderr",
      };

      const formatted = await formatValidationSpillover(
        validationResult,
        process.cwd(),
        3000,
      );
      assert.strictEqual(formatted.stdout, "small stdout");
      assert.strictEqual(formatted.stderr, "small stderr");
    },
  );

  await t.test("returns successful validationResult untouched", async () => {
    const validationResult = { success: true };
    const formatted = await formatValidationSpillover(
      validationResult,
      process.cwd(),
    );
    assert.strictEqual(formatted, validationResult);
  });
});

test("runRalphLoop core lifecycle and control flow", async (t) => {
  await t.test(
    "exits immediately when requireInitialDiff is true and diff is empty",
    async () => {
      let sessionCalled = false;
      const mockCapture = async () => ({
        diff: "",
        status: "",
        diffHash: "empty-hash",
        plan: "",
        changedFiles: [],
        modifiedFilesContext: "",
      });

      const result = await runRalphLoop({
        requireInitialDiff: true,
        captureState: mockCapture,
        sessionRunner: async () => {
          sessionCalled = true;
          return "output";
        },
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.iterations, 0);
      assert.strictEqual(sessionCalled, false);
      assert.strictEqual(result.reason, "No staged changes found to review.");
    },
  );

  await t.test(
    "does not exit early when requireInitialDiff is true but validator fails on empty diff",
    async () => {
      let capturedLastValidation = null;
      let validatorCalls = 0;
      const mockCapture = async () => ({
        diff: "",
        status: "",
        diffHash: "empty-hash",
        plan: "",
        changedFiles: [],
        modifiedFilesContext: "",
      });

      const result = await runRalphLoop({
        requireInitialDiff: true,
        maxIterations: 1,
        captureState: mockCapture,
        validator: async () => {
          validatorCalls++;
          return { success: false, stderr: "Tests failed" };
        },
        iterationRunner: async ({ lastValidation }) => {
          capturedLastValidation = lastValidation;
          return { shouldImplement: false, forceComplete: true };
        },
        logger: () => {},
      });

      assert.strictEqual(result.iterations, 1);
      assert.ok(
        validatorCalls >= 1,
        "Validator must be executed on empty diff",
      );
      assert.ok(
        capturedLastValidation,
        "lastValidation must be passed to iterationRunner on iteration 1",
      );
      assert.strictEqual(capturedLastValidation.success, false);
      assert.strictEqual(capturedLastValidation.stderr, "Tests failed");
    },
  );

  await t.test(
    "converges in 1 iteration when agent makes no changes and validation passes",
    async () => {
      let sessionCount = 0;
      let validatorCount = 0;

      const mockState = {
        diff: "initial diff",
        status: "M file.go",
        diffHash: "hash-initial",
        plan: "plan",
        changedFiles: ["file.go"],
        modifiedFilesContext: "",
      };

      const result = await runRalphLoop({
        captureState: async () => ({ ...mockState }),
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async () => {
          sessionCount++;
          return "Everything looks good. No changes needed.";
        },
        validator: async () => {
          validatorCount++;
          return { success: true };
        },
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.iterations, 1);
      assert.strictEqual(sessionCount, 1);
      assert.strictEqual(validatorCount, 1);
      assert.ok(result.reason.includes("No further changes suggested"));
    },
  );

  await t.test(
    "does not converge on stateChanged: false when iterationRunner returns shouldImplement: true",
    async () => {
      let runnerCalls = 0;
      const mockState = {
        diff: "initial diff",
        status: "M file.go",
        diffHash: "hash-initial",
        plan: "plan",
        changedFiles: ["file.go"],
        modifiedFilesContext: "",
      };

      const result = await runRalphLoop({
        captureState: async () => ({ ...mockState }),
        detectChange: () => ({ hasChanged: false }),
        maxIterations: 3,
        noProgressLimit: 2,
        iterationRunner: async () => {
          runnerCalls++;
          return {
            agentOutput: "No-op implementation",
            shouldImplement: true,
          };
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      // Should trip circuit breaker after noProgressLimit (2) calls rather than converging on iteration 1
      assert.strictEqual(runnerCalls, 2);
      assert.strictEqual(result.completed, false);
      assert.ok(
        result.reason.toLowerCase().includes("circuit breaker") ||
          result.reason.toLowerCase().includes("no progress"),
      );
    },
  );

  await t.test(
    "iterates when agent makes improvements in iteration 1 and converges in iteration 2",
    async () => {
      let currentIteration = 0;
      const states = [
        { diffHash: "hash-0", status: "M file.go", diff: "diff 0", plan: "p" }, // baseline
        { diffHash: "hash-1", status: "M file.go", diff: "diff 1", plan: "p" }, // after iter 1 (changed)
        { diffHash: "hash-1", status: "M file.go", diff: "diff 1", plan: "p" }, // iter 2 before
        { diffHash: "hash-1", status: "M file.go", diff: "diff 1", plan: "p" }, // after iter 2 (no change)
      ];
      let stateIdx = 0;

      const sessionPrompts = [];

      const result = await runRalphLoop({
        maxIterations: 5,
        captureState: async () => {
          const state = states[Math.min(stateIdx, states.length - 1)];
          stateIdx++;
          return state;
        },
        detectChange: (before, after) => {
          return { hasChanged: before.diffHash !== after.diffHash };
        },
        sessionRunner: async ({ initialPrompt }) => {
          currentIteration++;
          sessionPrompts.push(initialPrompt);
          return currentIteration === 1
            ? "Applied refactoring"
            : "No further changes";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.iterations, 2);
      assert.strictEqual(sessionPrompts.length, 2);
      assert.ok(result.reason.includes("No further changes suggested"));
    },
  );

  await t.test(
    "injects negative feedback from failed validation and heals in subsequent iteration",
    async () => {
      let currentIteration = 0;
      const sessionPrompts = [];

      // Iteration 1: changes state, validator fails
      // Iteration 2: fixes failure, validator succeeds
      // Iteration 3: no changes, validator succeeds -> converges!
      const changeSequence = [
        true, // iter 1 made changes
        true, // iter 2 made changes (fixed the issue)
        false, // iter 3 made no changes
      ];

      const validationSequence = [
        {
          success: false,
          stdout: "FAIL",
          stderr: "syntax error: unexpected token on line 42",
        },
        { success: true, stdout: "PASS", stderr: "" },
        { success: true, stdout: "PASS", stderr: "" },
      ];

      const result = await runRalphLoop({
        maxIterations: 5,
        captureState: async () => ({
          diff: "diff content",
          status: "M file.go",
          diffHash: `hash-${currentIteration}`,
          plan: "plan",
          changedFiles: ["file.go"],
        }),
        detectChange: () => ({
          hasChanged: changeSequence[currentIteration - 1] ?? false,
        }),
        sessionRunner: async ({ initialPrompt }) => {
          currentIteration++;
          sessionPrompts.push(initialPrompt);
          return `Turn output for iter ${currentIteration}`;
        },
        validator: async ({ iteration }) => {
          return validationSequence[iteration - 1] || { success: true };
        },
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.iterations, 3);

      // Verify negative feedback was injected into iteration 2 prompt
      assert.ok(
        sessionPrompts[1].includes(
          "⚠️ CRITICAL VALIDATION FEEDBACK (Negative Backpressure):",
        ),
      );
      assert.ok(
        sessionPrompts[1].includes("syntax error: unexpected token on line 42"),
      );

      // Verify iteration 3 prompt did not have validation failure
      assert.ok(!sessionPrompts[2].includes("<validation_failure>"));
      assert.ok(sessionPrompts[2].includes("<validation_status>"));
    },
  );

  await t.test(
    "halts via no-progress circuit breaker when validation fails and agent makes no changes",
    async () => {
      let sessionCount = 0;

      const result = await runRalphLoop({
        maxIterations: 10,
        noProgressLimit: 3,
        captureState: async () => ({
          diff: "failing diff",
          status: "M file.go",
          diffHash: "constant-hash",
          plan: "plan",
        }),
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async () => {
          sessionCount++;
          return "I could not figure out how to fix it.";
        },
        validator: async () => ({
          success: false,
          stdout: "",
          stderr: "Tests failed with exit code 1",
        }),
        logger: () => {},
      });

      assert.strictEqual(result.success, false);
      assert.strictEqual(result.completed, false);
      assert.strictEqual(result.iterations, 3);
      assert.strictEqual(sessionCount, 3);
      assert.ok(
        result.reason.includes("No progress made for 3 consecutive iterations"),
      );
    },
  );

  await t.test(
    "detects oscillating failing diff states and trips circuit breaker",
    async () => {
      let sessionCount = 0;
      // Agent alternates between state A and state B, but both fail validation
      const diffHashes = ["hash-A", "hash-B", "hash-A", "hash-B"];

      const result = await runRalphLoop({
        maxIterations: 10,
        noProgressLimit: 3,
        captureState: async () => ({
          diff: "failing diff",
          status: "M file.go",
          diffHash: diffHashes[sessionCount % diffHashes.length],
          plan: "plan",
        }),
        detectChange: () => ({ hasChanged: true }),
        sessionRunner: async () => {
          sessionCount++;
          return `Attempt ${sessionCount}`;
        },
        validator: async () => ({
          success: false,
          stderr: "Still failing",
        }),
        logger: () => {},
      });

      assert.strictEqual(result.success, false);
      assert.strictEqual(result.completed, false);
      assert.ok(
        result.reason.includes("No progress made for 3 consecutive iterations"),
      );
    },
  );

  await t.test(
    "clears failing diff hash counts when validation succeeds",
    async () => {
      let sessionCount = 0;
      // Iteration 1: fails with hash-A
      // Iteration 2: succeeds with hash-B (clears failingDiffHashCounts!)
      // Iteration 3: fails with hash-A (first occurrence after reset, so no immediate oscillation trip)
      // Iteration 4: fails with hash-A (second occurrence after reset, now trips)
      const states = [
        { diffHash: "hash-init", status: "", plan: "" },
        { diffHash: "hash-A", status: "", plan: "" },
        { diffHash: "hash-B", status: "", plan: "" },
        { diffHash: "hash-A", status: "", plan: "" },
        { diffHash: "hash-A", status: "", plan: "" },
      ];
      const validationResults = [
        { success: false, stderr: "Fail 1" },
        { success: true },
        { success: false, stderr: "Fail 2" },
        { success: false, stderr: "Fail 3" },
      ];

      let stateIdx = 0;
      const result = await runRalphLoop({
        maxIterations: 10,
        noProgressLimit: 2,
        captureState: async () => {
          const s = states[Math.min(stateIdx, states.length - 1)];
          stateIdx++;
          return s;
        },
        detectChange: (before, after) => ({
          hasChanged: before.diffHash !== after.diffHash,
        }),
        sessionRunner: async () => {
          sessionCount++;
          return `Iteration ${sessionCount}`;
        },
        validator: async ({ iteration }) =>
          validationResults[iteration - 1] || { success: false },
        logger: () => {},
      });

      assert.strictEqual(result.success, false);
      // Iteration 1 (fail A), Iteration 2 (pass B, clears count), Iteration 3 (fail A, count=1, no trip), Iteration 4 (fail A, count=2 -> trips noProgressLimit=2)
      assert.strictEqual(result.iterations, 4);
      assert.strictEqual(sessionCount, 4);
    },
  );

  await t.test(
    "enforces maxIterations ceiling when changes continue without full resolution",
    async () => {
      let sessionCount = 0;

      const result = await runRalphLoop({
        maxIterations: 4,
        captureState: async () => ({
          diff: "diff",
          status: "M file.go",
          diffHash: `hash-${sessionCount}`,
          plan: "plan",
        }),
        detectChange: () => ({ hasChanged: true }),
        sessionRunner: async () => {
          sessionCount++;
          return "Still working on it...";
        },
        validator: async () => ({
          success: false,
          stderr: "Still failing",
        }),
        logger: () => {},
      });

      assert.strictEqual(result.success, false);
      assert.strictEqual(result.completed, false);
      assert.strictEqual(result.iterations, 4);
      assert.strictEqual(sessionCount, 4);
      assert.ok(result.reason.includes("Reached maximum iterations (4)"));
    },
  );

  await t.test(
    "terminates immediately when agent outputs completion promise and validation passes",
    async () => {
      const result = await runRalphLoop({
        maxIterations: 10,
        captureState: async () => ({
          diff: "diff",
          status: "M file.go",
          diffHash: "hash-1",
          plan: "plan",
        }),
        detectChange: () => ({ hasChanged: true }),
        sessionRunner: async () => {
          return "All tasks done. <promise>COMPLETE</promise>";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.iterations, 1);
      assert.ok(result.reason.includes("completion promise"));
    },
  );

  await t.test(
    "executes onIterationStart and onIterationEnd lifecycle hooks",
    async () => {
      const starts = [];
      const ends = [];

      const result = await runRalphLoop({
        maxIterations: 2,
        captureState: async () => ({
          diff: "diff",
          status: "",
          diffHash: "h",
          plan: "",
        }),
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async () => "done",
        validator: async () => ({ success: true }),
        onIterationStart: async (ctx) => {
          starts.push(ctx.iteration);
        },
        onIterationEnd: async (ctx) => {
          ends.push(ctx.iteration);
        },
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.deepStrictEqual(starts, [1]);
      assert.deepStrictEqual(ends, [1]);
    },
  );

  await t.test("supports custom checkCompletion predicate", async () => {
    const result = await runRalphLoop({
      maxIterations: 5,
      captureState: async () => ({
        diff: "d",
        status: "",
        diffHash: "h",
        plan: "",
      }),
      detectChange: () => ({ hasChanged: true }),
      sessionRunner: async () => "custom response",
      validator: async () => ({ success: true }),
      checkCompletion: async ({ iteration }) => {
        if (iteration === 2) {
          return { isComplete: true, reason: "Custom early stop at 2" };
        }
        return { isComplete: false };
      },
      logger: () => {},
    });

    assert.strictEqual(result.success, true);
    assert.strictEqual(result.completed, true);
    assert.strictEqual(result.iterations, 2);
    assert.strictEqual(result.reason, "Custom early stop at 2");
  });

  await t.test(
    "handles agent session throwing an exception gracefully",
    async () => {
      let called = false;
      const result = await runRalphLoop({
        maxIterations: 2,
        noProgressLimit: 2,
        captureState: async () => ({
          diff: "d",
          status: "",
          diffHash: "h",
          plan: "",
        }),
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async () => {
          called = true;
          throw new Error("SDK network timeout");
        },
        validator: async () => ({ success: false, stderr: "Tests failing" }),
        logger: () => {},
      });

      assert.strictEqual(called, true);
      assert.strictEqual(result.success, false);
    },
  );

  await t.test(
    "does not falsely converge when agent session throws an error even if validator passes",
    async () => {
      let callCount = 0;
      const result = await runRalphLoop({
        maxIterations: 3,
        noProgressLimit: 2,
        captureState: async () => ({
          diff: "d",
          status: "",
          diffHash: "h",
          plan: "",
        }),
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async () => {
          callCount++;
          throw new Error("Session crashed unexpectedly");
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(callCount, 2);
      assert.strictEqual(result.success, false);
      assert.strictEqual(result.completed, false);
      assert.ok(result.reason.includes("Agent session failed with error"));
      assert.ok(result.reason.includes("Session crashed unexpectedly"));
    },
  );

  await t.test(
    "prevents unhandled rejection and handles errors gracefully when iterationRunner throws",
    async () => {
      let runnerCalled = false;
      const result = await runRalphLoop({
        maxIterations: 2,
        noProgressLimit: 2,
        captureState: async () => ({
          diff: "d",
          status: "",
          diffHash: "h",
          plan: "",
        }),
        detectChange: () => ({ hasChanged: false }),
        iterationRunner: async () => {
          runnerCalled = true;
          throw new Error("Multi-agent pipeline timeout");
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(runnerCalled, true);
      assert.strictEqual(result.success, false);
      assert.strictEqual(result.completed, false);
      assert.ok(result.reason.includes("Agent session failed with error"));
      assert.ok(result.reason.includes("Multi-agent pipeline timeout"));
    },
  );

  await t.test(
    "prevents completion and retries when iterationRunner requests forceComplete but validator fails",
    async () => {
      let iterationCount = 0;
      let lastValidationPassedToRunner = null;

      const result = await runRalphLoop({
        maxIterations: 3,
        noProgressLimit: 3,
        captureState: async () => ({
          diff: "d",
          status: "",
          diffHash: `h-${iterationCount}`,
          plan: "",
        }),
        detectChange: () => ({ hasChanged: true }),
        iterationRunner: async ({ iteration, lastValidation }) => {
          iterationCount = iteration;
          lastValidationPassedToRunner = lastValidation;
          if (iteration === 1) {
            // First iteration attempts to complete early
            return {
              forceComplete: true,
              completionReason: "Only minor issues seen",
              agentOutput: "Looks good",
            };
          }
          // Second iteration completes when validator succeeds
          return {
            forceComplete: true,
            completionReason: "Fixed and validated",
            agentOutput: "Resolved",
          };
        },
        validator: async ({ iteration }) => {
          if (iteration === 1) {
            return { success: false, stderr: "linter error: missing comment" };
          }
          return { success: true };
        },
        logger: () => {},
      });

      // Iteration 1 had forceComplete, but validator failed, so it proceeded to iteration 2
      assert.strictEqual(result.iterations, 2);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.success, true);
      assert.strictEqual(lastValidationPassedToRunner.success, false);
      assert.ok(
        lastValidationPassedToRunner.stderr.includes(
          "linter error: missing comment",
        ),
      );
    },
  );

  await t.test(
    "passes validationCommands option to buildPrompt in single-agent loops",
    async () => {
      let capturedPrompt = "";
      const result = await runRalphLoop({
        maxIterations: 1,
        captureState: async () => ({
          diff: "d",
          status: "",
          diffHash: "h",
          plan: "",
        }),
        detectChange: () => ({ hasChanged: false }),
        validationCommands: ["npm run lint:custom", "npm run test:unit"],
        sessionRunner: async ({ initialPrompt }) => {
          capturedPrompt = initialPrompt;
          return "Code looks good.";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.ok(
        capturedPrompt.includes(
          "'npm run lint:custom' and 'npm run test:unit'",
        ),
      );
    },
  );
});

test("multi-agent prompts and suggestion parsing", async (t) => {
  await t.test(
    "buildReviewerPrompt strives for minimum 5 suggestions with escape hatch and rating scale",
    () => {
      const prompt = buildReviewerPrompt({
        iteration: 1,
        maxIterations: 10,
        stateBefore: {
          diff: "diff content",
          modifiedFilesContext: "<file>code</file>",
        },
        activePlan: "## Active Plan",
      });

      assert.ok(prompt.includes("empty list [] if no suggestions exist"));
      assert.ok(prompt.includes("output: []"));
      assert.ok(prompt.includes("Strive for a minimum of 5 suggestions"));
      assert.ok(prompt.includes("Escape Hatch"));
      assert.ok(prompt.includes('"critical", "medium", or "minor"'));
      assert.ok(prompt.includes("Analyze for Impact"));
      assert.ok(prompt.includes("Bugs & Logic Flaws"));
      assert.ok(prompt.includes("Minimize Pedantic Noise"));
      assert.ok(prompt.includes("Be Concrete"));
      assert.ok(
        prompt.includes(
          "cite the exact file name and line reference from the diff, explain the risk, and provide a concrete replacement snippet or fix rationale",
        ),
      );
      assert.ok(
        prompt.includes(
          "Do NOT add markdown formatting (such as asterisks or bolding, e.g. **[critical]** or [**critical**])",
        ),
      );
      assert.ok(prompt.includes("Do NOT make file modifications"));
      assert.ok(
        prompt.includes("<active_plan>\n## Active Plan\n</active_plan>"),
      );
    },
  );

  await t.test(
    "buildReviewerPrompt injects deterministic validation errors",
    () => {
      const prompt = buildReviewerPrompt({
        lastValidation: { success: false, stderr: "undefined function Foo" },
      });
      assert.ok(prompt.includes("⚠️ PREVIOUS VALIDATION FAILURE"));
      assert.ok(prompt.includes("undefined function Foo"));
    },
  );

  await t.test(
    "buildTallyPrompt requests strict JSON with critical, medium, minor counts",
    () => {
      const prompt = buildTallyPrompt("1. [critical] Fix A\n2. [medium] Fix B");
      assert.ok(prompt.includes('"critical": <integer>'));
      assert.ok(prompt.includes('"medium": <integer>'));
      assert.ok(prompt.includes('"minor": <integer>'));
      assert.ok(prompt.includes("1. [critical] Fix A"));
    },
  );

  await t.test(
    "buildImplementationPrompt instructs implementing all suggestions when critical/medium present and includes update_topic",
    () => {
      const prompt = buildImplementationPrompt({
        suggestionsText: "1. [critical] Fix leak\n2. [minor] Rename var",
        tally: { critical: 1, medium: 0, minor: 1 },
        activePlan: "## Plan",
        stateBefore: { diff: "diff" },
        lastValidation: {
          success: false,
          stdout: "FAIL: TestSomething",
          stderr: "panic: nil pointer dereference",
        },
      });
      assert.ok(prompt.includes("1 critical, 0 medium, 1 minor"));
      assert.ok(prompt.includes("implement ALL suggestions"));
      assert.ok(prompt.includes("1. [critical] Fix leak"));
      assert.ok(
        prompt.includes(
          'call update_topic with a descriptive "title" and "summary"',
        ),
      );
      assert.ok(
        prompt.includes(
          "You have access to GrepSearch (and grep_search) to search across workspace file contents.",
        ),
      );
      assert.ok(
        prompt.includes(
          "Do NOT enter plan mode (do not call enter_plan_mode).",
        ),
      );
      assert.ok(
        prompt.includes("Do NOT invoke subagents (do not call invoke_agent)."),
      );
      assert.ok(
        prompt.includes(
          "The active plan from 'plans/current.md' is already embedded below in <active_plan>. Do NOT waste turns reading it from disk.",
        ),
      );
      assert.ok(
        prompt.includes(
          "Formulate a clear implementation plan for addressing the review suggestions",
        ),
      );
      assert.ok(
        prompt.includes(
          "CRITICAL PLAN INTEGRITY RULE: Do NOT overwrite or wipe 'plans/current.md'",
        ),
      );
      assert.ok(!prompt.includes('"strategic_intent"'));
      assert.ok(
        prompt.includes(
          "⚠️ PREVIOUS VALIDATION FAILURE (Deterministic Backpressure):",
        ),
      );
      assert.ok(prompt.includes("STDOUT:\nFAIL: TestSomething"));
      assert.ok(
        prompt.includes("STDERR / ERRORS:\npanic: nil pointer dereference"),
      );
    },
  );

  await t.test(
    "parseSuggestionCounts extracts final count when conversational padding or preamble exists",
    () => {
      const text =
        "I found 0 critical, 0 medium, and 0 minor issues initially, but wait, there are actually 2 critical, 1 medium, and 3 minor suggestions!";
      const tally = parseSuggestionCounts(text);
      assert.strictEqual(tally.critical, 2);
      assert.strictEqual(tally.medium, 1);
      assert.strictEqual(tally.minor, 3);
      assert.strictEqual(tally.total, 6);
    },
  );

  await t.test(
    "parseSuggestionCounts parses valid JSON object from flash-lite output",
    () => {
      const jsonText =
        '```json\n{\n  "critical": 5,\n  "medium": 2,\n  "minor": 1\n}\n```';
      const tally = parseSuggestionCounts(jsonText);
      assert.strictEqual(tally.critical, 5);
      assert.strictEqual(tally.medium, 2);
      assert.strictEqual(tally.minor, 1);
      assert.strictEqual(tally.total, 8);
    },
  );

  await t.test(
    "parseSuggestionCounts respects explicitly parsed zero counts { critical: 0, medium: 0, minor: 0 } without word-matching fallback",
    () => {
      const jsonText =
        '```json\n{\n  "critical": 0,\n  "medium": 0,\n  "minor": 0\n}\n```';
      const tally = parseSuggestionCounts(jsonText);
      assert.strictEqual(tally.critical, 0);
      assert.strictEqual(tally.medium, 0);
      assert.strictEqual(tally.minor, 0);
      assert.strictEqual(tally.total, 0);
    },
  );

  await t.test(
    "parseSuggestionCounts rejects incomplete JSON schemas and falls through to reviewer tags",
    () => {
      const incompleteJson = '```json\n{\n  "minor": 3\n}\n```';
      const reviewerText = `
1. [critical] Fix security bug
2. [medium] Fix performance bottleneck
3. [minor] Fix typo
`;
      const tally = parseSuggestionCounts(incompleteJson, reviewerText);
      assert.strictEqual(tally.critical, 1);
      assert.strictEqual(tally.medium, 1);
      assert.strictEqual(tally.minor, 1);
      assert.strictEqual(tally.total, 3);
    },
  );

  await t.test("parseSuggestionCounts parses explicit zero text counts", () => {
    const text = "Identified 0 critical, 0 medium, and 0 minor suggestions.";
    const tally = parseSuggestionCounts(text);
    assert.strictEqual(tally.critical, 0);
    assert.strictEqual(tally.medium, 0);
    assert.strictEqual(tally.minor, 0);
    assert.strictEqual(tally.total, 0);
  });

  await t.test(
    "parseSuggestionCounts returns zero counts for [] and code block arrays",
    () => {
      assert.deepStrictEqual(parseSuggestionCounts("[]"), {
        critical: 0,
        medium: 0,
        minor: 0,
        total: 0,
      });
      assert.deepStrictEqual(parseSuggestionCounts("```json\n[]\n```"), {
        critical: 0,
        medium: 0,
        minor: 0,
        total: 0,
      });
    },
  );

  await t.test(
    "parseSuggestionCounts parses JSON object with title-cased or upper-cased keys",
    () => {
      const jsonText =
        '```json\n{\n  "Critical": 3,\n  "MEDIUM": 2,\n  "Minor": 4\n}\n```';
      const tally = parseSuggestionCounts(jsonText);
      assert.strictEqual(tally.critical, 3);
      assert.strictEqual(tally.medium, 2);
      assert.strictEqual(tally.minor, 4);
      assert.strictEqual(tally.total, 9);
    },
  );

  await t.test(
    "parseSuggestionCounts parses text statements with counts",
    () => {
      const text =
        "Here is the summary: 4 critical, 3 medium, and 2 minor suggestions identified.";
      const tally = parseSuggestionCounts(text);
      assert.strictEqual(tally.critical, 4);
      assert.strictEqual(tally.medium, 3);
      assert.strictEqual(tally.minor, 2);
      assert.strictEqual(tally.total, 9);
    },
  );

  await t.test(
    "parseSuggestionCounts parses pluralized text category counts",
    () => {
      const text = "Identified 5 criticals, 2 mediums, and 1 minor in review.";
      const tally = parseSuggestionCounts(text);
      assert.strictEqual(tally.critical, 5);
      assert.strictEqual(tally.medium, 2);
      assert.strictEqual(tally.minor, 1);
      assert.strictEqual(tally.total, 8);
    },
  );

  await t.test(
    "parseSuggestionCounts falls back to scanning bracketed severity tags from raw reviewer text",
    () => {
      const reviewerText = `
1. [critical] Fix SQL injection
2. [critical] Close database connection
3. [medium] Add input validation
4. [minor] Fix typo in log
5. [minor] Rename variable x
`;
      const tally = parseSuggestionCounts("", reviewerText);
      assert.strictEqual(tally.critical, 2);
      assert.strictEqual(tally.medium, 1);
      assert.strictEqual(tally.minor, 2);
      assert.strictEqual(tally.total, 5);
    },
  );

  await t.test(
    "parseSuggestionCounts handles markdown bold formatting around severity tags",
    () => {
      const reviewerText = `
1. **[critical]** Fix SQL injection
2. [**critical**] Close database connection
3. **[medium]** Add input validation
4. [minor] Fix typo in log
5. *[minor]* Rename variable x
`;
      const tally = parseSuggestionCounts("", reviewerText);
      assert.strictEqual(tally.critical, 2);
      assert.strictEqual(tally.medium, 1);
      assert.strictEqual(tally.minor, 2);
      assert.strictEqual(tally.total, 5);
    },
  );

  await t.test(
    "parseSuggestionCounts does not mistakenly tally conversational sentences in list items as critical tags",
    () => {
      const reviewerText = `
1. [minor] Documentation update: We should fix the critical bug description in README.
2. [minor] Variable rename: Change temp to result.
3. [minor] Formatting: Clean up trailing spaces.
4. [minor] Add comments explaining the critical path algorithm.
5. [minor] Reorder import statements.
`;
      const tally = parseSuggestionCounts("", reviewerText);
      // Because all 5 items are tagged [minor], "critical" inside the descriptions must NOT be tallied as critical
      assert.strictEqual(tally.critical, 0);
      assert.strictEqual(tally.medium, 0);
      assert.strictEqual(tally.minor, 5);
      assert.strictEqual(tally.total, 5);
    },
  );

  await t.test(
    "parseSuggestionCounts handles parenthesized severity tags like (critical)",
    () => {
      const reviewerText = `
1. (critical) Fix null pointer exception in handler
2. (medium) Optimize loop performance
3. (minor) Fix typo in error message
4. (minor) Remove unused parameter
5. (minor) Add missing docstring
`;
      const tally = parseSuggestionCounts("", reviewerText);
      assert.strictEqual(tally.critical, 1);
      assert.strictEqual(tally.medium, 1);
      assert.strictEqual(tally.minor, 3);
      assert.strictEqual(tally.total, 5);
    },
  );

  await t.test(
    "parseSuggestionCounts returns zero counts for conversational text without structural tags or numbers",
    () => {
      const conversationalText =
        "This critical review discusses medium complexity topics and minor details in general conversation without any list or structured findings.";
      const tally = parseSuggestionCounts(conversationalText);
      assert.strictEqual(tally.critical, 0);
      assert.strictEqual(tally.medium, 0);
      assert.strictEqual(tally.minor, 0);
      assert.strictEqual(tally.total, 0);
    },
  );

  await t.test(
    "buildImplementationPrompt and buildReviewPrompt fall back gracefully when validationCommands is empty",
    () => {
      const implPrompt = buildImplementationPrompt({
        suggestionsText: "1. [critical] Fix A",
        tally: { critical: 1, medium: 0, minor: 0 },
        validationCommands: [],
      });
      assert.ok(
        implPrompt.includes("Ensure all tests and linters pass cleanly."),
      );
      assert.ok(!implPrompt.includes("(e.g.,"));

      const reviewPrompt = buildReviewPrompt({
        iteration: 1,
        maxIterations: 5,
        validationCommands: [],
      });
      assert.ok(
        reviewPrompt.includes("Ensure all tests and linters pass cleanly."),
      );
      assert.ok(!reviewPrompt.includes("(e.g.,"));
    },
  );

  await t.test(
    "formatValidationCommands formats single, pair, and 3+ commands cleanly",
    () => {
      assert.strictEqual(
        formatValidationCommands(["go test ./..."]),
        "'go test ./...'",
      );
      assert.strictEqual(
        formatValidationCommands(["go test ./...", "npm test"]),
        "'go test ./...' and 'npm test'",
      );
      assert.strictEqual(
        formatValidationCommands(["cmd1", "cmd2", "cmd3"]),
        "'cmd1', 'cmd2', and 'cmd3'",
      );
      assert.strictEqual(
        formatValidationCommands("npm run lint"),
        "npm run lint",
      );
      assert.strictEqual(formatValidationCommands([]), "");
      assert.strictEqual(formatValidationCommands(null), "");
    },
  );

  await t.test(
    "parseSuggestionCounts handles asterisk bullet points like * [critical]",
    () => {
      const reviewerText = `
* [critical] Fix memory leak in connection pool
* [medium] Add input validation on HTTP handler
* [minor] Clean up unused variable
* [minor] Fix typo in log comment
* [minor] Reorder import statements
`;
      const tally = parseSuggestionCounts("", reviewerText);
      assert.strictEqual(tally.critical, 1);
      assert.strictEqual(tally.medium, 1);
      assert.strictEqual(tally.minor, 3);
      assert.strictEqual(tally.total, 5);
    },
  );

  await t.test(
    "parseSuggestionCounts routes debug messages to injected custom logger",
    () => {
      const debugLogs = [];
      const customLogger = (msg) => debugLogs.push(msg);

      // Provide malformed JSON text to trigger debug logging
      const malformedJsonText = '{"critical": invalid}';
      const tally = parseSuggestionCounts(malformedJsonText, "", customLogger);

      assert.deepStrictEqual(tally, {
        critical: 0,
        medium: 0,
        minor: 0,
        total: 0,
      });
      assert.ok(debugLogs.length > 0);
      assert.ok(
        debugLogs.some((l) => l.includes("parseJSONFromText returned null")),
      );
    },
  );

  await t.test(
    "calculateScaledTurns calculates bounds and scales turns properly",
    () => {
      // Default base turns (20) and 5 suggestions -> 20 turns
      assert.strictEqual(calculateScaledTurns(20, 5), 20);
      // Lower suggestions (e.g. 1 suggestion) -> floor at baseTurns (20)
      assert.strictEqual(calculateScaledTurns(20, 1), 20);
      // 10 suggestions (2x multiplier) -> 40 turns
      assert.strictEqual(calculateScaledTurns(20, 10), 40);
      // 15 suggestions (3x multiplier) -> 60 turns (ceiling)
      assert.strictEqual(calculateScaledTurns(20, 15), 60);
      // 30 suggestions -> capped at ceiling (60)
      assert.strictEqual(calculateScaledTurns(20, 30), 60);
      // Custom maxCeiling
      assert.strictEqual(calculateScaledTurns(20, 30, 80), 80);
      // When baseTurns exceeds maxCeiling, base budget is clamped to maxCeiling
      assert.strictEqual(calculateScaledTurns(100, 1, 60), 60);
    },
  );

  await t.test(
    "tallySuggestionsWithFlashLite executes sessionRunner and parses counts",
    async () => {
      const tally = await tallySuggestionsWithFlashLite({
        suggestionsText:
          "1. [critical] Bug\n2. [medium] Perf\n3. [minor] Doc\n4. [minor] Style\n5. [minor] Lint",
        sessionRunner: async () => '{"critical": 1, "medium": 1, "minor": 3}',
        logger: () => {},
      });

      assert.strictEqual(tally.critical, 1);
      assert.strictEqual(tally.medium, 1);
      assert.strictEqual(tally.minor, 3);
      assert.strictEqual(tally.total, 5);
    },
  );
});

test("runReviewRalphLoop multi-agent pipeline execution", async (t) => {
  await t.test(
    "skips implementation and completes convergence when only minor suggestions exist",
    async () => {
      const sessionsInvoked = [];

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        captureState: mockCapture,
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async ({ requestedModel }) => {
          sessionsInvoked.push(requestedModel);
          if (requestedModel === "gemini-3.1-pro-preview") {
            return `
1. [minor] Suggestion 1: Documentation typo
2. [minor] Suggestion 2: Inline comment clarity
3. [minor] Suggestion 3: Code formatting
4. [minor] Suggestion 4: Reorder imports
5. [minor] Suggestion 5: Variable naming
`;
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            return JSON.stringify({ critical: 0, medium: 0, minor: 5 });
          }
          return "Should not be called";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.iterations, 1);
      // Verified Pro thinking reviewer and Flash-Lite were called, but Flash implementer was NOT called!
      assert.deepStrictEqual(sessionsInvoked, [
        "gemini-3.1-pro-preview",
        "gemini-3.1-flash-lite",
      ]);
      assert.ok(result.reason.includes("Only minor suggestions found"));
    },
  );

  await t.test(
    "retries when reviewer returns empty output and does not falsely converge",
    async () => {
      let reviewerInvocations = 0;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 2,
        captureState: mockCapture,
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async ({ requestedModel }) => {
          if (requestedModel === "gemini-3.1-pro-preview") {
            reviewerInvocations++;
            if (reviewerInvocations === 1) {
              return ""; // Empty response on first iteration
            }
            return `
1. [minor] Suggestion 1
2. [minor] Suggestion 2
3. [minor] Suggestion 3
4. [minor] Suggestion 4
5. [minor] Suggestion 5
`;
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            return JSON.stringify({ critical: 0, medium: 0, minor: 5 });
          }
          return "";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(
        reviewerInvocations,
        2,
        "Reviewer should have retried after returning empty output",
      );
      assert.strictEqual(result.success, true);
      assert.strictEqual(result.iterations, 2);
    },
  );

  await t.test(
    "converges immediately in 1 iteration when reviewer returns []",
    async () => {
      let reviewerInvocations = 0;
      let flashLiteInvocations = 0;
      let capturedAllowedTools = null;
      let capturedMaxTurns = null;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 5,
        captureState: mockCapture,
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async ({
          requestedModel,
          allowedTools,
          maxTurnsOverride,
        }) => {
          if (requestedModel === "gemini-3.1-pro-preview") {
            reviewerInvocations++;
            capturedAllowedTools = allowedTools;
            capturedMaxTurns = maxTurnsOverride;
            return "[]";
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            flashLiteInvocations++;
            return JSON.stringify({ critical: 0, medium: 0, minor: 0 });
          }
          return "";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(reviewerInvocations, 1);
      assert.strictEqual(
        flashLiteInvocations,
        0,
        "Flash-lite should be skipped when reviewer outputs []",
      );
      assert.deepStrictEqual(
        capturedAllowedTools,
        [],
        "Reviewer session must have allowedTools: []",
      );
      assert.strictEqual(
        capturedMaxTurns,
        1,
        "Reviewer session must use maxTurnsOverride: 1",
      );
      assert.strictEqual(result.success, true);
      assert.strictEqual(result.iterations, 1);
      assert.strictEqual(result.completed, true);
    },
  );

  await t.test(
    "implements all suggestions when critical or medium suggestions are present",
    async () => {
      const sessionsInvoked = [];
      let iteration = 0;
      let stageCalled = false;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: `hash-${iteration}`,
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 5,
        captureState: mockCapture,
        stageFn: async () => {
          stageCalled = true;
        },
        detectChange: () => ({ hasChanged: iteration === 1 }),
        sessionRunner: async ({ requestedModel, initialPrompt }) => {
          sessionsInvoked.push(requestedModel);

          if (initialPrompt.includes("senior technical reviewer")) {
            iteration++;
            if (iteration === 1) {
              return `
1. [critical] Memory leak in buffer
2. [medium] Missing validation check
3. [minor] Typo in doc
4. [minor] Formatting
5. [minor] Comment cleanup
`;
            }
            // Iteration 2: only minor suggestions remain
            return `
1. [minor] Comment cleanup
2. [minor] Typo in doc
3. [minor] Variable rename
4. [minor] Formatting
5. [minor] Indentation
`;
          }

          if (requestedModel === "gemini-3.1-flash-lite") {
            if (iteration === 1) {
              return JSON.stringify({ critical: 1, medium: 1, minor: 3 });
            }
            return JSON.stringify({ critical: 0, medium: 0, minor: 5 });
          }

          if (
            initialPrompt.includes("remaining suggestion(s) to implement") ||
            initialPrompt.includes("Implementation Phase")
          ) {
            const match =
              initialPrompt.match(/ID #(\d+)/i) ||
              initialPrompt.match(/suggestion.*?(\d+)/i);
            const id = match ? match[1] : "1";
            return `<completed_suggestion id="${id}">Implemented</completed_suggestion>`;
          }

          return "";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.iterations, 2);
      assert.strictEqual(
        stageCalled,
        true,
        "Stage function must be called after inner implementation loop completes",
      );
      const proReviewerInvocations = sessionsInvoked.filter(
        (m, idx) => m === "gemini-3.1-pro-preview" && (idx === 0 || idx === 7),
      );
      assert.strictEqual(
        proReviewerInvocations.length,
        2,
        "Pro reviewer should be invoked on review iterations",
      );
      const proImplementerInvocations = sessionsInvoked.filter(
        (m, idx) => m === "gemini-3.1-pro-preview" && idx === 2,
      );
      assert.strictEqual(
        proImplementerInvocations.length,
        1,
        "Pro model should be invoked for the first implementation iteration",
      );
      const flashInvocations = sessionsInvoked.filter(
        (m) => m === "gemini-3.8-flash",
      );
      assert.strictEqual(
        flashInvocations.length,
        4,
        "Flash implementer should be invoked for remaining 4 suggestions",
      );
    },
  );

  await t.test(
    "runs thinking reviewer in isolated session and executes inner implementation loop",
    async () => {
      let reviewerIsolate = null;
      let implementerTurns = null;
      let stageCalled = false;
      let reviewIteration = 0;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 2,
        maxTurnsPerSession: 20,
        captureState: mockCapture,
        stageFn: async () => {
          stageCalled = true;
        },
        detectChange: () => ({ hasChanged: reviewIteration === 1 }),
        sessionRunner: async ({
          requestedModel,
          isolate,
          maxTurnsOverride,
          initialPrompt,
        }) => {
          if (isolate === true) {
            reviewerIsolate = isolate;
            reviewIteration++;
            if (reviewIteration === 1) {
              return `
1. [critical] Item 1
2. [critical] Item 2
3. [medium] Item 3
4. [medium] Item 4
5. [minor] Item 5
`;
            }
            return `
1. [minor] Item 1
2. [minor] Item 2
3. [minor] Item 3
4. [minor] Item 4
5. [minor] Item 5
`;
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            if (reviewIteration === 1) {
              return JSON.stringify({ critical: 2, medium: 2, minor: 1 });
            }
            return JSON.stringify({ critical: 0, medium: 0, minor: 5 });
          }
          if (initialPrompt.includes("remaining suggestion(s) to implement")) {
            implementerTurns = maxTurnsOverride;
            const match =
              initialPrompt.match(/ID #(\d+)/i) ||
              initialPrompt.match(/suggestion.*?(\d+)/i);
            const id = match ? match[1] : "1";
            return `<completed_suggestion id="${id}">COMPLETE</completed_suggestion>`;
          }
          return "";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(
        reviewerIsolate,
        true,
        "Thinking reviewer should run in isolated session",
      );
      assert.strictEqual(implementerTurns, 20);
      assert.strictEqual(stageCalled, true, "Inner loop must stage code");
    },
  );

  await t.test(
    "injects workspace skills context into reviewer and implementer prompts and instructions",
    async () => {
      let capturedReviewerPrompt = "";
      let capturedImplPrompt = "";

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      await runReviewRalphLoop({
        maxIterations: 1,
        captureState: mockCapture,
        stageFn: async () => {},
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async ({ requestedModel, initialPrompt }) => {
          if (
            initialPrompt &&
            initialPrompt.includes("senior technical reviewer")
          ) {
            capturedReviewerPrompt = initialPrompt;
            return `
1. [critical] Item 1
2. [medium] Item 2
3. [minor] Item 3
4. [minor] Item 4
5. [minor] Item 5
`;
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            return JSON.stringify({ critical: 1, medium: 1, minor: 3 });
          }
          capturedImplPrompt = initialPrompt;
          return '<completed_suggestion id="1">COMPLETE</completed_suggestion>';
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.ok(
        capturedReviewerPrompt.includes("<available_skills>"),
        "Reviewer prompt must include skills",
      );
      assert.ok(
        capturedReviewerPrompt.includes("Skills-First Mandate"),
        "Reviewer prompt must include mandate",
      );
      assert.strictEqual(
        (capturedReviewerPrompt.match(/<available_skills>/g) || []).length,
        1,
        "Reviewer prompt must include skills context exactly once",
      );

      assert.ok(
        capturedImplPrompt.includes("<available_skills>"),
        "Implementer prompt must include skills",
      );
      assert.ok(
        capturedImplPrompt.includes("Skills-First Mandate"),
        "Implementer prompt must include mandate",
      );
      assert.strictEqual(
        (capturedImplPrompt.match(/<available_skills>/g) || []).length,
        1,
        "Implementer prompt must include skills context exactly once",
      );
    },
  );

  await t.test(
    "passes blockTools to implementer session to prevent subagents and plan mode",
    async () => {
      let capturedBlockTools = null;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      await runReviewRalphLoop({
        maxIterations: 1,
        captureState: mockCapture,
        stageFn: async () => {},
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async ({
          requestedModel,
          initialPrompt,
          blockTools,
        }) => {
          if (
            initialPrompt &&
            initialPrompt.includes("senior technical reviewer")
          ) {
            return "1. [critical] Fix leak\n2. [medium] Validate input\n3. [minor] A\n4. [minor] B\n5. [minor] C";
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            return JSON.stringify({ critical: 1, medium: 1, minor: 3 });
          }
          capturedBlockTools = blockTools;
          return '<completed_suggestion id="1">Done</completed_suggestion>';
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.deepStrictEqual(capturedBlockTools, [
        "invoke_agent",
        "enter_plan_mode",
      ]);
    },
  );

  await t.test(
    "propagates inner loop implementation failure and does not falsely converge",
    async () => {
      let reviewIterations = 0;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 2,
        captureState: mockCapture,
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async ({ requestedModel, initialPrompt }) => {
          if (
            initialPrompt &&
            initialPrompt.includes("senior technical reviewer")
          ) {
            reviewIterations++;
            return `
1. [critical] Item 1
2. [medium] Item 2
3. [minor] Item 3
4. [minor] Item 4
5. [minor] Item 5
`;
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            return JSON.stringify({ critical: 1, medium: 1, minor: 3 });
          }
          // Implementer fails to complete suggestions
          return "I could not make the requested changes.";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      // Should not converge as successful on iteration 1
      assert.strictEqual(result.completed, false);
      assert.strictEqual(reviewIterations, 2);
    },
  );

  await t.test(
    "skips implementation when only minor suggestions exist",
    async () => {
      const sessionsInvoked = [];

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 1,
        captureState: mockCapture,
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async ({ requestedModel }) => {
          sessionsInvoked.push(requestedModel);
          if (requestedModel === "gemini-3.1-pro-preview") {
            return "1. [minor] Maintainability improvement";
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            return JSON.stringify({ critical: 0, medium: 0, minor: 1 });
          }
          return "Should not be called";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.deepStrictEqual(sessionsInvoked, [
        "gemini-3.1-pro-preview",
        "gemini-3.1-flash-lite",
      ]);
      assert.strictEqual(
        result.reason,
        "Only minor suggestions found (1 minor, 0 critical, 0 medium). No implementation needed.",
      );
    },
  );

  await t.test(
    "retries when reviewer text yields 0 total categorized suggestions",
    async () => {
      let reviewerInvocations = 0;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 2,
        captureState: mockCapture,
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async ({ requestedModel }) => {
          if (requestedModel === "gemini-3.1-pro-preview") {
            reviewerInvocations++;
            if (reviewerInvocations === 1) {
              return "Some prose without severity tags";
            }
            return "[]";
          }
          if (requestedModel === "gemini-3.1-flash-lite") {
            return JSON.stringify({ critical: 0, medium: 0, minor: 0 });
          }
          return "";
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(
        reviewerInvocations,
        2,
        "Should retry on zero categorized suggestions",
      );
      assert.strictEqual(result.success, true);
      assert.strictEqual(result.iterations, 2);
    },
  );
});

test("inner implementation Ralph Loop utilities and execution", async (t) => {
  await t.test(
    "parseIndividualSuggestions parses numbered and bulleted suggestions",
    () => {
      const raw = `
1. [critical] Fix buffer overflow in reader
2. [medium] Validate nil pointer before lookup
* [minor] Documentation typo in README
- (minor) Reorder imports alphabetically
5. Suggestion without severity tag
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 5);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.ok(items[0].text.includes("Fix buffer overflow"));

      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");

      assert.strictEqual(items[2].id, 3);
      assert.strictEqual(items[2].severity, "minor");

      assert.strictEqual(items[3].id, 4);
      assert.strictEqual(items[3].severity, "minor");

      assert.strictEqual(items[4].id, 5);
      assert.strictEqual(items[4].severity, "unclassified");
    },
  );

  await t.test(
    "parseIndividualSuggestions tolerates internal whitespace within severity tags",
    () => {
      const raw = `
1. [ critical ] Fix buffer overflow with internal spaces
2. (  medium  ) Validate nil pointer with spaces in parens
3. [   minor   ] Typo fix with extra spaces
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 3);
      assert.strictEqual(items[0].severity, "critical");
      assert.strictEqual(items[1].severity, "medium");
      assert.strictEqual(items[2].severity, "minor");
    },
  );

  await t.test(
    "parseIndividualSuggestions tolerates leading whitespace for indented suggestions",
    () => {
      const raw = `
  1. [critical] Fix memory leak in buffer
  2. [medium] Validate nil pointer before lookup
    3. [minor] Rename identifier
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 3);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.strictEqual(items[0].text, "Fix memory leak in buffer");

      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.strictEqual(items[1].text, "Validate nil pointer before lookup");

      assert.strictEqual(items[2].id, 3);
      assert.strictEqual(items[2].severity, "minor");
      assert.strictEqual(items[2].text, "Rename identifier");
    },
  );

  await t.test(
    "parseIndividualSuggestions tolerates indented suggestions while preserving indented sub-items",
    () => {
      const raw = `
  1. [critical] Fix memory leak in buffer:
     - Check stream handle
     - Free allocated memory
  2. [medium] Validate nil pointer:
     1. Verify obj is not nil
     2. Return error if nil
  3. [minor] Rename identifier
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 3);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.ok(
        items[0].text.includes(
          "Fix memory leak in buffer:\n     - Check stream handle\n     - Free allocated memory",
        ),
      );
      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.ok(
        items[1].text.includes(
          "Validate nil pointer:\n     1. Verify obj is not nil\n     2. Return error if nil",
        ),
      );
      assert.strictEqual(items[2].id, 3);
      assert.strictEqual(items[2].severity, "minor");
      assert.strictEqual(items[2].text, "Rename identifier");
    },
  );

  await t.test(
    "parseIndividualSuggestions handles empty or invalid input gracefully",
    () => {
      assert.deepStrictEqual(parseIndividualSuggestions(""), []);
      assert.deepStrictEqual(parseIndividualSuggestions(null), []);
    },
  );

  await t.test(
    "parseIndividualSuggestions handles multiline suggestion when text starts on next line",
    () => {
      const raw = `
1. [critical]
Fix buffer overflow in reader when handling large chunks.
2. [medium]
Validate nil pointer before lookup.
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 2);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.strictEqual(
        items[0].text,
        "Fix buffer overflow in reader when handling large chunks.",
      );

      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.strictEqual(items[1].text, "Validate nil pointer before lookup.");
    },
  );

  await t.test(
    "parseIndividualSuggestions preserves leading indentation on the first line when suggestion body starts on next line",
    () => {
      const raw = `
1. [critical]
   \`\`\`go
   stream.Close()
   \`\`\`
2. [medium]
   - Sub-bullet 1
   - Sub-bullet 2
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 2);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.strictEqual(items[0].text, "   ```go\n   stream.Close()\n   ```");

      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.strictEqual(items[1].text, "   - Sub-bullet 1\n   - Sub-bullet 2");
    },
  );

  await t.test(
    "parseIndividualSuggestions preserves newlines and indentation in multiline suggestions",
    () => {
      const raw = `
1. [critical] Fix resource leak in stream:
   Ensure stream is closed in finally block:
   \`\`\`go
   if stream != nil {
       stream.Close()
   }
   \`\`\`
2. [medium] Add validation:
   Check error return.
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 2);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.ok(
        items[0].text.includes(
          "Fix resource leak in stream:\n   Ensure stream is closed in finally block:",
        ),
      );
      assert.ok(items[0].text.includes("       stream.Close()"));
      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.ok(
        items[1].text.includes("Add validation:\n   Check error return."),
      );
    },
  );

  await t.test(
    "parseIndividualSuggestions prevents redundant newline accumulation across multiple blank lines",
    () => {
      const raw = `
1. [critical] Fix resource leak:

   Paragraph 1.


   Paragraph 2.
2. [minor] Clean up comments.
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 2);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.strictEqual(
        items[0].text,
        "Fix resource leak:\n\n   Paragraph 1.\n\n   Paragraph 2.",
      );
    },
  );

  await t.test(
    "parseIndividualSuggestions preserves indented sub-bullets and numbered sub-lists without splitting",
    () => {
      const raw = `
1. [critical] Fix the memory leak:
   - In file A: close buffer
   - In file B: release lock
2. [medium] Validate nil pointer before lookup:
   1. Check if obj is nil
   2. Return error if nil
3. [minor] Documentation typo
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 3);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.ok(
        items[0].text.includes(
          "Fix the memory leak:\n   - In file A: close buffer\n   - In file B: release lock",
        ),
      );
      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.ok(
        items[1].text.includes(
          "Validate nil pointer before lookup:\n   1. Check if obj is nil\n   2. Return error if nil",
        ),
      );
      assert.strictEqual(items[2].id, 3);
      assert.strictEqual(items[2].severity, "minor");
      assert.strictEqual(items[2].text, "Documentation typo");
    },
  );

  await t.test(
    "parseIndividualSuggestions ignores header-like lines inside markdown code blocks",
    () => {
      const raw = `
1. [critical] Add database migration script:
   Here is the migration script:
   \`\`\`bash
   1. echo "migration step 1"
   2. echo "migration step 2"
   * echo "bullet item"
   \`\`\`
   Verify the logs afterwards.
2. [medium] Add validation for null user:
   \`\`\`go
   // 1. check pointer
   if user == nil {
       return ErrNotFound
   }
   \`\`\`
3. [minor] Update documentation
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 3);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.ok(items[0].text.includes('1. echo "migration step 1"'));
      assert.ok(items[0].text.includes('2. echo "migration step 2"'));
      assert.ok(items[0].text.includes('* echo "bullet item"'));
      assert.ok(items[0].text.includes("Verify the logs afterwards."));

      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.ok(items[1].text.includes("// 1. check pointer"));
      assert.ok(items[1].text.includes("if user == nil {"));

      assert.strictEqual(items[2].id, 3);
      assert.strictEqual(items[2].severity, "minor");
      assert.strictEqual(items[2].text, "Update documentation");
    },
  );

  await t.test(
    "parseIndividualSuggestions resets unclosed code block when next suggestion header is encountered",
    () => {
      const raw = `
1. [critical] Fix buffer overflow in reader:
   \`\`\`go
   buf := make([]byte, 1024)
   // note: unclosed code block without closing fence
2. [medium] Validate nil pointer before lookup:
   Check nil obj.
3. [minor] Rename identifier:
   Rename x to counter.
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 3);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.ok(items[0].text.includes("buf := make([]byte, 1024)"));

      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.ok(items[1].text.includes("Validate nil pointer before lookup"));

      assert.strictEqual(items[2].id, 3);
      assert.strictEqual(items[2].severity, "minor");
      assert.ok(items[2].text.includes("Rename identifier"));
    },
  );

  await t.test(
    "parseIndividualSuggestions supports tilde (~~~) code fences and ignores internal headers",
    () => {
      const raw = `
1. [critical] Add database migration script:
   Here is the migration script using tildes:
   ~~~bash
   1. echo "step 1 inside tilde block"
   2. echo "step 2 inside tilde block"
   * bullet item
   ~~~
   Verify the logs afterwards.
2. [medium] Validate nil pointer before lookup:
   Check nil obj.
3. [minor] Rename identifier
`;
      const items = parseIndividualSuggestions(raw);
      assert.strictEqual(items.length, 3);
      assert.strictEqual(items[0].id, 1);
      assert.strictEqual(items[0].severity, "critical");
      assert.ok(items[0].text.includes('1. echo "step 1 inside tilde block"'));
      assert.ok(items[0].text.includes('2. echo "step 2 inside tilde block"'));
      assert.ok(items[0].text.includes("* bullet item"));
      assert.ok(items[0].text.includes("Verify the logs afterwards."));

      assert.strictEqual(items[1].id, 2);
      assert.strictEqual(items[1].severity, "medium");
      assert.ok(items[1].text.includes("Validate nil pointer before lookup"));

      assert.strictEqual(items[2].id, 3);
      assert.strictEqual(items[2].severity, "minor");
      assert.strictEqual(items[2].text, "Rename identifier");
    },
  );

  await t.test(
    "extractCompletedSuggestionIds parses tags, comma-separated IDs, conversational text, and defaults",
    () => {
      const remaining = [{ id: 1 }, { id: 2 }, { id: 3 }];

      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          '<completed_suggestion id="2">Done</completed_suggestion>',
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          '<completed_suggestion id="1, 2">Done</completed_suggestion>',
          remaining,
        ),
        [1, 2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          '<completed_suggestion id="2">Implemented suggestion 1</completed_suggestion>',
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          '<completed_suggestion id="2"/>',
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          '<completed_suggestion id="3">Unclosed tag without closing delimiter',
          remaining,
        ),
        [3],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>3</completed_suggestion>",
          remaining,
        ),
        [3],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>1, 3</completed_suggestion>",
          remaining,
        ),
        [1, 3],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>ID: 2</completed_suggestion>",
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>#2</completed_suggestion>",
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion># 2</completed_suggestion>",
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>id 2</completed_suggestion>",
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>suggestion 2</completed_suggestion>",
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>suggestions 1, 2</completed_suggestion>",
          remaining,
        ),
        [1, 2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>suggestion #2</completed_suggestion>",
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "<completed_suggestion>Fixed bug on line 2</completed_suggestion>",
          remaining,
        ),
        [],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          '<completed_suggestion id="2">First</completed_suggestion><completed_suggestion id="3">Second</completed_suggestion>',
          remaining,
        ),
        [2, 3],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "I have implemented suggestions 2, 3.",
          remaining,
        ),
        [2, 3],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "Resolved suggestion #1 successfully",
          remaining,
        ),
        [1],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds("Addressing suggestion id 3", remaining),
        [3],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds("Fixed suggestion 2", remaining),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds("Completed suggestion #2", remaining),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "suggestion 2 could not be implemented",
          remaining,
        ),
        [],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "Working on suggestion #2 for this turn",
          remaining,
        ),
        [],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "Here is the fix for suggestion 2",
          remaining,
        ),
        [],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "I considered suggestion 99, but implemented suggestion 2",
          remaining,
        ),
        [2],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          "I made changes without specific tag",
          remaining,
        ),
        [],
      );
      assert.deepStrictEqual(
        extractCompletedSuggestionIds(
          '<completed_suggestion id="99">Unknown</completed_suggestion>',
          remaining,
        ),
        [],
      );
    },
  );

  await t.test(
    "buildBatchSuggestionImplementationPrompt formats remaining items and batch instructions",
    () => {
      const prompt = buildBatchSuggestionImplementationPrompt({
        remainingSuggestions: [
          {
            id: 1,
            severity: "critical",
            text: "Fix buffer leak\n<leak>details</leak>",
          },
          { id: 2, severity: "medium", text: "Add validation" },
        ],
        currentSuggestionIndex: 1,
        totalSuggestions: 2,
      });
      assert.ok(prompt.includes("Implementation Phase - Step 1 of 2."));
      assert.ok(
        prompt.includes("There are 2 remaining suggestion(s) to implement:"),
      );
      assert.ok(
        prompt.includes(
          "1. [critical] (Original #1) Fix buffer leak\n<leak>details</leak>",
        ),
      );
      assert.ok(prompt.includes("UP TO 3 SUGGESTIONS"));
      assert.ok(
        prompt.includes(
          "The active plan from 'plans/current.md' is already embedded below in <active_plan>. Do NOT waste turns reading it from disk.",
        ),
      );
      assert.ok(
        prompt.includes(
          "formulate a clear implementation plan for addressing the suggestions",
        ),
      );
      assert.ok(
        prompt.includes(
          "CRITICAL PLAN INTEGRITY RULE: Do NOT overwrite or wipe 'plans/current.md'",
        ),
      );
      assert.ok(prompt.includes("WHICH suggestions were implemented"));
    },
  );

  await t.test(
    "buildBatchSuggestionImplementationPrompt caps validation output length safely",
    () => {
      const massiveStdout = "stdout-line\n".repeat(300); // ~3600 chars
      const massiveStderr = "stderr-line\n".repeat(300); // ~3600 chars
      const prompt = buildBatchSuggestionImplementationPrompt({
        remainingSuggestions: [
          { id: 1, severity: "critical", text: "Fix bug" },
        ],
        lastValidation: {
          success: false,
          stdout: massiveStdout,
          stderr: massiveStderr,
        },
      });

      assert.ok(prompt.includes("PREVIOUS VALIDATION FAILURE"));
      assert.ok(prompt.includes("[truncated, showing last 2000 chars]"));
      // Ensure the prompt does not contain the entire 3600 character output
      assert.ok(!prompt.includes(massiveStdout));
      assert.ok(!prompt.includes(massiveStderr));
    },
  );

  await t.test(
    "runImplementationRalphLoop iterates per suggestion, removes completed items, and stages code",
    async () => {
      let stageCalled = false;
      const promptsReceived = [];
      const loggedMessages = [];
      let innerIter = 0;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runImplementationRalphLoop({
        suggestions: [
          { id: 1, severity: "critical", text: "Fix leak" },
          { id: 2, severity: "medium", text: "Add validation" },
        ],
        captureState: mockCapture,
        stageFn: async () => {
          stageCalled = true;
        },
        sessionRunner: async ({ initialPrompt }) => {
          innerIter++;
          promptsReceived.push(initialPrompt);
          if (innerIter === 1) {
            return '<completed_suggestion id="1">Fixed leak</completed_suggestion>';
          }
          return '<completed_suggestion id="2">Added validation</completed_suggestion>';
        },
        validator: async () => ({ success: true }),
        logger: (msg) => loggedMessages.push(msg),
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.implementedCount, 2);
      assert.strictEqual(result.totalSuggestions, 2);
      assert.strictEqual(result.iterations, 2);
      assert.strictEqual(
        stageCalled,
        true,
        "Code must be staged after completing all suggestions",
      );
      assert.ok(promptsReceived[0].includes("Fix leak"));
      assert.ok(promptsReceived[1].includes("Add validation"));
      assert.ok(
        loggedMessages.some((m) =>
          m.includes(
            "Waiting for implementation of suggestion(s) #1, #2 (e.g., ID #1 [critical]",
          ),
        ),
      );
    },
  );

  await t.test(
    "runImplementationRalphLoop scales turns dynamically using calculateScaledTurns",
    async () => {
      const capturedTurns = [];
      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      // 10 suggestions with base turns 20 scales to 40 turns
      const tenSuggestions = Array.from({ length: 10 }, (_, i) => ({
        id: i + 1,
        severity: "minor",
        text: `Fix issue ${i + 1}`,
      }));

      await runImplementationRalphLoop({
        suggestions: tenSuggestions,
        maxIterations: 1,
        maxTurnsPerSession: 20,
        captureState: mockCapture,
        sessionRunner: async ({ maxTurnsOverride }) => {
          capturedTurns.push(maxTurnsOverride);
          return '<completed_suggestion id="1">Done</completed_suggestion>';
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(
        capturedTurns[0],
        40,
        "First iteration turn capacity must be scaled based on 10 suggestions",
      );
      assert.strictEqual(
        capturedTurns[1],
        36,
        "Second iteration turn capacity must be scaled based on 9 suggestions",
      );
    },
  );

  await t.test(
    "runImplementationRalphLoop retries suggestion on validation failure without removing item",
    async () => {
      let validatorCallCount = 0;
      let stageCalled = false;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runImplementationRalphLoop({
        suggestions: [{ id: 1, severity: "critical", text: "Fix leak" }],
        maxIterations: 3,
        captureState: mockCapture,
        stageFn: async () => {
          stageCalled = true;
        },
        sessionRunner: async () =>
          '<completed_suggestion id="1">Fixed</completed_suggestion>',
        validator: async () => {
          validatorCallCount++;
          // Fail on iteration 1, succeed on iteration 2
          return validatorCallCount === 1
            ? { success: false, stdout: "", stderr: "Linter error" }
            : { success: true };
        },
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.implementedCount, 1);
      assert.strictEqual(result.iterations, 2);
      assert.strictEqual(stageCalled, true);
    },
  );

  await t.test(
    "runImplementationRalphLoop handles empty suggestions without iterating",
    async () => {
      let sessionRunnerCalled = false;
      let stageCalled = false;

      const result = await runImplementationRalphLoop({
        suggestions: [],
        suggestionsText: "",
        stageFn: async () => {
          stageCalled = true;
        },
        sessionRunner: async () => {
          sessionRunnerCalled = true;
          return "";
        },
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.completed, true);
      assert.strictEqual(result.totalSuggestions, 0);
      assert.strictEqual(result.implementedCount, 0);
      assert.strictEqual(result.iterations, 0);
      assert.strictEqual(sessionRunnerCalled, false);
      assert.strictEqual(stageCalled, false);
    },
  );

  await t.test(
    "runImplementationRalphLoop does not stage code if loop hits maxIterations and fails validation",
    async () => {
      let stageCalled = false;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runImplementationRalphLoop({
        suggestions: [{ id: 1, severity: "critical", text: "Fix leak" }],
        maxIterations: 2,
        captureState: mockCapture,
        stageFn: async () => {
          stageCalled = true;
        },
        sessionRunner: async () =>
          '<completed_suggestion id="1">Attempted</completed_suggestion>',
        validator: async () => ({
          success: false,
          stdout: "",
          stderr: "Compile error",
        }),
        logger: () => {},
      });

      assert.strictEqual(result.success, false);
      assert.strictEqual(result.completed, false);
      assert.strictEqual(result.implementedCount, 0);
      assert.strictEqual(result.iterations, 3); // 2 max iterations + 1 fallback attempt
      assert.strictEqual(
        stageCalled,
        false,
        "Code must not be staged when suggestions remain uncompleted",
      );
    },
  );

  await t.test(
    "runImplementationRalphLoop logs attempted suggestion ID on validation failure when agent chooses non-first suggestion",
    async () => {
      const loggedMessages = [];
      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runImplementationRalphLoop({
        suggestions: [
          { id: 1, severity: "critical", text: "Fix leak" },
          { id: 2, severity: "medium", text: "Add validation" },
        ],
        maxIterations: 1,
        captureState: mockCapture,
        stageFn: async () => {},
        sessionRunner: async () =>
          '<completed_suggestion id="2">Attempted fix 2</completed_suggestion>',
        validator: async () => ({
          success: false,
          stdout: "",
          stderr: "Validation error",
        }),
        logger: (msg) => loggedMessages.push(msg),
      });

      assert.strictEqual(result.success, false);
      assert.strictEqual(result.completed, false);
      assert.ok(
        loggedMessages.some((m) =>
          m.includes(
            "Running deterministic validation pipeline for suggestion ID(s) 2",
          ),
        ),
      );
      assert.ok(
        loggedMessages.some((m) =>
          m.includes("Validation failed for suggestion ID(s) 2"),
        ),
      );
    },
  );

  await t.test(
    "runImplementationRalphLoop does not dequeue suggestion when no completion tag and no repo changes",
    async () => {
      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runImplementationRalphLoop({
        suggestions: [{ id: 1, severity: "critical", text: "Fix leak" }],
        maxIterations: 2,
        captureState: mockCapture,
        detectChange: () => ({ hasChanged: false }),
        sessionRunner: async () => "I could not make the requested change.",
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, false);
      assert.strictEqual(result.completed, false);
      assert.strictEqual(result.implementedCount, 0);
    },
  );

  await t.test(
    "runImplementationRalphLoop uses proModel for iteration 1 and flashModel for iterations 2+",
    async () => {
      const modelsUsed = [];
      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      let iteration = 0;
      const result = await runImplementationRalphLoop({
        suggestions: [
          { id: 1, severity: "critical", text: "Fix leak" },
          { id: 2, severity: "medium", text: "Add validation" },
        ],
        model: "flash-model-custom",
        proModel: "pro-model-custom",
        captureState: mockCapture,
        stageFn: async () => {},
        sessionRunner: async ({ requestedModel }) => {
          iteration++;
          modelsUsed.push(requestedModel);
          return `<completed_suggestion id="${iteration}">Done</completed_suggestion>`;
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.iterations, 2);
      assert.strictEqual(
        modelsUsed[0],
        "pro-model-custom",
        "Iteration 1 must use proModel",
      );
      assert.strictEqual(
        modelsUsed[1],
        "flash-model-custom",
        "Iteration 2 must use flashModel",
      );
    },
  );

  await t.test(
    "runReviewRalphLoop uses proModel for review iterations 1-3 and flashModel for iterations 4+",
    async () => {
      const reviewModelsUsed = [];
      let reviewIteration = 0;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 4,
        models: {
          pro: "my-pro-model",
          flash: "my-flash-model",
          flash_lite: "my-flash-lite-model",
        },
        captureState: mockCapture,
        stageFn: async () => {},
        detectChange: () => ({ hasChanged: true }),
        sessionRunner: async ({ requestedModel, initialPrompt }) => {
          if (initialPrompt.includes("senior technical reviewer")) {
            reviewIteration++;
            reviewModelsUsed.push(requestedModel);
            if (reviewIteration < 4) {
              return "1. [critical] Issue\n2. [medium] Issue 2\n3. [minor] A\n4. [minor] B\n5. [minor] C";
            }
            // Converge on iteration 4
            return "1. [minor] A\n2. [minor] B\n3. [minor] C\n4. [minor] D\n5. [minor] E";
          }
          if (requestedModel === "my-flash-lite-model") {
            if (reviewIteration < 4) {
              return JSON.stringify({ critical: 1, medium: 1, minor: 3 });
            }
            return JSON.stringify({ critical: 0, medium: 0, minor: 5 });
          }
          // Implementer
          return '<completed_suggestion id="1">Fixed</completed_suggestion>';
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.iterations, 4);
      assert.deepStrictEqual(
        reviewModelsUsed,
        ["my-pro-model", "my-pro-model", "my-pro-model", "my-flash-model"],
        "Review iterations 1-3 must use proModel, iteration 4 must switch to flashModel",
      );
    },
  );

  await t.test(
    "runImplementationRalphLoop uses proModel for critical suggestions in subsequent iterations",
    async () => {
      const modelsUsed = [];
      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      let iteration = 0;
      const result = await runImplementationRalphLoop({
        suggestions: [
          { id: 1, severity: "minor", text: "Typo" },
          { id: 2, severity: "critical", text: "Security bug" },
        ],
        model: "flash-model-custom",
        proModel: "pro-model-custom",
        captureState: mockCapture,
        stageFn: async () => {},
        sessionRunner: async ({ requestedModel }) => {
          iteration++;
          modelsUsed.push(requestedModel);
          return `<completed_suggestion id="${iteration}">Done</completed_suggestion>`;
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(result.iterations, 2);
      assert.strictEqual(
        modelsUsed[0],
        "pro-model-custom",
        "Iteration 1 uses proModel tier",
      );
      assert.strictEqual(
        modelsUsed[1],
        "pro-model-custom",
        "Iteration 2 uses proModel for critical suggestion",
      );
    },
  );

  await t.test(
    "evaluates stateChanged per attempt rather than stale baseline across retries",
    async () => {
      let innerAttempts = 0;
      const states = [
        {
          diff: "original",
          status: "",
          diffHash: "h0",
          plan: "",
          changedFiles: [],
          modifiedFilesContext: "",
        },
        {
          diff: "after-attempt2-no-change",
          status: "",
          diffHash: "h1",
          plan: "",
          changedFiles: [],
          modifiedFilesContext: "",
        },
        {
          diff: "before-attempt2",
          status: "",
          diffHash: "h1",
          plan: "",
          changedFiles: [],
          modifiedFilesContext: "",
        },
        {
          diff: "after-attempt2-no-change",
          status: "",
          diffHash: "h1",
          plan: "",
          changedFiles: [],
          modifiedFilesContext: "",
        },
      ];

      let captureIndex = 0;
      const mockCapture = async () => {
        const s = states[Math.min(captureIndex, states.length - 1)];
        captureIndex++;
        return s;
      };

      let secondAttemptStateChanged = null;

      await runImplementationRalphLoop({
        suggestions: [{ id: 1, severity: "minor", text: "Task 1" }],
        maxIterations: 1,
        captureState: mockCapture,
        detectChange: (before, after) => ({
          hasChanged: before.diffHash !== after.diffHash,
        }),
        sessionRunner: async () => {
          innerAttempts++;
          return "Attempt " + innerAttempts;
        },
        validator: async ({ iteration, stateChanged }) => {
          if (iteration === 1) {
            return { success: false, stderr: "Validation error" };
          }
          secondAttemptStateChanged = stateChanged;
          return { success: true };
        },
        logger: () => {},
      });
      assert.strictEqual(innerAttempts, 2);
      assert.strictEqual(
        secondAttemptStateChanged,
        false,
        "Second attempt made no edits beyond attempt 1, so stateChanged must be false",
      );
    },
  );

  await t.test(
    "runImplementationRalphLoop respects custom implementerInstructions option",
    async () => {
      let capturedPrompt = "";
      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M test_file_test.go",
        diffHash: "diff-hash-1",
        plan: "## Active Plan",
        changedFiles: ["test_file_test.go"],
        modifiedFilesContext: '<file path="test_file_test.go">content</file>',
      });

      await runImplementationRalphLoop({
        suggestions: [{ id: 1, severity: "minor", text: "Task 1" }],
        maxIterations: 1,
        implementerInstructions:
          "CUSTOM_TEST_ONLY_CONSTRAINT: Modify only test files.",
        captureState: mockCapture,
        detectChange: () => ({ hasChanged: true }),
        sessionRunner: async ({ initialPrompt }) => {
          capturedPrompt = initialPrompt;
          return '<completed_suggestion id="1">Fixed</completed_suggestion>';
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.ok(
        capturedPrompt.includes(
          "CUSTOM_TEST_ONLY_CONSTRAINT: Modify only test files.",
        ),
        "Custom implementer instructions must be prepended to the initialPrompt",
      );
    },
  );

  await t.test(
    "runReviewRalphLoop respects custom proReviewIterations and proImplementationIterations",
    async () => {
      const reviewModelsUsed = [];
      let reviewIteration = 0;

      const mockCapture = async () => ({
        diff: "staged diff",
        status: "M main.go",
        diffHash: "diff-hash",
        plan: "## Review Plan",
        changedFiles: ["main.go"],
        modifiedFilesContext: '<file path="main.go">content</file>',
      });

      const result = await runReviewRalphLoop({
        maxIterations: 6,
        proReviewIterations: 5,
        proImplementationIterations: 2,
        models: {
          pro: "pro-5",
          flash: "flash-5",
          flash_lite: "flash-lite-5",
        },
        captureState: mockCapture,
        stageFn: async () => {},
        detectChange: () => ({ hasChanged: true }),
        sessionRunner: async ({ requestedModel, initialPrompt }) => {
          if (initialPrompt.includes("senior technical reviewer")) {
            reviewIteration++;
            reviewModelsUsed.push(requestedModel);
            if (reviewIteration < 6) {
              return "1. [critical] Issue\n2. [medium] Issue 2\n3. [minor] A\n4. [minor] B\n5. [minor] C";
            }
            return "1. [minor] A\n2. [minor] B\n3. [minor] C\n4. [minor] D\n5. [minor] E";
          }
          if (requestedModel === "flash-lite-5") {
            if (reviewIteration < 6) {
              return JSON.stringify({ critical: 1, medium: 1, minor: 3 });
            }
            return JSON.stringify({ critical: 0, medium: 0, minor: 5 });
          }
          return '<completed_suggestion id="1">Fixed</completed_suggestion>';
        },
        validator: async () => ({ success: true }),
        logger: () => {},
      });

      assert.strictEqual(result.success, true);
      assert.deepStrictEqual(
        reviewModelsUsed,
        ["pro-5", "pro-5", "pro-5", "pro-5", "pro-5", "flash-5"],
        "First 5 review iterations must use pro model, iteration 6 switches to flash",
      );
    },
  );
});
