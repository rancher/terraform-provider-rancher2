package rancher2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// tokenCommandTimeout bounds token_command. The command is expected to return
// a cached token quickly; an interactive login can't be completed from inside
// Terraform, so a hanging command is reported instead of waited on.
const tokenCommandTimeout = 60 * time.Second

// tokenCommandWaitDelay is how long to wait for the output pipes to close
// after the command is killed. A descendant that inherited them (e.g. from
// `sh -c`) would otherwise keep the provider waiting past the timeout.
const tokenCommandWaitDelay = 2 * time.Second

// maxTokenCommandStderr caps how much stderr is included in an error.
const maxTokenCommandStderr = 512

// tokenLikeRe matches Rancher tokens and other long secret-looking strings
// (bearer tokens, JWT segments). Dots are excluded so hostnames in errors
// such as "no such host" stay readable. This is best effort: short secrets
// or ones split by other characters are not caught.
var tokenLikeRe = regexp.MustCompile(`token-[a-z0-9]+:[A-Za-z0-9]+|[A-Za-z0-9_\-+=]{32,}`)

// sanitizeTokenCommandStderr redacts secret-looking strings from stderr on a
// best-effort basis and truncates it.
func sanitizeTokenCommandStderr(stderr []byte) string {
	s := tokenLikeRe.ReplaceAllString(strings.TrimSpace(string(stderr)), "[REDACTED]")
	if len(s) > maxTokenCommandStderr {
		s = s[:maxTokenCommandStderr] + "...(truncated)"
	}
	return s
}

// runTokenCommand runs command and returns the Rancher API token it prints on
// stdout. The output can be a raw token or a Kubernetes ExecCredential JSON
// (as printed by `rancher token`), in which case status.token is used.
// The token from stdout is never included in errors; stderr is included with
// secret-looking strings redacted on a best-effort basis.
func runTokenCommand(command []string, timeout time.Duration) (string, error) {
	if len(command) == 0 || command[0] == "" {
		return "", errors.New("[ERROR] token_command must contain at least the program to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.WaitDelay = tokenCommandWaitDelay

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("[ERROR] token_command %q timed out after %s. If it needs an interactive (browser) login, run it once in a terminal first so it can use a cached token", command[0], timeout)
	}
	if err != nil {
		return "", fmt.Errorf("[ERROR] token_command %q failed: %v: %s", command[0], err, sanitizeTokenCommandStderr(stderr.Bytes()))
	}

	token, err := parseTokenCommandOutput(stdout.Bytes())
	if err != nil {
		return "", fmt.Errorf("[ERROR] token_command %q: %v", command[0], err)
	}
	return token, nil
}

func parseTokenCommandOutput(out []byte) (string, error) {
	trimmed := bytes.TrimSpace(out)
	if len(trimmed) == 0 {
		return "", errors.New("printed no token")
	}

	if trimmed[0] == '{' {
		var cred struct {
			Kind   string `json:"kind"`
			Status struct {
				Token string `json:"token"`
			} `json:"status"`
		}
		if err := json.Unmarshal(trimmed, &cred); err != nil {
			return "", fmt.Errorf("output looks like JSON but can't be parsed: %v", err)
		}
		if cred.Status.Token == "" {
			return "", errors.New("JSON output has no status.token")
		}
		return cred.Status.Token, nil
	}

	if strings.ContainsAny(string(trimmed), " \t\r\n") {
		return "", errors.New("output is neither a single token nor an ExecCredential JSON")
	}
	return string(trimmed), nil
}
