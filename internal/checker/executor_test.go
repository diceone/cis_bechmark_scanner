package checker

import (
	"context"
	"testing"
	"time"
)

func TestNewCommandExecutor(t *testing.T) {
	t.Run("default timeout", func(t *testing.T) {
		e := NewCommandExecutor(0)
		if e.defaultTimeout != 30*time.Second {
			t.Errorf("expected default timeout 30s, got %v", e.defaultTimeout)
		}
	})

	t.Run("custom timeout", func(t *testing.T) {
		e := NewCommandExecutor(60)
		if e.defaultTimeout != 60*time.Second {
			t.Errorf("expected timeout 60s, got %v", e.defaultTimeout)
		}
	})
}

func TestExecute(t *testing.T) {
	t.Run("successful command", func(t *testing.T) {
		e := NewCommandExecutor(10)
		output, err := e.Execute(context.Background(), "echo hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output != "hello" {
			t.Errorf("expected 'hello', got %q", output)
		}
	})

	t.Run("command with arguments", func(t *testing.T) {
		e := NewCommandExecutor(10)
		output, err := e.Execute(context.Background(), "echo hello world")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output != "hello world" {
			t.Errorf("expected 'hello world', got %q", output)
		}
	})

	t.Run("command exit code non-zero", func(t *testing.T) {
		e := NewCommandExecutor(10)
		output, err := e.Execute(context.Background(), "exit 1")
		if err == nil {
			t.Fatal("expected error for exit code 1")
		}

		exitErr, ok := err.(*ExitError)
		if !ok {
			t.Errorf("expected ExitError, got %T", err)
		}
		if exitErr.Code != 1 {
			t.Errorf("expected exit code 1, got %d", exitErr.Code)
		}
		if output != "" {
			t.Errorf("expected empty output, got %q", output)
		}
	})

	t.Run("command not found", func(t *testing.T) {
		e := NewCommandExecutor(10)
		_, err := e.Execute(context.Background(), "nonexistentcommand12345")
		if err == nil {
			t.Error("expected error for nonexistent command")
		}
	})

	t.Run("command with stderr output", func(t *testing.T) {
		e := NewCommandExecutor(10)
		_, err := e.Execute(context.Background(), "echo error >&2 && exit 1")
		if err == nil {
			t.Fatal("expected error")
		}

		exitErr, ok := err.(*ExitError)
		if !ok {
			t.Errorf("expected ExitError, got %T", err)
		}
		if exitErr.Code != 1 {
			t.Errorf("expected exit code 1, got %d", exitErr.Code)
		}
		if exitErr.Message != "error" {
			t.Errorf("expected stderr 'error', got %q", exitErr.Message)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		e := NewCommandExecutor(10)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := e.Execute(ctx, "echo test")
		if err == nil {
			t.Error("expected error for cancelled context")
		}
	})

	t.Run("command timeout", func(t *testing.T) {
		e := NewCommandExecutor(1)
		ctx := context.Background()

		_, err := e.Execute(ctx, "sleep 10")
		if err == nil {
			t.Error("expected timeout error")
		}

		_, ok := err.(*TimeoutError)
		if !ok {
			t.Errorf("expected TimeoutError, got %T", err)
		}
	})
}

func TestExitError(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		err := &ExitError{Code: 1, Message: "error message"}
		if err.Error() != "exit code 1: error message" {
			t.Errorf("unexpected error string: %s", err.Error())
		}
		if err.ExitStatus() != 1 {
			t.Errorf("unexpected exit status: %d", err.ExitStatus())
		}
	})

	t.Run("without message", func(t *testing.T) {
		err := &ExitError{Code: 2}
		if err.Error() != "exit code 2" {
			t.Errorf("unexpected error string: %s", err.Error())
		}
	})
}

func TestTimeoutError(t *testing.T) {
	err := &TimeoutError{Duration: 30 * time.Second}
	expected := "command timed out after 30s"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
