#!/usr/bin/env node
import { execFile } from "node:child_process";
import crypto from "node:crypto";
import fs from "node:fs/promises";
import path from "node:path";
import { stdin as input, stdout as output } from "node:process";
import * as readline from "node:readline/promises";
import { promisify } from "node:util";

import { tool, z } from "@google/gemini-cli-sdk";

import { PREFERRED_DEFAULT_BRANCHES } from "../skills/git/scripts/create-branch.js";

// Import our deterministic modular libraries
import {
  flushLogs,
  initializeAgentRunner,
  isAbortError,
  runAgentSession,
} from "./lib/agent-runner.js";
import { stripDiffMetadata, validateMessage } from "./lib/git-release.js";
import { getPRComments } from "./lib/github-context.js";
import { runQAPipeline } from "./lib/qa-runner.js";
import {
  DEFAULT_IMPLEMENTER_INSTRUCTIONS,
  runReviewRalphLoop,
} from "./lib/ralph-loop.js";
import {
  exists,
  getRepoRoot,
  loadAgentInstructions,
  savePlanFromJSON,
  validateAgentOutput,
} from "./lib/utils.js";

const execFileAsync = promisify(execFile);

/**
 * Retrieves git diff relative to HEAD using the git skill script.
 * @param {string} [targetRepoRoot]
 * @returns {Promise<string>}
 */
async function getGitDiff(targetRepoRoot) {
  const root = targetRepoRoot || (await getRepoRoot());
  const diffScript = path.join(
    root,
    ".gemini/skills/git/scripts/get_diff_to_review.js",
  );
  const { stdout } = await execFileAsync(
    process.execPath,
    [diffScript, "--raw", "--head"],
    {
      cwd: root,
      maxBuffer: 10 * 1024 * 1024,
      timeout: 60000,
    },
  );
  return stdout ? stdout.trimEnd() : "";
}

/**
 * Stages changes and executes a validated commit using git skill scripts.
 * @param {string} commitMessage
 * @param {string} [targetRepoRoot]
 * @returns {Promise<{stdout: string, stderr: string}>}
 */
async function stageAndCommit(commitMessage, targetRepoRoot) {
  const root = targetRepoRoot || (await getRepoRoot());
  const stageScript = path.join(
    root,
    ".gemini/skills/git/scripts/stage-files.js",
  );
  await execFileAsync(
    process.execPath,
    [stageScript, "--all", "--include-untracked"],
    {
      cwd: root,
      maxBuffer: 10 * 1024 * 1024,
      timeout: 60000,
    },
  );
  const commitScript = path.join(root, ".gemini/skills/git/scripts/commit.js");
  const { stdout, stderr } = await execFileAsync(
    process.execPath,
    [commitScript, "-m", commitMessage],
    {
      cwd: root,
      maxBuffer: 10 * 1024 * 1024,
      timeout: 60000,
    },
  );
  return { stdout, stderr };
}

/**
 * Verifies that only test files were modified during a Test-only task.
 * @param {string} repoRoot
 * @returns {Promise<boolean>} True if changes are strictly within test scope boundaries, false otherwise.
 */
async function verifyTestScopeBoundaries(repoRoot) {
  const modifiedFiles = new Set();

  // Check if any product files were altered using get_modified_paths skill script via execFileAsync
  try {
    const scriptPath = path.join(
      repoRoot,
      ".gemini/skills/git/scripts/get_modified_paths.js",
    );
    const { stdout: modifiedPathsOutput } = await execFileAsync(
      process.execPath,
      [scriptPath],
      {
        cwd: repoRoot,
      },
    );
    modifiedPathsOutput
      .split("\n")
      .map((p) => p.trim())
      .filter(
        (f) =>
          Boolean(f) &&
          !f.startsWith("❌") &&
          !f.startsWith("===") &&
          !f.startsWith("⚠️"),
      )
      .forEach((p) => modifiedFiles.add(p));
  } catch (err) {
    console.error(
      `❌ Failed closed: Could not inspect modified paths for scope boundary check: ${err?.message || err}`,
    );
    return false;
  }

  // Also ensure staged files in the index are captured to prevent bypassing checks via staging
  try {
    const stagedScriptPath = path.join(
      repoRoot,
      ".gemini/skills/git/scripts/get_staged_files.js",
    );
    const { stdout: stagedPathsOutput } = await execFileAsync(
      process.execPath,
      [stagedScriptPath, "--name-only"],
      {
        cwd: repoRoot,
      },
    );
    stagedPathsOutput
      .split("\n")
      .map((p) => p.trim())
      .filter(
        (f) =>
          Boolean(f) &&
          !f.startsWith("❌") &&
          !f.startsWith("===") &&
          !f.startsWith("⚠️"),
      )
      .forEach((f) => modifiedFiles.add(f));
  } catch (err) {
    console.error(
      `❌ Failed closed: Could not inspect staged files for scope boundary check: ${err?.message || err}`,
    );
    return false;
  }

  // Also check any branch changes relative to base branch (e.g. if changes were committed)
  try {
    const branchChangesScript = path.join(
      repoRoot,
      ".gemini/skills/git/scripts/get_branch_changes.js",
    );
    const { stdout: branchChangesOutput } = await execFileAsync(
      process.execPath,
      [branchChangesScript, "--name-only"],
      {
        cwd: repoRoot,
      },
    );
    branchChangesOutput
      .split("\n")
      .map((p) => p.trim())
      .filter(
        (f) =>
          Boolean(f) &&
          !f.startsWith("❌") &&
          !f.startsWith("===") &&
          !f.startsWith("⚠️"),
      )
      .forEach((f) => modifiedFiles.add(f));
  } catch (err) {
    console.error(
      `❌ Failed closed: Could not inspect branch changes for scope boundary check: ${err?.message || err}`,
    );
    return false;
  }

  const invalidChanges = Array.from(modifiedFiles).filter((f) => {
    if (f === "dev/null") {
      return false;
    }
    const parts = f.split("/");
    const isInTestDir = parts
      .slice(0, -1)
      .some(
        (p) =>
          p === "test" ||
          p === "tests" ||
          p === "__tests__" ||
          p === "__snapshots__" ||
          p === "testdata" ||
          p === "mocks",
      );
    const isGoTest = f.endsWith("_test.go");
    const isJsTest =
      f.endsWith(".test.js") ||
      f.endsWith(".spec.js") ||
      f.endsWith(".test.mjs") ||
      f.endsWith(".spec.mjs") ||
      f.endsWith(".test.ts") ||
      f.endsWith(".spec.ts");

    return !isInTestDir && !isGoTest && !isJsTest;
  });

  if (invalidChanges.length > 0) {
    console.error(
      `❌ Safety violation: Product files were modified during a Test-only task: ${invalidChanges.join(", ")}`,
    );
    return false;
  }

  return true;
}

const rl = readline.createInterface({ input, output });

let isShuttingDown = false;

// Prevent unhandled AbortError from crashing Node process when SDK streams abort
process.on("uncaughtException", (err, origin) => {
  if (isAbortError(err)) {
    return;
  }
  if (isShuttingDown) {
    return;
  }
  isShuttingDown = true;
  console.error(`❌ Uncaught exception (${origin}):`, err);
  try {
    if (rl) {
      rl.close();
    }
  } catch (closeErr) {
    console.debug(
      `[DEBUG] Error closing readline interface: ${closeErr?.message || closeErr}`,
    );
  }

  try {
    flushLogs()
      .catch((cleanupErr) => {
        console.debug(
          `[DEBUG] Error during uncaughtException cleanup: ${cleanupErr?.message || cleanupErr}`,
        );
      })
      .finally(() => {
        process.exit(1);
      });
    const forceExitTimer = setTimeout(() => {
      process.exit(1);
    }, 2000);
    forceExitTimer.unref();
  } catch {
    process.exit(1);
  }
});

process.on("SIGINT", () => {
  if (isShuttingDown) {
    process.exit(130);
  }
  isShuttingDown = true;
  console.error("\n❌ Process interrupted by user (SIGINT).");
  try {
    if (rl) {
      rl.close();
    }
  } catch (closeErr) {
    console.debug(
      `[DEBUG] Error closing readline interface: ${closeErr?.message || closeErr}`,
    );
  }

  try {
    flushLogs()
      .catch((cleanupErr) => {
        console.debug(
          `[DEBUG] Error during SIGINT cleanup: ${cleanupErr?.message || cleanupErr}`,
        );
      })
      .finally(() => {
        process.exit(130);
      });
    const forceExitTimer = setTimeout(() => {
      process.exit(130);
    }, 2000);
    forceExitTimer.unref();
  } catch {
    process.exit(130);
  }
});

const USE_CASES = {
  1: { name: "Bugfix", desc: "Address a problem the user is facing" },
  2: { name: "Feature", desc: "Add a new feature the user would like" },
  3: { name: "Refactor", desc: "Refactor something without changing behavior" },
  4: { name: "Test", desc: "Add new tests without changing behavior" },
  5: { name: "PR Comments", desc: "Address PR comments programmatically" },
  6: {
    name: "Review Loop",
    desc: "Iterative loop to review staged code until no changes exist",
  },
};

// Define custom tool for human-in-the-loop clarification during planning
const askUserTool = tool(
  {
    name: "ask_user",
    description:
      "Ask the human user a clarifying question when critical setup or context details are missing.",
    inputSchema: z.object({
      question: z
        .string()
        .describe("The exact clarifying question to prompt the user with."),
    }),
  },
  async (params) => {
    console.log(`\n\n🤖 [Agent requested input]: ${params.question}`);
    const answer = await rl.question("👉 Your Answer: ");
    return { answer };
  },
);

/**
 * Main Orchestrated Workflow Execution
 */
async function main() {
  {
    const { models } = await initializeAgentRunner();
    const repoRoot = await getRepoRoot();
    console.log("\n==================================================");
    console.log("🤖 [Gemini CLI Orchestrator] - Modular Pipeline");
    console.log("==================================================");

    // Step 1: Use Case Selection
    console.log("\n🎯 Select your interaction use case:");
    for (const [key, uc] of Object.entries(USE_CASES)) {
      console.log(`  ${key}) [${uc.name}] - ${uc.desc}`);
    }

    let selection = "";
    while (!USE_CASES[selection]) {
      const inputSel = await rl.question("\n👉 Selection (1-6): ");
      selection = inputSel.trim();
      if (!USE_CASES[selection]) {
        console.log(
          "⚠️ Invalid selection. Please enter a number between 1 and 6.",
        );
      }
    }

    const useCase = USE_CASES[selection];
    console.log(`\n🚀 Starting Tailored Workflow for: [${useCase.name}]`);

    let objective;
    if (selection === "6") {
      const currentPlanPath = path.join(repoRoot, "plans/current.md");
      if (!(await exists(currentPlanPath))) {
        console.error(
          "❌ Error: plans/current.md not found. An existing plan is required to run the review loop.",
        );
        process.exitCode = 1;
        return;
      }
      const planText = await fs.readFile(currentPlanPath, "utf8");
      const match = planText.match(/## 🎯 OBJECTIVE\s+([\s\S]*?)(?=\n## |$)/);
      const rawObjective = match
        ? match[1].trim()
        : "Review staged code against existing plan in plans/current.md";
      // Take the first non-empty paragraph and apply a character ceiling to prevent excessively long strings / branch names
      const firstParagraph = rawObjective.split(/\n\s*\n/)[0].trim();
      if (firstParagraph.length <= 200) {
        objective = firstParagraph;
      } else {
        const truncated = Array.from(firstParagraph).slice(0, 200).join("");
        const lastSpace = truncated.lastIndexOf(" ");
        objective = (
          lastSpace > 50 ? truncated.slice(0, lastSpace) : truncated
        ).trim();
      }
      console.log(
        `ℹ️ Review Loop will use existing plan from plans/current.md (Objective: "${objective.slice(0, 80)}...")`,
      );
    } else {
      objective = await rl.question("\n🎯 Enter your objective/task: ");
      if (!objective || objective.trim() === "") {
        console.error("❌ Error: An objective is required.");
        process.exitCode = 1;
        return;
      }
    }

    const planSchema = z.object({
      title: z.string(),
      objective: z.string(),
      scope_boundaries: z.object({
        in_scope: z.array(z.string()),
        out_of_scope: z.array(z.string()),
      }),
      exit_criteria: z.array(z.string()),
      implementation_tasks: z.array(z.string()),
    });

    // ==========================================
    // PHASE 1: PLANNING (Read-Only)
    // ==========================================
    console.log("\n--- 📂 Phase 1: Planning ---");
    let planApproved = false;
    let planFileFound = null;

    // Tailored planning logic per Use Case
    if (selection === "5") {
      // PR COMMENTS: Skip planning completely and fetch comments programmatically
      console.log(
        "PR Comments workflow selected. Fetching comments programmatically via script...",
      );
      try {
        const comments = await getPRComments();
        console.log("\n=== RETRIEVED PR COMMENTS ===");
        console.log(comments);
        console.log("==============================\n");

        const autoPlan = {
          title: `Address PR Comments: ${objective}`,
          objective: objective,
          scope_boundaries: {
            in_scope: [
              "Address review findings and conversation comments in the fetched list.",
            ],
            out_of_scope: [
              "Any changes unrelated to the requested PR review feedback.",
            ],
          },
          exit_criteria: [
            "All identified issues in the comments are resolved.",
            "Tests and linters pass cleanly.",
          ],
          implementation_tasks: [
            "Analyze review feedback.",
            "Surgically update corresponding source files.",
          ],
        };
        if (!(await savePlanFromJSON(autoPlan))) {
          console.error("❌ Failed to save plan file to disk.");
          process.exitCode = 1;
          return;
        }
        planApproved = true;
        console.log(
          "✅ Generated deterministic implementation plan from PR Comments context.",
        );
      } catch (err) {
        console.error(`❌ Failed to fetch PR comments: ${err?.message || err}`);
        process.exitCode = 1;
        return;
      }
    } else if (selection === "3") {
      // REFACTOR: Generate a deterministic refactor plan programmatically
      console.log(
        "Refactor workflow selected. Generating deterministic plan...",
      );
      const autoPlan = {
        title: `Refactor: ${objective}`,
        objective: objective,
        scope_boundaries: {
          in_scope: [
            "Refactor targeted components to improve readability or structure.",
            "Ensure exact same behavioral outputs.",
          ],
          out_of_scope: [
            "Changing public APIs.",
            "Adding new features or fixing unrelated bugs.",
          ],
        },
        exit_criteria: [
          "Behavioral coverage remains identical.",
          "All tests and linters pass cleanly with zero new warnings.",
        ],
        implementation_tasks: [
          "Examine existing target modules.",
          "Apply clean architecture and DRY principles surgically.",
        ],
      };
      try {
        if (!(await savePlanFromJSON(autoPlan))) {
          console.error("❌ Failed to save plan file to disk.");
          process.exitCode = 1;
          return;
        }
      } catch (err) {
        console.error(
          `❌ Failed to save plan file to disk: ${err?.message || err}`,
        );
        process.exitCode = 1;
        return;
      }
      planApproved = true;
      console.log("✅ Generated deterministic Refactor implementation plan.");
    } else if (selection === "4") {
      // TEST: Generate a deterministic testing plan programmatically
      console.log(
        "Add Tests workflow selected. Generating deterministic plan...",
      );
      const autoPlan = {
        title: `Add Tests: ${objective}`,
        objective: objective,
        scope_boundaries: {
          in_scope: [
            "Adding or expanding automated test files (e.g. *_test.go, *.test.*, *.spec.*, test_*.*, or files under test directories).",
          ],
          out_of_scope: [
            "Modifying any non-test product files (files outside of test files).",
          ],
        },
        exit_criteria: [
          "Test coverage is expanded successfully.",
          "New tests pass cleanly.",
          "No changes made to product code.",
        ],
        implementation_tasks: [
          "Analyze target code paths requiring test coverage.",
          "Write or expand unit / integration tests.",
        ],
      };
      try {
        if (!(await savePlanFromJSON(autoPlan))) {
          console.error("❌ Failed to save plan file to disk.");
          process.exitCode = 1;
          return;
        }
      } catch (err) {
        console.error(
          `❌ Failed to save plan file to disk: ${err?.message || err}`,
        );
        process.exitCode = 1;
        return;
      }
      planApproved = true;
      console.log("✅ Generated deterministic Testing implementation plan.");
    } else if (selection === "6") {
      // REVIEW LOOP: Preserve existing plans/current.md and proceed directly to review
      console.log(
        "Review Loop workflow selected. Preserving existing plans/current.md...",
      );
      planApproved = true;
    } else {
      // BUGFIX & FEATURE: Run standard agent planning with user-interview capabilities
      console.log(
        "Bugfix/Feature workflow selected. Invoking @planner agent...",
      );
      const planConfig = await loadAgentInstructions("planner");
      const planPrompt = `Objective: "${objective}".\nPlease conduct your user interview using the \`ask_user\` tool (intent = "clarification"). Once you have all the context you need, generate the implementation plan as a strict JSON block according to your system instructions. Do not write any files to disk yourself.

Your final JSON response must strictly conform to this schema:
\`\`\`json
{
  "title": "Clear implementation plan title",
  "objective": ${JSON.stringify(objective)},
  "scope_boundaries": {
    "in_scope": ["Specific files or components to touch"],
    "out_of_scope": ["Things not to touch or refactor"]
  },
  "exit_criteria": ["Testable definition of done criteria"],
  "implementation_tasks": ["Actionable step 1", "Actionable step 2"]
}
\`\`\``;

      const planSystemInstructions =
        planConfig.instructions ||
        `You are strictly in PLANNING phase (Phase 1). Do NOT modify any source files. You must formulate a complete plan JSON and output it directly.`;

      try {
        const planResult = await runAgentSession({
          initialPrompt: planPrompt,
          systemInstructions: planSystemInstructions,
          requestedModel: planConfig.model,
          blockTools: [
            "write_file",
            "replace",
            "create_file",
            "edit_file",
            "run_shell_command",
          ],
          customTools: [askUserTool],
          rl,
        });

        const validatedPlan = await validateAgentOutput(
          planResult,
          planSchema,
          planConfig?.model || models?.flash,
        );
        if (!validatedPlan) {
          console.error(
            "❌ Failed to parse plan JSON output from planner agent. Raw output:",
          );
          console.error(
            `\n--- RAW OUTPUT START ---\n${planResult}\n--- RAW OUTPUT END ---\n`,
          );
          process.exitCode = 1;
          return;
        }
        const success = await savePlanFromJSON(validatedPlan, objective);
        if (!success) {
          console.error("❌ Failed to save plan file to disk.");
          process.exitCode = 1;
          return;
        }
      } catch (err) {
        console.error(`❌ Plan generation failed: ${err?.message || err}`);
        process.exitCode = 1;
        return;
      }
    }

    // Display the plan and gate approval
    const possiblePlanPaths = [path.join(repoRoot, "plans/current.md")];
    for (const p of possiblePlanPaths) {
      if (await exists(p)) {
        planFileFound = p;
        break;
      }
    }

    if (planFileFound) {
      console.log("\n======================================");
      console.log(
        `📄 CURRENT PLAN (${path.relative(repoRoot, planFileFound)}):`,
      );
      console.log("======================================");
      const planContent = await fs.readFile(planFileFound, "utf8");
      console.log(planContent);
      console.log("======================================\n");
    } else {
      console.error(
        "❌ Error: plans/current.md not found. Gating fails-closed.",
      );
      process.exitCode = 1;
      return;
    }

    while (!planApproved) {
      const userApprove = await rl.question(
        "👉 Do you approve the plan? (yes/no): ",
      );
      if (userApprove.trim().toLowerCase() === "yes") {
        planApproved = true;
        console.log("\n✅ Plan approved by developer!");
      } else {
        console.log("❌ Plan rejected. Exiting orchestrator.");
        return;
      }
    }

    // Check if on main or master, sync, and branch out
    try {
      const switchBranchScript = path.join(
        repoRoot,
        ".gemini/skills/git/scripts/switch-branch.js",
      );
      const { stdout: branchOutput } = await execFileAsync(
        process.execPath,
        [switchBranchScript, "--current"],
        {
          cwd: repoRoot,
        },
      );
      const currentBranch = branchOutput.trim();
      if (PREFERRED_DEFAULT_BRANCHES.includes(currentBranch)) {
        console.log(
          `\n🌿 Currently on "${currentBranch}" branch. Syncing and creating a new working branch...`,
        );

        const prefixMap = {
          1: "bugfix",
          2: "feature",
          3: "refactor",
          4: "test",
          5: "pr-comments",
          6: "review-loop",
        };
        const prefix = prefixMap[selection] || "task";
        const sanitizedObjective =
          objective
            .toLowerCase()
            .replace(/[^a-z0-9]+/g, "-")
            .substring(0, 50)
            .replace(/^-+|-+$/g, "") || "working-branch";
        let newBranch = `${prefix}/${sanitizedObjective}`;

        // If the branch already exists, append a random ID
        try {
          await execFileAsync(
            process.execPath,
            [switchBranchScript, "--check", newBranch],
            { cwd: repoRoot },
          );
          const randomId = crypto.randomBytes(4).toString("hex");
          newBranch = `${newBranch}-${randomId}`;
        } catch (err) {
          // Branch does not exist, safe to use as is
          console.debug(
            `[DEBUG] Branch refs/heads/${newBranch} does not exist: ${err?.message || err}`,
          );
        }

        const createBranchScript = path.join(
          repoRoot,
          ".gemini/skills/git/scripts/create-branch.js",
        );
        await execFileAsync(
          process.execPath,
          [createBranchScript, newBranch, currentBranch],
          { cwd: repoRoot },
        );
        console.log(`✅ Checked out new branch: ${newBranch}`);
      }
    } catch (err) {
      console.error(`❌ Branch creation failed: ${err?.message || err}`);
      process.exitCode = 1;
      return;
    }

    // ==========================================
    // PHASE 2: IMPLEMENTATION (Write Access)
    // ==========================================
    console.log("\n--- 🔨 Phase 2: Implementation ---");

    if (selection === "6") {
      console.log(
        "ℹ️ Review Loop selected - skipping Phase 2 initial implementation.",
      );
    } else {
      console.log("Transitioning to Implementation phase with write tools...");

      // Load active plan content so the agent has direct access without wasting turns
      const activePlanPath =
        planFileFound || path.join(repoRoot, "plans/current.md");
      let activePlanContent = "";
      try {
        activePlanContent = await fs.readFile(activePlanPath, "utf8");
      } catch (err) {
        console.warn(
          `[WARN] Could not load active plan from ${activePlanPath}: ${err?.message || err}`,
        );
      }

      let implementInstructions = `You are in the IMPLEMENTATION phase (Phase 2). 
Implement the approved plan meticulously and surgically. The full active plan from 'plans/current.md' is already embedded below in <active_plan>. Do NOT waste turns reading it from disk.`;

      if (selection === "4") {
        // Strictly restrict the agent to test files
        implementInstructions += `\n⚠️ STRICT CONSTRAINT: You are ONLY allowed to write or modify test files (e.g. *_test.go, *.test.*, *.spec.*, test_*.*, or files under test directories). DO NOT modify any product files.`;
      }

      const implementPrompt = `Objective: "${objective}".
Please implement the approved plan documented in 'plans/current.md' meticulously and surgically.
The full active plan from 'plans/current.md' is embedded directly below. Do NOT execute tools or commands to read 'plans/current.md' from disk.
Once you have fully finished your implementation, stop.

<active_plan>
${activePlanContent}
</active_plan>`;

      try {
        await runAgentSession({
          initialPrompt: implementPrompt,
          systemInstructions: implementInstructions,
          requestedModel: models.pro,
          rl,
        });
        console.log("\n✅ Implementation session completed.");
      } catch (err) {
        console.error(`❌ Implementation failed: ${err?.message || err}`);
        process.exitCode = 1;
        return;
      }

      // Programmatic verification of scope boundaries
      if (selection === "4") {
        const passed = await verifyTestScopeBoundaries(repoRoot);
        if (!passed) {
          process.exitCode = 1;
          return;
        }
      }
    }

    // ==========================================
    // PHASE 3: AUTOMATED QA REVIEW
    // ==========================================
    console.log("\n--- 🛡️ Phase 3: Automated QA Review ---");

    // Stage changes so any new untracked or modified files from Phase 2 are included in the baseline diff
    try {
      const stageScript = path.join(
        repoRoot,
        ".gemini/skills/git/scripts/stage-files.js",
      );
      await execFileAsync(
        process.execPath,
        [stageScript, "--all", "--include-untracked"],
        {
          cwd: repoRoot,
        },
      );
    } catch (err) {
      const errText =
        typeof err === "string"
          ? err
          : `${err?.message || ""} ${err?.stderr || ""}`;
      if (
        errText.includes("No files specified to stage") ||
        errText.includes("No files to stage")
      ) {
        console.debug(
          "[DEBUG] Working tree already clean or staged before QA review loop.",
        );
      } else {
        console.error(
          `❌ Failed closed: Could not stage files before QA review loop: ${errText.trim() || "Unknown error"}`,
        );
        process.exitCode = 1;
        return;
      }
    }

    const planPath = planFileFound || path.join(repoRoot, "plans/current.md");
    console.log(
      "\n🔄 Initiating Multi-Agent Ralph Review Loop (Thinking Reviewer -> Flash-Lite Tally -> Pro/Flash Implementer)...",
    );

    const loopResult = await runReviewRalphLoop({
      repoRoot,
      planPath,
      requireInitialDiff: selection === "6",
      proModel: models.pro,
      flashModel: models.flash,
      flashLiteModel: models.flash_lite,
      proReviewIterations: 3,
      proImplementationIterations: 1,
      implementerInstructions:
        selection === "4"
          ? `${DEFAULT_IMPLEMENTER_INSTRUCTIONS}\n⚠️ STRICT CONSTRAINT: You are ONLY allowed to write or modify test files (e.g. *_test.go, *.test.*, *.spec.*, test_*.*, or files under test directories). DO NOT modify any product files.`
          : DEFAULT_IMPLEMENTER_INSTRUCTIONS,
      sessionRunner: runAgentSession,
      validator: (args) =>
        runQAPipeline({
          ...args,
          fix: false,
          commands: [
            "node .gemini/skills/code-validation/scripts/lint.js all",
            "node .gemini/skills/code-validation/scripts/test.js all",
          ],
        }),
      validationCommands: [
        "node .gemini/skills/code-validation/scripts/lint.js all",
        "node .gemini/skills/code-validation/scripts/test.js all",
      ],
      rl,
      maxIterations: 15,
      maxTurnsPerSession: 20,
    });

    if (!loopResult.completed || !loopResult.success) {
      console.error(
        `\n❌ Ralph Review Loop finished without full convergence: ${loopResult.reason}`,
      );
      process.exitCode = 1;
      return;
    }

    console.log(
      `\n🟢 Ralph Review Loop completed successfully (${loopResult.iterations} iteration(s)).`,
    );

    // For test use case, re-verify scope boundaries after review loop implementations
    if (selection === "4") {
      const passed = await verifyTestScopeBoundaries(repoRoot);
      if (!passed) {
        process.exitCode = 1;
        return;
      }
    }

    if (selection === "6") {
      console.log("\n--- 📦 Final Review & Summary ---");
      try {
        await execFileAsync(
          process.execPath,
          [
            path.join(repoRoot, ".gemini/skills/git/scripts/stage-files.js"),
            "--all",
            "--include-untracked",
          ],
          {
            cwd: repoRoot,
          },
        );
      } catch (err) {
        const errText =
          typeof err === "string"
            ? err
            : `${err?.message || ""} ${err?.stderr || ""}`;
        if (
          errText.includes("No files specified to stage") ||
          errText.includes("No files to stage")
        ) {
          console.debug(
            "[DEBUG] Working tree already clean or staged during review loop Phase 3.",
          );
        } else {
          console.error(
            `❌ Failed closed: Could not stage files in review loop Phase 3: ${errText.trim() || "Unknown error"}`,
          );
          process.exitCode = 1;
          return;
        }
      }
      const diffText = await getGitDiff();
      if (!diffText || diffText.trim() === "") {
        console.log("✅ No changes were made in the review loop.");
        return;
      }

      console.log(
        "Asking Gemini to generate a summary, commit message, and PR description...",
      );
      const summaryPrompt = `Based on the following git diff, please generate a final summary of all changes made, a proposed single-line conventional commit message, and a PR description. Do NOT make any further code changes. Output the response clearly formatted:\n\n${diffText}`;

      const summaryMsg = await runAgentSession({
        initialPrompt: summaryPrompt,
        systemInstructions:
          "You are a technical reviewer. Summarize the changes provided without modifying any files or calling tools.",
        requestedModel: models.flash || models.pro,
        maxTurnsOverride: 1,
        requireTopic: false,
        isolate: true,
        rl,
      });

      console.log("\n======================================");
      console.log("📝 FINAL REVIEW SUMMARY:");
      console.log("======================================");
      console.log(summaryMsg);
      console.log("======================================\n");
      console.log(
        "✅ Review loop complete. Please review the changes, commit, and push manually.",
      );
      return;
    }

    // ==========================================
    // PHASE 4: FINAL USER REVIEW & COMMIT
    // ==========================================
    console.log("\n--- 📦 Phase 4: Final User Review & Commit ---");

    const diffText = await getGitDiff();
    if (!diffText || diffText.trim() === "") {
      console.log(
        "🟢 QA passed but no changes were detected in Git. Session complete!",
      );
      return;
    }

    const initialDiffHash = crypto
      .createHash("sha256")
      .update(stripDiffMetadata(diffText))
      .digest("hex");

    console.log("\n======================================");
    console.log("🔍 PROPOSED CHANGES (git diff):");
    console.log("======================================");
    console.log(diffText);
    console.log("======================================\n");
    console.log(`🔐 Cryptographic Diff Hash (SHA-256): ${initialDiffHash}\n`);

    const finalApproval = await rl.question(
      "👉 Do you approve these changes for commit? (yes/no): ",
    );
    if (finalApproval.trim().toLowerCase() !== "yes") {
      console.log(
        "❌ Commit cancelled by developer. Sourcing/modifications are left in working directory.",
      );
      return;
    }

    // Cryptographic Verification: Calculate the hash again right before staging
    const preStagingDiff = await getGitDiff();
    const preStagingDiffHash = crypto
      .createHash("sha256")
      .update(stripDiffMetadata(preStagingDiff))
      .digest("hex");

    if (preStagingDiffHash !== initialDiffHash) {
      console.error(
        "❌ Cryptographic Verification Error: Active diff hash changed after user approval!",
      );
      process.exitCode = 1;
      return;
    }
    console.log("🔒 Cryptographic diff hash successfully verified!");

    console.log("Asking Gemini to generate a commit message candidate...");
    const commitPrompt = `Generate a single-line, highly descriptive and concise git commit message conforming to Conventional Commits format (e.g., "feat: add feature X" or "fix: resolve bug Y") based strictly on this diff:
\n\n${diffText}`;

    const defaultMsgRaw = await runAgentSession({
      initialPrompt: commitPrompt,
      systemInstructions:
        "You are a professional software engineer. Generate a single-line conventional commit message with NO quotes, markdown, or preambles.",
      requestedModel: models.flash_lite || models.flash,
      maxTurnsOverride: 1,
      requireTopic: false,
      isolate: true,
      rl,
    });

    const defaultMsg = defaultMsgRaw.trim().replace(/^['"`]+|['"`]+$/g, "");

    let commitApproved = false;
    let finalMsg = "";

    while (!commitApproved) {
      console.log(`\nProposed commit message: "${defaultMsg}"`);
      const userMsg = await rl.question(
        "Enter commit message (or press enter for default): ",
      );
      const candidateMsg = userMsg.trim() || defaultMsg;

      if (!candidateMsg) {
        console.error(
          "❌ Error: Commit message is empty. Please enter a valid commit message.",
        );
        continue;
      }

      if (candidateMsg.length > 100) {
        console.error(
          `❌ Error: Commit message should be less than 100 characters (currently ${candidateMsg.length}).`,
        );
        continue;
      }

      // Programmatically validate message
      const validation = await validateMessage(candidateMsg);
      if (!validation.valid) {
        console.error(`❌ Error: ${validation.reason}`);
        console.log(
          "Please adjust your commit type or description according to conventional commit guidelines.",
        );
        continue;
      }

      finalMsg = candidateMsg;
      commitApproved = true;
    }

    console.log(`Committing: ${finalMsg}`);
    try {
      await stageAndCommit(finalMsg);
      console.log("🟢 Changes committed successfully!");
    } catch (err) {
      console.error(
        `❌ Commit failed:\n${err?.message || err}\n${err?.stderr || ""}`.trimEnd(),
      );
      process.exitCode = 1;
      return;
    }
  }
}

main()
  .catch((err) => {
    if (isShuttingDown) {
      return;
    }
    isShuttingDown = true;
    console.error(
      "❌ Fatal Orchestrator Error:",
      err?.stack || err?.message || err,
    );
    process.exitCode = 1;
  })
  .finally(async () => {
    try {
      if (rl) {
        rl.close();
      }
    } catch (closeErr) {
      console.debug(
        `[DEBUG] Error closing readline interface: ${closeErr?.message || closeErr}`,
      );
    }
    await flushLogs();
  });
