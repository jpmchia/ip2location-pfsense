package pfsense

import (
	"encoding/json"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"

	"log"
	"pfSense/cache"
	"pfSense/config"

	"strconv"
)

var PfSenseRedisPool *redis.Pool
var riPf cache.RedisCache

func init() {
	log.Println("Loading configuration for pfSense")
	config.LoadConfigProvider("IP2LOCATION")

}

func ProcessLog(c echo.Context) (*Ip2ResultId, error) {
	var err error
	var logEntries FilterLog

	body := c.Request().Body

	err = json.NewDecoder(body).Decode(&logEntries)

	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	res := Process(logEntries)
	var resultVal Ip2ResultId
	resultVal.Id = strconv.Itoa(res)

	return &resultVal, nil
}

func ProcessLogEntries(logEntries FilterLog) int {
	var result int

	//result = Process(logEntries)

	return result
}

func GetResult(resultid string) ([]Ip2Map, error) {
	var result []Ip2Map

	return result, nil
}

func Process(logEntries FilterLog) int {
	var result int

	//result = Process(logEntries)

	return result
}
