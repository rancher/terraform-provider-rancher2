# GitHub Copilot Review Instructions

When performing a code review or suggesting changes, adhere to the following guidelines to avoid "bikeshedding" and ensure feedback remains high-impact:

## Focus on Critical and Highly Important Issues
- **Security:** Highlight potential vulnerabilities, exposed secrets, or unsafe data handling.
- **Bugs & Logic Errors:** Point out broken logic, unhandled edge cases, nil pointer dereferences, or potential race conditions.
- **Performance:** Identify significant bottlenecks, severe memory leaks, or highly inefficient resource usage.
- **Architecture:** Flag major architectural flaws or severe violations of core design principles that will drastically harm maintainability.

## Avoid Bikeshedding (Trivial Suggestions)
- Do **not** suggest changes that minimally affect the functionality of the code.
- Ignore subjective styling, variable naming (unless dangerously misleading), and minor formatting adjustments.
- Do not recommend alternative language syntax or minor refactors if the current implementation is functional and readable.
- If a suggestion does not prevent a bug, fix a vulnerability, or drastically improve performance, omit it.

## Review Format
- Provide actionable, concrete feedback for the critical issues identified.
- If the pull request has no critical or highly important issues, explicitly state that the code looks good and approve the review. 
- Resist the urge to leave comments just for the sake of leaving comments.
