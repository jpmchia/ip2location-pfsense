package pfsense

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jpmchia/ip2location-pfsense/backend/cache"
	"github.com/jpmchia/ip2location-pfsense/backend/config"
	"github.com/jpmchia/ip2location-pfsense/backend/ip2location"
	"github.com/jpmchia/ip2location-pfsense/backend/util"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/nitishm/go-rejson/v4"

	"strconv"
)

const pfSenseCache string = "pfsense"

var LogIpMode bool = false

func init() {
	ActiveWatchList = NewWatchList()
	var Config = config.GetConfiguration()
	LogIpMode = Config.Service.LogIpMode
}

// Decodes the JSON request body into the FilterLog struct,
// it then calls the ProcessLogEntries function
func ProcessLog(c echo.Context) (pResultId Ip2ResultId, err error) {

	var LogEntries FilterLog

	body := c.Request().Body

	err = json.NewDecoder(body).Decode(&LogEntries)

	if err != nil {
		util.HandleFatalError(err, "[pfsense] Failed reading the request body %s", err)
		return pResultId, echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	// Pass the pointer to the filterlog entries for processing
	res := ProcessLogEntries(LogEntries)
	var resultVal Ip2ResultId
	resultVal.Id = strconv.FormatInt(res, 16)

	return resultVal, nil
}

// Extracts the IP address from the log entry and returns it as the key
// which is used to retrieve the IP2Location data from the cache
func ProcessLogEntries(logEntries FilterLog) int64 {
	var result int64
	var err error
	var ip2MapList []Ip2Map // Create the result set
	//var pIp2Map *Ip2Map     // Create a pointer to a IP2Map
	var ip2locationEntry *ip2location.Ip2LocationEntry
	var countProcessed int = 0
	var countLocated int = 0

	// Iterate through the log entries
	for _, logEntry := range logEntries {
		util.LogDebug("[pfsense] ProcessLogEntries:  Processing log entry: %v", logEntry)
		countProcessed++

		// Retrieve the public IP address from the log entry
		ip, dir := DeterminePublicIp(&logEntry)
		if ip == "" && dir == "" {
			// If the IP address is nil, then there is nothing to process
			util.LogDebug("[pfsense] ProcessLogEntries:  Skippping any enrichment and moving on to the next.")
			continue
		}
		util.LogDebug("[pfsense] ProcessLogEntries:  Public IP address: %v; and the direction is: %s", ip, dir)

		// Convert the IP address to a key
		key := strings.ReplaceAll(ip, ":", ".")
		util.LogDebug("[pfsense] ProcessLogEntries:  Generated key from the IP address: %s => %s", ip, key)

		watchList := ActiveWatchList.Contains(key)
		util.LogDebug("[pfsense] ProcessLogEntries:  IP address %v is on the watchlist: %v", key, watchList)

		// Process the log entry and produces an IP2Map struct that is returned as a pointer
		ip2locationEntry, err = RetrieveIp2LocationData(&logEntry, key, ip)

		if err != nil {
			util.HandleError(err, "[pfsense] ProcessLogEntries:  Failed to retrieve IP2Location data: %s", err)
			continue
		}

		err = ip2locationEntry.UpdateHits(logEntry.Time, watchList, key)
		if err == nil {
			util.LogDebug("[pfsense] ProcessLogEntries:  Updated IP2Location hits. Saving back to the cache")
			err = ip2locationEntry.SaveToCache(key)
			util.HandleError(err, "[pfsense] ProcessLogEntries:  Failed to save IP2Location data to cache: %s", err)
		}
		countLocated++

		util.LogDebug("[pfsense] ProcessLogEntries:  Adding the enriched log entry the result set.")

		ip2Map, err := CreateIp2Map(&logEntry, ip2locationEntry, watchList)
		if err != nil {
			util.HandleError(err, "[pfsense] ProcessLogEntries:  Failed to create IP2Map: %s", err)
			continue
		}

		// Append the IP2Map struct to the result set

		ip2MapList = append(ip2MapList, *ip2Map)
	}

	// Cache the result set
	log.Printf("[pfsense] ProcessLogEntries:  Processed %d log entries and located %d IP addresses", countProcessed, countLocated)
	result = CacheResult(ip2MapList)

	return result
}

// Checks the watchlist for the IP address, if it is on the watchlist
// then it updates the log entry with the watchlist data
func CheckWatchList(ipAddr string) bool {

	if ActiveWatchList.Contains(ipAddr) {

		util.LogDebug("[pfsense] CheckWatchList:  IP address %v is on the watchlist", ipAddr)
		return true
	}
	return false
}

// Caches the result set in Redis, returns the result set ID
// which is used to retrieve the results from the cache
func CacheResult(ip2MapList []Ip2Map) int64 {
	var rh *rejson.Handler = cache.Handler(pfSenseCache)

	log.Printf("[pfsense] Caching results for collection: %d results\n", len(ip2MapList))

	now := time.Now()
	resultSet := now.UnixNano()
	str := fmt.Sprintf("%d", resultSet)
	trunckey := str[0:13]

	log.Printf("[pfsense] Saving with the key: %s\n", trunckey)
	res, err := rh.JSONSet(trunckey, ".", ip2MapList)

	if err != nil {
		util.HandleFatalError(err, "[pfsense] Failed to store results in cache @ %v", err)
		return -1
	}

	log.Printf("[pfsense] ResultSet = %s %v", trunckey, res)
	return resultSet
}

// Retrieves the result set from the cache
func GetResult(resultid string) ([]Ip2Map, error) {
	var result []Ip2Map
	var rh *rejson.Handler = cache.Handler(pfSenseCache)
	var err error

	util.LogDebug("[pfsense] Attempting to retrieve from cache: %s", resultid)

	outJSON, err := redis.Bytes(rh.JSONGet(resultid, "."))
	if err != nil {

		util.LogDebug("[pfsense] Failed to JSONGet: %s", err.Error())
		return nil, nil

	} else if err == nil {

		result := []Ip2Map{}
		err = json.Unmarshal(outJSON, &result)

		if err == nil {
			log.Printf("[pfsense] Found in cache: %s", resultid)
			return result, nil
		}
	}
	return result, nil
}

// Retrieves the result set from the cache
func GetRawResult(resultid string) ([]byte, error) {
	var result []byte
	var rh *rejson.Handler = cache.Handler(pfSenseCache)
	var err error

	util.LogDebug("[pfsense] Attempting to retrieve from cache: %s", resultid)

	outBytes, err := redis.Bytes(rh.JSONGet(resultid, "."))
	if err != nil {
		util.LogDebug("[pfsense] Failed to JSONGet: %s", err.Error())
		return nil, nil
	} else if err == nil {
		log.Printf("[pfsense] Retrieved: %s", resultid)
		return outBytes, nil
	}
	return result, nil
}

// Checks the IP address to determine if it is a public or private IP address
func DeterminePublicIp(logEntry *LogEntry) (string, string) {

	util.LogDebug("[pfsense] DeterminePublicIp:  Checking %v", logEntry)

	ip, dir := DetermineIp(*logEntry)

	if ip == "" || dir == "" {
		util.LogDebug("[pfsense] DeterminePublicIp:  No IP address to process - both are private IP addresses.")
		return "", ""
	}

	logEntry.Direction = strings.ToLower(dir)

	return ip, dir
}

// Retrieves the result set from the cache and if not found, then
// it will retrieve the result set from the API
func RetrieveIp2LocationData(pLogentry *LogEntry, key string, ip string) (pIp2Location *ip2location.Ip2LocationEntry, err error) {

	// Try to retrieve the IP2Location data from the cache
	pIp2Location, err = ip2location.RetrieveIpLocationFromCache(ip, key)
	util.HandleError(err, "[pfsense] EnrichLogWithIp:  Failed to retrieve IP2Location data from cache: %s", err)

	if err == nil && pIp2Location != nil && !LogIpMode {
		return pIp2Location, err
	}

	pIp2Location, err = ip2location.RetrieveIpLocationFromApi(ip, key)
	util.HandleError(err, "[pfsense] EnrichLogWithIp:  Failed to retrieve IP2Location data from API: %s", err)

	if err == nil && pIp2Location != nil && !LogIpMode {
		return pIp2Location, err
	}

	if pIp2Location == nil {
		util.LogDebug("[pfsense] EnrichLogWithIp:  IP2Location data is nil")
		err = fmt.Errorf("unable to obtain ip2location data")
	}

	return nil, err
}

// Adds the log entry to the IP2Map struct
func AddLogEntryToIp2Map(pLogentry *LogEntry, pIp2Map *Ip2Map) {

	if pIp2Map.LastSeen == pLogentry.Time {
		return
	}

	var logentry = *pLogentry

	pIp2Map.Hits++
	pIp2Map.LogEntry = append(pIp2Map.LogEntry, logentry)
	pIp2Map.LastSeen = pLogentry.Time
}

func UpdateIp2Map(pLogentry *LogEntry, pIp2map *Ip2Map, watchList bool) error {

	var ip2map Ip2Map = *pIp2map
	var logentry LogEntry = *pLogentry

	util.LogDebug("[pfsense] UpdateIp2Map:  Updating IP2Map: %v", ip2map)

	if ip2map.LastSeen != pLogentry.Time {
		ip2map.Hits++
		util.LogDebug("[pfsense] UpdateIp2Map:  Updating IP2Map hits to: %v", ip2map.Hits)
	}
	ip2map.LastSeen = pLogentry.Time

	if watchList {
		var err error
		util.LogDebug("[pfsense] IP address %v is on the watchlist", ip2map.IP)
		AddLogEntryToIp2Map(&logentry, &ip2map)
		util.HandleError(err, "[pfsense] Failed to add log entry to watchlist: %v", err)
		return err
	}

	return nil
}

// Converts the log entry to an IP2Map struct with only minimal data
func ConvertToIp2Map(logEntry LogEntry, watchList bool) (*Ip2Map, error) {

	var ip2map Ip2Map

	ip2map.Time = logEntry.Time
	ip2map.IP = logEntry.Srcip
	ip2map.Direction = logEntry.Direction
	ip2map.Act = logEntry.Act
	ip2map.Reason = logEntry.Reason
	ip2map.Interface = logEntry.Interface
	ip2map.Realint = logEntry.Realint
	ip2map.Version = logEntry.Version
	ip2map.Srcip = logEntry.Srcip
	ip2map.Dstip = logEntry.Dstip
	ip2map.Srcport = logEntry.Srcport
	ip2map.Dstport = logEntry.Dstport
	ip2map.Proto = logEntry.Proto
	ip2map.Protoid = logEntry.Protoid
	ip2map.Length = logEntry.Length
	ip2map.Rulenum = logEntry.Rulenum
	ip2map.Subrulenum = logEntry.Subrulenum
	ip2map.Anchor = logEntry.Anchor
	ip2map.Tracker = logEntry.Tracker
	ip2map.CountryCode = ""
	ip2map.CountryName = ""
	ip2map.RegionName = ""
	ip2map.CityName = ""
	ip2map.ZipCode = ""
	ip2map.TimeZone = ""
	ip2map.Asn = ""
	ip2map.As = ""
	ip2map.IsProxy = false
	ip2map.WatchList = watchList
	ip2map.FirstSeen = logEntry.Time
	ip2map.LastSeen = logEntry.Time
	ip2map.Hits = 1
	ip2map.Minimal = true
	ip2map.WatchHits = 0 // Will be updated by the watchlist function

	if watchList {
		var err error
		util.LogDebug("[pfsense] IP address %v is on the watchlist", ip2map.IP)
		AddLogEntryToIp2Map(&logEntry, &ip2map)
		util.HandleError(err, "[pfsense] Failed to add log entry to watchlist: %v", err)
	}

	return &ip2map, nil
}

func CreateIp2Map(pLogentry *LogEntry, ip2Location *ip2location.Ip2LocationEntry, watchList bool) (*Ip2Map, error) {

	var ip2map Ip2Map
	var logentry LogEntry

	logentry = *pLogentry

	ip2map.Time = logentry.Time
	ip2map.IP = logentry.Srcip
	ip2map.Latitude = ip2Location.Latitude
	ip2map.Longitude = ip2Location.Longitude
	ip2map.Direction = logentry.Direction
	ip2map.Act = logentry.Act
	ip2map.Reason = logentry.Reason
	ip2map.Interface = logentry.Interface
	ip2map.Realint = logentry.Realint
	ip2map.Version = logentry.Version
	ip2map.Srcip = logentry.Srcip
	ip2map.Dstip = logentry.Dstip
	ip2map.Srcport = logentry.Srcport
	ip2map.Dstport = logentry.Dstport
	ip2map.Proto = logentry.Proto
	ip2map.Protoid = logentry.Protoid
	ip2map.Length = logentry.Length
	ip2map.Rulenum = logentry.Rulenum
	ip2map.Subrulenum = logentry.Subrulenum
	ip2map.Anchor = logentry.Anchor
	ip2map.Tracker = logentry.Tracker
	ip2map.CountryCode = ip2Location.CountryCode
	ip2map.CountryName = ip2Location.CountryName
	ip2map.RegionName = ip2Location.RegionName
	ip2map.CityName = ip2Location.CityName
	ip2map.ZipCode = ip2Location.ZipCode
	ip2map.TimeZone = ip2Location.TimeZone
	ip2map.Asn = ip2Location.Asn
	ip2map.As = ip2Location.As
	ip2map.IsProxy = ip2Location.IsProxy
	ip2map.WatchList = false
	ip2map.FirstSeen = logentry.Time
	ip2map.LastSeen = logentry.Time
	ip2map.Hits = 1
	ip2map.WatchHits = 0 // Will be updated by the watchlist function

	if watchList {
		var err error
		util.LogDebug("[pfsense] IP address %v is on the watchlist", ip2map.IP)
		AddLogEntryToIp2Map(&logentry, &ip2map)
		util.HandleError(err, "[pfsense] Failed to add log entry to watchlist: %v", err)
	}

	return &ip2map, nil
}

// 	if watchList {
// 		util.LogDebug("[pfsense] IP address %v is on the watchlist", key)
// 		//err = ActiveWatchList.AddLogEntry(ip2Map.IP, *ip2Map)
// 		util.HandleError(err, "[pfsense] Failed to add log entry to watchlist: %v", err)
// 	}
// 	if err != nil || ip2Location == nil {
// 		util.LogDebug("[pfsense] Unable to retrieve: %s", err)
// 		util.HandleError(err, "[pfsense] Unable to retrieve: %s", err)
// 	if err != nil {
// 		util.LogDebug("[pfsense] Unable to retrieve: %s", err)
// 		return nil, err
// 	}
// 	ip2Map, err := CreateIp2Map(logEntry, ip2Location)
// 	if err != nil {
// 		util.LogDebug("[pfsense] Unable to create IP2Map: %s", err)
// 		return nil, err
// 	}
// 	return &ip2Map, nil
// }
