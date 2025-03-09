package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	dir string

	dbDNS string

	httpPort string

	jwtSecret string
}

func (c *Config) HTTPPort() string { return c.httpPort }

func (c *Config) DBDNS() string { return c.dbDNS }

func (c *Config) JWTSecret() string { return c.jwtSecret }

func (c *Config) CredentialsFile() (string, error) {
	dir := c.dir + "/adminPanel"
	if !existPath(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
		return c.createCredentials(dir + "/credentials.json")
	}

	credentialsFile := dir + "/credentials.json"
	if !existPath(credentialsFile) {
		return c.createCredentials(credentialsFile)
	}
	return credentialsFile, nil
}

func (c *Config) createCredentials(credentialsFile string) (string, error) {
	password, err := generateRandom(12)
	if err != nil {
		return "", err
	}

	passwordHash, err := calculateHash(password)
	if err != nil {
		return "", err
	}

	credentials := struct {
		Username           string `json:"username"`
		Password           string `json:"password"`
		ForcePasswordReset bool   `json:"force_password_reset"`
	}{
		Username:           "Admin",
		Password:           passwordHash,
		ForcePasswordReset: true,
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(credentialsFile, data, 0644)
	if err != nil {
		return "", err
	}

	return credentialsFile, nil
}

func LoadConfig() (*Config, error) {
	dir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	if !existPath(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	jwtSecret, err := getJWTSecrt()
	if err != nil {
		return nil, err
	}

	dbDNS, err := getDBDNS()
	if err != nil {
		return nil, err
	}

	httpPort := getHTTPPort()

	return &Config{
		dir:       dir,
		dbDNS:     dbDNS,
		httpPort:  httpPort,
		jwtSecret: jwtSecret,
	}, nil
}
