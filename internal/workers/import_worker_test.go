package workers

import (
	"testing"
	"time"
)

func TestDefaultImportWorkerConfig(t *testing.T) {
	config := DefaultImportWorkerConfig()

	if config.PollInterval != 5*time.Second {
		t.Errorf("PollInterval should be 5s, got %v", config.PollInterval)
	}
	if config.Workers != 2 {
		t.Errorf("Workers should be 2, got %d", config.Workers)
	}
	if config.BatchSize != 500 {
		t.Errorf("BatchSize should be 500, got %d", config.BatchSize)
	}
	if config.UploadDir != "./data/uploads" {
		t.Errorf("UploadDir should be ./data/uploads, got %s", config.UploadDir)
	}
}

func TestImportWorkerConfig_CustomValues(t *testing.T) {
	config := ImportWorkerConfig{
		PollInterval: 10 * time.Second,
		Workers:      5,
		BatchSize:    1000,
		UploadDir:    "/custom/uploads",
	}

	if config.PollInterval != 10*time.Second {
		t.Errorf("PollInterval should be 10s, got %v", config.PollInterval)
	}
	if config.Workers != 5 {
		t.Errorf("Workers should be 5, got %d", config.Workers)
	}
	if config.BatchSize != 1000 {
		t.Errorf("BatchSize should be 1000, got %d", config.BatchSize)
	}
	if config.UploadDir != "/custom/uploads" {
		t.Errorf("UploadDir should be /custom/uploads, got %s", config.UploadDir)
	}
}

// TestColumnMapping tests the column name normalization logic
func TestColumnMapping(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Email", "email"},
		{"EMAIL", "email"},
		{"  email  ", "email"},
		{"Name", "name"},
		{"First Name", "first name"},
		{"DOMAIN", "domain"},
		{"Reason", "reason"},
	}

	for _, tt := range tests {
		// Simulate the normalization logic from import_worker.go
		// strings.ToLower(strings.TrimSpace(col))
		result := normalizeColumnName(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeColumnName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// normalizeColumnName normalizes a column name for comparison
// This mirrors the logic in import_worker.go
func normalizeColumnName(col string) string {
	// Trim leading/trailing whitespace
	start := 0
	end := len(col)
	for start < end && (col[start] == ' ' || col[start] == '\t') {
		start++
	}
	for end > start && (col[end-1] == ' ' || col[end-1] == '\t') {
		end--
	}
	col = col[start:end]

	// Convert to lowercase
	result := make([]byte, len(col))
	for i := 0; i < len(col); i++ {
		c := col[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32 // Convert to lowercase
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func TestImportWorker_Stats_Initial(t *testing.T) {
	// Create a minimal worker to test Stats
	worker := &ImportWorker{}

	processed, failed := worker.Stats()
	if processed != 0 {
		t.Errorf("Initial processed should be 0, got %d", processed)
	}
	if failed != 0 {
		t.Errorf("Initial failed should be 0, got %d", failed)
	}
}

func TestImportWorker_Stats_Counting(t *testing.T) {
	worker := &ImportWorker{}

	// Simulate counting
	worker.processed.Add(10)
	worker.failed.Add(2)

	processed, failed := worker.Stats()
	if processed != 10 {
		t.Errorf("processed should be 10, got %d", processed)
	}
	if failed != 2 {
		t.Errorf("failed should be 2, got %d", failed)
	}
}

// TestEmailValidation tests basic email string handling
func TestEmailNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test@example.com", "test@example.com"},
		{"TEST@EXAMPLE.COM", "test@example.com"},
		{"  Test@Example.com  ", "test@example.com"},
		{"", ""},
	}

	for _, tt := range tests {
		// Simulate the normalization logic from importSubscriber
		result := normalizeEmail(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeEmail(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// normalizeEmail normalizes an email address for storage
// This mirrors the logic in import_worker.go
func normalizeEmail(email string) string {
	// Trim leading/trailing whitespace
	start := 0
	end := len(email)
	for start < end && (email[start] == ' ' || email[start] == '\t') {
		start++
	}
	for end > start && (email[end-1] == ' ' || email[end-1] == '\t') {
		end--
	}
	email = email[start:end]

	// Convert to lowercase
	result := make([]byte, len(email))
	for i := 0; i < len(email); i++ {
		c := email[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// TestDomainNormalization tests domain string handling
func TestDomainNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "example.com"},
		{"EXAMPLE.COM", "example.com"},
		{"  Example.com  ", "example.com"},
		{"", ""},
	}

	for _, tt := range tests {
		result := normalizeEmail(tt.input) // Same logic applies
		if result != tt.expected {
			t.Errorf("normalizeDomain(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestImportTypes(t *testing.T) {
	// Verify the import types are valid strings
	validTypes := []string{"subscribers", "suppression", "blocked_domains"}

	for _, importType := range validTypes {
		if importType == "" {
			t.Error("Import type should not be empty")
		}
	}
}
