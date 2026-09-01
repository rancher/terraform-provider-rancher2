#!/usr/bin/env node

import { execFileSync } from 'child_process';
import fs from 'fs';
import path from 'path';

// Parse command-line arguments
const args = process.argv.slice(2);
let branch = 'release/v15';
let startRef = null;

for (let i = 0; i < args.length; i++) {
  if ((args[i] === '--branch' || args[i] === '-b') && args[i + 1]) {
    branch = args[i + 1];
    i++;
  } else if ((args[i] === '--start' || args[i] === '-s') && args[i + 1]) {
    startRef = args[i + 1];
    i++;
  }
}

console.log(`=========================================`);
console.log(`🔍 Backport QA Tracker Audit (GitHub API)`);
console.log(`Target Branch: ${branch}`);
console.log(`=========================================\n`);

// 1. Determine starting compare reference
const baseCompare = startRef || 'main';
console.log(`📌 [1/4] Comparing starting reference '${baseCompare}' with head branch '${branch}'...`);

// 2. Fetch unique commits using GitHub Compare API
console.log(`📋 [2/4] Fetching unique commits on '${branch}' since '${baseCompare}' from GitHub API...`);
let commits = [];
try {
  const apiPath = `repos/rancher/terraform-provider-rancher2/compare/${baseCompare}...${branch}`;
  const responseJson = execFileSync('gh', ['api', apiPath], { encoding: 'utf-8', stdio: ['ignore', 'pipe', 'pipe'] }).trim();
  const response = JSON.parse(responseJson);
  
  if (response && response.commits) {
    // Fix 1: Fail fast if GitHub Compare API truncated the returned commits array
    if (response.total_commits > response.commits.length) {
      console.error(`❌ [Error] GitHub API compare endpoint truncated the commits list (returned ${response.commits.length} of ${response.total_commits} commits). Rerun with --start/-s to narrow the range so the audit is complete.`);
      process.exit(2);
    }

    commits = response.commits.map(c => {
      const sha = c.sha;
      const fullMessage = (c.commit && c.commit.message) || 'No message';
      const subject = fullMessage.split('\n')[0];
      return { sha, subject };
    });
  }
} catch (err) {
  // Fix 2: Explicitly detect when GitHub CLI is not installed (ENOENT)
  if (err.code === 'ENOENT') {
    console.error(`❌ [Error] The GitHub CLI ('gh') was not found on your system. Please ensure 'gh' is installed and authenticated to run this script.`);
    process.exit(1);
  }
  
  // Fix 3: Log standard error stream from the CLI for rich debugging
  const errMsg = err.stderr ? err.stderr.toString().trim() : err.message;
  console.error(`❌ [Error] Failed to fetch compare logs from GitHub REST API: ${errMsg}`);
  console.error(`Please verify that you have internet access and that the gh CLI is authenticated.`);
  process.exit(1);
}

console.log(`   Found ${commits.length} commits to audit.\n`);

if (commits.length === 0) {
  console.log(`🎉 No commits found since ${baseCompare}. Audit complete!`);
  process.exit(0);
}

// 3. Trace each commit back to its PR and QA tracking issue
console.log(`🛰️ [3/4] Auditing commits against GitHub API (this may take a moment)...`);
const results = [];

for (const commit of commits) {
  process.stdout.write(`   - Auditing ${commit.sha.substring(0, 8)}... `);
  let pr = null;
  try {
    const prsJson = execFileSync('gh', ['api', `repos/rancher/terraform-provider-rancher2/commits/${commit.sha}/pulls`], { encoding: 'utf-8', stdio: ['ignore', 'pipe', 'pipe'] }).trim();
    const prs = JSON.parse(prsJson);
    
    // Fix 4: Safely ensure the parsed JSON is verified as an array before filtering
    if (prs && Array.isArray(prs) && prs.length > 0) {
      // Prioritize the backport PR base matching our target branch, filtering out release-please PRs
      const releasePleaseRegex = /^(chore\(release|release-please)/i;
      const nonReleasePleasePrs = prs.filter(p => !releasePleaseRegex.test(p.title || ''));
      
      if (nonReleasePleasePrs.length > 0) {
        pr = nonReleasePleasePrs.find(p => p.base && p.base.ref === branch) || nonReleasePleasePrs[0];
      } else {
        pr = prs.find(p => p.base && p.base.ref === branch) || prs[0];
      }
    }
  } catch (err) {
    const errMsg = err.stderr ? err.stderr.toString().trim() : err.message;
    console.error(`\n⚠️ [Debug] Failed to fetch PR for commit ${commit.sha}: ${errMsg}`);
  }

  let qaIssue = null;
  if (pr) {
    const body = pr.body || '';
    // Capture pattern #1234 or rancher/terraform-provider-rancher2/issues/1234
    const issueMatches = body.match(/#\d+/g) || [];
    const urlMatches = body.match(/issues\/\d+/g) || [];

    const issueNumbers = new Set();
    for (const m of issueMatches) {
      issueNumbers.add(parseInt(m.slice(1), 10));
    }
    for (const m of urlMatches) {
      const parts = m.split('/');
      issueNumbers.add(parseInt(parts[parts.length - 1], 10));
    }

    // Inspect each candidate issue to check labels
    for (const issueNum of issueNumbers) {
      try {
        const issueJson = execFileSync('gh', ['issue', 'view', String(issueNum), '-R', 'rancher/terraform-provider-rancher2', '--json', 'number,title,labels,state,url'], { encoding: 'utf-8', stdio: ['ignore', 'pipe', 'pipe'] }).trim();
        const issue = JSON.parse(issueJson);
        const isBackportQA = issue.labels && issue.labels.some(l => l.name === 'internal/backport');
        if (isBackportQA) {
          qaIssue = issue;
          break; // Found the QA tracking issue, no need to look further
        }
      } catch (err) {
        const errMsg = err.stderr ? err.stderr.toString().trim() : err.message;
        console.error(`\n⚠️ [Debug] Failed to fetch issue #${issueNum}: ${errMsg}`);
      }
    }
  }

  if (pr) {
    if (qaIssue) {
      process.stdout.write(`Matched to Backport PR #${pr.number} and QA Issue #${qaIssue.number} (${qaIssue.state})\n`);
    } else {
      process.stdout.write(`Matched to Backport PR #${pr.number} (⚠️ Missing QA Issue)\n`);
    }
  } else {
    process.stdout.write(`⚠️ No PR found!\n`);
  }

  results.push({
    sha: commit.sha,
    subject: commit.subject,
    pr: pr ? { number: pr.number, title: pr.title, url: pr.html_url, createdAt: pr.created_at } : null,
    qaIssue: qaIssue ? { number: qaIssue.number, title: qaIssue.title, state: qaIssue.state, url: qaIssue.url } : null,
  });
}

// Preserve User's Sorting: Sort results by PR creation date (ascending)
results.sort((a, b) => {
  const dateA = a.pr && a.pr.createdAt ? a.pr.createdAt : '9999-99-99T99:99:99Z';
  const dateB = b.pr && b.pr.createdAt ? b.pr.createdAt : '9999-99-99T99:99:99Z';
  return dateA.localeCompare(dateB);
});

// 4. Generate structured reports
console.log(`\n📝 [4/4] Generating reports...`);

// Render beautiful console table (including User's Created column)
console.log(`\n============================================================================================================================`);
console.log(`| SHA      | Backport PR  | QA Issue   | Status   | Created    | Title                                                     |`);
console.log(`============================================================================================================================`);
for (const r of results) {
  const shortSha = r.sha.substring(0, 8);
  const prCol = r.pr ? `#${r.pr.number}`.padEnd(12) : 'N/A'.padEnd(12);
  const issueCol = r.qaIssue ? `#${r.qaIssue.number}`.padEnd(10) : 'N/A'.padEnd(10);
  
  let statusCol = '❌ MISSING';
  if (r.qaIssue) {
    statusCol = r.qaIssue.state === 'OPEN' ? '🟢 OPEN   ' : '🔵 CLOSED ';
  }
  
  const createdCol = (r.pr && r.pr.createdAt ? r.pr.createdAt.substring(0, 10) : 'N/A').padEnd(10);
  
  const title = (r.pr ? r.pr.title : r.subject).substring(0, 56);
  const titleCol = title.padEnd(57);
  
  console.log(`| ${shortSha} | ${prCol} | ${issueCol} | ${statusCol} | ${createdCol} | ${titleCol} |`);
}
console.log(`============================================================================================================================\n`);

// Generate Markdown report
const reportPath = path.resolve('backport_qa_report.md');
let mdContent = `# 🔍 Release Audit Report: Backport QA Status

**Audit Date:** ${new Date().toLocaleDateString()}
**Target Branch:** \`${branch}\`
**Compare Ref:** \`${baseCompare}\`

---

## Summary Table

| Release Commit | Backport PR | QA Tracking Issue | QA Status | PR Created | Description / Title |
| :--- | :--- | :--- | :---: | :--- | :--- |
`;

for (const r of results) {
  const commitLink = `[\`${r.sha.substring(0, 8)}\`](https://github.com/rancher/terraform-provider-rancher2/commit/${r.sha})`;
  const prLink = r.pr ? `[#${r.pr.number}](${r.pr.url})` : '`N/A`';
  const qaLink = r.qaIssue ? `[#${r.qaIssue.number}](${r.qaIssue.url})` : '`N/A`';
  const createdDate = r.pr && r.pr.createdAt ? r.pr.createdAt.substring(0, 10) : '`N/A`';
  
  let statusBadge = '🔴 **Missing**';
  if (r.qaIssue) {
    statusBadge = r.qaIssue.state === 'OPEN' ? '🟢 **Open**' : '🔵 **Closed**';
  }
  
  // Fix 5: Sanitize PR title/subject for Markdown table safety (escaping pipes and removing newlines)
  const title = (r.pr ? r.pr.title : r.subject)
    .replace(/\|/g, '\\|')
    .replace(/\r?\n|\r/g, ' ');
    
  mdContent += `| ${commitLink} | ${prLink} | ${qaLink} | ${statusBadge} | ${createdDate} | ${title} |\n`;
}

mdContent += `
---
*Report generated automatically by \`scripts/trace_backports.js\`*
`;

try {
  fs.writeFileSync(reportPath, mdContent);
  console.log(`🎉 Markdown report successfully generated at: ${reportPath}`);
} catch (err) {
  console.error(`❌ Failed to write Markdown report:`, err.message);
}
