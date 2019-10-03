package core

import (
	prom "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	namespace = "imagespy"
)

type Collector struct {
	f     *Finder
	image *prom.Desc
	up    *prom.Desc
}

func (c *Collector) Describe(ch chan<- *prom.Desc) {
	ch <- c.image
	ch <- c.up
}

func (c *Collector) Collect(ch chan<- prom.Metric) {
	rows, err := c.f.Find()
	if err != nil {
		log.Errorf("find images: %v", err)
		ch <- prom.MustNewConstMetric(c.up, prom.GaugeValue, 0)
		return
	}

	for _, r := range rows {
		var result float64
		if r.Latest.Digest != r.Current.Digest {
			result = 1
		}

		currentDigest := r.Current.Digest[7:15]
		latestDigest := r.Latest.Digest[7:15]
		ch <- prom.MustNewConstMetric(c.image, prom.GaugeValue, result, r.Input, r.Instance, r.Current.Source, r.Current.Repository, r.Current.Tag, r.Latest.Tag, currentDigest, latestDigest)
	}

	ch <- prom.MustNewConstMetric(c.up, prom.GaugeValue, 1)
}

func NewCollector(f *Finder) *Collector {
	return &Collector{
		f: f,
		image: prom.NewDesc(
			prom.BuildFQName(namespace, "", "image_status"),
			"Update status of an image (0=no-update, 1=needs-update).",
			[]string{"input", "instance", "source", "repository", "current_tag", "latest_tag", "current_digest", "latest_digest"},
			nil,
		),
		up: prom.NewDesc(
			prom.BuildFQName(namespace, "", "up"),
			"Could images read from the database.",
			nil,
			nil,
		),
	}
}
