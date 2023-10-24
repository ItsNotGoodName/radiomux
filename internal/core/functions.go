package core

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken() (string, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(secret), nil
}
