#!/usr/bin/env node
import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { GeminiCliAgent } from "@google/gemini-cli-sdk";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Open/Closed Principle compliant: Load models dynamically via environment variables with a safe fallback
const models = process.env.GEMINI_EXERCISE_MODELS
  ? process.env.GEMINI_EXERCISE_MODELS.split(",").map((m) => m.trim())
  : [
    "gemini-3.1-pro-preview",
    "gemini-3.8-flash",
    "gemini-2.5-pro",
    "gemini-3.1-flash-lite",
  ];

async function main() {
  const tmpDir = path.resolve(__dirname, "..", "tmp");
  await fs.mkdir(tmpDir, { recursive: true });
  const logFilePath = path.join(tmpDir, "exercise-agents.log");
  let logContent = `Gemini Models Daily Quota Inception Log - ${new Date().toISOString()}\n=========================================\n`;

  const log = (msg) => {
    console.log(msg);
    logContent += msg + "\n";
  };
  const logError = (msg) => {
    console.error(msg);
    logContent += msg + "\n";
  };

  log(
    `📅 Morning Check-in: Initiating quota resets for available Gemini CLI models...`,
  );
  log(`Models to trigger: ${models.join(", ")}`);

  const summary = {
    passed: [],
    failed: [],
  };

  for (const model of models) {
    log(`\n=========================================`);
    log(`☀️ Triggering Model: ${model}`);
    log(`=========================================`);

    let timeoutId;
    try {
      // Instantiate fast agent using the SDK with the target model
      const agent = new GeminiCliAgent({
        model: model,
        instructions: "You are a polite, helpful assistant.",
      });

      const controller = new globalThis.AbortController();
      timeoutId = setTimeout(() => {
        controller.abort(new Error("Model quota exercise request timed out after 60s."));
      }, 60000);
      const greetingPrompt = `Good morning! It is 8:00 AM. Just saying a quick hello to confirm you are awake and ready for the day!`;

      const session = agent.session();
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
        log(`\n\n🟢 ${model} quota inception successful!`);
        summary.passed.push(model);
      } else {
        throw new Error("Received empty response from model.");
      }
    } catch (err) {
      logError(`\n\n🔴 ${model} quota inception failed: ${err.message}`);
      summary.failed.push({ model, reason: err.message });
    } finally {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    }
  }

  log("\n=========================================");
  log("📊 8:00 AM Quota Trigger Summary");
  log("=========================================");
  log(`Total Models Scanned    : ${models.length}`);
  log(`🟢 Successful Resets    : ${summary.passed.length}`);
  if (summary.passed.length > 0) {
    log(`   - ${summary.passed.join("\n   - ")}`);
  }
  log(`🔴 Failed Resets        : ${summary.failed.length}`);
  if (summary.failed.length > 0) {
    log(
      `   - ${summary.failed.map((f) => `${f.model} (Reason: ${f.reason})`).join("\n   - ")}`,
    );
  }
  log("=========================================\n");

  await fs.writeFile(logFilePath, logContent, "utf-8");
  console.log(`\n📝 Full morning check-in logged to: ${logFilePath}`);

  if (summary.passed.length === 0 && models.length > 0) {
    process.exitCode = 1;
  }
}

main().catch((err) => {
  console.error(`❌ Fatal Exercise Agents Error: ${err.stack || err.message}`);
  process.exit(1);
});
