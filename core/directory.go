package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imagespy/imagespy/discovery"
	log "github.com/sirupsen/logrus"
)

type Directory struct {
	d string
}

func (d *Directory) ReadAll() ([]*discovery.Input, error) {
	return readAll(d.d)
}

func readAll(d string) ([]*discovery.Input, error) {
	matches, err := filepath.Glob(d + "/*.json")
	if err != nil {
		return nil, err
	}

	var result []*discovery.Input
	for _, m := range matches {
		i, err := readJSON(m)
		if err != nil {
			return nil, err
		}

		result = append(result, i)
	}

	return result, nil
}

func readJSON(p string) (*discovery.Input, error) {
	log.Debugf("reading input from %s", p)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	i := &discovery.Input{}
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func NewDirectory(d string) (*Directory, error) {
	fi, err := os.Stat(d)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory to watch '%s' does not exist", d)
		}

		return nil, fmt.Errorf("create new directory: %w", err)
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("directory to watch '%s' is not a directory", d)
	}

	return &Directory{d: d}, nil
}
