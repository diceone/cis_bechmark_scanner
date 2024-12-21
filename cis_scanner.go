package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Check represents a CIS Benchmark check.
type Check struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
	Command     string `yaml:"command"`
	Expected    string `yaml:"expected"`
	Result      string // Stores the result (PASS/FAIL/ERROR).
}

// Checks stores all the CIS Benchmark checks.
var Checks []Check

// LoadChecks loads checks from a YAML file.
func LoadChecks(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Checks)
	if err != nil {
		return err
	}
	return nil
}

// runCommand runs a shell command and returns the output.
func runCommand(cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// performCheck validates a single CIS check.
func performCheck(check *Check) {
	output, err := runCommand(check.Command)
	if err != nil {
		check.Result = fmt.Sprintf("ERROR: %v", err)
		return
	}
	if strings.Contains(output, check.Expected) {
		check.Result = "PASS"
	} else {
		check.Result = "FAIL"
	}
}

// generateHTMLReport creates an HTML report summarizing the scan results.
func generateHTMLReport(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the HTML content
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	file.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<title>CIS Benchmark Scan Report</title>\n")
	file.WriteString("<style>\nbody { font-family: Arial, sans-serif; margin: 20px; }\n")
	file.WriteString("table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }\n")
	file.WriteString("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }\n")
	file.WriteString("th { background-color: #f4f4f4; }\n")
	file.WriteString(".pass { color: green; font-weight: bold; }\n")
	file.WriteString(".fail { color: red; font-weight: bold; }\n")
	file.WriteString(".error { color: orange; font-weight: bold; }\n")
	file.WriteString("</style>\n</head>\n<body>\n")
	file.WriteString(fmt.Sprintf("<h1>CIS Benchmark Scan Report</h1>\n<p>Generated on: %s</p>\n", timestamp))
	file.WriteString("<table>\n<tr><th>ID</th><th>Description</th><th>Result</th></tr>\n")

	for _, check := range Checks {
		resultClass := "fail"
		if check.Result == "PASS" {
			resultClass = "pass"
		} else if strings.HasPrefix(check.Result, "ERROR") {
			resultClass = "error"
		}
		file.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td class='%s'>%s</td></tr>\n",
			check.ID, check.Description, resultClass, check.Result))
	}

	file.WriteString("</table>\n</body>\n</html>\n")
	return nil
}

// listChecks displays all available checks.
func listChecks() {
	fmt.Println("Available CIS Benchmark Checks:")
	for _, check := range Checks {
		fmt.Printf("  %s: %s\n", check.ID, check.Description)
	}
}

func main() {
	var (
		showHelp  bool
		listFlag  bool
		scanFlag  bool
		reportFlag bool
		osType    string
	)

	flag.BoolVar(&showHelp, "help", false, "Display this help message")
	flag.BoolVar(&listFlag, "list", false, "List all available checks")
	flag.BoolVar(&scanFlag, "scan", false, "Run all checks")
	flag.BoolVar(&reportFlag, "report", false, "Generate an HTML report after scanning")
	flag.StringVar(&osType, "os", "rhel8", "Target OS for benchmark (e.g., rhel8, ubuntu20)")
	flag.Parse()

	if showHelp {
		fmt.Println("CIS Benchmark Scanner")
		fmt.Println("\nUsage:")
		fmt.Println("  cis_scanner [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  -help      Display this help message")
		fmt.Println("  -list      List all available checks")
		fmt.Println("  -scan      Run all checks")
		fmt.Println("  -report    Generate an HTML report after scanning")
		fmt.Println("  -os        Target OS for benchmark (e.g., rhel8, ubuntu20)")
		fmt.Println("\nExamples:")
		fmt.Println("  cis_scanner -scan -os rhel8")
		fmt.Println("  cis_scanner -list -os ubuntu20")
		fmt.Println("  cis_scanner -scan -report -os rhel8")
		return
	}

	checksFile := fmt.Sprintf("cis_checks_%s.yaml", osType)
	if err := LoadChecks(checksFile); err != nil {
		fmt.Printf("Error loading checks: %v\n", err)
		os.Exit(1)
	}

	if listFlag {
		listChecks()
		return
	}

	if scanFlag {
		fmt.Printf("Starting CIS Benchmark Scan for %s\n\n", strings.ToUpper(osType))
		for i := range Checks {
			fmt.Printf("Checking %s: %s\n", Checks[i].ID, Checks[i].Description)
			performCheck(&Checks[i])
			fmt.Printf("Result: %s\n\n", Checks[i].Result)
		}
		fmt.Println("Scan complete.")

		if reportFlag {
			reportFile := fmt.Sprintf("cis_scan_report_%s.html", osType)
			if err := generateHTMLReport(reportFile); err != nil {
				fmt.Printf("Error generating report: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("HTML report generated: %s\n", reportFile)
		}
		return
	}

	flag.Usage()
}
