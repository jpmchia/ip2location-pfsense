package ip2location

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jpmchia/ip2location-pfsense/backend/util"

	"github.com/spf13/viper"
)

type Counter struct {
	Max       int64
	Count     int64
	StartDate time.Time
	NextReset time.Time
}

type ActiveCounters struct {
	Monthly     Counter
	Daily       Counter
	Hourly      Counter
	Lifetime    Counter
	CacheHits   int64
	CacheMisses int64
	ApiMisses   int64
	Enabled     bool
}

type CounterValues struct {
	ApiCalls          int
	ApiErrors         int
	CacheHits         int
	CacheMisses       int
	ApiCallsRemaining int
	HourlyLimit       int
}

var appName string = "github.com/jpmchia/ip2location-pfsense/backend"
var appNameCounters string = "github.com/jpmchia/ip2location-pfsense/backend-counters"

var localFile string = "counters.yaml"
var Counters ActiveCounters
var CountersProvider *viper.Viper

func init() {
	util.LogDebug("Initialising counters %s", appName)

	CountersProvider = initViperCounters(appNameCounters)

	counters, err := readCounterValues(localFile, CountersProvider)
	util.HandleError(err, "Unable to read counters file.")

	Counters = *counters
}

func initViperCounters(appName string) *viper.Viper {
	v := viper.New()

	v.SetDefault("counters.monthly.max", 30000)
	v.SetDefault("counters.monthly.count", 0)
	v.SetDefault("counters.monthly.startdate", time.Now().UTC())
	v.SetDefault("counters.monthly.nextreset", time.Now().UTC().AddDate(0, 1, 0))

	v.SetDefault("counters.daily.max", 900)
	v.SetDefault("counters.daily.count", 0)
	v.SetDefault("counters.daily.startdate", time.Now().UTC())
	v.SetDefault("counters.daily.nextreset", time.Now().UTC().AddDate(0, 0, 1))

	v.SetDefault("counters.hourly.max", 900)
	v.SetDefault("counters.hourly.count", 0)
	v.SetDefault("counters.hourly.startdate", time.Now().UTC())
	v.SetDefault("counters.hourly.nextreset", time.Now().UTC().Add(time.Hour))

	v.SetDefault("counters.lifetime.max", 0)
	v.SetDefault("counters.lifetime.count", 0)
	v.SetDefault("counters.lifetime.startdate", time.Now().UTC())
	v.SetDefault("counters.lifetime.nextreset", time.Now().UTC())
	v.SetDefault("counters.apimisses", 0)
	v.SetDefault("counters.cachehits", 0)
	v.SetDefault("counters.cachemisses", 0)
	v.SetDefault("counters.enabled", true)
	v = setCounterLocations(v)
	return v
}

func setCounterLocations(v *viper.Viper) *viper.Viper {

	util.LogDebug("Setting config locations")
	// Use config file from the flag.
	v.SetConfigFile(localFile)
	v.AddConfigPath(".")
	v.AddConfigPath(fmt.Sprintf("/opt/%s", appName))
	v.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	v.AddConfigPath(fmt.Sprintf("/usr/local/%s", appName))
	v.AddConfigPath(fmt.Sprintf("/usr/local/etc/%s", appName))
	home, err := os.UserHomeDir()
	util.HandleError(err, "Unable to determine user's home directory")
	v.AddConfigPath(fmt.Sprintf("%s/.%s", home, appName))
	v.AddConfigPath(fmt.Sprintf("%s/.config/%s", home, appName))
	v.AddConfigPath(home)

	return v
}

func readCounterValues(filename string, v *viper.Viper) (*ActiveCounters, error) {

	v.SetConfigFile(filename)
	err := v.ReadInConfig()
	if err != nil {
		util.HandleError(err, "Unable to read counters file: %v\n", err.Error())
		// Create a new counters file
		InitialiseCounters(filename, 30000, 900, 900, false)
		// Read the new counters file
		err = v.ReadInConfig()
		util.HandleError(err, "Unable to read counters file: %v\n", err.Error())
	}

	err = v.Unmarshal(&Counters)
	util.HandleFatalError(err, "Unable to unmarshal counter values:\n")

	return &Counters, err
}

func InitialiseCounters(filename string, monthly int, daily int, hourly int, force bool) {

	viper := initViperCounters(appNameCounters)

	err := viper.ReadInConfig()
	if err == nil {
		log.Printf("Counter values found: %s", localFile)
		ShowCounters()
		if !force {
			log.Printf("\nTo reset counters, use the --force flag")
			os.Exit(0)
		}
	} else {
		util.LogDebug("No counters file found. Creating a new one.")
	}

	CountersProvider.Set("counters.monthly.max", monthly)
	CountersProvider.Set("counters.monthly.count", 0)
	CountersProvider.Set("counters.monthly.startdate", time.Now().UTC().Format(time.RFC3339))
	CountersProvider.Set("counters.monthly.nextreset", time.Now().UTC().AddDate(0, 1, 0).Format(time.RFC3339))
	CountersProvider.Set("counters.daily.max", daily)
	CountersProvider.Set("counters.daily.count", 0)
	CountersProvider.Set("counters.daily.startdate", time.Now().UTC().Format(time.RFC3339))
	CountersProvider.Set("counters.daily.nextreset", time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339))
	CountersProvider.Set("counters.hourly.max", hourly)
	CountersProvider.Set("counters.hourly.count", 0)
	CountersProvider.Set("counters.hourly.startdate", time.Now().UTC())
	CountersProvider.Set("counters.hourly.nextreset", time.Now().UTC().Add(time.Hour))
	CountersProvider.Set("counters.lifetime.max", 0)
	CountersProvider.Set("counters.lifetime.count", 0)
	CountersProvider.Set("counters.lifetime.startdate", time.Now().UTC())
	CountersProvider.Set("counters.lifetime.nextreset", time.Now().UTC())
	CountersProvider.Set("counters.apimisses", 0)
	CountersProvider.Set("counters.cachehits", 0)
	CountersProvider.Set("counters.cachemisses", 0)
	CountersProvider.Set("counters.enabled", true)

	// CountersProvider.Set("counters.startdate", time.Now().UTC().Format(time.RFC3339))
	// CountersProvider.Set("counters.lastreset", time.Now().UTC().Format(time.RFC3339))
	// CountersProvider.Set("counters.nextreset", time.Now().UTC().AddDate(0, 1, 0).Format(time.RFC3339))
	// CountersProvider.Set("counters.lastcheck", time.Now().UTC().AddDate(0, 1, 0).Format(time.RFC3339))
	// CountersProvider.Set("counters.count", 0)
	// CountersProvider.Set("counters.dailycount", 0)
	// CountersProvider.Set("counters.lifetime", 0)
	// CountersProvider.Set("counters.enabled", true)

	log.Printf("Creating the local counters file: %s", localFile)

	err = CountersProvider.SafeWriteConfigAs(localFile)
	util.HandleFatalError(err, "Unable to create counters file: %v\n", err.Error())

	err = CountersProvider.Unmarshal(&Counters)
	util.HandleFatalError(err, "Unable to unmarshal counter values:\n")

	util.LogDebug("Reinitialised counters")
}

func createCountersFile(filename string) {

	// TODO: Read in arguments for the monthly, daily and hourly limits
	util.LogDebug("Creating the local counters file: %s", localFile)

	err := CountersProvider.SafeWriteConfigAs(localFile)
	util.HandleFatalError(err, "Unable to write configuration:\n")

	fmt.Printf("Configuration file created: %s", localFile)
}

// Create a new counters file - --force, --filename <filename>, --monthly <max>, --daily <max>, --hourly <max>
func CreateCountersFile(args []string) {
	var filename string
	// var force bool = false

	if len(args) == 0 {
		filename = localFile
		fmt.Printf("No filename specified. Using the default file: %s", localFile)
	}
	if len(args) > 1 {
		// if args[0] == "--force" {
		// 	force = true
		// }
		if args[0] == "--filename" {
			filename = args[1]
		}
	}

	var provider = initViperCounters(appNameCounters)
	// Load existing
	_, err := readCounterValues(filename, provider)
	util.HandleError(err, "Unable to read counters file:\n")

	createCountersFile(filename)
}

func printCounters(v *viper.Viper) {
	for _, k := range v.AllKeys() {
		fmt.Printf("%s = %v\n", k, v.Get(k))
	}
}

func ShowCounters() {
	fmt.Printf("Counters file: %s\n", localFile)
	fmt.Println("-------------")

	printCounters(CountersProvider)

	os.Exit(0)
}

func IncrementCounters(api int, apimiss int, cache int, cachemiss int) error {

	if api != 0 {
		Counters.Lifetime.Count = Counters.Lifetime.Count + int64(api)
		Counters.Daily.Count = Counters.Daily.Count + int64(api)
		Counters.Hourly.Count = Counters.Hourly.Count + int64(api)
		Counters.Monthly.Count = Counters.Monthly.Count + int64(api)
	}
	if cache != 0 {
		Counters.CacheHits = Counters.CacheHits + int64(cache)
		Counters.ApiMisses = Counters.ApiMisses + int64(apimiss)
		Counters.CacheMisses = Counters.CacheMisses + int64(cache)
	}

	return nil
}

func WriteBackCounters(values CounterValues, andWrite bool) {

	CountersProvider.Set("counters.monthly.count", Counters.Monthly.Count+int64(values.ApiCalls))
	CountersProvider.Set("counters.daily.count", Counters.Daily.Count+int64(values.ApiCalls))
	CountersProvider.Set("counters.hourly.count", Counters.Hourly.Count+int64(values.ApiCalls))
	CountersProvider.Set("counters.lifetime.count", Counters.Lifetime.Count+int64(values.ApiCalls))
	CountersProvider.Set("counters.cachemisses", Counters.CacheHits+int64(values.CacheMisses))
	CountersProvider.Set("counters.cachehits", Counters.CacheHits+int64(values.CacheHits))
	CountersProvider.Set("counters.apimisses", Counters.CacheMisses+int64(values.ApiErrors))

	if andWrite {
		err := CountersProvider.WriteConfig()
		util.HandleFatalError(err, "Unable to write configuration:\n")

		err = CountersProvider.Unmarshal(&Counters)
		util.HandleFatalError(err, "Unable to unmarshal counter values:\n")
	}
}

func GetRemainingThisDay() int {

	if time.Now().UTC().After(Counters.Daily.NextReset) {
		Counters.Daily.Count = 0
		Counters.Daily.NextReset = time.Now().UTC().AddDate(0, 0, 1)
		go WriteBackCounters(CounterValues{}, true)
		util.HandleError(nil, "Unable to write back counters")
	}
	if Counters.Daily.Max == 0 {
		return -1
	}
	if Counters.Daily.Count >= Counters.Daily.Max {
		return 0
	}

	if Counters.Daily.Count < Counters.Daily.Max {
		return int(Counters.Daily.Max - Counters.Daily.Count)
	}
	return 0
}

func CalculateNewHourlyMax() int {
	// Calculate the number of hours until the next reset
	var hoursUntilReset int = int(Counters.Daily.NextReset.Sub(time.Now().UTC()).Hours())
	// Calcualte tehe number of increments remaining until the daily limit is reached
	var incrementsRemaining int = int(Counters.Daily.Max - Counters.Daily.Count)
	if (hoursUntilReset == 0) || (incrementsRemaining < 1) {
		return 0
	}
	// Apportion the remaining increments to the number of hours until the next reset
	var hourlyIncrements int = incrementsRemaining / hoursUntilReset
	return hourlyIncrements
}

func GetRemainingThisHour(counterVals CounterValues) int {
	if time.Now().UTC().After(Counters.Hourly.NextReset) {
		Counters.Hourly.NextReset = time.Now().UTC().Add(time.Hour)
		Counters.Hourly.Count = 0
		Counters.Hourly.Max = int64(CalculateNewHourlyMax())
		go WriteBackCounters(counterVals, true)
		return int(Counters.Hourly.Max)
	}
	if Counters.Hourly.Max == 0 {
		return -1
	}
	if Counters.Hourly.Count >= Counters.Hourly.Max {
		return 0
	}
	return int(Counters.Hourly.Max - Counters.Hourly.Count)
}

func GetRemainingToday() int {
	return int(Counters.Daily.Max - Counters.Daily.Count)
}

func GetRemainingCount() (int, int, int) {
	return CountersProvider.GetInt("counters.monthly.max") - CountersProvider.GetInt("counters.monthly.count"),
		CountersProvider.GetInt("counters.daily.max") - CountersProvider.GetInt("counters.daily.count"),
		CountersProvider.GetInt("counters.hourly.max") - CountersProvider.GetInt("counters.hourly.count")
}

func ResetCounters(args []string) {
	// TODO
	//InitailiseCounters(args[0], args[1], args[2], args[3], true)
	//provider = InitailiseCounters(filename, 30000, 900, 900, true)
	//Counters = *new_counters
}

/*
func CheckCounters() bool {
	var dateTimeNow time.Time = time.Now().UTC()
	var dateTimeLastReset time.Time = Counters.GetTime("counters.lastreset")
	var dateTimeNextReset time.Time = Counters.GetTime("counters.nextreset")
	var dateTimeLastCheck time.Time = Counters.GetTime("counters.lastcheck")

	// Monthly counter
	if dateTimeNow.After(dateTimeNextReset) {
		util.LogDebug("Resetting counters")

		Counters.Set("counters.lastreset", dateTimeNow.Format(time.RFC3339))
		Counters.Set("counters.nextreset", dateTimeNow.AddDate(0, 1, 0).Format(time.RFC3339))
		Counters.Set("counters.dailycount", 0)
		Counters.Set("counters.hourlycount", 0)
		Counters.Set("counters.lastcheck", dateTimeNow.Format(time.RFC3339))

		err := Counters.WriteConfig()
		util.HandleFatalError(err, "Unable to write configuration:\n")

		err = Counters.Unmarshal(&CounterValues)
		util.HandleFatalError(err, "Unable to unmarshal counter values:\n")

		return true
	}

	// Daily counter, reset if the last reset was more than 24 hours ago
	if dateTimeNow.After(dateTimeLastReset.AddDate(0, 0, 1)) {
		util.LogDebug("Resetting daily counters")

		Counters.Set("counters.lastreset", dateTimeNow.Format(time.RFC3339))
		Counters.Set("counters.dailycount", 0)
		Counters.Set("counters.hourlycount", 0)
		Counters.Set("counters.lastcheck", dateTimeNow.Format(time.RFC3339))
	}

	// Instead of reseting the hourly counter, we just check what time the last reset was and we check how many hows until the 24 hours passes, then we adjust the hourly limit accordingly
	if dateTimeNow.After(dateTimeLastCheck.Add(time.Hour)) {
		util.LogDebug("Resetting hourly counters")

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
	util.HandleFatalError(err, "Unable to write configuration:\n")

	err = Counters.Unmarshal(&CounterValues)
	util.HandleFatalError(err, "Unable to unmarshal counter values:\n")

	return false
}
*/
