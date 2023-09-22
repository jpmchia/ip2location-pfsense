package ip2location

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ip2location-pfsense/cache"
	"ip2location-pfsense/config"
	. "ip2location-pfsense/util"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
)

var ctx = context.Background()

type MapRequest struct {
	FromTime string `json:"time"`
	Action   string `json:"act"`
	Detail   string `json:"plan"`
}

type Ip2ResultId struct {
	Id int `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var Ip2ApiConfig config.Ip2ApiOptions
var Ip2ApiClient http.Client

const Ip2LocationCache = "ip2location"

func init() {
	LogDebug("Loading configuration for IP2Location cache.")

	var conf, err = config.LoadConfiguration()
	if err != nil {
		HandleFatalError(err, "Failed to load configuration: %s", err.Error())
	}

	Ip2ApiConfig = conf.IP2API
}

func RetrieveIpLocation(ipAddress string, key string) (*Ip2LocationBasic, error) {
	var ip2location *Ip2LocationBasic
	var rh *rejson.Handler = cache.Handler(Ip2LocationCache)
	var err error

	LogDebug("Checking in cache: %s", key)

	outJSON, err := redis.Bytes(rh.JSONGet(key, "."))
	if err != nil {
		LogDebug("Failed to JSONGet: %s", err.Error())
	} else if err == nil {
		readIp := Ip2LocationBasic{}
		err = json.Unmarshal(outJSON, &readIp)
		if err == nil {
			log.Printf("Found in cache: %s", key)
			return &readIp, nil
		}
	}

	LogDebug("Not found in cache: %s", key)

	ip2location, err = LookupIPLocation(ipAddress)

	if err != nil {
		HandleError("Unable to retrieve: %s", err.Error())
		return nil, err
	}

	LogDebug("Adding IP2Location API response the cache: %v", ip2location)

	rh = cache.Handler(Ip2LocationCache)
	_, err = rh.JSONSet(key, ".", *ip2location)
	if err != nil {
		HandleError(err, "Failed to store results in cache")
		return nil, err
	}

	b, err := json.MarshalIndent(*ip2location, "", "  ")
	if err != nil { // Handle the error
		HandleError(err, "Unable to marshal: %s", err)
		return nil, err
	}
	log.Printf("Added to the cache: %s", b)

	return ip2location, nil
}

func LookupIPLocation(ipAddress string) (*Ip2LocationBasic, error) {
	log.Default().Printf("Looking up IP Location API: %s", ipAddress)
	ip2_urlquery := Ip2ApiConfig.URL + "?key=" + Ip2ApiConfig.Key + "&ip=" + ipAddress + "&format=json"

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	Ip2ApiClient = http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	// Ip2ApiClient = http.NewRequest("GET", ip2_urlquery, nil)
	response, err := Ip2ApiClient.Get(ip2_urlquery)
	var ipLocation Ip2LocationBasic
	if err != nil {
		HandleError(err, "LookupIP: HTTP request failed.\n %s", err.Error())
		return &ipLocation, err
	}
	//json.Unmarshal([]byte(), &logEntries)
	err = json.NewDecoder(response.Body).Decode(&ipLocation)
	if err != nil {
		HandleError(err, "LookupIP: JSON decode failed.\n %s", err.Error())
		return &ipLocation, err
	}
	return &ipLocation, nil
}
