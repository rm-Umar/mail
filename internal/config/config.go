package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	// IMAP
	IMAPServer  string `json:"imap_server"`
	IMAPPort    int    `json:"imap_port"`
	InboxFolder string `json:"inbox_folder"`

	// SMTP
	SMTPServer string `json:"smtp_server"`
	SMTPPort   int    `json:"smtp_port"`

	// Credentials
	Username string `json:"username"`
	Password string `json:"password"`

	// Whether to use TLS for both IMAP & SMTP
	UseTLS bool `json:"use_tls"`
}

// DefaultConfig is used by Login() as the starting values.
var DefaultConfig = Config{
	IMAPServer: "imap.gmail.com",
	IMAPPort:   993,

	SMTPServer: "smtp.gmail.com",
	SMTPPort:   465,

	UseTLS: true,
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config) error {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %v", err)
	}

	appDir := filepath.Join(configDir, ".go-email")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(appDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	return nil
}

// LoadConfig loads the configuration from a file
func LoadConfig() (*Config, error) {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %v", err)
	}

	configFile := filepath.Join(configDir, ".go-email", "config.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return &cfg, nil
}
