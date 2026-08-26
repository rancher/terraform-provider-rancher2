# GitHub Copilot Review Guidelines

---

## Abstract

These guidelines prescribe the evaluation criteria, severity thresholds, and review formats executed by our automated subagents and Copilot assistants during our proactive pre-commit quality checks.

---

## Single-Pass Comprehensive Review

- **No Drip-Feeding:** You MUST provide ALL of your feedback in a single, comprehensive review. Do not hold back comments for future interactions.
- Analyze the entire pull request and consolidate every single finding into your initial response.

## Severity and Confidence

- **Severity Threshold:** ONLY generate comments, suggestions, or feedback for issues that you evaluate as **High** impact or critical security/safety risks.
- **Ignore Low/Medium:** Do not leave comments for Low or Medium severity issues (e.g., general style preferences, minor refactoring, or non-critical formatting).
- **Confidence:** Only provide feedback when you are highly confident ($\ge 80\%$) that an issue or bug exists.

## Focus Exclusively on Functionality-Blocking Issues

- **Security:** Highlight potential vulnerabilities, exposed secrets, or unsafe data handling.
- **Bugs & Logic Errors:** Point out broken logic, unhandled edge cases, nil pointer dereferences, or potential race conditions.
- **Performance:** Identify significant bottlenecks, severe memory leaks, or highly inefficient resource usage.
- **Architecture:** Flag major architectural flaws or severe violations of core design principles that will drastically harm maintainability.

## Review Format

- Provide actionable, concrete feedback ONLY for the critical, blocking issues identified.
- If the pull request has no critical or highly important issues, explicitly state that the code looks good and approve the review.
- Resist the urge to leave comments just for the sake of leaving comments. If in doubt about whether an issue blocks functionality, DO NOT comment.
