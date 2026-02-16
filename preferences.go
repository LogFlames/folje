package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Preferences struct {
	LastConfigPath       string `json:"lastConfigPath"`
	LastIpAddress        string `json:"lastIpAddress"`
	LastVideoSourceId    string `json:"lastVideoSourceId"`
	LastVideoSourceLabel string `json:"lastVideoSourceLabel"`
}

func getPreferencesPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "Folje", "preferences.json"), nil
}

func loadPreferences() Preferences {
	path, err := getPreferencesPath()
	if err != nil {
		LogError("Failed to get preferences path: %s", err.Error())
		return Preferences{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// File doesn't exist or can't be read - return empty preferences
		return Preferences{}
	}

	var prefs Preferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		LogError("Failed to parse preferences: %s", err.Error())
		return Preferences{}
	}

	return prefs
}

func savePreferences(prefs Preferences) error {
	path, err := getPreferencesPath()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (a *App) updateLastConfigPath(path string) {
	prefs := loadPreferences()
	prefs.LastConfigPath = path
	if err := savePreferences(prefs); err != nil {
		LogError("Failed to save preferences: %s", err.Error())
	}
}

func (a *App) updateLastIpAddress(ip string) {
	prefs := loadPreferences()
	prefs.LastIpAddress = ip
	if err := savePreferences(prefs); err != nil {
		LogError("Failed to save preferences: %s", err.Error())
	}
}

func (a *App) updateLastVideoSource(id, label string) {
	prefs := loadPreferences()
	prefs.LastVideoSourceId = id
	prefs.LastVideoSourceLabel = label
	if err := savePreferences(prefs); err != nil {
		LogError("Failed to save preferences: %s", err.Error())
	}
}
