package discovery

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Discoverer is implemented by a struct that is capable of discovering Docker images.
type Discoverer interface {
	Discover() (*Result, error)
	Name() string
}

// Exporter uses a Discoverer to find Docker images.
type Exporter struct {
	f         *Finder
	container *prometheus.Desc
	up        *prometheus.Desc
}

func (d *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- d.container
	ch <- d.up
}

func (d *Exporter) Collect(ch chan<- prometheus.Metric) {
	in, err := d.f.Find()
	if err != nil {
		log.Errorf("finder returned: %w", err)
		ch <- prometheus.MustNewConstMetric(d.up, prometheus.GaugeValue, 0)
		return
	}

	for _, r := range in {
		var needsUpdate float64 = 0
		if r.Latest.Digest != r.Container.Image.Digest {
			needsUpdate = 1
		}

		ch <- prometheus.MustNewConstMetric(d.container, prometheus.GaugeValue, needsUpdate, r.Container.Name, r.Container.Image.Tag, r.Latest.Tag)
	}

	ch <- prometheus.MustNewConstMetric(d.up, prometheus.GaugeValue, 1)
}

func NewExporter(f *Finder) *Exporter {
	return &Exporter{
		f: f,
		container: prometheus.NewDesc(
			prometheus.BuildFQName("imagespy", "", "container"),
			"Indicates that the image of the container needs an update (0=no update, 1=update).",
			[]string{"name", "current_tag", "latest_tag"},
			nil,
		),
		up: prometheus.NewDesc(
			prometheus.BuildFQName("imagespy", "", "up"),
			"Was the last scrape successful.",
			nil,
			nil,
		),
	}
}
