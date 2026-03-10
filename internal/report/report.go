package report

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/diceone/cis_benchmark_scanner/internal/checker"
)

type ReportFormat string

const (
	FormatHTML ReportFormat = "html"
	FormatJSON ReportFormat = "json"
)

type Report struct {
	Timestamp string           `json:"timestamp"`
	OSType    string           `json:"os_type"`
	Summary   Summary          `json:"summary"`
	Results   []checker.Result `json:"results"`
}

type Summary struct {
	Total       int `json:"total"`
	Pass        int `json:"pass"`
	Fail        int `json:"fail"`
	Error       int `json:"error"`
	PassPercent int `json:"pass_percent"`
}

func GenerateHTML(filename string, results []checker.Result, osType string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating report file: %w", err)
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	summary := calculateSummary(results)

	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString("<meta charset=\"utf-8\">\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	b.WriteString("<title>CIS Benchmark Scan Report</title>\n")
	b.WriteString("<style>\n")
	b.WriteString("body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }\n")
	b.WriteString(".container { max-width: 1200px; margin: 0 auto; }\n")
	b.WriteString("h1 { color: #333; margin-bottom: 10px; }\n")
	b.WriteString(".meta { color: #666; margin-bottom: 20px; }\n")
	b.WriteString(".summary { display: flex; gap: 20px; margin-bottom: 30px; flex-wrap: wrap; }\n")
	b.WriteString(".stat-box { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); min-width: 150px; }\n")
	b.WriteString(".stat-box h3 { margin: 0 0 10px; color: #666; font-size: 14px; }\n")
	b.WriteString(".stat-box .value { font-size: 32px; font-weight: bold; }\n")
	b.WriteString(".stat-box.total .value { color: #333; }\n")
	b.WriteString(".stat-box.pass .value { color: #28a745; }\n")
	b.WriteString(".stat-box.fail .value { color: #dc3545; }\n")
	b.WriteString(".stat-box.error .value { color: #ffc107; }\n")
	b.WriteString(".progress-bar { height: 8px; background: #e9ecef; border-radius: 4px; margin-top: 10px; overflow: hidden; }\n")
	b.WriteString(".progress-fill { height: 100%; background: #28a745; transition: width 0.3s; }\n")
	b.WriteString("table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }\n")
	b.WriteString("th { background: #f8f9fa; padding: 12px; text-align: left; font-weight: 600; border-bottom: 2px solid #e9ecef; }\n")
	b.WriteString("td { padding: 12px; border-bottom: 1px solid #e9ecef; }\n")
	b.WriteString("tr:hover { background: #f8f9fa; }\n")
	b.WriteString(".pass { color: #28a745; font-weight: 600; }\n")
	b.WriteString(".fail { color: #dc3545; font-weight: 600; }\n")
	b.WriteString(".error { color: #ffc107; font-weight: 600; }\n")
	b.WriteString(".id { font-family: monospace; font-size: 13px; }\n")
	b.WriteString("</style>\n</head>\n<body>\n")
	b.WriteString("<div class=\"container\">\n")
	b.WriteString("<h1>CIS Benchmark Scan Report</h1>\n")
	b.WriteString(fmt.Sprintf("<div class=\"meta\">\n"))
	b.WriteString(fmt.Sprintf("<div>Generated: %s</div>\n", timestamp))
	b.WriteString(fmt.Sprintf("<div>Target OS: %s</div>\n", osType))
	b.WriteString("</div>\n")
	b.WriteString("<div class=\"summary\">\n")
	b.WriteString(fmt.Sprintf("<div class=\"stat-box total\"><h3>Total Checks</h3><div class=\"value\">%d</div></div>\n", summary.Total))
	b.WriteString(fmt.Sprintf("<div class=\"stat-box pass\"><h3>Passed</h3><div class=\"value\">%d</div>\n", summary.Pass))
	if summary.Total > 0 {
		b.WriteString(fmt.Sprintf("<div class=\"progress-bar\"><div class=\"progress-fill\" style=\"width: %d%%\"></div></div>\n", summary.PassPercent))
	}
	b.WriteString("</div>\n")
	b.WriteString(fmt.Sprintf("<div class=\"stat-box fail\"><h3>Failed</h3><div class=\"value\">%d</div></div>\n", summary.Fail))
	b.WriteString(fmt.Sprintf("<div class=\"stat-box error\"><h3>Errors</h3><div class=\"value\">%d</div></div>\n", summary.Error))
	b.WriteString("</div>\n")
	b.WriteString("<table>\n<thead>\n<tr><th>ID</th><th>Description</th><th>Result</th></tr>\n</thead>\n<tbody>\n")

	for _, result := range results {
		resultClass := "fail"
		if result.Result == checker.ResultPass {
			resultClass = "pass"
		} else if strings.HasPrefix(result.Result, checker.ResultError) {
			resultClass = "error"
		}
		b.WriteString(fmt.Sprintf("<tr><td class=\"id\">%s</td><td>%s</td><td class=\"%s\">%s</td></tr>\n",
			escapeHTML(result.ID), escapeHTML(result.Description), resultClass, escapeHTML(result.Result)))
	}

	b.WriteString("</tbody>\n</table>\n")
	b.WriteString("</div>\n</body>\n</html>\n")

	_, err = file.WriteString(b.String())
	if err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	return nil
}

func GenerateJSON(filename string, results []checker.Result, osType string) error {
	summary := calculateSummary(results)

	report := Report{
		Timestamp: time.Now().Format(time.RFC3339),
		OSType:    osType,
		Summary:   summary,
		Results:   results,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("writing json file: %w", err)
	}

	return nil
}

func calculateSummary(results []checker.Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch {
		case r.Result == checker.ResultPass:
			s.Pass++
		case strings.HasPrefix(r.Result, checker.ResultError):
			s.Error++
		default:
			s.Fail++
		}
	}
	if s.Total > 0 {
		s.PassPercent = (s.Pass * 100) / s.Total
	}
	return s
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
