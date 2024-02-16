package settings

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	HTTPPort         string   `envconfig:"HTTP_PORT" default:"80"`
	RedisCluster     bool     `envconfig:"REDIS_CLUSTER"`
	RedisHosts       []string `envconfig:"REDIS_HOST"`
	RedisIdleTiemout int      `envconfig:"REDIS_IDLE_TIMEOUT" default:"-1"`
	RedisRetry       int      `envconfig:"REDIS_RETRY" default:"-1"`
}

func NewSettings() (Settings, error) {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
