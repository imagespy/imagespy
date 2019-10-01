package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Discoverer interface {
	Discover() (*Input, error)
}

type Logger interface {
	Debug(args ...interface{})
}

type NoopLogger struct{}

func (n *NoopLogger) Debug(args ...interface{}) {}

type Discovery struct {
	log Logger
}

func (d *Discovery) Log(l Logger) {
	d.log = l
}

func (d *Discovery) Run(dv Discoverer, interval time.Duration, path string) error {
	d.log.Debug("starting initial discovery run")
	err := d.run(dv, path)
	if err != nil {
		return err
	}

	d.log.Debug("finished initial discovery run")
	t := time.NewTicker(interval)
	for range t.C {
		d.log.Debug("starting scheduled discovery run")
		err := d.run(dv, path)
		if err != nil {
			return err
		}

		d.log.Debug("finished scheduled discovery run")
	}

	return nil
}

func (d *Discovery) run(dv Discoverer, path string) error {
	in, err := dv.Discover()
	if err != nil {
		return fmt.Errorf("discoverer returned: %w", err)
	}

	b, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal discoverer input json: %w", err)
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return fmt.Errorf("write discoverer input to file %s: %w", path, err)
	}

	return nil
}

var DefaultDiscovery = &Discovery{
	log: &NoopLogger{},
}

func Log(l Logger) {
	DefaultDiscovery.Log(l)
}

func Run(dv Discoverer, interval time.Duration, path string) error {
	return DefaultDiscovery.Run(dv, interval, path)
}
