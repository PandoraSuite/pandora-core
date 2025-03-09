package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func existPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func generateRandom(length int) (string, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}

func calculateHash(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), 12)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}
