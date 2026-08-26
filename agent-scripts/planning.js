import { execSync } from 'child_process';
import fs from 'fs';

/**
 * Checks git status to determine if there is an active (modified, added, untracked) blueprint in docs/development/
 * @param {string} cwd - The current working directory
 * @returns {boolean} - True if an active blueprint exists, false otherwise
 */
export function checkActiveBlueprint(cwd) {
  try {
    const statusOutput = execSync('git status --porcelain', {
      cwd: cwd || process.cwd(),
      stdio: ['ignore', 'pipe', 'ignore'],
    }).toString();

    return statusOutput.split('\n').some((line) => {
      const trimmed = line.trim();
      if (!trimmed.includes('docs/development/')) {
        return false;
      }
      const status = line.substring(0, 2);
      // Ensure the file is not deleted ('D') or ignored ('!')
      return !status.includes('D') && !status.includes('!');
    });
  } catch (err) {
    console.error('Failed to run git status inside blueprint check:', err.message || err);
    return false;
  }
}

/**
 * Programmatically validates that a plan file contains all of our strict requirements.
 * @param {string} planPath - The path to the active plan markdown file
 * @returns {object} - { valid: boolean, errors: string[] }
 */
export function validatePlanContent(planPath) {
  const errors = [];
  if (!fs.existsSync(planPath)) {
    return { valid: false, errors: ['Plan file does not exist.'] };
  }

  const content = fs.readFileSync(planPath, 'utf-8');

  // 1. Checklist check: must contain markdown checklist items "[ ]"
  const checklistMatch = /-\s*\[\s*\]/g.test(content);
  if (!checklistMatch) {
    errors.push('The plan must include each step in a checklist (using "- [ ]").');
  }

  // 2. Comprehensive tests check
  const testMatch = /test|testing|linter/i.test(content);
  if (!testMatch) {
    errors.push('The plan must include running comprehensive tests.');
  }

  // 3. Quality gates check
  const gateMatch = /gate|signature|seal|approval/i.test(content);
  if (!gateMatch) {
    errors.push('The plan must include our standard quality gates.');
  }

  // 4. Maintaining the agentic framework check
  const frameworkMatch = /agentic framework|system script|enforcer hook/i.test(content);
  if (!frameworkMatch) {
    errors.push('The plan must include maintaining the agentic framework if improvements or bugs are found in it.');
  }

  // 5. Updating documentation check
  const docMatch = /document|documentation|docs\//i.test(content);
  if (!docMatch) {
    errors.push('The plan must include updating documentation to describe the changes we plan to make.');
  }

  return {
    valid: errors.length === 0,
    errors,
  };
}

export { checkActiveBlueprint as checkActivePlan };
