package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type FanStatusResponse struct {
	Status     string `json:"status"`
	RPM        int    `json:"rpm"`
	PWMPercent int    `json:"pwm_percent"`
}

func getFanSpeed(fanURL string) (int, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fanURL + "/api/v0/fan/status")
	if err != nil {
		return 0, fmt.Errorf("failed to get fan status: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var fanStatus FanStatusResponse
	if err := json.Unmarshal(body, &fanStatus); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if fanStatus.Status != "ok" {
		return 0, fmt.Errorf("API error: status not ok")
	}

	return fanStatus.PWMPercent, nil
}

func setFanSpeed(fanURL string, speed int) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("%s/api/v0/fan/0/set?value=%d", fanURL, speed)
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to set fan speed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if apiResp.Status != "success" {
		return fmt.Errorf("API error: %s", apiResp.Message)
	}

	return nil
}

func getFanRPM(fanURL string) (int, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fanURL + "/api/v0/fan/status")
	if err != nil {
		return 0, fmt.Errorf("failed to get fan status: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var fanStatus FanStatusResponse
	if err := json.Unmarshal(body, &fanStatus); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if fanStatus.Status != "ok" {
		return 0, fmt.Errorf("API error: status not ok")
	}

	return fanStatus.RPM, nil
}