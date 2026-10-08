package signalnet

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func hash(email string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(email))))
	return hex.EncodeToString(sum[:])
}
