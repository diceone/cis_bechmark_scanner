    # CIS Benchmark Scanner for RHEL 8

A command-line tool that scans RHEL 8 systems for compliance with CIS (Center for Internet Security) benchmarks. The tool performs various security checks and generates detailed HTML reports.

## Features

- 30+ security checks based on CIS benchmarks
- HTML report generation with color-coded results
- YAML-based configuration for easy customization
- Support for both AMD64 and ARM64 architectures
- Available as binary, DEB, and RPM packages

## Installation

### Using Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/diceone/cis_benchmark_scanner/releases) for your platform:

```bash
# Linux (AMD64)
curl -LO https://github.com/diceone/cis_benchmark_scanner/releases/latest/download/cis_benchmark_scanner_Linux_x86_64.tar.gz
tar xzf cis_benchmark_scanner_Linux_x86_64.tar.gz

# Linux (ARM64)
curl -LO https://github.com/diceone/cis_benchmark_scanner/releases/latest/download/cis_benchmark_scanner_Linux_arm64.tar.gz
tar xzf cis_benchmark_scanner_Linux_arm64.tar.gz
```

### Using Package Managers

For Debian/Ubuntu:
```bash
curl -LO https://github.com/diceone/cis_benchmark_scanner/releases/latest/download/cis-benchmark-scanner_linux_amd64.deb
sudo dpkg -i cis-benchmark-scanner_linux_amd64.deb
```

For RHEL/CentOS:
```bash
curl -LO https://github.com/diceone/cis_benchmark_scanner/releases/latest/download/cis-benchmark-scanner_linux_amd64.rpm
sudo rpm -i cis-benchmark-scanner_linux_amd64.rpm
```

### Building from Source

```bash
git clone https://github.com/diceone/cis_benchmark_scanner.git
cd cis_benchmark_scanner
go build -o cis_scanner cis_scanner.go
```

## Usage

1. List all available checks for a specific OS:
```bash
./cis_scanner -list -os rhel8
```

2. Run security checks for RHEL 8:
```bash
./cis_scanner -scan -os rhel8
```

3. Run checks and generate HTML report for Ubuntu 20.04:
```bash
./cis_scanner -scan -report -os ubuntu20
```

The HTML report will be generated as `cis_scan_report_[os].html` in the current directory.

## Configuration

Security checks are defined in OS-specific YAML files (e.g., `cis_checks_rhel8.yaml`, `cis_checks_ubuntu20.yaml`). Each check has:
- Unique ID
- Description
- Command to execute
- Expected result

Example check from `cis_checks_rhel8.yaml`:
```yaml
- id: "1.1.1.1"
  description: "Ensure mounting of cramfs filesystems is disabled"
  command: "modprobe -n -v cramfs"
  expected: "install /bin/true"
```

## Supported Operating Systems

Currently supported OS configurations:
- RHEL 8 (`-os rhel8`)
  - Complete CIS Level 1 Server benchmark checks
  - Filesystem security
  - Network configuration
  - System hardening

- RHEL 9 (`-os rhel9`)
  - Enhanced SELinux configuration
  - Firewalld and nftables integration
  - RHEL 9 specific package verification
  - Systemd-based service management
  - Advanced crypto policies
  - DNF package management checks
  - Container security (Podman)
  - TPM2 and LUKS encryption
  - USBGuard protection
  - DNSSEC and DNS-over-TLS
  - SCAP compliance tools
  - Enterprise authentication (SSSD)
  - Network security enhancements

- Ubuntu 22.04 LTS (`-os ubuntu2204`)
  - Complete CIS Level 1 Server benchmark checks
  - AppArmor security
  - UFW firewall configuration
  - System hardening
  - Package verification using dpkg
  - Ubuntu-specific filesystem checks

Additional OS support can be added by creating new YAML files following the naming convention `cis_checks_[os].yaml`.

## Usage Examples

1. Run security checks for RHEL 9:
```bash
./cis_scanner -scan -os rhel9
```

2. Generate HTML report for Ubuntu 22.04:
```bash
./cis_scanner -scan -report -os ubuntu2204
```

3. List available checks for RHEL 9:
```bash
./cis_scanner -list -os rhel9
```

## Security Checks

The scanner includes checks for:
- Filesystem configurations
- System software
- Network settings
- Access control
- Password policies
- System auditing
- Firewall configuration
- SELinux settings
- System hardening
- File permissions

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- CIS Benchmarks™ by Center for Internet Security, Inc
- RHEL 8 Security Technical Implementation Guide (STIG)

### Steps to Create the Git Repository

1. Initialize the Git repository:
   ```bash
   git init cis-benchmark-scanner
   cd cis-benchmark-scanner
   ```

2. Add the files:
   ```bash
   git add .gitignore README.md go.mod go.sum cis_scanner.go cis_checks.yaml
   mkdir reports
   ```

3. Commit the changes:
   ```bash
   git commit -m "Initial commit with scanner implementation and YAML config"
   ```

4. Optionally, push to a remote repository:
   ```bash
   git remote add origin <repository_url>
   git branch -M main
   git push -u origin main
   ```

3. **`go.mod`**:
   - The Go module file, generated with `go mod init`.
   - Tracks dependencies for the project.
   - Example content:
     ```
     module cis-benchmark-scanner

     go 1.20

     require gopkg.in/yaml.v3 v3.0.1
     ```

4. **`go.sum`**:
   - The Go checksum file for tracking exact dependency versions.

5. **`cis_scanner.go`**:
   - The main Go program file containing the CLI logic, YAML loader, and HTML report generator.

6. **`cis_checks.yaml`**:
   - The YAML configuration file defining all the CIS benchmark checks.
   - Example content:
     ```yaml
     - id: "1.1.1.1"
       description: "Ensure mounting of 'cramfs' is disabled"
       command: "modprobe -n -v cramfs"
       expected: "install /bin/true"

     - id: "2.2.2"
       description: "Ensure 'xinetd' is not installed"
       command: "rpm -q xinetd"
       expected: "package xinetd is not installed"

     - id: "3.1.1"
       description: "Ensure IP forwarding is disabled"
       command: "sysctl net.ipv4.ip_forward"
       expected: "net.ipv4.ip_forward = 0"
     ```

7. **`reports/`**:
   - Directory for storing generated HTML reports.
   - Example: `cis_scan_report.html`.

8. **`cis_scan_report.html`**:
   - A generated HTML report containing the scan results.
   - This file is created dynamically during execution and saved to the `reports/` directory.
