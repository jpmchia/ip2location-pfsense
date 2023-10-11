package ip2location

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jpmchia/ip2location-pfsense/cache"
	"github.com/jpmchia/ip2location-pfsense/config"
	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/jpmchia/ip2location-pfsense/version"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
)

// var ctx = context.Background()

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
var ip2ApiClient http.Client
var errors int

const Ip2LocationCache = "ip2location"

type ApiErr struct {
	Err        error
	StatusCode int
	Response   *http.Response
}

func (e *ApiErr) Error() string {
	return e.Err.Error()
}

func init() {
	util.LogDebug("[ip2location] Loading configuration for IP2Location cache.")

	var conf = config.GetConfiguration()
	Ip2ApiConfig = conf.IP2API
	Ip2ApiConfig.Source = url.QueryEscape(version.AppName)
	Ip2ApiConfig.Version = url.QueryEscape(version.Version)

	util.LogDebug("[ip2location] IP2Location API configuration: %v", Ip2ApiConfig)
}

func RetrieveIpLocationFromCache(ipAddress string, key string) (ip2location *Ip2LocationEntry, err error) {

	var rh *rejson.Handler = cache.Handler(Ip2LocationCache)

	outJSON, err := redis.Bytes(rh.JSONGet(key, "."))
	if err != nil {
		util.LogDebug("[ip2location] RetrieveIpLocationFromCache:  Failed to JSONGet: %s", err.Error())
	} else if err == nil {
		//ip2location := Ip2LocationPlus{}
		err = json.Unmarshal(outJSON, &ip2location)
		if err == nil {
			log.Printf("[ip2location] RetrieveIpLocationFromCache:  Found in cache: %s", key)
			return ip2location, nil
		}
	}
	util.LogDebug("[ip2location] RetrieveIpLocationFromCache:  Not found in cache: %s", key)

	return nil, nil
}

func RetrieveIpLocationFromApi(ipAddress string, key string) (ip2location *Ip2LocationEntry, err error) {

	if errors > Ip2ApiConfig.MaxErrors {
		return nil, fmt.Errorf("not calling ip2location.io, many errors: %v", errors)
	}

	ip2location, err = LookupIPLocation(ipAddress)

	if err != nil || ip2location == nil {
		util.HandleError("[ip2location] RetrieveIpLocationFromApi:  Unable to retrieve: %s", err.Error())
		errors++
		return nil, err
	}

	util.LogDebug("[ip2location] RetrieveIpLocationFromApi:  Adding IP2Location API response the cache: %v", ip2location)

	var rh = cache.Handler(Ip2LocationCache)
	_, err = rh.JSONSet(key, ".", *ip2location)
	if err != nil {
		util.HandleError(err, "[ip2location] RetrieveIpLocationFromApi:  Failed to store results in cache")
		return nil, err
	}

	// b, err := json.MarshalIndent(*ip2location, "", "  ")
	// if err != nil { // Handle the error
	// 	util.HandleError(err, "[ip2location] RetrieveIpLocationFromApi:  Unable to marshal: %s", err)
	// 	return nil, err
	// }
	// log.Printf("[ip2location] RetrieveIpLocationFromApi:  Added to the cache: %s", b)

	return ip2location, nil
}

func (ip2location *Ip2LocationEntry) UpdateHits(lastSeen string, onWatchlist bool, key string) (err error) {

	// if ip2location == nil {
	// 	return fmt.Errorf("ip2location is nil")
	// }
	// ip2location := *pIp2location

	if len(ip2location.FirstSeen) == 0 || ip2location.FirstSeen == "" {
		util.LogDebug("[ip2location] UpdateHits:  FirstSeen is empty, setting to lastSeen: %s", lastSeen)
		ip2location.FirstSeen = lastSeen
		ip2location.LastSeen = lastSeen
		ip2location.Watched = onWatchlist
		ip2location.Hits = 1
		util.LogDebug("[ip2location] UpdateHits:  LastSeen set to lastSeen: %s %s", lastSeen, ip2location.LastSeen)
		util.LogDebug("[ip2location] UpdateHits:  Updated IP2Location hits to: %v", ip2location.Hits)
		return nil
	}

	if ip2location.LastSeen == lastSeen {
		util.LogDebug("[ip2location] UpdateHits:  LastSeen is equal to lastSeen: %s", lastSeen)
		ip2location.LastSeen = lastSeen
		return fmt.Errorf("no updates as lastSeen is equal to lastSeen: %s = %v", lastSeen, ip2location.LastSeen)
	}

	ip2location.LastSeen = lastSeen
	ip2location.Watched = onWatchlist
	ip2location.Hits++

	util.LogDebug("[ip2location] UpdateHits:  LastSeen set to lastSeen: %s %s", lastSeen, ip2location.LastSeen)
	util.LogDebug("[ip2location] UpdateHits:  Updated IP2Location hits to: %v", ip2location.Hits)

	err = ip2location.SaveToCache(key)
	util.HandleError(err, "[ip2location] UpdateHits:  Failed to store results in cache")

	return nil
}

func (ip2location Ip2LocationEntry) SaveToCache(key string) error {
	var rh *rejson.Handler = cache.Handler(Ip2LocationCache)

	_, err := rh.JSONSet(key, ".", ip2location)

	if err != nil {
		util.HandleError(err, "[ip2location] SaveToCache:  Failed to store results in cache")
		return err
	}

	// b, err := json.MarshalIndent(ip2location, "", "  ")
	// if err != nil { // Handle the error
	// 	util.HandleError(err, "[ip2location] SaveToCache:  Unable to marshal: %s", err)
	// 	return err
	// }
	//log.Printf("[ip2location] SaveToCache:  Added to the cache: %s", key)

	return nil
}

func AddToCache(ipAddress string, ip2location *Ip2LocationEntry) (err error) {
	var rh *rejson.Handler = cache.Handler(Ip2LocationCache)
	var key = strings.ReplaceAll(ipAddress, ":", ".")
	_, err = rh.JSONSet(key, ".", *ip2location)
	if err != nil {
		util.HandleError(err, "[ip2location] AddToCache:  Failed to store results in cache")
		return err
	}
	return nil
}

func RetrieveIpPlus(ipAddress string) (*Ip2LocationEntry, error) {
	var ip2location *Ip2LocationEntry
	var rh *rejson.Handler = cache.Handler(Ip2LocationCache)
	var err error
	var key = strings.ReplaceAll(ipAddress, ":", ".")

	ip2location, err = LookupIPLocation(ipAddress)

	if err != nil {
		util.HandleError("[ip2location] RetrieveIpPlus:  Unable to retrieve: %s", err.Error())
		return nil, err
	}

	if ip2location == nil {
		util.HandleError("[ip2location] RetrieveIpPlus:  Unable to retrieve: %s", err.Error())
		return nil, err
	}

	util.LogDebug("[ip2location] RetrieveIpPlus:  Adding IP2Location API response the cache: %v", ip2location)

	rh = cache.Handler(Ip2LocationCache)
	_, err = rh.JSONSet(key, ".", *ip2location)
	if err != nil {
		util.HandleError(err, "[ip2location] RetrieveIpPlus:  Failed to store results in cache")
		return nil, err
	}

	b, err := json.MarshalIndent(*ip2location, "", "  ")
	if err != nil { // Handle the error
		util.HandleError(err, "[ip2location] RetrieveIpPlus:  Unable to marshal: %s", err)
		return nil, err
	}
	log.Printf("[ip2location] RetrieveIpPlus:  Added to the cache: %s", b)

	return ip2location, nil
}

func apiError(err error, response *http.Response) error {
	statusCode := 0
	if response != nil {
		statusCode = response.StatusCode
	}

	apiErr := &ApiErr{
		err,
		statusCode,
		response,
	}

	return fmt.Errorf("%w", apiErr)
}

func LookupIPLocation(ipAddress string) (*Ip2LocationEntry, error) {

	ip2_urlquery := Ip2ApiConfig.URL + "?key=" + Ip2ApiConfig.Key + "&ip=" + ipAddress + "&format=json" + "&source=" + Ip2ApiConfig.Source + "&source_version=" + Ip2ApiConfig.Version
	log.Default().Printf("[ip2location] LookupIPLocation:  Looking up IP Location API: %v", ipAddress)

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	ip2ApiClient = http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	resp, err := ip2ApiClient.Get(ip2_urlquery)

	if err != nil {
		err = apiError(err.(*url.Error), resp)
		return nil, err
	}

	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		resperr := apiError(fmt.Errorf("bad status: %s", resp.Status), resp)
		return nil, resperr
	}

	var ipLocationPlus Ip2LocationPlus
	err = json.NewDecoder(resp.Body).Decode(&ipLocationPlus)
	resp.Body.Close()
	if err != nil {
		util.HandleError(err, "[ip2location] LookupIPLocation: JSON decode failed.\n %s", err.Error())
		return nil, err
	}
	var ipLocationEntry = Ip2LocationEntry{Ip2LocationPlus: ipLocationPlus}

	return &ipLocationEntry, nil
}

func LookupIPLocationBasic(ipAddress string) (ipLocation *Ip2LocationBasic, err error) {

	ip2_urlquery := Ip2ApiConfig.URL + "?key=" + Ip2ApiConfig.Key + "&ip=" + ipAddress + "&format=json" + "&source=" + Ip2ApiConfig.Source + "&source_version=" + Ip2ApiConfig.Version
	log.Default().Printf("[ip2location] LookupIPLocationBasic:  Looking up IP Location API: %v", ipAddress)

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	ip2ApiClient = http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	resp, err := ip2ApiClient.Get(ip2_urlquery)

	if err != nil {
		err = apiError(err.(*url.Error), resp)
		return nil, err
	}

	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		resperr := apiError(fmt.Errorf("bad status: %s", resp.Status), resp)
		return nil, resperr
	}

	err = json.NewDecoder(resp.Body).Decode(&ipLocation)
	resp.Body.Close()
	if err != nil {
		util.HandleError(err, "[ip2location] LookupIPLocationBasic: JSON decode failed.\n %s", err.Error())
		return nil, err
	}

	return ipLocation, nil
}
