package main

import (
	"github.com/google/uuid"
	"github.com/tidwall/evio"
	"github.com/zufardhiyaulhaq/echo-redis/pkg/settings"

	redis_client "github.com/zufardhiyaulhaq/echo-redis/pkg/redis"
)

func main() {
	var events evio.Events

	settings, err := settings.NewSettings()
	if err != nil {
		panic(err.Error())
	}

	var redisClient redis_client.RedisClient
	if settings.RedisCluster {
		redisClient = redis_client.NewCluster(settings.RedisHosts)
	} else {
		redisClient = redis_client.New(settings.RedisHosts[0])
	}

	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		key := uuid.NewString()
		value := string(in)

		err := redisClient.Set(key, value)
		if err != nil {
			out = []byte(err.Error())
			return
		}

		valueRedis, err := redisClient.Get(key)
		if err != nil {
			out = []byte(err.Error())
			return
		}

		out = []byte(valueRedis)

		return
	}

	if err := evio.Serve(events, "tcp://0.0.0.0:"+settings.RedisEventPort); err != nil {
		panic(err.Error())
	}

	defer func() {
		if err = redisClient.Close(); err != nil {
			panic(err)
		}
	}()

}
