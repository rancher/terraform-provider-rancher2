import { flushLogs, runAgentSession } from "./agent-runner.js";

async function testIsolate() {
  console.log("Testing isolated mode...");
  try {
    const text = await runAgentSession({
      initialPrompt:
        'Respond with exactly the word "Hello". Do not use any tools.',
      requestedModel: "gemini-3.8-flash",
      isolate: true,
      standalone: true,
    });
    console.log("\nOutput from isolated agent:", text);
    console.log("\nIsolate test passed (agent ran without crashing).");
  } catch (err) {
    console.error("\nIsolate test failed:", err);
    process.exitCode = 1;
  } finally {
    await flushLogs();
  }
}

testIsolate();
