package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/diceone/cis_benchmark_scanner/internal/checker"
)

func TestGenerateHTML(t *testing.T) {
	t.Run("basic report", func(t *testing.T) {
		results := []checker.Result{
			{ID: "1.1.1", Description: "Test check", Result: checker.ResultPass, Output: "test output"},
			{ID: "1.1.2", Description: "Failed check", Result: checker.ResultFail, Expected: "expected", Output: "actual"},
			{ID: "1.1.3", Description: "Error check", Result: "ERROR: something went wrong"},
		}

		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "report.html")

		err := GenerateHTML(outputFile, results, "rhel9")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		htmlContent := string(content)

		if !containsSubstring(htmlContent, "CIS Benchmark Scan Report") {
			t.Error("expected title in report")
		}
		if !containsSubstring(htmlContent, "rhel9") {
			t.Error("expected OS type in report")
		}
		if !containsSubstring(htmlContent, "1.1.1") {
			t.Error("expected check ID in report")
		}
		if !containsSubstring(htmlContent, checker.ResultPass) {
			t.Error("expected PASS result in report")
		}
		if !containsSubstring(htmlContent, checker.ResultFail) {
			t.Error("expected FAIL result in report")
		}
	})

	t.Run("with escape characters", func(t *testing.T) {
		results := []checker.Result{
			{ID: "1", Description: "<script>alert('xss')</script>", Result: checker.ResultPass},
		}

		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "report.html")

		err := GenerateHTML(outputFile, results, "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		if containsSubstring(string(content), "<script>") {
			t.Error("HTML should be escaped")
		}
		if !containsSubstring(string(content), "&lt;script&gt;") {
			t.Error("expected escaped HTML entities")
		}
	})

	t.Run("summary calculation", func(t *testing.T) {
		results := []checker.Result{
			{ID: "1", Result: checker.ResultPass},
			{ID: "2", Result: checker.ResultPass},
			{ID: "3", Result: checker.ResultFail},
			{ID: "4", Result: checker.ResultError + ": timeout"},
		}

		summary := calculateSummary(results)

		if summary.Total != 4 {
			t.Errorf("expected Total 4, got %d", summary.Total)
		}
		if summary.Pass != 2 {
			t.Errorf("expected Pass 2, got %d", summary.Pass)
		}
		if summary.Fail != 1 {
			t.Errorf("expected Fail 1, got %d", summary.Fail)
		}
		if summary.Error != 1 {
			t.Errorf("expected Error 1, got %d", summary.Error)
		}
		if summary.PassPercent != 50 {
			t.Errorf("expected PassPercent 50, got %d", summary.PassPercent)
		}
	})

	t.Run("empty results", func(t *testing.T) {
		results := []checker.Result{}

		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "report.html")

		err := GenerateHTML(outputFile, results, "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		summary := calculateSummary(results)
		if summary.Total != 0 {
			t.Errorf("expected Total 0, got %d", summary.Total)
		}
		if summary.PassPercent != 0 {
			t.Errorf("expected PassPercent 0, got %d", summary.PassPercent)
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		results := []checker.Result{{ID: "1", Result: checker.ResultPass}}

		err := GenerateHTML("/nonexistent/directory/report.html", results, "test")
		if err == nil {
			t.Error("expected error for invalid path")
		}
	})
}

func TestGenerateJSON(t *testing.T) {
	t.Run("basic report", func(t *testing.T) {
		results := []checker.Result{
			{ID: "1.1.1", Description: "Test check", Command: "echo test", Expected: "test", Result: checker.ResultPass, Output: "test"},
			{ID: "1.1.2", Description: "Failed check", Command: "echo fail", Expected: "expected", Result: checker.ResultFail, Output: "fail"},
		}

		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "report.json")

		err := GenerateJSON(outputFile, results, "ubuntu2204")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		var report Report
		if err := json.Unmarshal(content, &report); err != nil {
			t.Fatalf("failed to parse JSON: %v", err)
		}

		if report.OSType != "ubuntu2204" {
			t.Errorf("expected OSType 'ubuntu2204', got %s", report.OSType)
		}
		if report.Summary.Total != 2 {
			t.Errorf("expected Total 2, got %d", report.Summary.Total)
		}
		if report.Summary.Pass != 1 {
			t.Errorf("expected Pass 1, got %d", report.Summary.Pass)
		}
		if report.Summary.Fail != 1 {
			t.Errorf("expected Fail 1, got %d", report.Summary.Fail)
		}
		if len(report.Results) != 2 {
			t.Errorf("expected 2 results, got %d", len(report.Results))
		}
	})

	t.Run("with remediation", func(t *testing.T) {
		results := []checker.Result{
			{ID: "1", Description: "Test", Result: checker.ResultFail, Remediation: "Fix this issue"},
		}

		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "report.json")

		err := GenerateJSON(outputFile, results, "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		var report Report
		if err := json.Unmarshal(content, &report); err != nil {
			t.Fatalf("failed to parse JSON: %v", err)
		}

		if report.Results[0].Remediation != "Fix this issue" {
			t.Errorf("expected Remediation 'Fix this issue', got %s", report.Results[0].Remediation)
		}
	})

	t.Run("valid JSON structure", func(t *testing.T) {
		results := []checker.Result{
			{ID: "1", Description: "Test", Result: checker.ResultPass},
		}

		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "report.json")

		err := GenerateJSON(outputFile, results, "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		if !containsSubstring(string(content), "\"timestamp\"") {
			t.Error("expected timestamp field in JSON")
		}
		if !containsSubstring(string(content), "\"os_type\"") {
			t.Error("expected os_type field in JSON")
		}
		if !containsSubstring(string(content), "\"summary\"") {
			t.Error("expected summary field in JSON")
		}
		if !containsSubstring(string(content), "\"results\"") {
			t.Error("expected results field in JSON")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		results := []checker.Result{{ID: "1", Result: checker.ResultPass}}

		err := GenerateJSON("/nonexistent/directory/report.json", results, "test")
		if err == nil {
			t.Error("expected error for invalid path")
		}
	})
}

func TestCalculateSummary(t *testing.T) {
	tests := []struct {
		name          string
		results       []checker.Result
		expectedTotal int
		expectedPass  int
		expectedFail  int
		expectedError int
	}{
		{
			name:          "empty",
			results:       []checker.Result{},
			expectedTotal: 0,
			expectedPass:  0,
			expectedFail:  0,
			expectedError: 0,
		},
		{
			name: "all pass",
			results: []checker.Result{
				{Result: checker.ResultPass},
				{Result: checker.ResultPass},
				{Result: checker.ResultPass},
			},
			expectedTotal: 3,
			expectedPass:  3,
			expectedFail:  0,
			expectedError: 0,
		},
		{
			name: "mixed results",
			results: []checker.Result{
				{Result: checker.ResultPass},
				{Result: checker.ResultFail},
				{Result: checker.ResultError + ": timeout"},
				{Result: checker.ResultPass},
			},
			expectedTotal: 4,
			expectedPass:  2,
			expectedFail:  1,
			expectedError: 1,
		},
		{
			name: "all fail",
			results: []checker.Result{
				{Result: checker.ResultFail},
				{Result: checker.ResultFail},
			},
			expectedTotal: 2,
			expectedPass:  0,
			expectedFail:  2,
			expectedError: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := calculateSummary(tt.results)

			if summary.Total != tt.expectedTotal {
				t.Errorf("Total: expected %d, got %d", tt.expectedTotal, summary.Total)
			}
			if summary.Pass != tt.expectedPass {
				t.Errorf("Pass: expected %d, got %d", tt.expectedPass, summary.Pass)
			}
			if summary.Fail != tt.expectedFail {
				t.Errorf("Fail: expected %d, got %d", tt.expectedFail, summary.Fail)
			}
			if summary.Error != tt.expectedError {
				t.Errorf("Error: expected %d, got %d", tt.expectedError, summary.Error)
			}
		})
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"<script>", "&lt;script&gt;"},
		{"a & b", "a &amp; b"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"'single'", "&#39;single&#39;"},
		{"<>&\"'", "&lt;&gt;&amp;&quot;&#39;"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeHTML(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
