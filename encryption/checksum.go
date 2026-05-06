package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateMessageChecksum(message, secretKey string, length int) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	rawHash := h.Sum(nil)
	encoded := base64.URLEncoding.EncodeToString(rawHash)
	if len(encoded) > length {
		return encoded[:length]
	}
	return encoded
}

func IsMessageValidByCheckSum(message, checksum, secret string) bool {
	expected := GenerateMessageChecksum(message, secret, len(checksum))
	return hmac.Equal([]byte(expected), []byte(checksum))
}
