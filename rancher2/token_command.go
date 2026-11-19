package rancher2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
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

// maxTokenCommandOutput caps stdout. A token or ExecCredential JSON is a few
// KiB at most; more means the command is misconfigured.
const maxTokenCommandOutput = 64 * 1024

// errTokenCommandOutputTooLarge stops the copy from the command's stdout.
var errTokenCommandOutputTooLarge = errors.New("token_command output too large")

// cappedBuffer collects up to max bytes. When more is written it kills the
// command (via cancel) instead of growing without bound.
type cappedBuffer struct {
	buf      bytes.Buffer
	max      int
	cancel   context.CancelFunc
	exceeded bool
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	if b.buf.Len()+len(p) > b.max {
		b.exceeded = true
		b.cancel()
		return 0, errTokenCommandOutputTooLarge
	}
	return b.buf.Write(p)
}

func (b *cappedBuffer) Bytes() []byte { return b.buf.Bytes() }

// runTokenCommand runs command and returns the Rancher API token it prints on
// stdout. The output can be a raw token or a Kubernetes ExecCredential JSON
// (as printed by `rancher token`), in which case status.token is used.
// Neither stdout nor stderr is included in errors or logs, because either
// could contain the token; on failure the user is told to run the command in
// a terminal to see its output.
func runTokenCommand(command []string, timeout time.Duration) (string, error) {
	if len(command) == 0 || command[0] == "" {
		return "", errors.New("[ERROR] token_command must contain at least the program to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stdout := &cappedBuffer{max: maxTokenCommandOutput, cancel: cancel}
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Stdout = stdout
	// stderr is discarded (nil): it may contain the token.
	cmd.WaitDelay = tokenCommandWaitDelay

	err := cmd.Run()
	if stdout.exceeded {
		return "", fmt.Errorf("[ERROR] token_command %q printed more than %d bytes; it should print only a token or an ExecCredential JSON", command[0], maxTokenCommandOutput)
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("[ERROR] token_command %q timed out after %s. If it needs an interactive (browser) login, run it once in a terminal first so it can use a cached token", command[0], timeout)
	}
	if err != nil {
		return "", fmt.Errorf("[ERROR] token_command %q failed: %v. Its output is not shown because it may contain the token; run the command in a terminal to see why it failed", command[0], err)
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
