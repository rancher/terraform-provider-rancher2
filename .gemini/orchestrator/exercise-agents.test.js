import assert from "node:assert/strict";
import test from "node:test";
import {
  DEFAULT_EXERCISE_MODELS,
  closeSessionSafe,
  parseExerciseModels,
} from "./exercise-agents.js";

test("parseExerciseModels", async (t) => {
  await t.test(
    "falls back to DEFAULT_EXERCISE_MODELS when input is undefined or empty",
    () => {
      assert.deepStrictEqual(
        parseExerciseModels(undefined),
        DEFAULT_EXERCISE_MODELS,
      );
      assert.deepStrictEqual(parseExerciseModels(""), DEFAULT_EXERCISE_MODELS);
      assert.deepStrictEqual(
        parseExerciseModels(null),
        DEFAULT_EXERCISE_MODELS,
      );
    },
  );

  await t.test("parses and trims valid comma-separated models", () => {
    const res = parseExerciseModels("gemini-2.5-pro, gemini-3.8-flash ");
    assert.deepStrictEqual(res, ["gemini-2.5-pro", "gemini-3.8-flash"]);
  });

  await t.test(
    "filters out empty strings resulting from leading/trailing or consecutive commas",
    () => {
      const res = parseExerciseModels(",gemini-1.5, , gemini-3.0, ");
      assert.deepStrictEqual(res, ["gemini-1.5", "gemini-3.0"]);
    },
  );

  await t.test(
    "falls back to default models if input consists solely of whitespace and commas",
    () => {
      assert.deepStrictEqual(
        parseExerciseModels("   "),
        DEFAULT_EXERCISE_MODELS,
      );
      assert.deepStrictEqual(
        parseExerciseModels(", , , "),
        DEFAULT_EXERCISE_MODELS,
      );
    },
  );
});

test("closeSessionSafe", async (t) => {
  await t.test("successfully invokes session.close() if defined", async () => {
    let closed = false;
    const mockSession = {
      close: async () => {
        closed = true;
      },
    };
    await closeSessionSafe(mockSession);
    assert.strictEqual(closed, true);
  });

  await t.test("gracefully swallows errors or rejections from session.close()", async () => {
    let called = false;
    const mockSession = {
      close: async () => {
        called = true;
        throw new Error("Failed to close socket handle");
      },
    };
    await assert.doesNotReject(async () => {
      await closeSessionSafe(mockSession);
    });
    assert.strictEqual(called, true);
  });

  await t.test("handles null, undefined, or objects without a close method", async () => {
    await assert.doesNotReject(async () => {
      await closeSessionSafe(null);
      await closeSessionSafe(undefined);
      await closeSessionSafe({});
      await closeSessionSafe({ close: "not a function" });
    });
  });
});

