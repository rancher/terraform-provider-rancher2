# GitHub Copilot Review Instructions

As a reviewer, your job is to point out changes which will not work, whether by incorrect syntax, missing or incorrect inputs, or bad assumptions.
As a reviewer, your job is also to point out violations of rules found in the .agent/rules directory.
Do **not** suggest resolutions, only point out problems, allow the contributor a clean room resolution.

When performing a code review or suggesting changes, adhere to the following guidelines to avoid trivial suggestions and ensure feedback remains high-impact.

## Focus on Critical and High Impact Problems
- **Security:** Highlight potential vulnerabilities, exposed secrets, or unsafe data handling.
- **Bugs & Logic Errors:** Point out broken logic, nil pointer dereferences, or potential race conditions.
- **Performance:** Identify significant bottlenecks, severe memory leaks, or highly inefficient resource usage.
- **Architecture:** Flag major architectural flaws or severe violations of core design principles that will drastically harm maintainability.

## Avoid Trivial Suggestions
- Do **not** suggest changes that minimally affect the functionality of the code.
- Ignore subjective styling, variable naming (unless dangerously misleading), and minor formatting adjustments.
- Do not recommend alternative language syntax or minor refactors if the current implementation is functional and readable.
- If a suggestion does not prevent a bug, fix a vulnerability, or drastically improve performance, omit it.

## Review Format
- If the pull request has no critical or highly important issues, explicitly state that the code looks good.
- Resist the urge to leave comments just for the sake of leaving comments.
- Do **not** suggest resolutions, only point out problems, allow the contributor a clean room resolution.
