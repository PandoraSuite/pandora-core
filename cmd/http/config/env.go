package config

import (
	"errors"
	"log"
	"os"
)

func getJWTSecrt() (string, error) {
	if value, exists := os.LookupEnv("PANDORA_JWT_SECRET"); exists {
		return value, nil
	}

	log.Println("[WARNING] No JWT secret was provided. Using a randomly generated secret.")

	randomKey, err := generateRandom(64)
	if err != nil {
		return "", err
	}

	return randomKey, nil
}

func getConfigDir() string {
	if value, exists := os.LookupEnv("PANDORA_CONFIG_DIR"); exists {
		return value
	}

	log.Println("[WARNING] No configuration directory was provided. Using default path: /etc/pandora")
	return "/etc/pandora"
}

func getDBDNS() (string, error) {
	if value, exists := os.LookupEnv("PANDORA_DB_DNS"); exists {
		return value, nil
	}
	return "", errors.New("database DNS is required")
}

func getHTTPPort() string {
	if value, exists := os.LookupEnv("PANDORA_HTTP_PORT"); exists {
		return value
	}
	return "80"
}
