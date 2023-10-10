package apikey

import (
	"errors"
	"math/rand"
	"time"

	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/labstack/echo/v4"
)

type ApiKey struct {
	Key       string
	Issued    time.Time
	Expires   time.Time
	IpAddress string `json:"ip_address"`
}

var validKeys map[string]ApiKey

const keyLength = 10
const keyChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const keyExpiry = 5 * 60 * 1000 // 5 minutes

func init() {
	util.LogDebug("[apikey] Initialising API keys")
	validKeys = make(map[string]ApiKey)
}

func KeyHandler(c echo.Context) error {

	ipAddr := c.RealIP()

	if ipAddr == "" {
		util.Log("[apikeys] Unable to issue a key, no IP address found")
		return c.JSON(400, "No IP address found")
	}

	apiKey := AddKey(ipAddr)
	util.Log("[apikeys] Issued key: %s, for IP address: %s, Key expires at: %s", apiKey.Key, apiKey.IpAddress, apiKey.Expires.Format(time.RFC3339))

	return c.JSON(200, apiKey)
}

func ApiKeyHandler(c echo.Context) error {
	ipAddr := c.RealIP()
	key := AddKey(ipAddr)
	return c.JSON(200, key)
}

// Generates a random key, of length keyLength, from the characters in keyChars
func GenerateKey() string {
	key := ""
	for i := 0; i < keyLength; i++ {
		key += string(keyChars[rand.Intn(len(keyChars))])
	}
	return key
}

// Generates a random key, of length keyLength, from the characters in keyChars
func GenerateApiKey(ipAddr string, validFor int) ApiKey {
	key := GenerateKey()
	issued := time.Now()
	expires := time.Now().Add(time.Duration(validFor) * time.Millisecond)
	ipAddress := ipAddr
	return ApiKey{key, issued, expires, ipAddress}
}

// Generates and adds a new key to the validKeys map
func AddKey(ipAddr string) ApiKey {
	key := GenerateApiKey(ipAddr, keyExpiry)
	validKeys[key.Key] = key
	return key
}

// Removes a key from the validKeys map
func RemoveKey(key string) {
	delete(validKeys, key)
}

// Returns a slice of all the valid keys
func GetKeys() []ApiKey {
	keys := []ApiKey{}
	for _, v := range validKeys {
		keys = append(keys, v)
	}
	return keys
}

// Returns a key from the validKeys map
func GetKey(key string) (ApiKey, bool) {
	k, ok := validKeys[key]
	return k, ok
}

// Returns a key from the validKeys map by IP address
func GetKeyByIp(ip string) (ApiKey, bool) {
	for _, v := range validKeys {
		if v.IpAddress == ip {
			return v, true
		}
	}
	return ApiKey{}, false
}

// Returns a key from the validKeys map by IP address and key
func GetKeyByIpAndKey(ip string, key string) (ApiKey, bool) {
	for _, v := range validKeys {
		if v.IpAddress == ip && v.Key == key {
			return v, true
		}
	}
	return ApiKey{}, false
}

func RemoveExpiredKeys() {
	for _, v := range validKeys {
		if v.Expires.Before(time.Now()) {
			util.Log("[apikeys] Removing expired key: %s, for IP address: %s, Key expired at: %s", v.Key, v.IpAddress, v.Expires.Format(time.RFC3339))
			RemoveKey(v.Key)
		}
	}
}

// Validates a key by IP address and key
func ValidateIpKey(ip string, key string) bool {

	apiKey, found := GetKeyByIpAndKey(ip, key)

	if !found {
		util.Log("[apikeys] Key: %s; not found for ip: %s", key, ip)
		return false
	}

	if apiKey.Expires.Before(time.Now()) {
		util.Log("[apikeys] Key: %s; expired: %s; removing from valid keys", key, apiKey.Expires.Format(time.RFC3339))
		RemoveKey(key)
		return false
	}

	if apiKey.Expires.After(time.Now()) {
		util.Log("[apikeys] Valid key received: %s; removing from valid keys", key)
		RemoveKey(key)
		return true
	}

	return false
}

func ValidateKey(key string, c echo.Context) (bool, error) {
	if key == "" {
		util.LogDebug("[service] No key supplied")
		return false, errors.New("missing key parameter")
	}

	ipAddr := c.RealIP()
	if ValidateIpKey(ipAddr, key) {
		return true, nil
	}

	return false, nil
}
