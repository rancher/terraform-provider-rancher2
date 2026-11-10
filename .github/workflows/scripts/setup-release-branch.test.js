import assert from "node:assert/strict";
import { execFileSync } from "node:child_process";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { describe, it } from "node:test";

describe("setup-release-branch.sh", () => {
  const scriptPath = path.resolve(
    process.cwd(),
    ".github/workflows/scripts/setup-release-branch.sh",
  );

  it("passes bash syntax check (bash -n)", () => {
    assert.doesNotThrow(() => {
      execFileSync("bash", ["-n", scriptPath], { stdio: "pipe" });
    });
  });

  it("fails with error message when branch argument is invalid", () => {
    assert.throws(
      () => {
        execFileSync("bash", [scriptPath, "invalid-branch-name"], {
          stdio: "pipe",
        });
      },
      (error) => {
        assert.match(
          error.stderr.toString(),
          /does not match a recognized release branch format/,
        );
        return true;
      },
    );
  });

  it("configures stubs and release-please versions correctly in a target directory", () => {
    const tmpDir = fs.mkdtempSync(
      path.join(os.tmpdir(), "setup-release-branch-test-"),
    );
    try {
      // Create test environment
      const workflowsDir = path.join(tmpDir, ".github", "workflows");
      fs.mkdirSync(workflowsDir, { recursive: true });

      // Create dummy workflow files (both legacy and arbitrary future workflows)
      fs.writeFileSync(path.join(workflowsDir, "manual-release.yml"), "legacy");
      fs.writeFileSync(
        path.join(workflowsDir, "manual-rc-release.yml"),
        "legacy",
      );
      fs.writeFileSync(path.join(workflowsDir, "custom-ci.yml"), "future");
      fs.writeFileSync(
        path.join(workflowsDir, "arbitrary-workflow.yaml"),
        "future",
      );

      // Create initial release-please-config.json
      const initialConfig = {
        "release-as": "v15.3.0",
        packages: { ".": { "release-type": "go" } },
      };
      fs.writeFileSync(
        path.join(tmpDir, "release-please-config.json"),
        JSON.stringify(initialConfig, null, 2),
      );

      // Run setup-release-branch.sh inside tmpDir with branch "release/v16"
      execFileSync("bash", [scriptPath, "release/v16"], {
        cwd: tmpDir,
        stdio: "pipe",
      });

      // 1. Verify release-please-config.json
      const updatedConfig = JSON.parse(
        fs.readFileSync(
          path.join(tmpDir, "release-please-config.json"),
          "utf8",
        ),
      );
      assert.strictEqual(updatedConfig["release-as"], "v16.0.0");

      // 2. Verify .release-please-manifest.json
      const updatedManifest = JSON.parse(
        fs.readFileSync(
          path.join(tmpDir, ".release-please-manifest.json"),
          "utf8",
        ),
      );
      assert.strictEqual(updatedManifest["."], "16.0.0");

      // 3. Verify non-stub workflows removed
      assert.strictEqual(
        fs.existsSync(path.join(workflowsDir, "manual-release.yml")),
        false,
      );
      assert.strictEqual(
        fs.existsSync(path.join(workflowsDir, "manual-rc-release.yml")),
        false,
      );
      assert.strictEqual(
        fs.existsSync(path.join(workflowsDir, "custom-ci.yml")),
        false,
      );
      assert.strictEqual(
        fs.existsSync(path.join(workflowsDir, "arbitrary-workflow.yaml")),
        false,
      );

      // 4. Verify pull_request.yaml caller stub
      const prStub = fs.readFileSync(
        path.join(workflowsDir, "pull_request.yaml"),
        "utf8",
      );
      assert.match(
        prStub,
        /uses: rancher\/terraform-provider-rancher2\/\.github\/workflows\/pull_request\.yaml@main/,
      );
      assert.match(prStub, /merge_group:/);
      assert.match(prStub, /contents:\s+read/);
      assert.match(prStub, /pull-requests:\s+read/);

      // 5. Verify release.yml caller stub
      const releaseStub = fs.readFileSync(
        path.join(workflowsDir, "release.yml"),
        "utf8",
      );
      assert.match(
        releaseStub,
        /uses: rancher\/terraform-provider-rancher2\/\.github\/workflows\/release-please\.yml@main/,
      );
      assert.match(releaseStub, /branches:\s+- 'release\/v\*'/);
      assert.match(releaseStub, /contents:\s+write/);
      assert.match(releaseStub, /id-token:\s+write/);
      assert.match(releaseStub, /issues:\s+write/);
      assert.match(releaseStub, /pull-requests:\s+write/);
      assert.match(releaseStub, /actions:\s+read/);

      // 6. Verify e2e-pr-listener.yml is created/preserved
      assert.strictEqual(
        fs.existsSync(path.join(workflowsDir, "e2e-pr-listener.yml")),
        true,
      );
      const listenerContent = fs.readFileSync(
        path.join(workflowsDir, "e2e-pr-listener.yml"),
        "utf8",
      );
      assert.match(listenerContent, /name:\s+PR Executor Listener/);
      assert.match(listenerContent, /branches:\s+- 'release\/v\*'/);
    } finally {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    }
  });

  it('normalizes bare version like "v17" to v17.0.0', () => {
    const tmpDir = fs.mkdtempSync(
      path.join(os.tmpdir(), "setup-release-branch-test-"),
    );
    try {
      execFileSync("bash", [scriptPath, "v17"], {
        cwd: tmpDir,
        stdio: "pipe",
      });
      const updatedConfig = JSON.parse(
        fs.readFileSync(
          path.join(tmpDir, "release-please-config.json"),
          "utf8",
        ),
      );
      assert.strictEqual(updatedConfig["release-as"], "v17.0.0");
    } finally {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    }
  });
});
