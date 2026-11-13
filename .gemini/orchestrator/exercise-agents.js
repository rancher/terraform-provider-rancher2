import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { GeminiCliAgent } from "@google/gemini-cli-sdk";
import { log, logError } from "./lib/logger.js";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Open/Closed Principle compliant: Load models dynamically via environment variables with a safe fallback
export const DEFAULT_EXERCISE_MODELS = [
  "gemini-3.1-pro-preview",
  "gemini-3.8-flash",
  "gemini-2.5-pro",
  "gemini-3.1-flash-lite",
];

export function parseExerciseModels(
  envModels = process.env.GEMINI_EXERCISE_MODELS,
) {
  if (!envModels) {
    return DEFAULT_EXERCISE_MODELS;
  }
  const parsed = envModels
    .split(",")
    .map((m) => m.trim())
    .filter(Boolean);
  return parsed.length > 0 ? parsed : DEFAULT_EXERCISE_MODELS;
}

export const models = parseExerciseModels();

export async function closeSessionSafe(session) {
  if (session && typeof session.close === "function") {
    await session.close().catch(() => {});
  }
}

async function main() {
  const tmpDir = path.resolve(__dirname, "..", "tmp");
  await fs.mkdir(tmpDir, { recursive: true });
  const logFilePath = path.join(tmpDir, "exercise-agents.log");
  let logContent = "";
  const logWrapper = (msg) => {
    log(msg);
    logContent += `[INFO] ${new Date().toISOString()} - ${msg}\n`;
  };
  const logErrorWrapper = (msg) => {
    logError(msg);
    logContent += `[ERROR] ${new Date().toISOString()} - ${msg}\n`;
  };

  logWrapper(
    `📅 Morning Check-in: Initiating quota resets for available Gemini CLI models...`,
  );
  logWrapper(`Models to trigger: ${models.join(", ")}`);

  const summary = {
    passed: [],
    failed: [],
  };

  for (const model of models) {
    logWrapper(`\n=========================================`);
    logWrapper(`☀️ Triggering Model: ${model}`);
    logWrapper(`=========================================`);

    let timeoutId;
    let session;
    try {
      // Instantiate fast agent using the SDK with the target model
      const agent = new GeminiCliAgent({
        model: model,
        instructions: "You are a polite, helpful assistant.",
      });

      const controller = new globalThis.AbortController();
      timeoutId = setTimeout(() => {
        controller.abort(
          new Error("Model quota exercise request timed out after 60s."),
        );
      }, 60000);
      const greetingPrompt = `Good morning! It is 8:00 AM. Just saying a quick hello to confirm you are awake and ready for the day!`;

      session = agent.session();
      await session.initialize();
      const stream = session.sendStream(greetingPrompt, controller.signal);
      let responseText = "";

      for await (const chunk of stream) {
        if (chunk.type === "content") {
          const text = chunk.value || "";
          process.stdout.write(text);
          responseText += text;
          logContent += text;
        }
      }

      if (responseText.trim()) {
        logWrapper(`\n\n🟢 ${model} quota inception successful!`);
        summary.passed.push(model);
      } else {
        throw new Error("Received empty response from model.");
      }
    } catch (err) {
      logErrorWrapper(`\n\n🔴 ${model} quota inception failed: ${err.message}`);
      summary.failed.push({ model, reason: err.message });
    } finally {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
      await closeSessionSafe(session);
    }
  }

  logWrapper("\n=========================================");
  logWrapper("📊 8:00 AM Quota Trigger Summary");
  logWrapper("=========================================");
  logWrapper(`Total Models Scanned    : ${models.length}`);
  logWrapper(`🟢 Successful Resets    : ${summary.passed.length}`);
  if (summary.passed.length > 0) {
    logWrapper(`   - ${summary.passed.join("\n   - ")}`);
  }
  logWrapper(`🔴 Failed Resets        : ${summary.failed.length}`);
  if (summary.failed.length > 0) {
    logWrapper(
      `   - ${summary.failed.map((f) => `${f.model} (Reason: ${f.reason})`).join("\n   - ")}`,
    );
  }
  logWrapper("=========================================\n");

  await fs.writeFile(logFilePath, logContent, "utf-8");
  log(`\n📝 Full morning check-in logged to: ${logFilePath}`);

  if (summary.passed.length === 0 && models.length > 0) {
    process.exitCode = 1;
  }
}

if (
  process.argv[1] &&
  path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url))
) {
  main().catch((err) => {
    logError(`❌ Fatal Exercise Agents Error: ${err.stack || err.message}`);
    process.exit(1);
  });
}
