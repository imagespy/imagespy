package discovery

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Addr        string
	MetricsPath string
	Registry    Registry
}

type Server struct {
	cfg        Config
	Discoverer Discoverer
}

func (s *Server) Start() error {
	sc := NewScraper([]Registry{s.cfg.Registry})
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

func NewServer(cfg Config, d Discoverer) *Server {
	return &Server{
		cfg:        cfg,
		Discoverer: d,
	}
}
