package discovery

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricFinderErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "imagepy",
		Subsystem: "finder",
		Name:      "errors_total",
		Help:      "Number of total errors of the Finder.",
	})
)

type FinderResult struct {
	Container Container
	Latest    Image
	Input     string
}

type Finder struct {
	d Discoverer
	s *Scraper
}

func (m *Finder) Find() ([]*FinderResult, error) {
	discovererResult, err := m.d.Discover()
	if err != nil {
		return nil, fmt.Errorf("read all images in find: %w", err)
	}

	results := []*FinderResult{}
	for _, c := range discovererResult.Containers {
		latestImage, err := m.s.Scrape(c.Image)
		if err != nil {
			if errors.Is(err, ErrTagNotSupported) {
				log.Debugf("image '%s' of container '%s' has a tag that it not supported", c.Image, c.Name)
				continue
			}

			if errors.Is(err, ErrUnknownRegistry) {
				log.Debugf("registry of image '%s' of container '%s' unknown", c.Image, c.Name)
				continue
			}

			log.Errorf("scrape latest image of %s: %w", c.Image, err)
			metricFinderErrorsTotal.Inc()
			continue
		}

		results = append(results, &FinderResult{
			Container: c,
			Latest:    latestImage,
			Input:     m.d.Name(),
		})
	}

	return results, nil
}

func NewFinder(d Discoverer, s *Scraper) *Finder {
	return &Finder{d: d, s: s}
}
