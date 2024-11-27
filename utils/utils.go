package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashId(id int64, salt string) string {
	combined := fmt.Sprintf("%d%s", id, salt)
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:6])
}
