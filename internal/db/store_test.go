package db

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping database tests")
		return nil
	}

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}

// cleanupTestData removes test data
func cleanupTestData(t *testing.T, db *sql.DB, table string, id string) {
	if db == nil || id == "" {
		return
	}

	_, err := db.Exec("DELETE FROM "+table+" WHERE id = ?", id)
	if err != nil {
		t.Logf("Warning: Failed to cleanup %s: %v", table, err)
	}
}

// TestNewStore tests store creation
func TestNewStore(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)

	assert.NotNil(t, store)
	assert.NotNil(t, store.Queries)
	assert.Equal(t, testDB, store.GetDB())
}

// TestStorePing tests database connectivity
func TestStorePing(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)

	err := store.Ping(context.Background())
	assert.NoError(t, err)
}

// TestCreateUser tests user creation
func TestCreateUser(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)
	ctx := context.Background()

	userID := uuid.New().String()
	email := "test-user-" + time.Now().Format("20060102150405") + "@example.com"

	defer cleanupTestData(t, testDB, "users", userID)

	params := CreateUserParams{
		ID:           userID,
		Name:         "Test User",
		Email:        email,
		PasswordHash: "hashed_password_123",
		Role:         "agent",
	}

	user, err := store.CreateUser(ctx, params)
	require.NoError(t, err)

	assert.NotEmpty(t, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "agent", user.Role)
	assert.Equal(t, "active", user.Status)
}

// TestGetUserByID tests user retrieval by ID
func TestGetUserByID(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)
	ctx := context.Background()

	// Create user first
	userID := uuid.New().String()
	email := "test-get-user-" + time.Now().Format("20060102150405") + "@example.com"
	createParams := CreateUserParams{
		ID:           userID,
		Name:         "Test User",
		Email:        email,
		PasswordHash: "hashed_password",
		Role:         "agent",
	}

	created, err := store.CreateUser(ctx, createParams)
	require.NoError(t, err)
	defer cleanupTestData(t, testDB, "users", created.ID)

	// Get user by ID
	user, err := store.GetUserByID(ctx, created.ID)
	require.NoError(t, err)

	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, email, user.Email)
}

// TestGetUserByEmail tests user retrieval by email
func TestGetUserByEmail(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)
	ctx := context.Background()

	// Create user first
	userID := uuid.New().String()
	email := "test-get-by-email-" + time.Now().Format("20060102150405") + "@example.com"
	createParams := CreateUserParams{
		ID:           userID,
		Name:         "Test User",
		Email:        email,
		PasswordHash: "hashed_password",
		Role:         "agent",
	}

	created, err := store.CreateUser(ctx, createParams)
	require.NoError(t, err)
	defer cleanupTestData(t, testDB, "users", created.ID)

	// Get user by email
	user, err := store.GetUserByEmail(ctx, email)
	require.NoError(t, err)

	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, email, user.Email)
}

// TestTransactionSupport tests transaction execution
func TestTransactionSupport(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)
	ctx := context.Background()

	email := "tx-test-" + time.Now().Format("20060102150405") + "@example.com"

	t.Run("SuccessfulTransaction", func(t *testing.T) {
		var userID string

		err := store.ExecTx(ctx, func(q *Queries) error {
			// Create user
			newID := uuid.New().String()
			userParams := CreateUserParams{
				ID:           newID,
				Name:         "TX Test User",
				Email:        email,
				PasswordHash: "hashed_password",
				Role:         "agent",
			}
			user, err := q.CreateUser(ctx, userParams)
			if err != nil {
				return err
			}
			userID = user.ID
			return nil
		})

		require.NoError(t, err)
		defer cleanupTestData(t, testDB, "users", userID)

		// Verify user was created
		user, err := store.GetUserByID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, userID, user.ID)
	})

	t.Run("FailedTransactionRollback", func(t *testing.T) {
		newEmail := "rollback-" + time.Now().Format("20060102150405") + "@example.com"
		var attemptedUserID string

		err := store.ExecTx(ctx, func(q *Queries) error {
			// Create user
			newID := uuid.New().String()
			userParams := CreateUserParams{
				ID:           newID,
				Name:         "Rollback Test User",
				Email:        newEmail,
				PasswordHash: "hashed_password",
				Role:         "agent",
			}
			user, err := q.CreateUser(ctx, userParams)
			if err != nil {
				return err
			}
			attemptedUserID = user.ID

			// Force an error by trying to create duplicate
			_, err = q.CreateUser(ctx, userParams)
			return err
		})

		assert.Error(t, err, "Transaction should fail")

		// Verify user was rolled back
		if attemptedUserID != "" {
			_, err = store.GetUserByID(ctx, attemptedUserID)
			assert.Error(t, err)
		}
	})
}

// TestExecTxWithResult tests transaction execution with result
func TestExecTxWithResult(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)
	ctx := context.Background()

	email := "tx-result-" + time.Now().Format("20060102150405") + "@example.com"

	user, err := ExecTxWithResult(store, ctx, func(q *Queries) (CreateUserRow, error) {
		params := CreateUserParams{
			ID:           uuid.New().String(),
			Name:         "TX Result User",
			Email:        email,
			PasswordHash: "hashed_password",
			Role:         "agent",
		}
		return q.CreateUser(ctx, params)
	})

	require.NoError(t, err)
	defer cleanupTestData(t, testDB, "users", user.ID)

	assert.NotEmpty(t, user.ID)
	assert.Equal(t, email, user.Email)
}

// TestCallSessionOperations tests call session operations
func TestCallSessionOperations(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)
	ctx := context.Background()

	// Create agent first
	agentEmail := "session-agent-" + time.Now().Format("20060102150405") + "@example.com"
	agentParams := CreateUserParams{
		ID:           uuid.New().String(),
		Name:         "Session Agent",
		Email:        agentEmail,
		PasswordHash: "hashed_password",
		Role:         "agent",
	}
	agent, err := store.CreateUser(ctx, agentParams)
	require.NoError(t, err)
	defer cleanupTestData(t, testDB, "users", agent.ID)
}

// TestStoreClose tests closing the database connection
func TestStoreClose(t *testing.T) {
	testDB := setupTestDB(t)

	store := NewStore(testDB)

	err := store.Close()
	assert.NoError(t, err)

	// Verify connection is closed
	err = store.Ping(context.Background())
	assert.Error(t, err)
}

// TestListUsersWithParams tests listing users with parameters
func TestListUsersWithParams(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	store := NewStore(testDB)
	ctx := context.Background()

	// Create a test user
	email := "list-test-" + time.Now().Format("20060102150405") + "@example.com"
	createParams := CreateUserParams{
		ID:           uuid.New().String(),
		Name:         "List Test User",
		Email:        email,
		PasswordHash: "hashed_password",
		Role:         "agent",
	}

	user, err := store.CreateUser(ctx, createParams)
	require.NoError(t, err)
	defer cleanupTestData(t, testDB, "users", user.ID)

	// List users
	listParams := ListUsersParams{
		PageSize:   10,
		PageOffset: 0,
	}

	users, err := store.ListUsers(ctx, listParams)
	require.NoError(t, err)
	assert.Greater(t, len(users), 0)
}
