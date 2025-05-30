package config

import (
	"errors"
	"os"
)

func getDBDNS() (string, error) {
	if value, exists := os.LookupEnv("PANDORA_DB_DNS"); exists {
		return value, nil
	}
	return "", errors.New("database DNS is required")
}

func getPort() string {
	if value, exists := os.LookupEnv("PANDORA_GRPC_PORT"); exists {
		return value
	}
	return "50051"
}
