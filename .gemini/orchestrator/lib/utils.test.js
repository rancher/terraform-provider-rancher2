import { z } from "@google/gemini-cli-sdk";
import assert from "node:assert";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import {
  formatToolAction,
  getRepoRoot,
  getSkillsFingerprint,
  loadWorkspaceSkillsContext,
  normalizePlanObject,
  parseJSONFromText,
  renderBox,
  renderProseCard,
  renderToolCard,
  renderTopicBanner,
  resetSkillsContextCache,
  savePlanFromJSON,
  stripAnsi,
  validateAgentOutput,
  wrapLine,
} from "./utils.js";

test("parseJSONFromText resilient extraction", async (t) => {
  await t.test(
    "extracts the last valid JSON block when followed by non-JSON code blocks",
    () => {
      const sample = [
        "Initial explanation.",
        "```json",
        '{"approval_status": "APPROVED", "findings": []}',
        "```",
        "Here is a command to verify:",
        "```",
        "go test ./...",
        "```",
      ].join("\n");

      const parsed = parseJSONFromText(sample, "qa");
      assert.ok(parsed);
      assert.strictEqual(parsed.approval_status, "APPROVED");
      assert.deepStrictEqual(parsed.findings, []);
    },
  );

  await t.test(
    "fails closed to UNAPPROVED when no valid JSON is present in QA output",
    () => {
      const text = "The changes are APPROVED and ready to ship.";
      const parsed = parseJSONFromText(text, "qa");
      assert.ok(parsed);
      assert.strictEqual(parsed.approval_status, "UNAPPROVED");
      assert.strictEqual(parsed.findings.length, 1);

      // Empty string fallback
      const emptyParsed = parseJSONFromText("", "qa");
      assert.ok(emptyParsed);
      assert.strictEqual(emptyParsed.approval_status, "UNAPPROVED");
      assert.strictEqual(emptyParsed.findings[0].file, "unknown");

      // Non-string inputs (null, undefined, number) fallback
      const nullParsed = parseJSONFromText(null, "qa");
      assert.ok(nullParsed);
      assert.strictEqual(nullParsed.approval_status, "UNAPPROVED");

      const undefParsed = parseJSONFromText(undefined, "qa");
      assert.ok(undefParsed);
      assert.strictEqual(undefParsed.approval_status, "UNAPPROVED");

      // Non-qa fallbackType returns null for invalid input
      assert.strictEqual(parseJSONFromText("", "plan"), null);
      assert.strictEqual(parseJSONFromText(null, "plan"), null);
    },
  );

  await t.test(
    "extracts raw JSON when unadorned by code blocks but surrounded by commentary",
    () => {
      const text =
        'Here is the QA review result: {"approval_status": "APPROVED", "findings": []} and additional notes.';
      const parsed = parseJSONFromText(text, "qa");
      assert.ok(parsed);
      assert.strictEqual(parsed.approval_status, "APPROVED");
      assert.deepStrictEqual(parsed.findings, []);
    },
  );

  await t.test(
    "extracts latest JSON when multiple raw objects exist without code blocks",
    () => {
      const text =
        'Draft: {"draft": true}. Final QA review: {"approval_status": "APPROVED", "findings": []}';
      const parsed = parseJSONFromText(text, "qa");
      assert.ok(parsed);
      assert.strictEqual(parsed.approval_status, "APPROVED");
      assert.deepStrictEqual(parsed.findings, []);
    },
  );

  await t.test(
    "extracts unadorned JSON containing nested objects and arrays",
    () => {
      const text =
        'Draft: {"draft": {"a": 1}}. Final: {"approval_status": "UNAPPROVED", "findings": [{"file": "foo.go", "line_numbers": [12]}]}';
      const parsed = parseJSONFromText(text, "qa");
      assert.ok(parsed);
      assert.strictEqual(parsed.approval_status, "UNAPPROVED");
      assert.strictEqual(parsed.findings.length, 1);
      assert.strictEqual(parsed.findings[0].file, "foo.go");
    },
  );

  await t.test(
    "extracts top-level JSON when text contains intermingled array and object brackets",
    () => {
      const text =
        'prefix [1, 2] middle {"plan": [{"file": "main.go", "steps": ["step 1"]}]} suffix';
      const parsed = parseJSONFromText(text, "plan");
      assert.ok(parsed);
      assert.strictEqual(parsed.title, "[main.go] step 1");
    },
  );

  await t.test(
    "ignores trailing bracketed markdown citations when extracting plan object",
    () => {
      const text =
        'Here is the plan: {"plan": [{"file": "main.go", "title": "Implement feature"}]}. See [1] and [2].';
      const parsed = parseJSONFromText(text, "plan");
      assert.ok(parsed);
      assert.strictEqual(parsed.title, "Implement feature");
      assert.deepStrictEqual(parsed.implementation_tasks, [
        "[main.go] Implement feature",
      ]);
    },
  );

  await t.test(
    "normalizes missing approval_status when findings array is provided in qa mode",
    () => {
      const approvedText = '{"findings": []}';
      const parsedApproved = parseJSONFromText(approvedText, "qa");
      assert.ok(parsedApproved);
      assert.strictEqual(parsedApproved.approval_status, "APPROVED");
      assert.deepStrictEqual(parsedApproved.findings, []);

      const unapprovedText =
        '{"findings": [{"file": "foo.go", "line_numbers": [12]}]}';
      const parsedUnapproved = parseJSONFromText(unapprovedText, "qa");
      assert.ok(parsedUnapproved);
      assert.strictEqual(parsedUnapproved.approval_status, "UNAPPROVED");
      assert.strictEqual(parsedUnapproved.findings.length, 1);
      assert.strictEqual(parsedUnapproved.findings[0].file, "foo.go");
    },
  );

  await t.test(
    "normalizes planner output containing plan array into canonical plan schema",
    () => {
      const crashPayload = `\`\`\`json
{
  "plan": [
    {
      "file": "agent-scripts/lib/agent-runner.js",
      "description": "Enhance the tool execution logging to include key argument details for better visibility.",
      "steps": [
        "Locate the 'tool_call_request' handling block around line 609.",
        "Add logic to extract contextual information from args.",
        "Update the console.log statement."
      ]
    }
  ]
}
\`\`\``;

      const parsed = parseJSONFromText(crashPayload, "plan");
      assert.ok(parsed);
      assert.strictEqual(
        parsed.title,
        "Enhance the tool execution logging to include key argument details for better visibility.",
      );
      assert.strictEqual(
        parsed.objective,
        "Enhance the tool execution logging to include key argument details for better visibility.",
      );
      assert.deepStrictEqual(parsed.scope_boundaries.in_scope, [
        "Modify agent-scripts/lib/agent-runner.js",
      ]);
      assert.ok(parsed.exit_criteria.length > 0);
      assert.strictEqual(parsed.implementation_tasks.length, 3);
      assert.strictEqual(
        parsed.implementation_tasks[0],
        "[agent-scripts/lib/agent-runner.js] Locate the 'tool_call_request' handling block around line 609.",
      );
    },
  );

  await t.test("normalizes raw JSON array into canonical plan schema", () => {
    const rawArray = `[
      {"file": "main.go", "description": "Update main entrypoint"},
      {"file": "main_test.go", "description": "Add test coverage"}
    ]`;

    const parsed = parseJSONFromText(rawArray, "plan");
    assert.ok(parsed);
    assert.strictEqual(parsed.title, "Update main entrypoint");
    assert.deepStrictEqual(parsed.scope_boundaries.in_scope, [
      "Modify main.go",
      "Modify main_test.go",
    ]);
    assert.strictEqual(parsed.implementation_tasks.length, 2);
    assert.strictEqual(
      parsed.implementation_tasks[0],
      "[main.go] Update main entrypoint",
    );
  });

  await t.test("normalizes plan with alternate task and step keys", () => {
    const alternate = JSON.stringify({
      title: "Custom Title",
      tasks: ["Task A", "Task B"],
    });

    const parsed = parseJSONFromText(alternate, "plan");
    assert.ok(parsed);
    assert.strictEqual(parsed.title, "Custom Title");
    assert.deepStrictEqual(parsed.implementation_tasks, ["Task A", "Task B"]);
  });

  await t.test(
    "normalizePlanObject fallback values when fields are empty",
    () => {
      const emptyObj = normalizePlanObject({}, "Fallback Objective");
      assert.strictEqual(emptyObj.title, "Fallback Objective");
      assert.strictEqual(emptyObj.objective, "Fallback Objective");
      assert.deepStrictEqual(emptyObj.scope_boundaries.in_scope, [
        "Implement requested tasks as planned",
      ]);
      assert.strictEqual(emptyObj.implementation_tasks.length, 1);
    },
  );
});

test("savePlanFromJSON resilient saving", async (t) => {
  const repoRoot = await getRepoRoot();
  const planFile = path.join(repoRoot, "plans/current.md");
  let originalContent = null;
  try {
    originalContent = await fs.readFile(planFile, "utf8");
  } catch {
    // File did not exist prior to test
  }

  t.after(async () => {
    if (originalContent !== null) {
      await fs.writeFile(planFile, originalContent, "utf8");
    } else {
      try {
        await fs.unlink(planFile);
      } catch {
        // Ignore if already deleted
      }
    }
  });

  await t.test(
    "successfully saves normalized plan to plans/current.md",
    async () => {
      const crashPayload = `\`\`\`json
{
  "plan": [
    {
      "file": "agent-scripts/lib/agent-runner.js",
      "description": "Enhance tool execution logging",
      "steps": [
        "Step 1: Extract tool arguments",
        "Step 2: Log enriched tool details"
      ]
    }
  ]
}
\`\`\``;

      const saved = await savePlanFromJSON(crashPayload, "Fallback Title");
      assert.strictEqual(saved, true);

      const content = await fs.readFile(planFile, "utf8");

      assert.ok(
        content.includes(
          "# IMPLEMENTATION PLAN: Enhance tool execution logging",
        ),
      );
      assert.ok(content.includes("## 🎯 OBJECTIVE"));
      assert.ok(content.includes("**IN SCOPE:**"));
      assert.ok(content.includes("- Modify agent-scripts/lib/agent-runner.js"));
      assert.ok(content.includes("## 🛠️ IMPLEMENTATION TASKS"));
      assert.ok(
        content.includes(
          "- [ ] [agent-scripts/lib/agent-runner.js] Step 1: Extract tool arguments",
        ),
      );
    },
  );

  await t.test("returns false when text contains no valid JSON", async () => {
    const saved = await savePlanFromJSON("This is not JSON at all.");
    assert.strictEqual(saved, false);
  });

  await t.test(
    "successfully saves when passed a plan object directly",
    async () => {
      const planObj = {
        title: "Direct Object Plan",
        implementation_tasks: ["Task 1"],
      };
      const saved = await savePlanFromJSON(planObj);
      assert.strictEqual(saved, true);
    },
  );
});

test("validateAgentOutput schema validation and error handling", async (t) => {
  const schema = z.object({
    name: z.string(),
    count: z.number(),
  });

  await t.test(
    "successfully validates and returns parsed object matching schema",
    async () => {
      const raw = '```json\n{"name": "test-agent", "count": 42}\n```';
      const result = await validateAgentOutput(raw, schema);
      assert.deepStrictEqual(result, { name: "test-agent", count: 42 });
    },
  );

  await t.test(
    "successfully validates when passed a JavaScript object directly",
    async () => {
      const obj = { name: "direct-object", count: 7 };
      const result = await validateAgentOutput(obj, schema);
      assert.deepStrictEqual(result, { name: "direct-object", count: 7 });
    },
  );

  await t.test(
    "re-formats and recovers successfully when session runner returns valid schema JSON",
    async () => {
      const invalidRaw = "This is invalid text";
      let calls = 0;
      const mockRunner = async () => {
        calls++;
        return '```json\n{"name": "recovered-agent", "count": 99}\n```';
      };
      const result = await validateAgentOutput(invalidRaw, schema, mockRunner);
      assert.strictEqual(calls, 1);
      assert.deepStrictEqual(result, { name: "recovered-agent", count: 99 });
    },
  );

  await t.test(
    "gracefully returns null when validation fails and session errors",
    async () => {
      const invalidRaw = "This is completely invalid and not JSON";
      const mockRunner = async () => {
        throw new Error("SDK session failed");
      };
      const result = await validateAgentOutput(invalidRaw, schema, mockRunner);
      assert.strictEqual(result, null);
    },
  );

  await t.test(
    "gracefully returns null after initial parse and 4 validator agent retries fail",
    async () => {
      const invalidRaw = "This is completely invalid and not JSON";
      let calls = 0;
      const mockRunner = async () => {
        calls++;
        return "Still invalid";
      };
      const result = await validateAgentOutput(invalidRaw, schema, mockRunner);
      assert.strictEqual(calls, 4); // Initial try failed, then 4 retries via validator agent
      assert.strictEqual(result, null);
    },
  );

  await t.test(
    "passes isolate: true and maxTurns: 5 options to the validator session runner",
    async () => {
      let receivedOptions = null;
      const mockRunner = async (opts) => {
        receivedOptions = opts;
        return '```json\n{"name": "isolated-agent", "count": 10}\n```';
      };
      const result = await validateAgentOutput(
        "invalid text",
        schema,
        mockRunner,
      );
      assert.strictEqual(receivedOptions.isolate, true);
      assert.strictEqual(receivedOptions.maxTurns, 5);
      assert.deepStrictEqual(result, { name: "isolated-agent", count: 10 });
    },
  );

  await t.test(
    "handles nested schemas and unwraps optional inner types during retry",
    async () => {
      const nestedSchema = z.object({
        meta: z.object({
          tag: z.string().optional(),
        }),
      });
      const mockRunner = async (opts) => {
        assert.ok(opts.initialPrompt.includes("Expected Schema Structure:"));
        return '```json\n{"meta": {"tag": "beta"}}\n```';
      };
      const result = await validateAgentOutput(
        "not json",
        nestedSchema,
        mockRunner,
      );
      assert.deepStrictEqual(result, { meta: { tag: "beta" } });
    },
  );

  await t.test(
    "generates expected schema structure for top-level array schemas during retry",
    async () => {
      const arraySchema = z.array(z.string());
      const mockRunner = async (opts) => {
        assert.ok(opts.initialPrompt.includes("Expected Schema Structure:"));
        return '```json\n["item1", "item2"]\n```';
      };
      const result = await validateAgentOutput(
        "invalid",
        arraySchema,
        mockRunner,
      );
      assert.deepStrictEqual(result, ["item1", "item2"]);
    },
  );

  await t.test(
    "throws TypeError when schema is missing or invalid",
    async () => {
      await assert.rejects(async () => validateAgentOutput("data", null), {
        name: "TypeError",
        message: /valid Zod schema/,
      });
    },
  );

  await t.test(
    "introspects ZodLiteral and ZodUnion schemas in retry prompts",
    async () => {
      const unionSchema = z.object({
        status: z.union([z.literal("APPROVED"), z.literal("REJECTED")]),
      });
      const mockRunner = async (opts) => {
        assert.ok(
          opts.initialPrompt.includes('"APPROVED" | "REJECTED"') ||
            opts.initialPrompt.includes('\\"APPROVED\\" | \\"REJECTED\\"'),
        );
        return '```json\n{"status": "APPROVED"}\n```';
      };
      const result = await validateAgentOutput(
        "invalid",
        unionSchema,
        mockRunner,
      );
      assert.deepStrictEqual(result, { status: "APPROVED" });
    },
  );
});

test("renderBox and cli formatting utilities", async (t) => {
  await t.test("stripAnsi removes ANSI escape codes from styled text", () => {
    assert.strictEqual(
      stripAnsi("\u001b[32m" + "hello" + "\u001b[0m world"),
      "hello world",
    );
    assert.strictEqual(stripAnsi("plain text"), "plain text");
    assert.strictEqual(stripAnsi(123), "123");
    // OSC sequences and terminal control characters (F3)
    assert.strictEqual(
      stripAnsi("\x1b]0;Spoofed Title\x07clean text"),
      "clean text",
    );
    assert.strictEqual(
      stripAnsi("\x1b]0;Spoofed Title\x1b\\clean text"),
      "clean text",
    );
    assert.strictEqual(
      stripAnsi("\x1b[2J\x1b[H\x1b[?25hclean text"),
      "clean text",
    );
    assert.strictEqual(
      stripAnsi("text\x00with\x07bell\x08and\x1bcontrol"),
      "textwithbellandcontrol",
    );
    // C1 8-bit controls and CSI/OSC sequences (U+0080-U+009F)
    assert.strictEqual(stripAnsi("\u009b2J\u009bHclean text"), "clean text");
    assert.strictEqual(
      stripAnsi("\u009d0;Spoofed\u009cclean text"),
      "clean text",
    );
    assert.strictEqual(stripAnsi("text\u0080with\u009fc1"), "textwithc1");
  });

  await t.test(
    "wrapLine wraps long strings while preserving indentation",
    () => {
      const line =
        "  this is a long line that needs to wrap into multiple parts";
      const wrapped = wrapLine(line, 25);
      assert.ok(wrapped.length > 1);
      assert.ok(wrapped[0].startsWith("  "));
      assert.ok(wrapped[1].startsWith("  "));
    },
  );

  await t.test(
    "wrapLine splits oversized tokens that exceed maxWidth into width-limited chunks (F3)",
    () => {
      const url = "https://github.com/rancher/terraform-provider-file/pull/444";
      const wrapped = wrapLine(url, 20);
      assert.ok(wrapped.length > 1);
      for (const chunk of wrapped) {
        assert.ok(
          chunk.length <= 20,
          `Chunk "${chunk}" must not exceed maxWidth 20`,
        );
      }

      const indentedUrl = "  prefix " + url;
      const wrappedIndented = wrapLine(indentedUrl, 25);
      for (const chunk of wrappedIndented) {
        assert.ok(
          chunk.length <= 25,
          `Chunk "${chunk}" must not exceed maxWidth 25`,
        );
        assert.ok(
          chunk.startsWith("  "),
          "Each chunk must preserve the line indentation",
        );
      }
    },
  );

  await t.test("renderBox leaves single-line text untouched", () => {
    const single = "agent used run_shell_command";
    assert.strictEqual(renderBox(single), single);
  });

  await t.test("renderBox wraps multi-line text in unicode box borders", () => {
    const multiline = "line 1\nline 2 longer";
    const boxed = renderBox(multiline);
    assert.ok(boxed.startsWith("┌─"));
    assert.ok(boxed.endsWith("─┘"));
    assert.ok(boxed.includes("│ line 1"));
    assert.ok(boxed.includes("│ line 2 longer │"));
  });

  await t.test(
    "renderBox wraps oversized tokens without exceeding maxWidth (F3)",
    () => {
      const multilineWithOversized =
        "Header\nhttps://github.com/rancher/terraform-provider-file/pull/444/very/long/path/argument";
      const boxed = renderBox(multilineWithOversized, 46);
      const lines = boxed.split("\n");
      for (const line of lines) {
        assert.ok(
          line.length <= 46,
          `Box line "${line}" (len ${line.length}) must not exceed maxWidth 46`,
        );
      }
    },
  );

  await t.test(
    "wrapLine safely handles non-positive maxWidth and non-string inputs",
    () => {
      assert.deepStrictEqual(wrapLine(null, 10), [""]);
      assert.deepStrictEqual(wrapLine(undefined, 10), [""]);
      const wrappedZero = wrapLine("abc def", 0);
      assert.ok(wrappedZero.length > 0);
      const wrappedNegative = wrapLine("abc def", -5);
      assert.ok(wrappedNegative.length > 0);
    },
  );

  await t.test(
    "wrapLine does not exceed safeMaxWidth when indentation exceeds maxWidth",
    () => {
      const indented = "        hello world";
      const wrapped = wrapLine(indented, 4);
      for (const chunk of wrapped) {
        assert.ok(
          chunk.length <= 4,
          `Chunk "${chunk}" must not exceed maxWidth 4`,
        );
      }
    },
  );

  await t.test("renderBox safely handles non-string and falsy inputs", () => {
    assert.strictEqual(renderBox(null), "");
    assert.strictEqual(renderBox(undefined), "");
    assert.strictEqual(renderBox(123), "123");
  });

  await t.test("renderBox supports round style borders", () => {
    const multiline = "line A\nline B";
    const boxed = renderBox(multiline, 40, "round");
    assert.ok(boxed.startsWith("╭─"));
    assert.ok(boxed.endsWith("─╯"));
    assert.ok(boxed.includes("│ line A"));
    assert.ok(boxed.includes("│ line B"));
  });

  await t.test(
    "formatToolAction formats known tools and generic fallbacks cleanly",
    () => {
      assert.strictEqual(
        formatToolAction("run_shell_command", { command: "git status" }),
        "Shell git status",
      );
      assert.strictEqual(
        formatToolAction("read_file", { file_path: "README.md" }),
        "ReadFile README.md",
      );
      assert.strictEqual(
        formatToolAction("write_file", { file_path: "foo.txt" }),
        "WriteFile foo.txt",
      );
      assert.strictEqual(
        formatToolAction("replace", { file_path: "bar.js" }),
        "Replace bar.js",
      );
      assert.strictEqual(
        formatToolAction("grep_search", { pattern: "test_func" }),
        "GrepSearch 'test_func'",
      );
      assert.strictEqual(
        formatToolAction("GrepSearch", { pattern: "test_func" }),
        "GrepSearch 'test_func'",
      );
      assert.strictEqual(
        formatToolAction("glob", { pattern: "*.go" }),
        "Glob '*.go'",
      );
      assert.strictEqual(
        formatToolAction("list_directory", { dir_path: "src" }),
        "ListDirectory src",
      );
      assert.strictEqual(
        formatToolAction("list_directory", {}),
        "ListDirectory .",
      );
      assert.strictEqual(formatToolAction("ask_user"), "AskUser");
      assert.strictEqual(
        formatToolAction("update_topic", { title: "New Topic" }),
        "UpdateTopic New Topic",
      );
      assert.strictEqual(
        formatToolAction("custom_tool", { param: "value" }),
        "CustomTool value",
      );
      assert.strictEqual(formatToolAction("no_op"), "NoOp");
      assert.strictEqual(formatToolAction(null), "Tool");
    },
  );

  await t.test(
    "renderToolCard formats rounded box with header and execution result",
    () => {
      const card = renderToolCard(
        {
          toolName: "run_shell_command",
          args: { command: "git status" },
          status: "success",
          result: "On branch main\nnothing to commit",
        },
        80,
      );
      const lines = card.split("\n");
      assert.ok(lines[0].startsWith("╭─"));
      assert.ok(lines[0].endsWith("─╮"));
      assert.strictEqual(lines[0].length, 80);
      assert.ok(lines[1].includes("✓  Shell git status"));
      assert.strictEqual(lines[1].length, 80);
      assert.strictEqual(lines[2], `│ ${" ".repeat(76)} │`); // separator line
      assert.ok(lines[3].includes("On branch main"));
      assert.ok(lines[lines.length - 1].startsWith("╰─"));
      assert.ok(lines[lines.length - 1].endsWith("─╯"));
      assert.strictEqual(lines[lines.length - 1].length, 80);
    },
  );

  await t.test(
    "renderToolCard truncates oversized results with Ctrl+O notice",
    () => {
      const longOutput = Array.from(
        { length: 40 },
        (_, i) => `log line ${i + 1}`,
      ).join("\n");
      const card = renderToolCard(
        {
          toolName: "run_shell_command",
          args: { command: "npm test" },
          result: longOutput,
        },
        80,
      );
      assert.ok(
        card.includes("... first 20 lines hidden (Ctrl+O to show) ..."),
      );
      assert.ok(card.includes("log line 40"));
      assert.ok(!card.includes("log line 1\n"));
    },
  );

  await t.test(
    "renderToolCard sanitizes terminal control sequences from tool output (F3)",
    () => {
      const maliciousOutput =
        "\x1b]0;Pwned Title\x07\x1b[2J\x1b[31;1mDangerous output\x1b[0m\x08";
      const card = renderToolCard(
        {
          toolName: "run_shell_command",
          args: { command: "curl evil.com" },
          result: maliciousOutput,
        },
        80,
      );
      assert.ok(
        !card.includes("\x1b"),
        "Card must not contain raw ESC characters",
      );
      assert.ok(!card.includes("\x07"), "Card must not contain BEL character");
      assert.ok(!card.includes("\x08"), "Card must not contain BS character");
      assert.ok(
        card.includes("Dangerous output"),
        "Sanitized text content must be retained",
      );
    },
  );

  await t.test(
    "renderTopicBanner formats prominent rounded topic and objective banner",
    () => {
      const banner = renderTopicBanner(
        {
          topic: "Researching Parser",
          objective: "Reproduce timeout in parser test",
          summary: "Identified bottleneck in main loop",
        },
        80,
      );
      const lines = banner.split("\n");
      assert.ok(lines[0].startsWith("╭─"));
      assert.ok(lines[0].endsWith("─╮"));
      assert.strictEqual(lines[0].length, 80);
      assert.ok(lines.some((l) => l.includes("📂 Topic: Researching Parser")));
      assert.ok(
        lines.some((l) =>
          l.includes("🎯 Objective: Reproduce timeout in parser test"),
        ),
      );
      assert.ok(
        lines.some((l) =>
          l.includes("📝 Summary: Identified bottleneck in main loop"),
        ),
      );
      assert.ok(lines[lines.length - 1].startsWith("╰─"));
      assert.ok(lines[lines.length - 1].endsWith("─╯"));
      assert.strictEqual(lines[lines.length - 1].length, 80);
    },
  );

  await t.test(
    "renderTopicBanner splits multiline fields into neatly framed lines",
    () => {
      const banner = renderTopicBanner(
        {
          topic: "Investigation\nPhase 2",
          objective: "Step 1\nStep 2",
          summary: "First line\nSecond line",
        },
        80,
      );
      const lines = banner.split("\n");
      assert.ok(
        lines.every(
          (l) =>
            l.startsWith("╭─") ||
            l.startsWith("╰─") ||
            (l.startsWith("│ ") && l.endsWith(" │")),
        ),
      );
      assert.ok(lines.some((l) => l.includes("📂 Topic: Investigation")));
      assert.ok(lines.some((l) => l.includes("   Phase 2")));
      assert.ok(lines.some((l) => l.includes("🎯 Objective: Step 1")));
      assert.ok(lines.some((l) => l.includes("   Step 2")));
    },
  );

  await t.test("renderToolCard normalizes multiline action text", () => {
    const card = renderToolCard(
      {
        toolName: "run_shell_command",
        args: { command: 'echo "line1\nline2"' },
      },
      80,
    );
    const lines = card.split("\n");
    assert.ok(
      lines.every(
        (l) =>
          l.startsWith("╭─") ||
          l.startsWith("╰─") ||
          (l.startsWith("│ ") && l.endsWith(" │")),
      ),
    );
    assert.ok(lines[1].includes('Shell echo "line1 line2"'));
  });

  await t.test(
    "renderTopicBanner normalizes bare carriage returns to prevent terminal spoofing",
    () => {
      const banner = renderTopicBanner(
        {
          topic: "Investigation\rSpoofed Topic",
          objective: "Step 1\r\nStep 2\rSpoofed Objective",
          summary: "Summary\rSpoofed Summary",
        },
        80,
      );
      assert.ok(
        !banner.includes("\r"),
        "Banner must not contain any bare carriage return characters",
      );
      assert.ok(banner.includes("Spoofed Topic"));
      assert.ok(banner.includes("Step 1"));
    },
  );

  await t.test("renderToolCard truncates massive single-line results", () => {
    const card = renderToolCard(
      {
        toolName: "run_shell_command",
        args: { command: "cat huge.txt" },
        result: "x".repeat(25000),
      },
      80,
    );
    assert.ok(card.includes("[truncated remaining 5000 characters]"));
  });

  await t.test(
    "renderProseCard formats multiline prose in rounded box matching tool card width",
    () => {
      const card = renderProseCard(
        "Paragraph 1\n\nParagraph 2 with more details",
        80,
      );
      assert.ok(card.startsWith("╭─"));
      assert.ok(card.endsWith("─╯"));
      assert.ok(card.includes("│ Paragraph 1"));
      assert.ok(card.includes("│ Paragraph 2 with more details"));
    },
  );

  await t.test(
    "renderProseCard safely handles empty or non-string inputs",
    () => {
      assert.strictEqual(renderProseCard(""), "");
      assert.strictEqual(renderProseCard("   \n\r  "), "");
      assert.strictEqual(renderProseCard(null), "");
      assert.strictEqual(renderProseCard(undefined), "");
    },
  );

  await t.test(
    "renderProseCard sanitizes ANSI escape and terminal control sequences",
    () => {
      const card = renderProseCard(
        "\x1b[31m" + "Red text" + "\x1b[0m with green info",
        80,
      );
      assert.ok(!card.includes("\x1b"));
      assert.ok(card.includes("Red text with green info"));
    },
  );
});

test("loadWorkspaceSkillsContext discovery and formatting", async (t) => {
  await t.test(
    "discovers skills in .gemini/skills and formats available_skills block with mandate",
    async () => {
      resetSkillsContextCache();
      const skillsContext = await loadWorkspaceSkillsContext();

      assert.ok(skillsContext.includes("# Available Workspace Skills"));
      assert.ok(skillsContext.includes("<available_skills>"));
      assert.ok(skillsContext.includes("<name>git</name>"));
      assert.ok(skillsContext.includes("<name>github-pr</name>"));
      assert.ok(skillsContext.includes("<name>github-ci</name>"));
      assert.ok(
        skillsContext.includes(
          "node .gemini/skills/git/scripts/read-ignored-file.js",
        ),
      );
      assert.ok(skillsContext.includes("Skills-First Mandate"));
    },
  );

  await t.test(
    "caches skills context across repeated calls and resetSkillsContextCache clears it",
    async () => {
      resetSkillsContextCache();
      const first = await loadWorkspaceSkillsContext();
      const second = await loadWorkspaceSkillsContext();
      assert.strictEqual(first, second);

      resetSkillsContextCache();
      const third = await loadWorkspaceSkillsContext();
      assert.strictEqual(typeof third, "string");
      assert.ok(third.includes("<available_skills>"));
    },
  );

  await t.test(
    "handles non-existent skills directory safely without throwing",
    async () => {
      resetSkillsContextCache();
      const skillsContext = await loadWorkspaceSkillsContext(
        "/nonexistent/repo/root",
      );
      assert.strictEqual(typeof skillsContext, "string");
      assert.ok(!skillsContext.includes("<available_skills>"));
      assert.ok(skillsContext.includes("Skills-First Mandate"));
    },
  );

  await t.test(
    "getSkillsFingerprint changes when a file is added, modified, or removed",
    async () => {
      const tempDir = await fs.mkdtemp(path.join(os.tmpdir(), "skills-test-"));
      try {
        const skillSubdir = path.join(tempDir, "demo");
        await fs.mkdir(skillSubdir, { recursive: true });
        const testFile = path.join(skillSubdir, "SKILL.md");
        await fs.writeFile(testFile, "initial content");

        const fp1 = await getSkillsFingerprint(tempDir);
        assert.ok(fp1.includes("demo/SKILL.md"));

        // Modifying file changes fingerprint
        await fs.writeFile(
          testFile,
          "updated content that is much longer to guarantee size change",
        );
        const fp2 = await getSkillsFingerprint(tempDir);
        assert.notStrictEqual(fp1, fp2);

        // Deleting file changes fingerprint
        await fs.unlink(testFile);
        const fp3 = await getSkillsFingerprint(tempDir);
        assert.notStrictEqual(fp2, fp3);
        assert.ok(!fp3.includes("SKILL.md"));
      } finally {
        await fs.rm(tempDir, { recursive: true, force: true });
      }
    },
  );

  await t.test(
    "loadWorkspaceSkillsContext safely handles primitive error in fingerprint getter",
    async () => {
      resetSkillsContextCache();
      // Path that is a file, not a directory, so readdir or stat might throw or behave unexpectedly
      const tempFile = path.join(os.tmpdir(), `test-file-${Date.now()}.txt`);
      await fs.writeFile(tempFile, "dummy");
      try {
        const skillsContext = await loadWorkspaceSkillsContext(tempFile);
        assert.strictEqual(typeof skillsContext, "string");
      } finally {
        await fs.unlink(tempFile).catch(() => {});
      }
    },
  );

  await t.test(
    "loadWorkspaceSkillsContext memoizes context within TTL without re-reading fingerprint",
    async () => {
      resetSkillsContextCache();
      const first = await loadWorkspaceSkillsContext();
      const second = await loadWorkspaceSkillsContext();
      assert.strictEqual(first, second);
    },
  );
});
