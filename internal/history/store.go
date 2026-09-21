package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Store manages loading and atomically persisting session history to local disk.
type Store struct {
	mu      sync.Mutex
	Path    string
	Records []Record
}

// DefaultPath returns the OS-standard configuration path: os.UserConfigDir()/goarcade/history.json.
func DefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	return filepath.Join(configDir, "goarcade", "history.json"), nil
}

// Open loads history from the given path.
// If the file does not exist, an empty store is returned.
// If the file exists but is corrupted, it is backed up to <path>.corrupt-<timestamp>
// and an empty store is initialized so the app can continue.
func Open(path string) (*Store, error) {
	if path == "" {
		var err error
		path, err = DefaultPath()
		if err != nil {
			return nil, err
		}
	}

	store := &Store{
		Path:    path,
		Records: make([]Record, 0),
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	if len(data) == 0 {
		return store, nil
	}

	if err := json.Unmarshal(data, &store.Records); err != nil {
		// File is corrupted: backup and start clean
		timestamp := time.Now().Format("20060102-150405")
		backupPath := fmt.Sprintf("%s.corrupt-%s", path, timestamp)
		if renameErr := os.Rename(path, backupPath); renameErr != nil {
			return nil, fmt.Errorf("history is corrupted and backup failed: %w (original unmarshal error: %v)", renameErr, err)
		}
		store.Records = make([]Record, 0)
	}

	return store, nil
}

// Append adds a record to the store and saves it to disk atomically.
func (s *Store) Append(r Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Records = append(s.Records, r)
	return s.saveLocked()
}

// Save writes all in-memory records to disk using an atomic write (temp file + rename).
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked()
}

func (s *Store) saveLocked() error {
	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(s.Records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode history JSON: %w", err)
	}

	tmpPath := s.Path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write temporary history file: %w", err)
	}

	if err := os.Rename(tmpPath, s.Path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to atomically rename history file: %w", err)
	}

	return nil
}
