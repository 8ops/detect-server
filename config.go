package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// EmailConfig represents the email configuration
type EmailConfig struct {
	Nickname   string   `yaml:"nickname"`
	Username   string   `yaml:"username"`
	Passport   string   `yaml:"passport"`
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	To         []string `yaml:"to"`
	CC         []string `yaml:"cc"`
	Attachment []string `yaml:"attachment"`
}

// Config represents the application configuration
type Config struct {
	Email EmailConfig `yaml:"email"`
}

// loadConfig loads the configuration from a YAML file
func loadConfig(filename string) (*Config, error) {
	config := &Config{}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return config, fmt.Errorf("config file does not exist: %s", filename)
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}