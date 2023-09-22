package ip2location

import (
	"fmt"
	. "ip2location-pfsense/util"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type CounterConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	StartDate  string `mapstructure:"startdate"`  // The date the counters were started
	LastReset  string `mapstructure:"lastreset"`  // The date the counters were last reset
	ResetLimit string `mapstructure:"limitreset"` // The date and time the hourly and daily limits will be reset
	LastCheck  string `mapstructure:"lastcheck"`  // The last time we checked the counters and limits
	NextReset  string `mapstructure:"nextreset"`  // The date the monthly counters will be reset
	Limits     struct {
		Monthly string `mapstructure:"monthly"`
		Daily   string `mapstructure:"daily"`
		Hourly  string `mapstructure:"hourly"`
	} `mapstructure:"limits"`
	Lifetime string `mapstructure:"lifetime"`
	Action   string `mapstructure:"action"`
}

const appName string = "IP2Location-pfSense"
const appNameCounters string = "IP2Location-pfSense-Counters"

var localFile string = "counters.yaml"
var CounterValues CounterConfig
var Counters *viper.Viper

func init() {
	LogDebug("Initialising counters")

	Counters = initViperCounters(appNameCounters)
	err := viper.Unmarshal(&CounterValues)

	HandleFatalError(err, "Unable to unmarshal counter values:\n")
}

func initViperCounters(appName string) *viper.Viper {
	v := viper.New()

	v.SetDefault("counters.limits.monthly", "30000")
	v.SetDefault("counters.limits.daily", "900")
	v.SetDefault("counters.limits.hourly", "60")
	v.SetDefault("counters.startdate", "")
	v.SetDefault("counters.lastreset", "")
	v.SetDefault("counters.lastcheck", "")
	v.SetDefault("counters.nextreset", "")
	v.SetDefault("counters.lifetime", "")
	v.SetDefault("counters.count", 0)
	v.SetDefault("counters.dailycount", 0)
	v.SetDefault("counters.lifetime", 0)
	v.SetDefault("counters.enabled", true)

	v = setConfigLocations(v)

	return v
}

func LoadCounters() (CounterConfig, error) {
	Counters = initViperCounters(appNameCounters)

	err := Counters.ReadInConfig()
	HandleError(err, "Unable to read counters file: %v\n", err.Error())

	err = Counters.Unmarshal(&CounterValues)
	HandleFatalError(err, "Unable to unmarshal counter values:\n")

	return CounterValues, err
}

func InitailiseCounters(force bool) {
	Counters = initViperCounters(appNameCounters)

	err := Counters.ReadInConfig()
	if err == nil {
		log.Printf("Counter values found: %s", localFile)
		ShowCounters()
		if !force {
			log.Printf("\nTo reset counters, use the --force flag")
			os.Exit(0)
		}
	} else {
		LogDebug("No counters file found. Creating a new one.")
	}

	Counters.Set("counters.limits.monthly", "30000")
	Counters.Set("counters.limits.daily", "900")
	Counters.Set("counters.limits.hourly", "60")
	Counters.Set("counters.startdate", time.Now().UTC().Format(time.RFC3339))
	Counters.Set("counters.lastreset", time.Now().UTC().Format(time.RFC3339))
	Counters.Set("counters.nextreset", time.Now().UTC().AddDate(0, 1, 0).Format(time.RFC3339))
	Counters.Set("counters.lastcheck", time.Now().UTC().AddDate(0, 1, 0).Format(time.RFC3339))
	Counters.Set("counters.count", 0)
	Counters.Set("counters.dailycount", 0)
	Counters.Set("counters.lifetime", 0)
	Counters.Set("counters.enabled", true)

	log.Printf("Creating the local counters file: %s", localFile)

	err = Counters.SafeWriteConfigAs(localFile)
	HandleFatalError(err, "Unable to create counters file: %v\n", err.Error())

	err = Counters.Unmarshal(&CounterValues)
	HandleFatalError(err, "Unable to unmarshal counter values:\n")

	LogDebug("Reinitialised counters")
}

func CreateCountersFile(args []string) {
	if len(args) == 0 {
		fmt.Printf("No filename specified. Using the default file: %s", localFile)
	}
	if len(args) > 1 {
		fmt.Printf("Too many arguments specified. Using the default file: %s", localFile)
	}
	if len(args) == 1 {
		localFile = args[0]
	}

	LogDebug("Creating the local counters file: %s", localFile)

	err := Counters.SafeWriteConfigAs(localFile)
	HandleFatalError(err, "Unable to write configuration:\n")

	fmt.Printf("Configuration file created: %s", localFile)
}

func setConfigLocations(v *viper.Viper) *viper.Viper {
	LogDebug("Setting config locations")
	// Use config file from the flag.
	v.SetConfigFile(localFile)
	v.AddConfigPath(".")
	v.AddConfigPath(fmt.Sprintf("/opt/%s", appName))
	v.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	v.AddConfigPath(fmt.Sprintf("/usr/local/%s", appName))
	v.AddConfigPath(fmt.Sprintf("/usr/local/etc/%s", appName))
	home, err := os.UserHomeDir()
	HandleError(err, "Unable to determine user's home directory")
	v.AddConfigPath(fmt.Sprintf("%s/.%s", home, appName))
	v.AddConfigPath(fmt.Sprintf("%s/.config/%s", home, appName))
	v.AddConfigPath(home)

	return v
}

func printCounters(v *viper.Viper) {
	for _, k := range v.AllKeys() {
		fmt.Printf("%s = %v\n", k, v.Get(k))
	}
}

func ShowCounters() {
	fmt.Printf("Counters file: %s\n", localFile)
	fmt.Println("-------------")

	printCounters(Counters)

	os.Exit(0)
}

func UpdateCounters(addIncrements int) int {
	Counters.Set("counters.count", Counters.GetInt("counters.count")+addIncrements)
	Counters.Set("counters.dailycount", Counters.GetInt("counters.dailycount")+addIncrements)
	Counters.Set("counters.lifetime", Counters.GetInt("counters.lifetime")+addIncrements)

	err := Counters.WriteConfig()
	HandleFatalError(err, "Unable to write configuration:\n")

	err = Counters.Unmarshal(&CounterValues)
	HandleFatalError(err, "Unable to unmarshal counter values:\n")

	return Counters.GetInt("counters.count")
}

func GetRemainingCount() (int, int, int) {
	return Counters.GetInt("counters.limits.monthly") - Counters.GetInt("counters.count"),
		Counters.GetInt("counters.limits.daily") - Counters.GetInt("counters.dailycount"),
		Counters.GetInt("counters.limits.hourly") - Counters.GetInt("counters.hourlycount")
}

func CheckCounters() bool {
	var dateTimeNow time.Time = time.Now().UTC()
	var dateTimeLastReset time.Time = Counters.GetTime("counters.lastreset")
	var dateTimeNextReset time.Time = Counters.GetTime("counters.nextreset")
	var dateTimeLastCheck time.Time = Counters.GetTime("counters.lastcheck")

	// Monthly counter
	if dateTimeNow.After(dateTimeNextReset) {
		LogDebug("Resetting counters")

		Counters.Set("counters.lastreset", dateTimeNow.Format(time.RFC3339))
		Counters.Set("counters.nextreset", dateTimeNow.AddDate(0, 1, 0).Format(time.RFC3339))
		Counters.Set("counters.dailycount", 0)
		Counters.Set("counters.hourlycount", 0)
		Counters.Set("counters.lastcheck", dateTimeNow.Format(time.RFC3339))

		err := Counters.WriteConfig()
		HandleFatalError(err, "Unable to write configuration:\n")

		err = Counters.Unmarshal(&CounterValues)
		HandleFatalError(err, "Unable to unmarshal counter values:\n")

		return true
	}

	// Daily counter, reset if the last reset was more than 24 hours ago
	if dateTimeNow.After(dateTimeLastReset.AddDate(0, 0, 1)) {
		LogDebug("Resetting daily counters")

		Counters.Set("counters.lastreset", dateTimeNow.Format(time.RFC3339))
		Counters.Set("counters.dailycount", 0)
		Counters.Set("counters.hourlycount", 0)
		Counters.Set("counters.lastcheck", dateTimeNow.Format(time.RFC3339))
	}

	// Instead of reseting the hourly counter, we just check what time the last reset was and we check how many hows until the 24 hours passes, then we adjust the hourly limit accordingly
	if dateTimeNow.After(dateTimeLastCheck.Add(time.Hour)) {
		LogDebug("Resetting hourly counters")

		Counters.Set("counters.lastcheck", dateTimeNow.Format(time.RFC3339))

		// Calculate the number of hours until the next reset
		var hoursUntilReset int = int(dateTimeNextReset.Sub(dateTimeNow).Hours())
		var hourlyLimit int = Counters.GetInt("counters.limits.hourly")

		// Calcualte tehe number of increments remaining until the daily limit is reached
		var incrementsRemaining int = Counters.GetInt("counters.limits.daily") - Counters.GetInt("counters.dailycount")

		// Apportion the remaining increments to the number of hours until the next reset
		var hourlyIncrements int = incrementsRemaining / hoursUntilReset

		// If the hourly increments is less than the hourly limit, then we adjust the hourly limit accordingly
		if hourlyIncrements < hourlyLimit {
			Counters.Set("counters.limits.hourly", hourlyIncrements)
		}

		// If the number of hours until the next reset is equal to the hourly limit, then we adjust the hourly limit accordingly
		if hoursUntilReset == 0 {
			Counters.Set("counters.limits.hourly", 0)
			Counters.Set("counters.hourlycount", 0)
		}
	}

	err := Counters.WriteConfig()
	HandleFatalError(err, "Unable to write configuration:\n")

	err = Counters.Unmarshal(&CounterValues)
	HandleFatalError(err, "Unable to unmarshal counter values:\n")

	return false
}
