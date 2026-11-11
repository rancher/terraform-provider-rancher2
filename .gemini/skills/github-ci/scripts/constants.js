export const DEFAULT_PATTERN =
  /\b(error|fail|failed|failure|exception|fatal|panic|exit (?:code|status) [1-9]\d*|undefined|cannot find|timed out)\b|##\[error\]/i;

/**
 * Robustly parses a workflow run ID string or number, ensuring it is a safe integer and within positive bounds.
 * @param {string|number} val
 * @returns {number|null}
 */
export function parseRunId(val) {
  if (typeof val !== "string" && typeof val !== "number") return null;
  const cleanVal = String(val).replace(/^#/, "");
  if (!/^\d+$/.test(cleanVal)) return null;

  const parsed = parseInt(cleanVal, 10);
  if (
    !Number.isNaN(parsed) &&
    Number.isSafeInteger(parsed) &&
    parsed > 0 &&
    parsed < 1e15
  ) {
    return parsed;
  }
  return null;
}

