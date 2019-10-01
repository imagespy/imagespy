package config

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Registry struct {
	Address       string
	Auth          string
	BasicPassword string
	BasicUsername string
	Protocol      string
}

type Config struct {
	CacheExpiration    time.Duration
	DiscoveryDirectory string
	DockerRegistries   []Registry
	HTTPAddress        string
	LogLevel           log.Level
	PrometheusPath     string
}

func Parse(path string) (cfg Config, _ error) {
	viper.SetEnvPrefix("imagespy")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBufferString(defaultConfig))
	if err != nil {
		return cfg, fmt.Errorf("parse default config: %v", err)
	}

	if path != "" {
		viper.SetConfigFile(path)
		err = viper.MergeInConfig()
		if err != nil {
			return cfg, fmt.Errorf("parse config: %v", err)
		}
	}

	logLvl, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		return cfg, fmt.Errorf("parse log level: %w", err)
	}

	var registries []Registry
	for rKey, v := range viper.GetStringMap("docker.registries") {
		r, ok := v.(map[string]interface{})
		if !ok {
			return cfg, fmt.Errorf("docker.registries.%s is invalid", rKey)
		}

		registries = append(registries, Registry{
			Address:  r["address"].(string),
			Auth:     r["auth"].(string),
			Protocol: r["protocol"].(string),
		})
	}

	return Config{
		CacheExpiration:    viper.GetDuration("cache.expiration"),
		DiscoveryDirectory: viper.GetString("discovery.directory"),
		DockerRegistries:   registries,
		HTTPAddress:        viper.GetString("http.address"),
		LogLevel:           logLvl,
		PrometheusPath:     viper.GetString("prometheus.path"),
	}, nil
}
