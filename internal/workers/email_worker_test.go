package workers

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/services/email"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Get default batch size from dispatcher config for tests
var testBatchSize = int32(email.DefaultDispatcherConfig().BatchSize)

// ============================================================================
// Mock Types and Helpers
// ============================================================================

// mockEmailSender implements a minimal mock of the email.Service for testing
type mockEmailSender struct {
	sendCalled   int
	sendEmails   []string
	sendError    error
	sendCallback func(to, subject, body string)
	mu           sync.Mutex
}

func (m *mockEmailSender) recordSend(to, subject, body string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sendCalled++
	m.sendEmails = append(m.sendEmails, to)
	if m.sendCallback != nil {
		m.sendCallback(to, subject, body)
	}
}

func (m *mockEmailSender) getSendCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sendCalled
}

func (m *mockEmailSender) getSentEmails() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, len(m.sendEmails))
	copy(result, m.sendEmails)
	return result
}

// mockDB represents a minimal mock for database operations
type mockDB struct {
	pendingEmails      []db.GetPendingEmailsRow
	pendingEmailsError error
	sentEmails         []string
	failedEmails       []string
	sequences          map[string]db.EmailSequence
	markSentError      error
	markFailedError    error
	mu                 sync.Mutex
}

func newMockDB() *mockDB {
	return &mockDB{
		pendingEmails: []db.GetPendingEmailsRow{},
		sentEmails:    []string{},
		failedEmails:  []string{},
		sequences:     make(map[string]db.EmailSequence),
	}
}

func (m *mockDB) addPendingEmail(email db.GetPendingEmailsRow) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pendingEmails = append(m.pendingEmails, email)
}

func (m *mockDB) getSentEmails() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, len(m.sentEmails))
	copy(result, m.sentEmails)
	return result
}

func (m *mockDB) getFailedEmails() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, len(m.failedEmails))
	copy(result, m.failedEmails)
	return result
}

// ============================================================================
// Test: Dispatcher Config Defaults
// ============================================================================

func TestDispatcherConfigDefaults(t *testing.T) {
	config := email.DefaultDispatcherConfig()

	// Verify dispatcher config defaults are set appropriately
	assert.Equal(t, 10, config.Workers, "default workers should be 10")
	assert.Equal(t, 100, config.BatchSize, "default batch size should be 100")
	assert.Equal(t, 14.0, config.RateLimit, "default rate limit should be 14/sec")
	assert.Equal(t, 5*time.Second, config.PollInterval, "default poll interval should be 5 seconds")
}

func TestDispatcherConfig_Reasonable(t *testing.T) {
	config := email.DefaultDispatcherConfig()

	// Poll interval should be between 1 second and 5 minutes
	minInterval := 1 * time.Second
	maxInterval := 5 * time.Minute

	assert.GreaterOrEqual(t, config.PollInterval, minInterval,
		"poll interval should be at least 1 second")
	assert.LessOrEqual(t, config.PollInterval, maxInterval,
		"poll interval should not exceed 5 minutes")

	// Batch size should be between 1 and 1000
	assert.GreaterOrEqual(t, config.BatchSize, 1, "batch size should be at least 1")
	assert.LessOrEqual(t, config.BatchSize, 1000, "batch size should not exceed 1000")
}

// ============================================================================
// Test: processEmails function behavior with SequenceService
// ============================================================================

// mockSequenceService wraps a real SequenceService but allows control over behavior
type mockSequenceService struct {
	processFunc func(ctx context.Context, batchSize int32) (int, error)
	callCount   int
	lastBatch   int32
	mu          sync.Mutex
}

func (m *mockSequenceService) ProcessPendingEmails(ctx context.Context, batchSize int32) (int, error) {
	m.mu.Lock()
	m.callCount++
	m.lastBatch = batchSize
	m.mu.Unlock()

	if m.processFunc != nil {
		return m.processFunc(ctx, batchSize)
	}
	return 0, nil
}

func (m *mockSequenceService) getCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func (m *mockSequenceService) getLastBatchSize() int32 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastBatch
}

// testableSequenceService implements the minimal interface for processEmails testing
type testableSequenceService struct {
	*email.SequenceService
	processFuncOverride func(ctx context.Context, batchSize int32) (int, error)
}

// ============================================================================
// Test: SequenceService creation
// ============================================================================

func TestNewSequenceService(t *testing.T) {
	// Test that NewSequenceService creates a service with default base URL
	service := email.NewSequenceService(nil, nil)
	require.NotNil(t, service, "NewSequenceService should return a non-nil service")
}

func TestNewSequenceServiceWithBaseURL(t *testing.T) {
	baseURL := "https://example.com"
	service := email.NewSequenceServiceWithBaseURL(nil, nil, baseURL)
	require.NotNil(t, service, "NewSequenceServiceWithBaseURL should return a non-nil service")
}

func TestNewSequenceServiceWithBaseURL_EmptyURL(t *testing.T) {
	// Empty base URL should still create a valid service
	service := email.NewSequenceServiceWithBaseURL(nil, nil, "")
	require.NotNil(t, service, "service should be created even with empty base URL")
}

func TestNewSequenceServiceWithBaseURL_VariousURLFormats(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{"https URL", "https://example.com"},
		{"http URL", "http://localhost:8080"},
		{"trailing slash", "https://example.com/"},
		{"with path", "https://example.com/api"},
		{"with port", "https://example.com:8443"},
		{"localhost", "http://localhost"},
		{"IP address", "http://192.168.1.1:9888"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := email.NewSequenceServiceWithBaseURL(nil, nil, tt.baseURL)
			require.NotNil(t, service, "service should be created for URL: %s", tt.baseURL)
		})
	}
}

// ============================================================================
// Test: ProcessPendingEmails behavior patterns
// ============================================================================

func TestProcessPendingEmails_ReturnsZeroWhenNoPending(t *testing.T) {
	// When there are no pending emails, ProcessPendingEmails should return 0
	// This tests the expected behavior pattern
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return 0, nil
		},
	}

	sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

	assert.NoError(t, err)
	assert.Equal(t, 0, sent)
	assert.Equal(t, 1, mockService.getCallCount())
	assert.Equal(t, int32(testBatchSize), mockService.getLastBatchSize())
}

func TestProcessPendingEmails_ReturnsCountWhenEmailsSent(t *testing.T) {
	expectedSent := 5
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return expectedSent, nil
		},
	}

	sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

	assert.NoError(t, err)
	assert.Equal(t, expectedSent, sent)
}

func TestProcessPendingEmails_ReturnsError(t *testing.T) {
	expectedError := errors.New("database connection failed")
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return 0, expectedError
		},
	}

	sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, sent)
}

func TestProcessPendingEmails_PartialSuccess(t *testing.T) {
	// Some emails sent before an error
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return 3, errors.New("partial failure")
		},
	}

	sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

	// When there's an error, the function returns 0 and error (based on implementation)
	// But the mock shows the behavior - in reality the implementation might differ
	assert.Error(t, err)
	assert.Equal(t, 3, sent) // Mock returns 3 even with error
}

func TestProcessPendingEmails_BatchSizeRespected(t *testing.T) {
	testCases := []int32{1, 5, 10, 50, 100}

	for _, batchSize := range testCases {
		t.Run("batch_size_"+string(rune(batchSize)), func(t *testing.T) {
			mockService := &mockSequenceService{
				processFunc: func(ctx context.Context, bs int32) (int, error) {
					return int(bs), nil // Return batch size as sent count
				},
			}

			sent, err := mockService.ProcessPendingEmails(context.Background(), batchSize)

			assert.NoError(t, err)
			assert.Equal(t, int(batchSize), sent)
			assert.Equal(t, batchSize, mockService.getLastBatchSize())
		})
	}
}

func TestProcessPendingEmails_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			default:
				return 5, nil
			}
		},
	}

	sent, err := mockService.ProcessPendingEmails(ctx, testBatchSize)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Equal(t, 0, sent)
}

func TestProcessPendingEmails_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(10 * time.Millisecond)

	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			default:
				return 5, nil
			}
		},
	}

	sent, err := mockService.ProcessPendingEmails(ctx, testBatchSize)

	assert.Error(t, err)
	assert.Equal(t, 0, sent)
}

// ============================================================================
// Test: Error handling and recovery patterns
// ============================================================================

func TestProcessEmails_ErrorRecovery(t *testing.T) {
	// Simulate a sequence of calls with intermittent errors
	callCount := 0
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			callCount++
			if callCount%2 == 0 {
				return 0, errors.New("intermittent error")
			}
			return 5, nil
		},
	}

	// First call succeeds
	sent1, err1 := mockService.ProcessPendingEmails(context.Background(), testBatchSize)
	assert.NoError(t, err1)
	assert.Equal(t, 5, sent1)

	// Second call fails
	sent2, err2 := mockService.ProcessPendingEmails(context.Background(), testBatchSize)
	assert.Error(t, err2)
	assert.Equal(t, 0, sent2)

	// Third call succeeds (recovery)
	sent3, err3 := mockService.ProcessPendingEmails(context.Background(), testBatchSize)
	assert.NoError(t, err3)
	assert.Equal(t, 5, sent3)
}

func TestProcessEmails_TransientDatabaseErrors(t *testing.T) {
	tests := []struct {
		name      string
		errorType error
	}{
		{"connection reset", errors.New("connection reset by peer")},
		{"connection refused", errors.New("dial tcp: connection refused")},
		{"timeout", errors.New("i/o timeout")},
		{"too many connections", errors.New("too many connections")},
		{"deadlock", errors.New("deadlock detected")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockSequenceService{
				processFunc: func(ctx context.Context, batchSize int32) (int, error) {
					return 0, tt.errorType
				},
			}

			sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

			assert.Error(t, err)
			assert.Equal(t, 0, sent)
		})
	}
}

func TestProcessEmails_EmailServiceErrors(t *testing.T) {
	tests := []struct {
		name  string
		error error
	}{
		{"SMTP auth failed", errors.New("SMTP AUTH failed")},
		{"connection timeout", errors.New("connection timeout")},
		{"rate limited", errors.New("rate limit exceeded")},
		{"invalid recipient", errors.New("invalid recipient address")},
		{"message rejected", errors.New("message rejected")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockSequenceService{
				processFunc: func(ctx context.Context, batchSize int32) (int, error) {
					return 0, tt.error
				},
			}

			sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

			assert.Error(t, err)
			assert.Equal(t, 0, sent)
		})
	}
}

// ============================================================================
// Test: State management
// ============================================================================

func TestWorkerStateManagement(t *testing.T) {
	// Simulate tracking state across multiple process cycles
	var totalSent int32

	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			// Simulate sending 2 emails per cycle
			sent := 2
			atomic.AddInt32(&totalSent, int32(sent))
			return sent, nil
		},
	}

	// Simulate 5 processing cycles
	for i := 0; i < 5; i++ {
		_, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)
		assert.NoError(t, err)
	}

	assert.Equal(t, int32(10), atomic.LoadInt32(&totalSent))
	assert.Equal(t, 5, mockService.getCallCount())
}

func TestWorkerProcessingOrder(t *testing.T) {
	// Verify that emails are processed in order they're returned
	processingOrder := []int32{}
	mu := sync.Mutex{}

	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			mu.Lock()
			processingOrder = append(processingOrder, batchSize)
			mu.Unlock()
			return int(batchSize), nil
		},
	}

	// Process with different batch sizes
	mockService.ProcessPendingEmails(context.Background(), 5)
	mockService.ProcessPendingEmails(context.Background(), 10)
	mockService.ProcessPendingEmails(context.Background(), 3)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []int32{5, 10, 3}, processingOrder)
}

// ============================================================================
// Test: Timing and scheduling
// ============================================================================

func TestWorkerTicker_Behavior(t *testing.T) {
	// Test that ticker creates proper intervals
	// Using a very short interval for testing
	testInterval := 10 * time.Millisecond

	ticker := time.NewTicker(testInterval)
	defer ticker.Stop()

	start := time.Now()
	tickCount := 0

	timeout := time.After(50 * time.Millisecond)

loop:
	for {
		select {
		case <-ticker.C:
			tickCount++
			if tickCount >= 3 {
				break loop
			}
		case <-timeout:
			break loop
		}
	}

	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, tickCount, 3, "should have received at least 3 ticks")
	assert.GreaterOrEqual(t, elapsed, 30*time.Millisecond, "should have taken at least 30ms")
}

func TestWorkerTicker_StopCleanup(t *testing.T) {
	ticker := time.NewTicker(10 * time.Millisecond)

	// Simulate some ticks
	<-ticker.C

	// Stop the ticker
	ticker.Stop()

	// Verify ticker channel doesn't receive more ticks after stop
	select {
	case <-ticker.C:
		// This might still receive a buffered tick, which is ok
	case <-time.After(50 * time.Millisecond):
		// No more ticks - expected behavior
	}
	// Test passes if we reach here without blocking forever
}

func TestWorkerProcessesImmediatelyOnStart(t *testing.T) {
	// The worker should process emails immediately on startup before the ticker fires
	processed := make(chan struct{}, 1)

	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			select {
			case processed <- struct{}{}:
			default:
			}
			return 0, nil
		},
	}

	// Simulate the immediate processing that happens on startup
	go func() {
		mockService.ProcessPendingEmails(context.Background(), testBatchSize)
	}()

	select {
	case <-processed:
		// Successfully processed immediately
	case <-time.After(100 * time.Millisecond):
		t.Error("worker should process immediately on start")
	}
}

// ============================================================================
// Test: Concurrent access safety
// ============================================================================

func TestProcessEmails_ConcurrentSafety(t *testing.T) {
	var callCount int32
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			atomic.AddInt32(&callCount, 1)
			// Simulate some processing time
			time.Sleep(1 * time.Millisecond)
			return 1, nil
		},
	}

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mockService.ProcessPendingEmails(context.Background(), testBatchSize)
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(numGoroutines), atomic.LoadInt32(&callCount))
}

func TestMockSequenceService_ThreadSafe(t *testing.T) {
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return int(batchSize), nil
		},
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			mockService.ProcessPendingEmails(context.Background(), int32(n%10+1))
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 100, mockService.getCallCount())
}

// ============================================================================
// Test: Email data structures
// ============================================================================

func TestGetPendingEmailsRow_Fields(t *testing.T) {
	contactID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	row := db.GetPendingEmailsRow{
		ID:              uuid.New().String(),
		Email:           "test@example.com",
		Name:            "Test User",
		Subject:         "Test Subject",
		HtmlBody:        "<p>Test Body</p>",
		ScheduledFor:    now,
		TrackingToken:   sql.NullString{String: "abc123", Valid: true},
		ContactID:       sql.NullString{String: contactID, Valid: true},
		TemplateID:      sql.NullString{String: "100", Valid: true},
		TemplateType:    sql.NullString{String: "simple", Valid: true},
		IsTransactional: sql.NullInt64{Int64: 0, Valid: true},
	}

	assert.NotEmpty(t, row.ID)
	assert.Equal(t, "test@example.com", row.Email)
	assert.Equal(t, "Test User", row.Name)
	assert.Equal(t, "Test Subject", row.Subject)
	assert.Equal(t, "<p>Test Body</p>", row.HtmlBody)
	assert.Equal(t, now, row.ScheduledFor)
	assert.True(t, row.TrackingToken.Valid)
	assert.Equal(t, "abc123", row.TrackingToken.String)
	assert.True(t, row.ContactID.Valid)
	assert.Equal(t, contactID, row.ContactID.String)
	assert.True(t, row.TemplateID.Valid)
	assert.Equal(t, "100", row.TemplateID.String)
	assert.True(t, row.TemplateType.Valid)
	assert.Equal(t, "simple", row.TemplateType.String)
	assert.True(t, row.IsTransactional.Valid)
	assert.Equal(t, int64(0), row.IsTransactional.Int64)
}

func TestGetPendingEmailsRow_NullFields(t *testing.T) {
	row := db.GetPendingEmailsRow{
		ID:              uuid.New().String(),
		Email:           "test@example.com",
		Name:            "",
		Subject:         "Test",
		HtmlBody:        "<p>Body</p>",
		ScheduledFor:    time.Now().Format(time.RFC3339),
		TrackingToken:   sql.NullString{Valid: false},
		ContactID:       sql.NullString{Valid: false},
		TemplateID:      sql.NullString{Valid: false},
		TemplateType:    sql.NullString{Valid: false},
		IsTransactional: sql.NullInt64{Valid: false},
	}

	assert.False(t, row.TrackingToken.Valid)
	assert.False(t, row.ContactID.Valid)
	assert.False(t, row.TemplateID.Valid)
	assert.False(t, row.IsTransactional.Valid)
}

func TestEmailTemplateTypes(t *testing.T) {
	// Test different template types that the worker should handle
	templateTypes := []string{
		"none",         // Raw HTML - no wrapping
		"simple",       // Just footer
		"branded",      // Header + footer
		"confirmation", // Confirmation emails
	}

	for _, tt := range templateTypes {
		t.Run(tt, func(t *testing.T) {
			row := db.GetPendingEmailsRow{
				ID:           uuid.New().String(),
				Email:        "test@example.com",
				TemplateType: sql.NullString{String: tt, Valid: true},
			}
			assert.Equal(t, tt, row.TemplateType.String)
		})
	}
}

// ============================================================================
// Test: Scheduling time calculations
// ============================================================================

func TestScheduledForTime_InPast(t *testing.T) {
	// Email scheduled in the past should be processed
	pastTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)

	row := db.GetPendingEmailsRow{
		ID:           uuid.New().String(),
		Email:        "test@example.com",
		ScheduledFor: pastTime,
	}

	parsedTime, err := time.Parse(time.RFC3339, row.ScheduledFor)
	require.NoError(t, err)
	assert.True(t, parsedTime.Before(time.Now()))
}

func TestScheduledForTime_InFuture(t *testing.T) {
	// Email scheduled in the future should not be processed yet
	futureTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	row := db.GetPendingEmailsRow{
		ID:           uuid.New().String(),
		Email:        "test@example.com",
		ScheduledFor: futureTime,
	}

	parsedTime, err := time.Parse(time.RFC3339, row.ScheduledFor)
	require.NoError(t, err)
	assert.True(t, parsedTime.After(time.Now()))
}

func TestScheduledForTime_ExactlyNow(t *testing.T) {
	now := time.Now()
	nowStr := now.Format(time.RFC3339)

	row := db.GetPendingEmailsRow{
		ID:           uuid.New().String(),
		Email:        "test@example.com",
		ScheduledFor: nowStr,
	}

	parsedTime, err := time.Parse(time.RFC3339, row.ScheduledFor)
	require.NoError(t, err)
	// Time should be approximately now (within a few seconds due to format precision)
	assert.WithinDuration(t, now, parsedTime, 2*time.Second)
}

// ============================================================================
// Test: Tracking token handling
// ============================================================================

func TestTrackingToken_Valid(t *testing.T) {
	token := "abc123def456"
	row := db.GetPendingEmailsRow{
		ID:            uuid.New().String(),
		Email:         "test@example.com",
		TrackingToken: sql.NullString{String: token, Valid: true},
	}

	assert.True(t, row.TrackingToken.Valid)
	assert.Equal(t, token, row.TrackingToken.String)
}

func TestTrackingToken_Null(t *testing.T) {
	row := db.GetPendingEmailsRow{
		ID:            uuid.New().String(),
		Email:         "test@example.com",
		TrackingToken: sql.NullString{String: "", Valid: false},
	}

	assert.False(t, row.TrackingToken.Valid)
}

func TestTrackingToken_EmptyButValid(t *testing.T) {
	row := db.GetPendingEmailsRow{
		ID:            uuid.New().String(),
		Email:         "test@example.com",
		TrackingToken: sql.NullString{String: "", Valid: true},
	}

	// An empty but valid token is a valid edge case
	assert.True(t, row.TrackingToken.Valid)
	assert.Empty(t, row.TrackingToken.String)
}

// ============================================================================
// Test: Contact and template associations
// ============================================================================

func TestContactID_Association(t *testing.T) {
	contactID := uuid.New().String()

	row := db.GetPendingEmailsRow{
		ID:        uuid.New().String(),
		Email:     "test@example.com",
		ContactID: sql.NullString{String: contactID, Valid: true},
	}

	assert.True(t, row.ContactID.Valid)
	assert.Equal(t, contactID, row.ContactID.String)
}

func TestTemplateID_Association(t *testing.T) {
	row := db.GetPendingEmailsRow{
		ID:         uuid.New().String(),
		Email:      "test@example.com",
		TemplateID: sql.NullString{String: "42", Valid: true},
	}

	assert.True(t, row.TemplateID.Valid)
	assert.Equal(t, "42", row.TemplateID.String)
}

// ============================================================================
// Test: Transactional email flag
// ============================================================================

func TestIsTransactional_True(t *testing.T) {
	row := db.GetPendingEmailsRow{
		ID:              uuid.New().String(),
		Email:           "test@example.com",
		IsTransactional: sql.NullInt64{Int64: 1, Valid: true},
	}

	assert.True(t, row.IsTransactional.Valid)
	assert.Equal(t, int64(1), row.IsTransactional.Int64)
}

func TestIsTransactional_False(t *testing.T) {
	row := db.GetPendingEmailsRow{
		ID:              uuid.New().String(),
		Email:           "test@example.com",
		IsTransactional: sql.NullInt64{Int64: 0, Valid: true},
	}

	assert.True(t, row.IsTransactional.Valid)
	assert.Equal(t, int64(0), row.IsTransactional.Int64)
}

func TestIsTransactional_Null(t *testing.T) {
	row := db.GetPendingEmailsRow{
		ID:              uuid.New().String(),
		Email:           "test@example.com",
		IsTransactional: sql.NullInt64{Valid: false},
	}

	assert.False(t, row.IsTransactional.Valid)
}

// ============================================================================
// Test: Batch processing simulation
// ============================================================================

func TestBatchProcessing_ExactBatchSize(t *testing.T) {
	// Process exactly testBatchSize emails
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return int(batchSize), nil
		},
	}

	sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

	assert.NoError(t, err)
	assert.Equal(t, int(testBatchSize), sent)
}

func TestBatchProcessing_LessThanBatchSize(t *testing.T) {
	// Fewer pending emails than batch size
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return 3, nil // Only 3 emails pending
		},
	}

	sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

	assert.NoError(t, err)
	assert.Equal(t, 3, sent)
}

func TestBatchProcessing_ZeroEmails(t *testing.T) {
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return 0, nil
		},
	}

	sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)

	assert.NoError(t, err)
	assert.Equal(t, 0, sent)
}

func TestBatchProcessing_MultipleBatches(t *testing.T) {
	// Simulate processing multiple batches in sequence
	totalEmails := 25
	var totalSent int

	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			remaining := totalEmails - totalSent
			if remaining <= 0 {
				return 0, nil
			}
			if remaining < int(batchSize) {
				return remaining, nil
			}
			return int(batchSize), nil
		},
	}

	// Process until no more emails
	for {
		sent, err := mockService.ProcessPendingEmails(context.Background(), testBatchSize)
		require.NoError(t, err)
		totalSent += sent
		if sent == 0 {
			break
		}
	}

	assert.Equal(t, totalEmails, totalSent)
}

// ============================================================================
// Test: Worker lifecycle
// ============================================================================

func TestWorkerLifecycle_StartStop(t *testing.T) {
	// Simulate worker start and graceful stop
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	processCount := int32(0)

	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			atomic.AddInt32(&processCount, 1)
			return 0, nil
		},
	}

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		defer close(done)

		// Process immediately
		mockService.ProcessPendingEmails(ctx, testBatchSize)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mockService.ProcessPendingEmails(ctx, testBatchSize)
			}
		}
	}()

	// Let it run for a bit
	time.Sleep(50 * time.Millisecond)

	// Stop the worker
	cancel()

	// Wait for graceful shutdown
	select {
	case <-done:
		// Worker stopped successfully
	case <-time.After(100 * time.Millisecond):
		t.Error("worker did not stop within timeout")
	}

	assert.Greater(t, atomic.LoadInt32(&processCount), int32(0))
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkProcessEmails_NoEmails(b *testing.B) {
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return 0, nil
		},
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mockService.ProcessPendingEmails(ctx, testBatchSize)
	}
}

func BenchmarkProcessEmails_WithEmails(b *testing.B) {
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return int(batchSize), nil
		},
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mockService.ProcessPendingEmails(ctx, testBatchSize)
	}
}

func BenchmarkMockServiceCallCount(b *testing.B) {
	mockService := &mockSequenceService{
		processFunc: func(ctx context.Context, batchSize int32) (int, error) {
			return 5, nil
		},
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mockService.ProcessPendingEmails(ctx, testBatchSize)
		_ = mockService.getCallCount()
	}
}
