package random

import (
	"crypto/rand"
	"encoding/base32"
)

func String(length int) (string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(randomBytes)[:length], nil
}
