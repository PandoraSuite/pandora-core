package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
)

func getDir() string {
	if value, exists := os.LookupEnv("PANDORA_DIR"); exists {
		return value
	}

	return "/etc/pandora"
}

func getDBDNS() string {
	if value, exists := os.LookupEnv("PANDORA_DB_DNS"); exists {
		return value
	}
	return "host=localhost port=5436 user=pandora password= dbname=pandora sslmode=disable timezone=UTC"
}

func getTaskEngineDBDNS() string {
	if value, exists := os.LookupEnv("PANDORA_TASKENGINE_DB_DNS"); exists {
		return value
	}
	return getDBDNS()
}

func getJWTSecrt() string {
	if value, exists := os.LookupEnv("PANDORA_JWT_SECRET"); exists {
		return value
	}

	log.Println("[WARNING] No JWT secret was provided. Using a randomly generated secret.")

	key := make([]byte, 64)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(key)
}

func getHTTPPort() string {
	if value, exists := os.LookupEnv("PANDORA_HTTP_PORT"); exists {
		return value
	}
	return "80"
}

func getGRPCPort() string {
	if value, exists := os.LookupEnv("PANDORA_GRPC_PORT"); exists {
		return value
	}
	return "50051"
}

func getExposeVersion() bool {
	if value, exists := os.LookupEnv("PANDORA_EXPOSE_VERSION"); exists {
		return value == "true"
	}
	return true
}
