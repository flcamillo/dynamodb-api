package main

import (
	"api/repositories"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type MockLogForConfig struct {
	messages []string
}

func (m *MockLogForConfig) Info(format string, a ...any) {
	m.messages = append(m.messages, format)
}

func (m *MockLogForConfig) Error(format string, a ...any) {
	m.messages = append(m.messages, format)
}

func (m *MockLogForConfig) Warn(format string, a ...any) {
	m.messages = append(m.messages, format)
}

func (m *MockLogForConfig) Debug(format string, a ...any) {
	m.messages = append(m.messages, format)
}

func TestNewConfig(t *testing.T) {
	config := NewConfig("test-config.json")

	if config == nil {
		t.Fatal("NewConfig() returned nil")
	}

	if config.File != "test-config.json" {
		t.Errorf("File = %s, want test-config.json", config.File)
	}

	if config.Address != "0.0.0.0" {
		t.Errorf("Address = %s, want 0.0.0.0", config.Address)
	}

	if config.Port != 7000 {
		t.Errorf("Port = %d, want 7000", config.Port)
	}

	if config.RecordTTLMinutes != 24*60 {
		t.Errorf("RecordTTLMinutes = %d, want %d", config.RecordTTLMinutes, 24*60)
	}
}

func TestConfigSave(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")

	config := NewConfig(configFile)
	config.Port = 8080
	config.Address = "127.0.0.1"
	config.RecordTTLMinutes = 48 * 60

	err := config.Save()
	if err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load and verify
	loaded, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v, want nil", err)
	}

	if loaded.Port != 8080 {
		t.Errorf("Loaded Port = %d, want 8080", loaded.Port)
	}

	if loaded.Address != "127.0.0.1" {
		t.Errorf("Loaded Address = %s, want 127.0.0.1", loaded.Address)
	}

	if loaded.RecordTTLMinutes != 48*60 {
		t.Errorf("Loaded RecordTTLMinutes = %d, want %d", loaded.RecordTTLMinutes, 48*60)
	}
}

func TestLoadConfigNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "new-config.json")

	// File does not exist
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		t.Fatal("Config file should not exist before test")
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v, want nil", err)
	}

	if config == nil {
		t.Fatal("LoadConfig() returned nil")
	}

	// Check if file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Config file should have been created")
	}

	// Verify default values
	if config.Port != 7000 {
		t.Errorf("Default Port = %d, want 7000", config.Port)
	}
}

func TestLoadConfigExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "existing-config.json")

	// Create initial config and save
	config := NewConfig(configFile)
	config.Port = 9000
	config.Address = "192.168.1.1"
	config.RecordTTLMinutes = 12 * 60

	err := config.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load it back
	loaded, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v, want nil", err)
	}

	if loaded.Port != 9000 {
		t.Errorf("Port = %d, want 9000", loaded.Port)
	}

	if loaded.Address != "192.168.1.1" {
		t.Errorf("Address = %s, want 192.168.1.1", loaded.Address)
	}

	if loaded.RecordTTLMinutes != 12*60 {
		t.Errorf("RecordTTLMinutes = %d, want %d", loaded.RecordTTLMinutes, 12*60)
	}
}

func TestLoadConfigInvalidPath(t *testing.T) {
	// Try to load from a directory that doesn't exist
	configFile := "/nonexistent/path/config.json"

	_, err := LoadConfig(configFile)
	if err == nil {
		t.Error("LoadConfig() should error for invalid path")
	}
}

func TestConfigFields(t *testing.T) {
	config := NewConfig("test.json")

	// Test that interfaces fields can be set
	config.Log = &MockLogForConfig{}
	config.Repository = repositories.NewMemoryDB(&repositories.MemoryDBConfig{
		Log: config.Log,
		TTL: time.Hour,
	})

	if config.Log == nil {
		t.Error("Log should not be nil")
	}

	if config.Repository == nil {
		t.Error("Repository should not be nil")
	}
}

func TestConfigRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "roundtrip.json")

	// Create original config
	original := NewConfig(configFile)
	original.Port = 5555
	original.Address = "10.0.0.1"
	original.RecordTTLMinutes = 36 * 60

	// Save
	err := original.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load
	loaded, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify all fields match
	if loaded.Port != original.Port {
		t.Errorf("Port mismatch: %d != %d", loaded.Port, original.Port)
	}

	if loaded.Address != original.Address {
		t.Errorf("Address mismatch: %s != %s", loaded.Address, original.Address)
	}

	if loaded.RecordTTLMinutes != original.RecordTTLMinutes {
		t.Errorf("RecordTTLMinutes mismatch: %d != %d", loaded.RecordTTLMinutes, original.RecordTTLMinutes)
	}
}
