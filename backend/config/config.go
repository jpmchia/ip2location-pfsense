package config

import (
	"fmt"
	"os"
	"time"

	"ip2location-pfsense/util"

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

var Config Options

const appName string = "IP2Location-pfSense"

var CfgFile string = "config.yaml"
var defaultConfig *viper.Viper

// Default initialiser for the applicaiton's configuration
func init() {

	util.LogDebug("Initialising configuration")

	Configure()
}

func Configure() {
	// Load the default configuration
	defaultConfig = initViperConfig(appName)

	// Set and load additional configuration locations
	setConfigLocations(CfgFile)

	// Read the configuration file
	err := defaultConfig.ReadInConfig()
	util.HandleError(err, "Unable to read configuration:\n")

	// Unmarshal the configuration into the Config struct
	err = defaultConfig.Unmarshal(&Config)
	util.HandleError(err, "Unable to unmarshal configuration:\n")
}

// Config returns a default config provider
func ConfigProvider() *viper.Viper {
	return defaultConfig
}

func GetConfiguration() Options {
	err := defaultConfig.Unmarshal(&Config)
	util.HandleError(err, "Unable to unmarshal configuration:\n")
	return Config
}

// func LoadConfiguration() (Options, error) {
// 	defaultConfig = initViperConfig(appName)
// 	setConfigLocations(CfgFile)
// 	err := defaultConfig.ReadInConfig()
// 	if err != nil {
// 		util.HandleError(err, "Unable to read configuration:\n")
// 		return Config, err
// 	}
// 	util.LogDebug("Using configuration file: %s", defaultConfig.ConfigFileUsed())
// 	err = defaultConfig.Unmarshal(&Config)
// 	util.HandleError(err, "Unable to unmarshal configuration:\n")
// 	return Config, err
// }

// LoadConfigProvider returns a configured viper instance
// func LoadConfigProvider(appName string) Provider {
// 	return initViperConfig(appName)
// }

func SetValue(key string, value interface{}) {
	defaultConfig.Set(key, value)
}

// SetConfigFile sets the config file to use and reloads the configuraton
func SetConfigFile(file string) {
	setConfigLocations(file)
	err := defaultConfig.ReadInConfig()
	util.HandleError(err, "Unable to read configuration:\n")
}

// Initialises viper with default values
func initViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults
	v.SetDefault("jsonlogs", false)
	v.SetDefault("loglevel", "debug")
	v.SetDefault("installation_path", "/usr/local/ip2location")
	v.SetDefault("use_redis", true)
	v.SetDefault("redis.ip2location.host", "127.0.0.1")
	v.SetDefault("redis.ip2location.port", "6379")
	v.SetDefault("redis.ip2location.db", 1)
	v.SetDefault("redis.ip2location.auth", "ip2location")
	v.SetDefault("redis.ip2location.pass", "password")
	v.SetDefault("redis.pfsense.host", "127.0.0.1")
	v.SetDefault("redis.pfsense.port", "6379")
	v.SetDefault("redis.pfsense.db", 2)
	v.SetDefault("redis.pfsense.auth", "ip2location")
	v.SetDefault("redis.pfsense.pass", "password")

	v.SetDefault("ip2api.url", "https://api.ip2location.io/")
	v.SetDefault("ip2api.key", "")
	v.SetDefault("ip2api.plan", "Free")
	v.SetDefault("ip2api.max_errors", 5)

	v.SetDefault("service.bind_host", "127.0.0.1")
	v.SetDefault("service.bind_port", "9999")
	v.SetDefault("service.use_ssl", false)
	v.SetDefault("service.ssl_cert", "cert.pem")
	v.SetDefault("service.ssl_key", "cert.key")
	v.SetDefault("service.ingest_logs", "/api/filterlog")
	v.SetDefault("service.ip2l_results", "/api/results")
	v.SetDefault("service.ip2geomap", "/index.html")
	v.SetDefault("service.healthcheck", "/health")
	v.SetDefault("service.ip_requests", "/api/ip2location")

	v.SetDefault("use_cache", true)
	v.SetDefault("debug", false)

	return v
}

// SetConfigLocations sets the locations to search for the configuration file
func setConfigLocations(file string) {
	util.LogDebug("Setting config locations")
	CfgFile = file
	// Use config file from the flag.
	defaultConfig.SetConfigType("yaml")
	defaultConfig.SetConfigFile(CfgFile)
	defaultConfig.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	defaultConfig.AddConfigPath(fmt.Sprintf("/usr/local/%s", appName))
	defaultConfig.AddConfigPath(fmt.Sprintf("/usr/local/etc/%s", appName))
	defaultConfig.AddConfigPath(fmt.Sprintf("/opt/%s", appName))
	home, err := os.UserHomeDir()
	util.HandleError(err, "Unable to determine user's home directory")
	defaultConfig.AddConfigPath(fmt.Sprintf("%s/.%s", home, appName))
	defaultConfig.AddConfigPath(fmt.Sprintf("%s/.config/%s", home, appName))
	defaultConfig.AddConfigPath(home)
}

// PrintConfig prints the configuration to stdout
func printConfig(v *viper.Viper) {
	for _, k := range v.AllKeys() {
		fmt.Printf("%s = %v\n", k, v.Get(k))
	}
}

// WriteConfigValue writes a value to the configuration file
func WriteConfigValue(key string, value any) {
	defaultConfig.Set(key, value)

	err := defaultConfig.WriteConfig()
	util.HandleFatalError("Unable to write configuration:\n", err.Error())
}

// ShowConfig prints the configuration to stdout
func ShowConfig() {
	fmt.Println("Configuration:")
	fmt.Println("-------------")
	fmt.Printf("Configuration file: %s\n", defaultConfig.ConfigFileUsed())
	fmt.Println("-------------")

	printConfig(defaultConfig)

	os.Exit(0)
}

// CreateConfigFile creates a new configuration file with the default values and exits,
// if no filename is specified, the default filename is used
func CreateConfigFile(args []string) {
	if len(args) == 0 {
		fmt.Printf("No filename specified. Using the default file: %s", CfgFile)
	}
	if len(args) > 1 {
		fmt.Printf("Too many arguments specified. Using the default file: %s", CfgFile)
	}
	if len(args) == 1 {
		CfgFile = args[0]
	}

	util.LogDebug("Creating configuration file: %s", CfgFile)

	// If a config file is found, read it in.
	err := defaultConfig.SafeWriteConfigAs(CfgFile)
	util.HandleFatalError(err, "Unable to write configuration:\n")

	fmt.Printf("Configuration file created: %s", CfgFile)
}
