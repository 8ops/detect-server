package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"gopkg.in/gomail.v2"
)

// Report represents the detection report
type Report struct {
	Type           string
	Timestamp      time.Time
	BasicInfo      *BasicInfo
	USEMetrics     *USEMetrics
	FileIntegrity  *FileIntegrity
	RunningAssets  *RunningAssets
	NetworkStatus  *NetworkStatus
}

// BasicInfo represents basic server information
type BasicInfo struct {
	HardwareInfo string
	OSInfo       string
	SoftwareInfo string
}

// USEMetrics represents Utilization/Saturation/Errors metrics
type USEMetrics struct {
	CPUUsage      float64
	MemoryUsage   float64
	DiskUsage     float64
	NetworkUsage  float64
	CPUSaturation float64
	MemorySaturation float64
	DiskSaturation   float64
	NetworkSaturation float64
	AbnormalEvents   []string
}

// FileIntegrity represents file integrity check results
type FileIntegrity struct {
	BinFiles []string
	LibFiles []string
}

// RunningAssets represents running assets like ports and processes
type RunningAssets struct {
	OpenPorts []string
	Processes []string
}

// NetworkStatus represents network connectivity status
type NetworkStatus struct {
	Connectivity bool
	Samples      []NetworkSample
}

// NetworkSample represents a network sample
type NetworkSample struct {
	Timestamp time.Time
	Status    bool
}

// ToStdout outputs the report to stdout
func (r *Report) ToStdout() error {
	fmt.Println("=== Server Detection Report ===")
	fmt.Printf("Type: %s\n", r.Type)
	fmt.Printf("Timestamp: %s\n", r.Timestamp.Format("2006-01-02 15:04:05"))
	
	fmt.Println("\n--- Basic Information ---")
	if r.BasicInfo != nil {
		fmt.Printf("Hardware: %s\n", r.BasicInfo.HardwareInfo)
		fmt.Printf("OS: %s\n", r.BasicInfo.OSInfo)
		fmt.Printf("Software: %s\n", r.BasicInfo.SoftwareInfo)
	}
	
	fmt.Println("\n--- USE Metrics ---")
	if r.USEMetrics != nil {
		fmt.Printf("CPU Usage: %.2f%%\n", r.USEMetrics.CPUUsage)
		fmt.Printf("Memory Usage: %.2f%%\n", r.USEMetrics.MemoryUsage)
		fmt.Printf("Disk Usage: %.2f%%\n", r.USEMetrics.DiskUsage)
		fmt.Printf("Network Usage: %.2f%%\n", r.USEMetrics.NetworkUsage)
		fmt.Printf("CPU Saturation: %.2f%%\n", r.USEMetrics.CPUSaturation)
		fmt.Printf("Memory Saturation: %.2f%%\n", r.USEMetrics.MemorySaturation)
		fmt.Printf("Disk Saturation: %.2f%%\n", r.USEMetrics.DiskSaturation)
		fmt.Printf("Network Saturation: %.2f%%\n", r.USEMetrics.NetworkSaturation)
		
		if len(r.USEMetrics.AbnormalEvents) > 0 {
			fmt.Println("Abnormal Events:")
			for _, event := range r.USEMetrics.AbnormalEvents {
				fmt.Printf("  - %s\n", event)
			}
		}
	}
	
	fmt.Println("\n--- File Integrity ---")
	if r.FileIntegrity != nil {
		fmt.Printf("Bin Files: %s\n", strings.Join(r.FileIntegrity.BinFiles, ", "))
		fmt.Printf("Lib Files: %s\n", strings.Join(r.FileIntegrity.LibFiles, ", "))
	}
	
	if r.Type == "More Detection" {
		fmt.Println("\n--- Running Assets ---")
		if r.RunningAssets != nil {
			fmt.Printf("Open Ports: %s\n", strings.Join(r.RunningAssets.OpenPorts, ", "))
			fmt.Printf("Processes: %s\n", strings.Join(r.RunningAssets.Processes, ", "))
		}
		
		fmt.Println("\n--- Network Status ---")
		if r.NetworkStatus != nil {
			fmt.Printf("Connectivity: %t\n", r.NetworkStatus.Connectivity)
			fmt.Println("Samples:")
			for _, sample := range r.NetworkStatus.Samples {
				status := "OK"
				if !sample.Status {
					status = "FAIL"
				}
				fmt.Printf("  %s: %s\n", sample.Timestamp.Format("15:04:05"), status)
			}
		}
	}
	
	return nil
}

// ToWeb outputs the report to a web interface (placeholder)
func (r *Report) ToWeb() error {
	return fmt.Errorf("web output not implemented yet")
}

// ToEmail sends the report via email
func (r *Report) ToEmail(config EmailConfig) error {
	// Create message
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Username)
	msg.SetHeader("To", config.To...)
	if len(config.CC) > 0 {
		msg.SetHeader("Cc", config.CC...)
	}
	msg.SetHeader("Subject", fmt.Sprintf("Server Detection Report - %s", r.Timestamp.Format("2006-01-02")))
	
	// Build body
	body := fmt.Sprintf("Server Detection Report\n\nType: %s\nTimestamp: %s\n", 
		r.Type, r.Timestamp.Format("2006-01-02 15:04:05"))
	
	// Add basic info
	if r.BasicInfo != nil {
		body += fmt.Sprintf("\nBasic Information:\nHardware: %s\nOS: %s\nSoftware: %s\n", 
			r.BasicInfo.HardwareInfo, r.BasicInfo.OSInfo, r.BasicInfo.SoftwareInfo)
	}
	
	// Add USE metrics
	if r.USEMetrics != nil {
		body += fmt.Sprintf("\nUSE Metrics:\nCPU Usage: %.2f%%\nMemory Usage: %.2f%%\nDisk Usage: %.2f%%\nNetwork Usage: %.2f%%\n", 
			r.USEMetrics.CPUUsage, r.USEMetrics.MemoryUsage, r.USEMetrics.DiskUsage, r.USEMetrics.NetworkUsage)
		
		if len(r.USEMetrics.AbnormalEvents) > 0 {
			body += "Abnormal Events:\n"
			for _, event := range r.USEMetrics.AbnormalEvents {
				body += fmt.Sprintf("  - %s\n", event)
			}
		}
	}
	
	// Add file integrity
	if r.FileIntegrity != nil {
		body += fmt.Sprintf("\nFile Integrity:\nBin Files: %s\nLib Files: %s\n", 
			strings.Join(r.FileIntegrity.BinFiles, ", "), strings.Join(r.FileIntegrity.LibFiles, ", "))
	}
	
	// Add more detection items if applicable
	if r.Type == "More Detection" {
		// Add running assets
		if r.RunningAssets != nil {
			body += fmt.Sprintf("\nRunning Assets:\nOpen Ports: %s\nProcesses: %s\n", 
				strings.Join(r.RunningAssets.OpenPorts, ", "), strings.Join(r.RunningAssets.Processes, ", "))
		}
		
		// Add network status
		if r.NetworkStatus != nil {
			body += fmt.Sprintf("\nNetwork Status:\nConnectivity: %t\n", r.NetworkStatus.Connectivity)
		}
	}
	
	msg.SetBody("text/plain", body)
	
	// Create dialer
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Passport)
	
	// Send email
	return dialer.DialAndSend(msg)
}

// ToHTML generates an HTML report
func (r *Report) ToHTML() error {
	// Define HTML template
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
	<title>Server Detection Report</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 20px; }
		h1, h2, h3 { color: #333; }
		.section { margin-bottom: 20px; }
		.metric { margin: 5px 0; }
		.progress-bar {
			width: 100%;
			background-color: #f0f0f0;
			border-radius: 5px;
			margin: 5px 0;
		}
		.progress {
			height: 20px;
			border-radius: 5px;
			text-align: center;
			line-height: 20px;
			color: white;
		}
		.cpu { background-color: #ff6b6b; width: {{.USEMetrics.CPUUsage}}%; }
		.memory { background-color: #4ecdc4; width: {{.USEMetrics.MemoryUsage}}%; }
		.disk { background-color: #45b7d1; width: {{.USEMetrics.DiskUsage}}%; }
		.network { background-color: #96ceb4; width: {{.USEMetrics.NetworkUsage}}%; }
	</style>
</head>
<body>
	<h1>Server Detection Report</h1>
	<div class="section">
		<h2>Report Information</h2>
		<p><strong>Type:</strong> {{.Type}}</p>
		<p><strong>Timestamp:</strong> {{.Timestamp.Format "2006-01-02 15:04:05"}}</p>
	</div>
	
	<div class="section">
		<h2>Basic Information</h2>
		<p><strong>Hardware:</strong> {{.BasicInfo.HardwareInfo}}</p>
		<p><strong>OS:</strong> {{.BasicInfo.OSInfo}}</p>
		<p><strong>Software:</strong> {{.BasicInfo.SoftwareInfo}}</p>
	</div>
	
	<div class="section">
		<h2>USE Metrics</h2>
		<div class="metric">
			<strong>CPU Usage: {{printf "%.2f" .USEMetrics.CPUUsage}}%</strong>
			<div class="progress-bar">
				<div class="progress cpu">{{printf "%.2f" .USEMetrics.CPUUsage}}%</div>
			</div>
		</div>
		<div class="metric">
			<strong>Memory Usage: {{printf "%.2f" .USEMetrics.MemoryUsage}}%</strong>
			<div class="progress-bar">
				<div class="progress memory">{{printf "%.2f" .USEMetrics.MemoryUsage}}%</div>
			</div>
		</div>
		<div class="metric">
			<strong>Disk Usage: {{printf "%.2f" .USEMetrics.DiskUsage}}%</strong>
			<div class="progress-bar">
				<div class="progress disk">{{printf "%.2f" .USEMetrics.DiskUsage}}%</div>
			</div>
		</div>
		<div class="metric">
			<strong>Network Usage: {{printf "%.2f" .USEMetrics.NetworkUsage}}%</strong>
			<div class="progress-bar">
				<div class="progress network">{{printf "%.2f" .USEMetrics.NetworkUsage}}%</div>
			</div>
		</div>
		
		{{if .USEMetrics.AbnormalEvents}}
		<h3>Abnormal Events</h3>
		<ul>
			{{range .USEMetrics.AbnormalEvents}}
			<li>{{.}}</li>
			{{end}}
		</ul>
		{{end}}
	</div>
	
	<div class="section">
		<h2>File Integrity</h2>
		<p><strong>Bin Files:</strong> {{range .FileIntegrity.BinFiles}}{{.}} {{end}}</p>
		<p><strong>Lib Files:</strong> {{range .FileIntegrity.LibFiles}}{{.}} {{end}}</p>
	</div>
	
	{{if eq .Type "More Detection"}}
	<div class="section">
		<h2>Running Assets</h2>
		<p><strong>Open Ports:</strong> {{range .RunningAssets.OpenPorts}}{{.}} {{end}}</p>
		<p><strong>Processes:</strong> {{range .RunningAssets.Processes}}{{.}} {{end}}</p>
	</div>
	
	<div class="section">
		<h2>Network Status</h2>
		<p><strong>Connectivity:</strong> {{.NetworkStatus.Connectivity}}</p>
		<h3>Samples</h3>
		<ul>
			{{range .NetworkStatus.Samples}}
			<li>{{.Timestamp.Format "15:04:05"}}: {{if .Status}}OK{{else}}FAIL{{end}}</li>
			{{end}}
		</ul>
	</div>
	{{end}}
</body>
</html>
`

	// Parse template
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	// Create output file
	file, err := os.Create("report.html")
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute template
	return tmpl.Execute(file, r)
}

// ToPDF generates a PDF report
func (r *Report) ToPDF() error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(40, 10, "Server Detection Report")
	pdf.Ln(12)

	// Report info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Type: %s", r.Type))
	pdf.Ln(4)
	pdf.Cell(40, 10, fmt.Sprintf("Timestamp: %s", r.Timestamp.Format("2006-01-02 15:04:05")))
	pdf.Ln(8)

	// Basic info section
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Basic Information")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 12)
	if r.BasicInfo != nil {
		pdf.Cell(40, 10, fmt.Sprintf("Hardware: %s", r.BasicInfo.HardwareInfo))
		pdf.Ln(4)
		pdf.Cell(40, 10, fmt.Sprintf("OS: %s", r.BasicInfo.OSInfo))
		pdf.Ln(4)
		pdf.Cell(40, 10, fmt.Sprintf("Software: %s", r.BasicInfo.SoftwareInfo))
		pdf.Ln(8)
	}

	// USE metrics section
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "USE Metrics")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 12)
	if r.USEMetrics != nil {
		pdf.Cell(40, 10, fmt.Sprintf("CPU Usage: %.2f%%", r.USEMetrics.CPUUsage))
		pdf.Ln(4)
		pdf.Cell(40, 10, fmt.Sprintf("Memory Usage: %.2f%%", r.USEMetrics.MemoryUsage))
		pdf.Ln(4)
		pdf.Cell(40, 10, fmt.Sprintf("Disk Usage: %.2f%%", r.USEMetrics.DiskUsage))
		pdf.Ln(4)
		pdf.Cell(40, 10, fmt.Sprintf("Network Usage: %.2f%%", r.USEMetrics.NetworkUsage))
		pdf.Ln(8)
	}

	// File integrity section
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "File Integrity")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 12)
	if r.FileIntegrity != nil {
		pdf.Cell(40, 10, fmt.Sprintf("Bin Files: %s", strings.Join(r.FileIntegrity.BinFiles, ", ")))
		pdf.Ln(4)
		pdf.Cell(40, 10, fmt.Sprintf("Lib Files: %s", strings.Join(r.FileIntegrity.LibFiles, ", ")))
		pdf.Ln(8)
	}

	// Additional sections for more detection
	if r.Type == "More Detection" {
		// Running assets section
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Running Assets")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 12)
		if r.RunningAssets != nil {
			pdf.Cell(40, 10, fmt.Sprintf("Open Ports: %s", strings.Join(r.RunningAssets.OpenPorts, ", ")))
			pdf.Ln(4)
			pdf.Cell(40, 10, fmt.Sprintf("Processes: %s", strings.Join(r.RunningAssets.Processes, ", ")))
			pdf.Ln(8)
		}

		// Network status section
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Network Status")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 12)
		if r.NetworkStatus != nil {
			pdf.Cell(40, 10, fmt.Sprintf("Connectivity: %t", r.NetworkStatus.Connectivity))
			pdf.Ln(8)
		}
	}

	// Save to file
	return pdf.OutputFileAndClose("report.pdf")
}