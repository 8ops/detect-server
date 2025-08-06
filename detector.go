package main

// Detector defines the interface for server detection
type Detector interface {
	Detect() (*Report, error)
}

// QuickDetector implements quick detection
type QuickDetector struct{}

// NewQuickDetector creates a new QuickDetector
func NewQuickDetector() *QuickDetector {
	return &QuickDetector{}
}

// Detect performs quick detection
func (qd *QuickDetector) Detect() (*Report, error) {
	report := &Report{
		Type: "Quick Detection",
	}

	// Basic information
	basicInfo, err := getBasicInfo()
	if err != nil {
		return nil, err
	}
	report.BasicInfo = basicInfo

	// USE metrics
	useMetrics, err := getUSEMetrics()
	if err != nil {
		return nil, err
	}
	report.USEMetrics = useMetrics

	// File integrity
	fileIntegrity, err := checkFileIntegrity()
	if err != nil {
		return nil, err
	}
	report.FileIntegrity = fileIntegrity

	return report, nil
}

// MoreDetector implements more comprehensive detection
type MoreDetector struct{}

// NewMoreDetector creates a new MoreDetector
func NewMoreDetector() *MoreDetector {
	return &MoreDetector{}
}

// Detect performs more comprehensive detection
func (md *MoreDetector) Detect() (*Report, error) {
	report := &Report{
		Type: "More Detection",
	}

	// Basic information
	basicInfo, err := getBasicInfo()
	if err != nil {
		return nil, err
	}
	report.BasicInfo = basicInfo

	// USE metrics
	useMetrics, err := getUSEMetrics()
	if err != nil {
		return nil, err
	}
	report.USEMetrics = useMetrics

	// File integrity
	fileIntegrity, err := checkFileIntegrity()
	if err != nil {
		return nil, err
	}
	report.FileIntegrity = fileIntegrity

	// Running assets (ports and processes)
	runningAssets, err := getRunningAssets()
	if err != nil {
		return nil, err
	}
	report.RunningAssets = runningAssets

	// Network status
	networkStatus, err := checkNetworkStatus()
	if err != nil {
		return nil, err
	}
	report.NetworkStatus = networkStatus

	return report, nil
}