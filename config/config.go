package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ConfigDirectory returns ~/.altbins/.altnotepad
func ConfigDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".altbins", ".altnotepad")
}

// BackupDirectory holds the texts of the unsaved documents of the session
func BackupDirectory() string {
	return filepath.Join(ConfigDirectory(), "backup")
}

// Init loads the settings and the history
func Init() {
	LoadSettings()
	loadHistory()
}

// readJSON reads the file of the config directory into v; false if there is none
func readJSON(name string, v any) bool {
	bs, err := os.ReadFile(filepath.Join(ConfigDirectory(), name))
	if err != nil {
		return false
	}
	return json.Unmarshal(bs, v) == nil
}

// writeJSON writes v to the file of the config directory, replacing it at once
func writeJSON(name string, v any) error {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(ConfigDirectory(), 0755); err != nil {
		return err
	}
	path := filepath.Join(ConfigDirectory(), name)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, bs, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
