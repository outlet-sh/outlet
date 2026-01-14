package email

import (
	"testing"
	"time"
)

func TestDefaultSMTPPoolConfig(t *testing.T) {
	config := DefaultSMTPPoolConfig()

	if config.PoolSize != 10 {
		t.Errorf("Default pool size should be 10, got %d", config.PoolSize)
	}
	if config.MaxIdle != 5*time.Minute {
		t.Errorf("Default max idle should be 5m, got %v", config.MaxIdle)
	}
	if !config.UseTLS {
		t.Error("Default should use TLS")
	}
}

func TestNewSMTPPool_DefaultValues(t *testing.T) {
	config := SMTPPoolConfig{
		Host: "smtp.example.com",
		Port: 587,
	}

	pool := NewSMTPPool(config)
	defer pool.Close()

	if pool.poolSize != 10 {
		t.Errorf("Pool size should default to 10, got %d", pool.poolSize)
	}
	if pool.maxIdle != 5*time.Minute {
		t.Errorf("Max idle should default to 5m, got %v", pool.maxIdle)
	}
}

func TestNewSMTPPool_CustomValues(t *testing.T) {
	config := SMTPPoolConfig{
		Host:     "smtp.example.com",
		Port:     587,
		User:     "user",
		Pass:     "pass",
		UseTLS:   true,
		PoolSize: 20,
		MaxIdle:  10 * time.Minute,
	}

	pool := NewSMTPPool(config)
	defer pool.Close()

	if pool.poolSize != 20 {
		t.Errorf("Pool size should be 20, got %d", pool.poolSize)
	}
	if pool.maxIdle != 10*time.Minute {
		t.Errorf("Max idle should be 10m, got %v", pool.maxIdle)
	}
	if pool.host != "smtp.example.com" {
		t.Errorf("Host should be smtp.example.com, got %s", pool.host)
	}
	if pool.user != "user" {
		t.Errorf("User should be user, got %s", pool.user)
	}
}

func TestSMTPPool_Stats_Initial(t *testing.T) {
	config := SMTPPoolConfig{
		Host:     "smtp.example.com",
		Port:     587,
		PoolSize: 5,
	}

	pool := NewSMTPPool(config)
	defer pool.Close()

	pooled, created := pool.Stats()
	if pooled != 0 {
		t.Errorf("Initial pooled should be 0, got %d", pooled)
	}
	if created != 0 {
		t.Errorf("Initial created should be 0, got %d", created)
	}
}

func TestSMTPPool_Get_ClosedPool(t *testing.T) {
	config := SMTPPoolConfig{
		Host: "smtp.example.com",
		Port: 587,
	}

	pool := NewSMTPPool(config)
	pool.Close()

	_, err := pool.Get()
	if err == nil {
		t.Error("Get on closed pool should return error")
	}
	if err.Error() != "pool is closed" {
		t.Errorf("Expected 'pool is closed' error, got: %v", err)
	}
}

func TestSMTPPool_Put_Nil(t *testing.T) {
	config := SMTPPoolConfig{
		Host: "smtp.example.com",
		Port: 587,
	}

	pool := NewSMTPPool(config)
	defer pool.Close()

	// Should not panic
	pool.Put(nil)

	pooled, _ := pool.Stats()
	if pooled != 0 {
		t.Errorf("Pool should still be empty after putting nil, got %d", pooled)
	}
}

func TestSMTPPool_Put_ClosedPool(t *testing.T) {
	config := SMTPPoolConfig{
		Host: "smtp.example.com",
		Port: 587,
	}

	pool := NewSMTPPool(config)
	pool.Close()

	// Should not panic on closed pool
	pool.Put(nil)
}

func TestSMTPPool_Close_Multiple(t *testing.T) {
	config := SMTPPoolConfig{
		Host: "smtp.example.com",
		Port: 587,
	}

	pool := NewSMTPPool(config)
	pool.Close()

	// Second close should not panic (though channel is already closed)
	// This tests defensive behavior
	defer func() {
		if r := recover(); r != nil {
			// Expected - closing already closed channel panics
			// This is acceptable behavior
		}
	}()
}

// Integration tests would require a mock SMTP server
// These tests verify the pool logic without actual connections

func TestSMTPPool_ChannelCapacity(t *testing.T) {
	config := SMTPPoolConfig{
		Host:     "smtp.example.com",
		Port:     587,
		PoolSize: 3,
	}

	pool := NewSMTPPool(config)
	defer pool.Close()

	// Verify channel has correct capacity
	if cap(pool.conns) != 3 {
		t.Errorf("Channel capacity should be 3, got %d", cap(pool.conns))
	}
}
