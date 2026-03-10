package checker

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestNewChecker(t *testing.T) {
	executor := &mockExecutor{}

	t.Run("default options", func(t *testing.T) {
		c := NewChecker(executor)
		if c == nil {
			t.Fatal("expected checker, got nil")
		}
		if c.parallel != 1 {
			t.Errorf("expected parallel 1, got %d", c.parallel)
		}
		if c.matchMode != MatchExact {
			t.Errorf("expected MatchExact, got %v", c.matchMode)
		}
	})

	t.Run("with parallel option", func(t *testing.T) {
		c := NewChecker(executor, WithParallel(4))
		if c.parallel != 4 {
			t.Errorf("expected parallel 4, got %d", c.parallel)
		}
	})

	t.Run("with match mode option", func(t *testing.T) {
		c := NewChecker(executor, WithMatchMode(MatchContains))
		if c.matchMode != MatchContains {
			t.Errorf("expected MatchContains, got %v", c.matchMode)
		}
	})

	t.Run("parallel zero uses default", func(t *testing.T) {
		c := NewChecker(executor, WithParallel(0))
		if c.parallel != 1 {
			t.Errorf("expected parallel 1, got %d", c.parallel)
		}
	})
}

func TestLoadChecks(t *testing.T) {
	executor := &mockExecutor{}
	c := NewChecker(executor)

	t.Run("non-existent file", func(t *testing.T) {
		err := c.LoadChecks("nonexistent.yaml")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})

	t.Run("valid yaml file", func(t *testing.T) {
		content := `- id: "1.1.1"
  description: "Test check"
  command: "echo test"
  expected: "test"
- id: "1.1.2"
  description: "Another test"
  command: "echo pass"
  expected: "pass"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		err := c.LoadChecks(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checks := c.GetChecks()
		if len(checks) != 2 {
			t.Errorf("expected 2 checks, got %d", len(checks))
		}
	})

	t.Run("duplicate IDs", func(t *testing.T) {
		content := `- id: "1.1.1"
  description: "First"
  command: "echo first"
  expected: "first"
- id: "1.1.1"
  description: "Duplicate"
  command: "echo dup"
  expected: "dup"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		err := c.LoadChecks(tmpFile)
		if err == nil {
			t.Error("expected error for duplicate IDs")
		}
		if err != nil && !containsString(err.Error(), "duplicate") {
			t.Errorf("expected duplicate error, got: %v", err)
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		content := `- description: "No ID"
  command: "echo test"
  expected: "test"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		err := c.LoadChecks(tmpFile)
		if err == nil {
			t.Error("expected error for missing ID")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		content := `- id: "1.1.1"
  description: "Test"
  command: "echo test
  expected: "test"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		err := c.LoadChecks(tmpFile)
		if err == nil {
			t.Error("expected error for invalid yaml")
		}
	})
}

func TestRunChecks(t *testing.T) {
	t.Run("sequential execution", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{
				"echo test": "test",
				"echo pass": "pass",
				"false":     "",
			},
		}
		c := NewChecker(executor)

		content := `- id: "1.1.1"
  description: "Test check"
  command: "echo test"
  expected: "test"
- id: "1.1.2"
  description: "Another test"
  command: "echo pass"
  expected: "pass"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), nil)
		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}

		for _, r := range results {
			if r.Result != ResultPass {
				t.Errorf("expected PASS for %s, got %s", r.ID, r.Result)
			}
		}
	})

	t.Run("parallel execution", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{
				"check1": "ok",
				"check2": "ok",
				"check3": "ok",
			},
		}
		c := NewChecker(executor, WithParallel(3))

		content := `- id: "1"
  description: "Check 1"
  command: "check1"
  expected: "ok"
- id: "2"
  description: "Check 2"
  command: "check2"
  expected: "ok"
- id: "3"
  description: "Check 3"
  command: "check3"
  expected: "ok"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), nil)
		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}
	})

	t.Run("specific check IDs", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{
				"echo test": "test",
			},
		}
		c := NewChecker(executor)

		content := `- id: "1.1.1"
  description: "First"
  command: "echo first"
  expected: "first"
- id: "1.1.2"
  description: "Second"
  command: "echo test"
  expected: "test"
- id: "1.1.3"
  description: "Third"
  command: "echo third"
  expected: "third"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), []string{"1.1.2"})
		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if results[0].ID != "1.1.2" {
			t.Errorf("expected ID 1.1.2, got %s", results[0].ID)
		}
	})

	t.Run("duplicate check IDs in request", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{
				"echo test": "test",
			},
		}
		c := NewChecker(executor)

		content := `- id: "1.1.1"
  description: "Test"
  command: "echo test"
  expected: "test"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), []string{"1.1.1", "1.1.1", "1.1.1"})
		if len(results) != 1 {
			t.Errorf("expected 1 result (deduplicated), got %d", len(results))
		}
	})

	t.Run("empty checks", func(t *testing.T) {
		c := NewChecker(&mockExecutor{})
		results := c.RunChecks(context.Background(), nil)
		if results != nil {
			t.Errorf("expected nil for empty checks, got %v", results)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		executor := &mockExecutor{shouldTimeout: true}
		c := NewChecker(executor)

		content := `- id: "1"
  description: "Timeout test"
  command: "sleep 10"
  expected: ""
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		results := c.RunChecks(ctx, nil)
		for _, r := range results {
			if !containsString(r.Result, ResultError) {
				t.Errorf("expected ERROR result, got %s", r.Result)
			}
		}
	})
}

func TestMatchModes(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{
				"cmd": "exact output value",
			},
		}
		c := NewChecker(executor, WithMatchMode(MatchExact))

		content := `- id: "1"
  description: "Test"
  command: "cmd"
  expected: "exact output value"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), nil)
		if results[0].Result != ResultPass {
			t.Errorf("expected PASS for exact match, got %s", results[0].Result)
		}
	})

	t.Run("contains match", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{
				"cmd": "this contains the expected value somewhere",
			},
		}
		c := NewChecker(executor, WithMatchMode(MatchContains))

		content := `- id: "1"
  description: "Test"
  command: "cmd"
  expected: "expected value"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), nil)
		if results[0].Result != ResultPass {
			t.Errorf("expected PASS for contains match, got %s", results[0].Result)
		}
	})

	t.Run("exact match failure", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{
				"cmd": "output value",
			},
		}
		c := NewChecker(executor, WithMatchMode(MatchExact))

		content := `- id: "1"
  description: "Test"
  command: "cmd"
  expected: "different value"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), nil)
		if results[0].Result != ResultFail {
			t.Errorf("expected FAIL for mismatch, got %s", results[0].Result)
		}
	})
}

func TestExitCodeHandling(t *testing.T) {
	t.Run("exit code non-zero with expected empty", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{},
			exitCodes: map[string]int{
				"grep pattern": 1,
			},
		}
		c := NewChecker(executor)

		content := `- id: "1"
  description: "Test"
  command: "grep pattern"
  expected: ""
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), nil)
		if results[0].Result != ResultPass {
			t.Errorf("expected PASS for exit 1 with empty expected, got %s", results[0].Result)
		}
	})

	t.Run("exit code non-zero with expected value", func(t *testing.T) {
		executor := &mockExecutor{
			responses: map[string]string{},
			exitCodes: map[string]int{
				"grep pattern": 1,
			},
		}
		c := NewChecker(executor)

		content := `- id: "1"
  description: "Test"
  command: "grep pattern"
  expected: "found"
`
		tmpFile := createTempFile(t, content)
		defer osRemove(t, tmpFile)

		if err := c.LoadChecks(tmpFile); err != nil {
			t.Fatalf("failed to load checks: %v", err)
		}

		results := c.RunChecks(context.Background(), nil)
		if results[0].Result != ResultFail {
			t.Errorf("expected FAIL for exit 1 with expected value, got %s", results[0].Result)
		}
	})
}

func TestRemediation(t *testing.T) {
	executor := &mockExecutor{
		responses: map[string]string{
			"cmd": "output",
		},
	}
	c := NewChecker(executor)

	content := `- id: "1"
  description: "Test"
  command: "cmd"
  expected: "output"
  remediation: "Run this to fix: apt install package"
`
	tmpFile := createTempFile(t, content)
	defer osRemove(t, tmpFile)

	if err := c.LoadChecks(tmpFile); err != nil {
		t.Fatalf("failed to load checks: %v", err)
	}

	checks := c.GetChecks()
	if len(checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(checks))
	}

	if checks[0].Remediation != "Run this to fix: apt install package" {
		t.Errorf("expected remediation, got %s", checks[0].Remediation)
	}
}

type mockExecutor struct {
	responses     map[string]string
	exitCodes     map[string]int
	shouldTimeout bool
}

func (m *mockExecutor) Execute(ctx context.Context, cmd string) (string, error) {
	if m.shouldTimeout {
		return "", &TimeoutError{Duration: 30}
	}

	if code, ok := m.exitCodes[cmd]; ok && code != 0 {
		return m.responses[cmd], &ExitError{Code: code}
	}

	if resp, ok := m.responses[cmd]; ok {
		return resp, nil
	}

	return "", errors.New("unexpected command")
}

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "cis_test_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}
	return f.Name()
}

func osRemove(t *testing.T, path string) {
	t.Helper()
	os.Remove(path)
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
