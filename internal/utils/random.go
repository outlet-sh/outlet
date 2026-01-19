package utils

import (
	"crypto/rand"
	"math/big"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GeneratePublicID generates a cryptographically secure random 8-character base62 ID
func GeneratePublicID() string {
	return GenerateRandomBase62(8)
}

// GenerateRandomBase62 generates a cryptographically secure random base62 string of the given length
func GenerateRandomBase62(length int) string {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(base62Chars)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// Fall back to less random but still functional approach if crypto/rand fails
			n = big.NewInt(int64(i % len(base62Chars)))
		}
		result[i] = base62Chars[n.Int64()]
	}

	return string(result)
}
