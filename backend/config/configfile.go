package config

type Options struct {
	InstallationPath string                  `mapstructure:"installation_path"`
	IP2API           Ip2ApiOptions           `mapstructure:"ip2api"`
	Redis            map[string]RedisOptions `mapstructure:"redis"`
	JsonLogs         bool                    `mapstructure:"jsonlogs"`
	LogLevel         string                  `mapstructure:"loglevel"`
	UseRedis         bool                    `mapstructure:"use_redis"`
	Service          ServiceOptions          `mapstructure:"service"`
}

type RedisOptions struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Db   int    `mapstructure:"db"`
	Auth string `mapstructure:"auth"`
	Pass string `mapstructure:"pass"`
}

type Ip2ApiOptions struct {
	URL       string `mapstructure:"url"`
	Key       string `mapstructure:"key"`
	Plan      string `mapstructure:"plan"`
	MaxErrors int    `mapstructure:"max_errors"`
}

type APIKeys map[string]string

type ServiceOptions struct {
	BindHost    string `mapstructure:"bind_host"`
	BindPort    string `mapstructure:"bind_port"`
	SSLCert     string `mapstructure:"ssl_cert"`
	SSLKey      string `mapstructure:"ssl_key"`
	UseSSL      bool   `mapstructure:"use_ssl"`
	IngestLogs  string `mapstructure:"ingest_logs"`
	Results     string `mapstructure:"ip2l_results"`
	DetailPage  string `mapstructure:"ip2geomap"`
	IPRequests  string `mapstructure:"ip_requests"`
	HealthCheck string `mapstructure:"healthcheck"`
}
