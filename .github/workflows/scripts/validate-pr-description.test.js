import assert from 'node:assert/strict';
import test from 'node:test';
import run, { validateIssueExists, validatePrDescription } from './validate-pr-description.js';

test('validatePrDescription - succeeds with valid format', () => {
  const prBody = 'Some initial context\n- Addresses: #1234\nSome other details';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '1234'},
  ]);
});

test('validatePrDescription - succeeds with asterisk bullet', () => {
  const prBody = 'Some initial context\n* Addresses: #1234\nSome other details';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '1234'},
  ]);
});

test('validatePrDescription - succeeds with whitespace around line', () => {
  const prBody = '  - Addresses: #999  ';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '999'},
  ]);
});

test('validatePrDescription - succeeds with case-insensitive and flexible spacing', () => {
  const prBody = '- addresses:  #4321';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '4321'},
  ]);
});

test('validatePrDescription - succeeds with multiple issue numbers', () => {
  const prBody = 'Context\n- Addresses: #1234\n- Addresses: #5678\nFooter';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '1234'},
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '5678'},
  ]);
});

test('validatePrDescription - succeeds with trailing period', () => {
  const prBody = '- Addresses: #1234.';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '1234'},
  ]);
});

test('validatePrDescription - succeeds with inline code backticks around issue', () => {
  const prBody = '- Addresses: `#1234`';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '1234'},
  ]);
});

test('validatePrDescription - succeeds with owner and repo', () => {
  const prBody = '- Addresses: fooo/baar#1234';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'fooo', repo: 'baar', issue: '1234'},
  ]);
});

test('validatePrDescription - succeeds multiple different issues', () => {
  const prBody = '- Addresses: rancher/rancher#1234\n- Addresses: harvester/harvester#123';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'rancher', issue: '1234'},
    {owner: 'harvester', repo: 'harvester', issue: '123'},
  ]);
});

test('validatePrDescription - succeeds with CRLF line endings', () => {
  const prBody = 'Some initial context\r\n- Addresses: #1234\r\nSome other details';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '1234'},
  ]);
});

test('validatePrDescription - ignores issue reference inside HTML comments', () => {
  const prBody = '<!--\n- Addresses: #1234\n-->\nSome other details';
  assert.throws(
    () => validatePrDescription('rancher', 'terraform-provider-rancher2', prBody),
    {
      message: 'Please add a description line that matches this format: - Addresses: #<issue number>'
    }
  );
});

test('validatePrDescription - matches valid issue outside HTML comments', () => {
  const prBody = '<!--\n- Addresses: #1234\n-->\n- Addresses: #5678\nSome other details';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '5678'},
  ]);
});

test('validatePrDescription - ignores issue reference inside fenced code blocks', () => {
  const prBody = '```markdown\n- Addresses: #1234\n```\nSome other details';
  assert.throws(
    () => validatePrDescription('rancher', 'terraform-provider-rancher2', prBody),
    {
      message: 'Please add a description line that matches this format: - Addresses: #<issue number>'
    }
  );
});

test('validatePrDescription - matches valid issue outside fenced code blocks', () => {
  const prBody = '```markdown\n- Addresses: #1234\n```\n- Addresses: #5678\nSome other details';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '5678'},
  ]);
});

test('validatePrDescription - ignores issue reference inside tilde fenced code blocks', () => {
  const prBody = '~~~markdown\n- Addresses: #1234\n~~~\nSome other details';
  assert.throws(
    () => validatePrDescription('rancher', 'terraform-provider-rancher2', prBody),
    {
      message: 'Please add a description line that matches this format: - Addresses: #<issue number>'
    }
  );
});

test('validatePrDescription - matches valid issue outside tilde fenced code blocks', () => {
  const prBody = '~~~markdown\n- Addresses: #1234\n~~~\n- Addresses: #5678\nSome other details';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '5678'},
  ]);
});

test('validatePrDescription - ignores issue reference inside 4-backtick fenced code blocks', () => {
  const prBody = '````markdown\n- Addresses: #1234\n````\nSome other details';
  assert.throws(
    () => validatePrDescription('rancher', 'terraform-provider-rancher2', prBody),
    {
      message: 'Please add a description line that matches this format: - Addresses: #<issue number>'
    }
  );
});

test('validatePrDescription - matches valid issue outside 4-backtick fenced code blocks', () => {
  const prBody = '````markdown\n- Addresses: #1234\n````\n- Addresses: #5678\nSome other details';
  const issueNumbers = validatePrDescription('rancher', 'terraform-provider-rancher2', prBody);
  assert.deepStrictEqual(issueNumbers, [
    {owner: 'rancher', repo: 'terraform-provider-rancher2', issue: '5678'},
  ]);
});

test('validatePrDescription - throws on empty description', () => {
  assert.throws(
    () => validatePrDescription('rancher', 'terraform-provider-rancher2', ''),
    {
      message: 'Please add a description line that matches this format: - Addresses: #<issue number>'
    }
  );
  assert.throws(
    () => validatePrDescription('rancher', 'terraform-provider-rancher2', null),
    {
      message: 'Please add a description line that matches this format: - Addresses: #<issue number>'
    }
  );
});

test('validatePrDescription - throws when format does not match', () => {
  const invalidBodies = [
    'Addresses: #1234',
    '- Addresses #1234',
    '- Addresses: 1234',
    '- Addresses: #1234 extra text',
    'Some text - Addresses: #1234',
    '- Addresses: `#1234',
    '- Addresses: #1234`',
  ];

  for (const body of invalidBodies) {
    assert.throws(
      () => validatePrDescription('rancher', 'terraform-provider-rancher2', body),
      {
        message: 'Please add a description line that matches this format: - Addresses: #<issue number>'
      }
    );
  }
});

test('validateIssueExists - succeeds when issue is found', async () => {
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ owner, repo, issue_number }) => {
          assert.strictEqual(owner, 'rancher');
          assert.strictEqual(repo, 'terraform-provider-rancher2');
          assert.strictEqual(issue_number, 1234);
          return { data: { id: 1, number: 1234, state: 'open' } };
        }
      }
    }
  };

  const exists = await validateIssueExists(mockGithub, 'rancher', 'terraform-provider-rancher2', '1234');
  assert.strictEqual(exists, true);
});

test('validateIssueExists - throws when issue is not found (404)', async () => {
  const mockGithub = {
    rest: {
      issues: {
        get: async () => {
          const error = new Error('Not Found');
          error.status = 404;
          throw error;
        }
      }
    }
  };

  await assert.rejects(
    async () => validateIssueExists(mockGithub, 'rancher', 'terraform-provider-rancher2', '9999'),
    {
      message: 'Issue #9999 not found in rancher/terraform-provider-rancher2 repo.'
    }
  );
});

test('validateIssueExists - throws when referenced item is a pull request', async () => {
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ owner, repo, issue_number }) => {
          assert.strictEqual(owner, 'rancher');
          assert.strictEqual(repo, 'terraform-provider-rancher2');
          assert.strictEqual(issue_number, 5678);
          return { data: { id: 2, number: 5678, pull_request: {} } };
        }
      }
    }
  };

  await assert.rejects(
    async () => validateIssueExists(mockGithub, 'rancher', 'terraform-provider-rancher2', '5678'),
    {
      message: 'Issue #5678 is a pull request, not an issue.'
    }
  );
});

test('validateIssueExists - throws when referenced item is closed', async () => {
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ owner, repo, issue_number }) => {
          assert.strictEqual(owner, 'rancher');
          assert.strictEqual(repo, 'terraform-provider-rancher2');
          assert.strictEqual(issue_number, 4321);
          return { data: { id: 3, number: 4321, state: 'closed' } };
        }
      }
    }
  };

  await assert.rejects(
    async () => validateIssueExists(mockGithub, 'rancher', 'terraform-provider-rancher2', '4321'),
    {
      message: 'Issue #4321 is already closed.'
    }
  );
});

test('validateIssueExists - rethrows unexpected API errors', async () => {
  const mockGithub = {
    rest: {
      issues: {
        get: async () => {
          const error = new Error('Internal Server Error');
          error.status = 500;
          throw error;
        }
      }
    }
  };

  await assert.rejects(
    async () => validateIssueExists(mockGithub, 'rancher', 'terraform-provider-rancher2', '1234'),
    {
      message: 'Internal Server Error'
    }
  );
});

test('run - succeeds when PR body contains valid open issue', async () => {
  let failedMessage = null;
  const mockCore = {
    setFailed: (msg) => { failedMessage = msg; }
  };
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ owner, repo, issue_number }) => {
          assert.strictEqual(owner, 'rancher');
          assert.strictEqual(repo, 'rancher');
          assert.strictEqual(issue_number, 1234);
          return { data: { id: 1, number: 1234, state: 'open' } };
        }
      }
    }
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'rancher' },
    payload: {
      pull_request: {
        body: 'Context\n- Addresses: #1234\nNotes'
      }
    }
  };

  await run({ github: mockGithub, context: mockContext, core: mockCore });
  assert.strictEqual(failedMessage, null);
});

test('run - fails when PR body is missing valid format', async () => {
  let failedMessage = null;
  const mockCore = {
    setFailed: (msg) => { failedMessage = msg; }
  };
  const mockGithub = { rest: { issues: {} } };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        body: 'No addresses line here'
      }
    }
  };

  await run({ github: mockGithub, context: mockContext, core: mockCore });
  assert.strictEqual(
    failedMessage,
    'Please add a description line that matches this format: - Addresses: #<issue number>'
  );
});

test('run - fails when referenced issue does not exist', async () => {
  let failedMessage = null;
  const mockCore = {
    setFailed: (msg) => { failedMessage = msg; }
  };
  const mockGithub = {
    rest: {
      issues: {
        get: async () => {
          const err = new Error('Not Found');
          err.status = 404;
          throw err;
        }
      }
    }
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        body: '- Addresses: #9999'
      }
    }
  };

  await run({ github: mockGithub, context: mockContext, core: mockCore });
  assert.strictEqual(
    failedMessage,
    'Issue #9999 not found in rancher/terraform-provider-rancher2 repo.'
  );
});

test('run - deduplicates multiple occurrences of same issue and succeeds', async () => {
  let failedMessage = null;
  let callCount = 0;
  const mockCore = {
    setFailed: (msg) => { failedMessage = msg; }
  };
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          callCount++;
          assert.strictEqual(issue_number, 1234);
          return { data: { id: 1, number: 1234, state: 'open' } };
        }
      }
    }
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        body: '- Addresses: #1234\n* Addresses: #1234'
      }
    }
  };

  await run({ github: mockGithub, context: mockContext, core: mockCore });
  assert.strictEqual(failedMessage, null);
  assert.strictEqual(callCount, 1);
});

test('run - aggregates errors when multiple referenced issues fail validation', async () => {
  let failedMessage = null;
  const mockCore = {
    setFailed: (msg) => { failedMessage = msg; }
  };
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 101) {
            const err = new Error('Not Found');
            err.status = 404;
            throw err;
          }
          if (issue_number === 102) {
            return { data: { id: 2, number: 102, state: 'closed' } };
          }
          throw new Error('Unexpected issue');
        }
      }
    }
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        body: '- Addresses: #101\n- Addresses: #102'
      }
    }
  };

  await run({ github: mockGithub, context: mockContext, core: mockCore });
  assert.strictEqual(
    failedMessage,
    'Issue #101 not found in rancher/terraform-provider-rancher2 repo.\nIssue #102 is already closed.'
  );
});

test('run - fails when referenced item is a pull request', async () => {
  let failedMessage = null;
  const mockCore = {
    setFailed: (msg) => { failedMessage = msg; }
  };
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 5678);
          return { data: { id: 2, number: 5678, pull_request: {} } };
        }
      }
    }
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        body: '- Addresses: #5678'
      }
    }
  };

  await run({ github: mockGithub, context: mockContext, core: mockCore });
  assert.strictEqual(
    failedMessage,
    'Issue #5678 is a pull request, not an issue.'
  );
});

test('run - fails when referenced issue is already closed', async () => {
  let failedMessage = null;
  const mockCore = {
    setFailed: (msg) => { failedMessage = msg; }
  };
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 4321);
          return { data: { id: 3, number: 4321, state: 'closed' } };
        }
      }
    }
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        body: '- Addresses: #4321'
      }
    }
  };

  await run({ github: mockGithub, context: mockContext, core: mockCore });
  assert.strictEqual(
    failedMessage,
    'Issue #4321 is already closed.'
  );
});

test('run - fails and reports message on unexpected API error', async () => {
  let failedMessage = null;
  const mockCore = {
    setFailed: (msg) => { failedMessage = msg; }
  };
  const mockGithub = {
    rest: {
      issues: {
        get: async () => {
          const err = new Error('Internal Server Error');
          err.status = 500;
          throw err;
        }
      }
    }
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        body: '- Addresses: #1234'
      }
    }
  };

  await run({ github: mockGithub, context: mockContext, core: mockCore });
  assert.strictEqual(failedMessage, 'Internal Server Error');
});
