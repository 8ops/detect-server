package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	configFile = flag.String("config", ".config.yaml", "Configuration file path")
	checkType  = flag.String("c", "quick", "Check type: quick or more")
	outputType = flag.String("s", "stdout", "Output type: stdout, web, email, html, pdf")
)

func main() {
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Printf("Warning: Could not load config file: %v", err)
	}

	// Create detector based on check type
	var detector Detector
	switch *checkType {
	case "quick":
		detector = NewQuickDetector()
	case "more":
		detector = NewMoreDetector()
	default:
		log.Fatalf("Invalid check type: %s. Use 'quick' or 'more'", *checkType)
	}

	// Run detection
	report, err := detector.Detect()
	if err != nil {
		log.Fatalf("Detection failed: %v", err)
	}

	// Generate output based on output type
	switch *outputType {
	case "stdout":
		err = report.ToStdout()
	case "web":
		err = report.ToWeb()
	case "email":
		err = report.ToEmail(config.Email)
	case "html":
		err = report.ToHTML()
	case "pdf":
		err = report.ToPDF()
	default:
		log.Fatalf("Invalid output type: %s. Use 'stdout', 'web', 'email', 'html', or 'pdf'", *outputType)
	}

	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	fmt.Println("Detection completed successfully!")
}