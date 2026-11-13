#!/usr/bin/env node

import { Buffer } from "node:buffer";
import path from "node:path";
import { fileURLToPath } from "node:url";
import {
  DESTRUCTIVE_BINS,
  SHELL_WRAPPERS,
  MUTATING_GIT_SUBCOMMANDS,
  MUTATING_GH_VERBS,
  GH_FLAGS_CONSUMING_ARG,
  BLACKLIST_PATHS,
  TMUX_EXEC_SUBCOMMANDS,
  MUTATING_REMOTE_SUBCOMMANDS,
} from "./constants.js";

const originalWrite = process.stdout.write;
console.log = () => {};
console.info = () => {};
console.debug = () => {};
console.warn = () => {};
process.stdout.write = function (chunk, encoding, callback) {
  const cb = typeof encoding === "function" ? encoding : callback;
  if (typeof cb === "function") {
    cb();
  }
  return true;
};

function outputResult(obj) {
  originalWrite.call(process.stdout, JSON.stringify(obj) + "\n");
}

function isFlagConsumingArg(wrapper, flag) {
  if (wrapper === "xargs") {
    return (
      flag === "-I" ||
      flag === "-d" ||
      flag === "-E" ||
      flag === "-L" ||
      flag === "-n" ||
      flag === "-s" ||
      flag === "-a" ||
      flag === "-P" ||
      flag === "-R" ||
      flag === "-S" ||
      flag === "--delimiter" ||
      flag === "--max-args" ||
      flag === "--max-procs" ||
      flag === "--max-chars" ||
      flag === "--arg-file" ||
      flag === "--process-slot-var"
    );
  }
  if (wrapper === "sudo") {
    return (
      flag === "-u" ||
      flag === "-g" ||
      flag === "-C" ||
      flag === "-D" ||
      flag === "-h" ||
      flag === "-p" ||
      flag === "-U" ||
      flag === "-r" ||
      flag === "-t" ||
      flag === "--user" ||
      flag === "--group"
    );
  }
  if (wrapper === "env") {
    return (
      flag === "-u" || flag === "-C" || flag === "--chdir" || flag === "--unset"
    );
  }
  if (wrapper === "exec") {
    return flag === "-a";
  }
  if (wrapper === "timeout") {
    return (
      flag === "-k" ||
      flag === "--kill-after" ||
      flag === "-s" ||
      flag === "--signal"
    );
  }
  if (wrapper === "nice") {
    return flag === "-n" || flag === "--adjustment";
  }
  if (wrapper === "time") {
    return (
      flag === "-f" ||
      flag === "--format" ||
      flag === "-o" ||
      flag === "--output"
    );
  }
  if (wrapper === "watch") {
    return flag === "-n" || flag === "--interval";
  }
  if (wrapper === "su") {
    return flag === "-s" || flag === "--shell";
  }
  if (wrapper === "doas") {
    return flag === "-u" || flag === "-a" || flag === "-C";
  }
  if (wrapper === "chroot") {
    return flag === "--userspec" || flag === "--groups";
  }
  if (wrapper === "unshare") {
    return (
      flag === "--map-user" ||
      flag === "--map-group" ||
      flag === "--map-users" ||
      flag === "--map-groups" ||
      flag === "-S" ||
      flag === "--setuid" ||
      flag === "-G" ||
      flag === "--setgid" ||
      flag === "-R" ||
      flag === "--root" ||
      flag === "-w" ||
      flag === "--wd"
    );
  }
  if (wrapper === "screen") {
    return flag === "-c" || flag === "-S" || flag === "-p";
  }
  if (wrapper === "tmux") {
    return (
      flag === "-c" ||
      flag === "-S" ||
      flag === "-f" ||
      flag === "-L" ||
      flag === "-T" ||
      flag === "-s" ||
      flag === "-n" ||
      flag === "-t" ||
      flag === "-p"
    );
  }
  return false;
}

/**
 * Splits a compound command string into individual statements
 * separated by ;, &&, ||, |, &, or newlines, ignoring separators inside quotes.
 */
function splitStatements(cmdStr) {
  const statements = [];
  let current = "";
  let inSingleQuote = false;
  let inDoubleQuote = false;
  let escape = false;

  for (let i = 0; i < cmdStr.length; i++) {
    const char = cmdStr[i];

    if (escape) {
      if (char === "\n") {
        escape = false;
        continue;
      }
      current += char;
      escape = false;
      continue;
    }

    if (char === "\\" && !inSingleQuote) {
      escape = true;
      current += char;
      continue;
    }

    if (char === "'" && !inDoubleQuote) {
      inSingleQuote = !inSingleQuote;
      current += char;
      continue;
    }

    if (char === '"' && !inSingleQuote) {
      inDoubleQuote = !inDoubleQuote;
      current += char;
      continue;
    }

    if (!inSingleQuote && !inDoubleQuote) {
      if (char === "\n" || char === ";") {
        if (current.trim()) {
          statements.push(current.trim());
        }
        current = "";
        continue;
      }
      if (char === "&" || char === "|") {
        if (cmdStr[i + 1] === char) {
          if (current.trim()) {
            statements.push(current.trim());
          }
          current = "";
          i++;
          continue;
        } else {
          if (current.trim()) {
            statements.push(current.trim());
          }
          current = "";
          continue;
        }
      }
    }

    current += char;
  }

  if (current.trim()) {
    statements.push(current.trim());
  }

  return statements;
}

/**
 * Tokenizes a single command statement into words, respecting quotes.
 */
function tokenize(statement) {
  const tokens = [];
  let current = "";
  let inSingleQuote = false;
  let inDoubleQuote = false;
  let escape = false;

  for (let i = 0; i < statement.length; i++) {
    const char = statement[i];

    if (escape) {
      if (char === "\n") {
        escape = false;
        continue;
      }
      current += char;
      escape = false;
      continue;
    }

    if (char === "\\" && !inSingleQuote) {
      escape = true;
      continue;
    }

    if (char === "'" && !inDoubleQuote) {
      inSingleQuote = !inSingleQuote;
      continue;
    }

    if (char === '"' && !inSingleQuote) {
      inDoubleQuote = !inDoubleQuote;
      continue;
    }

    if (!inSingleQuote && !inDoubleQuote) {
      if (char === "{" && statement[i + 1] === "}") {
        current += "{}";
        i++;
        continue;
      }
      if (
        /\s/.test(char) ||
        char === "(" ||
        char === ")" ||
        char === "<" ||
        char === ">" ||
        char === "{" ||
        char === "}" ||
        char === ","
      ) {
        if (current) {
          tokens.push(current);
          current = "";
        }
        continue;
      }
    }

    current += char;
  }

  if (current) {
    tokens.push(current);
  }

  return tokens;
}

/**
 * Extracts the executable base name and recursively unwraps
 * shell interpreters (-c, -e) and command wrappers (command, sudo, env, etc.).
 * Returns a list of all identified executables to check.
 */
function extractExecutables(tokens, depth = 0) {
  if (depth > 5) {
    throw new Error("Recursion depth exceeded");
  }
  const executables = [];
  let idx = 0;

  while (idx < tokens.length) {
    const token = tokens[idx];

    // Skip environment variable assignments like FOO=bar
    if (/^[A-Za-z_][A-Za-z0-9_]*=/.test(token)) {
      idx++;
      continue;
    }

    const baseName = path.basename(token).toLowerCase();

    if (baseName.includes("$")) {
      executables.push({ executable: "rm", baseName: "rm", args: [] });
      return executables;
    }

    // If it's a wrapper like sudo, env, command, exec, bash -c, etc.
    if (SHELL_WRAPPERS.has(baseName)) {
      executables.push({
        executable: token,
        baseName,
        args: tokens.slice(idx + 1),
      });

      idx++;
      /**
       * Indicates whether a command payload was successfully found and unwrapped
       * from within the current wrapper's arguments. When true, the parser can safely
       * exit the inner argument loop and recursively process the unwrapped payload,
       * preventing the wrapper's own flags from being misidentified as bare executables.
       */
      let payloadFoundInWrapper = false;
      while (idx < tokens.length) {
        const flag = tokens[idx];
        if (/^[A-Za-z_][A-Za-z0-9_]*=/.test(flag)) {
          idx++;
          continue;
        }
        if (baseName === "su" && !flag.startsWith("-")) {
          const hasCmdFlag = tokens
            .slice(idx)
            .some(
              (t) =>
                Boolean(t.match(/^-[a-zA-Z]*?c/)) ||
                t === "--command" ||
                t.startsWith("--command="),
            );
          if (hasCmdFlag) {
            idx++;
            continue;
          }
        }
        const isShell =
          baseName === "bash" ||
          baseName === "sh" ||
          baseName === "zsh" ||
          baseName === "su" ||
          baseName === "dash" ||
          baseName === "ash" ||
          baseName === "fish" ||
          baseName === "ksh";
        const isNonShellInterpreter =
          baseName === "node" ||
          baseName === "python" ||
          baseName === "python2" ||
          baseName === "python3" ||
          baseName === "ruby" ||
          baseName === "perl";
        const shellClusterMatch =
          isShell && flag.match(/^-[a-zA-Z]*?c([\s\S]*)$/);
        const isSuCommand =
          baseName === "su" &&
          (flag === "--command" || flag.startsWith("--command="));
        const isNonShellEval =
          isNonShellInterpreter &&
          (flag.startsWith("-c") ||
            flag.startsWith("-e") ||
            flag.startsWith("--eval") ||
            flag.startsWith("-p"));
        const isEnvSplit =
          baseName === "env" &&
          (flag.startsWith("-S") || flag.startsWith("--split-string"));
        if (
          Boolean(shellClusterMatch) ||
          isNonShellEval ||
          isEnvSplit ||
          isSuCommand
        ) {
          const isAttached = shellClusterMatch
            ? shellClusterMatch[1].length > 0
            : isEnvSplit
              ? flag !== "-S" && flag !== "--split-string"
              : isSuCommand
                ? flag.startsWith("--command=")
                : flag !== "-c" &&
                  flag !== "-e" &&
                  flag !== "--eval" &&
                  flag !== "-p";
          if (!isAttached) {
            idx++;
          }
          if (idx < tokens.length || isAttached) {
            // For standard sh/bash/zsh/su we try to parse the payload as bash statements
            if (isShell || isEnvSplit || isSuCommand) {
              let payloadText;
              if (isAttached) {
                if (shellClusterMatch) {
                  payloadText = shellClusterMatch[1];
                } else if (flag.startsWith("--split-string=")) {
                  payloadText = flag.slice(15);
                } else if (flag.startsWith("-S")) {
                  payloadText = flag.slice(2);
                } else if (flag.startsWith("--command=")) {
                  payloadText = flag.slice(10);
                } else {
                  payloadText = flag;
                }
              } else {
                payloadText = tokens[idx];
              }
              const payloadStmts = splitStatements(payloadText);
              for (const stmt of payloadStmts) {
                const subExecs = extractExecutables(tokenize(stmt), depth + 1);
                executables.push(...subExecs);
              }
            } else {
              // For non-shell interpreters (node, python, perl, ruby) evaluating
              // arbitrary string payloads is unsafe since we can't parse their AST.
              // We add a synthetic "rm" executable to ensure it is denied.
              executables.push({ executable: "rm", baseName: "rm", args: [] });
            }
            payloadFoundInWrapper = true;
          }
          break; // Stop processing flags for this wrapper
        }

        if (baseName === "tmux" && TMUX_EXEC_SUBCOMMANDS.has(flag)) {
          idx++;
          continue;
        }
        if (baseName === "find") {
          let findIdx = idx;
          while (findIdx < tokens.length) {
            const fToken = tokens[findIdx];
            if (fToken === "-delete") {
              executables.push({ executable: "rm", baseName: "rm", args: [] });
            } else if (
              fToken === "-exec" ||
              fToken === "-execdir" ||
              fToken === "-ok" ||
              fToken === "-okdir"
            ) {
              const terminatorIdx = tokens
                .slice(findIdx + 1)
                .findIndex((t) => t === ";" || t === "+");
              if (terminatorIdx !== -1) {
                // Ensure we only process the payload up to the terminator
                const subExecs = extractExecutables(
                  tokens.slice(findIdx + 1, findIdx + 1 + terminatorIdx),
                  depth + 1,
                );
                executables.push(...subExecs);
                findIdx += terminatorIdx + 1;
              } else {
                // No terminator found, treat as invalid/denied to prevent payload bypass
                executables.push({
                  executable: "rm",
                  baseName: "rm",
                  args: [],
                });
                findIdx = tokens.length;
              }
            }
            findIdx++;
          }
          payloadFoundInWrapper = true;
          break;
        }

        // Check if it's a combined flag and argument, e.g., -Oid_rsa
        // We must check if the flag is a known consuming flag.
        if (flag.startsWith("-")) {
          // 1. Handle explicitly assigned flags like --output=file or -O=file
          if (flag.includes("=")) {
            const assignment = flag
              .slice(flag.indexOf("=") + 1)
              .replace(/^["']|["']$/g, "");
            const canonicalToken = path.posix.normalize(
              assignment.replace(/\\/g, "/").toLowerCase(),
            );
            if (isPathBlacklisted(canonicalToken)) {
              executables.push({ executable: "rm", baseName: "rm", args: [] });
              payloadFoundInWrapper = true;
              break;
            }
            idx++;
            continue;
          }

          // 2. Handle cases like -Oid_rsa where O is the flag and id_rsa is the argument
          // Check if it's a single-letter flag followed by an argument without space
          if (/^-[a-zA-Z](.+)$/.test(flag)) {
            const flagName = flag.slice(0, 2);
            const arg = flag.slice(2).replace(/^["']|["']$/g, "");

            // If the flag is one that consumes an argument, check the argument
            if (isFlagConsumingArg(baseName, flagName)) {
              const canonicalToken = path.posix.normalize(
                arg.replace(/\\/g, "/").toLowerCase(),
              );
              if (isPathBlacklisted(canonicalToken)) {
                executables.push({
                  executable: "rm",
                  baseName: "rm",
                  args: [],
                });
                payloadFoundInWrapper = true;
                break;
              }
              idx++;
              continue;
            }
          }

          // 3. Handle standard consuming flags with space (idx += 2)
          if (isFlagConsumingArg(baseName, flag)) {
            if (idx + 1 < tokens.length) {
              const arg = tokens[idx + 1].replace(/^["']|["']$/g, "");
              const canonicalToken = path.posix.normalize(
                arg.replace(/\\/g, "/").toLowerCase(),
              );
              if (isPathBlacklisted(canonicalToken)) {
                executables.push({
                  executable: "rm",
                  baseName: "rm",
                  args: [],
                });
                payloadFoundInWrapper = true;
                break;
              }
            }
            idx += 2;
            continue;
          }

          // Check if it is a flag that consumes argument, if so, we should consume the next token
          // Check if it is a flag that *might* be a blacklist path.
          // This seems to be the problematic part causing false negatives?
          // If flag starts with -, skip it but don't just increment idx once.
          idx++;
          continue;
        }
        if (baseName === "timeout" || baseName === "chroot") {
          idx++;
        }
        if (baseName === "tmux") {
          const remaining = tokens.slice(idx).join(" ");
          const payloadStmts = splitStatements(remaining);
          for (const stmt of payloadStmts) {
            const subExecs = extractExecutables(tokenize(stmt), depth + 1);
            executables.push(...subExecs);
          }
          payloadFoundInWrapper = true;
          break;
        }
        break; // Reached next executable token
      }

      if (payloadFoundInWrapper) {
        return executables;
      }
      continue; // Continue unwrapping outer shell wrapper
    }

    if (/\s/.test(token)) {
      const subTokens = tokenize(token);
      if (subTokens.length > 0) {
        executables.push(...extractExecutables(subTokens, depth + 1));
        break;
      }
    }

    // Standard executable found
    const args = tokens.slice(idx + 1);
    for (const rawToken of args) {
      let argToken = rawToken;
      if (argToken.startsWith("-")) {
        if (argToken.includes("=")) {
          argToken = argToken
            .slice(argToken.indexOf("=") + 1)
            .replace(/^["']|["']$/g, "");
        } else if (/^-[a-zA-Z](.+)$/.test(argToken)) {
          const candidate = argToken.slice(2).replace(/^["']|["']$/g, "");
          const candidateNormalized = path.posix
            .normalize(candidate.replace(/\\/g, "/").toLowerCase());
          if (isPathBlacklisted(candidateNormalized)) {
            argToken = candidate;
          } else {
            continue;
          }
        } else {
          continue;
        }
      }
      if (!argToken || argToken.startsWith("-")) {
        continue;
      }
      const canonicalToken = path.posix.normalize(
        argToken.replace(/\\/g, "/").toLowerCase(),
      );
      if (isPathBlacklisted(canonicalToken)) {
        executables.push({ executable: "rm", baseName: "rm", args: [] });
        return executables;
      }
    }

    executables.push({
      executable: token,
      baseName,
      args,
    });
    break;
  }

  return executables;
}

function isStatementDestructive(statement) {
  const tokens = tokenize(statement);
  if (tokens.length === 0) {
    return false;
  }
  try {
    const execs = extractExecutables(tokens);
    return execs.some((exec) => DESTRUCTIVE_BINS.has(exec.baseName));
  } catch {
    return true; // Deny if recursion depth exceeded
  }
}

function isStatementMutatingGit(statement) {
  const tokens = tokenize(statement);
  if (tokens.length === 0) {
    return false;
  }
  let execs;
  try {
    execs = extractExecutables(tokens);
  } catch {
    return true; // Deny if recursion depth exceeded
  }

  for (const { baseName, args } of execs) {
    if (baseName !== "git" && baseName !== "gh") {
      continue;
    }

    if (baseName === "gh") {
      let ghIdx = 0;
      while (ghIdx < args.length) {
        const arg = args[ghIdx];
        if (GH_FLAGS_CONSUMING_ARG.has(arg)) {
          ghIdx += 2;
          continue;
        }
        if (arg.startsWith("-")) {
          ghIdx++;
          continue;
        }
        const verb = arg.toLowerCase();
        if (
          verb.includes("$") ||
          MUTATING_GH_VERBS.has(verb) ||
          MUTATING_GIT_SUBCOMMANDS.has(verb)
        ) {
          return true;
        }
        ghIdx++;
      }
      continue;
    }

    let subCmdIdx = 0;
    while (subCmdIdx < args.length) {
      const arg = args[subCmdIdx];
      if (
        arg === "-C" ||
        arg === "-c" ||
        arg === "--git-dir" ||
        arg === "--work-tree"
      ) {
        subCmdIdx += 2;
        continue;
      }
      if (arg.startsWith("-")) {
        subCmdIdx++;
        continue;
      }
      break;
    }

    if (subCmdIdx < args.length) {
      const subCmd = args[subCmdIdx].toLowerCase();

      if (subCmd === "config") {
        const subArgs = args.slice(subCmdIdx + 1);
        const isQuery = subArgs.some(
          (a) =>
            a === "-l" ||
            a === "--list" ||
            a === "--get" ||
            a === "--get-all" ||
            a === "--get-regexp",
        );
        if (!isQuery) {
          return true;
        }
      } else if (subCmd === "remote") {
        const subArgs = args.slice(subCmdIdx + 1);
        const firstNonFlag = subArgs.find((a) => !a.startsWith("-"));
        if (
          firstNonFlag &&
          MUTATING_REMOTE_SUBCOMMANDS.has(firstNonFlag.toLowerCase())
        ) {
          return true;
        }
      } else if (subCmd.includes("$") || MUTATING_GIT_SUBCOMMANDS.has(subCmd)) {
        return true;
      }

      if (subCmd === "branch") {
        const subArgs = args.slice(subCmdIdx + 1);
        if (
          subArgs.some(
            (a) =>
              a === "-d" ||
              a === "-D" ||
              a === "--delete" ||
              a.startsWith("-d") ||
              a.startsWith("-D") ||
              a === "-m" ||
              a === "-M" ||
              a === "--move" ||
              a.startsWith("-m") ||
              a.startsWith("-M") ||
              a === "-c" ||
              a === "-C" ||
              a === "--copy" ||
              a.startsWith("-c") ||
              a.startsWith("-C") ||
              a === "-u" ||
              a === "--set-upstream-to" ||
              a.startsWith("--set-upstream-to=") ||
              a === "--unset-upstream" ||
              a === "--edit-description" ||
              a === "-f" ||
              a === "--force",
          )
        ) {
          return true;
        }

        const hasExplicitList = subArgs.some(
          (a) => a === "--list" || a.startsWith("--list="),
        );
        const BRANCH_OPTION_WITH_ARG = new Set([
          "--sort",
          "--format",
          "--points-at",
          "--merged",
          "--no-merged",
          "--contains",
          "--no-contains",
          "--abbrev",
        ]);
        let hasPositional = false;
        for (let i = 0; i < subArgs.length; i++) {
          const arg = subArgs[i];
          if (arg.startsWith("-")) {
            if (BRANCH_OPTION_WITH_ARG.has(arg)) {
              i++;
            }
            continue;
          }
          hasPositional = true;
          break;
        }

        if (hasPositional && !hasExplicitList) {
          return true;
        }
      }
      if (subCmd === "tag") {
        const subArgs = args.slice(subCmdIdx + 1);
        if (
          subArgs.some(
            (a) =>
              a === "-d" ||
              a === "-D" ||
              a === "--delete" ||
              a.startsWith("-d") ||
              a.startsWith("-D") ||
              a === "-a" ||
              a === "-s" ||
              a === "-f" ||
              a === "--force",
          )
        ) {
          return true;
        }
        const isList = subArgs.some((a) => a === "-l" || a === "--list");
        if (!isList && subArgs.some((a) => !a.startsWith("-"))) {
          return true;
        }
      }
    }
  }
  return false;
}

const ALLOWED_ENV_FILES = new Set([".envrc"]);

function isPathBlacklisted(target) {
  if (!target) {
    return false;
  }
  const canonicalTarget = path.posix.normalize(target);
  const normalizedTarget = canonicalTarget.endsWith("/")
    ? canonicalTarget
    : canonicalTarget + "/";

  // Check for exact match or sub-directory match
  for (const b of BLACKLIST_PATHS) {
    const canonicalBlacklist = path.posix.normalize(b);
    const normalizedBlacklist = canonicalBlacklist.endsWith("/")
      ? canonicalBlacklist
      : canonicalBlacklist + "/";

    if (
      normalizedTarget === normalizedBlacklist ||
      normalizedTarget.startsWith(normalizedBlacklist) ||
      (!normalizedBlacklist.startsWith("/") &&
        normalizedTarget.includes("/" + normalizedBlacklist))
    ) {
      return true;
    }
  }

  if (/(?:^|\/)\.git(?:$|\/)/.test(canonicalTarget)) {
    return true;
  }

  if (/(?:^|\/)\.ssh(?:$|\/)/.test(canonicalTarget)) {
    return true;
  }

  if (/(?:^|\/)id_(?:rsa|dsa|ecdsa|ed25519)(?:$|[._-])/.test(path.basename(canonicalTarget))) {
    return true;
  }

  if (
    /(?:^|\/)\.gemini\/settings\.json(?:$|\/)/.test(canonicalTarget) ||
    /(?:^|\/)block-restricted-commands\.js(?:$|\/)/.test(canonicalTarget)
  ) {
    return true;
  }
  // Check for ".env" as an exact filename or path segment (e.g. .env, dir/.env, .env/file, .env.local).
  if (/(?:^|\/)\.env(?:\.[^/]*|-[^/]*)*(?:$|\/)/.test(canonicalTarget)) {
    if (ALLOWED_ENV_FILES.has(path.basename(canonicalTarget))) {
      return false;
    }
    return true;
  }
  return false;
}

/**
 * Checks whether a tool call is safe or should be blocked.
 * Fast, 100% deterministic, zero-network evaluation.
 */
function evaluateRules(tool_name, tool_input) {
  let cmdStr = "";
  if (
    tool_name === "run_shell_command" &&
    tool_input &&
    typeof tool_input.command === "string"
  ) {
    cmdStr = tool_input.command.trim();
  }

  // 0. Check for command substitutions (which could bypass static analysis)
  if (cmdStr) {
    let inSingleQuote = false;
    let inDoubleQuote = false;
    let escape = false;
    for (let i = 0; i < cmdStr.length; i++) {
      const char = cmdStr[i];
      if (escape) {
        escape = false;
        continue;
      }
      if (char === "\\" && !inSingleQuote) {
        escape = true;
        continue;
      }
      if (char === '"' && !inSingleQuote) {
        inDoubleQuote = !inDoubleQuote;
        continue;
      }
      if (char === "'" && !inDoubleQuote) {
        inSingleQuote = !inSingleQuote;
        continue;
      }
      if (!inSingleQuote) {
        if (char === "`") return "deny";
        if (
          (char === "$" || char === "<" || char === ">") &&
          cmdStr[i + 1] === "("
        )
          return "deny";
      }
    }
  }

  // 1. Check destructive shell commands
  if (
    cmdStr &&
    splitStatements(cmdStr).some((stmt) => isStatementDestructive(stmt))
  ) {
    return "deny";
  }

  // 2. Check mutating git commands
  if (
    cmdStr &&
    splitStatements(cmdStr).some((stmt) => isStatementMutatingGit(stmt))
  ) {
    return "deny";
  }

  // 3. Check sensitive file paths - canonicalize paths to defeat traversal
  const rawTargetPath = (
    tool_input?.file_path ||
    tool_input?.path ||
    tool_input?.dir_path ||
    ""
  ).replace(/\\/g, "/");

  if (rawTargetPath) {
    const canonicalTarget = path.posix.normalize(rawTargetPath.toLowerCase());
    if (isPathBlacklisted(canonicalTarget)) {
      return "deny";
    }
  }

  if (cmdStr) {
    let execs;
    try {
      execs = splitStatements(cmdStr).flatMap((stmt) =>
        extractExecutables(tokenize(stmt)),
      );
    } catch {
      return "deny";
    }
    for (const exec of execs) {
      const canonicalExec = path.posix.normalize(
        exec.executable.replace(/\\/g, "/").toLowerCase(),
      );
      const execStripped = canonicalExec.replace(
        /^\/usr\/(?:local\/)?(?:bin|sbin)\//,
        "",
      );
      if (isPathBlacklisted(execStripped)) {
        return "deny";
      }
      for (const rawArg of exec.args) {
        let arg = rawArg;
        if (arg.startsWith("-")) {
          if (arg.includes("=")) {
            arg = arg.slice(arg.indexOf("=") + 1).replace(/^["']|["']$/g, "");
          } else if (/^-[a-zA-Z](.+)$/.test(arg)) {
            const candidate = arg.slice(2).replace(/^["']|["']$/g, "");
            const candidateNormalized = path.posix.normalize(
              candidate.replace(/\\/g, "/").toLowerCase(),
            );
            if (isPathBlacklisted(candidateNormalized)) {
              arg = candidate;
            } else {
              continue;
            }
          } else {
            continue;
          }
        }
        if (!arg || arg.startsWith("-")) {
          continue;
        }
        const canonicalArg = path.posix.normalize(
          arg.replace(/\\/g, "/").toLowerCase(),
        );
        if (isPathBlacklisted(canonicalArg)) {
          return "deny";
        }
      }
    }
  }

  // 4. Block unauthorized subagent spawning via hook
  if (tool_name === "invoke_agent") {
    return "deny";
  }

  // 5. Block external GitHub PR/issue web fetches
  if (tool_name === "web_fetch" && tool_input) {
    const targetStr =
      `${tool_input.url || ""} ${tool_input.prompt || ""}`.toLowerCase();
    if (
      targetStr.includes("github.com") &&
      (targetStr.includes("/pull") || targetStr.includes("/issues"))
    ) {
      return "deny";
    }
  }

  return "allow";
}

async function main() {
  let inputData;
  try {
    const buffers = [];
    for await (const chunk of process.stdin) {
      buffers.push(chunk);
    }
    const rawData = Buffer.concat(buffers).toString("utf-8");
    inputData = JSON.parse(rawData);
  } catch {
    outputResult({
      decision: "deny",
      reason: "Failed to parse input parameters.",
    });
    process.exit(0);
  }

  const { tool_name, tool_input } = inputData;
  const decision = evaluateRules(tool_name, tool_input);

  if (decision === "deny") {
    const proposedCall = `${tool_name}(${JSON.stringify(tool_input || {})})`;
    outputResult({
      decision: "deny",
      reason: `tool call ${proposedCall} was detected as potentially destructive`,
    });
    process.exit(0);
  }

  outputResult({ decision: "allow" });
  process.exit(0);
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  main().catch(() => {
    // Fail-open on unexpected crash to avoid locking developer workflow
    outputResult({ decision: "allow" });
    process.exit(0);
  });
}
