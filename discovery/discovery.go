package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Discoverer is implemented by a struct that is capable of discovering Docker images.
type Discoverer interface {
	Discover() (*Input, error)
}

// Logger is implemented by a struct that can receive messages.
type Logger interface {
	Debug(args ...interface{})
}

type noopLogger struct{}

func (n *noopLogger) Debug(args ...interface{}) {}

// Discovery uses a Discoverer to find Docker images.
type Discovery struct {
	client *http.Client
	log    Logger
}

// Log sets the Logger to use.
func (d *Discovery) Log(l Logger) {
	d.log = l
}

// Run executes the given Discoverer on an interval and writes the result to the given file.
func (d *Discovery) Run(dv Discoverer, interval time.Duration, address string) error {
	d.log.Debug("starting initial discovery run")
	err := d.run(dv, address)
	if err != nil {
		return err
	}

	d.log.Debug("finished initial discovery run")
	t := time.NewTicker(interval)
	for range t.C {
		d.log.Debug("starting scheduled discovery run")
		err := d.run(dv, address)
		if err != nil {
			return err
		}

		d.log.Debug("finished scheduled discovery run")
	}

	return nil
}

func (d *Discovery) run(dv Discoverer, address string) error {
	in, err := dv.Discover()
	if err != nil {
		return fmt.Errorf("discoverer returned: %w", err)
	}

	b, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal discoverer input json: %w", err)
	}

	var url string
	if strings.HasSuffix(address, "/") {
		url = address + "discover"
	} else {
		url = address + "/discover"
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("create request for '%s': %w", url, err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("send discoverer input to '%s': %w", url, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("send discoverer input to '%s': endpoint returned status code %d", url, resp.StatusCode)
	}

	return nil
}

// DefaultDiscovery is an instance of Discovery with a no-op logger.
var DefaultDiscovery = &Discovery{
	client: &http.Client{
		Timeout: 2 * time.Second,
	},
	log: &noopLogger{},
}

// Log sets the Logger to use by the default Discovery instance.
func Log(l Logger) {
	DefaultDiscovery.Log(l)
}

// Run use the default Discovery instance to execute the given Discoverer on an interval and writes the result to the given file.
func Run(dv Discoverer, interval time.Duration, path string) error {
	return DefaultDiscovery.Run(dv, interval, path)
}
