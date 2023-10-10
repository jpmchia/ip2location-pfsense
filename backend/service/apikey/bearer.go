package apikey

import (
	"errors"
	"log"

	"github.com/jpmchia/ip2location-pfsense/config"

	"github.com/jpmchia/ip2location-pfsense/util"

	"github.com/labstack/echo/v4"
)

var validTokens []string

func init() {
	util.LogDebug("[apikey] Initialising API bearer tokens")

	config.Configure()

	validTokens = config.ConfigProvider().GetStringSlice("apitokens")

	if len(validTokens) == 0 {
		log.Printf("[apikeys] No 'apitokens' configured. The API beaker token authentication will not be available and will subsequently deny all requests.")
		log.Printf("%v", config.ConfigProvider().AllSettings())
	} else {
		util.Log("[apikeys] API bearer tokens: %v", validTokens)
	}
}

func ValidateToken(key string, c echo.Context) (bool, error) {

	if key == "" {
		util.Log("[apikey] Missing API key")
		return false, errors.New("missing api key")
	}

	validTokens = config.ConfigProvider().GetStringSlice("apitokens")

	if len(validTokens) == 0 {
		util.Log("[apikey] No API bearer tokens configured")
		return true, errors.New("no api bearer tokens configured")
	}

	for _, validToken := range validTokens {
		if key == validToken {
			util.Log("[apikey] Validated API bearer token.")
			return true, nil
		}
	}
	util.Log("[apikey] Invalid API token: %s", key)
	return false, nil
}
