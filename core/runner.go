package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imagespy/imagespy/core/config"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Runner struct {
	cfg config.Config
	srv *http.Server
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

	router := mux.NewRouter()
	router.Handle(r.cfg.PrometheusPath, promhttp.Handler())
	discoverHandler := &discoverHandler{s: watcher}
	router.HandleFunc("/discover", discoverHandler.discover).Methods("POST")
	if r.cfg.UIEnabled {
		uiHandler, err := NewUIHandler(finder, r.cfg.UITemplatePath)
		if err != nil {
			log.Fatalf("create ui handler: %v", err)
		}

		router.HandleFunc("/ui", uiHandler.handle).Methods("GET")
		router.Handle("/", http.RedirectHandler("/ui", http.StatusMovedPermanently)).Methods("GET")
		fs := http.FileServer(http.Dir(r.cfg.UIStaticPath))
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fs))
	}

	r.srv = &http.Server{Addr: r.cfg.HTTPAddress, Handler: router}
	log.Infof("Starting HTTP Server on %s", r.cfg.HTTPAddress)
	err = r.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (r *Runner) Stop() error {
	if r.srv == nil {
		return nil
	}

	log.Info("Shutting down HTTP Server")
	err := r.srv.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("shut down http server: %w", err)
	}

	return nil
}

func NewRunnerFromConfig(p string) (*Runner, error) {
	cfg, err := config.Parse(p)
	if err != nil {
		return nil, err
	}

	return &Runner{cfg: cfg}, nil
}
