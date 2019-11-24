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
	http.Handle(s.cfg.MetricsPath, promhttp.Handler())
	return http.ListenAndServe(s.cfg.Addr, nil)
}

func NewServer(cfg Config, d Discoverer) *Server {
	return &Server{
		cfg:        cfg,
		Discoverer: d,
	}
}
