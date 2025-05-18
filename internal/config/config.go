package config

type Config struct {
	http *HTTPConfig
	grpc *GRPCConfig
}

func (c *Config) DBDNS() string {
	if c.http != nil {
		return c.http.dbDNS
	}

	if c.grpc != nil {
		return c.grpc.dbDNS
	}

	return ""
}

func (c *Config) HTTPConfig() *HTTPConfig { return c.http }

func (c *Config) GRPCConfig() *GRPCConfig { return c.grpc }

type baseConfig struct {
	dbDNS string
}

type HTTPConfig struct {
	*baseConfig

	dir string

	port string

	exposeVersion bool

	jwtSecret string

	credentialsFile string
}

func (c *HTTPConfig) DBDNS() string { return c.dbDNS }

func (c *HTTPConfig) Dir() string { return c.dir }

func (c *HTTPConfig) Port() string { return c.port }

func (c *HTTPConfig) ExposeVersion() bool { return c.exposeVersion }

func (c *HTTPConfig) JWTSecret() string { return c.jwtSecret }

func (c *HTTPConfig) CredentialsFile() string { return c.credentialsFile }

type GRPCConfig struct {
	*baseConfig

	port string
}

func (c *GRPCConfig) DBDNS() string { return c.dbDNS }

func (c *GRPCConfig) Port() string { return c.port }

func loadBaseConfig() *baseConfig {
	return &baseConfig{
		dbDNS: getDBDNS(),
	}
}

func LoadConfig() *Config {
	return &Config{
		http: LoadHTTPConfig(),
		grpc: LoadGRPCConfig(),
	}
}

func LoadHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		dir:             getDir(),
		port:            getHTTPPort(),
		jwtSecret:       getJWTSecrt(),
		baseConfig:      loadBaseConfig(),
		exposeVersion:   getExposeVersion(),
		credentialsFile: getCredentialsFilePath(getDir()),
	}
}

func LoadGRPCConfig() *GRPCConfig {
	return &GRPCConfig{
		port:       getGRPCPort(),
		baseConfig: loadBaseConfig(),
	}
}
