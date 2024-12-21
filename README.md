    # CIS Benchmark Scanner

A comprehensive security compliance scanner for Linux systems, implementing CIS (Center for Internet Security) Benchmark checks.

## Features

- **Multiple Distribution Support**:
  - RHEL 8 (Complete CIS Level 1 Implementation)
  - RHEL 9
  - SLES 15
  - Ubuntu (planned)

- **Security Checks**:
  - Initial System Settings
  - Services Configuration
  - Network Security
  - Auditing and Logging
  - Access Controls
  - System Maintenance
  - User and Group Settings

- **Comprehensive Coverage**:
  - File System Configuration
  - Boot Settings
  - Process Accounting
  - User Environment
  - Password Policies
  - Network Parameters
  - Audit Rules
  - Time Synchronization

## Installation

```bash
go install github.com/diceone/cis_benchmark_scanner@latest
```

## Usage

```bash
# Run all checks
cis_benchmark_scanner scan

# Run checks for specific OS
cis_benchmark_scanner scan --os rhel8

# Run specific check by ID
cis_benchmark_scanner scan --check 1.1.1
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

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

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

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Center for Internet Security (CIS) for their comprehensive security benchmarks
- The Go community for excellent security-related packages
- Contributors who help improve and maintain this project
