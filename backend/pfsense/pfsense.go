package pfsense

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"ip2location-pfsense/cache"
	"ip2location-pfsense/ip2location"
	"ip2location-pfsense/util"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
	"github.com/nitishm/go-rejson/v4"

	"strconv"
)

const pfSenseCache string = "pfsense"

//const pfSenseCacheConfig string = "redis.pfsense"

// var ctx = context.Background()
// var addrpf *string

func init() {
	// LogDebug("Loading configuration for pfSense")
	// c := config.LoadConfigProvider("IP2Location-pfSense")
	// c.Get(pfSenseCache)
	//cache.LoadConfiguration(pfSenseCacheConfig)
}

func ProcessLog(c echo.Context) (*Ip2ResultId, error) {
	var err error
	var logEntries FilterLog

	body := c.Request().Body
	err = json.NewDecoder(body).Decode(&logEntries)

	if err != nil {
		util.HandleFatalError(err, "[pfsense] Failed reading the request body %s", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	res := ProcessLogEntries(logEntries)
	var resultVal Ip2ResultId
	resultVal.Id = strconv.FormatInt(res, 16)

	return &resultVal, nil
}

func ProcessLogEntries(logEntries FilterLog) int64 {
	var result int64
	var err error
	var ip2MapList []Ip2Map
	var ip2Map *Ip2Map

	for _, logEntry := range logEntries {
		util.LogDebug("[pfsense] Processing log entry: %v", logEntry)

		ip2Map, err = EnrichLogWithIp(logEntry)

		if ip2Map == nil && err == nil {
			util.HandleError(err, "[pfsense] Failed to process log entry. ip2Map came back nil.")
			continue
		} else if err != nil {
			util.HandleError(err, "[pfsense] Failed to process log entry: %v", err)
			continue
		}

		ip2MapList = append(ip2MapList, *ip2Map)
	}

	if err != nil {
		log.Fatal(err)
	}

	result = CacheResult(ip2MapList)

	return result
}

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

	//log.Printf("[pfsense] ResultSet = %s %v", strconv.FormatInt(resultSet, 16), res)
	log.Printf("[pfsense] ResultSet = %s %v", trunckey, res)
	return resultSet
}

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

func EnrichLogWithIp(logEntry LogEntry) (*Ip2Map, error) {
	var key string

	key, dir := DetermineIp(logEntry)

	if key == "" || dir == "" {
		util.LogDebug("[pfsense] No IP address to process - both private.")
		return nil, nil
	}

	logEntry.Direction = strings.ToLower(dir)
	nkey := strings.ReplaceAll(key, ":", ".")
	util.LogDebug("Key: %v => %v %s %v", nkey, logEntry.Srcip, dir, logEntry.Dstip)

	ip2Location, err := ip2location.RetrieveIpLocation(key, nkey)
	util.HandleError(err, "[pfsense] Unable to retrieve: %s", err)

	if err != nil {
		util.LogDebug("[pfsense] Unable to retrieve: %s", err)
		return nil, err
	}

	ip2Map, err := CreateIp2Map(logEntry, ip2Location)

	if err != nil {
		util.LogDebug("[pfsense] Unable to create IP2Map: %s", err)
		return nil, err
	}

	return &ip2Map, nil
}

func CreateIp2Map(logentry LogEntry, ip2Location *ip2location.Ip2LocationBasic) (Ip2Map, error) {
	var ip2map Ip2Map

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

	return ip2map, nil
}
