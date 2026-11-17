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

// runTokenCommand runs command and returns the Rancher API token it prints on
// stdout. The output can be a raw token or a Kubernetes ExecCredential JSON
// (as printed by `rancher token`), in which case status.token is used.
// The token is never included in errors.
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

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("[ERROR] token_command %q timed out after %s. If it needs an interactive (browser) login, run it once in a terminal first so it can use a cached token", command[0], timeout)
	}
	if err != nil {
		return "", fmt.Errorf("[ERROR] token_command %q failed: %v: %s", command[0], err, strings.TrimSpace(stderr.String()))
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
