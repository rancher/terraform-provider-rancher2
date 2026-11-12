import assert from "node:assert";
import fs from "node:fs/promises";
import path from "node:path";
import test from "node:test";
import {
  checkSpelling,
  createCommit,
  createGitClient,
  DEFAULT_PRODUCT_ROOTS,
  fileAffectsProduct,
  getProductRoots,
  parseCommitArgs,
  resolveGlobalGitConfig,
  validateCommit,
  validateCommitMessage,
  validateCommitTitle,
} from "../scripts/commit.js";
import {
  createBranch,
  parseCreateBranchArgs,
} from "../scripts/create-branch.js";
import {
  getBranchChanges,
  parseBranchChangesArgs,
} from "../scripts/get_branch_changes.js";
import {
  getCommitLog,
  parseCommitLogArgs,
  parseCommitLogOptions,
} from "../scripts/get_commit_log.js";
import { getDeletedFile } from "../scripts/get_deleted_file.js";
import {
  getDiffToReview,
  parseDiffArgs,
} from "../scripts/get_diff_to_review.js";
import {
  getFilesChanged,
  parseFilesChangedArgs,
} from "../scripts/get_files_changed.js";
import { getModifiedPaths } from "../scripts/get_modified_paths.js";
import {
  getStagedFiles,
  parseStageArgs,
  parseStageOutput,
} from "../scripts/get_staged_files.js";
import {
  isExcludedByAIExclude,
  readIgnoredFile,
} from "../scripts/read-ignored-file.js";
import { stageFiles } from "../scripts/stage-files.js";
import {
  checkBranchExists,
  getCurrentBranch,
  switchBranch,
} from "../scripts/switch-branch.js";
import { unstageFiles } from "../scripts/unstage-files.js";

test("git skill - commit message validation", async (t) => {
  await t.test("rejects empty or missing commit messages", () => {
    const res1 = validateCommitMessage({ message: "", affectsProduct: true });
    assert.strictEqual(res1.valid, false);
    assert.match(res1.reason, /Empty commit message|Missing commit message/);

    const res2 = validateCommitMessage({
      message: "   ",
      affectsProduct: true,
    });
    assert.strictEqual(res2.valid, false);
    assert.match(res2.reason, /Empty commit message|Missing commit message/);
  });

  await t.test("rejects non-conventional commit messages", () => {
    const res = validateCommitMessage({
      message: "just a commit message",
      affectsProduct: true,
    });
    assert.strictEqual(res.valid, false);
    assert.match(
      res.reason,
      /Commit message does not start with the required prefix/,
    );
  });

  await t.test("rejects subject lines longer than 70 characters", () => {
    const longSubject = "fix: " + "a".repeat(70);
    const res = validateCommitMessage({
      message: longSubject,
      affectsProduct: true,
    });
    assert.strictEqual(res.valid, false);
    assert.match(res.reason, /less than 70 characters/);
  });

  await t.test(
    "enforces product boundary rules for non-product changes",
    () => {
      // Non-product changes (outside rancher2/) cannot use feat, refactor, or !
      const resFeat = validateCommitMessage({
        message: "feat: add new feature",
        affectsProduct: false,
      });
      assert.strictEqual(resFeat.valid, false);
      assert.match(resFeat.reason, /cannot use a product prefix/);

      const resRefactor = validateCommitMessage({
        message: "refactor: clean docs",
        affectsProduct: false,
      });
      assert.strictEqual(resRefactor.valid, false);
      assert.match(resRefactor.reason, /cannot use a product prefix/);

      const resBreaking = validateCommitMessage({
        message: "feat!: critical break",
        affectsProduct: false,
      });
      assert.strictEqual(resBreaking.valid, false);
      assert.match(resBreaking.reason, /cannot use a product prefix/);

      // Non-product changes CANNOT use fix: because fix: is a product prefix
      const resFix = validateCommitMessage({
        message: "fix: address typo in docs",
        affectsProduct: false,
      });
      assert.strictEqual(resFix.valid, false);
      assert.match(resFix.reason, /cannot use a product prefix/);

      // Non-product changes CAN use docs, ci, dev, env, test
      const resDocs = validateCommitMessage({
        message: "docs: address typo in readme",
        affectsProduct: false,
      });
      assert.strictEqual(resDocs.valid, true);

      const resCi = validateCommitMessage({
        message: "ci: update github action workflow",
        affectsProduct: false,
      });
      assert.strictEqual(resCi.valid, true);
    },
  );

  await t.test(
    "allows feat and breaking changes when product files are affected",
    () => {
      const resFeat = validateCommitMessage({
        message: "feat: new local provider option",
        affectsProduct: true,
      });
      assert.strictEqual(resFeat.valid, true);

      const resBreaking = validateCommitMessage({
        message: "feat!: breaking change to provider",
        affectsProduct: true,
      });
      assert.strictEqual(resBreaking.valid, true);
    },
  );
});

test("git skill - readIgnoredFile with line slicing", async (t) => {
  await t.test("throws an error if file does not exist", async () => {
    await assert.rejects(async () => {
      await readIgnoredFile({ targetPath: "non-existent-file-xyz.txt" });
    }, /not found/);
  });

  await t.test(
    "reads full file content when no lines are specified",
    async () => {
      const content = await readIgnoredFile({ targetPath: "package.json" });
      assert.match(content, /"devDependencies":/);
    },
  );

  await t.test(
    "slices specific lines when startLine and endLine are provided",
    async () => {
      const content = await readIgnoredFile({
        targetPath: "package.json",
        startLine: 1,
        endLine: 3,
      });
      const lines = content.split("\n");
      assert.strictEqual(lines.length, 3);
      assert.strictEqual(lines[0].trim(), "{");
    },
  );

  await t.test(
    "throws when targetPath is outside repository root",
    async () => {
      await assert.rejects(async () => {
        await readIgnoredFile({ targetPath: "../../outside-repo-file.txt" });
      }, /Cannot read file outside repository root/);
    },
  );

  await t.test("throws when startLine is greater than endLine", async () => {
    await assert.rejects(async () => {
      await readIgnoredFile({
        targetPath: "package.json",
        startLine: 10,
        endLine: 5,
      });
    }, /startLine .* cannot be greater than endLine/);
  });

  await t.test("throws when targetPath is a directory", async () => {
    await assert.rejects(async () => {
      await readIgnoredFile({ targetPath: "rancher2" });
    }, /is a directory, not a file/);
  });

  await t.test("throws when startLine or endLine is NaN", async () => {
    await assert.rejects(async () => {
      await readIgnoredFile({
        targetPath: "package.json",
        startLine: "invalid",
      });
    }, /Invalid startLine: 'invalid'/);

    await assert.rejects(async () => {
      await readIgnoredFile({
        targetPath: "package.json",
        startLine: 1,
        endLine: "invalid",
      });
    }, /Invalid endLine: 'invalid'/);
  });

  await t.test(
    "throws when targetPath is a symlink pointing outside repository root",
    async () => {
      const symlinkPath = path.resolve("temp-external-symlink.txt");
      try {
        // Point symlink to an external file
        await fs.symlink("/etc/hosts", symlinkPath);
        await assert.rejects(async () => {
          await readIgnoredFile({ targetPath: "temp-external-symlink.txt" });
        }, /Cannot read file outside repository root/);
      } finally {
        try {
          await fs.unlink(symlinkPath);
        } catch {
          // cleanup
        }
      }
    },
  );
});

test("git skill - getModifiedPaths parsing", async (t) => {
  await t.test("parses porcelain status format cleanly", async () => {
    const mockGit = {
      raw: async (args) => {
        assert.deepStrictEqual(args, [
          "-c",
          "core.quotepath=false",
          "status",
          "--porcelain",
        ]);
        return " M file1.js\n?? file2.txt\nR  old.js -> new.js\n";
      },
    };

    const paths = await getModifiedPaths({ git: mockGit });
    assert.deepStrictEqual(paths, ["file1.js", "file2.txt", "new.js"]);
  });

  await t.test("unquotes quoted renamed files", async () => {
    const mockGit = {
      raw: async () => 'R  "old name.js" -> "new name.js"\n',
    };
    const paths = await getModifiedPaths({ git: mockGit });
    assert.deepStrictEqual(paths, ["new name.js"]);
  });

  await t.test("captures both staged and unstaged file paths", async () => {
    const mockGit = {
      raw: async () =>
        "M  staged.js\n M unstaged.js\nA  added.js\n?? untracked.js\n",
    };
    const paths = await getModifiedPaths({ git: mockGit });
    assert.deepStrictEqual(paths, [
      "staged.js",
      "unstaged.js",
      "added.js",
      "untracked.js",
    ]);
  });
});

test("git skill - isExcludedByAIExclude fail closed", async (t) => {
  await t.test("returns false when .aiexclude does not exist", async () => {
    // Non-existent repo directory without .aiexclude
    const result = await isExcludedByAIExclude("/tmp/somefile.txt", "/tmp");
    assert.strictEqual(result, false);
  });
});

test("git skill - stageFiles unit tests", async (t) => {
  await t.test(
    "throws when no files are specified and all is false",
    async () => {
      const mockGit = {
        revparse: async () => "/mock/root",
      };
      await assert.rejects(async () => {
        await stageFiles({ git: mockGit, files: [], all: false });
      }, /No files specified to stage/);
    },
  );

  await t.test("stages specified files when not blocked", async () => {
    let capturedAdd = null;
    const mockGit = {
      revparse: async () => process.cwd(),
      add: async (files) => {
        capturedAdd = files;
      },
    };

    const staged = await stageFiles({ git: mockGit, files: ["package.json"] });
    assert.deepStrictEqual(staged, ["package.json"]);
    assert.deepStrictEqual(capturedAdd, ["--", "package.json"]);
  });

  await t.test(
    "deduplicates explicit file arguments when all is false",
    async () => {
      let capturedAdd = null;
      const mockGit = {
        revparse: async () => process.cwd(),
        add: async (files) => {
          capturedAdd = files;
        },
      };

      const staged = await stageFiles({
        git: mockGit,
        files: ["package.json", "package.json", "main.go", "package.json"],
      });
      assert.deepStrictEqual(staged, ["package.json", "main.go"]);
      assert.deepStrictEqual(capturedAdd, ["--", "package.json", "main.go"]);
    },
  );

  await t.test(
    "stages all tracked modified files when all is true",
    async () => {
      let capturedAdd = null;
      const mockGit = {
        revparse: async () => process.cwd(),
        raw: async (args) => {
          assert.deepStrictEqual(args, [
            "-c",
            "core.quotepath=false",
            "ls-files",
            "-m",
            "-d",
          ]);
          return "package.json\n";
        },
        add: async (files) => {
          capturedAdd = files;
        },
      };

      const staged = await stageFiles({ git: mockGit, all: true });
      assert.deepStrictEqual(staged, ["package.json"]);
      assert.deepStrictEqual(capturedAdd, ["--", "package.json"]);
    },
  );

  await t.test(
    "stages modified and untracked files when includeUntracked is true",
    async () => {
      let capturedAdd = null;
      const mockGit = {
        revparse: async () => process.cwd(),
        raw: async (args) => {
          assert.deepStrictEqual(args, [
            "-c",
            "core.quotepath=false",
            "ls-files",
            "-m",
            "-d",
            "-o",
            "--exclude-standard",
          ]);
          return "package.json\nnewfile.txt\n";
        },
        add: async (files) => {
          capturedAdd = files;
        },
      };

      const staged = await stageFiles({
        git: mockGit,
        all: true,
        includeUntracked: true,
      });
      assert.deepStrictEqual(staged, ["package.json", "newfile.txt"]);
      assert.deepStrictEqual(capturedAdd, [
        "--",
        "package.json",
        "newfile.txt",
      ]);
    },
  );
});

test("git skill - unstageFiles unit tests", async (t) => {
  await t.test(
    "throws when no files are specified and all is false",
    async () => {
      await assert.rejects(async () => {
        await unstageFiles({ files: [], all: false });
      }, /No files specified to unstage/);
    },
  );

  await t.test("unstages specified files with git.reset", async () => {
    let capturedReset = null;
    const mockGit = {
      reset: async (args) => {
        capturedReset = args;
      },
    };

    const unstaged = await unstageFiles({
      git: mockGit,
      files: ["package.json"],
    });
    assert.deepStrictEqual(unstaged, ["package.json"]);
    assert.deepStrictEqual(capturedReset, ["HEAD", "--", "package.json"]);
  });

  await t.test("unstages all staged files when all is true", async () => {
    let capturedReset = null;
    const mockGit = {
      raw: async (args) => {
        assert.deepStrictEqual(args, [
          "-c",
          "core.quotepath=false",
          "diff",
          "--name-only",
          "--cached",
        ]);
        return "package.json\nmain.go\n";
      },
      reset: async (args) => {
        capturedReset = args;
      },
    };

    const unstaged = await unstageFiles({ git: mockGit, all: true });
    assert.deepStrictEqual(unstaged, ["package.json", "main.go"]);
    assert.deepStrictEqual(capturedReset, ["HEAD", "--"]);
  });

  await t.test(
    "returns empty array when all is true but nothing is staged",
    async () => {
      const mockGit = {
        raw: async () => "",
      };
      const unstaged = await unstageFiles({ git: mockGit, all: true });
      assert.deepStrictEqual(unstaged, []);
    },
  );

  await t.test(
    "falls back to git rm --cached when HEAD reset throws",
    async () => {
      let capturedRaw = null;
      const mockGit = {
        reset: async () => {
          throw new Error("fatal: ambiguous argument HEAD");
        },
        raw: async (args) => {
          capturedRaw = args;
        },
      };

      const unstaged = await unstageFiles({
        git: mockGit,
        files: ["package.json"],
      });
      assert.deepStrictEqual(unstaged, ["package.json"]);
      assert.deepStrictEqual(capturedRaw, [
        "rm",
        "--cached",
        "--",
        "package.json",
      ]);
    },
  );

  await t.test(
    "falls back to ls-files when diff throws on unborn branch during unstage all",
    async () => {
      let capturedReset = null;
      const mockGit = {
        raw: async (args) => {
          if (args.includes("diff")) {
            throw new Error("fatal: ambiguous argument HEAD");
          }
          if (args.includes("ls-files")) {
            return "100644 abc 0\tunborn_file.js\n";
          }
          return "";
        },
        reset: async (args) => {
          capturedReset = args;
        },
      };

      const unstaged = await unstageFiles({ git: mockGit, all: true });
      assert.deepStrictEqual(unstaged, ["unborn_file.js"]);
      assert.deepStrictEqual(capturedReset, ["HEAD", "--"]);
    },
  );
});

test("git skill - commit.js spell checking and validateCommit", async (t) => {
  await t.test("checkSpelling passes clean message", async () => {
    const res = await checkSpelling("fix: update documentation and comments");
    assert.strictEqual(res.valid, true);
  });

  await t.test("checkSpelling catches misspelled words", async () => {
    const misspelledWord = "u" + "pdte";
    const res = await checkSpelling(`fix: ${misspelledWord} documentation`);
    // If cspell CLI is installed, it identifies words; if skipped, valid is true
    if (!res.skipped) {
      assert.strictEqual(res.valid, false);
      assert.match(res.reason, /Commit message contains spelling errors on:/);
    }
  });

  await t.test(
    "validateCommit validates format and spell check together",
    async () => {
      const misspelledWord = "inccor" + "rect";
      const mockCleanSpeller = async () => ({ valid: true });
      const mockDirtySpeller = async () => ({
        valid: false,
        reason: `Commit message contains spelling errors on: ${misspelledWord}`,
      });

      const resClean = await validateCommit({
        message: "fix: correct typo in provider",
        affectsProduct: true,
        checkSpellingFn: mockCleanSpeller,
      });
      assert.strictEqual(resClean.valid, true);

      const resDirty = await validateCommit({
        message: `fix: ${misspelledWord} behavior in provider`,
        affectsProduct: true,
        checkSpellingFn: mockDirtySpeller,
      });
      assert.strictEqual(resDirty.valid, false);
      assert.match(resDirty.reason, /spelling errors on:/);
    },
  );

  await t.test(
    "validateCommit allows standard merge commits without spell check",
    async () => {
      let spellerCalled = false;
      const mockSpeller = async () => {
        spellerCalled = true;
        return { valid: true };
      };

      const res = await validateCommit({
        message:
          "Merge branch 'main' of github.com:rancher/terraform-provider-rancher2",
        affectsProduct: true,
        checkSpellingFn: mockSpeller,
      });
      assert.strictEqual(res.valid, true);
      assert.strictEqual(spellerCalled, false);
    },
  );
});

test("git skill - parseCommitArgs unit tests", async (t) => {
  await t.test("returns empty string for empty args", () => {
    assert.strictEqual(parseCommitArgs([]), "");
    assert.strictEqual(parseCommitArgs(), "");
  });

  await t.test("parses positional words into space-separated string", () => {
    const msg = parseCommitArgs(["fix:", "update", "readme", "docs"]);
    assert.strictEqual(msg, "fix: update readme docs");
  });

  await t.test("parses single -m flag", () => {
    const msg = parseCommitArgs(["-m", "fix: update readme"]);
    assert.strictEqual(msg, "fix: update readme");
  });

  await t.test(
    "joins multiple -m flags with double newlines as paragraphs",
    () => {
      const msg = parseCommitArgs([
        "-m",
        "fix: update readme",
        "-m",
        "Added detailed explanation of skills and CLI options.",
      ]);
      assert.strictEqual(
        msg,
        "fix: update readme\n\nAdded detailed explanation of skills and CLI options.",
      );
    },
  );

  await t.test("parses -m=... and --message=... flags", () => {
    const msg = parseCommitArgs(["-m=fix: title", "--message=Paragraph 2"]);
    assert.strictEqual(msg, "fix: title\n\nParagraph 2");
  });
});

test("git skill - get_staged_files unit tests", async (t) => {
  await t.test("parseStageOutput parses git ls-files --stage output", () => {
    const sample =
      "100644 3716302c 0\t.gemini/settings.json\n" +
      "100644 a1b2c3d4 1\tfile_conflict.txt\n" +
      "100644 e5f6a7b8 2\tfile_conflict.txt\n" +
      '100644 99887766 0\t"file with spaces.txt"\n';

    const entries = parseStageOutput(sample);
    assert.strictEqual(entries.length, 4);
    assert.strictEqual(entries[0].path, ".gemini/settings.json");
    assert.strictEqual(entries[0].stage, 0);
    assert.strictEqual(entries[1].path, "file_conflict.txt");
    assert.strictEqual(entries[1].stage, 1);
    assert.strictEqual(entries[2].path, "file_conflict.txt");
    assert.strictEqual(entries[2].stage, 2);
    assert.strictEqual(entries[3].path, "file with spaces.txt");
    assert.strictEqual(entries[3].stage, 0);
  });

  await t.test(
    "getStagedFiles retrieves stages, statuses, and detects conflicts",
    async () => {
      const mockGit = {
        raw: async (args) => {
          if (args.includes("ls-files")) {
            return "100644 aabbcc 0\tfile1.js\n100755 ddeeff 1\tconflict.js\n100644 112233 0\tclean.js\n";
          }
          if (args.includes("diff")) {
            return "M\tfile1.js\nU\tconflict.js\n";
          }
          return "";
        },
      };

      const result = await getStagedFiles({ git: mockGit });
      assert.strictEqual(result.hasConflicts, true);
      assert.deepStrictEqual(result.paths, ["file1.js", "conflict.js"]);
      assert.strictEqual(result.entries.length, 2);
      assert.strictEqual(result.entries[0].status, "M");
      assert.strictEqual(result.entries[1].status, "U");
    },
  );

  await t.test("parseStageArgs extracts nameOnly flag", () => {
    assert.deepStrictEqual(parseStageArgs([]), { nameOnly: false });
    assert.deepStrictEqual(parseStageArgs(["--name-only"]), { nameOnly: true });
    assert.deepStrictEqual(parseStageArgs(["-n"]), { nameOnly: true });
  });
});

test("git skill - createBranch unit tests", async (t) => {
  await t.test("throws when branch name is empty", async () => {
    await assert.rejects(async () => {
      await createBranch({ branchName: "" });
    }, /Branch name cannot be empty/);
  });

  await t.test("throws when branch name starts with a hyphen", async () => {
    await assert.rejects(async () => {
      await createBranch({ branchName: "-b" });
    }, /Branch name cannot start with a hyphen/);
  });

  await t.test("throws when branch already exists", async () => {
    const mockGit = {
      branchLocal: async () => ({ all: ["main", "existing-branch"] }),
      getRemotes: async () => [{ name: "origin" }],
      fetch: async () => {},
      branch: async () => ({ current: "main" }),
      merge: async () => {},
    };

    await assert.rejects(async () => {
      await createBranch({ git: mockGit, branchName: "existing-branch" });
    }, /Branch 'existing-branch' already exists/);
  });

  await t.test("syncs upstream and checks out branch", async () => {
    let checkedOut = null;
    let syncedRemote = null;
    const mockGit = {
      branchLocal: async () => ({ all: ["main"] }),
      getRemotes: async () => [{ name: "upstream" }, { name: "origin" }],
      fetch: async (remote) => {
        syncedRemote = remote;
      },
      branch: async () => ({ current: "main" }),
      merge: async () => {},
      checkoutBranch: async (branchName, base) => {
        checkedOut = { branchName, base };
      },
    };

    const result = await createBranch({
      git: mockGit,
      branchName: "new-feature",
    });
    assert.strictEqual(result.branchName, "new-feature");
    assert.strictEqual(result.baseBranch, "main");
    assert.strictEqual(result.syncedRemote, "upstream");
    assert.strictEqual(syncedRemote, "upstream");
    assert.deepStrictEqual(checkedOut, {
      branchName: "new-feature",
      base: "main",
    });
  });

  await t.test(
    "strips remote prefix from base branch when syncing remote ref",
    async () => {
      let syncedRemote = null;
      let syncedRef = null;
      const mockGit = {
        branchLocal: async () => ({ all: ["main"] }),
        getRemotes: async () => [{ name: "upstream" }, { name: "origin" }],
        fetch: async (remote, ref) => {
          syncedRemote = remote;
          syncedRef = ref;
        },
        branch: async () => ({ current: "feature" }),
        checkoutBranch: async () => {},
      };

      const result = await createBranch({
        git: mockGit,
        branchName: "from-upstream",
        baseBranch: "upstream/main",
      });
      assert.strictEqual(result.branchName, "from-upstream");
      assert.strictEqual(syncedRemote, "upstream");
      assert.strictEqual(syncedRef, "main:main");
    },
  );

  await t.test(
    "dynamically identifies default base branch via symbolic-ref when omitted",
    async () => {
      let checkedOut = null;
      const mockGit = {
        raw: async (args) => {
          if (args[0] === "symbolic-ref") {
            return "origin/trunk\n";
          }
          return "";
        },
        branchLocal: async () => ({ all: ["feature", "trunk"] }),
        getRemotes: async () => [{ name: "origin" }],
        fetch: async () => {},
        branch: async () => ({ current: "trunk" }),
        merge: async () => {},
        checkoutBranch: async (branchName, base) => {
          checkedOut = { branchName, base };
        },
      };

      const result = await createBranch({
        git: mockGit,
        branchName: "feature-trunk",
      });
      assert.strictEqual(result.branchName, "feature-trunk");
      assert.strictEqual(result.baseBranch, "trunk");
      assert.deepStrictEqual(checkedOut, {
        branchName: "feature-trunk",
        base: "trunk",
      });
    },
  );

  await t.test(
    "aborts and throws error when synchronization fails",
    async () => {
      const mockGit = {
        branchLocal: async () => ({ all: ["main"] }),
        getRemotes: async () => [{ name: "upstream" }],
        fetch: async () => {
          throw new Error("Network unreachable");
        },
        branch: async () => ({ current: "main" }),
      };

      await assert.rejects(async () => {
        await createBranch({ git: mockGit, branchName: "failing-sync" });
      }, /Failed to sync base branch 'main' from 'upstream': Network unreachable/);
    },
  );
});

test("git skill - switchBranch unit tests", async (t) => {
  await t.test("throws when branch name is empty", async () => {
    await assert.rejects(async () => {
      await switchBranch({ branchName: "" });
    }, /Branch name cannot be empty/);
  });

  await t.test("throws when branch name starts with a hyphen", async () => {
    await assert.rejects(async () => {
      await switchBranch({ branchName: "--detach" });
    }, /Branch name cannot start with a hyphen/);
  });

  await t.test("throws when branch does not exist", async () => {
    const mockGit = {
      branch: async () => ({ all: ["main"] }),
    };

    await assert.rejects(async () => {
      await switchBranch({ git: mockGit, branchName: "non-existent" });
    }, /Branch 'non-existent' does not exist/);
  });

  await t.test("switches when branch exists locally or remotely", async () => {
    let checkedOut = null;
    const mockGit = {
      branch: async () => ({ all: ["main", "remotes/origin/feat-remote"] }),
      checkout: async (name) => {
        checkedOut = name;
      },
    };

    const result = await switchBranch({
      git: mockGit,
      branchName: "feat-remote",
    });
    assert.strictEqual(result.branchName, "feat-remote");
    assert.strictEqual(checkedOut, "feat-remote");
  });

  await t.test(
    "tracks upstream remote branch when branch only exists on upstream",
    async () => {
      let checkedOut = null;
      const mockGit = {
        branch: async () => ({
          all: ["main", "remotes/upstream/feat-upstream"],
        }),
        checkout: async (args) => {
          checkedOut = args;
        },
      };

      const result = await switchBranch({
        git: mockGit,
        branchName: "feat-upstream",
      });
      assert.strictEqual(result.branchName, "feat-upstream");
      assert.deepStrictEqual(checkedOut, [
        "-b",
        "feat-upstream",
        "--track",
        "upstream/feat-upstream",
      ]);
    },
  );

  await t.test(
    "normalizes remote prefixes to prevent detached HEAD checkout",
    async () => {
      let checkedOutOrigin = null;
      const mockGitOrigin = {
        branch: async () => ({ all: ["main", "remotes/origin/feat-1"] }),
        checkout: async (name) => {
          checkedOutOrigin = name;
        },
      };
      const resOrigin = await switchBranch({
        git: mockGitOrigin,
        branchName: "origin/feat-1",
      });
      assert.strictEqual(resOrigin.branchName, "feat-1");
      assert.strictEqual(checkedOutOrigin, "feat-1");

      let checkedOutUpstream = null;
      const mockGitUpstream = {
        branch: async () => ({ all: ["main", "remotes/upstream/feat-2"] }),
        checkout: async (args) => {
          checkedOutUpstream = args;
        },
      };
      const resUpstream = await switchBranch({
        git: mockGitUpstream,
        branchName: "remotes/upstream/feat-2",
      });
      assert.strictEqual(resUpstream.branchName, "feat-2");
      assert.deepStrictEqual(checkedOutUpstream, [
        "-b",
        "feat-2",
        "--track",
        "upstream/feat-2",
      ]);
    },
  );
});

test("git skill - getCurrentBranch unit tests", async (t) => {
  await t.test("returns current branch name from branch summary", async () => {
    const mockGit = {
      branch: async () => ({ current: "feature/ralph-loop" }),
    };
    const current = await getCurrentBranch({ git: mockGit });
    assert.strictEqual(current, "feature/ralph-loop");
  });

  await t.test(
    "returns empty string when current branch is missing",
    async () => {
      const mockGit = {
        branch: async () => ({ current: "" }),
      };
      const current = await getCurrentBranch({ git: mockGit });
      assert.strictEqual(current, "");
    },
  );
});

test("git skill - checkBranchExists unit tests", async (t) => {
  await t.test(
    "returns false when branchName is empty or missing",
    async () => {
      assert.strictEqual(await checkBranchExists({ branchName: "" }), false);
      assert.strictEqual(await checkBranchExists({ branchName: null }), false);
    },
  );

  await t.test("identifies existing local branch", async () => {
    const mockGit = {
      branch: async () => ({ all: ["main", "feat-branch"] }),
    };
    assert.strictEqual(
      await checkBranchExists({ git: mockGit, branchName: "feat-branch" }),
      true,
    );
  });

  await t.test("identifies existing remote tracking branch", async () => {
    const mockGit = {
      branch: async () => ({
        all: [
          "main",
          "remotes/origin/feat-remote",
          "remotes/upstream/feat-upstream",
        ],
      }),
    };
    assert.strictEqual(
      await checkBranchExists({ git: mockGit, branchName: "feat-remote" }),
      true,
    );
    assert.strictEqual(
      await checkBranchExists({
        git: mockGit,
        branchName: "origin/feat-remote",
      }),
      true,
    );
    assert.strictEqual(
      await checkBranchExists({ git: mockGit, branchName: "feat-upstream" }),
      true,
    );
    assert.strictEqual(
      await checkBranchExists({
        git: mockGit,
        branchName: "upstream/feat-upstream",
      }),
      true,
    );
    assert.strictEqual(
      await checkBranchExists({ git: mockGit, branchName: "refs/heads/main" }),
      true,
    );
    assert.strictEqual(
      await checkBranchExists({
        git: mockGit,
        branchName: "refs/remotes/origin/feat-remote",
      }),
      true,
    );
  });

  await t.test("returns false when branch does not exist", async () => {
    const mockGit = {
      branch: async () => ({ all: ["main"] }),
    };
    assert.strictEqual(
      await checkBranchExists({ git: mockGit, branchName: "missing" }),
      false,
    );
  });

  await t.test(
    "handles git branch errors safely by returning false",
    async () => {
      const mockGit = {
        branch: async () => {
          throw new Error("git error");
        },
      };
      assert.strictEqual(
        await checkBranchExists({ git: mockGit, branchName: "any" }),
        false,
      );
    },
  );
});

test("git skill - parseCreateBranchArgs unit tests", async (t) => {
  await t.test("parses positional branch and base branch", () => {
    const { branchName, baseBranch } = parseCreateBranchArgs([
      "my-feature",
      "main",
    ]);
    assert.strictEqual(branchName, "my-feature");
    assert.strictEqual(baseBranch, "main");
  });

  await t.test("parses --base and -b flags", () => {
    const res1 = parseCreateBranchArgs(["my-feature", "--base", "develop"]);
    assert.strictEqual(res1.branchName, "my-feature");
    assert.strictEqual(res1.baseBranch, "develop");

    const res2 = parseCreateBranchArgs(["my-feature", "-b", "develop"]);
    assert.strictEqual(res2.branchName, "my-feature");
    assert.strictEqual(res2.baseBranch, "develop");

    const res3 = parseCreateBranchArgs(["--base=develop", "my-feature"]);
    assert.strictEqual(res3.branchName, "my-feature");
    assert.strictEqual(res3.baseBranch, "develop");

    const res4 = parseCreateBranchArgs(["-b=develop", "my-feature"]);
    assert.strictEqual(res4.branchName, "my-feature");
    assert.strictEqual(res4.baseBranch, "develop");
  });

  await t.test("handles empty or missing arguments", () => {
    const res = parseCreateBranchArgs([]);
    assert.strictEqual(res.branchName, "");
    assert.strictEqual(res.baseBranch, undefined);
  });

  await t.test(
    "throws when --base or -b is provided as a trailing flag without a value",
    () => {
      assert.throws(
        () => parseCreateBranchArgs(["my-feature", "--base"]),
        /Flag '--base' requires a value\./,
      );
      assert.throws(
        () => parseCreateBranchArgs(["my-feature", "-b"]),
        /Flag '-b' requires a value\./,
      );
    },
  );
});

test("git skill - getBranchChanges unit tests", async (t) => {
  await t.test(
    "gets diff relative to base branch using merge-base",
    async () => {
      let capturedArgs = null;
      const mockGit = {
        branchLocal: async () => ({ all: ["main"] }),
        raw: async (args) => {
          if (args[0] === "merge-base") {
            return "abc1234\n";
          }
          if (args.includes("diff")) {
            throw new Error("fallback");
          }
          return "";
        },
        diff: async (args) => {
          capturedArgs = args;
          return "diff --git a/file.go b/file.go\n+new line\n";
        },
      };

      const diff = await getBranchChanges({ git: mockGit });
      assert.match(diff, /\+new line/);
      assert.deepStrictEqual(capturedArgs, [
        "--no-ext-diff",
        "--no-color",
        "abc1234",
        "--",
      ]);
    },
  );

  await t.test("uses specified baseBranch when provided", async () => {
    let capturedArgs = null;
    const mockGit = {
      raw: async (args) => {
        if (args[0] === "merge-base") {
          assert.strictEqual(args[1], "upstream/main");
          return "def5678\n";
        }
        if (args.includes("diff")) {
          throw new Error("fallback");
        }
        return "";
      },
      diff: async (args) => {
        capturedArgs = args;
        return "diff output";
      },
    };

    const diff = await getBranchChanges({
      git: mockGit,
      baseBranch: "upstream/main",
    });
    assert.strictEqual(diff, "diff output");
    assert.deepStrictEqual(capturedArgs, [
      "--no-ext-diff",
      "--no-color",
      "def5678",
      "--",
    ]);
  });

  await t.test(
    "automatically resolves trunk or dev as base branch from PREFERRED_DEFAULT_BRANCHES",
    async () => {
      let checkedBase = null;
      const mockGit = {
        branchLocal: async () => ({ all: ["feature", "trunk"] }),
        raw: async (args) => {
          if (args[0] === "merge-base") {
            checkedBase = args[1];
            return "trunk123\n";
          }
          return "";
        },
        diff: async () => "diff output",
      };

      await getBranchChanges({ git: mockGit });
      assert.strictEqual(checkedBase, "trunk");
    },
  );

  await t.test("supports optional files parameter with raw diff", async () => {
    let capturedArgs = null;
    const mockGit = {
      raw: async (args) => {
        if (args[0] === "merge-base") {
          return "abc1234\n";
        }
        if (args.includes("diff")) {
          capturedArgs = args;
          return "filtered diff";
        }
        return "";
      },
    };

    const diff = await getBranchChanges({
      git: mockGit,
      baseBranch: "main",
      files: ["file1.go", "file2.go"],
    });
    assert.strictEqual(diff, "filtered diff");
    assert.deepStrictEqual(capturedArgs, [
      "-c",
      "core.quotepath=false",
      "diff",
      "--no-ext-diff",
      "--no-color",
      "abc1234",
      "--",
      "file1.go",
      "file2.go",
    ]);
  });

  await t.test("filters out -- delimiter from files parameter", async () => {
    let capturedArgs = null;
    const mockGit = {
      raw: async (args) => {
        if (args[0] === "merge-base") {
          return "abc1234\n";
        }
        if (args.includes("diff")) {
          capturedArgs = args;
          return "filtered diff";
        }
        return "";
      },
    };

    await getBranchChanges({
      git: mockGit,
      baseBranch: "main",
      files: ["--", "file1.go", "--", "file2.go"],
    });
    assert.deepStrictEqual(capturedArgs, [
      "-c",
      "core.quotepath=false",
      "diff",
      "--no-ext-diff",
      "--no-color",
      "abc1234",
      "--",
      "file1.go",
      "file2.go",
    ]);
  });

  await t.test(
    "supports statOnly and nameOnly flags in getBranchChanges",
    async () => {
      const captured = [];
      const mockGit = {
        raw: async (args) => {
          if (args[0] === "merge-base") {
            return "abc1234\n";
          }
          captured.push(args);
          return "stat or name output";
        },
      };

      const statRes = await getBranchChanges({
        git: mockGit,
        baseBranch: "main",
        statOnly: true,
      });
      assert.strictEqual(statRes, "stat or name output");
      assert.ok(captured[0].includes("--stat"));

      captured.length = 0;
      const nameRes = await getBranchChanges({
        git: mockGit,
        baseBranch: "main",
        nameOnly: true,
      });
      assert.strictEqual(nameRes, "stat or name output");
      assert.ok(captured[0].includes("--name-only"));
    },
  );

  await t.test(
    "parseBranchChangesArgs extracts flags and options correctly",
    () => {
      assert.deepStrictEqual(parseBranchChangesArgs([]), {
        baseBranch: undefined,
        files: [],
        statOnly: false,
        nameOnly: false,
      });
      assert.deepStrictEqual(parseBranchChangesArgs(["--stat"]), {
        baseBranch: undefined,
        files: [],
        statOnly: true,
        nameOnly: false,
      });
      assert.deepStrictEqual(
        parseBranchChangesArgs(["--name-only", "-b", "upstream/main"]),
        {
          baseBranch: "upstream/main",
          files: [],
          statOnly: false,
          nameOnly: true,
        },
      );
      assert.deepStrictEqual(
        parseBranchChangesArgs(["--base=main", "file1.go", "file2.go"]),
        {
          baseBranch: "main",
          files: ["file1.go", "file2.go"],
          statOnly: false,
          nameOnly: false,
        },
      );
      assert.deepStrictEqual(parseBranchChangesArgs(["main", "file1.go"]), {
        baseBranch: "main",
        files: ["file1.go"],
        statOnly: false,
        nameOnly: false,
      });
      assert.deepStrictEqual(parseBranchChangesArgs(["--", "file1.go"]), {
        baseBranch: undefined,
        files: ["file1.go"],
        statOnly: false,
        nameOnly: false,
      });
      assert.deepStrictEqual(parseBranchChangesArgs(["GNUmakefile"]), {
        baseBranch: undefined,
        files: ["GNUmakefile"],
        statOnly: false,
        nameOnly: false,
      });
    },
  );
});

test("git skill - getCommitLog unit tests", async (t) => {
  await t.test("formats commit log with limit", async () => {
    let capturedLimit = null;
    const mockGit = {
      log: async (opts) => {
        capturedLimit = opts.maxCount;
        return {
          all: [
            {
              hash: "1234567890abcdef",
              message: "feat: add feature",
              date: "2026-10-15",
            },
            {
              hash: "abcdef1234567890",
              message: "fix: resolve bug",
              date: "2026-10-14",
            },
          ],
        };
      },
    };

    const commits = await getCommitLog({ git: mockGit, limit: 2 });
    assert.strictEqual(capturedLimit, 2);
    assert.strictEqual(commits.length, 2);
    assert.strictEqual(commits[0].hash, "1234567");
    assert.strictEqual(commits[0].message, "feat: add feature");
  });

  await t.test(
    "returns empty array when git log throws or repo has no commits",
    async () => {
      const mockGit = {
        log: async () => {
          throw new Error(
            "fatal: your current branch does not have any commits yet",
          );
        },
      };

      const commits = await getCommitLog({ git: mockGit });
      assert.deepStrictEqual(commits, []);
    },
  );

  await t.test(
    "parseCommitLogArgs supports positional, -n, and --limit flags",
    () => {
      assert.strictEqual(parseCommitLogArgs([]), 5);
      assert.strictEqual(parseCommitLogArgs(["10"]), 10);
      assert.strictEqual(parseCommitLogArgs(["-n", "7"]), 7);
      assert.strictEqual(parseCommitLogArgs(["--limit", "12"]), 12);
      assert.strictEqual(parseCommitLogArgs(["-n=8"]), 8);
      assert.strictEqual(parseCommitLogArgs(["--limit=15"]), 15);
    },
  );

  await t.test("supports file and branch scoping in getCommitLog", async () => {
    let capturedOpts = null;
    const mockGit = {
      log: async (opts) => {
        capturedOpts = opts;
        return {
          all: [{ hash: "1234567890", message: "test", date: "2026-10-15" }],
        };
      },
    };

    await getCommitLog({
      git: mockGit,
      limit: 3,
      file: "main.go",
      branch: "feat",
    });
    assert.strictEqual(capturedOpts.maxCount, 3);
    assert.strictEqual(capturedOpts.file, "main.go");
    assert.strictEqual(capturedOpts.from, "feat");
  });

  await t.test("parseCommitLogOptions parses file and branch options", () => {
    assert.deepStrictEqual(parseCommitLogOptions([]), {
      limit: 5,
      file: undefined,
      branch: undefined,
    });
    assert.deepStrictEqual(
      parseCommitLogOptions(["-n", "10", "-f", "file.go", "-b", "main"]),
      {
        limit: 10,
        file: "file.go",
        branch: "main",
      },
    );
    assert.deepStrictEqual(
      parseCommitLogOptions(["--file=main.go", "--branch=dev"]),
      {
        limit: 5,
        file: "main.go",
        branch: "dev",
      },
    );
  });
});

test("git skill - getDeletedFile unit tests", async (t) => {
  await t.test("throws when file path is missing", async () => {
    await assert.rejects(async () => {
      await getDeletedFile({ filePath: "" });
    }, /Missing file path/);
  });

  await t.test("throws when file was never found in git history", async () => {
    const mockGit = {
      log: async () => ({ all: [] }),
    };

    await assert.rejects(async () => {
      await getDeletedFile({ git: mockGit, filePath: "missing.go" });
    }, /never found in git history/);
  });

  await t.test("retrieves file contents from last commit show", async () => {
    const mockGit = {
      log: async () => ({ all: [{ hash: "commit-hash-123" }] }),
      show: async (args) => {
        assert.deepStrictEqual(args, ["commit-hash-123:old-file.txt"]);
        return "historical content";
      },
    };

    const content = await getDeletedFile({
      git: mockGit,
      filePath: "old-file.txt",
    });
    assert.strictEqual(content, "historical content");
  });

  await t.test(
    "retrieves file contents from specific commit when commit is provided",
    async () => {
      let capturedArg = null;
      const mockGit = {
        show: async (args) => {
          capturedArg = args[0];
          return "commit-specific content\n";
        },
      };

      const content = await getDeletedFile({
        git: mockGit,
        filePath: "file.go",
        commit: "abc1234",
      });
      assert.strictEqual(content, "commit-specific content\n");
      assert.strictEqual(capturedArg, "abc1234:file.go");
    },
  );
});

test("git skill - getFilesChanged unit tests", async (t) => {
  await t.test("returns status and branch diff text", async () => {
    const mockGit = {
      branchLocal: async () => ({ all: ["main"] }),
      raw: async (args) => {
        if (args.includes("status")) {
          return " M rancher2/provider.go\n";
        }
        if (args.includes("diff")) {
          return "M\trancher2/provider.go\n";
        }
        return "";
      },
    };

    const result = await getFilesChanged({ git: mockGit });
    assert.strictEqual(result.baseBranch, "main");
    assert.strictEqual(result.statusText, " M rancher2/provider.go\n");
    assert.strictEqual(result.branchDiffText, "M\trancher2/provider.go\n");
  });

  await t.test("supports explicit baseBranch and remote fallback", async () => {
    const mockGit = {
      branchLocal: async () => ({
        all: ["feature-branch", "remotes/origin/main"],
      }),
      raw: async (args) => {
        if (args.includes("status")) {
          return "";
        }
        if (args.includes("diff")) {
          return "A\tnew_file.txt\n";
        }
        return "";
      },
    };

    const resultExplicit = await getFilesChanged({
      git: mockGit,
      baseBranch: "upstream/main",
    });
    assert.strictEqual(resultExplicit.baseBranch, "upstream/main");
    assert.strictEqual(resultExplicit.branchDiffText, "A\tnew_file.txt\n");

    const resultFallback = await getFilesChanged({ git: mockGit });
    assert.strictEqual(resultFallback.baseBranch, "origin/main");

    const mockGitUpstream = {
      branchLocal: async () => ({
        all: ["feature-branch", "remotes/origin/main", "remotes/upstream/main"],
      }),
      raw: async () => "",
    };
    const resultUpstream = await getFilesChanged({ git: mockGitUpstream });
    assert.strictEqual(resultUpstream.baseBranch, "upstream/main");
  });

  await t.test(
    "automatically resolves trunk or dev as base branch from PREFERRED_DEFAULT_BRANCHES in getFilesChanged",
    async () => {
      let checkedBase = null;
      const mockGit = {
        branchLocal: async () => ({ all: ["feature", "trunk"] }),
        raw: async (args) => {
          if (args.includes("status")) {
            return "";
          }
          if (args.includes("merge-base")) {
            checkedBase = args[1];
            return "trunk123\n";
          }
          return "";
        },
      };

      const res = await getFilesChanged({ git: mockGit });
      assert.strictEqual(res.baseBranch, "trunk");
      assert.strictEqual(checkedBase, "trunk");
    },
  );

  await t.test(
    "resolves refs/remotes/ tracking branches in getFilesChanged",
    async () => {
      const mockGitRefs = {
        branchLocal: async () => ({
          all: ["feature", "refs/remotes/upstream/dev"],
        }),
        raw: async () => "",
      };
      const res = await getFilesChanged({ git: mockGitRefs });
      assert.strictEqual(res.baseBranch, "upstream/dev");
    },
  );

  await t.test(
    "parseFilesChangedArgs parses positional and flagged base branch",
    () => {
      assert.deepStrictEqual(parseFilesChangedArgs([]), {
        baseBranch: undefined,
      });
      assert.deepStrictEqual(parseFilesChangedArgs(["main"]), {
        baseBranch: "main",
      });
      assert.deepStrictEqual(parseFilesChangedArgs(["-b", "upstream/main"]), {
        baseBranch: "upstream/main",
      });
      assert.deepStrictEqual(parseFilesChangedArgs(["--base=main"]), {
        baseBranch: "main",
      });
      assert.deepStrictEqual(
        parseFilesChangedArgs(["--base", "origin/master"]),
        { baseBranch: "origin/master" },
      );
      assert.deepStrictEqual(
        parseFilesChangedArgs(["-b", "main", "extra_arg"]),
        { baseBranch: "main" },
      );
      assert.deepStrictEqual(parseFilesChangedArgs(["feature-branch"]), {
        baseBranch: undefined,
      });
      assert.deepStrictEqual(parseFilesChangedArgs(["--", "upstream/main"]), {
        baseBranch: undefined,
      });
    },
  );
});

test("git skill - getDiffToReview unit tests", async (t) => {
  await t.test(
    "calls git raw with --no-ext-diff, --no-color, and core.quotepath=false",
    async () => {
      const calls = [];
      const mockGit = {
        raw: async (args) => {
          calls.push(args);
          if (args.includes("--staged")) {
            return "staged diff output";
          }
          return "unstaged diff output";
        },
      };

      const result = await getDiffToReview({ git: mockGit });
      assert.strictEqual(result.staged, "staged diff output");
      assert.strictEqual(result.unstaged, "unstaged diff output");
      assert.strictEqual(calls.length, 2);
      assert.deepStrictEqual(calls[0], [
        "-c",
        "core.quotepath=false",
        "diff",
        "--no-ext-diff",
        "--no-color",
        "--staged",
        "--",
      ]);
      assert.deepStrictEqual(calls[1], [
        "-c",
        "core.quotepath=false",
        "diff",
        "--no-ext-diff",
        "--no-color",
        "--",
      ]);
    },
  );

  await t.test("passes specific file paths when provided", async () => {
    const calls = [];
    const mockGit = {
      raw: async (args) => {
        calls.push(args);
        return "file diff";
      },
    };

    const result = await getDiffToReview({
      git: mockGit,
      files: ["main.go", "README.md"],
    });
    assert.strictEqual(result.staged, "file diff");
    assert.deepStrictEqual(calls[0], [
      "-c",
      "core.quotepath=false",
      "diff",
      "--no-ext-diff",
      "--no-color",
      "--staged",
      "--",
      "main.go",
      "README.md",
    ]);
  });

  await t.test(
    "supports stagedOnly, unstagedOnly, statOnly, and nameOnly options",
    async () => {
      const calls = [];
      const mockGit = {
        raw: async (args) => {
          calls.push(args);
          return "diff data";
        },
      };

      const stagedRes = await getDiffToReview({
        git: mockGit,
        stagedOnly: true,
      });
      assert.strictEqual(stagedRes.staged, "diff data");
      assert.strictEqual(stagedRes.unstaged, "");
      assert.strictEqual(calls.length, 1);
      assert.ok(calls[0].includes("--staged"));

      calls.length = 0;
      const unstagedRes = await getDiffToReview({
        git: mockGit,
        unstagedOnly: true,
      });
      assert.strictEqual(unstagedRes.staged, "");
      assert.strictEqual(unstagedRes.unstaged, "diff data");
      assert.strictEqual(calls.length, 1);
      assert.ok(!calls[0].includes("--staged"));

      calls.length = 0;
      const statRes = await getDiffToReview({ git: mockGit, statOnly: true });
      assert.strictEqual(statRes.staged, "diff data");
      assert.strictEqual(statRes.unstaged, "diff data");
      assert.ok(calls[0].includes("--stat"));
      assert.ok(calls[1].includes("--stat"));

      calls.length = 0;
      const nameRes = await getDiffToReview({ git: mockGit, nameOnly: true });
      assert.strictEqual(nameRes.staged, "diff data");
      assert.strictEqual(nameRes.unstaged, "diff data");
      assert.ok(calls[0].includes("--name-only"));
      assert.ok(calls[1].includes("--name-only"));

      calls.length = 0;
      const headRes = await getDiffToReview({ git: mockGit, head: true });
      assert.strictEqual(headRes.head, "diff data");
      assert.strictEqual(headRes.staged, "");
      assert.strictEqual(headRes.unstaged, "");
      assert.ok(calls[0].includes("HEAD"));

      const unbornGit = {
        raw: async () => {
          throw new Error(
            'fatal: ambiguous argument "HEAD": unknown revision or path not in the working tree.',
          );
        },
      };
      const unbornRes = await getDiffToReview({ git: unbornGit, head: true });
      assert.strictEqual(unbornRes.head, "");

      const corruptGit = {
        raw: async () => {
          throw new Error("fatal: index file corrupt");
        },
      };
      await assert.rejects(async () => {
        await getDiffToReview({ git: corruptGit, head: true });
      }, /index file corrupt/);
    },
  );

  await t.test("parseDiffArgs separates flags and files correctly", () => {
    assert.deepStrictEqual(parseDiffArgs(["--staged"]), {
      files: [],
      stagedOnly: true,
      unstagedOnly: false,
      statOnly: false,
      nameOnly: false,
      head: false,
      raw: false,
    });
    assert.deepStrictEqual(
      parseDiffArgs(["--unstaged", "file1.txt", "file2.txt"]),
      {
        files: ["file1.txt", "file2.txt"],
        stagedOnly: false,
        unstagedOnly: true,
        statOnly: false,
        nameOnly: false,
        head: false,
        raw: false,
      },
    );
    assert.deepStrictEqual(parseDiffArgs(["--cached", "file1.txt"]), {
      files: ["file1.txt"],
      stagedOnly: true,
      unstagedOnly: false,
      statOnly: false,
      nameOnly: false,
      head: false,
      raw: false,
    });
    assert.deepStrictEqual(
      parseDiffArgs(["--staged", "--unstaged", "file1.txt"]),
      {
        files: ["file1.txt"],
        stagedOnly: false,
        unstagedOnly: false,
        statOnly: false,
        nameOnly: false,
        head: false,
        raw: false,
      },
    );
    assert.deepStrictEqual(parseDiffArgs(["--stat", "file1.txt"]), {
      files: ["file1.txt"],
      stagedOnly: false,
      unstagedOnly: false,
      statOnly: true,
      nameOnly: false,
      head: false,
      raw: false,
    });
    assert.deepStrictEqual(parseDiffArgs(["--name-only", "file1.txt"]), {
      files: ["file1.txt"],
      stagedOnly: false,
      unstagedOnly: false,
      statOnly: false,
      nameOnly: true,
      head: false,
      raw: false,
    });
    assert.deepStrictEqual(parseDiffArgs(["--head", "--raw", "file1.txt"]), {
      files: ["file1.txt"],
      stagedOnly: false,
      unstagedOnly: false,
      statOnly: false,
      nameOnly: false,
      head: true,
      raw: true,
    });
    assert.deepStrictEqual(
      parseDiffArgs(["--", "--special-file.txt", "file2.txt"]),
      {
        files: ["--special-file.txt", "file2.txt"],
        stagedOnly: false,
        unstagedOnly: false,
        statOnly: false,
        nameOnly: false,
        head: false,
        raw: false,
      },
    );
  });
});

test("git skill - createCommit unit tests", async (t) => {
  await t.test("throws when on a detached HEAD", async () => {
    const mockGit = {
      branch: async () => ({ current: "HEAD" }),
    };

    await assert.rejects(async () => {
      await createCommit({ git: mockGit, message: "fix: test message" });
    }, /Cannot commit on a detached HEAD/);
  });

  await t.test("throws when on protected branch main or master", async () => {
    const mockGitMain = {
      branch: async () => ({ current: "main" }),
    };
    await assert.rejects(async () => {
      await createCommit({ git: mockGitMain, message: "fix: test message" });
    }, /Commits directly to the 'main' branch are not allowed/);

    const mockGitMaster = {
      branch: async () => ({ current: "master" }),
    };
    await assert.rejects(async () => {
      await createCommit({ git: mockGitMaster, message: "fix: test message" });
    }, /Commits directly to the 'master' branch are not allowed/);
  });

  await t.test("throws when staging area is empty", async () => {
    const mockGit = {
      branch: async () => ({ current: "feature/test-branch" }),
      diffSummary: async (args) => {
        if (args && args.includes("--cached")) {
          return { files: [] };
        }
        return { files: [{ file: "modified.js" }] };
      },
    };

    await assert.rejects(async () => {
      await createCommit({ git: mockGit, message: "fix: test message" });
    }, /No staged changes found/);
  });

  await t.test(
    "creates signed commit when staged changes exist and validation passes",
    async () => {
      let capturedCommit = null;
      const mockGit = {
        branch: async () => ({ current: "feature/test-branch" }),
        diffSummary: async (args) => {
          if (args && args.includes("--cached")) {
            return { files: [{ file: "internal/provider.go" }] };
          }
          return { files: [] };
        },
        commit: async (msg, paths, opts) => {
          capturedCommit = { msg, paths, opts };
          return { commit: "abc1234", branch: "feature/test-branch" };
        },
      };

      const mockSpeller = async () => ({ valid: true });
      const result = await createCommit({
        git: mockGit,
        message: "feat: add support for new option",
        checkSpellingFn: mockSpeller,
      });

      assert.strictEqual(result.commit, "abc1234");
      assert.strictEqual(result.branch, "feature/test-branch");
      assert.strictEqual(
        capturedCommit.msg,
        "feat: add support for new option",
      );
      assert.deepStrictEqual(capturedCommit.opts, { "-s": null, "-S": null });
    },
  );

  await t.test(
    "rejects feat commit when staged file contains internal/ but does not start with internal/",
    async () => {
      const mockGit = {
        branch: async () => ({ current: "feature/test-branch" }),
        diffSummary: async (args) => {
          if (args && args.includes("--cached")) {
            return { files: [{ file: "docs/internal/example.md" }] };
          }
          return { files: [] };
        },
      };

      const mockSpeller = async () => ({ valid: true });
      await assert.rejects(async () => {
        await createCommit({
          git: mockGit,
          message: "feat: add support for new option",
          checkSpellingFn: mockSpeller,
        });
      }, /cannot use a product prefix/);
    },
  );
});

test("git skill - resolveGlobalGitConfig and createGitClient hermeticity", async (t) => {
  await t.test(
    "resolveGlobalGitConfig respects explicit non-devnull GIT_CONFIG_GLOBAL",
    () => {
      const res = resolveGlobalGitConfig({
        env: { GIT_CONFIG_GLOBAL: "/custom/isolated/.gitconfig" },
      });
      assert.strictEqual(res, "/custom/isolated/.gitconfig");
    },
  );

  await t.test(
    "resolveGlobalGitConfig falls back to userConfigPath when GIT_CONFIG_GLOBAL is /dev/null",
    () => {
      const mockFs = {
        existsSync: (p) => p === "/mock/home/.gitconfig",
      };
      const res = resolveGlobalGitConfig({
        userConfigPath: "/mock/home/.gitconfig",
        fsModule: mockFs,
        env: { GIT_CONFIG_GLOBAL: "/dev/null" },
      });
      assert.strictEqual(res, "/mock/home/.gitconfig");
    },
  );

  await t.test("resolveGlobalGitConfig returns path when exists", () => {
    const mockFs = {
      existsSync: (p) => p === "/mock/home/.gitconfig",
    };
    const res = resolveGlobalGitConfig({
      userConfigPath: "/mock/home/.gitconfig",
      fsModule: mockFs,
    });
    assert.strictEqual(res, "/mock/home/.gitconfig");
  });

  await t.test(
    "resolveGlobalGitConfig returns null when file does not exist",
    () => {
      const mockFs = {
        existsSync: () => false,
      };
      const res = resolveGlobalGitConfig({
        userConfigPath: "/mock/home/.gitconfig",
        fsModule: mockFs,
      });
      assert.strictEqual(res, null);
    },
  );

  await t.test(
    "createCommit uses mock git and does not call resolveGlobalGitConfig",
    async () => {
      let resolverCalled = false;
      const mockResolver = () => {
        resolverCalled = true;
        return "/mock/home/.gitconfig";
      };
      const mockGit = {
        branch: async () => ({ current: "feature/test-branch" }),
        diffSummary: async (args) => {
          if (args && args.includes("--cached")) {
            return { files: [{ file: "rancher2/provider.go" }] };
          }
          return { files: [] };
        },
        commit: async () => ({
          commit: "abc1234",
          branch: "feature/test-branch",
        }),
      };

      const result = await createCommit({
        git: mockGit,
        message: "fix: resolve issue",
        checkSpellingFn: async () => ({ valid: true }),
        resolveGlobalConfigFn: mockResolver,
      });

      assert.strictEqual(result.commit, "abc1234");
      assert.strictEqual(resolverCalled, false);
    },
  );

  await t.test(
    "createGitClient attaches GIT_CONFIG_GLOBAL to simpleGit env when provided",
    () => {
      const gitClient = createGitClient({
        globalConfigPath: "/mock/custom/.gitconfig",
      });
      assert.strictEqual(
        gitClient._executor?.env?.GIT_CONFIG_GLOBAL,
        "/mock/custom/.gitconfig",
      );
    },
  );

  await t.test(
    "createGitClient does not set GIT_CONFIG_GLOBAL when path is null",
    () => {
      const gitClient = createGitClient({
        globalConfigPath: null,
      });
      assert.strictEqual(
        gitClient._executor?.env?.GIT_CONFIG_GLOBAL,
        undefined,
      );
    },
  );
});

test("git skill - getProductRoots and fileAffectsProduct", async (t) => {
  await t.test(
    "returns DEFAULT_PRODUCT_ROOTS when settings file is missing or lacks project.productPaths",
    async () => {
      const roots = await getProductRoots("/non/existent/path");
      assert.deepStrictEqual(roots, DEFAULT_PRODUCT_ROOTS);
    },
  );

  await t.test(
    "fileAffectsProduct matches standard and custom roots correctly",
    () => {
      assert.strictEqual(fileAffectsProduct("rancher2/provider.go"), true);
      assert.strictEqual(fileAffectsProduct("internal/provider.go"), true);
      assert.strictEqual(
        fileAffectsProduct("rancher2/provider_test.go"),
        false,
      );
      assert.strictEqual(fileAffectsProduct("docs/README.md"), false);
      assert.strictEqual(
        fileAffectsProduct("custom/app.py", ["custom/"]),
        true,
      );
      assert.strictEqual(
        fileAffectsProduct("other/app.py", ["custom/"]),
        false,
      );
    },
  );

  await t.test(
    "validateCommitMessage uses shared productRoots and agrees with fileAffectsProduct",
    () => {
      const resCustomFeat = validateCommitMessage({
        message: "feat: add feature",
        changedFiles: ["custom/module.go"],
        productRoots: ["custom/"],
      });
      assert.strictEqual(resCustomFeat.valid, true);

      const resCustomDocs = validateCommitMessage({
        message: "docs: update docs",
        changedFiles: ["custom/module.go"],
        productRoots: ["custom/"],
      });
      assert.strictEqual(resCustomDocs.valid, false);
      assert.match(resCustomDocs.reason, /does not indicate a product change/);
    },
  );

  await t.test(
    "validateCommitMessage preserves product-directory test file boundaries and enforces test: prefix",
    () => {
      // Editing product directory test files requires test: prefix
      const resTest = validateCommitMessage({
        message: "test: update unit test for cluster",
        changedFiles: ["rancher2/resource_cluster_test.go"],
      });
      assert.strictEqual(resTest.valid, true);

      // docs: or ci: on product directory test file is rejected
      const resDocs = validateCommitMessage({
        message: "docs: update docstrings",
        changedFiles: ["rancher2/resource_cluster_test.go"],
      });
      assert.strictEqual(resDocs.valid, false);
      assert.match(resDocs.reason, /does not indicate a product change/);

      // fix: or feat: on product directory test file without product code is rejected
      const resFix = validateCommitMessage({
        message: "fix: update cluster test",
        changedFiles: ["rancher2/resource_cluster_test.go"],
      });
      assert.strictEqual(resFix.valid, false);
      assert.match(resFix.reason, /cannot use a product prefix/);
    },
  );

  await t.test(
    "validateCommitTitle aliases validateCommitMessage cleanly",
    () => {
      const res = validateCommitTitle("feat: add feature", true);
      assert.strictEqual(res.valid, true);
    },
  );
});
