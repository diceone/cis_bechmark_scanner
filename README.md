    # CIS Benchmark Scanner

A comprehensive security compliance scanner for Linux systems, implementing CIS (Center for Internet Security) Benchmark checks.

## Features

- **Multiple Distribution Support**:
  - RHEL 8 (Complete CIS Level 1 Implementation)
  - RHEL 9 (Complete CIS Level 1 Implementation + Enterprise Features)
  - SLES 15
  - Ubuntu 22.04 LTS (Complete CIS Level 1 Implementation)

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

- **Ubuntu 22.04 Features**:
  - AppArmor Profiles
  - Snap Security
  - UFW Configuration
  - Systemd-journald
  - Time Synchronization
  - Package Verification

## Installation

```bash
go install github.com/diceone/cis_benchmark_scanner@latest
```

## Usage

```bash
# Run all checks
cis_benchmark_scanner scan

# Run checks for specific OS
cis_benchmark_scanner scan --os rhel9
cis_benchmark_scanner scan --os ubuntu2204

# Run specific check by ID
cis_benchmark_scanner scan --check 1.1.1

# Generate HTML report
cis_benchmark_scanner scan --report
```

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

### RHEL 9 Enterprise Features

#### Hardware Security
- TPM2 Device Management
- Secure Boot Configuration
- IMA Policy Settings

#### Disk Encryption
- LUKS2 with TPM2
- Automated Decryption
- Strong Cipher Settings

#### System Security
- SELinux Enforcing Mode
- FIPS 140-3 Compliance
- USBGuard Protection
- Enterprise Authentication

### Ubuntu 22.04 Features

#### Access Control
- AppArmor Profiles
- Snap Package Security
- Service Hardening
- SSH Configuration

#### System Security
- UFW Firewall Rules
- Systemd Services
- Time Synchronization
- Package Verification

#### Auditing & Logging
- Audit Rules
- Journald Settings
- Log Forwarding
- File Integrity

## Security Checks Implementation

Each security check is defined in YAML format with:
- Unique identifier
- Description
- Command to execute
- Expected result
- Optional remediation steps

Example:
```yaml
- id: "1.1.1"
  description: "Ensure mounting of cramfs filesystems is disabled"
  command: "modprobe -n -v cramfs"
  expected: "install /bin/true"
```

## HTML Reports

The scanner generates detailed reports including:
- Overall compliance status
- Individual check results
- Remediation suggestions
- System information

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
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Center for Internet Security (CIS) for their comprehensive security benchmarks
- The Go community for excellent security-related packages
- Contributors who help improve and maintain this project
