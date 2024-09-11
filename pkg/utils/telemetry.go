package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var telemetryFilePath = filepath.Join(os.Getenv("HOME"), ".pizza-cli", "telemetry.json")

// userTelemetryConfig is the config for the user's anonymous telemetry data
type userTelemetryConfig struct {
	ID string `json:"id"`
}

// getOrCreateUniqueID reads the telemetry.json file to fetch the user's anonymous, unique ID.
// In case of error (i.e., if the file doesn't exist or is invalid) it generates
// a new UUID and stores it in the telemetry.json file
func getOrCreateUniqueID() (string, error) {
	if _, err := os.Stat(telemetryFilePath); os.IsNotExist(err) {
		return createTelemetryUUID()
	}

	data, err := os.ReadFile(telemetryFilePath)
	if err != nil {
		return createTelemetryUUID()
	}

	// Try parsing the telemetry file
	var teleData userTelemetryConfig
	err = json.Unmarshal(data, &teleData)
	if err != nil || teleData.ID == "" {
		return createTelemetryUUID()
	}

	return teleData.ID, nil
}

// createTelemetryUUID generates a new UUID and writes it to the user's telemetry.json file
func createTelemetryUUID() (string, error) {
	newUUID := uuid.New().String()

	teleData := userTelemetryConfig{
		ID: newUUID,
	}

	data, err := json.Marshal(teleData)
	if err != nil {
		return "", fmt.Errorf("error creating telemetry data: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(telemetryFilePath), 0755)
	if err != nil {
		return "", fmt.Errorf("error creating directory for telemetry file: %w", err)
	}

	err = os.WriteFile(telemetryFilePath, data, 0600)
	if err != nil {
		return "", fmt.Errorf("error writing telemetry file: %w", err)
	}

	return newUUID, nil
}
