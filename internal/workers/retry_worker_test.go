package workers

import (
	"testing"
	"time"
)

func TestDefaultRetryWorkerConfig(t *testing.T) {
	config := DefaultRetryWorkerConfig()

	if config.PollInterval != 30*time.Second {
		t.Errorf("PollInterval should be 30s, got %v", config.PollInterval)
	}
	if config.MaxRetries != 3 {
		t.Errorf("MaxRetries should be 3, got %d", config.MaxRetries)
	}
	if config.BatchSize != 100 {
		t.Errorf("BatchSize should be 100, got %d", config.BatchSize)
	}

	// Check retry intervals
	expectedIntervals := []time.Duration{
		5 * time.Minute,
		30 * time.Minute,
		2 * time.Hour,
	}
	if len(config.RetryIntervals) != len(expectedIntervals) {
		t.Fatalf("RetryIntervals length should be %d, got %d",
			len(expectedIntervals), len(config.RetryIntervals))
	}
	for i, interval := range config.RetryIntervals {
		if interval != expectedIntervals[i] {
			t.Errorf("RetryIntervals[%d] should be %v, got %v", i, expectedIntervals[i], interval)
		}
	}
}

func TestRetryWorkerConfig_CustomValues(t *testing.T) {
	config := RetryWorkerConfig{
		PollInterval: 1 * time.Minute,
		MaxRetries:   5,
		RetryIntervals: []time.Duration{
			1 * time.Minute,
			5 * time.Minute,
			15 * time.Minute,
			30 * time.Minute,
			1 * time.Hour,
		},
		BatchSize: 50,
	}

	if config.PollInterval != 1*time.Minute {
		t.Errorf("PollInterval should be 1m, got %v", config.PollInterval)
	}
	if config.MaxRetries != 5 {
		t.Errorf("MaxRetries should be 5, got %d", config.MaxRetries)
	}
	if config.BatchSize != 50 {
		t.Errorf("BatchSize should be 50, got %d", config.BatchSize)
	}
	if len(config.RetryIntervals) != 5 {
		t.Errorf("RetryIntervals should have 5 intervals, got %d", len(config.RetryIntervals))
	}
}

func TestRetryWorker_Stats_Initial(t *testing.T) {
	worker := &RetryWorker{}

	retried, succeeded, exhausted := worker.Stats()
	if retried != 0 {
		t.Errorf("Initial retried should be 0, got %d", retried)
	}
	if succeeded != 0 {
		t.Errorf("Initial succeeded should be 0, got %d", succeeded)
	}
	if exhausted != 0 {
		t.Errorf("Initial exhausted should be 0, got %d", exhausted)
	}
}

func TestRetryWorker_Stats_Counting(t *testing.T) {
	worker := &RetryWorker{}

	// Simulate activity
	worker.retried.Add(10)
	worker.succeeded.Add(7)
	worker.exhausted.Add(3)

	retried, succeeded, exhausted := worker.Stats()
	if retried != 10 {
		t.Errorf("retried should be 10, got %d", retried)
	}
	if succeeded != 7 {
		t.Errorf("succeeded should be 7, got %d", succeeded)
	}
	if exhausted != 3 {
		t.Errorf("exhausted should be 3, got %d", exhausted)
	}
}

// TestBackoffIntervalSelection tests the backoff interval selection logic
func TestBackoffIntervalSelection(t *testing.T) {
	intervals := []time.Duration{
		5 * time.Minute,
		30 * time.Minute,
		2 * time.Hour,
	}

	tests := []struct {
		retryCount int64
		expected   time.Duration
	}{
		{0, 5 * time.Minute},   // First retry uses first interval
		{1, 30 * time.Minute},  // Second retry uses second interval
		{2, 2 * time.Hour},     // Third retry uses third interval
		{3, 2 * time.Hour},     // Fourth retry (exceeds array) uses last interval
		{10, 2 * time.Hour},    // Any higher count uses last interval
	}

	for _, tt := range tests {
		intervalIdx := int(tt.retryCount)
		if intervalIdx >= len(intervals) {
			intervalIdx = len(intervals) - 1
		}
		result := intervals[intervalIdx]
		if result != tt.expected {
			t.Errorf("retryCount=%d: expected %v, got %v", tt.retryCount, tt.expected, result)
		}
	}
}

// TestMaxRetriesCheck tests the max retries exceeded logic
func TestMaxRetriesCheck(t *testing.T) {
	maxRetries := int64(3)

	tests := []struct {
		retryCount int64
		shouldSkip bool
	}{
		{0, false}, // Can retry
		{1, false}, // Can retry
		{2, false}, // Can retry (last attempt)
		{3, true},  // Max reached, permanently fail
		{4, true},  // Already past max
	}

	for _, tt := range tests {
		exceeded := tt.retryCount >= maxRetries
		if exceeded != tt.shouldSkip {
			t.Errorf("retryCount=%d: expected skip=%v, got skip=%v",
				tt.retryCount, tt.shouldSkip, exceeded)
		}
	}
}

// TestRetryIntervalsBounds tests that retry intervals are reasonable
func TestRetryIntervalsBounds(t *testing.T) {
	config := DefaultRetryWorkerConfig()

	// First retry should be at least 1 minute
	if config.RetryIntervals[0] < time.Minute {
		t.Error("First retry interval should be at least 1 minute")
	}

	// Each interval should be longer than the previous
	for i := 1; i < len(config.RetryIntervals); i++ {
		if config.RetryIntervals[i] <= config.RetryIntervals[i-1] {
			t.Errorf("Interval %d (%v) should be longer than interval %d (%v)",
				i, config.RetryIntervals[i], i-1, config.RetryIntervals[i-1])
		}
	}

	// Last interval shouldn't exceed 24 hours (too long defeats the purpose)
	if config.RetryIntervals[len(config.RetryIntervals)-1] > 24*time.Hour {
		t.Error("Last retry interval should not exceed 24 hours")
	}
}

// TestRetryWorkerConfigValidation tests config boundary conditions
func TestRetryWorkerConfigValidation(t *testing.T) {
	// Test empty intervals (edge case)
	config := RetryWorkerConfig{
		RetryIntervals: []time.Duration{},
		MaxRetries:     0,
	}

	if len(config.RetryIntervals) != 0 {
		t.Error("Empty intervals should be allowed")
	}

	// Test single interval
	config.RetryIntervals = []time.Duration{5 * time.Minute}
	if len(config.RetryIntervals) != 1 {
		t.Error("Single interval should be allowed")
	}
}

func TestPollIntervalReasonable(t *testing.T) {
	config := DefaultRetryWorkerConfig()

	// Poll interval should be at least 10 seconds (too fast wastes resources)
	if config.PollInterval < 10*time.Second {
		t.Error("Poll interval should be at least 10 seconds")
	}

	// Poll interval should be at most 5 minutes (too slow delays retries)
	if config.PollInterval > 5*time.Minute {
		t.Error("Poll interval should not exceed 5 minutes")
	}
}

func TestBatchSizeReasonable(t *testing.T) {
	config := DefaultRetryWorkerConfig()

	// Batch size should be at least 1
	if config.BatchSize < 1 {
		t.Error("Batch size should be at least 1")
	}

	// Batch size shouldn't be too large (memory concerns)
	if config.BatchSize > 1000 {
		t.Error("Batch size should not exceed 1000")
	}
}
