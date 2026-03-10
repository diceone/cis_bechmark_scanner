package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/diceone/cis_benchmark_scanner/internal/checker"
	"github.com/diceone/cis_benchmark_scanner/internal/report"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	var (
		showHelp    bool
		scanFlag    bool
		reportFlag  bool
		jsonFlag    bool
		listFlag    bool
		versionFlag bool
		osType      string
		checkIDs    string
		parallel    int
		timeout     int
		outputDir   string
	)

	flag.BoolVar(&showHelp, "help", false, "Display this help message")
	flag.BoolVar(&showHelp, "h", false, "Display this help message")
	flag.BoolVar(&scanFlag, "scan", false, "Run all checks")
	flag.BoolVar(&reportFlag, "report", false, "Generate an HTML report after scanning")
	flag.BoolVar(&jsonFlag, "json", false, "Generate a JSON report (use with -report)")
	flag.BoolVar(&listFlag, "list", false, "List all available checks")
	flag.BoolVar(&versionFlag, "version", false, "Display version information")
	flag.StringVar(&osType, "os", "", "Target OS for benchmark (e.g., rhel8, rhel9, sles15, ubuntu2204, ubuntu2404)")
	flag.StringVar(&checkIDs, "check", "", "Run specific check(s) by ID (comma-separated)")
	flag.IntVar(&parallel, "parallel", 1, "Number of parallel workers (default: 1)")
	flag.IntVar(&timeout, "timeout", 30, "Timeout in seconds for each check (default: 30)")
	flag.StringVar(&outputDir, "output", ".", "Output directory for reports (default: current directory)")
	flag.Parse()

	if versionFlag {
		fmt.Printf("cis_scanner version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built: %s\n", buildDate)
		fmt.Printf("  go version: %s\n", runtime.Version())
		os.Exit(0)
	}

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if osType == "" && (scanFlag || listFlag) {
		fmt.Fprintf(os.Stderr, "Error: -os flag is required\n\n")
		printUsage()
		os.Exit(1)
	}

	if !scanFlag && !listFlag {
		printUsage()
		os.Exit(1)
	}

	executor := checker.NewCommandExecutor(timeout)
	matchMode := checker.MatchContains
	ckr := checker.NewChecker(executor, checker.WithParallel(parallel), checker.WithMatchMode(matchMode))

	checksFile := fmt.Sprintf("cis_checks_%s.yaml", osType)
	if err := ckr.LoadChecks(checksFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: No checks file found for OS '%s'\n", osType)
			fmt.Fprintf(os.Stderr, "Supported OS: rhel8, rhel9, sles15, ubuntu2204, ubuntu2404\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error loading checks: %v\n", err)
		}
		os.Exit(1)
	}

	if listFlag {
		listChecks(ckr)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintf(os.Stderr, "\nReceived interrupt signal, stopping...\n")
		cancel()
	}()

	var ids []string
	if checkIDs != "" {
		ids = strings.Split(checkIDs, ",")
		for i, id := range ids {
			ids[i] = strings.TrimSpace(id)
		}
	}

	fmt.Printf("Starting CIS Benchmark Scan for %s\n", strings.ToUpper(osType))
	fmt.Printf("Running with %d worker(s), timeout: %ds\n\n", parallel, timeout)

	results := ckr.RunChecks(ctx, ids)

	summary := printResults(results)

	if reportFlag || jsonFlag {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}

		baseFilename := fmt.Sprintf("cis_scan_report_%s", osType)

		if reportFlag {
			htmlFile := filepath.Join(outputDir, baseFilename+".html")
			if err := report.GenerateHTML(htmlFile, results, osType); err != nil {
				fmt.Fprintf(os.Stderr, "Error generating HTML report: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\nHTML report generated: %s\n", htmlFile)
		}

		if jsonFlag {
			jsonFile := filepath.Join(outputDir, baseFilename+".json")
			if err := report.GenerateJSON(jsonFile, results, osType); err != nil {
				fmt.Fprintf(os.Stderr, "Error generating JSON report: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("JSON report generated: %s\n", jsonFile)
		}
	}

	if summary.fail > 0 || summary.errCount > 0 {
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("CIS Benchmark Scanner")
	fmt.Println()
	fmt.Println("A comprehensive security compliance scanner for Linux systems")
	fmt.Println("implementing CIS (Center for Internet Security) Benchmark checks.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  cis_scanner [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, -help          Display this help message")
	fmt.Println("  -version           Display version information")
	fmt.Println("  -os <type>         Target OS for benchmark")
	fmt.Println("                     Supported: rhel8, rhel9, sles15, ubuntu2204, ubuntu2404")
	fmt.Println("  -scan              Run all checks")
	fmt.Println("  -check <ids>       Run specific check(s) by ID (comma-separated)")
	fmt.Println("  -list              List all available checks for the specified OS")
	fmt.Println("  -parallel <n>      Number of parallel workers (default: 1)")
	fmt.Println("  -timeout <sec>     Timeout in seconds for each check (default: 30)")
	fmt.Println("  -report            Generate an HTML report")
	fmt.Println("  -json              Generate a JSON report (use with -report)")
	fmt.Println("  -output <dir>      Output directory for reports (default: current directory)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cis_scanner -scan -os rhel9")
	fmt.Println("  cis_scanner -scan -os ubuntu2204 -report -json")
	fmt.Println("  cis_scanner -scan -os rhel8 -parallel 4")
	fmt.Println("  cis_scanner -check 1.1.1.1,1.1.1.2 -os rhel9")
	fmt.Println("  cis_scanner -list -os ubuntu2404")
}

func listChecks(ckr *checker.Checker) {
	checks := ckr.GetChecks()
	fmt.Printf("Available CIS Benchmark Checks (%d total):\n\n", len(checks))

	sort.Slice(checks, func(i, j int) bool {
		return checks[i].ID < checks[j].ID
	})

	for _, check := range checks {
		fmt.Printf("  %s: %s\n", check.ID, check.Description)
	}
}

type summary struct {
	total, pass, fail, errCount int
}

func printResults(results []checker.Result) summary {
	var s summary
	s.total = len(results)

	for _, r := range results {
		fmt.Printf("Checking %s: %s\n", r.ID, r.Description)

		switch {
		case r.Result == checker.ResultPass:
			s.pass++
			fmt.Printf("  Result: %s\n\n", r.Result)
		case strings.HasPrefix(r.Result, checker.ResultError):
			s.errCount++
			fmt.Printf("  Result: %s\n", r.Result)
			if r.Output != "" {
				fmt.Printf("  Output: %s\n\n", truncate(r.Output, 200))
			} else {
				fmt.Println()
			}
		default:
			s.fail++
			fmt.Printf("  Result: %s\n", r.Result)
			if r.Expected != "" {
				fmt.Printf("  Expected: %s\n", truncate(r.Expected, 100))
			}
			if r.Output != "" {
				fmt.Printf("  Got: %s\n", truncate(r.Output, 200))
			}
			if r.Remediation != "" {
				fmt.Printf("  Remediation: %s\n", r.Remediation)
			}
			fmt.Println()
		}
	}

	fmt.Println("Scan Summary:")
	fmt.Printf("  Total:  %d\n", s.total)
	fmt.Printf("  Passed: %d\n", s.pass)
	fmt.Printf("  Failed: %d\n", s.fail)
	fmt.Printf("  Errors: %d\n", s.errCount)

	return s
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
