package main

import (
	"log"
	"pfSense/cache"
	"pfSense/cmd"
	"pfSense/config"
)

var Ip2LocationCache cache.Instance
var PfSenseResultsCache cache.Instance

func configure() {
	log.Print("Configuring Redis cache ...")

	config.LoadConfigProvider("IP2LOCATION")

	if config.Config().GetBool("use_cache") {

		Ip2LocationCache := new(cache.Instance)
		Ip2LocationCache.RedisInstance.RedisHostPort = config.Config().GetString("redis.ip2location.hostport")
		Ip2LocationCache.RedisInstance.RedisDb = config.Config().GetInt("redis.ip2location.db")
		Ip2LocationCache.RedisInstance.RedisAuth = config.Config().GetString("redis.ip2location.auth")
		Ip2LocationCache.RedisInstance.RedisPass = config.Config().GetString("redis.ip2location.pass")

		PfSenseResultsCache := new(cache.Instance)
		PfSenseResultsCache.RedisInstance.RedisHostPort = config.Config().GetString("redis.pfsense.hostport")
		PfSenseResultsCache.RedisInstance.RedisDb = config.Config().GetInt("redis.pfsense.db")
		PfSenseResultsCache.RedisInstance.RedisAuth = config.Config().GetString("redis.pfsense.auth")
		PfSenseResultsCache.RedisInstance.RedisPass = config.Config().GetString("redis.pfsense.pass")

	}
}

func main() {

	configure()

	Ip2LocationCache.RedisPool = cache.NewPool(Ip2LocationCache.RedisInstance)
	PfSenseResultsCache.RedisPool = cache.NewPool(PfSenseResultsCache.RedisInstance)

	cmd.Execute()
}
