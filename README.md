# detect-server

A comprehensive server detection tool for Linux systems written in Go.

## Features

- Quick detection mode for basic server health checks
- More detection mode for comprehensive server analysis
- Multiple output formats: stdout, web, email, HTML, PDF
- Configurable email reporting

## Installation

1. Make sure you have Go installed (version 1.21 or later)
2. Clone this repository
3. Run `go mod tidy` to download dependencies
4. Build the tool with `go build`

## Usage

```bash
# Quick detection (default)
./detect-server

# More comprehensive detection
./detect-server -c more

# Specify output format
./detect-server -s html

# Combine options
./detect-server -c more -s pdf
```

## Configuration

The tool reads configuration from `.config.yaml` in the current directory:

```yaml
email:
  nickname: "检测报告"
  username: "detect@8ops.top"
  passport: "xx"
  host: "smtp.exmail.qq.com"
  port: 465
  to: ["admin@8ops.top"]
  cc: []
  attachment: []
```

## Detection Items

### Quick Detection

1. Basic Information
   - Hardware information (CPU, memory, disk, network, power)
   - Operating system information
   - Software information (installed packages and running status)

2. USE Metrics
   - Resource utilization (CPU, memory, disk, network)
   - Resource saturation
   - Abnormal events (SSH brute force attempts)

3. File Integrity
   - Binary files
   - Library files

### More Detection

Includes all quick detection items plus:

4. Running Assets
   - Open ports and associated services
   - Process list

5. Network Status
   - Network connectivity
   - 5-minute sampling