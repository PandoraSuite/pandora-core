package config

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func getCredentialsFilePath(dir string) string {
	dir = dir + "/adminPanel"
	if !existPath(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}

	credentialsFile := dir + "/credentials.json"
	if !existPath(credentialsFile) {
		createCredentials(credentialsFile)
	}
	return credentialsFile
}

func createCredentials(credentialsPath string) {
	log.Println("[WARNING] No admin credentials were found, generating default credentials.")

	key := make([]byte, 12)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	password, err := calculateHash(base64.URLEncoding.EncodeToString(key))
	if err != nil {
		panic(err)
	}

	credentials := struct {
		Username           string `json:"username"`
		Password           string `json:"password"`
		ForcePasswordReset bool   `json:"force_password_reset"`
	}{
		Username:           "Admin",
		Password:           password,
		ForcePasswordReset: true,
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(credentialsPath, data, 0644)
	if err != nil {
		panic(err)
	}

	log.Println("[INFO] Default admin credentials:")
	log.Printf("        - Username: Admin\n")
	log.Printf("        - Password: %s\n", password)
	log.Println("[SECURITY] It is highly recommended to change the default password immediately.")
}

func calculateHash(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), 12)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func existPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
