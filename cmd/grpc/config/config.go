package config

type Config struct {
	dbDNS string

	port string
}

func (c *Config) Port() string { return c.port }

func (c *Config) DBDNS() string { return c.dbDNS }

func LoadConfig() (*Config, error) {
	dbDNS, err := getDBDNS()
	if err != nil {
		return nil, err
	}

	port := getPort()
	return &Config{
		dbDNS: dbDNS,
		port:  port,
	}, nil
}
