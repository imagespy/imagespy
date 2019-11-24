package discovery

import (
	"net/http"

	"github.com/imagespy/imagespy/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	cfg        config.Config
	Discoverer Discoverer
}

func (s *Server) Start() error {
	sc := NewScraper(s.cfg.Registries)
	ex := NewExporter(NewFinder(s.Discoverer, sc))
	prometheus.MustRegister(ex)
	metricsPath := s.cfg.MetricsPath
	if metricsPath == "" {
		metricsPath = "/metrics"
	}

	http.Handle(metricsPath, promhttp.Handler())
	addr := s.cfg.Addr
	if addr == "" {
		addr = ":8080"
	}

	return http.ListenAndServe(addr, nil)
}

func NewServer(cfg config.Config, d Discoverer) *Server {
	return &Server{
		cfg:        cfg,
		Discoverer: d,
	}
}
