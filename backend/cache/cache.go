package cache

import (
	"context"
	"fmt"
	"ip2location-pfsense/config"
	. "ip2location-pfsense/util"

	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
)

type RedisInstance struct {
	Config  redis.Options
	Rdb     *redis.Client
	Rh      *rejson.Handler
	KeysSet int
}

var ctx = context.Background()
var instances map[string]RedisInstance

func init() {
	LogDebug("Initialising cache service")

	instances = make(map[string]RedisInstance)
}

// Create Redis cache instances based on the configuration
func CreateInstances() {

	LogDebug("Executing CreateInstances")

	LogDebug("Mapping Redis config")
	conf := config.GetConfig().Get("redis")

	for key, val := range conf.(map[string]interface{}) {
		LogDebug("Redis config: %v = %v", key, val)

		subkey := fmt.Sprintf("redis.%s", key)
		rc, err := LoadConfiguration(subkey)
		HandleError(err, "Unable to load configuration for %s", key)
		CreateInstance(key, rc)
	}
}

// Creates an instance of the Redis cache, stores a reference to it in the instances map and returns it
func CreateInstance(name string, config RedisCacheConfig) RedisInstance {

	options := redis.Options{
		Addr:     config.HostPort,
		Password: config.Pass,
		DB:       config.Db,
	}

	instance := RedisInstance{
		Config: options,
		Rdb:    redis.NewClient(&options),
		Rh:     rejson.NewReJSONHandler(),
	}

	instances[name] = instance

	*instances[name].Rh = *rejson.NewReJSONHandler()
	instances[name].Rh.SetGoRedisClientWithContext(ctx, instances[name].Rdb)

	return instances[name]
}

// Returns a Redis instance from the instances map
func Instance(name string) RedisInstance {
	LogDebug("Getting Redis instance: %v", name)

	return instances[name]
}

// Gets a value from the Redis instance
func (ri RedisInstance) Get(key string) (interface{}, error) {
	LogDebug("Getting key: %v from Redis instance: %v", key, ri)

	return ri.Rh.JSONGet(key, ".")
}

// Sets	a value in the Redis instance
func Set(name string, key string, value interface{}) (interface{}, error) {
	LogDebug("Setting key: %v to value: %v in Redis instance: %v", key, value, name)
	ri := instances[name]

	ri.KeysSet++

	return ri.Rh.JSONSet(key, ".", value)
}

func Handler(name string) *rejson.Handler {
	LogDebug("Getting Redis handler for instance: %v", name)
	ri := instances[name]
	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClientWithContext(ctx, ri.Rdb)

	return rh
}
