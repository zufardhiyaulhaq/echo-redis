package main

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/zufardhiyaulhaq/echo-redis/pkg/settings"

	redis_client "github.com/zufardhiyaulhaq/echo-redis/pkg/redis"
)

func main() {

	settings, err := settings.NewSettings()
	if err != nil {
		panic(err.Error())
	}

	log.Info().Msg("creating redis client")
	var client redis_client.RedisClient
	if settings.RedisCluster {
		client = redis_client.NewCluster(settings)
	} else {
		client = redis_client.New(settings)
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	log.Info().Msg("starting server")
	server := NewServer(settings, client)

	go func() {
		log.Info().Msg("starting HTTP server")
		server.ServeHTTP()
		wg.Done()
	}()

	wg.Wait()

	defer func() {
		log.Info().Msg("closing redis client")
		if err = client.Close(); err != nil {
			panic(err)
		}
	}()
}
