package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Status struct {
	LastSpeeds map[string]int `json:"last_speeds"`
}

func LoadStatus() (*Status, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	statusPath := filepath.Join(homeDir, ".config", "openfan", "fan-status")

	data, err := os.ReadFile(statusPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Status{
				LastSpeeds: make(map[string]int),
			}, nil
		}
		return nil, fmt.Errorf("failed to read status file: %w", err)
	}

	var status Status
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("failed to parse status file: %w", err)
	}

	if status.LastSpeeds == nil {
		status.LastSpeeds = make(map[string]int)
	}

	return &status, nil
}

func (s *Status) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "openfan")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	statusPath := filepath.Join(configDir, "fan-status")

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}

	if err := os.WriteFile(statusPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write status file: %w", err)
	}

	return nil
}