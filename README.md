# Network Scanner

A high-performance network scanning and monitoring tool built in Go, featuring concurrent host discovery, port scanning, MAC vendor lookup, and a REST API for integration with NOC dashboards and monitoring systems.

## Features

### Core Capabilities
- **Concurrent Network Scanning**: Leverages Go goroutines for parallel host discovery across entire subnets
- **Host Discovery**: Ping sweep with automatic ARP cache analysis for MAC address retrieval
- **Port Scanning**: TCP connect scanning of common ports (SSH, HTTP, HTTPS, SMB, RDP, etc.)
- **Vendor Identification**: MAC address vendor lookup via IEEE OUI database
- **Change Detection**: Compare scan results to identify new/removed hosts
- **Persistent Storage**: JSON export/import for historical tracking
- **REST API**: HTTP endpoints for triggering scans and retrieving inventory

### Technical Highlights
- Written in Go 1.25+ for maximum performance and concurrency
- Modular package architecture for maintainability
- Sub-minute scan times for /24 networks
- Zero external dependencies for core scanning (uses system tools)
- Systemd integration for continuous monitoring

## Architecture
network-scanner/

├── main.go # CLI scanner entry point

├── server.go # REST API server entry point

├── pkg/

│ ├── models/ # Data structures (Host)

│ ├── scanner/ # Core scanning logic

│ ├── lookup/ # MAC vendor lookup

│ ├── storage/ # JSON persistence & comparison

│ └── api/ # REST API handlers

└── bin/ # Compiled binaries


## Installation

### Prerequisites
- Go 1.25 or higher
- Linux system (tested on Arch Linux)
- Standard network utilities: `ping`, `ip`

### Build from Source

Clone repository

git clone https://github.com/kragoolk/network-scanner.git
cd network-scanner

Build binaries

go build -o bin/scanner main.go
go build -o bin/scanner-api server.go
chmod +x bin/scanner bin/scanner-api

## Usage

### CLI Scanner

**Basic scan:**

./bin/scanner

**Save results:**

./bin/scanner -save baseline.json

**Compare scans (detect new devices):**

./bin/scanner -save latest.json -compare baseline.json


**Example Output:**

=== Network Scanner v0.7 ===
Scanning network: 192.168.1.0/24
=== Found 8 active hosts ===
IP Address MAC Address Vendor Open Ports

192.168.1.1 a0:91:ca:04:80:72 Nokia Solutions and Net 80,443
192.168.1.50 38:8b:59:09:cd:de Google, Inc. 8443
192.168.1.100 22:31:ca:8d:7f:c9 Unknown Vendor 445

Scan results saved to: baseline.json

### REST API Server

**Start server:**

./bin/scanner-api -network 192.168.1.0/24 -port 8080

**API Endpoints:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/`      | API information |
| POST   | `/scan`  | Trigger new network scan |
| GET    | `/hosts` | Retrieve scan results (JSON) |
| GET    | `/status`| Check scan status |

**Example API Usage:**

Trigger scan

curl -X POST http://localhost:8080/scan
Get results

curl http://localhost:8080/hosts | jq
Check status

curl http://localhost:8080/status

**Sample JSON Response:**

{
	"timestamp": "2025-11-12T12:32:13.876279857-06:00",
	"network": "192.168.1.0/24",
	"hosts": [
	    {
	    "IP": "192.168.1.1",
	    "MAC": "a0:91:ca:04:80:72",
	    "Vendor": "Nokia Solutions and Networks",
	"OpenPorts":
	    }
	]
}

## Systemd Service (Production Deployment)

**Install as system service:**

1. Create service file at `/etc/systemd/system/network-scanner.service`:

[Unit]
Description=Network Scanner API Service
After=network-online.target

[Service]
Type=simple
User=youruser
WorkingDirectory=/path/to/network-scanner
ExecStart=/path/to/network-scanner/bin/scanner-api -network 192.168.1.0/24 -port 8080
Restart=on-failure

[Install]
WantedBy=multi-user.target

2. Enable and start:

sudo systemctl daemon-reload
sudo systemctl enable network-scanner.service
sudo systemctl start network-scanner.service

3. View logs:

sudo journalctl -u network-scanner -f

## Performance

- **Scan Speed**: ~60 seconds for /24 network (256 addresses)
- **Concurrency**: Parallel goroutines for host and port scanning
- **Memory Usage**: ~28MB runtime footprint
- **Binary Size**: 8-9MB standalone executables

## Security Considerations

- Uses TCP connect scanning (no raw sockets/root required for basic scans)
- ARP cache reading requires standard user permissions
- API server has no authentication (deploy behind firewall/VPN)
- Systemd service runs with `NoNewPrivileges=true`

## Use Cases

- **Network Administration**: Quick inventory of active devices
- **Security Monitoring**: Detect unauthorized devices joining network
- **NOC Integration**: REST API feeds into dashboards (Grafana, etc.)
- **Incident Response**: Historical comparison of network state
- **IoT Management**: Track smart home device connectivity

## Future Enhancements

- [ ] SNMP device interrogation
- [ ] Webhook/email alerting for new devices
- [ ] SQLite database for long-term history
- [ ] Web dashboard UI (React frontend)
- [ ] Prometheus metrics export
- [ ] Docker containerization
- [ ] IPv6 support
- [ ] Service fingerprinting (banner grabbing)

## License

MIT License - See LICENSE file for details

## Author

Oliver Krauss - [GitHub](https://github.com/kragoolk) | [LinkedIn](https://linkedin.com/in/oliverkrauss)

## Acknowledgments

- MAC vendor data: macvendors.com API
- Inspired by nmap, arp-scan, and Enterprise network monitoring tools
