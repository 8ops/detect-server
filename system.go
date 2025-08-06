package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// getBasicInfo collects basic server information
func getBasicInfo() (*BasicInfo, error) {
	info := &BasicInfo{}

	// Get hardware info
	hardware, err := getHardwareInfo()
	if err != nil {
		return nil, err
	}
	info.HardwareInfo = hardware

	// Get OS info
	osInfo, err := getOSInfo()
	if err != nil {
		return nil, err
	}
	info.OSInfo = osInfo

	// Get software info
	software, err := getSoftwareInfo()
	if err != nil {
		return nil, err
	}
	info.SoftwareInfo = software

	return info, nil
}

// getHardwareInfo gets hardware information
func getHardwareInfo() (string, error) {
	// Try to get CPU info
	cpuInfo, err := exec.Command("lscpu").Output()
	if err != nil {
		// Fallback for systems without lscpu
		cpuInfo, err = exec.Command("cat", "/proc/cpuinfo").Output()
		if err != nil {
			return "Unknown", nil
		}
	}

	// Try to get memory info
	memInfo, err := exec.Command("free", "-h").Output()
	if err != nil {
		// Fallback for systems without free command
		memInfo, err = exec.Command("cat", "/proc/meminfo").Output()
		if err != nil {
			return "Unknown", nil
		}
	}

	// Try to get disk info
	diskInfo, err := exec.Command("df", "-h").Output()
	if err != nil {
		return "Unknown", nil
	}

	// Combine information
	result := fmt.Sprintf("CPU: %s, Memory: %s, Disk: %s", 
		strings.Split(string(cpuInfo), "\n")[0],
		strings.Split(string(memInfo), "\n")[1],
		strings.Split(string(diskInfo), "\n")[1])

	return result, nil
}

// getOSInfo gets operating system information
func getOSInfo() (string, error) {
	// Try to get OS info from /etc/os-release
	if _, err := os.Stat("/etc/os-release"); err == nil {
		file, err := os.Open("/etc/os-release")
		if err != nil {
			return "Unknown", nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var osName, osVersion string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "NAME=") {
				osName = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
			}
			if strings.HasPrefix(line, "VERSION=") {
				osVersion = strings.Trim(strings.TrimPrefix(line, "VERSION="), "\"")
			}
		}

		if osName != "" && osVersion != "" {
			return fmt.Sprintf("%s %s", osName, osVersion), nil
		}
	}

	// Fallback to uname
	uname, err := exec.Command("uname", "-a").Output()
	if err != nil {
		return "Unknown", nil
	}

	return strings.TrimSpace(string(uname)), nil
}

// getSoftwareInfo gets software information
func getSoftwareInfo() (string, error) {
	// Try to get installed packages (varies by distribution)
	var cmd *exec.Cmd
	if _, err := exec.LookPath("dpkg"); err == nil {
		// Debian/Ubuntu
		cmd = exec.Command("dpkg", "-l")
	} else if _, err := exec.LookPath("rpm"); err == nil {
		// RHEL/CentOS/Fedora
		cmd = exec.Command("rpm", "-qa")
	} else if _, err := exec.LookPath("pacman"); err == nil {
		// Arch Linux
		cmd = exec.Command("pacman", "-Q")
	} else {
		return "Unknown package manager", nil
	}

	output, err := cmd.Output()
	if err != nil {
		return "Error getting software info", nil
	}

	lines := strings.Split(string(output), "\n")
	return fmt.Sprintf("%d packages installed", len(lines)-1), nil
}

// getUSEMetrics collects USE (Utilization/Saturation/Errors) metrics
func getUSEMetrics() (*USEMetrics, error) {
	metrics := &USEMetrics{}

	// Get CPU usage
	cpuUsage, err := getCPUUsage()
	if err != nil {
		return nil, err
	}
	metrics.CPUUsage = cpuUsage

	// Get memory usage
	memoryUsage, err := getMemoryUsage()
	if err != nil {
		return nil, err
	}
	metrics.MemoryUsage = memoryUsage

	// Get disk usage
	diskUsage, err := getDiskUsage()
	if err != nil {
		return nil, err
	}
	metrics.DiskUsage = diskUsage

	// Get network usage
	networkUsage, err := getNetworkUsage()
	if err != nil {
		return nil, err
	}
	metrics.NetworkUsage = networkUsage

	// Get saturation metrics (simplified)
	metrics.CPUSaturation = cpuUsage * 0.8 // Simplified calculation
	metrics.MemorySaturation = memoryUsage * 0.9
	metrics.DiskSaturation = diskUsage * 0.7
	metrics.NetworkSaturation = networkUsage * 0.6

	// Check for abnormal events
	abnormalEvents, err := checkAbnormalEvents()
	if err != nil {
		return nil, err
	}
	metrics.AbnormalEvents = abnormalEvents

	return metrics, nil
}

// getCPUUsage gets CPU usage percentage
func getCPUUsage() (float64, error) {
	// Use top command to get CPU usage
	output, err := exec.Command("top", "-bn1").Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "%Cpu(s)") {
			// Parse the line to extract CPU usage
			fields := strings.Fields(line)
			for i, field := range fields {
				if strings.Contains(field, "id") && i > 0 {
					idle, err := parseFloat(fields[i-1])
					if err != nil {
						return 0, err
					}
					return 100 - idle, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("could not parse CPU usage")
}

// getMemoryUsage gets memory usage percentage
func getMemoryUsage() (float64, error) {
	// Use free command to get memory usage
	output, err := exec.Command("free").Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 3 {
		return 0, fmt.Errorf("unexpected output from free command")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 3 {
		return 0, fmt.Errorf("unexpected output from free command")
	}

	total, err := parseFloat(fields[1])
	if err != nil {
		return 0, err
	}

	used, err := parseFloat(fields[2])
	if err != nil {
		return 0, err
	}

	return (used / total) * 100, nil
}

// getDiskUsage gets disk usage percentage
func getDiskUsage() (float64, error) {
	// Use df command to get disk usage
	output, err := exec.Command("df", "/").Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("unexpected output from df command")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return 0, fmt.Errorf("unexpected output from df command")
	}

	// The 5th field contains usage percentage
	usageStr := strings.TrimSuffix(fields[4], "%")
	usage, err := parseFloat(usageStr)
	if err != nil {
		return 0, err
	}

	return usage, nil
}

// getNetworkUsage gets network usage (simplified)
func getNetworkUsage() (float64, error) {
	// This is a simplified implementation
	// In a real-world scenario, you would monitor network interfaces over time
	return 45.5, nil // Return a placeholder value
}

// checkAbnormalEvents checks for abnormal events like SSH brute force attempts
func checkAbnormalEvents() ([]string, error) {
	var events []string

	// Check auth.log for SSH brute force attempts (Debian/Ubuntu)
	if _, err := os.Stat("/var/log/auth.log"); err == nil {
		cmd := exec.Command("grep", "Failed password", "/var/log/auth.log")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			if len(lines) > 1 {
				events = append(events, fmt.Sprintf("Found %d SSH failed login attempts", len(lines)-1))
			}
		}
	}

	// Check secure log for SSH brute force attempts (RHEL/CentOS)
	if _, err := os.Stat("/var/log/secure"); err == nil {
		cmd := exec.Command("grep", "Failed password", "/var/log/secure")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			if len(lines) > 1 {
				events = append(events, fmt.Sprintf("Found %d SSH failed login attempts", len(lines)-1))
			}
		}
	}

	return events, nil
}

// checkFileIntegrity checks integrity of bin and lib files
func checkFileIntegrity() (*FileIntegrity, error) {
	integrity := &FileIntegrity{}

	// List bin files
	binFiles, err := listFiles("/bin")
	if err != nil {
		return nil, err
	}
	integrity.BinFiles = binFiles

	// List lib files
	libFiles, err := listFiles("/lib")
	if err != nil {
		return nil, err
	}
	integrity.LibFiles = libFiles

	return integrity, nil
}

// listFiles lists files in a directory
func listFiles(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// getRunningAssets gets running assets (ports and processes)
func getRunningAssets() (*RunningAssets, error) {
	assets := &RunningAssets{}

	// Get open ports
	ports, err := getOpenPorts()
	if err != nil {
		return nil, err
	}
	assets.OpenPorts = ports

	// Get processes
	processes, err := getProcesses()
	if err != nil {
		return nil, err
	}
	assets.Processes = processes

	return assets, nil
}

// getOpenPorts gets open ports
func getOpenPorts() ([]string, error) {
	// Use ss command to get open ports
	output, err := exec.Command("ss", "-tuln").Output()
	if err != nil {
		// Fallback to netstat
		output, err = exec.Command("netstat", "-tuln").Output()
		if err != nil {
			return nil, err
		}
	}

	lines := strings.Split(string(output), "\n")
	var ports []string
	for _, line := range lines {
		if strings.Contains(line, "LISTEN") {
			fields := strings.Fields(line)
			if len(fields) > 3 {
				ports = append(ports, fields[3])
			}
		}
	}

	return ports, nil
}

// getProcesses gets running processes
func getProcesses() ([]string, error) {
	// Use ps command to get processes
	output, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var processes []string
	for i, line := range lines {
		// Skip header line
		if i == 0 {
			continue
		}
		if line != "" {
			fields := strings.Fields(line)
			if len(fields) > 10 {
				// Get process name (11th field)
				processes = append(processes, fields[10])
			}
		}
	}

	return processes, nil
}

// checkNetworkStatus checks network connectivity
func checkNetworkStatus() (*NetworkStatus, error) {
	status := &NetworkStatus{
		Samples: make([]NetworkSample, 5), // 5 minutes of samples
	}

	// Check connectivity by pinging a well-known host
	for i := 0; i < 5; i++ {
		sample := NetworkSample{
			Timestamp: time.Now(),
		}

		// Ping Google DNS
		cmd := exec.Command("ping", "-c", "1", "-W", "3", "8.8.8.8")
		err := cmd.Run()
		sample.Status = (err == nil)

		status.Samples[i] = sample
		time.Sleep(60 * time.Second) // Wait 1 minute between samples
	}

	// Calculate overall connectivity
	successCount := 0
	for _, sample := range status.Samples {
		if sample.Status {
			successCount++
		}
	}
	status.Connectivity = (successCount >= 3) // Consider connected if at least 3 out of 5 succeed

	return status, nil
}

// parseFloat is a helper to parse float64 from string
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}