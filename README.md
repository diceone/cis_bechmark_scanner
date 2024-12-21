    # CIS Benchmark Scanner

A powerful command-line tool for scanning systems against CIS (Center for Internet Security) benchmarks. This tool supports multiple operating systems and provides detailed HTML reports of compliance status.

## Features

- Multi-OS support (RHEL 8/9, Ubuntu 22.04/24.04, SLES 15)
- Comprehensive security checks based on CIS benchmarks
- Detailed HTML report generation
- Easy-to-use command-line interface
- Configurable check definitions in YAML format
- Real-time scan status updates

## Installation

### Using Pre-built Binaries

Download the latest release from the [releases page](https://github.com/diceone/cis_benchmark_scanner/releases).

For RPM-based systems:
```bash
sudo rpm -i cis_benchmark_scanner_*.rpm
```

For DEB-based systems:
```bash
sudo dpkg -i cis_benchmark_scanner_*.deb
```

### Building from Source

1. Ensure Go 1.21 or later is installed
2. Clone the repository:
```bash
git clone https://github.com/diceone/cis_benchmark_scanner.git
cd cis_benchmark_scanner
```

3. Build the binary:
```bash
go build
```

## Usage

Basic command syntax:
```bash
cis_scanner [options]
```

### Options
- `-scan`: Run security checks
- `-report`: Generate HTML report
- `-list`: List available checks
- `-os`: Specify target OS (rhel8, rhel9, ubuntu2204, ubuntu2404, sles15)
- `-help`: Display help information

### Examples

1. Run security checks for RHEL 9:
```bash
./cis_scanner -scan -os rhel9
```

2. Generate HTML report for Ubuntu 24.04:
```bash
./cis_scanner -scan -report -os ubuntu2404
```

3. List available checks for Ubuntu 22.04:
```bash
./cis_scanner -list -os ubuntu2204
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

- Ubuntu 24.04 LTS (`-os ubuntu2404`)
  - Next-generation security features
  - Nftables-based firewall backend
  - Systemd-oomd memory protection
  - Wayland display server security
  - TPM2 integration
  - FIPS 140-3 compliance options
  - Wireguard VPN support
  - Pipewire audio security
  - Cloud-init hardening
  - Ubuntu Pro security features
  - Snap and LXD container security
  - DNS-over-TLS support
  - Automated service management

- SLES 15 (`-os sles15`)
  - Complete CIS Level 1 Server benchmark checks
  - AppArmor security profiles
  - Wicked Network Manager configuration
  - YAST2 security settings
  - Snapper system snapshots
  - Firewalld integration
  - RPM package verification
  - Time synchronization with chrony
  - System file permissions
  - Password policies
  - Warning banners
  - Enterprise Features:
    - SUSE Manager integration
    - SUSE Connect registration
    - Enterprise storage with Btrfs
    - FIPS 140-2 compliance
    - Security hardening patterns
    - Container security (Podman)
    - Advanced audit configuration
    - OpenSCAP compliance
    - Enterprise authentication (SSSD)
    - Automated patch management
    - System rollback capabilities

## Configuration

Security checks are defined in OS-specific YAML files:
- `cis_checks_rhel8.yaml`
- `cis_checks_rhel9.yaml`
- `cis_checks_ubuntu2204.yaml`
- `cis_checks_ubuntu2404.yaml`
- `cis_checks_sles15.yaml`

Each check has:
- Unique ID
- Description
- Command to execute
- Expected result

### Example Check Definition
```yaml
- id: "1.1.1.1"
  description: "Ensure mounting of cramfs filesystems is disabled"
  command: "modprobe -n -v cramfs"
  expected: "install /bin/true"
```

## HTML Reports

The scanner generates detailed HTML reports including:
- Overall compliance status
- Individual check results
- Color-coded pass/fail indicators
- Timestamp and system information
- Detailed check descriptions

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

### Development Requirements
- Go 1.21 or later
- GoReleaser for building releases
- Access to target OS for testing

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- CIS Benchmarks by Center for Internet Security
- Go community for excellent tooling
- All contributors and users of this project
