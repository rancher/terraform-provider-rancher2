import { TerminalQuotaError } from "@google/gemini-cli-core";
import { GeminiCliAgent } from "@google/gemini-cli-sdk";
import assert from "node:assert";
import { Buffer } from "node:buffer";
import fs from "node:fs/promises";
import path from "node:path";
import test from "node:test";
import {
  flushLogs,
  initializeAgentRunner,
  isAbortError,
  isMaxTurnsError,
  runAgentSession,
} from "./agent-runner.js";
import { promptTurnBudgetExhaustion } from "./turn-accounting.js";

test("TerminalQuotaError massive delay interceptor patch", async (t) => {
  await t.test("keeps small retry delay and original message untouched", () => {
    const rawMsg =
      "You have exhausted your capacity on this model. Your quota will reset after 5s.";
    const err = new TerminalQuotaError(
      rawMsg,
      { code: 429 },
      5,
      "MODEL_CAPACITY_EXHAUSTED",
    );

    assert.strictEqual(err.retryDelayMs, 5000);
    assert.strictEqual(err.message, rawMsg);
    assert.strictEqual(err.reason, "MODEL_CAPACITY_EXHAUSTED");
  });

  await t.test(
    "censors message and reason when retry delay is massive (> 5 minutes)",
    () => {
      const rawMsg =
        "You have exhausted your capacity on this model. Your quota will reset after 12h53m47s.";
      const err = new TerminalQuotaError(
        rawMsg,
        { code: 429 },
        46427,
        "MODEL_CAPACITY_EXHAUSTED",
      );

      assert.strictEqual(err.retryDelayMs, 46427000);
      assert.match(err.message, /exhausted capacity \(immediate fallback\)/);
      assert.doesNotMatch(err.message, /exhausted your capacity/);
      assert.strictEqual(
        err.reason,
        "MODEL_CAPACITY_EXHAUSTED_IMMEDIATE_FALLBACK",
      );
    },
  );
});

test("isMaxTurnsError detection", async (t) => {
  await t.test(
    "identifies max_turns errors in standard strings and Error objects",
    () => {
      assert.strictEqual(isMaxTurnsError("turn limit reached"), true);
      assert.strictEqual(isMaxTurnsError("max_turns exceeded"), true);
      assert.strictEqual(isMaxTurnsError("max-turns-exceeded"), true);
      assert.strictEqual(isMaxTurnsError("turn-limit-reached"), true);
      assert.strictEqual(isMaxTurnsError("maximum turns reached"), true);
      assert.strictEqual(
        isMaxTurnsError(new Error("Agent exceeded maximum turns limit")),
        true,
      );
      assert.strictEqual(isMaxTurnsError(new Error("Hit turn limit")), true);
    },
  );

  await t.test(
    "identifies max_turns when err.error is a string, object, or Error instance",
    () => {
      assert.strictEqual(
        isMaxTurnsError({ error: "max_turns_exceeded" }),
        true,
      );
      assert.strictEqual(
        isMaxTurnsError({
          error: { code: 400, message: "max_turns exceeded" },
        }),
        true,
      );
      assert.strictEqual(
        isMaxTurnsError({ error: { code: 400, message: "Invalid argument" } }),
        false,
      );

      // Top-level message alongside nested error object (e.g. HTTP client error)
      assert.strictEqual(
        isMaxTurnsError({
          message: "Request failed with status code 400",
          error: { code: 400, message: "max_turns exceeded" },
        }),
        true,
      );

      // Empty message property with valid error property
      assert.strictEqual(
        isMaxTurnsError({ message: "", error: "max_turns exceeded" }),
        true,
      );

      // Error instance with attached nested plain object details
      const customErr = new Error("Request failed with status code 400");
      customErr.error = { message: "maximum turns reached" };
      assert.strictEqual(isMaxTurnsError(customErr), true);

      // Nested Error instance attached to err.error (where Error properties are non-enumerable)
      const nestedErrInstance = new Error(
        "Request failed with status code 400",
      );
      nestedErrInstance.error = new Error("max_turns exceeded");
      assert.strictEqual(isMaxTurnsError(nestedErrInstance), true);

      const unrelatedNestedErr = new Error(
        "Request failed with status code 400",
      );
      unrelatedNestedErr.error = new Error("database connection lost");
      assert.strictEqual(isMaxTurnsError(unrelatedNestedErr), false);

      // Standard properties (message, reason) formatted as nested objects
      assert.strictEqual(
        isMaxTurnsError({ message: { text: "max_turns exceeded" } }),
        true,
      );
      assert.strictEqual(
        isMaxTurnsError({ reason: { details: "turn limit reached" } }),
        true,
      );
    },
  );

  await t.test("identifies max_turns in chained errors via err.cause", () => {
    // Cause as an Error instance
    const causeErr = new Error("Agent execution failed", {
      cause: new Error("turn limit reached"),
    });
    assert.strictEqual(isMaxTurnsError(causeErr), true);

    // Cause as a string
    const causeStringErr = new Error("Agent execution failed", {
      cause: "maximum turns reached",
    });
    assert.strictEqual(isMaxTurnsError(causeStringErr), true);

    // Cause as a plain object with message
    const causeObjErr = new Error("Agent execution failed", {
      cause: { message: "max_turns exceeded" },
    });
    assert.strictEqual(isMaxTurnsError(causeObjErr), true);

    // Unrelated cause
    const unrelatedCauseErr = new Error("Agent execution failed", {
      cause: new Error("network timeout"),
    });
    assert.strictEqual(isMaxTurnsError(unrelatedCauseErr), false);
  });

  await t.test(
    "identifies max_turns in HTTP response payloads (err.response.data) with circular references",
    () => {
      // Response data as nested object
      const httpLikeErr = new Error("Request failed with status code 400");
      httpLikeErr.response = {
        data: { error: { message: "max_turns exceeded" } },
      };
      httpLikeErr.circularRef = httpLikeErr;
      assert.strictEqual(isMaxTurnsError(httpLikeErr), true);

      // Response body as nested object (e.g. fetch / got / octokit client error)
      const fetchLikeErr = new Error("HTTP 400 Bad Request");
      fetchLikeErr.response = {
        body: { error: { message: "max_turns exceeded" } },
      };
      assert.strictEqual(isMaxTurnsError(fetchLikeErr), true);

      // Response body as raw string
      const fetchStringErr = new Error("HTTP 504 Gateway Timeout");
      fetchStringErr.response = { body: "turn limit reached" };
      assert.strictEqual(isMaxTurnsError(fetchStringErr), true);

      // Response data as raw string
      const httpStringErr = new Error("Bad Request");
      httpStringErr.response = {
        data: "Agent hit turn limit during processing",
      };
      httpStringErr.self = httpStringErr;
      assert.strictEqual(isMaxTurnsError(httpStringErr), true);

      // Response statusText matching turn limit (e.g. gateway error with empty body)
      const statusTextErr = new Error("Gateway Error");
      statusTextErr.response = {
        status: 504,
        statusText: "Turn Limit Exceeded",
      };
      assert.strictEqual(isMaxTurnsError(statusTextErr), true);

      // Response as direct raw string
      const rawStringResponseErr = new Error("Gateway Error");
      rawStringResponseErr.response = "504 Gateway Error: turn limit reached";
      assert.strictEqual(isMaxTurnsError(rawStringResponseErr), true);

      // Response data as Buffer containing plain text
      const httpBufferTextErr = new Error("Gateway Error");
      httpBufferTextErr.response = {
        data: Buffer.from("Agent hit turn limit during processing"),
      };
      assert.strictEqual(isMaxTurnsError(httpBufferTextErr), true);

      // Response data as Buffer containing JSON error object
      const httpBufferJsonErr = new Error("Bad Request");
      httpBufferJsonErr.response = {
        data: Buffer.from(
          JSON.stringify({ error: { message: "max_turns exceeded" } }),
        ),
      };
      assert.strictEqual(isMaxTurnsError(httpBufferJsonErr), true);

      // Response data as Buffer containing unrelated error
      const httpBufferUnrelatedErr = new Error("Internal Error");
      httpBufferUnrelatedErr.response = {
        data: Buffer.from(
          JSON.stringify({ error: { message: "Internal Server Error" } }),
        ),
      };
      assert.strictEqual(isMaxTurnsError(httpBufferUnrelatedErr), false);

      // Response data as Buffer containing ignored keys like config with turn mentions
      const httpBufferPromptIgnoredErr = new Error("Internal Error");
      httpBufferPromptIgnoredErr.response = {
        data: Buffer.from(
          JSON.stringify({
            config: { prompt: "fix max_turns" },
            error: { message: "Internal Server Error" },
          }),
        ),
      };
      assert.strictEqual(isMaxTurnsError(httpBufferPromptIgnoredErr), false);

      // Oversized buffer (>64KB) is skipped safely
      const oversizedBufferErr = new Error("Oversized Error");
      oversizedBufferErr.response = {
        data: Buffer.concat([
          Buffer.from("turn limit reached "),
          Buffer.alloc(65 * 1024, 0),
        ]),
      };
      assert.strictEqual(isMaxTurnsError(oversizedBufferErr), false);

      // Non-Buffer typed arrays are safely ignored
      const typedArrayErr = new Error("Binary Error");
      typedArrayErr.response = { data: new Uint8Array([1, 2, 3]) };
      assert.strictEqual(isMaxTurnsError(typedArrayErr), false);

      // Unrelated HTTP error
      const unrelatedHttpErr = new Error("Request failed with status code 500");
      unrelatedHttpErr.response = { data: { error: "Internal Server Error" } };
      assert.strictEqual(isMaxTurnsError(unrelatedHttpErr), false);
    },
  );

  await t.test(
    "handles non-error objects, null, undefined, and primitives safely",
    () => {
      assert.strictEqual(isMaxTurnsError(null), false);
      assert.strictEqual(isMaxTurnsError(undefined), false);
      assert.strictEqual(isMaxTurnsError(""), false);
      assert.strictEqual(isMaxTurnsError(123), false);
      assert.strictEqual(isMaxTurnsError(new Error("Connection reset")), false);

      // Object with arbitrary property matching turn keywords
      assert.strictEqual(isMaxTurnsError({ details: "turn limit hit" }), true);
      assert.strictEqual(isMaxTurnsError({ details: "other error" }), false);
    },
  );

  await t.test(
    "identifies max_turns when specified via custom error name",
    () => {
      const customNameErr = new Error("Execution halted prematurely");
      customNameErr.name = "MaxTurnsExceededError";
      assert.strictEqual(isMaxTurnsError(customNameErr), true);

      const maxTurnsNameErr = new Error("");
      maxTurnsNameErr.name = "MaxTurnsError";
      assert.strictEqual(isMaxTurnsError(maxTurnsNameErr), true);
    },
  );

  await t.test("traverses deeply chained causes and array details", () => {
    const deepCauseErr = new Error("Level 1", {
      cause: new Error("Level 2", {
        cause: new Error("turn limit reached"),
      }),
    });
    assert.strictEqual(isMaxTurnsError(deepCauseErr), true);

    const arrayDetailsErr = new Error("Resource exhausted");
    arrayDetailsErr.details = [{ message: "max_turns exceeded" }];
    assert.strictEqual(isMaxTurnsError(arrayDetailsErr), true);
  });

  await t.test(
    "avoids false positives when request payload or client config mentions turn keywords",
    () => {
      const httpErrWithPrompt = new Error(
        "Request failed with status code 500",
      );
      httpErrWithPrompt.config = {
        data: JSON.stringify({
          prompt: "Please fix the max_turns and turn limit bug in agent-runner",
        }),
      };
      httpErrWithPrompt.request = {
        body: "turn limit prompt content",
      };
      httpErrWithPrompt.response = {
        status: 500,
        data: { error: { message: "Internal Server Error" } },
      };
      assert.strictEqual(isMaxTurnsError(httpErrWithPrompt), false);
    },
  );

  await t.test(
    "identifies max_turns when signaled via property keys and snake_case properties",
    () => {
      assert.strictEqual(isMaxTurnsError({ max_turns: true }), true);
      assert.strictEqual(isMaxTurnsError({ max_turns_exceeded: true }), true);
      assert.strictEqual(isMaxTurnsError({ turn_limit: 10 }), true);
      assert.strictEqual(isMaxTurnsError({ turn_limit_exceeded: true }), true);
      assert.strictEqual(isMaxTurnsError({ details: { max_turns: 5 } }), true);

      // Explicit false boolean flags and falsy/zero values should not trigger false positives
      assert.strictEqual(isMaxTurnsError({ max_turns_exceeded: false }), false);
      assert.strictEqual(isMaxTurnsError({ turn_limit_reached: false }), false);
      assert.strictEqual(
        isMaxTurnsError({ status: "OK", details: { max_turns: false } }),
        false,
      );
      assert.strictEqual(isMaxTurnsError({ max_turns: 0 }), false);
      assert.strictEqual(isMaxTurnsError({ turn_limit: 0 }), false);
      assert.strictEqual(isMaxTurnsError({ max_turns: null }), false);
      assert.strictEqual(
        isMaxTurnsError({ details: { max_turns: null } }),
        false,
      );
    },
  );

  await t.test("identifies max_turns with singular phrasing variant", () => {
    assert.strictEqual(isMaxTurnsError("max turn reached"), true);
    assert.strictEqual(
      isMaxTurnsError(new Error("Hit max turn count limit")),
      true,
    );
    assert.strictEqual(isMaxTurnsError("exceeded maximum turn count"), true);
    assert.strictEqual(
      isMaxTurnsError(new Error("Agent reached maximum turn threshold")),
      true,
    );
  });

  await t.test(
    "handles throwing property getters gracefully without uncaught exceptions",
    () => {
      const throwingErr = {
        get message() {
          throw new Error("Exploding getter");
        },
        get response() {
          throw new Error("Exploding response");
        },
        error: "max_turns exceeded",
      };
      assert.strictEqual(isMaxTurnsError(throwingErr), true);

      const throwingArbitraryGetterErr = {
        get brokenProp() {
          throw new Error("Exploding arbitrary getter");
        },
        customTurnMessage: "turn limit reached",
      };
      assert.strictEqual(isMaxTurnsError(throwingArbitraryGetterErr), true);

      const completelyBrokenErr = {
        get message() {
          throw new Error("Exploding getter");
        },
        get details() {
          throw new Error("Exploding details");
        },
      };
      assert.strictEqual(isMaxTurnsError(completelyBrokenErr), false);
    },
  );

  await t.test(
    "avoids false positives for return limits, turnaround times, and adjacent keys",
    () => {
      // Return limit mentions in query and pagination errors
      assert.strictEqual(isMaxTurnsError("return limit reached"), false);
      assert.strictEqual(isMaxTurnsError("Query return limit exceeded"), false);
      assert.strictEqual(
        isMaxTurnsError("Failed to return limit parameter"),
        false,
      );
      assert.strictEqual(
        isMaxTurnsError("Invalid return_limit specified"),
        false,
      );
      assert.strictEqual(
        isMaxTurnsError(
          new Error(
            "API error: parameter return_limit must be a positive integer",
          ),
        ),
        false,
      );
      assert.strictEqual(isMaxTurnsError({ return_limit: 50 }), false);

      // Turnaround and turnout wording
      assert.strictEqual(
        isMaxTurnsError("Maximum turnaround time exceeded"),
        false,
      );
      assert.strictEqual(isMaxTurnsError("Max turnaround time reached"), false);
      assert.strictEqual(
        isMaxTurnsError(new Error("Request exceeded maximum turnaround SLA")),
        false,
      );
      assert.strictEqual(isMaxTurnsError("Max turnout reached"), false);

      // Adjacent property keys and values that could accidentally form phrases if joined indiscriminately
      assert.strictEqual(isMaxTurnsError({ turn: 1, limit: 100 }), false);
      assert.strictEqual(isMaxTurnsError({ max: 10, turn: 2 }), false);
      assert.strictEqual(isMaxTurnsError({ return: 1, limit: 10 }), false);
      assert.strictEqual(
        isMaxTurnsError({ prefix: "max", suffix: "turn" }),
        false,
      );
      assert.strictEqual(
        isMaxTurnsError({ header: "turn", detail: "limit" }),
        false,
      );
      assert.strictEqual(
        isMaxTurnsError({ action: "turn to retry", name: "max" }),
        false,
      );

      // Configuration and context containers on unrelated errors should not trigger false positives
      const errWithOptions = new Error("Gateway timeout");
      errWithOptions.options = { max_turns: 25 };
      assert.strictEqual(isMaxTurnsError(errWithOptions), false);

      const errWithAgent = new Error("Database unreachable");
      errWithAgent.agent = { max_turns: 25 };
      assert.strictEqual(isMaxTurnsError(errWithAgent), false);

      const errWithParams = new Error("Service unavailable");
      errWithParams.params = { turn_limit: 10 };
      assert.strictEqual(isMaxTurnsError(errWithParams), false);
    },
  );
});

test("promptTurnBudgetExhaustion caller readline reuse", async (t) => {
  await t.test(
    "reuses caller readline interface without closing it",
    async () => {
      let closed = false;
      const fakeRl = {
        async question() {
          return "continue";
        },
        close() {
          closed = true;
        },
      };
      const decision = await promptTurnBudgetExhaustion(5, 5, 5, null, {
        rl: fakeRl,
      });
      assert.strictEqual(decision.shouldContinue, true);
      assert.strictEqual(decision.newBudget, 10);
      assert.strictEqual(
        closed,
        false,
        "Caller readline interface must not be closed",
      );
    },
  );

  await t.test(
    "aborts session when user responds stop via reused interface",
    async () => {
      let closed = false;
      let aborted = false;
      const fakeRl = {
        async question() {
          return "stop";
        },
        close() {
          closed = true;
        },
      };
      const controller = {
        abort() {
          aborted = true;
        },
      };
      const decision = await promptTurnBudgetExhaustion(5, 5, 5, controller, {
        rl: fakeRl,
      });
      assert.strictEqual(decision.shouldContinue, false);
      assert.strictEqual(aborted, true);
      assert.strictEqual(closed, false);
    },
  );
});

test("agent-runner initialization and logging", async (t) => {
  await t.test(
    "initializeAgentRunner loads configuration without errors and targets orchestrator.log",
    async () => {
      const config = await initializeAgentRunner();
      assert.ok(config.repoRoot);
      assert.ok(config.logPath);
      assert.ok(
        config.logPath.endsWith(
          path.join(".gemini", "tmp", "orchestrator.log"),
        ),
      );
      assert.ok(config.models);
      assert.ok(config.models.pro);
      assert.ok(config.models.flash);
      assert.ok(config.models.flash_lite);
      await flushLogs();
      // Verify flushLogs is idempotent and safe to call when stream is already flushed
      await flushLogs();
    },
  );
});

test("isAbortError and uncaughtException handling", async (t) => {
  await t.test(
    "isAbortError identifies standard abort errors and variants",
    () => {
      assert.strictEqual(isAbortError({ name: "AbortError" }), true);
      assert.strictEqual(isAbortError({ type: "aborted" }), true);
      assert.strictEqual(isAbortError({ code: "ABORT_ERR" }), true);
      assert.strictEqual(isAbortError(new Error("regular error")), false);
      assert.strictEqual(isAbortError(null), false);
      assert.strictEqual(isAbortError(undefined), false);
      assert.strictEqual(isAbortError("AbortError"), false);
    },
  );

  await t.test(
    "activeAbortHandler suppresses AbortError and preserves fatal handling for other errors",
    async () => {
      await flushLogs();
      await initializeAgentRunner();

      const listeners = process.listeners("uncaughtException");
      const handler = listeners[listeners.length - 1];
      assert.ok(typeof handler === "function");

      // AbortError should be silently handled without removing handler
      const abortErr = new Error("aborted");
      abortErr.name = "AbortError";
      assert.doesNotThrow(() => handler(abortErr, "uncaughtException"));
      assert.ok(process.listeners("uncaughtException").includes(handler));

      // Non-AbortError should throw when no other listeners exist
      const otherListeners = process
        .rawListeners("uncaughtException")
        .filter((l) => l !== handler);
      for (const l of otherListeners) {
        process.removeListener("uncaughtException", l);
      }

      try {
        const nonAbortErr = new Error("unexpected fatal failure");
        assert.throws(
          () => handler(nonAbortErr, "uncaughtException"),
          /unexpected fatal failure/,
        );
        assert.strictEqual(
          process.listeners("uncaughtException").includes(handler),
          false,
        );
      } finally {
        for (const l of otherListeners) {
          process.on("uncaughtException", l);
        }
      }

      // After flushLogs, handler should be cleanly removed
      await initializeAgentRunner();
      await flushLogs();
      assert.strictEqual(
        process.listeners("uncaughtException").includes(handler),
        false,
      );
    },
  );

  await t.test(
    "activeAbortHandler forwards non-AbortError when another listener exists",
    async () => {
      await flushLogs();
      await initializeAgentRunner();

      const listeners = process.listeners("uncaughtException");
      const handler = listeners[listeners.length - 1];

      const otherListeners = process
        .rawListeners("uncaughtException")
        .filter((l) => l !== handler);
      for (const l of otherListeners) {
        process.removeListener("uncaughtException", l);
      }

      let customListenerFired = false;
      const customListener = (err) => {
        if (err.message === "forwarded failure") {
          customListenerFired = true;
        }
      };

      process.on("uncaughtException", customListener);
      try {
        const regularErr = new Error("forwarded failure");
        assert.doesNotThrow(() => handler(regularErr, "uncaughtException"));
        assert.strictEqual(customListenerFired, true);
      } finally {
        process.removeListener("uncaughtException", customListener);
        for (const l of otherListeners) {
          process.on("uncaughtException", l);
        }
        await flushLogs();
      }
    },
  );
});

test("runAgentSession maxTurns override integration", async (t) => {
  const originalSession = GeminiCliAgent.prototype.session;
  t.after(() => {
    GeminiCliAgent.prototype.session = originalSession;
  });

  const attachMockSendStream = (agent, sendStream) => {
    const session = originalSession.call(agent);
    session.sendStream = sendStream;
    const originalInitialize = session.initialize.bind(session);
    session.initialize = async () => {
      // Allow internal side-effects that setup config/toolRegistry if needed,
      // but override config.refreshAuth to prevent real metadata checks.
      if (session.config) {
        session.config.refreshAuth = async () => {};
      }
      return originalInitialize();
    };
    return session;
  };

  await t.test(
    "verifies maxTurns override reaches SDK agent session and enforces hard stop",
    async () => {
      let capturedAgentOptions = null;
      try {
        const config = await initializeAgentRunner();
        const initialLogContent = await fs
          .readFile(config.logPath, "utf-8")
          .catch(() => "");
        const initialOffset = initialLogContent.length;

        GeminiCliAgent.prototype.session = function () {
          capturedAgentOptions = this.options;
          return attachMockSendStream(this, async function* () {
            yield { type: "content", value: "Turn 1 initial response\n" };
            for (let turn = 1; turn <= 5; turn++) {
              yield {
                type: "tool_call_request",
                value: { name: "no_op", args: "{}" },
              };
              yield { type: "tool_call_result", value: "success" };
            }
            yield {
              type: "content",
              value: "Turn 6 response should never be reached\n",
            };
          });
        };

        const result = await runAgentSession({
          initialPrompt: "Run validator session",
          maxTurns: 5,
          isolate: true,
        });

        assert.strictEqual(
          capturedAgentOptions?.max_turns,
          5,
          "maxTurns override must reach GeminiCliAgent options",
        );
        assert.match(result, /Turn 1 initial response/);
        assert.doesNotMatch(
          result,
          /Turn 6 response should never be reached/,
          "Session must halt at 5 turns without processing subsequent turns",
        );

        await flushLogs();
        const logContent = await fs.readFile(config.logPath, "utf-8");
        const sessionLog = logContent.slice(initialOffset);
        assert.match(
          sessionLog,
          /\[Turn Limit Reached\] Override limit enforced as a hard stop\./,
          "Log must record hard stop enforcement for turn override",
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
      }
    },
  );

  await t.test(
    "runAgentSession aborts and throws timeout error when streamIdleTimeoutMs expires",
    async () => {
      try {
        GeminiCliAgent.prototype.session = function () {
          return attachMockSendStream(this, async function* () {
            await new Promise((resolve) => setTimeout(resolve, 150));
            yield { type: "content", value: "Too late" };
          });
        };

        await assert.rejects(async () => {
          await runAgentSession({
            initialPrompt: "Timeout test",
            isolate: true,
            streamIdleTimeoutMs: 50,
          });
        }, /Stream idle timeout exceeded/);
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
      }
    },
  );

  await t.test(
    "runAgentSession resets streamIdleTimeoutMs when chunks arrive and completes",
    async () => {
      try {
        GeminiCliAgent.prototype.session = function () {
          return attachMockSendStream(this, async function* () {
            await new Promise((resolve) => setTimeout(resolve, 25));
            yield { type: "content", value: "Chunk 1 " };
            await new Promise((resolve) => setTimeout(resolve, 25));
            yield { type: "content", value: "Chunk 2" };
          });
        };

        const result = await runAgentSession({
          initialPrompt: "Reset timeout test",
          isolate: true,
          streamIdleTimeoutMs: 50,
        });

        assert.strictEqual(result, "Chunk 1 Chunk 2");
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
      }
    },
  );

  await t.test(
    "respects precedence of maxTurnsOverride over maxTurns",
    async () => {
      let capturedAgentOptions = null;
      try {
        GeminiCliAgent.prototype.session = function () {
          capturedAgentOptions = this.options;
          return attachMockSendStream(this, async function* () {
            yield { type: "content", value: "Turn 1 initial response\n" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run precedence check",
          maxTurns: 10,
          maxTurnsOverride: 3,
          isolate: true,
        });

        assert.strictEqual(
          capturedAgentOptions?.max_turns,
          3,
          "maxTurnsOverride must take precedence over maxTurns",
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
      }
    },
  );

  await t.test(
    "ignores boolean maxTurns inputs and avoids coercing to 1",
    async () => {
      let capturedAgentOptions = null;
      try {
        GeminiCliAgent.prototype.session = function () {
          capturedAgentOptions = this.options;
          return attachMockSendStream(this, async function* () {
            yield { type: "content", value: "Turn 1 initial response\n" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run boolean maxTurns check",
          maxTurns: true,
          isolate: true,
        });

        assert.notStrictEqual(
          capturedAgentOptions?.max_turns,
          1,
          "maxTurns: true must not be coerced to 1",
        );
        assert.ok(
          capturedAgentOptions?.max_turns > 1,
          "Default model turn capacity must be preserved",
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
      }
    },
  );

  await t.test(
    "handles non-positive maxTurns by falling back to default capacity",
    async () => {
      let capturedAgentOptions = null;
      try {
        GeminiCliAgent.prototype.session = function () {
          capturedAgentOptions = this.options;
          return attachMockSendStream(this, async function* () {
            yield { type: "content", value: "Turn 1 initial response\n" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run zero maxTurns check",
          maxTurns: 0,
          isolate: true,
        });

        assert.ok(
          capturedAgentOptions?.max_turns > 0,
          "Zero maxTurns must fall back to positive default turn budget",
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
      }
    },
  );
});

test("runAgentSession enriched tool logging", async (t) => {
  await t.test(
    "formats tool usage log with tool name, arguments, and execution result in rounded cards",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      const attachMockSendStream = (agent, sendStream) => {
        const session = originalSession.call(agent);
        session.sendStream = sendStream;
        const originalInitialize = session.initialize.bind(session);
        session.initialize = async () => {
          if (session.config) {
            session.config.refreshAuth = async () => {};
          }
          return originalInitialize();
        };
        return session;
      };

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          return attachMockSendStream(this, async function* () {
            yield {
              type: "tool_call_request",
              value: {
                name: "run_shell_command",
                args: { command: "git status" },
              },
            };
            yield { type: "tool_call_result", value: "On branch main" };

            yield {
              type: "tool_call_request",
              value: { name: "read_file", args: { file_path: "README.md" } },
            };
            yield { type: "tool_call_result", value: "# Documentation" };

            yield {
              type: "tool_call_request",
              value: {
                name: "replace",
                args: { file_path: "foo.js", new_string: "const a = 1;" },
              },
            };
            yield { type: "tool_call_result", value: "success" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run enriched tool logging test",
          maxTurns: 3,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const cards = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("╰─"),
        );
        assert.strictEqual(cards.length, 3);
        assert.ok(cards[0].includes("✓  Shell git status"));
        assert.ok(cards[0].includes("On branch main"));
        assert.ok(cards[1].includes("✓  ReadFile README.md"));
        assert.ok(cards[1].includes("# Documentation"));
        assert.ok(cards[2].includes("✓  Replace foo.js"));
        assert.ok(cards[2].includes("success"));
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "formats tool usage log as a clean card when agent emits no prose",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      const attachMockSendStream = (agent, sendStream) => {
        const session = originalSession.call(agent);
        session.sendStream = sendStream;
        const originalInitialize = session.initialize.bind(session);
        session.initialize = async () => {
          if (session.config) {
            session.config.refreshAuth = async () => {};
          }
          return originalInitialize();
        };
        return session;
      };

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          return attachMockSendStream(this, async function* () {
            yield {
              type: "tool_call_request",
              value: {
                name: "run_shell_command",
                args: { command: "git status" },
              },
            };
            yield { type: "tool_call_result", value: "On branch main" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run no intent test",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const cards = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("Shell"),
        );
        assert.strictEqual(cards.length, 1);
        assert.ok(cards[0].includes("✓  Shell git status"));
        assert.ok(cards[0].includes("On branch main"));
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "renders topic updates inside a prominent rounded topic banner",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      const attachMockSendStream = (agent, sendStream) => {
        const session = originalSession.call(agent);
        session.sendStream = sendStream;
        const originalInitialize = session.initialize.bind(session);
        session.initialize = async () => {
          if (session.config) {
            session.config.refreshAuth = async () => {};
          }
          return originalInitialize();
        };
        return session;
      };

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          return attachMockSendStream(this, async function* () {
            yield {
              type: "tool_call_request",
              value: {
                name: "update_topic",
                args: { strategic_intent: "Reviewing the staged changes" },
              },
            };
            yield { type: "tool_call_result", value: "updated" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run multiline box test",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const banners = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("Objective"),
        );
        assert.strictEqual(banners.length, 1);
        assert.ok(banners[0].includes("╭─"));
        assert.ok(banners[0].includes("╰─"));
        assert.ok(
          banners[0].includes("🎯 Objective: Reviewing the staged changes"),
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "formats complete topic banner with title, objective, and summary",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      const attachMockSendStream = (agent, sendStream) => {
        const session = originalSession.call(agent);
        session.sendStream = sendStream;
        const originalInitialize = session.initialize.bind(session);
        session.initialize = async () => {
          if (session.config) {
            session.config.refreshAuth = async () => {};
          }
          return originalInitialize();
        };
        return session;
      };

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          return attachMockSendStream(this, async function* () {
            yield {
              type: "tool_call_request",
              value: {
                name: "update_topic",
                args: {
                  title: "Researching Bug",
                  strategic_intent: "Reproduce parser failure",
                  summary: "Starting investigation",
                },
              },
            };
            yield { type: "tool_call_result", value: "updated" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run complete topic banner test",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const banners = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("Topic"),
        );
        assert.strictEqual(banners.length, 1);
        assert.ok(banners[0].includes("📂 Topic: Researching Bug"));
        assert.ok(
          banners[0].includes("🎯 Objective: Reproduce parser failure"),
        );
        assert.ok(banners[0].includes("📝 Summary: Starting investigation"));
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "does not treat ordinary prose as tool card content",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      const attachMockSendStream = (agent, sendStream) => {
        const session = originalSession.call(agent);
        session.sendStream = sendStream;
        const originalInitialize = session.initialize.bind(session);
        session.initialize = async () => {
          if (session.config) {
            session.config.refreshAuth = async () => {};
          }
          return originalInitialize();
        };
        return session;
      };

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          return attachMockSendStream(this, async function* () {
            yield { type: "content", value: "Turn 1 initial response\n" };
            yield {
              type: "tool_call_request",
              value: { name: "no_op", args: "{}" },
            };
            yield { type: "tool_call_result", value: "ok" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run prose intent test",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const cards = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("NoOp"),
        );
        assert.strictEqual(cards.length, 1);
        assert.ok(cards[0].includes("✓  NoOp"));
        assert.ok(!cards[0].includes("Turn 1 initial response"));
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "renders distinct rounded cards for sequential tool calls within the same turn",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      const attachMockSendStream = (agent, sendStream) => {
        const session = originalSession.call(agent);
        session.sendStream = sendStream;
        const originalInitialize = session.initialize.bind(session);
        session.initialize = async () => {
          if (session.config) {
            session.config.refreshAuth = async () => {};
          }
          return originalInitialize();
        };
        return session;
      };

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          return attachMockSendStream(this, async function* () {
            yield {
              type: "tool_call_request",
              value: {
                name: "run_shell_command",
                args: { command: "git status" },
              },
            };
            yield {
              type: "tool_call_request",
              value: { name: "read_file", args: { file_path: "README.md" } },
            };
            yield { type: "tool_call_result", value: "clean" };
            yield { type: "tool_call_result", value: "content" };
          });
        };

        await runAgentSession({
          initialPrompt: "Run sequential multi-tool test",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const cards = logs.filter((msg) => msg.includes("╭─"));
        assert.strictEqual(cards.length, 2);
        assert.ok(cards[0].includes("✓  Shell git status"));
        assert.ok(cards[0].includes("clean"));
        assert.ok(cards[1].includes("✓  ReadFile README.md"));
        assert.ok(cards[1].includes("content"));
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "automatically registers update_topic tool in non-isolated sessions",
    async () => {
      let registeredToolNames = [];
      const originalSession = GeminiCliAgent.prototype.session;

      try {
        GeminiCliAgent.prototype.session = function () {
          const session = originalSession.call(this);
          session.sendStream = async function* () {
            yield { type: "content", value: "hello" };
          };
          const origInit = session.initialize.bind(session);
          session.initialize = async () => {
            if (session.config) {
              session.config.refreshAuth = async () => {};
            }
            await origInit();
            registeredToolNames =
              session.config?.toolRegistry?.getAllToolNames() || [];
          };
          return session;
        };

        await runAgentSession({
          initialPrompt: "Check registered tools",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        assert.ok(
          registeredToolNames.includes("update_topic"),
          "update_topic tool must be registered in non-isolated sessions",
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
      }
    },
  );

  await t.test(
    "renders update_topic banner only once when onToolExecution runs",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          const session = originalSession.call(this);
          session.sendStream = async function* () {
            yield {
              type: "tool_call_request",
              value: {
                name: "update_topic",
                args: {
                  title: "Investigation",
                  strategic_intent: "Find root cause",
                },
              },
            };
            // Simulate tool execution through toolRegistry triggering onToolExecution
            const updateTopic =
              session.config?.toolRegistry?.getTool("update_topic");
            if (updateTopic) {
              const invocation = updateTopic.createInvocation
                ? updateTopic.createInvocation()
                : updateTopic;
              await invocation.execute({
                title: "Investigation",
                strategic_intent: "Find root cause",
              });
            }
            yield { type: "tool_call_result", value: "Topic updated" };
          };
          const origInit = session.initialize.bind(session);
          session.initialize = async () => {
            if (session.config) {
              session.config.refreshAuth = async () => {};
            }
            return origInit();
          };
          return session;
        };

        await runAgentSession({
          initialPrompt: "Run single topic banner test",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const banners = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("Topic"),
        );
        assert.strictEqual(
          banners.length,
          1,
          "Topic banner must only be rendered once without duplicates",
        );
        assert.ok(banners[0].includes("📂 Topic: Investigation"));
        assert.ok(banners[0].includes("🎯 Objective: Find root cause"));
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "renders error card when fallback tool_call_result contains error: true",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          const session = originalSession.call(this);
          session.sendStream = async function* () {
            yield {
              type: "tool_call_request",
              value: {
                name: "read_file",
                args: { file_path: "secret.txt" },
              },
            };
            yield {
              type: "tool_call_result",
              value: {
                error: true,
                llmContent:
                  "Tool call deferred: Please call update_topic first.",
              },
            };
          };
          const origInit = session.initialize.bind(session);
          session.initialize = async () => {
            if (session.config) {
              session.config.refreshAuth = async () => {};
            }
            return origInit();
          };
          return session;
        };

        await runAgentSession({
          initialPrompt: "Run error card test",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const cards = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("ReadFile"),
        );
        assert.strictEqual(cards.length, 1);
        assert.ok(
          cards[0].includes("❌  ReadFile secret.txt"),
          "Card must show error symbol ❌ instead of success symbol ✓",
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "throws explicit error when update_topic is blocked while requireTopic is enabled",
    async () => {
      await assert.rejects(async () => {
        await runAgentSession({
          initialPrompt: "Run blocked update_topic session",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
          blockTools: ["update_topic"],
        });
      }, /Cannot block "update_topic" tool when requireTopic is enabled/);
    },
  );

  await t.test(
    "renders distinct cards for identical repeated tool calls without deduplication collisions",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          const session = originalSession.call(this);
          session.sendStream = async function* () {
            // First tool call
            yield {
              type: "tool_call_request",
              value: {
                name: "read_file",
                args: { file_path: "repeated.txt" },
              },
            };
            yield {
              type: "tool_call_result",
              value: "First read content",
            };
            // Second identical tool call with same name and same arguments
            yield {
              type: "tool_call_request",
              value: {
                name: "read_file",
                args: { file_path: "repeated.txt" },
              },
            };
            yield {
              type: "tool_call_result",
              value: "Second read content",
            };
          };
          const origInit = session.initialize.bind(session);
          session.initialize = async () => {
            if (session.config) {
              session.config.refreshAuth = async () => {};
            }
            return origInit();
          };
          return session;
        };

        await runAgentSession({
          initialPrompt: "Run repeated tool calls test",
          maxTurns: 5,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const cards = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("ReadFile repeated.txt"),
        );
        assert.strictEqual(
          cards.length,
          2,
          "Both identical tool calls must render their own card",
        );
        assert.ok(
          cards[0].includes("First read content"),
          "First card contains first result",
        );
        assert.ok(
          cards[1].includes("Second read content"),
          "Second card contains second result",
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );

  await t.test(
    "throws explicit error when requireTopic is true with isolate true",
    async () => {
      await assert.rejects(async () => {
        await runAgentSession({
          initialPrompt: "Run isolated session with topic requirement",
          maxTurns: 1,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: true,
          requireTopic: true,
        });
      }, /Cannot enable requireTopic when isolate is true/);
    },
  );

  await t.test(
    "correlates parallel tool results by callId when completed out of order",
    async () => {
      const originalSession = GeminiCliAgent.prototype.session;
      const originalConsoleLog = console.log;
      const logs = [];

      try {
        console.log = (...args) => {
          logs.push(args.join(" "));
        };

        GeminiCliAgent.prototype.session = function () {
          const session = originalSession.call(this);
          session.sendStream = async function* () {
            // Two parallel tool call requests with unique callIds
            yield {
              type: "tool_call_request",
              value: {
                callId: "call_alpha",
                name: "read_file",
                args: { file_path: "alpha.txt" },
              },
            };
            yield {
              type: "tool_call_request",
              value: {
                callId: "call_beta",
                name: "read_file",
                args: { file_path: "beta.txt" },
              },
            };
            // Results complete out of order: beta arrives before alpha
            yield {
              type: "tool_call_result",
              callId: "call_beta",
              value: "Content of beta",
            };
            yield {
              type: "tool_call_result",
              callId: "call_alpha",
              value: "Content of alpha",
            };
          };
          const origInit = session.initialize.bind(session);
          session.initialize = async () => {
            if (session.config) {
              session.config.refreshAuth = async () => {};
            }
            return origInit();
          };
          return session;
        };

        await runAgentSession({
          initialPrompt: "Run parallel out of order test",
          maxTurns: 5,
          requestedModel: "gemini-3.1-pro-preview",
          isolate: false,
        });

        const cards = logs.filter(
          (msg) => msg.includes("╭─") && msg.includes("ReadFile"),
        );
        assert.strictEqual(cards.length, 2, "Both parallel calls render cards");
        // Verify beta was rendered with beta content, and alpha with alpha content
        const betaCard = cards.find((c) => c.includes("beta.txt"));
        const alphaCard = cards.find((c) => c.includes("alpha.txt"));
        assert.ok(betaCard, "beta card exists");
        assert.ok(alphaCard, "alpha card exists");
        assert.ok(
          betaCard.includes("Content of beta"),
          "beta card contains beta content",
        );
        assert.ok(
          alphaCard.includes("Content of alpha"),
          "alpha card contains alpha content",
        );
      } finally {
        GeminiCliAgent.prototype.session = originalSession;
        console.log = originalConsoleLog;
      }
    },
  );
});
