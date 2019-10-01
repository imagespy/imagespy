package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/imagespy/imagespy/core"
)

var (
	configPath = flag.String("config.path", "", "Path to the config file to load")
)

func main() {
	flag.Parse()
	r, err := core.NewRunnerFromConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(r.Run())
}
