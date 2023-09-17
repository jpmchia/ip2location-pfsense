package cache

import (
	"ip2location-pfsense/config"
	. "ip2location-pfsense/util"

	"github.com/spf13/viper"
)

type RedisCacheConfig struct {
	HostPort string `mapstructure:"host"`
	Db       int    `mapstructure:"db"`
	Auth     string `mapstructure:"auth"`
	Pass     string `mapstructure:"pass"`
}

const RedisConfigKey = "redis"
const ip2LocationKey = "redis.ip2location"
const pfSenseKey = "redis.pfsense"

var configuration *viper.Viper

func LoadConfiguration(subkey string) (RedisCacheConfig, error) {
	_, err := config.LoadConfiguration()
	HandleFatalError(err, "Unable to load configuration")

	// redisConfigKey := options.Redis.Key
	configuration = config.GetConfig()
	LogDebug("Loading configuration for %s", subkey)
	LogDebug("Configuration: %v", configuration.AllSettings())

	subconfig := configuration.Sub(subkey)
	if subconfig == nil {
		HandleFatalError(nil, "Unable to find configuration for %s", subkey)
	}

	redisConfig, err := loadConfig(subconfig)
	HandleFatalError(err, "Unable to load configuration for %s", subkey)

	return *redisConfig, err
}

func loadConfig(v *viper.Viper) (conf *RedisCacheConfig, err error) {
	return &RedisCacheConfig{
		HostPort: v.GetString("host") + ":" + v.GetString("port"),
		Db:       v.GetInt("db"),
		Auth:     v.GetString("auth"),
		Pass:     v.GetString("pass"),
	}, err
}
