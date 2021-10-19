package main

import (
	"sync"

	"github.com/zufardhiyaulhaq/echo-redis/pkg/settings"

	redis_client "github.com/zufardhiyaulhaq/echo-redis/pkg/redis"
)

func main() {

	settings, err := settings.NewSettings()
	if err != nil {
		panic(err.Error())
	}

	var client redis_client.RedisClient
	if settings.RedisCluster {
		client = redis_client.NewCluster(settings)
	} else {
		client = redis_client.New(settings)
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	server := NewServer(settings, client)

	go func() {
		server.ServeEcho()
		wg.Done()
	}()

	go func() {
		server.ServeHTTP()
		wg.Done()
	}()

	wg.Wait()

	defer func() {
		if err = client.Close(); err != nil {
			panic(err)
		}
	}()
}
