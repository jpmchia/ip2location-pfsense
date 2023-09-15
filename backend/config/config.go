package config

import (
	"time"

	"github.com/spf13/viper"
)

type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var defaultConfig *viper.Viper

// Config returns a default config providers
func Config() Provider {
	return defaultConfig
}

// LoadConfigProvider returns a configured viper instance
func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func WriteConfigString(key string, value string) {
	defaultConfig.Set(key, value)
}

func SetConfigLocations() {
	defaultConfig.SetConfigName("config")
	defaultConfig.AddConfigPath("/usr/local/ip2location")
	defaultConfig.AddConfigPath(".")
	defaultConfig.AddConfigPath("/etc/ip2location")
	defaultConfig.AddConfigPath("/usr/local/etc/ip2location")
	defaultConfig.AddConfigPath("$HOME/.ip2location")
	defaultConfig.AddConfigPath("$HOME/.config/ip2location")
}

func init() {
	defaultConfig = readViperConfig("IP2LOCATION")
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults
	v.SetDefault("jsonlogs", false)
	v.SetDefault("loglevel", "debug")
	v.SetDefault("installation_path", "/usr/local/ip2location")
	v.SetDefault("use_redis", true)
	v.SetDefault("redis_host", "127.0.0.1")
	v.SetDefault("redis_port", "6379")
	v.SetDefault("redis_db_ip_log", 1)
	v.SetDefault("redis_db_location", 2)
	v.SetDefault("redis_auth", "ip2location")
	v.SetDefault("redis_pass", "password")
	v.SetDefault("bind_host", "192.168.0.51")
	v.SetDefault("bind_port", "9999")
	v.SetDefault("ip2api_url", "https://api.ip2location.io/")
	v.SetDefault("ip2api_key", "")
	v.SetDefault("ip2api_plan", "Free")
	v.SetDefault("use_cache", true)
	return v
}
