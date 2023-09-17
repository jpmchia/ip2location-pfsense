package config

import (
	"fmt"
	"os"
	"time"

	. "ip2location-pfsense/util"

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
	LogDebug("Initialising configuration")

	defaultConfig = initViperConfig(appName)

	SetConfigFile(CfgFile)

	err := viper.Unmarshal(&Config)
	HandleFatalError(err, "Unable to unmarshal configuration:\n")
}

// Config returns a default config provider
func ConfigProvider() Provider {
	return defaultConfig
}
func GetConfig() *viper.Viper {
	return defaultConfig
}

func LoadConfiguration() (Options, error) {
	viper := initViperConfig(appName)
	SetConfigFile(CfgFile)
	err := viper.Unmarshal(&Config)
	HandleFatalError(err, "Unable to unmarshal configuration:\n")
	return Config, err
}

// LoadConfigProvider returns a configured viper instance
func LoadConfigProvider(appName string) Provider {
	return initViperConfig(appName)
}

func SetValue(key string, value interface{}) {
	defaultConfig.Set(key, value)
}

// SetConfigFile sets the config file to use and reloads the configuraton
func SetConfigFile(file string) {
	CfgFile = file

	setConfigLocations()

	// If a config file is found, read it in.
	err := defaultConfig.ReadInConfig()
	HandleFatalError(err, "Unable to read configuration:\n")
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
	v.SetDefault("service.bind_host", "127.0.0.1")
	v.SetDefault("service.bind_port", "9999")
	v.SetDefault("ip2api.url", "https://api.ip2location.io/")
	v.SetDefault("ip2api.key", "4C9057DC19ADD10E8CEDF123C74C77CE")
	v.SetDefault("ip2api.plan", "Free")
	// v.SetDefault("counters.limits.monthly", "30000")
	// v.SetDefault("counters.limits.daily", "900")
	// v.SetDefault("counters.limits.hourly", "60")
	// v.SetDefault("counters.startdate", "")
	// v.SetDefault("counters.nextreset", "")
	// v.SetDefault("counters.lifetime", "")
	// v.SetDefault("counters.enabled", true)
	v.SetDefault("use_cache", true)
	v.SetDefault("debug", true)

	return v
}

// SetConfigLocations sets the locations to search for the configuration file
func setConfigLocations() {
	LogDebug("Setting config locations")
	// Use config file from the flag.
	defaultConfig.SetConfigFile(CfgFile)
	defaultConfig.AddConfigPath(".")
	defaultConfig.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	defaultConfig.AddConfigPath(fmt.Sprintf("/usr/local/%s", appName))
	defaultConfig.AddConfigPath(fmt.Sprintf("/usr/local/etc/%s", appName))
	defaultConfig.AddConfigPath(fmt.Sprintf("/opt/%s", appName))
	home, err := os.UserHomeDir()
	HandleError(err, "Unable to determine user's home directory")
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
	HandleFatalError("Unable to write configuration:\n", err.Error())
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

	LogDebug("Creating configuration file: %s", CfgFile)

	// If a config file is found, read it in.
	err := defaultConfig.SafeWriteConfigAs(CfgFile)
	HandleFatalError(err, "Unable to write configuration:\n")

	fmt.Printf("Configuration file created: %s", CfgFile)
}
