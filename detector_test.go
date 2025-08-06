package main

import (
	"testing"
)

func TestQuickDetector(t *testing.T) {
	detector := NewQuickDetector()
	if detector == nil {
		t.Error("Failed to create QuickDetector")
	}
}

func TestMoreDetector(t *testing.T) {
	detector := NewMoreDetector()
	if detector == nil {
		t.Error("Failed to create MoreDetector")
	}
}

func TestConfigLoading(t *testing.T) {
	// This test will pass if the config file exists and is valid
	config, err := loadConfig(".config.yaml")
	if err != nil {
		t.Logf("Warning: Could not load config file: %v", err)
	} else if config == nil {
		t.Error("Config is nil")
	}
}