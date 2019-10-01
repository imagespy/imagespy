package core

import (
	"net/http"

	"github.com/imagespy/imagespy/core/config"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Runner struct {
	cfg config.Config
}

func (r *Runner) Run() error {
	log.SetLevel(r.cfg.LogLevel)
	scraper := NewScraper(r.cfg.DockerRegistries)
	watcher, err := NewDirectory(r.cfg.DiscoveryDirectory)
	if err != nil {
		log.Fatalf("create directory watcher: %v", err)
	}

	finder := NewFinder(watcher, scraper, NewGoCache(r.cfg.CacheExpiration))
	pc := NewCollector(finder)
	prom.MustRegister(pc)

	http.Handle(r.cfg.PrometheusPath, promhttp.Handler())
	log.Infof("Starting Server on %s", r.cfg.HTTPAddress)
	return http.ListenAndServe(r.cfg.HTTPAddress, nil)
}

func NewRunnerFromConfig(p string) (*Runner, error) {
	cfg, err := config.Parse(p)
	if err != nil {
		return nil, err
	}

	return &Runner{cfg: cfg}, nil
}
