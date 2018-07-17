package security

import (
	"crypto/sha1"
	"encoding/hex"
)

func ShaOneEncrypt(s string) (string) {
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))

	return sha1_hash
}