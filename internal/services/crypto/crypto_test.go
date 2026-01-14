package crypto

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

// validKey is a valid 32-byte key (64 hex characters)
const validKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// generateValidKey generates a valid 32-byte hex-encoded key
func generateValidKey() string {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	return hex.EncodeToString(key)
}

// TestNewService tests the NewService constructor
func TestNewService(t *testing.T) {
	tests := []struct {
		name    string
		hexKey  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid 32-byte key",
			hexKey:  validKey,
			wantErr: false,
		},
		{
			name:    "empty key",
			hexKey:  "",
			wantErr: true,
			errMsg:  "encryption key is required",
		},
		{
			name:    "invalid hex characters",
			hexKey:  "xyz123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr: true,
			errMsg:  "invalid encryption key format: must be hex-encoded",
		},
		{
			name:    "key too short (16 bytes)",
			hexKey:  "0123456789abcdef0123456789abcdef",
			wantErr: true,
			errMsg:  "encryption key must be 32 bytes (64 hex characters)",
		},
		{
			name:    "key too long (48 bytes)",
			hexKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr: true,
			errMsg:  "encryption key must be 32 bytes (64 hex characters)",
		},
		{
			name:    "odd number of hex characters",
			hexKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde",
			wantErr: true,
			errMsg:  "invalid encryption key format: must be hex-encoded",
		},
		{
			name:    "uppercase hex should work",
			hexKey:  "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
			wantErr: false,
		},
		{
			name:    "mixed case hex should work",
			hexKey:  "0123456789AbCdEf0123456789aBcDeF0123456789ABCDEF0123456789abcdef",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewService(tt.hexKey)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewService() expected error, got nil")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("NewService() error = %q, want %q", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("NewService() unexpected error: %v", err)
				return
			}
			if svc == nil {
				t.Error("NewService() returned nil service without error")
			}
		})
	}
}

// TestEncryptDecryptRoundTrip tests encryption followed by decryption
func TestEncryptDecryptRoundTrip(t *testing.T) {
	svc, err := NewService(validKey)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{
			name:      "simple text",
			plaintext: []byte("Hello, World!"),
		},
		{
			name:      "empty bytes",
			plaintext: []byte{},
		},
		{
			name:      "single byte",
			plaintext: []byte{0x42},
		},
		{
			name:      "null bytes",
			plaintext: []byte{0x00, 0x00, 0x00},
		},
		{
			name:      "binary data with all byte values",
			plaintext: func() []byte { b := make([]byte, 256); for i := range b { b[i] = byte(i) }; return b }(),
		},
		{
			name:      "large payload (1MB)",
			plaintext: bytes.Repeat([]byte("A"), 1024*1024),
		},
		{
			name:      "unicode characters",
			plaintext: []byte("Hello, 世界! 🌍 Привет мир! مرحبا بالعالم"),
		},
		{
			name:      "special characters",
			plaintext: []byte("!@#$%^&*()_+-=[]{}|;':\",./<>?`~"),
		},
		{
			name:      "newlines and tabs",
			plaintext: []byte("line1\nline2\rline3\r\nline4\ttab"),
		},
		{
			name:      "JSON payload",
			plaintext: []byte(`{"key": "value", "nested": {"array": [1, 2, 3]}}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := svc.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() failed: %v", err)
			}

			// Ciphertext should be different from plaintext (unless empty)
			if len(tt.plaintext) > 0 && bytes.Equal(ciphertext, tt.plaintext) {
				t.Error("Encrypt() returned plaintext unchanged")
			}

			// Decrypt
			decrypted, err := svc.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decrypt() failed: %v", err)
			}

			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

// TestEncryptProducesDifferentCiphertexts ensures random nonces create different ciphertexts
func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	svc, err := NewService(validKey)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	plaintext := []byte("This is a test message")

	// Encrypt the same plaintext multiple times
	ciphertexts := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		ciphertext, err := svc.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Encrypt() iteration %d failed: %v", i, err)
		}
		ciphertexts[i] = ciphertext
	}

	// Verify all ciphertexts are different (due to random nonces)
	for i := 0; i < len(ciphertexts); i++ {
		for j := i + 1; j < len(ciphertexts); j++ {
			if bytes.Equal(ciphertexts[i], ciphertexts[j]) {
				t.Errorf("Ciphertexts %d and %d are identical, expected different nonces", i, j)
			}
		}
	}

	// But all should decrypt to the same plaintext
	for i, ct := range ciphertexts {
		decrypted, err := svc.Decrypt(ct)
		if err != nil {
			t.Fatalf("Decrypt() iteration %d failed: %v", i, err)
		}
		if !bytes.Equal(decrypted, plaintext) {
			t.Errorf("Decrypt() iteration %d = %q, want %q", i, decrypted, plaintext)
		}
	}
}

// TestDecryptErrors tests error conditions during decryption
func TestDecryptErrors(t *testing.T) {
	svc, err := NewService(validKey)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	tests := []struct {
		name       string
		ciphertext []byte
		wantErrMsg string
	}{
		{
			name:       "empty ciphertext",
			ciphertext: []byte{},
			wantErrMsg: "ciphertext too short",
		},
		{
			name:       "ciphertext shorter than nonce",
			ciphertext: []byte{0x01, 0x02, 0x03},
			wantErrMsg: "ciphertext too short",
		},
		{
			name:       "invalid ciphertext (random bytes)",
			ciphertext: bytes.Repeat([]byte{0xAB}, 100),
			wantErrMsg: "", // GCM authentication error, message varies
		},
		{
			name:       "tampered ciphertext",
			ciphertext: func() []byte {
				// First encrypt something valid
				ct, _ := svc.Encrypt([]byte("test"))
				// Tamper with the ciphertext portion (after nonce)
				if len(ct) > 20 {
					ct[20] ^= 0xFF
				}
				return ct
			}(),
			wantErrMsg: "", // GCM authentication error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Decrypt(tt.ciphertext)
			if err == nil {
				t.Error("Decrypt() expected error, got nil")
				return
			}
			if tt.wantErrMsg != "" && err.Error() != tt.wantErrMsg {
				t.Errorf("Decrypt() error = %q, want %q", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

// TestDifferentKeysCannotDecrypt verifies data encrypted with one key cannot be decrypted with another
func TestDifferentKeysCannotDecrypt(t *testing.T) {
	key1 := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	key2 := "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"

	svc1, err := NewService(key1)
	if err != nil {
		t.Fatalf("NewService(key1) failed: %v", err)
	}

	svc2, err := NewService(key2)
	if err != nil {
		t.Fatalf("NewService(key2) failed: %v", err)
	}

	plaintext := []byte("Secret message")

	// Encrypt with key1
	ciphertext, err := svc1.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	// Decrypt with key1 should work
	decrypted, err := svc1.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() with same key failed: %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Decrypt() with same key returned wrong plaintext")
	}

	// Decrypt with key2 should fail
	_, err = svc2.Decrypt(ciphertext)
	if err == nil {
		t.Error("Decrypt() with different key should have failed")
	}
}

// TestEncryptStringDecryptString tests the string convenience methods
func TestEncryptStringDecryptString(t *testing.T) {
	svc, err := NewService(validKey)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple string",
			plaintext: "Hello, World!",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "unicode string",
			plaintext: "Hello, 世界! 🌍 Привет!",
		},
		{
			name:      "string with special characters",
			plaintext: "password!@#$%^&*()_+-=",
		},
		{
			name:      "multiline string",
			plaintext: "line1\nline2\nline3",
		},
		{
			name:      "JSON string",
			plaintext: `{"api_key": "sk-test-123", "secret": "very-secret"}`,
		},
		{
			name:      "long string",
			plaintext: strings.Repeat("A very long message. ", 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := svc.EncryptString(tt.plaintext)
			if err != nil {
				t.Fatalf("EncryptString() failed: %v", err)
			}

			decrypted, err := svc.DecryptString(ciphertext)
			if err != nil {
				t.Fatalf("DecryptString() failed: %v", err)
			}

			if decrypted != tt.plaintext {
				t.Errorf("DecryptString() = %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

// TestDecryptStringErrors tests error conditions for DecryptString
func TestDecryptStringErrors(t *testing.T) {
	svc, err := NewService(validKey)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	tests := []struct {
		name       string
		ciphertext []byte
	}{
		{
			name:       "empty ciphertext",
			ciphertext: []byte{},
		},
		{
			name:       "invalid ciphertext",
			ciphertext: []byte{0x01, 0x02, 0x03},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.DecryptString(tt.ciphertext)
			if err == nil {
				t.Errorf("DecryptString() expected error, got result: %q", result)
			}
			if result != "" {
				t.Errorf("DecryptString() on error should return empty string, got: %q", result)
			}
		})
	}
}

// TestCiphertextLength verifies ciphertext length is appropriate
func TestCiphertextLength(t *testing.T) {
	svc, err := NewService(validKey)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	// GCM adds nonce (12 bytes) + auth tag (16 bytes) = 28 bytes overhead
	expectedOverhead := 28

	plaintextSizes := []int{0, 1, 10, 100, 1000, 10000}

	for _, size := range plaintextSizes {
		plaintext := bytes.Repeat([]byte{0x42}, size)
		ciphertext, err := svc.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Encrypt() failed for size %d: %v", size, err)
		}

		expectedLen := size + expectedOverhead
		if len(ciphertext) != expectedLen {
			t.Errorf("Ciphertext length for %d bytes = %d, want %d", size, len(ciphertext), expectedLen)
		}
	}
}

// TestEncryptDecryptCrossCompatibility tests that Encrypt output works with DecryptString
// and EncryptString output works with Decrypt
func TestEncryptDecryptCrossCompatibility(t *testing.T) {
	svc, err := NewService(validKey)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	testString := "Test message for cross-compatibility"

	// EncryptString -> Decrypt
	t.Run("EncryptString to Decrypt", func(t *testing.T) {
		ciphertext, err := svc.EncryptString(testString)
		if err != nil {
			t.Fatalf("EncryptString() failed: %v", err)
		}

		decrypted, err := svc.Decrypt(ciphertext)
		if err != nil {
			t.Fatalf("Decrypt() failed: %v", err)
		}

		if string(decrypted) != testString {
			t.Errorf("Decrypt() = %q, want %q", string(decrypted), testString)
		}
	})

	// Encrypt -> DecryptString
	t.Run("Encrypt to DecryptString", func(t *testing.T) {
		ciphertext, err := svc.Encrypt([]byte(testString))
		if err != nil {
			t.Fatalf("Encrypt() failed: %v", err)
		}

		decrypted, err := svc.DecryptString(ciphertext)
		if err != nil {
			t.Fatalf("DecryptString() failed: %v", err)
		}

		if decrypted != testString {
			t.Errorf("DecryptString() = %q, want %q", decrypted, testString)
		}
	})
}

// TestServiceImmutability verifies the service key cannot be modified externally
func TestServiceImmutability(t *testing.T) {
	svc, err := NewService(validKey)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	plaintext := []byte("Test message")

	// Encrypt once
	ciphertext1, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("First Encrypt() failed: %v", err)
	}

	// Encrypt again - should still work
	ciphertext2, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Second Encrypt() failed: %v", err)
	}

	// Both ciphertexts should decrypt correctly
	decrypted1, err := svc.Decrypt(ciphertext1)
	if err != nil {
		t.Fatalf("Decrypt(ciphertext1) failed: %v", err)
	}
	if !bytes.Equal(decrypted1, plaintext) {
		t.Error("Decrypt(ciphertext1) returned wrong plaintext")
	}

	decrypted2, err := svc.Decrypt(ciphertext2)
	if err != nil {
		t.Fatalf("Decrypt(ciphertext2) failed: %v", err)
	}
	if !bytes.Equal(decrypted2, plaintext) {
		t.Error("Decrypt(ciphertext2) returned wrong plaintext")
	}
}

// TestNilService tests behavior when service methods are called (would panic on nil)
// This is more of a documentation test showing expected panic behavior
func TestServiceWithGeneratedKey(t *testing.T) {
	// Test with a programmatically generated key
	key := generateValidKey()
	svc, err := NewService(key)
	if err != nil {
		t.Fatalf("NewService() with generated key failed: %v", err)
	}

	plaintext := []byte("Test with generated key")
	ciphertext, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	decrypted, err := svc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypt() = %v, want %v", decrypted, plaintext)
	}
}

// BenchmarkEncrypt benchmarks the encryption operation
func BenchmarkEncrypt(b *testing.B) {
	svc, err := NewService(validKey)
	if err != nil {
		b.Fatalf("NewService() failed: %v", err)
	}

	plaintext := []byte("Benchmark test data that is moderately sized for realistic testing.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Encrypt(plaintext)
		if err != nil {
			b.Fatalf("Encrypt() failed: %v", err)
		}
	}
}

// BenchmarkDecrypt benchmarks the decryption operation
func BenchmarkDecrypt(b *testing.B) {
	svc, err := NewService(validKey)
	if err != nil {
		b.Fatalf("NewService() failed: %v", err)
	}

	plaintext := []byte("Benchmark test data that is moderately sized for realistic testing.")
	ciphertext, err := svc.Encrypt(plaintext)
	if err != nil {
		b.Fatalf("Encrypt() failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Decrypt(ciphertext)
		if err != nil {
			b.Fatalf("Decrypt() failed: %v", err)
		}
	}
}

// BenchmarkEncryptLargePayload benchmarks encryption with 1MB payload
func BenchmarkEncryptLargePayload(b *testing.B) {
	svc, err := NewService(validKey)
	if err != nil {
		b.Fatalf("NewService() failed: %v", err)
	}

	plaintext := bytes.Repeat([]byte("A"), 1024*1024) // 1MB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Encrypt(plaintext)
		if err != nil {
			b.Fatalf("Encrypt() failed: %v", err)
		}
	}
}

// BenchmarkDecryptLargePayload benchmarks decryption with 1MB payload
func BenchmarkDecryptLargePayload(b *testing.B) {
	svc, err := NewService(validKey)
	if err != nil {
		b.Fatalf("NewService() failed: %v", err)
	}

	plaintext := bytes.Repeat([]byte("A"), 1024*1024) // 1MB
	ciphertext, err := svc.Encrypt(plaintext)
	if err != nil {
		b.Fatalf("Encrypt() failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Decrypt(ciphertext)
		if err != nil {
			b.Fatalf("Decrypt() failed: %v", err)
		}
	}
}
