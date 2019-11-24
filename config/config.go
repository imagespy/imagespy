package config

import (
	"flag"
	"fmt"
	"net/url"
)

var (
	address     string
	metricsPath string
	registries  registryList
)

func init() {
	flag.StringVar(&address, "imagespy.exporter.address", ":8080", "The address on which the Prometheus exporter listens.")
	flag.StringVar(&metricsPath, "imagespy.exporter.metricsPath", "/metrics", "The path on which the Prometheus exporter exposes metrics.")
	flag.Var(&registries, "imagespy.registry", "A connection string to a registry")
}

type registryList []Registry

func (r *registryList) String() string {
	return ""
}

func (r *registryList) Set(v string) error {
	parsed, err := url.Parse(v)
	if err != nil {
		return err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("Unknown scheme '%s' - only 'http' and 'https' are supported", parsed.Scheme)
	}

	registry := Registry{
		Address:  parsed.Host,
		Protocol: parsed.Scheme,
	}
	username := parsed.User.Username()
	password, isSet := parsed.User.Password()
	if username == "" {
		registry.Auth = NoAuth
	} else {
		if username == "token" && !isSet {
			registry.Auth = TokenAuth
		} else {
			registry.Auth = BasicAuth
			registry.BasicPassword = password
			registry.BasicUsername = username
		}
	}

	*r = append(*r, registry)
	return nil
}

type Config struct {
	Addr        string
	MetricsPath string
	Registries  []Registry
}

type AuthMethod int

const (
	NoAuth AuthMethod = iota + 1
	BasicAuth
	TokenAuth
)

type Registry struct {
	Address       string
	Auth          AuthMethod
	BasicPassword string
	BasicUsername string
	Protocol      string
}

func (r *Registry) String() string {
	return r.Address
}

func FromFlags() Config {
	c := Config{
		Addr:        address,
		MetricsPath: metricsPath,
		Registries:  []Registry{},
	}
	for _, r := range registries {
		c.Registries = append(c.Registries, r)
	}

	return c
}
