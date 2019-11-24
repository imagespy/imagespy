package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/imagespy/imagespy/discovery"
	"github.com/mitchellh/hashstructure"
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type GoCache struct {
	cache *gocache.Cache
}

func (g *GoCache) Get(i *discovery.Image) (*Result, error) {
	hash, err := hashImage(i)
	if err != nil {
		return nil, fmt.Errorf("hash image %s for to read from go-cache: %w", i, err)
	}

	log.Debugf("getting result cache for %s and hash %s", i, hash)
	r, found := g.cache.Get(ha, sh)
	if !found {
		return nil, nil
	}

	return r.(*Result), nil
}

func (g *GoCache) Set(i *discovery.Image, r *Result) error {
	hash, err := hashImage(i)
	if err != nil {
		return fmt.Errorf("hash image %s for to write to go-cache: %w", i, err)
	}

	log.Debugf("setting result cache for %s and hash %s", i, hash)
	g.cache.SetDefault(hash, r)
	return nil
}

func NewGoCache(exp time.Duration) *GoCache {
	c := gocache.New(exp, 2*exp)
	return &GoCache{cache: c}
}

func hashImage(i *discovery.Image) (string, error) {
	hash, err := hashstructure.Hash(i, nil)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(hash, 10), nil
}
