package sha

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sum256String(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
