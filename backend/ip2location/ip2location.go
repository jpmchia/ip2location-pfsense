package ip2location

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"ip2location-pfsense/cache"
	"ip2location-pfsense/config"
	"ip2location-pfsense/util"

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

	util.LogDebug("[ip2location] IP2Location API configuration: %v", Ip2ApiConfig)
}

func RetrieveIpLocation(ipAddress string, key string) (*Ip2LocationBasic, error) {
	var ip2location *Ip2LocationBasic
	var rh *rejson.Handler = cache.Handler(Ip2LocationCache)
	var err error

	outJSON, err := redis.Bytes(rh.JSONGet(key, "."))
	if err != nil {
		util.LogDebug("Failed to JSONGet: %s", err.Error())
	} else if err == nil {
		readIp := Ip2LocationBasic{}
		err = json.Unmarshal(outJSON, &readIp)
		if err == nil {
			log.Printf("[ip2location] Found in cache: %s", key)
			return &readIp, nil
		}
	}

	util.LogDebug("[ip2location] Not found in cache: %s", key)

	if errors > Ip2ApiConfig.MaxErrors {
		return nil, fmt.Errorf("not calling ip2location.io, many errors: %v", errors)
	}

	ip2location, err = LookupIPLocation(ipAddress)

	if err != nil {
		util.HandleError("[ip2location] Unable to retrieve: %s", err.Error())
		return nil, err
	}

	if ip2location == nil {
		util.HandleError("[ip2location] Unable to retrieve: %s", err.Error())
		return nil, err
	}

	util.LogDebug("[ip2location] Adding IP2Location API response the cache: %v", ip2location)

	rh = cache.Handler(Ip2LocationCache)
	_, err = rh.JSONSet(key, ".", *ip2location)
	if err != nil {
		util.HandleError(err, "[ip2location] Failed to store results in cache")
		return nil, err
	}

	b, err := json.MarshalIndent(*ip2location, "", "  ")
	if err != nil { // Handle the error
		util.HandleError(err, "[ip2location] Unable to marshal: %s", err)
		return nil, err
	}
	log.Printf("[ip2location] Added to the cache: %s", b)

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

func LookupIPLocation(ipAddress string) (*Ip2LocationBasic, error) {

	ip2_urlquery := Ip2ApiConfig.URL + "?key=" + Ip2ApiConfig.Key + "&ip=" + ipAddress + "&format=json"
	log.Default().Printf("[ip2location] Looking up IP Location API: %v", ipAddress)

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

	var ipLocation Ip2LocationBasic
	err = json.NewDecoder(resp.Body).Decode(&ipLocation)
	resp.Body.Close()
	if err != nil {
		util.HandleError(err, "[ip2location] LookupIPLocation: JSON decode failed.\n %s", err.Error())
		return nil, err
	}

	return &ipLocation, nil
}
