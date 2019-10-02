package main

import (
	"flag"
	"os"
	"os/signal"

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

	go func() {
		err := r.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	err = r.Stop()
	if err != nil {
		log.Fatal(err)
	}
}
