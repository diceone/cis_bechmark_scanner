# CIS Benchmark Scanner

A comprehensive security compliance scanner for Linux systems, implementing CIS (Center for Internet Security) Benchmark checks.

## Features

- **Multiple Distribution Support**:
  - RHEL 8 (Complete CIS Level 1 Implementation)
  - RHEL 9 (Complete CIS Level 1 Implementation + Enterprise Features)
  - SLES 15
  - Ubuntu 22.04 LTS (Complete CIS Level 1 Implementation)
  - Ubuntu 24.04 LTS

- **Security Checks**:
  - Initial System Settings
  - Services Configuration
  - Network Security
  - Auditing and Logging
  - Access Controls
  - System Maintenance
  - User and Group Settings

- **Enterprise Features** (RHEL 9):
  - TPM2 and Secure Boot
  - LUKS2 Encryption
  - FIPS 140-3 Compliance
  - SELinux Controls
  - USBGuard Security
  - Enterprise Authentication

- **Ubuntu 22.04/24.04 Features**:
  - AppArmor Profiles
  - Snap Security
  - UFW Configuration
  - Systemd-journald
  - Time Synchronization
  - Package Verification

## Installation

```bash
go install github.com/diceone/cis_benchmark_scanner/cmd/cis_scanner@latest
```

Or download the latest release from [GitHub Releases](https://github.com/diceone/cis_benchmark_scanner/releases).

## Usage

```bash
# Show help
cis_scanner -h

# Show version
cis_scanner -version

# List all available checks for an OS
cis_scanner -list -os rhel9

# Run all checks for an OS
cis_scanner -scan -os rhel9

# Run checks for different OS
cis_scanner -scan -os ubuntu2204
cis_scanner -scan -os sles15

# Run specific check(s) by ID
cis_scanner -scan -os rhel9 -check 1.1.1.1
cis_scanner -scan -os rhel9 -check 1.1.1.1,1.1.1.2,1.1.1.3

# Generate HTML report
cis_scanner -scan -os rhel9 -report

# Generate JSON report
cis_scanner -scan -os rhel9 -json

# Generate both reports
cis_scanner -scan -os rhel9 -report -json

# Specify output directory for reports
cis_scanner -scan -os rhel9 -report -output ./reports

# Run with parallel workers for faster execution
cis_scanner -scan -os rhel9 -parallel 4

# Set timeout for each check (default: 30 seconds)
cis_scanner -scan -os rhel9 -timeout 60
```

## Command Line Options

| Flag | Description |
|------|-------------|
| `-h`, `-help` | Display help message |
| `-version` | Display version information |
| `-os <type>` | Target OS (rhel8, rhel9, sles15, ubuntu2204, ubuntu2404) |
| `-scan` | Run security checks |
| `-list` | List available checks for specified OS |
| `-check <ids>` | Run specific check(s) by ID (comma-separated) |
| `-parallel <n>` | Number of parallel workers (default: 1) |
| `-timeout <sec>` | Timeout per check in seconds (default: 30) |
| `-report` | Generate HTML report |
| `-json` | Generate JSON report |
| `-output <dir>` | Output directory for reports (default: current) |

## Check Categories

### 1. Initial Setup
- File System Configuration
- Software Updates
- Boot Settings
- Process Hardening
- Mandatory Access Control

### 2. Services
- inetd Services
- Special Purpose Services
- Service Clients
- Time Synchronization

### 3. Network Configuration
- Network Parameters
- IPv4 & IPv6 Settings
- TCP Wrappers
- Firewall Configuration
- Wireless Interfaces

### 4. Logging and Auditing
- Configure System Accounting
- Configure Logging
- Ensure rsyslog Configured
- Configure journald

### 5. Access, Authentication and Authorization
- Configure time-based job schedulers
- Configure SSH
- Configure PAM
- User Accounts and Environment

### 6. System Maintenance
- System File Permissions
- User and Group Settings
- Root Login Controls
- Default umask Configuration

## Security Checks Implementation

Each security check is defined in YAML format with:
- `id`: Unique identifier (required)
- `description`: Human-readable description (required)
- `command`: Shell command to execute (required)
- `expected`: Expected output to match (required, can be empty for "no output" checks)
- `remediation`: Optional instructions for fixing the issue

Example:
```yaml
- id: "1.1.1"
  description: "Ensure mounting of cramfs filesystems is disabled"
  command: "modprobe -n -v cramfs"
  expected: "install /bin/true"
  remediation: "Add 'install cramfs /bin/true' to /etc/modprobe.d/cramfs.conf"
```

### Matching Behavior

The scanner uses substring matching by default:
- If `expected` is empty and command exits with non-zero: **PASS** (useful for "should not exist" checks)
- If `expected` is found in command output: **PASS**
- Otherwise: **FAIL**

## Reports

### HTML Report
Generates a detailed, styled HTML report with:
- Summary statistics (total, passed, failed, errors)
- Pass percentage
- Individual check results with color coding

### JSON Report
Generates a machine-readable JSON report suitable for:
- CI/CD integration
- Automated processing
- Historical tracking

Example JSON structure:
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "os_type": "rhel9",
  "summary": {
    "total": 200,
    "pass": 185,
    "fail": 12,
    "error": 3,
    "pass_percent": 92
  },
  "results": [
    {
      "id": "1.1.1",
      "description": "Ensure cramfs is disabled",
      "command": "modprobe -n -v cramfs",
      "expected": "install /bin/true",
      "result": "PASS",
      "output": "install /bin/true"
    }
  ]
}
```

## Output

The scanner outputs:
- Real-time progress to console
- Pass/fail status for each check
- Actual vs expected output for failures
- Summary statistics
- Exit code 0 if all checks pass, 1 otherwise

## Project Structure

```
cis_benchmark_scanner/
├── cmd/cis_scanner/       # Main entry point
│   └── main.go
├── internal/
│   ├── checker/           # Core checking logic
│   │   ├── checker.go     # Check execution and validation
│   │   ├── checker_test.go
│   │   ├── executor.go     # Command execution with timeout
│   │   └── executor_test.go
│   └── report/            # Report generation
│       ├── report.go      # HTML and JSON reports
│       └── report_test.go
├── cis_checks_rhel8.yaml   # RHEL 8 checks
├── cis_checks_rhel9.yaml   # RHEL 9 checks
├── cis_checks_sles15.yaml # SLES 15 checks
├── cis_checks_ubuntu2204.yaml  # Ubuntu 22.04 checks
├── cis_checks_ubuntu2404.yaml  # Ubuntu 24.04 checks
├── go.mod
├── go.sum
└── README.md
```

## Changelog

### v0.5.0 (Refactored)

**Architecture:**
- Separated CLI, checker logic, and report generation into distinct packages
- Added parallel execution support for faster scanning
- Added context-based timeout and cancellation handling

**New Features:**
- `-parallel N` flag for concurrent check execution
- `-check <ids>` flag to run specific checks by ID
- `-json` flag for JSON report output
- `-timeout <sec>` flag for per-check timeout configuration
- `-version` flag displays build information
- Remediation field support in YAML check definitions

**Improvements:**
- Comprehensive test coverage for internal packages
- Better HTML report styling with summary statistics
- Pass percentage calculation in reports
- Exit code semantics: 0 = all pass, 1 = any fail/error

**Bug Fixes:**
- Fixed duplicate check IDs in RHEL 8 and RHEL 9 YAML files
- Fixed module path consistency (`github.com/diceone/...`)
- Improved error handling for command timeouts

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Development

### Requirements
- Go 1.21 or later
- YAML knowledge for check definitions
- Linux system for testing

### Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/checker/...

# Run with verbose output
go test -v ./...
```

### Adding New Checks
1. Find the appropriate YAML file for the target OS
2. Add the check definition following the existing format
3. Ensure the check ID is unique within the file
4. Test the check command manually before committing

## Security Considerations

- Commands are executed via shell (`/bin/sh -c`)
- Ensure YAML check files come from trusted sources
- Run with appropriate privileges for system-level checks
- Review reports before sharing (may contain system information)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Center for Internet Security (CIS) for their comprehensive security benchmarks
- The Go community for excellent security-related packages
- Contributors who help improve and maintain this project