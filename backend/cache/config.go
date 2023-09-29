package cache

import (
	"github.com/jpmchia/ip2location-pfsense/backend/config"
	"github.com/jpmchia/ip2location-pfsense/backend/util"

	"github.com/spf13/viper"
)

type RedisCacheConfig struct {
	HostPort string `mapstructure:"host"`
	Db       int    `mapstructure:"db"`
	Auth     string `mapstructure:"auth"`
	Pass     string `mapstructure:"pass"`
	Type     string `mapstructure:"type"`
}

const RedisConfigKey = "redis"

var configuration *viper.Viper

func LoadConfiguration(subkey string) (RedisCacheConfig, error) {

	config.Configure()

	configuration = config.ConfigProvider()

	util.LogDebug("[cache] Loading configuration for %s", subkey)
	util.LogDebug("[cache] Configuration: %v", configuration.AllSettings())

	subconfig := configuration.Sub(subkey)
	if subconfig == nil {
		util.HandleFatalError(nil, "[cache] Unable to find configuration for %s", subkey)
	}

	redisConfig, err := loadConfig(subconfig)
	util.HandleFatalError(err, "[cache] Unable to load configuration for %s", subkey)

	return *redisConfig, err
}

func loadConfig(v *viper.Viper) (conf *RedisCacheConfig, err error) {
	return &RedisCacheConfig{
		HostPort: v.GetString("host") + ":" + v.GetString("port"),
		Db:       v.GetInt("db"),
		Auth:     v.GetString("auth"),
		Pass:     v.GetString("pass"),
		Type:     v.GetString("type"),
	}, err
}
