package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// Service provides encryption/decryption for sensitive data.
type Service struct {
	key []byte
}

// NewService creates a new crypto service with the given hex-encoded 32-byte key.
func NewService(hexKey string) (*Service, error) {
	if hexKey == "" {
		return nil, errors.New("encryption key is required")
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, errors.New("invalid encryption key format: must be hex-encoded")
	}

	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes (64 hex characters)")
	}

	return &Service{key: key}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM.
func (s *Service) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM.
func (s *Service) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns encrypted bytes.
func (s *Service) EncryptString(plaintext string) ([]byte, error) {
	return s.Encrypt([]byte(plaintext))
}

// DecryptString decrypts ciphertext and returns a string.
func (s *Service) DecryptString(ciphertext []byte) (string, error) {
	plaintext, err := s.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
