package core

import (
	"fmt"

	"github.com/imagespy/imagespy/discovery"
	log "github.com/sirupsen/logrus"
)

type Result struct {
	Current *discovery.Image
	Latest  *discovery.Image
	Input   string
}

type ResultCache interface {
	Get(*discovery.Image) (*Result, error)
	Set(*discovery.Image, *Result) error
}

type Finder struct {
	d  *Directory
	rs ResultCache
	s  *Scraper
}

func (m *Finder) Find() ([]*Result, error) {
	inputs, err := m.d.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read all images in find: %w", err)
	}

	results := []*Result{}
	for _, i := range inputs {
		for _, image := range i.Images {
			valid, err := isValidImage(image)
			if !valid {
				log.Warnf("validation of image from input %s failed: %v", i.Name, err)
				continue
			}

			r, err := m.rs.Get(image)
			if err != nil {
				return nil, fmt.Errorf("read result of image %s from cache: %w", image, err)
			}

			if r == nil {
				log.Debugf("scraping latest image of %s", image)
				latestImage, err := m.s.Scrape(image)
				if err != nil {
					if err == ErrTagNotSupported {
						continue
					}

					return nil, fmt.Errorf("scrape latest image of %s: %w", image, err)
				}

				r = &Result{
					Current: image,
					Latest:  latestImage,
					Input:   i.Name,
				}
				err = m.rs.Set(image, r)
				if err != nil {
					return nil, fmt.Errorf("set result cache of image %s: %w", image, err)
				}
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func isValidImage(i *discovery.Image) (bool, error) {
	if i.Digest == "" {
		return false, fmt.Errorf("digest of image '%s' empty", i)
	}

	if i.Repository == "" {
		return false, fmt.Errorf("repository of image '%s' empty", i)
	}

	if i.Source == "" {
		return false, fmt.Errorf("source of image '%s' empty", i)
	}

	if i.Tag == "" {
		return false, fmt.Errorf("tag of image '%s' empty", i)
	}

	return true, nil
}

func NewFinder(d *Directory, s *Scraper, rs ResultCache) *Finder {
	return &Finder{d: d, s: s, rs: rs}
}
