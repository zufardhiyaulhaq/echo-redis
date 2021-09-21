package settings

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	RedisEventPort string   `envconfig:"PORT"`
	RedisCluster   bool     `envconfig:"REDIS_CLUSTER"`
	RedisHosts     []string `envconfig:"REDIS_HOST"`
	RedisRetry     int      `envconfig:"REDIS_RETRY" default:"-1"`
}

func NewSettings() (Settings, error) {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
