package config

import (
	"errors"
	"os"
)

func getJWTSecrt() (string, error) {
	if value, exists := os.LookupEnv("PANDORA_JWT_SECRET"); exists {
		return value, nil
	}

	randomKey, err := generateRandom(64)
	if err != nil {
		return "", err
	}

	if err := os.Setenv("PANDORA_JWT_SECRET", randomKey); err != nil {
		return "", err
	}

	return randomKey, nil
}

func getConfigDir() (string, error) {
	if value, exists := os.LookupEnv("PANDORA_CONFIG_DIR"); exists {
		return value, nil
	}

	defaultDir := "/etc/pandora"
	if err := os.Setenv("PANDORA_JWT_SECRET", defaultDir); err != nil {
		return "", err
	}

	return defaultDir, nil
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
