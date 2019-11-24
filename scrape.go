package discovery

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/imagespy/api/versionparser"
	"github.com/imagespy/imagespy/config"
	"github.com/imagespy/registry-client"
)

var (
	ErrTagNotSupported = errors.New("cannot parse a version from tag")
	ErrUnknownRegistry = errors.New("registry unknown")
)

type Scraper struct {
	client     *http.Client
	registries []config.Registry
}

func (s *Scraper) Scrape(i Image) (li Image, err error) {
	currentVP := versionparser.FindForVersion(i.Tag)
	if isUnknownVersion(currentVP) {
		return li, ErrTagNotSupported
	}

	reg, err := s.findRegforRepo(i.Repository)
	if err != nil {
		return li, err
	}

	repo, err := reg.RepositoryFromString(i.Repository)
	if err != nil {
		return li, fmt.Errorf("create repository from string: %w", err)
	}

	lastestVP := currentVP
	latestTag := i.Tag
	tags, err := repo.Tags().GetAll()
	if err != nil {
		return li, fmt.Errorf("get all tags for %s:%s: %v", i.Repository, i.Tag, err)
	}

	for _, t := range tags {
		vp := versionparser.FindForVersion(t)
		if vp.Distinction() != lastestVP.Distinction() {
			continue
		}

		greater, err := vp.IsGreaterThan(lastestVP)
		if err != nil {
			if err == versionparser.ErrWrongVPType {
				continue
			}

			return li, err
		}

		if greater || t == i.Tag {
			lastestVP = vp
			latestTag = t
		}
	}

	ii, err := repo.Images().GetByTag(latestTag)
	if err != nil {
		return li, fmt.Errorf("get image %s:%s from registry: %v", i.Repository, i.Tag, err)
	}

	li = Image{
		Digest:     ii.Digest,
		Repository: ii.Domain + "/" + ii.Repository,
		Tag:        ii.Tag,
	}
	return li, nil
}

func (s *Scraper) findRegforRepo(repo string) (*registry.Registry, error) {
	for _, r := range s.registries {
		if strings.HasPrefix(repo, r.Address) {
			var auth registry.Authenticator
			switch r.Auth {
			case config.BasicAuth:
				auth = registry.NewBasicAuthenticator(r.BasicUsername, r.BasicPassword)
			case config.TokenAuth:
				auth = registry.NewTokenAuthenticator()
			default:
				auth = registry.NewNullAuthenticator()
			}

			return registry.New(registry.Options{
				Authenticator: auth,
				Client:        s.client,
				Domain:        r.Address,
				Protocol:      r.Protocol,
			}), nil
		}
	}

	return nil, ErrUnknownRegistry
}

func isUnknownVersion(vp versionparser.VersionParser) bool {
	_, is := vp.(*versionparser.Unknown)
	return is
}

func NewScraper(r []config.Registry) *Scraper {
	return &Scraper{
		client:     registry.DefaultClient(),
		registries: r,
	}
}
