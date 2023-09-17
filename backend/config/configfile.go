package config

type Options struct {
	InstallationPath string                  `mapstructure:"installation_path"`
	IP2API           Ip2ApiOptions           `mapstructure:"ip2api"`
	Redis            map[string]RedisOptions `mapstructure:"redis"`
	JsonLogs         bool                    `mapstructure:"jsonlogs"`
	LogLevel         string                  `mapstructure:"loglevel"`
	UseRedis         bool                    `mapstructure:"use_redis"`
}

type RedisOptions struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Db   int    `mapstructure:"db"`
	Auth string `mapstructure:"auth"`
	Pass string `mapstructure:"pass"`
}

type Ip2ApiOptions struct {
	URL  string `mapstructure:"url"`
	Key  string `mapstructure:"key"`
	Auth string `mapstructure:"auth"`
	Plan string `mapstructure:"plan"`
}
