package history

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResetAllData removes all saved session history and achievements from disk.
// If configDir is empty, the default os.UserConfigDir()/goarcade directory is targeted.
func ResetAllData(configDir string) error {
	if configDir == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to locate config directory: %w", err)
		}
		configDir = filepath.Join(base, "goarcade")
	}

	// Remove known data files
	targets := []string{
		filepath.Join(configDir, "history.json"),
		filepath.Join(configDir, "history.json.tmp"),
		filepath.Join(configDir, "achievements.json"),
		filepath.Join(configDir, "achievements.json.tmp"),
	}

	for _, target := range targets {
		if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", target, err)
		}
	}

	// Also clean up any .corrupt- backup files
	entries, err := os.ReadDir(configDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".tmp" || len(entry.Name()) > 15) {
				_ = os.Remove(filepath.Join(configDir, entry.Name()))
			}
		}
	}

	return nil
}
