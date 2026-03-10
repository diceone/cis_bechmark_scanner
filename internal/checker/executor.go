package checker

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type CommandExecutor struct {
	defaultTimeout time.Duration
	shell          string
}

func NewCommandExecutor(timeoutSeconds int) *CommandExecutor {
	timeout := time.Duration(timeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &CommandExecutor{
		defaultTimeout: timeout,
		shell:          "/bin/sh",
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, e.defaultTimeout)
	defer cancel()

	execCmd := exec.CommandContext(ctx, e.shell, "-c", cmd)

	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	err := execCmd.Run()
	output := strings.TrimSpace(stdout.String())

	if ctx.Err() == context.DeadlineExceeded {
		return output, &TimeoutError{Duration: e.defaultTimeout}
	}

	if err != nil {
		if execCmd.ProcessState != nil && execCmd.ProcessState.ExitCode() != 0 {
			return output, &ExitError{
				Code:    execCmd.ProcessState.ExitCode(),
				Message: strings.TrimSpace(stderr.String()),
			}
		}
		return output, fmt.Errorf("command execution failed: %w", err)
	}

	return output, nil
}

type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("exit code %d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("exit code %d", e.Code)
}

func (e *ExitError) ExitStatus() int {
	return e.Code
}

type TimeoutError struct {
	Duration time.Duration
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("command timed out after %v", e.Duration)
}
