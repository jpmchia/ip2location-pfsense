package config

type (
	Options struct {
		InstallationPath string                  `mapstructure:"installation_path"`
		IP2API           Ip2ApiOptions           `mapstructure:"ip2api"`
		Redis            map[string]RedisOptions `mapstructure:"redis"`
		JsonLogs         bool                    `mapstructure:"jsonlogs"`
		LogLevel         string                  `mapstructure:"loglevel"`
		UseRedis         bool                    `mapstructure:"use_redis"`
		Service          ServiceOptions          `mapstructure:"service"`
	}

	RedisOptions struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
		Db   int    `mapstructure:"db"`
		Auth string `mapstructure:"auth"`
		Pass string `mapstructure:"pass"`
	}

	Ip2ApiOptions struct {
		URL       string `mapstructure:"url"`
		Key       string `mapstructure:"key"`
		Plan      string `mapstructure:"plan"`
		MaxErrors int    `mapstructure:"max_errors"`
	}

	ServiceOptions struct {
		BindHost string `mapstructure:"bind_host"`
		BindPort string `mapstructure:"bind_port"`
		// CORS hosts to allow
		AllowHosts []string `mapstructure:"allow_hosts"`
		SSLCert    string   `mapstructure:"ssl_cert"`
		SSLKey     string   `mapstructure:"ssl_key"`
		UseSSL     bool     `mapstructure:"use_ssl"`
		// Filterlogs
		IngestLogs string `mapstructure:"ingest_logs"`
		Results    string `mapstructure:"ip2l_results"`
		// Content
		EnableWeb    bool     `mapstructure:"enable_web"`
		HomePage     string   `mapstructure:"homepage"`
		DetailPage   string   `mapstructure:"ip2geomap"`
		IPRequests   string   `mapstructure:"ip_requests"`
		HealthCheck  string   `mapstructure:"healthcheck"`
		BearerTokens []string `mapstructure:"service.bearertokens"`
		// API Keys
		ApiTimeout int `mapstructure:"api_timeout"`
	}
)
