package config

type Config struct {
	http       *HTTPConfig
	grpc       *GRPCConfig
	taskEngine *TaskEngineConfig
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

func (c *Config) TaskEngineConfig() *TaskEngineConfig { return c.taskEngine }

type baseConfig struct {
	dbDNS string
}

func (c *baseConfig) DBDNS() string { return c.dbDNS }

type HTTPConfig struct {
	*baseConfig

	dir string

	port string

	exposeVersion bool

	jwtSecret string

	credentialsFile string
}

func (c *HTTPConfig) Dir() string { return c.dir }

func (c *HTTPConfig) Port() string { return c.port }

func (c *HTTPConfig) ExposeVersion() bool { return c.exposeVersion }

func (c *HTTPConfig) JWTSecret() string { return c.jwtSecret }

func (c *HTTPConfig) CredentialsFile() string { return c.credentialsFile }

type GRPCConfig struct {
	*baseConfig

	port string
}

func (c *GRPCConfig) Port() string { return c.port }

type TaskEngineConfig struct {
	*baseConfig
}

func LoadConfig() *Config {
	return &Config{
		http:       LoadHTTPConfig(),
		grpc:       LoadGRPCConfig(),
		taskEngine: LoadTaskEngineConfig(),
	}
}

func LoadHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		dir:       getDir(),
		port:      getHTTPPort(),
		jwtSecret: getJWTSecrt(),
		baseConfig: &baseConfig{
			dbDNS: getDBDNS(),
		},
		exposeVersion:   getExposeVersion(),
		credentialsFile: getCredentialsFilePath(getDir()),
	}
}

func LoadGRPCConfig() *GRPCConfig {
	return &GRPCConfig{
		port: getGRPCPort(),
		baseConfig: &baseConfig{
			dbDNS: getDBDNS(),
		},
	}
}

func LoadTaskEngineConfig() *TaskEngineConfig {
	return &TaskEngineConfig{
		baseConfig: &baseConfig{
			dbDNS: getTaskEngineDBDNS(),
		},
	}
}
